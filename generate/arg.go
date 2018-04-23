package generate

import "go/ast"

type Arg struct {
	Name, AsField *ast.Ident
	Typ           ast.Expr
	IsStar        bool //是否是指针
	IsError       bool
	IsBody        bool
}

func (a Arg) Field(scope *ast.Scope) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{a.chooseName(scope)},
		Type:  a.Typ,
	}
}

func (a Arg) Result() *ast.Field {
	return &ast.Field{
		Names: nil,
		Type:  a.Typ,
	}
}

func (a Arg) chooseName(scope *ast.Scope) *ast.Ident {
	if a.Name == nil || scope.Lookup(a.Name.Name) != nil {
		return InventName(a.Typ, scope)
	}
	return a.Name
}

func (a Arg) Exported() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(export(a.AsField.Name))},
		Type:  a.Typ,
	}
}
