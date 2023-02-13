package generator

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/codegen"
	"github.com/eugenenosenko/gopoly/config"
	"github.com/eugenenosenko/gopoly/templates"
)

type dummyCreator struct{ buf *bytes.Buffer }

func (d *dummyCreator) Write(b []byte) (n int, err error)        { return d.buf.Write(b) }
func (d *dummyCreator) Provide(_ string) (io.WriteCloser, error) { return d, nil }
func (d *dummyCreator) Close() error                             { return nil }

func TestGenerator(t *testing.T) {
	t.Run("should correctly generate code for provided configuration", func(t *testing.T) {
		var b bytes.Buffer
		gen, err := NewTemplateGenerator(&Config{
			Provider: &dummyCreator{&b},
			Logf:     func(_ string, _ ...any) {},
		})
		require.NoError(t, err)

		i := &code.Interface{
			Name:         "Advert",
			MarkerMethod: "IsAdvert",
			Pkg:          "github.com/eugenenosenko/gopoly/internal/models",
		}
		sell := &code.Variant{Name: "SellAdvert", Fields: code.PolyFieldList{
			{
				Name:      "Runner",
				Tags:      "`json:\"sell\"`",
				Interface: i,
				Kind:      code.KindScalar,
			},
		}, Interface: i}
		i.Variants = code.VariantList{sell}

		err = gen.Generate(&codegen.Task{
			Filename: "_",
			Template: templates.DefaultJSONTemplate(),
			Input: &codegen.Input{
				Package: "github.com/eugenenosenko/gopoly/internal/models",
				Types: []*codegen.Type{
					{
						Name:               "Advert",
						Variants:           map[string]*code.Variant{"SELL": sell},
						DecodingStrategy:   config.DecodingStrategyDiscriminator.String(),
						DiscriminatorField: "type",
					},
				},
			},
		})
		require.NoError(t, err)

		data, err := os.ReadFile("testdata/output.golden")
		require.NoError(t, err)
		require.Equal(t, string(data), b.String())
	})
}
