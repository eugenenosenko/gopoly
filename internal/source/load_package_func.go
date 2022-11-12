package source

import (
	"go/ast"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

type Package struct {
	Name  string
	Path  string
	Files []*ast.File
}

func LoadFromPackage(paths []string) ([]*Package, error) {
	packs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedSyntax |
				packages.NeedName |
				packages.NeedCompiledGoFiles,
		},
		paths...,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "loading package path %v", paths)
	}

	var err2 error
	packages.Visit(packs, nil, func(pkg *packages.Package) {
		for _, e := range pkg.Errors {
			err2 = multierr.Append(err2, e)
		}
	})
	if err2 != nil {
		return nil, errors.Wrapf(err2, "reading package with path %v", paths)
	}

	var res []*Package
	for _, pack := range packs {
		res = append(res, &Package{
			Name:  pack.Name,
			Path:  pack.PkgPath,
			Files: pack.Syntax,
		})
	}
	return res, nil
}
