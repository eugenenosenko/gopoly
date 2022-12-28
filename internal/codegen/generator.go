package codegen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/internal/xfs"
)

//go:embed template.gotpl
var DefaultTemplate string

type Generator interface {
	Generate(*Task) error
}

type Config struct {
	FS   xfs.StreamCreator
	Logf func(format string, args ...any)
}

type generator struct {
	fs   xfs.StreamCreator
	logf func(format string, args ...any)
}

func NewGenerator(c *Config) (*generator, error) {
	if c == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &generator{fs: c.FS, logf: c.Logf}, nil
}

func (g *generator) Generate(t *Task) error {
	out, err := g.fs.Create(t.Filename)
	if err != nil {
		return errors.Wrapf(err, "creating %s file", t.Filename)
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "closing file %s", t.Filename)
		}
	}(out)

	temp := template.Must(template.New("gopoly").
		Funcs(DefaultFuncs()).
		Parse(t.Template))
	if err := temp.Execute(out, t.Data); err != nil {
		return errors.Wrapf(err, "generating code to file %s", t.Filename)
	}
	return nil
}

var _ Generator = (*generator)(nil)
