package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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

type TypeInfo struct {
	Type             string            `yaml:"type"`
	Subtypes         []string          `yaml:"subtypes"`
	MarkerMethod     string            `yaml:"marker_method"`
	Discriminator    DiscriminatorInfo `yaml:"discriminator"`
	DecodingStrategy DecodingStrategy  `yaml:"decoding_strategy"`
	TemplateFile     string            `yaml:"template_file"`
	Filename         string            `yaml:"filename"`
	PackagePath      string            `yaml:"package_path"`
}

type Config struct {
	Types            []*TypeInfo      `yaml:"types"`
	DecodingStrategy DecodingStrategy `yaml:"decoding_strategy"`
	MarkerMethod     string           `yaml:"marker_method"`
	TemplateFile     string           `yaml:"template_file"`
	Output           string           `yaml:"filename"`
	PackagePath      string           `yaml:"package_path"`
}

type DiscriminatorInfo struct {
	Field   string            `yaml:"field"`
	Mapping map[string]string `yaml:"mapping"`
}

func (c *Config) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "poly.Config.String: failed build string representation"
	}
	return string(data)
}

func NewConfigFromYAML(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "poly.NewConfigFromYAML: reading filename %s", filename)
	}
	var c Config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrapf(err, "poly.NewConfigFromYAML: unmarshaling filename %s", filename)
	}
	return &c, nil
}

var (
	_ fmt.Stringer = (*Config)(nil)
	_ fmt.Stringer = (*DecodingStrategy)(nil)
)
