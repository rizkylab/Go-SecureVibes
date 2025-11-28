package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/models"
)

// CompareScans compares two scans and returns the differences
func CompareScans(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanAID := c.Query("scan_a")
		scanBID := c.Query("scan_b")

		if scanAID == "" || scanBID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Both scan_a and scan_b parameters are required",
			})
		}

		// Get scan A
		var scanA models.Scan
		var summaryAJSON string
		err := db.QueryRow(`
			SELECT id, timestamp, project_path, commit_hash, branch, duration, 
			       dast_enabled, summary_json, status, created_at, updated_at
			FROM scans WHERE id = ?
		`, scanAID).Scan(&scanA.ID, &scanA.Timestamp, &scanA.ProjectPath, &scanA.CommitHash,
			&scanA.Branch, &scanA.Duration, &scanA.DASTEnabled, &summaryAJSON, &scanA.Status,
			&scanA.CreatedAt, &scanA.UpdatedAt)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan A not found",
			})
		}

		if summaryAJSON != "" {
			var summary models.Summary
			json.Unmarshal([]byte(summaryAJSON), &summary)
			scanA.Summary = &summary
		}

		// Get scan B
		var scanB models.Scan
		var summaryBJSON string
		err = db.QueryRow(`
			SELECT id, timestamp, project_path, commit_hash, branch, duration, 
			       dast_enabled, summary_json, status, created_at, updated_at
			FROM scans WHERE id = ?
		`, scanBID).Scan(&scanB.ID, &scanB.Timestamp, &scanB.ProjectPath, &scanB.CommitHash,
			&scanB.Branch, &scanB.Duration, &scanB.DASTEnabled, &summaryBJSON, &scanB.Status,
			&scanB.CreatedAt, &scanB.UpdatedAt)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Scan B not found",
			})
		}

		if summaryBJSON != "" {
			var summary models.Summary
			json.Unmarshal([]byte(summaryBJSON), &summary)
			scanB.Summary = &summary
		}

		// Get findings for scan A
		findingsA := getFindingsForScan(db, scanAID)
		findingsAMap := make(map[string]models.Finding)
		for _, f := range findingsA {
			key := f.FilePath + ":" + f.Title
			findingsAMap[key] = f
		}

		// Get findings for scan B
		findingsB := getFindingsForScan(db, scanBID)
		findingsBMap := make(map[string]models.Finding)
		for _, f := range findingsB {
			key := f.FilePath + ":" + f.Title
			findingsBMap[key] = f
		}

		// Calculate differences
		var newFindings []models.Finding
		var fixedFindings []models.Finding
		var existingFindings []models.Finding

		// Find new findings (in B but not in A)
		for key, finding := range findingsBMap {
			if _, exists := findingsAMap[key]; !exists {
				newFindings = append(newFindings, finding)
			} else {
				existingFindings = append(existingFindings, finding)
			}
		}

		// Find fixed findings (in A but not in B)
		for key, finding := range findingsAMap {
			if _, exists := findingsBMap[key]; !exists {
				fixedFindings = append(fixedFindings, finding)
			}
		}

		// Calculate severity trend
		severityTrend := &models.SeverityTrend{}
		if scanA.Summary != nil && scanB.Summary != nil {
			severityTrend.CriticalDelta = scanB.Summary.Critical - scanA.Summary.Critical
			severityTrend.HighDelta = scanB.Summary.High - scanA.Summary.High
			severityTrend.MediumDelta = scanB.Summary.Medium - scanA.Summary.Medium
			severityTrend.LowDelta = scanB.Summary.Low - scanA.Summary.Low
			severityTrend.TotalDelta = scanB.Summary.TotalIssues - scanA.Summary.TotalIssues
		}

		comparison := models.ScanComparison{
			ScanA:            &scanA,
			ScanB:            &scanB,
			NewFindings:      newFindings,
			FixedFindings:    fixedFindings,
			ExistingFindings: existingFindings,
			SeverityTrend:    severityTrend,
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    comparison,
		})
	}
}

// getFindingsForScan retrieves all findings for a scan
func getFindingsForScan(db *database.DB, scanID string) []models.Finding {
	rows, err := db.Query(`
		SELECT id, scan_id, title, description, severity, cwe, category, 
		       file_path, line_number, line_content, remediation, confidence, 
		       first_detected, last_seen, status
		FROM findings WHERE scan_id = ?
	`, scanID)

	if err != nil {
		return []models.Finding{}
	}
	defer rows.Close()

	findings := []models.Finding{}
	for rows.Next() {
		var finding models.Finding
		rows.Scan(&finding.ID, &finding.ScanID, &finding.Title, &finding.Description,
			&finding.Severity, &finding.CWE, &finding.Category, &finding.FilePath,
			&finding.LineNumber, &finding.LineContent, &finding.Remediation,
			&finding.Confidence, &finding.FirstDetected, &finding.LastSeen, &finding.Status)
		findings = append(findings, finding)
	}

	return findings
}
