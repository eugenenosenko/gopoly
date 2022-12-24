package xfs

import (
	"os"
)

type CreateFunc func(name string) (*os.File, error)

func (f CreateFunc) Create(name string) (*os.File, error) {
	return f(name)
}

type Creator interface {
	Create(name string) (*os.File, error)
}
