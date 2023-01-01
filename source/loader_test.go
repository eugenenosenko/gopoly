package source

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/config"
)

func TestLoad(t *testing.T) {
	t.Run("should correctly load source from packages", func(t *testing.T) {
		l, err := NewLoader(&Config{
			Logf:     func(_ string, _ ...any) {},
			LoadFunc: LoadFromPackage,
		})
		require.NoError(t, err)

		got, err := l.Load(context.Background(), []*config.TypeDefinition{
			{
				Type:             "Runner",
				Subtypes:         []string{"FastRunner", "SlowRunner"},
				MarkerMethod:     "IsRunner",
				DecodingStrategy: config.DecodingStrategyStrict,
				Package:          "github.com/eugenenosenko/gopoly/source/testdata",
				Output:           &config.OutputConfig{Filename: "out.gen.go"},
			},
		})
		runner := &code.Interface{
			Name:         "Runner",
			MarkerMethod: "IsRunner",
			Pkg:          "github.com/eugenenosenko/gopoly/source/testdata",
		}
		a := &code.Variant{Name: "FastRunner", Fields: code.PolyFieldList{
			{
				Name:      "Runner",
				Tags:      "`json:\"runner\"`",
				Interface: runner,
				Kind:      code.KindScalar,
			},
		}, Interface: runner}
		b := &code.Variant{Name: "SlowRunner", Fields: code.PolyFieldList{}, Interface: runner}

		runner.Variants = code.VariantList{a, b}
		want := code.SourceList{
			{
				Package:    "github.com/eugenenosenko/gopoly/source/testdata",
				Interfaces: code.InterfaceList{runner},
				Imports:    nil,
			},
		}
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
