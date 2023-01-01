package testdata

import (
	"github.com/eugenenosenko/gopoly/source/testdata/c"
)

type Runner interface {
	IsRunner()
}

type FastRunner struct {
	Name   string `json:"name"`
	Runner Runner `json:"runner"`
	C      c.C    `json:"c"`
}

func (a *FastRunner) IsRunner() {}

type SlowRunner struct {
	Name string
}

func (a *SlowRunner) IsRunner() {}
