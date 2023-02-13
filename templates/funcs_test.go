package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/codegen"
	"github.com/eugenenosenko/gopoly/config"
)

func TestLookUpImport(t *testing.T) {
	t.Run("should correctly look up imports without duplicates", func(t *testing.T) {
		i := &code.Interface{
			Name:         "Advert",
			MarkerMethod: "IsAdvert",
			Pkg:          "github.com/eugenenosenko/gopoly/internal/models",
		}
		sell := &code.Variant{Name: "SellAdvert", Fields: code.PolyFieldList{
			{
				Name:      "Runner",
				Tags:      "`json:\"runner\"`",
				Interface: i,
				Kind:      code.KindScalar,
			},
		}, Interface: i}
		rent := &code.Variant{Name: "RentAdvert", Fields: code.PolyFieldList{
			{
				Name:      "Runner",
				Tags:      "`json:\"runner\"`",
				Interface: i,
				Kind:      code.KindScalar,
			},
		}, Interface: i}
		i.Variants = code.VariantList{sell, rent}

		data := &codegen.Input{
			Package: "github.com/eugenenosenko/gopoly/internal/models",
			Types: []*codegen.Type{
				{
					Name:               "Advert",
					Variants:           map[string]*code.Variant{"SELL": sell, "RENT": rent},
					DecodingStrategy:   config.DecodingStrategyDiscriminator.String(),
					DiscriminatorField: "type",
				},
			},
			Imports: []*code.Import{
				{
					ShortName: "c",
					Path:      "github.com/eugenenosenko/gopoly/internal/code",
					Aliased:   true,
				},
				{
					ShortName: "",
					Path:      "github.com/eugenenosenko/gopoly/internal/code",
					Aliased:   false,
				},
				{
					ShortName: "m",
					Path:      "github.com/eugenenosenko/gopoly/internal/models",
					Aliased:   false,
				},
			},
		}
		got := lookupImports(data)
		assert.Equal(t, `"github.com/eugenenosenko/gopoly/internal/models"`, strings.TrimSpace(got))
	})
}
