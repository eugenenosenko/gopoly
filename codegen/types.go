package codegen

import (
	"github.com/eugenenosenko/gopoly/code"
)

// Task represent a code-gen task that Generator accepts.
// It's not advisable to have multiple tasks with the same Task.Filename since the Generator will just
// overwrite the data in it.
type Task struct {
	// Filename output file
	Filename string

	// Template string that is then parsed into template.Template
	Template string

	// Input represents Types, Interfaces, and Variants from a specific package,
	// that need to be generated.
	Input *Input
}

// Input represents Types, Interfaces, and Variants from a specific package,
type Input struct {
	Package string
	Imports []*code.Import
	Types   []*Type
}

// Type represents an interface for which unmarshal method needs to be created
type Type struct {
	Name               string
	Variants           map[string]*code.Variant
	DecodingStrategy   string
	DiscriminatorField string
}
