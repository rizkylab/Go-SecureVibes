package threatmodel

// Severity levels
const (
	SeverityCritical = "Critical"
	SeverityHigh     = "High"
	SeverityMedium   = "Medium"
	SeverityLow      = "Low"
)

// STRIDE Categories
const (
	StrideSpoofing              = "Spoofing"
	StrideTampering             = "Tampering"
	StrideRepudiation           = "Repudiation"
	StrideInformationDisclosure = "Information Disclosure"
	StrideDenialOfService       = "Denial of Service"
	StrideElevationOfPrivilege  = "Elevation of Privilege"
)

// Threat represents a potential security threat
type Threat struct {
	ID          string `json:"id"`
	Category    string `json:"category"` // STRIDE category
	Title       string `json:"title"`
	Description string `json:"description"`
	Target      string `json:"target"` // Component or Endpoint affected
	Severity    string `json:"severity"`
	Mitigation  string `json:"mitigation"`
}

// Result holds the list of identified threats
type Result struct {
	Threats []Threat `json:"threats"`
}
