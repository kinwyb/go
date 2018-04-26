package generate

import (
	"go/ast"
	"strings"
	"unicode"
)

type Iface struct {
	Name, StubName, RcvrName *ast.Ident
	Methods                  []Method
	Prefix                   string
	Comments                 *ast.CommentGroup
}

func (i Iface) Reciever() *ast.Field {
	return Field(i.receiverName(), i.StubName)
}

func (i Iface) StubStructDecl(root *ast.File) ast.Decl {
	ds := StructDecl(i.StubName, &ast.FieldList{
		List: []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent("serv")},
			Type:  ast.NewIdent("endPoints." + i.StubName.Name),
		}},
	})
	if root != nil {
		root.Decls = append(root.Decls, ds)
	}
	return ds
}

func (i Iface) receiverName() *ast.Ident {
	if i.RcvrName != nil {
		return i.RcvrName
	}
	scope := ast.NewScope(nil)
	for _, meth := range i.Methods {
		for _, arg := range meth.Params {
			if arg.Name != nil {
				scope.Insert(ast.NewObj(ast.Var, arg.Name.Name))
			}
		}
		for _, arg := range meth.Results {
			if arg.Name != nil {
				scope.Insert(ast.NewObj(ast.Var, arg.Name.Name))
			}
		}
	}
	return ast.NewIdent(unexport(InventName(i.Name, scope).Name))
}

func unexport(s string) string {
	first := true
	return strings.Map(func(r rune) rune {
		if first {
			first = false
			return unicode.ToLower(r)
		}
		return r
	}, s)
}
