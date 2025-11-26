package dast

// Finding represents a dynamic analysis finding
type Finding struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Target      string `json:"target"` // URL or Endpoint
	Evidence    string `json:"evidence"`
	Remediation string `json:"remediation"`
}

// Result holds the DAST findings
type Result struct {
	Findings []Finding `json:"findings"`
}
