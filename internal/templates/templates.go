package templates

import (
	_ "embed"
)

//go:embed default_template.gotpl
var defaultJSONTemplate string

// DefaultJSONTemplate returns predefined JSON go-template with default configuration for code generation
func DefaultJSONTemplate() string {
	return defaultJSONTemplate
}
