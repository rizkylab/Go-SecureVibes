package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/models"
)

// ExportScan exports scan results in various formats
func ExportScan(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scanID := c.Params("id")
		format := c.Query("format", "json")

		// Get scan
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

		if summaryJSON != "" {
			var summary models.Summary
			json.Unmarshal([]byte(summaryJSON), &summary)
			scan.Summary = &summary
		}

		// Get findings
		findings := getFindingsForScan(db, scanID)

		switch format {
		case "json":
			return exportJSON(c, scan, findings)
		case "markdown":
			return exportMarkdown(c, scan, findings, cfg, scanID)
		case "sarif":
			return exportSARIF(c, scan, findings)
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Unsupported format. Use: json, markdown, or sarif",
			})
		}
	}
}

// exportJSON exports scan as JSON
func exportJSON(c *fiber.Ctx, scan models.Scan, findings []models.Finding) error {
	data := fiber.Map{
		"scan":     scan,
		"findings": findings,
	}

	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=scan_%s.json", scan.ID))

	return c.JSON(data)
}

// exportMarkdown exports scan as Markdown
func exportMarkdown(c *fiber.Ctx, scan models.Scan, findings []models.Finding, cfg *config.Config, scanID string) error {
	md := fmt.Sprintf("# Security Scan Report\n\n")
	md += fmt.Sprintf("**Scan ID:** %s\n", scan.ID)
	md += fmt.Sprintf("**Project:** %s\n", scan.ProjectPath)
	md += fmt.Sprintf("**Timestamp:** %s\n", scan.Timestamp.Format("2006-01-02 15:04:05"))
	if scan.CommitHash != "" {
		md += fmt.Sprintf("**Commit:** %s\n", scan.CommitHash)
	}
	if scan.Branch != "" {
		md += fmt.Sprintf("**Branch:** %s\n", scan.Branch)
	}
	md += "\n---\n\n"

	// Summary
	if scan.Summary != nil {
		md += "## Summary\n\n"
		md += fmt.Sprintf("- **Total Issues:** %d\n", scan.Summary.TotalIssues)
		md += fmt.Sprintf("- **Critical:** %d\n", scan.Summary.Critical)
		md += fmt.Sprintf("- **High:** %d\n", scan.Summary.High)
		md += fmt.Sprintf("- **Medium:** %d\n", scan.Summary.Medium)
		md += fmt.Sprintf("- **Low:** %d\n", scan.Summary.Low)
		md += fmt.Sprintf("- **Info:** %d\n", scan.Summary.Info)
		md += "\n---\n\n"
	}

	// Findings by severity
	severities := []string{"critical", "high", "medium", "low", "info"}
	for _, severity := range severities {
		severityFindings := []models.Finding{}
		for _, f := range findings {
			if f.Severity == severity {
				severityFindings = append(severityFindings, f)
			}
		}

		if len(severityFindings) > 0 {
			md += fmt.Sprintf("## %s Severity Findings\n\n", capitalize(severity))
			for i, f := range severityFindings {
				md += fmt.Sprintf("### %d. %s\n\n", i+1, f.Title)
				md += fmt.Sprintf("**Severity:** %s\n", f.Severity)
				if f.CWE != "" {
					md += fmt.Sprintf("**CWE:** %s\n", f.CWE)
				}
				if f.Category != "" {
					md += fmt.Sprintf("**Category:** %s\n", f.Category)
				}
				md += fmt.Sprintf("**File:** %s:%d\n\n", f.FilePath, f.LineNumber)
				if f.Description != "" {
					md += fmt.Sprintf("**Description:**\n%s\n\n", f.Description)
				}
				if f.Remediation != "" {
					md += fmt.Sprintf("**Remediation:**\n%s\n\n", f.Remediation)
				}
				md += "---\n\n"
			}
		}
	}

	c.Set("Content-Type", "text/markdown")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=scan_%s.md", scan.ID))

	return c.SendString(md)
}

// exportSARIF exports scan as SARIF format
func exportSARIF(c *fiber.Ctx, scan models.Scan, findings []models.Finding) error {
	// SARIF 2.1.0 format
	sarif := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":            "SecureVibes",
						"informationUri":  "https://github.com/rizkylab/Go-SecureVibes",
						"version":         "1.0.0",
						"semanticVersion": "1.0.0",
					},
				},
				"results": convertFindingsToSARIF(findings),
			},
		},
	}

	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=scan_%s.sarif", scan.ID))

	return c.JSON(sarif)
}

// convertFindingsToSARIF converts findings to SARIF results
func convertFindingsToSARIF(findings []models.Finding) []map[string]interface{} {
	results := []map[string]interface{}{}

	for _, f := range findings {
		level := "warning"
		switch f.Severity {
		case "critical", "high":
			level = "error"
		case "medium":
			level = "warning"
		case "low", "info":
			level = "note"
		}

		result := map[string]interface{}{
			"ruleId":  f.CWE,
			"level":   level,
			"message": map[string]interface{}{"text": f.Title},
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]interface{}{
							"uri": f.FilePath,
						},
						"region": map[string]interface{}{
							"startLine": f.LineNumber,
						},
					},
				},
			},
		}

		if f.Description != "" {
			result["message"] = map[string]interface{}{
				"text":     f.Title,
				"markdown": f.Description,
			}
		}

		results = append(results, result)
	}

	return results
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
