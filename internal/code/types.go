package code

import (
	"go/build"
	"path"

	"github.com/eugenenosenko/gopoly/internal/xslices"
)

const (
	KindScalar = iota
	KindMap
	KindSlice
)

type (
	InterfaceList []*Interface
	ImportList    []*Import
	VariantList   []*Variant
	PolyFieldList []*PolyField
)

type PolyField struct {
	Name      string
	Tags      string
	Interface *Interface
	Kind      int
	Prefix    string
}

type Variant struct {
	Name      string
	Fields    PolyFieldList
	Interface *Interface
}

type Interface struct {
	Name         string
	MarkerMethod string
	Variants     VariantList
}

func (vvs VariantList) AssociateByVariantName() map[string]*Variant {
	res := make(map[string]*Variant, 0)
	for _, variant := range vvs {
		res[variant.Name] = variant
	}
	return res
}

type Package string

func (p Package) String() string {
	return string(p)
}

func (p Package) Path() string {
	return string(p)
}

func (p Package) Name() string {
	return path.Base(p.Path())
}

func (p Package) Dir() string {
	return packageNameToDir(p.Path())
}

func packageNameToDir(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return p.Dir
}

type SourceList []*Source

type Source struct {
	Package    Package
	Interfaces InterfaceList
	Imports    ImportList
}

func (ss SourceList) AssociateByPkgName() map[Package]*Source {
	res := make(map[Package]*Source, 0)
	for _, s := range ss {
		res[s.Package] = s
	}
	return res
}

func (ii InterfaceList) AssociateByName() map[string]*Interface {
	return xslices.ToMap[[]*Interface, map[string]*Interface](ii, func(t *Interface) string {
		return t.Name
	}, nil)
}

func (ii ImportList) AssociateByShortName() map[string]*Import {
	return xslices.ToMap[[]*Import, map[string]*Import](ii, func(t *Import) string {
		return t.ShortName
	}, nil)
}

type Import struct {
	ShortName string
	Path      string
	Aliased   bool
	Source    *Source
}
