package poly

import (
	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/generator"
	"github.com/eugenenosenko/gopoly/source"
)

type Config struct {
	SourceLoader  source.Loader
	CodeGenerator generator.Generator
	Logf          func(format string, args ...any)
}

// Client is the main component of the application. It takes source.Loader and codegen.Generator
type Client struct {
	Loader    source.Loader
	Generator generator.Generator
	Logf      func(format string, args ...any)
}

func NewClient(c *Config) (*Client, error) {
	if c == nil {
		return nil, errors.New("poly.NewClient: config is nil")
	}
	return &Client{
		Loader:    c.SourceLoader,
		Logf:      c.Logf,
		Generator: c.CodeGenerator,
	}, nil
}
