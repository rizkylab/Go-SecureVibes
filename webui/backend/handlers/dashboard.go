package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
)

// DashboardSummary represents dashboard summary data
type DashboardSummary struct {
	TotalScans        int             `json:"total_scans"`
	RecentScans       int             `json:"recent_scans"`
	TotalFindings     int             `json:"total_findings"`
	SeverityBreakdown SeverityCount   `json:"severity_breakdown"`
	RecentFindings    []RecentFinding `json:"recent_findings"`
	TrendData         []TrendPoint    `json:"trend_data"`
}

// SeverityCount represents count by severity
type SeverityCount struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
}

// RecentFinding represents a recent finding
type RecentFinding struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Severity  string `json:"severity"`
	FilePath  string `json:"file_path"`
	ScanID    string `json:"scan_id"`
	Timestamp string `json:"timestamp"`
}

// TrendPoint represents a point in trend chart
type TrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetDashboardSummary returns dashboard summary
func GetDashboardSummary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		summary := DashboardSummary{}

		// Get total scans
		db.QueryRow("SELECT COUNT(*) FROM scans").Scan(&summary.TotalScans)

		// Get recent scans (last 7 days)
		db.QueryRow(`
			SELECT COUNT(*) FROM scans 
			WHERE timestamp >= datetime('now', '-7 days')
		`).Scan(&summary.RecentScans)

		// Get total findings
		db.QueryRow("SELECT COUNT(*) FROM findings").Scan(&summary.TotalFindings)

		// Get severity breakdown
		rows, err := db.Query(`
			SELECT severity, COUNT(*) as count 
			FROM findings 
			GROUP BY severity
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var severity string
				var count int
				if err := rows.Scan(&severity, &count); err == nil {
					switch severity {
					case "critical":
						summary.SeverityBreakdown.Critical = count
					case "high":
						summary.SeverityBreakdown.High = count
					case "medium":
						summary.SeverityBreakdown.Medium = count
					case "low":
						summary.SeverityBreakdown.Low = count
					case "info":
						summary.SeverityBreakdown.Info = count
					}
				}
			}
		}

		// Get recent findings (last 10)
		rows, err = db.Query(`
			SELECT f.id, f.title, f.severity, f.file_path, f.scan_id, s.timestamp
			FROM findings f
			JOIN scans s ON f.scan_id = s.id
			ORDER BY s.timestamp DESC
			LIMIT 10
		`)
		if err == nil {
			defer rows.Close()
			summary.RecentFindings = []RecentFinding{}
			for rows.Next() {
				var finding RecentFinding
				if err := rows.Scan(&finding.ID, &finding.Title, &finding.Severity,
					&finding.FilePath, &finding.ScanID, &finding.Timestamp); err == nil {
					summary.RecentFindings = append(summary.RecentFindings, finding)
				}
			}
		}

		// Get trend data (last 30 days)
		rows, err = db.Query(`
			SELECT DATE(timestamp) as date, COUNT(*) as count
			FROM scans
			WHERE timestamp >= datetime('now', '-30 days')
			GROUP BY DATE(timestamp)
			ORDER BY date
		`)
		if err == nil {
			defer rows.Close()
			summary.TrendData = []TrendPoint{}
			for rows.Next() {
				var point TrendPoint
				if err := rows.Scan(&point.Date, &point.Count); err == nil {
					summary.TrendData = append(summary.TrendData, point)
				}
			}
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    summary,
		})
	}
}
