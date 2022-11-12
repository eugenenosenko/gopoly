package source

import (
	"context"
	"fmt"
	"go/ast"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/eugenenosenko/gopoly/internal/config"
	"github.com/eugenenosenko/gopoly/internal/xslices"
	"github.com/eugenenosenko/gopoly/internal/xtypes"
)

type Source struct {
	Interfaces []any
}

type Loader interface {
	Load(c context.Context, path []*config.TypeInfo) (*Source, error)
}

type loader struct {
	loadFunc func(paths []string) ([]*Package, error)
	logf     func(format string, args ...any)
}

func NewLoader(c *Config) (*loader, error) {
	if c == nil {
		return nil, errors.New("source.NewLoader: config is nil")
	}
	return &loader{loadFunc: c.LoadFunc, logf: c.Logf}, nil
}

func (l *loader) Load(ctx context.Context, c []*config.TypeInfo) (*Source, error) {
	packages := xslices.Map[[]*config.TypeInfo, []string](c, func(t *config.TypeInfo) string {
		return t.PackagePath
	})
	packages = xslices.Distinct(packages)

	files, err := l.loadFunc(packages)
	if err != nil {
		return nil, errors.Wrapf(err, "collecting ast.Files from %v", packages)
	}

	// var stypes StructTypeMap
	// var fdecls FuncDecMap

	// interface matched
	// fetch codegen configuration for it
	// validate
	// look up the
	decs := make([]*DecInfo, 0, len(files)*10)
	for d := range declarations(ctx, files) {
		decs = append(decs, d)
	}

	_, err = interfaces(decs, c)
	if err != nil {
		return nil, err
	}
	// for _, _ := range decs {
	// 	fmt.Println(ifaces)
	// }

	return nil, nil
}

// interfaces iterates decs and tries to match ifaces that were passed with the config
func interfaces(decs []*DecInfo, c []*config.TypeInfo) (InterfaceTypeMap, error) {
	types := xslices.ToMap[[]*config.TypeInfo, map[string]*config.TypeInfo](
		c,
		func(t *config.TypeInfo) string {
			return (&xtypes.Tuple[string, string]{
				First:  t.Type,
				Second: t.PackagePath,
			}).String()
		},
		nil,
	)

	ifaces := make(map[string]Entry[*ast.InterfaceType])
	for _, dec := range decs {
		if d, ok := dec.Dec.(*ast.GenDecl); ok {
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					// we are only interested in interfaces at this point
					if i, ok := ts.Type.(*ast.InterfaceType); ok {
						// check if this interface is among those provided by the configuration
						key := (&xtypes.Tuple[string, string]{
							First:  ts.Name.Name,
							Second: dec.PkgPath,
						}).String()
						if t, present := types[key]; present {
							ifaces[key] = Entry[*ast.InterfaceType]{Type: t, Dec: i}
						}
					}
				}
			}
		}
	}

	if diff := xslices.Difference(maps.Keys(ifaces), maps.Keys(types)); len(diff) > 0 {
		return nil, fmt.Errorf("failed to locate following interfaces: %v", diff)
	}

	for k, v := range ifaces {
		found := false
		for _, field := range v.Dec.Methods.List {
			if f, ok := field.Type.(*ast.FuncType); ok {
				if len(f.Params.List) != 0 || f.Results != nil {
					continue
				}
			}
			name, ok := xslices.First(field.Names)
			if !ok {
				continue
			}
			if name.Name == v.Type.MarkerMethod {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("interface '%s.%s' missing zero arg/returns marker-method '%s'", v.Type.PackagePath, k, v.Type.MarkerMethod)
		}
	}

	variants := make(map[string][]string)
	for _, t := range c {
		var expected []string
		if t.DecodingStrategy.IsDiscriminator() {
			expected = maps.Values(t.Discriminator.Mapping)
		}

		for _, dec := range decs {
			if fun, ok := dec.Dec.(*ast.FuncDecl); ok {
				if t.MarkerMethod != fun.Name.Name {
					continue
				}
				if dec.PkgPath != t.PackagePath {
					continue
				}
				receiver, ok := xslices.First(fun.Recv.List)
				if !ok {
					continue
				}

				var name string
				if ident, ok := receiver.Type.(*ast.Ident); ok {
					name = ident.Name
				}
				if t.DecodingStrategy.IsDiscriminator() && slices.Contains(expected, name) {
					variants[t.MarkerMethod] = append(variants[t.MarkerMethod], name)
				} else {
					variants[t.MarkerMethod] = append(variants[t.MarkerMethod], name)
				}
			}
		}

		// validate that all variants are accounted for
		if t.DecodingStrategy.IsDiscriminator() {
			if missing := xslices.Difference(variants[t.MarkerMethod], expected); len(missing) > 0 {
				return nil, fmt.Errorf("failed to match following %v variants of type %s in package %s", missing, t.Type, t.PackagePath)
			}
		}
	}

	return nil, nil
}

func declarations(ctx context.Context, packs []*Package) <-chan *DecInfo {
	c := make(chan *DecInfo)
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
							case c <- &DecInfo{
								Dec:     d,
								PkgName: p.Name,
								PkgPath: p.Path,
								Imports: xslices.Map[[]*ast.ImportSpec, []string](
									f.Imports,
									func(t *ast.ImportSpec) string { return t.Path.Value },
								),
							}:
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
