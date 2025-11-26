package assessment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Agent handles the architecture assessment phase
type Agent struct {
	ProjectPath string
	Excludes    []string
}

// New creates a new Assessment Agent
func New(projectPath string, excludes []string) *Agent {
	return &Agent{
		ProjectPath: projectPath,
		Excludes:    excludes,
	}
}

// Run executes the assessment
func (a *Agent) Run() (*Result, error) {
	result := &Result{
		Components:   []Component{},
		Endpoints:    []Endpoint{},
		DataFlows:    []DataFlow{},
		Dependencies: []string{},
	}

	err := filepath.Walk(a.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if info.IsDir() {
			for _, exclude := range a.Excludes {
				if strings.Contains(path, exclude) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Only analyze Go files for now
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		return a.analyzeGoFile(path, result)
	})

	if err != nil {
		return nil, err
	}

	// Deduplicate dependencies
	result.Dependencies = uniqueStrings(result.Dependencies)

	return result, nil
}

func (a *Agent) analyzeGoFile(path string, result *Result) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		// Log error but continue
		color.Yellow("Failed to parse file %s: %v", path, err)
		return nil
	}

	// Analyze imports for dependencies
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		if !strings.Contains(path, ".") {
			continue // Skip standard library (heuristic)
		}
		result.Dependencies = append(result.Dependencies, path)
	}

	// Analyze functions for potential endpoints
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for function calls that resemble HTTP handlers
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check for http.HandleFunc, r.GET, r.POST, etc.
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			funcName := sel.Sel.Name

			// Basic heuristic for HTTP methods
			if isHTTPMethod(funcName) || funcName == "HandleFunc" {
				// Try to extract path from first argument
				if len(call.Args) > 0 {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok {
						endpointPath := strings.Trim(lit.Value, "\"")
						result.Endpoints = append(result.Endpoints, Endpoint{
							Path:    endpointPath,
							Method:  strings.ToUpper(funcName),
							Handler: "Detected in " + filepath.Base(path),
							File:    path,
							Line:    fset.Position(call.Pos()).Line,
						})
					}
				}
			}
		}
		return true
	})

	return nil
}

func isHTTPMethod(name string) bool {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	name = strings.ToUpper(name)
	for _, m := range methods {
		if name == m {
			return true
		}
	}
	return false
}

func uniqueStrings(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
