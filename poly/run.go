package poly

import (
	"context"
	"path"

	"github.com/pkg/errors"

	"github.com/eugenenosenko/gopoly/code"
	"github.com/eugenenosenko/gopoly/codegen"
	"github.com/eugenenosenko/gopoly/config"
	"github.com/eugenenosenko/gopoly/internal/xmaps"
	"github.com/eugenenosenko/gopoly/templates"
)

// Run is the entry point of the application. Takes context.Context and config.Config and uses that to generate
// the tasks for the codegen.Generator.
//
// Codegen files will be created in the corresponding packages that is because if your config.TypesList contain
// polymorphic fields, codegen.Generator needs to create custom UnmarshalJSON methods for those.
//
// If multiple output files have been provided and the types are located in the same package, then declarations
// will be split between those files. If multiple types share the same output file but are located in the same package
// declarations will be created in the same file.
func (r *Client) Run(ctx context.Context, c *config.Config) error {
	if len(c.Types) == 0 {
		r.Logf("No types provided or configuration is incorrect.")
		return nil
	}

	sources, err := r.Loader.Load(ctx, c.Types)
	if err != nil {
		return errors.Wrapf(err, "loading source from packages")
	}

	// each definition needs to be generated in its own package
	// if definition has a separate output filename defined then output should go there
	tasks := make([]*codegen.Task, 0)
	psources := sources.AssociateByPkgName()
	for filename, types := range c.Types.AssociateByOutput() {
		for pkg, tts := range types.AssociateByPkgName() {
			p := code.Package(pkg)
			src := psources[p]
			d := codegen.Input{
				Package: p.Name(),
				Imports: src.Imports,
			}

			ntype := tts.AssociateByTypeName()
			for _, iface := range src.Interfaces {
				def, ok := ntype[iface.Name]
				if !ok {
					continue
				}

				variants := make(map[string]*code.Variant, 0)
				if def.DecodingStrategy.IsDiscriminator() {
					nvars := iface.Variants.AssociateByVariantName()
					for v, name := range def.Discriminator.Mapping { // iterate over discriminator mappings
						variants[v] = nvars[name]
					}
				} else {
					variants = xmaps.Merge(variants, iface.Variants.AssociateByVariantName())
				}
				// add type to to-be-generated data with its variants
				d.Types = append(d.Types, &codegen.Type{
					Name:               iface.Name,
					Variants:           variants,
					DecodingStrategy:   def.DecodingStrategy.String(),
					DiscriminatorField: def.Discriminator.Field,
				})
			}
			tasks = append(tasks, &codegen.Task{
				Filename: outputFilename(p, filename),
				Template: templates.DefaultJSONTemplate(),
				Input:    &d,
			})
		}
	}

	for _, task := range tasks {
		if err = r.Generator.Generate(task); err != nil {
			return errors.Wrapf(err, "generating codegen")
		}
	}
	return nil
}

func outputFilename(pkg code.Package, filename string) string {
	return path.Join(pkg.Dir(), path.Base(filename))
}
