package codegen

import (
	"fmt"
	"strings"
	"text/template"

	"golang.org/x/exp/maps"

	"github.com/eugenenosenko/gopoly/internal/code"
	"github.com/eugenenosenko/gopoly/internal/xslices"
)

func DefaultFuncs() template.FuncMap {
	return template.FuncMap{
		"lower":         strings.ToLower,
		"upper":         strings.ToUpper,
		"prefixed":      prefixedField,
		"lookupImports": lookupImports,
		"dedupTypes":    dedupTypes,
	}
}

// lookupImports will fetch required imports for the required Data, without duplicates and
// returns it in a proper string format.
func lookupImports(d *Data) string {
	m := xslices.ToMap[[]*code.Import, map[string]*code.Import](d.Imports, func(i *code.Import) string {
		return i.Path
	}, nil)

	iis := make(map[string]*code.Import, 0)
	for _, t := range d.Types {
		for _, v := range maps.Values(t.Variants) {
			for _, field := range v.Fields {
				if i, ok := m[field.Interface.Pkg]; ok {
					iis[i.ShortName+i.Path] = i
				}
			}
		}
	}

	if len(iis) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n\n")
	for _, ii := range iis {
		sb.WriteString("\t")
		if ii.Aliased {
			sb.WriteString(ii.ShortName)
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%q", ii.Path))
		sb.WriteString("\n")
	}
	return sb.String()
}

// dedupTypes filters out duplicated variants.
// User can define multiple discriminator mappings that match to same type.
// Dedup is required in order to not re-define Unmarshal method for the same code.Variant type.
func dedupTypes(vars map[string]*code.Variant) []*code.Variant {
	variants := xslices.ToMap[[]*code.Variant, map[string]*code.Variant](
		maps.Values(vars), func(v *code.Variant) string {
			return v.Name
		}, nil)
	return maps.Values(variants)
}

// prefixedField checks whether code.PolyField has an import-prefix and if it has one
// returns a composed field name, i.e. m.MyModel or models.MyModel.
func prefixedField(f code.PolyField) string {
	if f.Prefix != "" {
		return fmt.Sprintf("%s.", f.Prefix)
	}
	return ""
}
