package generate

import (
	"go/ast"
)

type Method struct {
	Prefix          string
	Name            *ast.Ident
	Params          []Arg
	Results         []Arg
	Comments        []string
	StructsResolved bool
}

func (m Method) FuncParams(scope *ast.Scope) *ast.FieldList {
	parms := &ast.FieldList{}
	if m.hasContext() {
		parms.List = []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent("ctx")},
			Type:  Sel(ast.NewIdent("context"), ast.NewIdent("Context")),
		}}
		scope.Insert(ast.NewObj(ast.Var, "ctx"))
	}
	parms.List = append(parms.List, MappedFieldList(func(a Arg) *ast.Field {
		return a.Field(scope)
	}, m.nonContextParams()...).List...)
	return parms
}

func (m Method) FuncResults() *ast.FieldList {
	return MappedFieldList(func(a Arg) *ast.Field {
		return a.Result()
	}, m.Results...)
}

func (m Method) hasContext() bool {
	if len(m.Params) < 1 {
		return false
	}
	carg := m.Params[0].Typ
	// ugh. this is maybe okay for the one-off, but a general case for matching
	// types would be helpful
	if sel, is := carg.(*ast.SelectorExpr); is && sel.Sel.Name == "Context" {
		if id, is := sel.X.(*ast.Ident); is && id.Name == "context" {
			return true
		}
	}
	return false
}

func (m Method) nonContextParams() []Arg {
	if m.hasContext() {
		return m.Params[1:]
	}
	return m.Params
}

func Sel(ids ...*ast.Ident) ast.Expr {
	switch len(ids) {
	default:
		return &ast.SelectorExpr{
			X:   Sel(ids[:len(ids)-1]...),
			Sel: ids[len(ids)-1],
		}
	case 1:
		return ids[0]
	case 0:
		panic("zero ids to Sel()")
	}
}

func (m Method) resolveStructNames() {
	if m.StructsResolved {
		return
	}
	m.StructsResolved = true
	scope := ast.NewScope(nil)
	for i, p := range m.Params {
		p.AsField = p.chooseName(scope)
		m.Params[i] = p
	}
	scope = ast.NewScope(nil)
	for i, r := range m.Results {
		r.AsField = r.chooseName(scope)
		m.Results[i] = r
	}
}

func (m Method) RequestStruct() ast.Decl {
	m.resolveStructNames()
	if len(m.Params) == 1 {
		return nil
	}
	name, _ := m.RequestStructName()
	return StructDecl(name, m.requestStructFields())
}

func (m Method) ResponseStruct() ast.Decl {
	m.resolveStructNames()
	return StructDecl(m.ResponseStructName(), m.ResponseStructFields())
}

func (m Method) RequestStructName() (*ast.Ident, ast.Expr) {
	if len(m.Params) == 1 {
		return ast.NewIdent(m.Params[0].Name.Name), m.Params[0].Exported().Type
	}
	ret := ast.NewIdent(m.Prefix + export(m.Name.Name) + "Request")
	return ret, &ast.StarExpr{
		X: ret,
	}
}

func (m Method) requestStructFields() *ast.FieldList {
	return MappedFieldList(func(a Arg) *ast.Field {
		return a.Exported()
	}, m.nonContextParams()...)
}

func (m Method) ResponseStructName() *ast.Ident {
	return ast.NewIdent(m.Prefix + export(m.Name.Name) + "Response")
}

func (m Method) ResponseStructFields() *ast.FieldList {
	return MappedFieldList(func(a Arg) *ast.Field {
		return a.Exported()
	}, m.Results...)
}

func (m Method) WrapRequest() ast.Expr {
	var kvs []ast.Expr
	m.resolveStructNames()
	for _, a := range m.Params {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(export(a.AsField.Name)),
			Value: ast.NewIdent(export(a.Name.Name)),
		})
	}
	name, _ := m.RequestStructName()
	return &ast.CompositeLit{
		Type: name,
		Elts: kvs,
	}
}
