package codegen

import (
	"github.com/eugenenosenko/gopoly/internal/code"
)

type Task struct {
	Filename string
	Template string
	Data     *Data
}

type Data struct {
	Package string
	Imports []*code.Import
	Types   []*Type
}

type Type struct {
	Name               string
	Variants           map[string]*code.Variant
	DecodingStrategy   string
	DiscriminatorField string
}
