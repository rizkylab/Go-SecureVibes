package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/models"
)

// ListFindings returns paginated list of findings for a scan
func ListFindings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")

		// Parse query parameters
		page, _ := strconv.Atoi(c.Query("page", "1"))
		pageSize, _ := strconv.Atoi(c.Query("size", "50"))
		severity := c.Query("severity", "")
		category := c.Query("category", "")
		cwe := c.Query("cwe", "")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 50
		}

		offset := (page - 1) * pageSize

		// Build query
		query := `SELECT id, scan_id, title, description, severity, cwe, category, 
		                 file_path, line_number, line_content, remediation, confidence, 
		                 first_detected, last_seen, status 
		          FROM findings WHERE scan_id = ?`
		countQuery := "SELECT COUNT(*) FROM findings WHERE scan_id = ?"
		args := []interface{}{scanID}

		if severity != "" {
			query += " AND severity = ?"
			countQuery += " AND severity = ?"
			args = append(args, severity)
		}

		if category != "" {
			query += " AND category = ?"
			countQuery += " AND category = ?"
			args = append(args, category)
		}

		if cwe != "" {
			query += " AND cwe = ?"
			countQuery += " AND cwe = ?"
			args = append(args, cwe)
		}

		query += " ORDER BY CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 ELSE 5 END, line_number LIMIT ? OFFSET ?"

		// Get total count
		var totalItems int
		countArgs := args
		db.QueryRow(countQuery, countArgs...).Scan(&totalItems)

		// Get findings
		args = append(args, pageSize, offset)
		rows, err := db.Query(query, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch findings",
			})
		}
		defer rows.Close()

		findings := []models.Finding{}
		for rows.Next() {
			var finding models.Finding
			err := rows.Scan(&finding.ID, &finding.ScanID, &finding.Title, &finding.Description,
				&finding.Severity, &finding.CWE, &finding.Category, &finding.FilePath,
				&finding.LineNumber, &finding.LineContent, &finding.Remediation,
				&finding.Confidence, &finding.FirstDetected, &finding.LastSeen, &finding.Status)
			if err != nil {
				continue
			}
			findings = append(findings, finding)
		}

		totalPages := (totalItems + pageSize - 1) / pageSize

		return c.JSON(models.PaginatedResponse{
			Success:    true,
			Data:       findings,
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		})
	}
}

// ListAllFindings returns paginated list of findings across all scans
func ListAllFindings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse query parameters
		page, _ := strconv.Atoi(c.Query("page", "1"))
		pageSize, _ := strconv.Atoi(c.Query("size", "50"))
		severity := c.Query("severity", "")
		category := c.Query("category", "")
		cwe := c.Query("cwe", "")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 50
		}

		offset := (page - 1) * pageSize

		// Build query
		query := `SELECT f.id, f.scan_id, f.title, f.description, f.severity, f.cwe, f.category, 
		                 f.file_path, f.line_number, f.line_content, f.remediation, f.confidence, 
		                 f.first_detected, f.last_seen, f.status, s.project_path, s.branch
		          FROM findings f
		          JOIN scans s ON f.scan_id = s.id
		          WHERE 1=1`
		countQuery := "SELECT COUNT(*) FROM findings f WHERE 1=1"
		args := []interface{}{}

		if severity != "" {
			query += " AND f.severity = ?"
			countQuery += " AND f.severity = ?"
			args = append(args, severity)
		}

		if category != "" {
			query += " AND f.category = ?"
			countQuery += " AND f.category = ?"
			args = append(args, category)
		}

		if cwe != "" {
			query += " AND f.cwe = ?"
			countQuery += " AND f.cwe = ?"
			args = append(args, cwe)
		}

		query += " ORDER BY f.last_seen DESC LIMIT ? OFFSET ?"

		// Get total count
		var totalItems int
		countArgs := args
		db.QueryRow(countQuery, countArgs...).Scan(&totalItems)

		// Get findings
		args = append(args, pageSize, offset)
		rows, err := db.Query(query, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch findings",
			})
		}
		defer rows.Close()

		type FindingWithContext struct {
			models.Finding
			ProjectPath string `json:"project_path"`
			Branch      string `json:"branch"`
		}

		findings := []FindingWithContext{}
		for rows.Next() {
			var f FindingWithContext
			err := rows.Scan(&f.ID, &f.ScanID, &f.Title, &f.Description,
				&f.Severity, &f.CWE, &f.Category, &f.FilePath,
				&f.LineNumber, &f.LineContent, &f.Remediation,
				&f.Confidence, &f.FirstDetected, &f.LastSeen, &f.Status,
				&f.ProjectPath, &f.Branch)
			if err != nil {
				continue
			}
			findings = append(findings, f)
		}

		totalPages := (totalItems + pageSize - 1) / pageSize

		return c.JSON(models.PaginatedResponse{
			Success:    true,
			Data:       findings,
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		})
	}
}

// GetFinding returns a single finding by ID
func GetFinding(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")
		findingID := c.Params("finding_id")

		var finding models.Finding

		err := db.QueryRow(`
			SELECT id, scan_id, title, description, severity, cwe, category, 
			       file_path, line_number, line_content, remediation, confidence, 
			       first_detected, last_seen, status
			FROM findings 
			WHERE id = ? AND scan_id = ?
		`, findingID, scanID).Scan(&finding.ID, &finding.ScanID, &finding.Title,
			&finding.Description, &finding.Severity, &finding.CWE, &finding.Category,
			&finding.FilePath, &finding.LineNumber, &finding.LineContent,
			&finding.Remediation, &finding.Confidence, &finding.FirstDetected,
			&finding.LastSeen, &finding.Status)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Finding not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    finding,
		})
	}
}

// UpdateFindingStatus updates the status of a finding
func UpdateFindingStatus(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")
		findingID := c.Params("finding_id")

		var req struct {
			Status string `json:"status"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			})
		}

		// Validate status
		validStatuses := map[string]bool{
			"new":            true,
			"existing":       true,
			"fixed":          true,
			"false_positive": true,
		}

		if !validStatuses[req.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid status value",
			})
		}

		// Update status
		result, err := db.Exec(`
			UPDATE findings 
			SET status = ? 
			WHERE id = ? AND scan_id = ?
		`, req.Status, findingID, scanID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to update finding status",
			})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Finding not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Finding status updated successfully",
		})
	}
}
