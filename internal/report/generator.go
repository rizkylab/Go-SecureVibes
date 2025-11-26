package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rizkylab/Go-SecureVibes/internal/agents/architecture"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/dast"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/staticanalysis"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/threatmodel"
)

// Generator handles report generation
type Generator struct {
	OutputFile   string
	OutputFormat string
}

// New creates a new Report Generator
func New(outputFile, outputFormat string) *Generator {
	return &Generator{
		OutputFile:   outputFile,
		OutputFormat: outputFormat,
	}
}

// Generate creates the final report
func (g *Generator) Generate(
	assess *architecture.Result,
	threats *threatmodel.Result,
	review *staticanalysis.Result,
	dastRes *dast.Result,
) error {
	if g.OutputFormat == "json" || g.OutputFormat == "both" {
		if err := g.generateJSON(assess, threats, review, dastRes); err != nil {
			return err
		}
	}

	if g.OutputFormat == "markdown" || g.OutputFormat == "both" {
		if err := g.generateMarkdown(assess, threats, review, dastRes); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateJSON(assess *architecture.Result, threats *threatmodel.Result, review *staticanalysis.Result, dastRes *dast.Result) error {
	report := map[string]interface{}{
		"scan_date":       time.Now(),
		"assessment":      assess,
		"threat_modeling": threats,
		"code_review":     review,
		"dast":            dastRes,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	filename := g.OutputFile
	if !isJSON(filename) {
		filename = filename + ".json"
	}

	return os.WriteFile(filename, data, 0644)
}

func (g *Generator) generateMarkdown(assess *architecture.Result, threats *threatmodel.Result, review *staticanalysis.Result, dastRes *dast.Result) error {
	f, err := os.Create(g.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Header
	fmt.Fprintf(f, "# 🔒 Security Audit Report\n\n")
	fmt.Fprintf(f, "**Date:** %s\n\n", time.Now().Format(time.RFC1123))

	// Executive Summary
	totalThreats := 0
	if threats != nil {
		totalThreats = len(threats.Threats)
	}
	totalVulns := 0
	if review != nil {
		totalVulns = len(review.Vulnerabilities)
	}
	totalDast := 0
	if dastRes != nil {
		totalDast = len(dastRes.Findings)
	}

	fmt.Fprintf(f, "## 📊 Executive Summary\n\n")
	fmt.Fprintf(f, "| Category | Count |\n")
	fmt.Fprintf(f, "|----------|-------|\n")
	fmt.Fprintf(f, "| Architecture Components | %d |\n", len(assess.Components))
	fmt.Fprintf(f, "| Identified Threats (STRIDE) | %d |\n", totalThreats)
	fmt.Fprintf(f, "| Code Vulnerabilities | %d |\n", totalVulns)
	fmt.Fprintf(f, "| DAST Findings | %d |\n", totalDast)
	fmt.Fprintf(f, "\n---\n\n")

	// Architecture
	fmt.Fprintf(f, "## 🏗️ Architecture Assessment\n\n")
	if len(assess.Endpoints) > 0 {
		fmt.Fprintf(f, "### Detected Endpoints\n\n")
		fmt.Fprintf(f, "| Method | Path | File |\n")
		fmt.Fprintf(f, "|--------|------|------|\n")
		for _, e := range assess.Endpoints {
			fmt.Fprintf(f, "| `%s` | `%s` | `%s:%d` |\n", e.Method, e.Path, e.File, e.Line)
		}
		fmt.Fprintf(f, "\n")
	}
	if len(assess.Dependencies) > 0 {
		fmt.Fprintf(f, "### External Dependencies\n\n")
		for _, d := range assess.Dependencies {
			fmt.Fprintf(f, "- `%s`\n", d)
		}
		fmt.Fprintf(f, "\n")
	}

	// Threat Model
	if threats != nil && len(threats.Threats) > 0 {
		fmt.Fprintf(f, "## 🎯 Threat Model (STRIDE)\n\n")
		for _, t := range threats.Threats {
			fmt.Fprintf(f, "### [%s] %s\n", t.Severity, t.Title)
			fmt.Fprintf(f, "- **Category:** %s\n", t.Category)
			fmt.Fprintf(f, "- **Target:** `%s`\n", t.Target)
			fmt.Fprintf(f, "- **Description:** %s\n", t.Description)
			fmt.Fprintf(f, "- **Mitigation:** %s\n\n", t.Mitigation)
		}
	}

	// Code Review
	if review != nil && len(review.Vulnerabilities) > 0 {
		fmt.Fprintf(f, "## 🔍 Code Review Findings\n\n")
		for _, v := range review.Vulnerabilities {
			fmt.Fprintf(f, "### [%s] %s\n", v.Severity, v.Title)
			fmt.Fprintf(f, "- **File:** `%s:%d`\n", v.File, v.Line)
			fmt.Fprintf(f, "- **CWE:** %s\n", v.CWE)
			fmt.Fprintf(f, "- **Description:** %s\n", v.Description)
			fmt.Fprintf(f, "- **Match:** `%s`\n", v.Match)
			fmt.Fprintf(f, "- **Suggestion:** %s\n\n", v.Suggestion)
		}
	}

	// DAST
	if dastRes != nil && len(dastRes.Findings) > 0 {
		fmt.Fprintf(f, "## 🚀 Dynamic Analysis Findings\n\n")
		for _, d := range dastRes.Findings {
			fmt.Fprintf(f, "### [%s] %s\n", d.Severity, d.Title)
			fmt.Fprintf(f, "- **Target:** `%s`\n", d.Target)
			fmt.Fprintf(f, "- **Description:** %s\n", d.Description)
			fmt.Fprintf(f, "- **Remediation:** %s\n\n", d.Remediation)
		}
	}

	return nil
}

func isJSON(s string) bool {
	return len(s) > 5 && s[len(s)-5:] == ".json"
}
