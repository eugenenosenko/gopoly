package codegen

import (
	_ "embed"
	"fmt"
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
	FS   xfs.Creator
	Logf func(format string, args ...any)
}

type generator struct {
	fs   xfs.Creator
	logf func(format string, args ...any)
}

func NewGenerator(c *Config) (*generator, error) {
	if c == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &generator{fs: c.FS, logf: c.Logf}, nil
}

func (g *generator) Generate(t *Task) error {
	file, err := g.fs.Create(t.Filename)
	if err != nil {
		return errors.Wrapf(err, "creating %s file", t.Filename)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "closing file %s", f.Name())
		}
	}(file)

	temp := template.Must(template.New("gopoly").
		Funcs(funcs).
		Parse(t.Template))
	if err = temp.Execute(file, t.Data); err != nil {
		return errors.Wrapf(err, "generating code to file %s", t.Filename)
	}
	return nil
}

var _ Generator = (*generator)(nil)
