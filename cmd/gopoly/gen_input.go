package main

import (
	"flag"
	"strings"

	"gopkg.in/yaml.v3"

	config2 "github.com/eugenenosenko/gopoly/internal/config"
)

type generateInput struct {
	usage string
	types []*config2.TypeInfo
}

// Set takes in command input and converts it into []*poly.TypeInfo structure i.e.
// -types "AdvertBase\
// subtypes=RentAdvert,SellAdvert\
// marker_method=i{{.Type}}\
// discriminator.field=type\
// discriminator.mapping=rent:RentAdvert,sell:SellAdvert\
// decoding_strategy=discriminator
// template_file=template.txt
// filename=out.gen.go;
func (t *generateInput) Set(s string) error {
	for _, info := range strings.Split(s, ";") {
		if info == "" {
			continue
		}
		fields := strings.Fields(info)
		typ, entries := fields[0], fields[1:]
		if err := t.set(typ, entries); err != nil {
			return err
		}
	}
	return nil
}

func (t *generateInput) set(typ string, details []string) error {
	genInfo := &config2.TypeInfo{
		Type: typ,
	}
	for _, detail := range details {
		kv := strings.Split(detail, "=")
		key, value := kv[0], kv[1]
		switch key {
		case "subtypes":
			for _, sub := range strings.Split(value, ",") {
				genInfo.Subtypes = append(genInfo.Subtypes, sub)
			}
		case "marker_method":
			genInfo.MarkerMethod = value
		case "discriminator.field":
			genInfo.Discriminator.Field = value
		case "discriminator.mapping":
			genInfo.Discriminator.Mapping = make(map[string]string, 0)
			for _, mapping := range strings.Split(value, ",") {
				m := strings.Split(mapping, ":")
				kind, implementation := m[0], m[1]
				genInfo.Discriminator.Mapping[kind] = implementation
			}
		case "decoding_strategy":
			genInfo.DecodingStrategy = config2.DecodingStrategy(value)
		case "template_file":
			genInfo.TemplateFile = value
		case "filename":
			genInfo.Filename = value
		}
	}
	t.types = append(t.types, genInfo)
	return nil
}

func (t *generateInput) isEmpty() bool {
	return len(t.types) == 0
}

func (t *generateInput) String() string {
	data, err := yaml.Marshal(t.types)
	if err != nil {
		return "main.generateInput.String: failed build string representation of 'typesInput'"
	}
	return string(data)
}

var _ flag.Value = (*generateInput)(nil)
