package codegen

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/internal/templates"
	"github.com/eugenenosenko/gopoly/poly"
)

type WriterProvider interface {
	Provide(name string) (io.WriteCloser, error)
}

type Config struct {
	Logf     func(format string, args ...any)
	Provider WriterProvider
}

type generator struct {
	logf     func(format string, args ...any)
	provider WriterProvider
}

func NewGenerator(c *Config) (*generator, error) {
	if c == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &generator{provider: c.Provider, logf: c.Logf}, nil
}

func (g *generator) Generate(t code.Task) error {
	w, err := g.provider.Provide(t.Filename())
	if err != nil {
		return errors.Wrapf(err, "creating %s file", t.Filename())
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "closing file %s", t.Filename())
		}
	}(w)

	temp := template.Must(template.New("gopoly").
		Funcs(templates.DefaultFuncs()).
		Parse(t.Template()))
	if err := temp.Execute(w, t.Data); err != nil {
		return errors.Wrapf(err, "generating code to file %s", t.Filename())
	}
	return nil
}

var _ poly.CodeGenerator = (*generator)(nil)
