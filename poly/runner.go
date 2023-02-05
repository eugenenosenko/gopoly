package poly

import (
	"context"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/code"
)

// SourceLoader
type SourceLoader interface {
	// Load
	Load(c context.Context, defs []code.Type) (code.SourceList, error)
}

// CodeGenerator
type CodeGenerator interface {
	// Generate
	Generate(code.Task) error
}

type RunnerConfig struct {
	SourceLoader  SourceLoader
	CodeGenerator CodeGenerator
	Logf          func(format string, args ...any)
}

// Runner is the main component of the application. It takes source.Loader and codegen.Generator
type Runner struct {
	Loader    SourceLoader
	Generator CodeGenerator
	Logf      func(format string, args ...any)
}

func NewRunner(c *RunnerConfig) (*Runner, error) {
	if c == nil {
		return nil, errors.New("poly.NewRunner: config is nil")
	}
	return &Runner{
		Loader:    c.SourceLoader,
		Logf:      c.Logf,
		Generator: c.CodeGenerator,
	}, nil
}
