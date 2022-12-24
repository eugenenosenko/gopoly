package config

import (
	"fmt"
	"go/build"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

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
	Subtypes         []string                `yaml:"subtypes"`
	MarkerMethod     string                  `yaml:"marker_method"`
	DecodingStrategy DecodingStrategy        `yaml:"decoding_strategy"`
	Discriminator    DiscriminatorDefinition `yaml:"discriminator"`
	Package          string                  `yaml:"package"`
	Output           OutputConfig            `yaml:"output"`
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
	Output           OutputConfig     `yaml:"output"`
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
	data, err := yaml.Marshal(c)
	if err != nil {
		return "config.Config.String: failed build string representation"
	}
	return string(data)
}

func NewConfigFromYAML(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "config.NewConfigFromYAML: reading filename %s", filename)
	}
	var c Config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrapf(err, "config.NewConfigFromYAML: unmarshaling filename %s", filename)
	}
	return &c, nil
}

var (
	_ fmt.Stringer = (*Config)(nil)
	_ fmt.Stringer = (*DecodingStrategy)(nil)
	_ fmt.Stringer = (*Package)(nil)
)
