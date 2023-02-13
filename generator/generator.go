package generator

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/codegen"
	"github.com/eugenenosenko/gopoly/templates"
)

type WriterProvider interface {
	Provide(name string) (io.WriteCloser, error)
}

type Generator interface {
	Generate(*codegen.Task) error
}

type Config struct {
	Logf     func(format string, args ...any)
	Provider WriterProvider
}

type templateGenerator struct {
	logf     func(format string, args ...any)
	provider WriterProvider
}

func NewTemplateGenerator(c *Config) (*templateGenerator, error) {
	if c == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &templateGenerator{provider: c.Provider, logf: c.Logf}, nil
}

func (g *templateGenerator) Generate(t *codegen.Task) error {
	w, err := g.provider.Provide(t.Filename)
	if err != nil {
		return errors.Wrapf(err, "creating %s file", t.Filename)
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "closing file %s", t.Filename)
		}
	}(w)

	temp := template.Must(template.New("gopoly").
		Funcs(templates.DefaultFuncs()).
		Parse(t.Template))
	if err := temp.Execute(w, t.Input); err != nil {
		return errors.Wrapf(err, "generating code to file %s", t.Filename)
	}
	return nil
}

var _ Generator = (*templateGenerator)(nil)
