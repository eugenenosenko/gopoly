package templates

import (
	_ "embed"
)

//go:embed default_template.gotpl
var defaultJSONTemplate string

func DefaultJSONTemplate() string {
	return defaultJSONTemplate
}
