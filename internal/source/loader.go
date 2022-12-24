package source

import (
	"context"
	"fmt"
	"go/ast"
	"path"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/eugenenosenko/gopoly/internal/code"
	"github.com/eugenenosenko/gopoly/internal/config"
	"github.com/eugenenosenko/gopoly/internal/xslices"
	"github.com/eugenenosenko/gopoly/internal/xtypes"
)

type PkgPath = string

type Loader interface {
	Load(c context.Context, path []*config.TypeDefinition) (code.SourceList, error)
}

type loader struct {
	loadFunc func(paths []string) ([]*Package, error)
	logf     func(format string, args ...any)
}

type Config struct {
	LoadFunc func(paths []string) ([]*Package, error)
	Logf     func(format string, args ...any)
}

type Definition struct {
	Dec     ast.Decl
	PkgName string
	PkgPath string
}

func NewLoader(c *Config) (*loader, error) {
	if c == nil {
		return nil, errors.New("source.NewLoader: config is nil")
	}
	return &loader{loadFunc: c.LoadFunc, logf: c.Logf}, nil
}

func (l *loader) Load(ctx context.Context, c []*config.TypeDefinition) (code.SourceList, error) {
	packages := xslices.ToSetFunc[[]*config.TypeDefinition, map[string]struct{}](
		c,
		func(t *config.TypeDefinition) string { return t.Package },
	)
	packageNames := maps.Keys(packages)

	files, err := l.loadFunc(packageNames)
	if err != nil {
		return nil, errors.Wrapf(err, "collecting ast.Files from %v", packageNames)
	}

	pdecs := make(map[PkgPath][]*Definition, 0)
	for d := range declarations(ctx, files) {
		pdecs[d.PkgPath] = append(pdecs[d.PkgPath], d)
	}

	ifaces, err := interfaces(pdecs, c)
	if err != nil {
		return nil, errors.Wrapf(err, "collecting interfaces from %v", packageNames)
	}

	psources := make(map[PkgPath]*code.Source, 0)
	for pkg, i := range ifaces {
		psources[pkg] = &code.Source{Package: code.Package(pkg), Interfaces: i}
	}

	for pkg, decs := range pdecs {
		source := psources[pkg]
		for _, dec := range decs {
			d, ok := dec.Dec.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range d.Specs {
				i, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}
				// imports are quoted, i.e. "time", "bytes"
				importpath, err := strconv.Unquote(i.Path.Value)
				if err != nil {
					return nil, err
				}
				// check if imported types are in the supplied config
				if _, present := packages[importpath]; present {
					var aliased bool
					sname := path.Base(importpath)
					if alias := i.Name; alias != nil {
						sname = alias.Name
						aliased = true
					}
					source.Imports = append(source.Imports, &code.Import{
						Source:    psources[importpath],
						ShortName: sname,
						Path:      importpath,
						Aliased:   aliased,
					})
				}
			}
		}
	}

	for pkg, decs := range pdecs {
		var types = xslices.Flatten(xslices.Map[[]*code.Interface, [][]*code.Variant](
			psources[pkg].Interfaces,
			func(i *code.Interface) []*code.Variant { return i.Variants },
		))
		variants := code.VariantList(types).AssociateByVariantName()
		pifaces := ifaces[pkg].AssociateByName()
		imports := psources[pkg].Imports.AssociateByShortName()

		for _, dec := range decs {
			d, ok := dec.Dec.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range d.Specs {
				// look for type declarations
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// filter out only structs
				strct, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				variant, ok := variants[ts.Name.Name]
				if !ok {
					continue
				}

				pfields := make([]*code.PolyField, 0)
				for _, field := range strct.Fields.List {
					var iface *code.Interface
					var kind int
					var prefix string
					if st, ok := field.Type.(*ast.Ident); ok { // scalar type field
						iface, kind = pifaces[st.Name], code.KindScalar
					} else if imp, ok := field.Type.(*ast.SelectorExpr); ok { // imported scalar type field
						iface, kind = matchImportedPFIface(imp, imports), code.KindScalar
						prefix = importPrefix(imp.X)
					}

					if arr, ok := field.Type.(*ast.ArrayType); ok {
						// imported slice type
						if imp, ok := arr.Elt.(*ast.SelectorExpr); ok {
							iface, kind = matchImportedPFIface(imp, imports), code.KindSlice
							prefix = importPrefix(imp.X)
						} else if st, ok := arr.Elt.(*ast.Ident); ok {
							iface, kind = pifaces[st.Name], code.KindSlice
						}
					}

					if mp, ok := field.Type.(*ast.MapType); ok {
						// imported map field-type
						if imp, ok := mp.Value.(*ast.SelectorExpr); ok {
							iface, kind = matchImportedPFIface(imp, imports), code.KindMap
							prefix = importPrefix(imp.X)
						} else if st, ok := mp.Value.(*ast.Ident); ok {
							iface, kind = pifaces[st.Name], code.KindMap
						}
					}

					if iface != nil {
						pfields = append(pfields, &code.PolyField{
							Name:      field.Names[0].Name,
							Tags:      field.Tag.Value,
							Interface: iface,
							Kind:      kind,
							Prefix:    prefix,
						})
					}
				}
				variant.Fields = pfields
			}
		}
	}

	return maps.Values(psources), nil
}

func importPrefix(e ast.Expr) string {
	if ident, ok := e.(*ast.Ident); ok {
		return ident.Name // short name or alias
	}
	return ""
}

func matchImportedPFIface(se *ast.SelectorExpr, imports map[string]*code.Import) *code.Interface {
	name, sname := se.Sel.Name, importPrefix(se.X)
	i, ok := imports[sname]
	if !ok {
		return nil
	}
	nifaces := i.Source.Interfaces.AssociateByName()
	iface, ok := nifaces[name]
	if !ok {
		return nil
	}
	return iface
}

// interfaces iterates decs and tries to match ifaces that were passed with the config
func interfaces(defs map[PkgPath][]*Definition, tts []*config.TypeDefinition) (map[PkgPath]code.InterfaceList, error) {
	if err := validateDefinitions(defs, tts); err != nil {
		return nil, errors.Wrap(err, "validating tts definitions")
	}
	pinterfaces := make(map[PkgPath]code.InterfaceList, 0)
	for _, t := range tts {
		i := &code.Interface{
			Name:         t.Type,
			MarkerMethod: t.MarkerMethod,
			Variants:     []*code.Variant{},
		}
		expected := xslices.ToSet[[]string](maps.Values(t.Discriminator.Mapping))

		for _, dec := range defs[t.Package] {
			// check only functions and methods
			fun, ok := dec.Dec.(*ast.FuncDecl)
			if !ok {
				continue
			}
			// check whether function name matches marker-method
			if t.MarkerMethod != fun.Name.Name {
				continue
			}
			// get the receiver field
			receiver, ok := xslices.First(fun.Recv.List)
			if !ok {
				continue
			}
			// fetch the name of the receiver type
			var name string
			if ident, ok := receiver.Type.(*ast.Ident); ok {
				name = ident.Name
			}
			// if discriminator check whether the receiver type is defined as a variant of the i-face
			if t.DecodingStrategy.IsDiscriminator() {
				if _, ok := expected[name]; ok {
					i.Variants = append(i.Variants, &code.Variant{Name: name, Interface: i})
				}
			} else {
				i.Variants = append(i.Variants, &code.Variant{Name: name, Interface: i})
			}
		}

		if t.DecodingStrategy.IsStrict() {
			pinterfaces[t.Package] = append(pinterfaces[t.Package], i) // just add all variants found
		} else {
			err := validateNoMissingVariants(i.Variants, expected) // validate if all were found
			if err != nil {
				return nil, errors.Wrapf(err, "mathing variants for '%s.%s'", t.Package, t.Type)
			}
			pinterfaces[t.Package] = append(pinterfaces[t.Package], i)
		}
	}
	return pinterfaces, nil
}

func validateNoMissingVariants(vars code.VariantList, expected map[string]struct{}) error {
	variants := xslices.Map[[]*code.Variant, []string](
		vars,
		func(t *code.Variant) string { return t.Name },
	)
	// check if all variants are accounted for
	if diff := xslices.Difference(variants, maps.Keys(expected)); len(diff) > 0 {
		return fmt.Errorf("failed to match following %v variants", diff)
	}
	return nil
}

func validateDefinitions(pdefs map[PkgPath][]*Definition, types []*config.TypeDefinition) error {
	provided := xslices.ToMap[[]*config.TypeDefinition, map[string]*config.TypeDefinition](
		types, func(t *config.TypeDefinition) string {
			return (&xtypes.Tuple[string, string]{First: t.Type, Second: t.Package}).String()
		}, nil,
	)

	ifaces := make(map[string]*xtypes.Tuple[*config.TypeDefinition, *ast.InterfaceType])
	for pkg, defs := range pdefs {
		for _, def := range defs {
			d, ok := def.Dec.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// we are only interested in interfaces at this point
				i, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				// composite map key Type.Name + PkgPath
				key := (&xtypes.Tuple[string, string]{First: ts.Name.Name, Second: pkg}).String()
				// check if this interface is among those provided by the configuration
				if t, present := provided[key]; present {
					ifaces[key] = &xtypes.Tuple[*config.TypeDefinition, *ast.InterfaceType]{First: t, Second: i}
				}
			}
		}
	}

	// check if all interfaces are accounted for
	if diff := xslices.Difference(maps.Keys(ifaces), maps.Keys(provided)); len(diff) > 0 {
		return fmt.Errorf("failed to locate following interfaces: %v", diff)
	}

	// validate signatures on the marker methods
	for _, v := range ifaces {
		found := false
		for _, field := range v.Second.Methods.List {
			if f, ok := field.Type.(*ast.FuncType); ok {
				if len(f.Params.List) != 0 || f.Results != nil {
					continue
				}
			}
			name, ok := xslices.First(field.Names)
			if !ok {
				continue
			}
			if name.Name == v.First.MarkerMethod {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("interface '%s.%s' missing zero arg/returns marker-method '%s'",
				v.First.Package,
				v.First.Type,
				v.First.MarkerMethod,
			)
		}
	}
	return nil
}

// declarations goes through the supplied packages and collects all ast-declarations
// for each package a new goroutine is started which starts a new goroutine for each file
// in the package, returns a channel of Definition.
func declarations(ctx context.Context, packs []*Package) <-chan *Definition {
	c := make(chan *Definition)
	go func() {
		defer close(c)
		var wgOuter sync.WaitGroup
		wgOuter.Add(len(packs))

		for _, pack := range packs {
			go func(p *Package) {
				defer wgOuter.Done()

				var wgInner sync.WaitGroup
				wgInner.Add(len(p.Files))

				for _, file := range p.Files {
					go func(f *ast.File, p *Package) {
						defer wgInner.Done()

						for _, d := range f.Decls {
							select {
							case <-ctx.Done():
								return
							case c <- &Definition{Dec: d, PkgName: p.Name, PkgPath: p.Path}:
							}
						}
					}(file, p)
				}
				wgInner.Wait()
			}(pack)
		}
		wgOuter.Wait()
	}()

	return c
}

var _ Loader = (*loader)(nil)
