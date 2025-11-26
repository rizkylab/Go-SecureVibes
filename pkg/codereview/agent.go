package codereview

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Rule defines a pattern to search for
type Rule struct {
	ID          string
	Pattern     *regexp.Regexp
	Title       string
	Description string
	Severity    string
	CWE         string
	Suggestion  string
}

// Agent handles the static code analysis phase
type Agent struct {
	ProjectPath string
	Excludes    []string
	Rules       []Rule
}

// New creates a new Code Review Agent
func New(projectPath string, excludes []string) *Agent {
	agent := &Agent{
		ProjectPath: projectPath,
		Excludes:    excludes,
	}
	agent.loadRules()
	return agent
}

func (a *Agent) loadRules() {
	a.Rules = []Rule{
		{
			ID:          "SA-001",
			Pattern:     regexp.MustCompile(`(?i)(password|secret|api[_]?key|access[_]?token)\s*[:=]+\s*["'][^"']+["']`),
			Title:       "Hardcoded Credential",
			Description: "Potential hardcoded credential detected.",
			Severity:    "High",
			CWE:         "CWE-798",
			Suggestion:  "Use environment variables or a secrets manager.",
		},
		{
			ID:          "SA-002",
			Pattern:     regexp.MustCompile(`(?i)md5\.New\(\)|sha1\.New\(\)`),
			Title:       "Weak Cryptography",
			Description: "Usage of weak hashing algorithm (MD5/SHA1).",
			Severity:    "Medium",
			CWE:         "CWE-327",
			Suggestion:  "Use stronger algorithms like SHA-256 or bcrypt/argon2 for passwords.",
		},
		{
			ID:          "SA-003",
			Pattern:     regexp.MustCompile(`(?i)(Select|Insert|Update|Delete).*\+.*`),
			Title:       "Potential SQL Injection",
			Description: "Dynamic SQL query construction using string concatenation.",
			Severity:    "High",
			CWE:         "CWE-89",
			Suggestion:  "Use parameterized queries.",
		},
		{
			ID:          "SA-004",
			Pattern:     regexp.MustCompile(`(?i)fmt\.Printf\(.*%v.*,.*(password|secret).*`),
			Title:       "Sensitive Data Leak in Logs",
			Description: "Logging sensitive data.",
			Severity:    "Medium",
			CWE:         "CWE-532",
			Suggestion:  "Ensure sensitive data is redacted before logging.",
		},
	}
}

// Run executes the code review
func (a *Agent) Run() (*Result, error) {
	result := &Result{
		Vulnerabilities: []Vulnerability{},
	}

	err := filepath.Walk(a.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			for _, exclude := range a.Excludes {
				if strings.Contains(path, exclude) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check supported extensions (currently mostly Go focused, but regex works on text)
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".js" && ext != ".ts" && ext != ".py" {
			return nil
		}

		return a.scanFile(path, result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Agent) scanFile(path string, result *Result) error {
	file, err := os.Open(path)
	if err != nil {
		color.Yellow("Could not open file %s: %v", path, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comments (basic heuristic)
		if strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		for _, rule := range a.Rules {
			if rule.Pattern.MatchString(line) {
				result.Vulnerabilities = append(result.Vulnerabilities, Vulnerability{
					ID:          rule.ID,
					Title:       rule.Title,
					Description: rule.Description,
					File:        path,
					Line:        lineNum,
					Severity:    rule.Severity,
					CWE:         rule.CWE,
					Suggestion:  rule.Suggestion,
					Match:       strings.TrimSpace(line),
				})
			}
		}
	}
	return nil
}
