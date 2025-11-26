package codereview

// Vulnerability represents a security issue found in the code
type Vulnerability struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Severity    string `json:"severity"`
	CWE         string `json:"cwe"`
	Suggestion  string `json:"suggestion"`
	Match       string `json:"match_content"`
}

// Result holds the code review findings
type Result struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}
