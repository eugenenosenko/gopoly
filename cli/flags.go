package cli

import (
	"flag"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/config"
)

var (
	cfg      = flag.String("c", ".gopoly.yaml", "config file that contains gopoly configuration")
	pack     = flag.String("p", "", "scoped package path where models are located")
	strategy = flag.String("d", "strict", "decoding strategy, either 'strict' or 'discriminator'")
	out      = flag.String("o", "", "output filename that will contain generated code")
	method   = flag.String("m", "Is{{.Name}}", "marker method or template that is used to identify polymorphic relations")

	types         = &TypesInput{usage: "codegen configuration for the polymorphic types"}
	defaultConfig = &config.Config{
		Types:            []*config.TypeDefinition{},
		DecodingStrategy: code.DecodingStrategyStrict,
		MarkerMethod:     "Is{{.Name}}",
		Package:          "",
	}
)

func init() {
	flag.Var(types, "t", types.usage)
}

type TypesInput struct {
	usage string
	types []*config.TypeDefinition
}

func (t *TypesInput) Set(s string) error {
	for _, def := range strings.Split(s, ";") {
		if def == "" {
			continue
		}
		fields := strings.Fields(def)
		typ, entries := fields[0], fields[1:]
		if err := t.set(typ, entries); err != nil {
			return err
		}
	}
	return nil
}

func (t *TypesInput) set(typ string, details []string) error {
	genDef := &config.TypeDefinition{
		Name: typ,
	}
	for _, detail := range details {
		kv := strings.Split(detail, "=")
		key, value := kv[0], kv[1]
		switch key {
		case "subtypes":
			genDef.Variants = append(genDef.Variants, strings.Split(value, ",")...)
		case "marker_method":
			genDef.MarkerMethod = value
		case "discriminator.field":
			genDef.Discriminator.Field = value
		case "discriminator.mapping":
			genDef.Discriminator.Mapping = make(map[string]string, 0)
			for _, mapping := range strings.Split(value, ",") {
				m := strings.Split(mapping, ":")
				kind, implementation := m[0], m[1]
				genDef.Discriminator.Mapping[kind] = implementation
			}
		case "decoding_strategy":
			genDef.DecodingStrategy = code.DecodingStrategy(value)
		case "filename":
			genDef.Output.Filename = value
		}
	}
	t.types = append(t.types, genDef)
	return nil
}

func (t *TypesInput) String() string {
	data, err := yaml.Marshal(t.types)
	if err != nil {
		return "main.TypesInput.String: failed build string representation of 'typesInput'"
	}
	return string(data)
}

var _ flag.Value = (*TypesInput)(nil)
