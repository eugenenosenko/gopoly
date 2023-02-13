package cli

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/eugenenosenko/gopoly/config"
	"github.com/eugenenosenko/gopoly/internal/xfs"
)

func MergedConfigProvider() (*config.Config, error) {
	flag.Parse()
	flag.Usage = usage
	return newMergedConfig()
}

func newMergedConfig() (*config.Config, error) {
	var target *config.Config
	if filename := *cfg; filename != "" && xfs.FileExists(filename) {
		yamlConf, err := newConfigFromYAML(filename)
		if err != nil {
			return nil, err
		}
		target = yamlConf
	} else {
		target = defaultConfig
	}
	if filename := *out; filename != "" {
		target.Output.Filename = filename
	}
	if m := *method; m != "" {
		target.MarkerMethod = m
	}
	if p := *pack; p != "" {
		target.Package = p
	}
	if s := *strategy; s != "" {
		target.DecodingStrategy = config.DecodingStrategy(s)
	}

	ntype := target.Types.AssociateByTypeName()
	// overwrite what is defined in config file with input from CLI
	for _, t := range types.types {
		ntype[t.Name] = t
	}

	for _, t := range ntype {
		if t.Package == "" {
			t.Package = target.Package
		}
		if t.MarkerMethod == "" {
			t.MarkerMethod = target.MarkerMethod
		}
		if t.Output == nil || t.Output.Filename == "" {
			t.Output = target.Output
		}
		// validate decoding strategy inputs
		if t.DecodingStrategy == config.DecodingStrategyStrict && len(t.Discriminator.Mapping) > 0 {
			return nil, errors.New("can't have discriminator mapping & strict decoding")
		}
		if t.DecodingStrategy == config.DecodingStrategyDiscriminator && len(t.Discriminator.Mapping) == 0 {
			return nil, errors.New("can't have discriminator decoding and empty mapping")
		}
		if t.DecodingStrategy == "" && len(t.Discriminator.Mapping) > 0 {
			t.DecodingStrategy = config.DecodingStrategyDiscriminator
		} else if t.DecodingStrategy == "" {
			t.DecodingStrategy = target.DecodingStrategy
		}
		if ds := t.DecodingStrategy; !ds.IsValid() {
			return nil, fmt.Errorf("not a valid decoding-strategy %s", ds)
		}
	}
	target.Types = maps.Values(ntype)

	rx := regexp.MustCompile(`{{\s?\.(\w+)\s?}}`)
	for _, def := range target.Types {
		if m := def.MarkerMethod; rx.MatchString(m) {
			parse, err := template.New("").Parse(m)
			if err != nil {
				return nil, errors.Wrap(err, "parsing marker-method template field")
			}
			var b bytes.Buffer
			if err = parse.Execute(&b, map[string]string{"Name": def.Name}); err != nil {
				return nil, errors.Wrap(err, "executing marker-method template")
			}
			def.MarkerMethod = b.String()
		}
	}

	return target, nil
}

func newConfigFromYAML(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "config.NewFromYAML: reading filename %s", filename)
	}
	var c config.Config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrapf(err, "config.NewFromYAML: unmarshaling filename %s", filename)
	}
	return &c, nil
}
