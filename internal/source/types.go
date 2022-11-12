package source

import (
	"go/ast"

	"github.com/eugenenosenko/gopoly/internal/config"
)

type (
	InterfaceTypeMap map[string]Entry[*ast.InterfaceType]
	StructsTypeMap   map[string]Entry[*ast.StructType]
	FuncDecMap       map[string]Entry[*ast.FuncDecl]

	DecType interface {
		*ast.InterfaceType | *ast.FuncDecl | *ast.StructType
	}

	Entry[T DecType] struct {
		Type *config.TypeInfo
		Dec  T
	}

	DecInfo struct {
		Dec     ast.Decl
		PkgName string
		PkgPath string
		Imports []string
	}
)

type Config struct {
	LoadFunc func(paths []string) ([]*Package, error)
	Logf     func(format string, args ...any)
}
