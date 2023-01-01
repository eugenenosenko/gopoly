package codegen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

type WriterProvider interface {
	Provide(name string) (io.WriteCloser, error)
}

//go:embed template.gotpl
var DefaultJSONTemplate string

type Generator interface {
	Generate(*Task) error
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

func (g *generator) Generate(t *Task) error {
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
		Funcs(DefaultFuncs()).
		Parse(t.Template))
	if err := temp.Execute(w, t.Data); err != nil {
		return errors.Wrapf(err, "generating code to file %s", t.Filename)
	}
	return nil
}

var _ Generator = (*generator)(nil)
