package scanner

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/yourusername/gosecvibes/pkg/assessment"
	"github.com/yourusername/gosecvibes/pkg/codereview"
	"github.com/yourusername/gosecvibes/pkg/dast"
	"github.com/yourusername/gosecvibes/pkg/report"
	"github.com/yourusername/gosecvibes/pkg/threatmodel"
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
func (s *Scanner) Run() error {
	startTime := time.Now()

	// 1. Architecture Assessment
	color.Blue("🏗️  Phase 1: Architecture Assessment...")
	assessAgent := assessment.New(s.Config.ProjectPath, s.Config.Excludes)
	assessResult, err := assessAgent.Run()
	if err != nil {
		return fmt.Errorf("assessment failed: %w", err)
	}
	color.Green("   Assessment complete.")

	// 2. Threat Modeling
	var threatResult *threatmodel.Result
	if !s.Config.SkipThreats {
		color.Blue("🎯 Phase 2: Threat Modeling (STRIDE)...")
		threatAgent := threatmodel.New()
		threatResult, err = threatAgent.Run(assessResult)
		if err != nil {
			return fmt.Errorf("threat modeling failed: %w", err)
		}
		color.Green("   Threat model generated.")
	} else {
		color.Yellow("   Skipping Threat Modeling.")
	}

	// 3. Code Review
	var reviewResult *codereview.Result
	if !s.Config.SkipReview {
		color.Blue("🔍 Phase 3: Static Code Review...")
		reviewAgent := codereview.New(s.Config.ProjectPath, s.Config.Excludes)
		reviewResult, err = reviewAgent.Run()
		if err != nil {
			return fmt.Errorf("code review failed: %w", err)
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
			return fmt.Errorf("DAST failed: %w", err)
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
		return fmt.Errorf("report generation failed: %w", err)
	}

	duration := time.Since(startTime)
	color.HiGreen("\n✨ Scan finished in %s", duration)
	color.HiGreen("   Report saved to: %s", s.Config.OutputFile)

	return nil
}
