package poly

import (
	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/internal/codegen"
	"github.com/eugenenosenko/gopoly/internal/source"
)

type RunnerConfig struct {
	SourceLoader  source.Loader
	CodeGenerator codegen.Generator
	Logf          func(format string, args ...any)
}

// Runner is the main component of the application. It takes source.Loader and codegen.Generator
type Runner struct {
	loader    source.Loader
	generator codegen.Generator
	logf      func(format string, args ...any)
}

func NewRunner(c *RunnerConfig) (*Runner, error) {
	if c == nil {
		return nil, errors.New("poly.NewRunner: config is nil")
	}
	return &Runner{
		loader:    c.SourceLoader,
		logf:      c.Logf,
		generator: c.CodeGenerator,
	}, nil
}
