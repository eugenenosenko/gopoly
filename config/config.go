package config

import (
	"encoding/json"
	"fmt"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/internal/xslices"
)

// TypeDefinition
type TypeDefinition struct {
	Name             string                  `yaml:"name"`
	Variants         []string                `yaml:"variants,omitempty"`
	MarkerMethod     string                  `yaml:"marker_method,omitempty"`
	DecodingStrategy code.DecodingStrategy   `yaml:"decoding_strategy,omitempty"`
	Discriminator    DiscriminatorDefinition `yaml:"discriminator,omitempty"`
	Package          string                  `yaml:"package,omitempty"`
	Output           *OutputConfig           `yaml:"output,omitempty"`
}

type OutputConfig struct {
	Filename string `yaml:"filename"`
}

type TypesList []*TypeDefinition

func (tts TypesList) AssociateByTypeName() map[string]*TypeDefinition {
	return xslices.ToMap[TypesList, map[string]*TypeDefinition](
		tts,
		func(t *TypeDefinition) string {
			return t.Name
		},
		nil,
	)
}

type Config struct {
	Types            TypesList             `yaml:"types"`
	DecodingStrategy code.DecodingStrategy `yaml:"decoding_strategy"`
	MarkerMethod     string                `yaml:"marker_method"`
	Output           *OutputConfig         `yaml:"output"`
	Package          string                `yaml:"package"`
}

func (tts TypesList) AssociateByPkgName() map[string]TypesList {
	res := make(map[string]TypesList, 0)
	for _, def := range tts {
		res[def.Package] = append(res[def.Package], def)
	}
	return res
}

func (tts TypesList) AssociateByOutput() map[string]TypesList {
	res := make(map[string]TypesList, 0)
	for _, def := range tts {
		res[def.Output.Filename] = append(res[def.Output.Filename], def)
	}
	return res
}

type DiscriminatorDefinition struct {
	Field   string            `yaml:"field"`
	Mapping map[string]string `yaml:"mapping"`
}

func (c *Config) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return "config.Config.String: failed build string representation"
	}
	return string(data)
}

var _ fmt.Stringer = (*Config)(nil)
