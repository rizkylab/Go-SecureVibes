package scanner

// Config holds the configuration for the security scan
type Config struct {
	ProjectPath  string
	OutputFile   string
	OutputFormat string
	MinSeverity  string
	SkipDAST     bool
	SkipThreats  bool
	SkipReview   bool
	Verbose      bool
	FailOn       string
	Excludes     []string
	TargetURL    string // New field for DAST target
	CIMode       bool
}
