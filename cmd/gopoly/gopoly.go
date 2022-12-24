package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/scanner"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime/debug"
	"strings"
	"syscall"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/eugenenosenko/gopoly/internal/codegen"
	"github.com/eugenenosenko/gopoly/internal/config"
	"github.com/eugenenosenko/gopoly/internal/poly"
	"github.com/eugenenosenko/gopoly/internal/source"
	"github.com/eugenenosenko/gopoly/internal/xfs"
)

var (
	cfg      = flag.String("c", "", "config file that contains gopoly configuration")
	pack     = flag.String("p", "", "scoped package path where models are located")
	strategy = flag.String("d", "strict", "decoding strategy, either 'strict' or 'discriminator'")
	out      = flag.String("o", "", "output filename that will contain generated code")
	method   = flag.String("m", "Is{{.Type}}", "marker method or template that is used to identify polymorphic relations")

	types         = &generateInput{usage: "codegen configuration for the polymorphic types"}
	defaultConfig = &config.Config{
		Types:            []*config.TypeDefinition{},
		DecodingStrategy: config.DecodingStrategyStrict,
		MarkerMethod:     "Is{{.Type}}",
		Output:           config.OutputConfig{Filename: *out},
		Package:          "",
	}
	exit = 0
)

type generateInput struct {
	usage string
	types []*config.TypeDefinition
}

// Set takes in command input and converts it into []*poly.TypeDefinition structure i.e.
// -types "AdvertBase\
// subtypes=RentAdvert,SellAdvert\
// marker_method=i{{.Type}}\
// discriminator.field=type\
// discriminator.mapping=rent:RentAdvert,sell:SellAdvert\
// decoding_strategy=discriminator
// template_file=template.txt
// filename=out.gen.go;.
func (t *generateInput) Set(s string) error {
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

func (t *generateInput) set(typ string, details []string) error {
	genDef := &config.TypeDefinition{
		Type: typ,
	}
	for _, detail := range details {
		kv := strings.Split(detail, "=")
		key, value := kv[0], kv[1]
		switch key {
		case "subtypes":
			genDef.Subtypes = append(genDef.Subtypes, strings.Split(value, ",")...)
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
			genDef.DecodingStrategy = config.DecodingStrategy(value)
		case "filename":
			genDef.Output.Filename = value
		}
	}
	t.types = append(t.types, genDef)
	return nil
}

func (t *generateInput) String() string {
	data, err := yaml.Marshal(t.types)
	if err != nil {
		return "main.generateInput.String: failed build string representation of 'typesInput'"
	}
	return string(data)
}

var _ flag.Value = (*generateInput)(nil)

func init() {
	flag.Var(types, "g", types.usage)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: gopoly [-p package-path] [-c config-file] [-d decoding-strategy]"+
		" [-o output-file] [-m marker-method] [-t template-file] [-g generate-type-info] [-v verbose]\n")
	flag.PrintDefaults()
	exit = 2
}

func reportAndExit(err error) {
	scanner.PrintError(os.Stderr, err)
	exit = 2
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			exit = 3
			_, _ = fmt.Fprintf(os.Stderr, "panic occurred; %v, stack %s", r, debug.Stack())
		}
		os.Exit(exit)
	}()

	flag.Usage = usage
	flag.Parse()

	conf, err := newMergedConfig() // combine configuration
	if err != nil {
		reportAndExit(err)
	}
	s, err := source.NewLoader(&source.Config{
		LoadFunc: source.LoadFromPackage,
		Logf:     log.Printf,
	})
	if err != nil {
		reportAndExit(errors.Wrapf(err, "creating source.Loader"))
	}

	generator, err := codegen.NewGenerator(&codegen.Config{
		FS:   xfs.CreateFunc(os.Create),
		Logf: log.Printf,
	})
	if err != nil {
		reportAndExit(errors.Wrapf(err, "creating codegen.Generator"))
	}

	runner, err := poly.NewRunner(&poly.RunnerConfig{
		SourceLoader:  s,
		CodeGenerator: generator,
		Logf:          log.Printf,
	})
	if err != nil {
		reportAndExit(errors.Wrapf(err, "creating poly.Runner"))
	}

	// trap cancellation signal on the context
	ctx, cancel := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(term)
		cancel()
	}()

	go func() {
		select {
		case <-term:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err = runner.Run(ctx, conf); err != nil {
		reportAndExit(err)
	}
}

func newMergedConfig() (*config.Config, error) {
	var target *config.Config
	if filename := *cfg; filename != "" {
		yamlConf, err := config.NewConfigFromYAML(filename)
		if err != nil {
			return nil, err
		}
		target = yamlConf
	} else {
		target = defaultConfig
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
		ntype[t.Type] = t
	}

	for _, t := range ntype {
		if t.Package == "" {
			t.Package = target.Package
		}
		if t.MarkerMethod == "" {
			t.MarkerMethod = target.MarkerMethod
		}
		if t.Output.Filename == "" {
			t.Output.Filename = target.Output.Filename
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
	}
	target.Types = maps.Values(ntype)

	rx := regexp.MustCompile(`{{\s?\.(\w+)\s?}}`)
	for _, def := range target.Types {
		if m := def.MarkerMethod; rx.MatchString(m) {
			parse, err := template.New("").Parse(m)
			if err != nil {
				return nil, errors.New("can't have discriminator decoding and empty mapping")
			}
			var b bytes.Buffer
			if err = parse.Execute(&b, map[string]string{"Type": def.Type}); err != nil {
				return nil, errors.New("can't have discriminator decoding and empty mapping")
			}
			def.MarkerMethod = b.String()
		}
	}

	return target, nil
}
