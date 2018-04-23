package generate

import (
	"fmt"
	"go/ast"
)

type SourceContext struct {
	Pkg        *ast.Ident
	Imports    []*ast.ImportSpec
	Interfaces []Iface
	Types      []*ast.TypeSpec
	Prefix     string
	Common     []*ast.Comment
}

func (sc *SourceContext) validate() error {
	if len(sc.Interfaces) != 1 {
		return fmt.Errorf("found %d interfaces, expecting exactly 1", len(sc.Interfaces))
	}
	for _, i := range sc.Interfaces {
		for _, m := range i.Methods {
			if len(m.Results) < 1 {
				return fmt.Errorf("Method %q of interface %q has no result types", m.Name, i.Name)
			}
		}
	}
	return nil
}

func (sc *SourceContext) ImportDecls(root *ast.File) (decls []ast.Decl) {
	have := map[string]struct{}{}
	notHave := func(is *ast.ImportSpec) bool {
		if _, has := have[is.Path.Value]; has {
			return false
		}
		have[is.Path.Value] = struct{}{}
		return true
	}
	for _, is := range sc.Imports {
		if notHave(is) {
			decls = append(decls, importFor(is))
		}
	}
	if root != nil {
		root.Decls = append(root.Decls, decls...)
	}
	return
}
