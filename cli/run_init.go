package cli

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/eugenenosenko/gopoly/config"
)

func (a *App) RunInit() error {
	a.Logf("Creating example configuration YAML file .gopoly.yaml")
	file, err := os.Create(".gopoly.yaml")
	if err != nil {
		return errors.Wrap(err, "creating .gopoly.yaml file")
	}

	body, err := yaml.Marshal(config.Config{
		Types: config.TypesList{
			{
				Name:             "A",
				Variants:         []string{"AImpl", "BImpl"},
				MarkerMethod:     "Is{{ .Name }}",
				DecodingStrategy: "discriminator",
				Package:          "github.com/username/example/events",
				Discriminator: config.DiscriminatorDefinition{
					Field:   "type",
					Mapping: map[string]string{"A": "AImpl", "B": "BImpl"},
				},
				Output: &config.OutputConfig{
					Filename: "a.gen.go",
				},
			},
			{
				Name: "C",
			},
		},
		MarkerMethod:     "Is{{ .Name }}",
		DecodingStrategy: "strict",
		Output: &config.OutputConfig{
			Filename: "gopoly.gen.go",
		},
		Package: "github.com/username/example/models",
	})
	if err != nil {
		return errors.Wrap(err, "marshaling config")
	}

	if _, err = file.Write(body); err != nil {
		return errors.Wrap(err, "writing config to file")
	}
	return file.Close()
}
