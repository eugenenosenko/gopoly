package code

import (
	"fmt"
	"go/build"
	"path"

	"github.com/eugenenosenko/gopoly/internal/xslices"
)

// DecodingStrategy represents a specific decoding algorithm that will be used
// when unmarshaling polymorphic structure, it can be either DecodingStrategyStrict or DecodingStrategyDiscriminator.
type DecodingStrategy string

func (s DecodingStrategy) String() string {
	return string(s)
}

func (s DecodingStrategy) IsDiscriminator() bool {
	return s == DecodingStrategyDiscriminator
}

func (s DecodingStrategy) IsStrict() bool {
	return s == DecodingStrategyStrict
}

func (s DecodingStrategy) IsValid() bool {
	switch s {
	case DecodingStrategyStrict, DecodingStrategyDiscriminator:
		return true
	default:
		return false
	}
}

const (
	// DecodingStrategyStrict is the decoding algorithm that will try to match 1:1 fields on JSON payload and on struct and
	// if JSON payload does contain a field that is not present on structure, it will return error.
	DecodingStrategyStrict = DecodingStrategy("strict")

	// DecodingStrategyDiscriminator is the decoding algorithm that will try to unmarshal payload based
	// on discriminator field value.
	DecodingStrategyDiscriminator = DecodingStrategy("discriminator")
)

const (
	KindScalar = iota
	KindMap
	KindSlice
)

// Task represent a code-gen task that Generator accepts.
// It's not advisable to have multiple tasks with the same Task.Filename since the Generator will just
// overwrite the data in it.
type Task interface {
	// Filename is the task's output file
	Filename() string

	// Template string that is then parsed into template.Template
	Template() string

	// Data represents Types, Interfaces, and Variants from a specific package,
	// that need to be generated.
	Data() Data
}

// Data represents Type, Interfaces, and Variants from a specific package.
type Data interface {
	Package() Package
	Imports() []Import
	Types() []Type
}

// Type represents an interface for which unmarshal method needs to be created.
type Type interface {
	Name() string
	Variants() map[string]Variant
	MarkerMethod() string
	DecodingStrategy() DecodingStrategy
	DiscriminatorField() string
	Package() Package
}

// PolyField is the representation of a polymorphic field (read that has interface type)
// if you have the following structure:
//
//	type A interface {
//		IsA()
//	}
//
//	type Payload struct {
//		Value A `json: "value"` // AField is a PolyField
//	}
type PolyField interface {
	Name() string
	Tags() string
	Interface() Interface
	Kind() int
	Prefix() string
}

// PolyFieldList
type PolyFieldList []PolyField

// Variant
type Variant interface {
	Name() string
	Fields() PolyFieldList
	Interface() Interface
}

// Interface
type Interface interface {
	Name() string
	MarkerMethod() string
	Variants() VariantList
	Pkg() Package
}

type InterfaceList []Interface

// Import
type Import interface {
	ShortName() string
	Path() string
	Aliased() bool
	Source() Source
}

type ImportList []Import

// AssociateByShortName
func (ii ImportList) AssociateByShortName() map[string]Import {
	return xslices.ToMap[[]Import, map[string]Import](ii, func(t Import) string {
		return t.ShortName()
	}, nil)
}

// Source
type Source interface {
	Package() Package
	Interfaces() InterfaceList
	Imports() ImportList
}

// Package
type Package string

func (p Package) String() string {
	return string(p)
}

// Path
func (p Package) Path() string {
	return string(p)
}

// Name
func (p Package) Name() string {
	return path.Base(p.Path())
}

// Dir
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

// VariantList
type VariantList []Variant

func (vvs VariantList) AssociateByVariantName() map[string]Variant {
	res := make(map[string]Variant, 0)
	for _, variant := range vvs {
		res[variant.Name()] = variant
	}
	return res
}

// SourceList
type SourceList []Source

// AssociateByPkgName
func (ss SourceList) AssociateByPkgName() map[Package]Source {
	res := make(map[Package]Source, 0)
	for _, s := range ss {
		res[s.Package()] = s
	}
	return res
}

// AssociateByName
func (ii InterfaceList) AssociateByName() map[string]Interface {
	return xslices.ToMap[[]Interface, map[string]Interface](ii, func(t Interface) string {
		return t.Name()
	}, nil)
}

var _ fmt.Stringer = (*DecodingStrategy)(nil)
