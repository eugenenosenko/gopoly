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
	"syscall"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	config2 "github.com/eugenenosenko/gopoly/internal/config"
	"github.com/eugenenosenko/gopoly/internal/poly"
	"github.com/eugenenosenko/gopoly/internal/source"
	"github.com/eugenenosenko/gopoly/internal/xslices"
)

var (
	config   = flag.String("c", "", "config file that contains gopoly configuration")
	pack     = flag.String("p", "", "scoped package path where models are located")
	strategy = flag.String("d", "strict", "decoding strategy, either 'strict' or 'discriminator'")
	out      = flag.String("o", "poly.gen.go", "output filename that will contain generated code")
	method   = flag.String("m", "Is{{.Type}}", "marker method or template that is used to identify polymorphic relations")
	tmpl     = flag.String("t", "codegen.tmpl", "template file that is used to generate marshalling/unmarshalling logic")

	types         = &generateInput{usage: "codegen configuration for the polymorphic types"}
	defaultConfig = &config2.Config{
		Types:            []*config2.TypeInfo{},
		DecodingStrategy: config2.DecodingStrategyStrict,
		MarkerMethod:     "Is{{.Type}}",
		TemplateFile:     "codegen.gotmpl",
		Output:           "poly.gen.go",
		PackagePath:      "",
	}
)

func init() {
	flag.Var(types, "g", types.usage)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: gopoly [-p package-path] [-c config-file] [-d decoding-strategy]"+
		" [-o output-file] [-m marker-method] [-t template-file] [-g generate-type-info] [-v verbose]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func reportAndExit(err error) {
	scanner.PrintError(os.Stderr, err)
	os.Exit(2)
}

func main() {
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

	runner, err := poly.NewRunner(&poly.RunnerConfig{
		SourceLoader: s,
		Logf:         log.Printf,
	})
	if err != nil {
		reportAndExit(errors.Wrapf(err, "creating poly.Runner"))
	}

	// trap cancellation signal on the context
	ctx, cancel := context.WithCancel(context.Background())
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(terminate)
		cancel()
	}()

	go func() {
		select {
		case <-terminate:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err = runner.Run(ctx, conf); err != nil {
		reportAndExit(err)
	}
	os.Exit(0)
}

func newMergedConfig() (*config2.Config, error) {
	target := &config2.Config{}
	if filename := *config; filename != "" {
		yamlConf, err := config2.NewConfigFromYAML(filename)
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
		target.PackagePath = p
	}
	if s := *strategy; s != "" {
		target.DecodingStrategy = config2.DecodingStrategy(s)
	}
	if o := *out; o != "" {
		target.Output = o
	}
	if t := *tmpl; t != "" {
		target.TemplateFile = t
	}

	tts := xslices.ToMap[[]*config2.TypeInfo, map[string]*config2.TypeInfo](
		target.Types,
		func(t *config2.TypeInfo) string { return t.Type }, nil,
	)

	// overwrite what is defined in config file with input from CLI
	for _, t := range types.types {
		tts[t.Type] = t
	}

	for _, t := range tts {
		if t.PackagePath == "" {
			t.PackagePath = target.PackagePath
		}
		if t.MarkerMethod == "" {
			t.MarkerMethod = target.MarkerMethod
		}
		if t.Filename == "" {
			t.Filename = target.TemplateFile
		}
		if t.TemplateFile == "" {
			t.TemplateFile = target.TemplateFile
		}
		// validate decoding strategy inputs
		if t.DecodingStrategy == config2.DecodingStrategyStrict && len(t.Discriminator.Mapping) > 0 {
			return nil, errors.New("can't have discriminator mapping & strict decoding")
		}
		if t.DecodingStrategy == config2.DecodingStrategyDiscriminator && len(t.Discriminator.Mapping) == 0 {
			return nil, errors.New("can't have discriminator decoding and empty mapping")
		}
		if t.DecodingStrategy == "" && len(t.Discriminator.Mapping) > 0 {
			t.DecodingStrategy = config2.DecodingStrategyDiscriminator
		} else if t.DecodingStrategy == "" {
			t.DecodingStrategy = target.DecodingStrategy
		}
	}
	target.Types = maps.Values(tts)

	rx := regexp.MustCompile(`{{\s?\.(\w+)\s?}}`)
	for _, info := range target.Types {
		if mthd := info.MarkerMethod; rx.MatchString(mthd) {
			parse, err := template.New("").Parse(mthd)
			if err != nil {
				return nil, errors.New("can't have discriminator decoding and empty mapping")
			}
			var buf bytes.Buffer
			if err = parse.Execute(&buf, map[string]string{"Type": info.Type}); err != nil {
				return nil, errors.New("can't have discriminator decoding and empty mapping")
			}
			info.MarkerMethod = buf.String()
		}
	}

	return target, nil
}
