package codegen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/eugenenosenko/gopoly/internal/code"
)

var funcs = template.FuncMap{
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"prefixed": prefixedField,
}

func prefixedField(f code.PolyField) string {
	if f.Prefix != "" {
		return fmt.Sprintf("%s.", f.Prefix)
	}
	return ""
}
