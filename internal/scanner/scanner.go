package scanner

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/architecture"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/dast"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/staticanalysis"
	"github.com/rizkylab/Go-SecureVibes/internal/agents/threatmodel"
	"github.com/rizkylab/Go-SecureVibes/internal/report"
)

// Scanner orchestrates the security scan
type Scanner struct {
	Config Config
}

// New creates a new Scanner instance
func New(config Config) *Scanner {
	return &Scanner{
		Config: config,
	}
}

// Run executes the full scan pipeline
func (s *Scanner) Run() (int, error) {
	startTime := time.Now()

	// 1. Architecture Assessment
	color.Blue("🏗️  Phase 1: Architecture Assessment...")
	assessAgent := architecture.New(s.Config.ProjectPath, s.Config.Excludes)
	assessResult, err := assessAgent.Run()
	if err != nil {
		return 0, fmt.Errorf("assessment failed: %w", err)
	}
	color.Green("   Assessment complete.")

	// 2. Threat Modeling
	var threatResult *threatmodel.Result
	if !s.Config.SkipThreats {
		color.Blue("🎯 Phase 2: Threat Modeling (STRIDE)...")
		threatAgent := threatmodel.New()
		threatResult, err = threatAgent.Run(assessResult)
		if err != nil {
			return 0, fmt.Errorf("threat modeling failed: %w", err)
		}
		color.Green("   Threat model generated.")
	} else {
		color.Yellow("   Skipping Threat Modeling.")
	}

	// 3. Code Review
	var reviewResult *staticanalysis.Result
	if !s.Config.SkipReview {
		color.Blue("🔍 Phase 3: Static Code Review...")
		reviewAgent := staticanalysis.New(s.Config.ProjectPath, s.Config.Excludes)
		reviewResult, err = reviewAgent.Run()
		if err != nil {
			return 0, fmt.Errorf("code review failed: %w", err)
		}
		color.Green("   Code review complete.")
	} else {
		color.Yellow("   Skipping Code Review.")
	}

	// 4. DAST (Optional)
	var dastResult *dast.Result
	if !s.Config.SkipDAST {
		color.Blue("🚀 Phase 4: Dynamic Analysis (DAST)...")
		dastAgent := dast.New(s.Config.TargetURL)
		dastResult, err = dastAgent.Run()
		if err != nil {
			return 0, fmt.Errorf("DAST failed: %w", err)
		}
		color.Green("   DAST complete.")
	} else {
		color.Yellow("   Skipping DAST.")
	}

	// 5. Report Generation
	color.Blue("📊 Phase 5: Generating Report...")
	reporter := report.New(s.Config.OutputFile, s.Config.OutputFormat)
	err = reporter.Generate(assessResult, threatResult, reviewResult, dastResult)
	if err != nil {
		return 0, fmt.Errorf("report generation failed: %w", err)
	}

	duration := time.Since(startTime)
	color.HiGreen("\n✨ Scan finished in %s", duration)
	color.HiGreen("   Report saved to: %s", s.Config.OutputFile)

	// Calculate Max Severity for CI/CD
	maxSeverity := 0
	if reviewResult != nil {
		for _, v := range reviewResult.Vulnerabilities {
			if v.Severity == "High" || v.Severity == "Critical" {
				maxSeverity = 2
			} else if maxSeverity < 1 {
				maxSeverity = 1
			}
		}
	}
	if threatResult != nil {
		for _, t := range threatResult.Threats {
			if t.Severity == "High" || t.Severity == "Critical" {
				maxSeverity = 2
			} else if maxSeverity < 1 {
				maxSeverity = 1
			}
		}
	}

	return maxSeverity, nil
}
