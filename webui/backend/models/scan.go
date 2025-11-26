package models

import (
	"time"
)

// Scan represents a security scan session
type Scan struct {
	ID          string    `json:"id" db:"id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	ProjectPath string    `json:"project_path" db:"project_path"`
	CommitHash  string    `json:"commit_hash,omitempty" db:"commit_hash"`
	Branch      string    `json:"branch,omitempty" db:"branch"`
	Duration    int64     `json:"duration" db:"duration"` // milliseconds
	DASTEnabled bool      `json:"dast_enabled" db:"dast_enabled"`
	SummaryJSON string    `json:"-" db:"summary_json"`
	Summary     *Summary  `json:"summary"`
	Status      string    `json:"status" db:"status"` // completed, failed, running
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Summary represents the vulnerability summary
type Summary struct {
	TotalIssues int `json:"total_issues"`
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
	Info        int `json:"info"`
}

// Architecture represents the application architecture
type Architecture struct {
	Components []Component `json:"components"`
	DataFlows  []DataFlow  `json:"data_flows"`
}

// Component represents a system component
type Component struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"` // handler, service, database, external
	Description  string   `json:"description"`
	FilePath     string   `json:"file_path"`
	LineNumber   int      `json:"line_number"`
	Dependencies []string `json:"dependencies"`
}

// DataFlow represents data flow between components
type DataFlow struct {
	From      string `json:"from"`
	To        string `json:"to"`
	DataType  string `json:"data_type"`
	Protocol  string `json:"protocol"`
	Encrypted bool   `json:"encrypted"`
}

// ThreatModel represents STRIDE threat model
type ThreatModel struct {
	Threats []Threat `json:"threats"`
}

// Threat represents a single threat
type Threat struct {
	ID          string `json:"id"`
	Component   string `json:"component"`
	Type        string `json:"type"` // Spoofing, Tampering, Repudiation, etc.
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Mitigation  string `json:"mitigation"`
	Status      string `json:"status"` // open, mitigated, accepted
}

// Finding represents a static analysis finding
type Finding struct {
	ID            string    `json:"id"`
	ScanID        string    `json:"scan_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Severity      string    `json:"severity"`
	CWE           string    `json:"cwe"`
	Category      string    `json:"category"` // crypto, auth, injection, config, secrets
	FilePath      string    `json:"file_path"`
	LineNumber    int       `json:"line_number"`
	LineContent   string    `json:"line_content"`
	Remediation   string    `json:"remediation"`
	References    []string  `json:"references"`
	Confidence    string    `json:"confidence"` // high, medium, low
	FirstDetected time.Time `json:"first_detected"`
	LastSeen      time.Time `json:"last_seen"`
	Status        string    `json:"status"` // new, existing, fixed
}

// DASTFinding represents a dynamic analysis finding
type DASTFinding struct {
	ID             string            `json:"id"`
	ScanID         string            `json:"scan_id"`
	Endpoint       string            `json:"endpoint"`
	Method         string            `json:"method"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Severity       string            `json:"severity"`
	RequestPayload string            `json:"request_payload"`
	ResponseCode   int               `json:"response_code"`
	ResponseBody   string            `json:"response_body"`
	Evidence       string            `json:"evidence"`
	Screenshot     string            `json:"screenshot,omitempty"`
	Headers        map[string]string `json:"headers"`
	Remediation    string            `json:"remediation"`
	Timestamp      time.Time         `json:"timestamp"`
}

// ScanComparison represents comparison between two scans
type ScanComparison struct {
	ScanA            *Scan          `json:"scan_a"`
	ScanB            *Scan          `json:"scan_b"`
	NewFindings      []Finding      `json:"new_findings"`
	FixedFindings    []Finding      `json:"fixed_findings"`
	ExistingFindings []Finding      `json:"existing_findings"`
	SeverityTrend    *SeverityTrend `json:"severity_trend"`
}

// SeverityTrend represents severity trend between scans
type SeverityTrend struct {
	CriticalDelta int `json:"critical_delta"`
	HighDelta     int `json:"high_delta"`
	MediumDelta   int `json:"medium_delta"`
	LowDelta      int `json:"low_delta"`
	TotalDelta    int `json:"total_delta"`
}

// ScanRequest represents a request to create/upload a scan
type ScanRequest struct {
	ProjectPath  string                 `json:"project_path"`
	CommitHash   string                 `json:"commit_hash,omitempty"`
	Branch       string                 `json:"branch,omitempty"`
	Architecture *Architecture          `json:"architecture"`
	ThreatModel  *ThreatModel           `json:"threat_model"`
	Findings     []Finding              `json:"findings"`
	DASTFindings []DASTFinding          `json:"dast_findings,omitempty"`
	Summary      *Summary               `json:"summary"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// FilterOptions represents filter options for findings
type FilterOptions struct {
	Severity   []string `json:"severity"`
	CWE        []string `json:"cwe"`
	Category   []string `json:"category"`
	Status     []string `json:"status"`
	FilePath   string   `json:"file_path"`
	SearchTerm string   `json:"search_term"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}
