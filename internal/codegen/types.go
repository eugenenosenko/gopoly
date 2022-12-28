package codegen

import (
	"github.com/eugenenosenko/gopoly/internal/code"
)

// Task represent a code-gen task that Generator accepts.
// It's not advisable to have multiple tasks with the same Task.Filename since the Generator will just
// overwrite the data in it.
type Task struct {
	// Filename output file
	Filename string

	// Template string that is then parsed into template.Template
	Template string

	// Data represents Types, Interfaces, and Variants from a specific package,
	// that need to be generated.
	Data *Data
}

// Data represents Types, Interfaces, and Variants from a specific package,
type Data struct {
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
