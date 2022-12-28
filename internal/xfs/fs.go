package xfs

import (
	"io"
	"os"
)

type FileCreatorFunc func(name string) (*os.File, error)

func (f FileCreatorFunc) Create(name string) (io.WriteCloser, error) {
	return f(name)
}

type StreamCreator interface {
	Create(name string) (io.WriteCloser, error)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

var _ StreamCreator = FileCreatorFunc(os.Create)
