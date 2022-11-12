package poly

import (
	"context"
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/internal/config"
	"github.com/eugenenosenko/gopoly/internal/source"
)

var jsonTagRegex = regexp.MustCompile(`json:"(.+),?.*"`)

type RunnerConfig struct {
	SourceLoader source.Loader
	Logf         func(format string, args ...any)
}

type Runner struct {
	sourceLoader source.Loader
	logf         func(format string, args ...any)
}

func NewRunner(c *RunnerConfig) (*Runner, error) {
	if c == nil {
		return nil, errors.New("poly.NewRunner: config is nil")
	}
	return &Runner{sourceLoader: c.SourceLoader, logf: c.Logf}, nil
}

func (r *Runner) Run(ctx context.Context, c *config.Config) error {
	i := newRunInfo()
	r.logf("Running info: %s", i)
	r.logf("Running config:\n%s", c)

	source, err := r.sourceLoader.Load(ctx, c.Types)
	if err != nil {
		return errors.Wrapf(err, "loading source from packages")
	}
	fmt.Println(source)
	return nil
}
