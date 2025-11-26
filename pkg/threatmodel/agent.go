package threatmodel

import (
	"fmt"
	"strings"

	"github.com/yourusername/gosecvibes/pkg/assessment"
)

// Agent handles the threat modeling phase
type Agent struct {
}

// New creates a new Threat Modeling Agent
func New() *Agent {
	return &Agent{}
}

// Run executes threat modeling based on assessment results
func (a *Agent) Run(assessment *assessment.Result) (*Result, error) {
	result := &Result{
		Threats: []Threat{},
	}

	// Analyze Endpoints
	for _, endpoint := range assessment.Endpoints {
		// Rule 1: Public Endpoints are subject to DoS
		result.Threats = append(result.Threats, Threat{
			ID:          fmt.Sprintf("T-DOS-%s", hash(endpoint.Path)),
			Category:    StrideDenialOfService,
			Title:       fmt.Sprintf("Potential Denial of Service on %s", endpoint.Path),
			Description: "Publicly accessible endpoints can be flooded with requests, causing resource exhaustion.",
			Target:      endpoint.Path,
			Severity:    SeverityMedium,
			Mitigation:  "Implement rate limiting and request validation.",
		})

		// Rule 2: Data Input endpoints (POST/PUT) are subject to Tampering
		if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
			result.Threats = append(result.Threats, Threat{
				ID:          fmt.Sprintf("T-TMP-%s", hash(endpoint.Path)),
				Category:    StrideTampering,
				Title:       fmt.Sprintf("Potential Data Tampering on %s", endpoint.Path),
				Description: "Endpoints accepting data are vulnerable to injection attacks and malformed input.",
				Target:      endpoint.Path,
				Severity:    SeverityHigh,
				Mitigation:  "Ensure strict input validation and sanitization. Use parameterized queries if interacting with DB.",
			})
		}

		// Rule 3: Information Disclosure on GET
		if endpoint.Method == "GET" {
			result.Threats = append(result.Threats, Threat{
				ID:          fmt.Sprintf("T-INF-%s", hash(endpoint.Path)),
				Category:    StrideInformationDisclosure,
				Title:       fmt.Sprintf("Potential Information Disclosure on %s", endpoint.Path),
				Description: "Ensure this endpoint does not leak sensitive data (PII, internal IDs, stack traces).",
				Target:      endpoint.Path,
				Severity:    SeverityMedium,
				Mitigation:  "Implement proper access controls and output filtering.",
			})
		}
	}

	// Analyze Dependencies
	for _, dep := range assessment.Dependencies {
		if strings.Contains(dep, "sql") || strings.Contains(dep, "pgx") || strings.Contains(dep, "gorm") || strings.Contains(dep, "mongo") {
			result.Threats = append(result.Threats, Threat{
				ID:          fmt.Sprintf("T-INJ-%s", hash(dep)),
				Category:    StrideTampering,
				Title:       fmt.Sprintf("SQL/NoSQL Injection Risk via %s", dep),
				Description: "Usage of database drivers implies data storage. Improper query construction can lead to injection.",
				Target:      dep,
				Severity:    SeverityCritical,
				Mitigation:  "Use ORM features correctly or parameterized queries. Avoid string concatenation in queries.",
			})
		}
	}

	return result, nil
}

func hash(s string) string {
	// Simple hash for ID generation (not crypto secure, just for ID)
	sum := 0
	for _, c := range s {
		sum += int(c)
	}
	return fmt.Sprintf("%d", sum)
}
