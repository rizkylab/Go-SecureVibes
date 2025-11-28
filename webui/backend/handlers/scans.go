package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/models"
)

// ListScans returns paginated list of scans
func ListScans(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse query parameters
		page, _ := strconv.Atoi(c.Query("page", "1"))
		pageSize, _ := strconv.Atoi(c.Query("size", "20"))
		status := c.Query("status", "")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		offset := (page - 1) * pageSize

		// Build query
		query := "SELECT id, timestamp, project_path, commit_hash, branch, duration, dast_enabled, summary_json, status, created_at FROM scans"
		countQuery := "SELECT COUNT(*) FROM scans"
		args := []interface{}{}

		if status != "" {
			query += " WHERE status = ?"
			countQuery += " WHERE status = ?"
			args = append(args, status)
		}

		query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		// Get total count
		var totalItems int
		countArgs := args
		db.QueryRow(countQuery, countArgs...).Scan(&totalItems)

		// Get scans
		args = append(args, pageSize, offset)
		rows, err := db.Query(query, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch scans",
			})
		}
		defer rows.Close()

		scans := []models.Scan{}
		for rows.Next() {
			var scan models.Scan
			var summaryJSON string
			err := rows.Scan(&scan.ID, &scan.Timestamp, &scan.ProjectPath, &scan.CommitHash,
				&scan.Branch, &scan.Duration, &scan.DASTEnabled, &summaryJSON, &scan.Status, &scan.CreatedAt)
			if err != nil {
				continue
			}

			// Parse summary JSON
			if summaryJSON != "" {
				var summary models.Summary
				if err := json.Unmarshal([]byte(summaryJSON), &summary); err == nil {
					scan.Summary = &summary
				}
			}

			scans = append(scans, scan)
		}

		totalPages := (totalItems + pageSize - 1) / pageSize

		return c.JSON(models.PaginatedResponse{
			Success:    true,
			Data:       scans,
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		})
	}
}

// GetScan returns a single scan by ID
func GetScan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		var scan models.Scan
		var summaryJSON string

		err := db.QueryRow(`
			SELECT id, timestamp, project_path, commit_hash, branch, duration, 
			       dast_enabled, summary_json, status, created_at, updated_at
			FROM scans WHERE id = ?
		`, scanID).Scan(&scan.ID, &scan.Timestamp, &scan.ProjectPath, &scan.CommitHash,
			&scan.Branch, &scan.Duration, &scan.DASTEnabled, &summaryJSON, &scan.Status,
			&scan.CreatedAt, &scan.UpdatedAt)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan not found",
			})
		}

		// Parse summary JSON
		if summaryJSON != "" {
			var summary models.Summary
			if err := json.Unmarshal([]byte(summaryJSON), &summary); err == nil {
				scan.Summary = &summary
			}
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    scan,
		})
	}
}

// UploadScanResult handles scan result upload
func UploadScanResult(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ScanRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			})
		}

		// Validate required fields
		if req.ProjectPath == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "project_path is required",
			})
		}

		// Generate scan ID
		scanID := uuid.New().String()
		timestamp := time.Now()

		// Serialize summary
		summaryJSON, _ := json.Marshal(req.Summary)

		// Insert scan
		_, err := db.Exec(`
			INSERT INTO scans (id, timestamp, project_path, commit_hash, branch, 
			                   duration, dast_enabled, summary_json, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, scanID, timestamp, req.ProjectPath, req.CommitHash, req.Branch,
			0, len(req.DASTFindings) > 0, string(summaryJSON), "completed")

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create scan",
			})
		}

		// Insert findings
		for _, finding := range req.Findings {
			findingID := uuid.New().String()
			_, err := db.Exec(`
				INSERT INTO findings (id, scan_id, title, description, severity, cwe, 
				                      category, file_path, line_number, line_content, 
				                      remediation, confidence, first_detected, last_seen, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, findingID, scanID, finding.Title, finding.Description, finding.Severity,
				finding.CWE, finding.Category, finding.FilePath, finding.LineNumber,
				finding.LineContent, finding.Remediation, finding.Confidence,
				timestamp, timestamp, "new")

			if err != nil {
				// Log error but continue
				fmt.Printf("Failed to insert finding: %v\n", err)
			}
		}

		// Insert DAST findings
		for _, dastFinding := range req.DASTFindings {
			findingID := uuid.New().String()
			_, err := db.Exec(`
				INSERT INTO dast_findings (id, scan_id, endpoint, method, title, description,
				                           severity, request_payload, response_code, response_body,
				                           evidence, screenshot, remediation, timestamp)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, findingID, scanID, dastFinding.Endpoint, dastFinding.Method, dastFinding.Title,
				dastFinding.Description, dastFinding.Severity, dastFinding.RequestPayload,
				dastFinding.ResponseCode, dastFinding.ResponseBody, dastFinding.Evidence,
				dastFinding.Screenshot, dastFinding.Remediation, timestamp)

			if err != nil {
				fmt.Printf("Failed to insert DAST finding: %v\n", err)
			}
		}

		// Save architecture and threat model to files
		scanDir := filepath.Join(cfg.Storage.OutputDir, scanID)
		os.MkdirAll(scanDir, 0755)

		if req.Architecture != nil {
			archData, _ := json.MarshalIndent(req.Architecture, "", "  ")
			os.WriteFile(filepath.Join(scanDir, "architecture.json"), archData, 0644)
		}

		if req.ThreatModel != nil {
			threatData, _ := json.MarshalIndent(req.ThreatModel, "", "  ")
			os.WriteFile(filepath.Join(scanDir, "threat_model.json"), threatData, 0644)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Scan uploaded successfully",
			"data": fiber.Map{
				"scan_id": scanID,
			},
		})
	}
}

// DeleteScan deletes a scan and its associated data
func DeleteScan(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		// Delete from database (cascade will delete findings)
		result, err := db.Exec("DELETE FROM scans WHERE id = ?", scanID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to delete scan",
			})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan not found",
			})
		}

		// Delete scan directory
		scanDir := filepath.Join(cfg.Storage.OutputDir, scanID)
		os.RemoveAll(scanDir)

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Scan deleted successfully",
		})
	}
}
