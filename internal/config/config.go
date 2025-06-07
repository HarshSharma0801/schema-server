package config

// Config represents the server configuration
type Config struct {
	Contract Contract `json:"contract" yaml:"contract"`
}

// Contract represents contract-specific configuration
type Contract struct {
	Path     string   `json:"path" yaml:"path"`
	Generate bool     `json:"generate" yaml:"generate"`
	Download bool     `json:"download" yaml:"download"`
	Driven   string   `json:"driven" yaml:"driven"`
	Mappings Mappings `json:"mappings" yaml:"mappings"`
}

// Mappings represents service mappings configuration
type Mappings struct {
	Self            string              `json:"self" yaml:"self"`
	ServicesMapping map[string][]string `json:"servicesMapping" yaml:"servicesMapping"`
}
