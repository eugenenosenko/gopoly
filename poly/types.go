package poly

import (
	"github.com/eugenenosenko/gopoly/code"
)

type data struct {
	pkg     code.Package
	imports []code.Import
	types   []code.Type
}

func (d *data) Package() code.Package  { return d.pkg }
func (d *data) Imports() []code.Import { return d.imports }
func (d *data) Types() []code.Type     { return d.types }

var _ code.Data = (*data)(nil)

type iType struct {
	name               string
	variants           map[string]code.Variant
	markerMethod       string
	decodingStrategy   code.DecodingStrategy
	discriminatorField string
	pkg                code.Package
}

func (i *iType) Name() string                            { return i.name }
func (i *iType) Variants() map[string]code.Variant       { return i.variants }
func (i *iType) MarkerMethod() string                    { return i.markerMethod }
func (i *iType) DecodingStrategy() code.DecodingStrategy { return i.decodingStrategy }
func (i *iType) DiscriminatorField() string              { return i.discriminatorField }
func (i *iType) Package() code.Package                   { return i.pkg }

var _ code.Type = (*iType)(nil)

type codegenTask struct {
	filename string
	template string
	data     code.Data
}

func (c *codegenTask) Filename() string { return c.filename }
func (c *codegenTask) Template() string { return c.template }
func (c *codegenTask) Data() code.Data  { return c.data }

var _ code.Task = (*codegenTask)(nil)
