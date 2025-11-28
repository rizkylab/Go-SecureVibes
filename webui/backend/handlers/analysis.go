package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/models"
)

// GetArchitecture returns architecture data for a scan
func GetArchitecture(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		// Verify scan exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM scans WHERE id = ?)", scanID).Scan(&exists)
		if err != nil || !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan not found",
			})
		}

		// Read architecture file
		archPath := filepath.Join(cfg.Storage.OutputDir, scanID, "architecture.json")
		data, err := os.ReadFile(archPath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Architecture data not found",
			})
		}

		var architecture models.Architecture
		if err := json.Unmarshal(data, &architecture); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to parse architecture data",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    architecture,
		})
	}
}

// GetThreatModel returns threat model data for a scan
func GetThreatModel(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		// Verify scan exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM scans WHERE id = ?)", scanID).Scan(&exists)
		if err != nil || !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan not found",
			})
		}

		// Read threat model file
		threatPath := filepath.Join(cfg.Storage.OutputDir, scanID, "threat_model.json")
		data, err := os.ReadFile(threatPath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Threat model data not found",
			})
		}

		var threatModel models.ThreatModel
		if err := json.Unmarshal(data, &threatModel); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to parse threat model data",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    threatModel,
		})
	}
}

// GetDASTFindings returns DAST findings for a scan
func GetDASTFindings(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		rows, err := db.Query(`
			SELECT id, scan_id, endpoint, method, title, description, severity,
			       request_payload, response_code, response_body, evidence, 
			       screenshot, remediation, timestamp
			FROM dast_findings
			WHERE scan_id = ?
			ORDER BY CASE severity 
				WHEN 'critical' THEN 1 
				WHEN 'high' THEN 2 
				WHEN 'medium' THEN 3 
				WHEN 'low' THEN 4 
				ELSE 5 
			END
		`, scanID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch DAST findings",
			})
		}
		defer rows.Close()

		findings := []models.DASTFinding{}
		for rows.Next() {
			var finding models.DASTFinding
			var screenshot *string

			err := rows.Scan(&finding.ID, &finding.ScanID, &finding.Endpoint, &finding.Method,
				&finding.Title, &finding.Description, &finding.Severity, &finding.RequestPayload,
				&finding.ResponseCode, &finding.ResponseBody, &finding.Evidence,
				&screenshot, &finding.Remediation, &finding.Timestamp)

			if err != nil {
				continue
			}

			if screenshot != nil {
				finding.Screenshot = *screenshot
			}

			findings = append(findings, finding)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    findings,
		})
	}
}
