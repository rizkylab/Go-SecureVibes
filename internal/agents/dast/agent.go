package dast

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
)

// Agent handles the dynamic application security testing phase
type Agent struct {
	TargetURL string
}

// New creates a new DAST Agent
func New(targetURL string) *Agent {
	return &Agent{
		TargetURL: targetURL,
	}
}

// Run executes the DAST scan
func (a *Agent) Run() (*Result, error) {
	result := &Result{
		Findings: []Finding{},
	}

	if a.TargetURL == "" {
		color.Yellow("   DAST skipped: No target URL provided (use --target)")
		return result, nil
	}

	color.Cyan("   Targeting: %s", a.TargetURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 1. Basic Health Check & Header Analysis
	resp, err := client.Get(a.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target: %v", err)
	}
	defer resp.Body.Close()

	// Check Security Headers
	headers := map[string]string{
		"X-Frame-Options":           "Clickjacking protection",
		"Content-Security-Policy":   "XSS protection",
		"Strict-Transport-Security": "Man-in-the-Middle protection",
		"X-Content-Type-Options":    "MIME-sniffing protection",
	}

	for header, desc := range headers {
		if val := resp.Header.Get(header); val == "" {
			result.Findings = append(result.Findings, Finding{
				ID:          "DAST-MISSING-HEADER-" + header,
				Title:       fmt.Sprintf("Missing Security Header: %s", header),
				Description: fmt.Sprintf("The response is missing the %s header, which provides %s.", header, desc),
				Severity:    "Low",
				Target:      a.TargetURL,
				Evidence:    "Header not found in response",
				Remediation: fmt.Sprintf("Configure the server to send the %s header.", header),
			})
		}
	}

	// 2. Check for Server Information Leakage
	if server := resp.Header.Get("Server"); server != "" {
		result.Findings = append(result.Findings, Finding{
			ID:          "DAST-INFO-SERVER",
			Title:       "Server Information Disclosure",
			Description: fmt.Sprintf("The 'Server' header exposes version information: %s", server),
			Severity:    "Low",
			Target:      a.TargetURL,
			Evidence:    fmt.Sprintf("Server: %s", server),
			Remediation: "Remove or obfuscate the 'Server' header in production.",
		})
	}

	return result, nil
}
