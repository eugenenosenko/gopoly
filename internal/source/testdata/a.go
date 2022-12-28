package testdata

import (
	"github.com/eugenenosenko/gopoly/internal/source/testdata/c"
)

type Runner interface {
	IsRunner()
}

type A struct {
	Name   string `json:"name"`
	Runner Runner `json:"runner"`
	C      c.C    `json:"c"`
}

func (a *A) IsRunner() {}

type B struct {
	Name string
}

func (a *B) IsRunner() {}
