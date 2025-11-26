package assessment

// ComponentType defines the type of a component
type ComponentType string

const (
	ComponentTypeService  ComponentType = "service"
	ComponentTypeDatabase ComponentType = "database"
	ComponentTypeExternal ComponentType = "external_api"
	ComponentTypeLibrary  ComponentType = "library"
)

// Component represents a logical part of the application
type Component struct {
	Name         string        `json:"name"`
	Type         ComponentType `json:"type"`
	Path         string        `json:"path"`
	Description  string        `json:"description"`
	Dependencies []string      `json:"dependencies"` // List of other components it interacts with
}

// Endpoint represents an API endpoint or entry point
type Endpoint struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Handler string `json:"handler"` // Function name
	Line    int    `json:"line"`
	File    string `json:"file"`
}

// DataFlow represents a flow of data between components
type DataFlow struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	DataType    string `json:"data_type"` // e.g., "user_input", "config", "pii"
}

// Result holds the detailed findings from the assessment
type Result struct {
	Components   []Component `json:"components"`
	Endpoints    []Endpoint  `json:"endpoints"`
	DataFlows    []DataFlow  `json:"data_flows"`
	Dependencies []string    `json:"external_dependencies"`
}
