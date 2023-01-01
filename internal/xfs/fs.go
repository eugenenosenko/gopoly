package xfs

import (
	"io"
	"os"
)

type FileWriterProviderFunc func(name string) (*os.File, error)

func (f FileWriterProviderFunc) Provide(name string) (io.WriteCloser, error) {
	return f(name)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
