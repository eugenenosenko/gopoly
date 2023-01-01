package config

import (
	"encoding/json"
	"fmt"
	"go/build"

	"github.com/eugenenosenko/gopoly/internal/xslices"
)

type DecodingStrategy string

func (s DecodingStrategy) String() string {
	return string(s)
}

func (s DecodingStrategy) IsDiscriminator() bool {
	return s == DecodingStrategyDiscriminator
}

func (s DecodingStrategy) IsStrict() bool {
	return s == DecodingStrategyStrict
}

func (s DecodingStrategy) IsValid() bool {
	switch s {
	case DecodingStrategyStrict, DecodingStrategyDiscriminator:
		return true
	default:
		return false
	}
}

const (
	DecodingStrategyStrict        = DecodingStrategy("strict")
	DecodingStrategyDiscriminator = DecodingStrategy("discriminator")
)

type TypeDefinition struct {
	Type             string                  `yaml:"type"`
	Subtypes         []string                `yaml:"subtypes,omitempty"`
	MarkerMethod     string                  `yaml:"marker_method,omitempty"`
	DecodingStrategy DecodingStrategy        `yaml:"decoding_strategy,omitempty"`
	Discriminator    DiscriminatorDefinition `yaml:"discriminator,omitempty"`
	Package          string                  `yaml:"package,omitempty"`
	Output           *OutputConfig           `yaml:"output,omitempty"`
}

type Package string

func (p Package) String() string {
	return string(p)
}

func (p Package) Name() string {
	return string(p)
}

func (p Package) Dir() string {
	return packageNameToDir(p.Name())
}

func packageNameToDir(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return p.Dir
}

type OutputConfig struct {
	Filename string `yaml:"filename"`
}

type TypesList []*TypeDefinition

func (tts TypesList) AssociateByTypeName() map[string]*TypeDefinition {
	return xslices.ToMap[TypesList, map[string]*TypeDefinition](
		tts,
		func(t *TypeDefinition) string {
			return t.Type
		},
		nil,
	)
}

type Config struct {
	Types            TypesList        `yaml:"types"`
	DecodingStrategy DecodingStrategy `yaml:"decoding_strategy"`
	MarkerMethod     string           `yaml:"marker_method"`
	Output           *OutputConfig    `yaml:"output"`
	Package          string           `yaml:"package"`
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

var (
	_ fmt.Stringer = (*Config)(nil)
	_ fmt.Stringer = (*DecodingStrategy)(nil)
	_ fmt.Stringer = (*Package)(nil)
)
