package rpcx

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/kinwyb/go/generate"
)

//生成服务代码
func addMethodService(root *ast.File, ifc *generate.Iface, m *generate.Method) {
	notImpl := generate.FetchFuncDecl("Example")
	notImpl.Name = m.Name
	notImpl.Recv = generate.FieldList(ifc.Reciever())
	//生成请求参数
	parms := &ast.FieldList{}
	parms.List = []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("ctx")},
		Type:  generate.Sel(ast.NewIdent("context"), ast.NewIdent("Context")),
	}, {
		Names: []*ast.Ident{ast.NewIdent("arg")},
		Type: &ast.StarExpr{
			X: m.RequestStructName(),
		},
	}, {
		Names: []*ast.Ident{ast.NewIdent("resp")},
		Type: &ast.StarExpr{
			X: m.ResponseStructName(),
		},
	}}
	notImpl.Type.Params = parms
	notImpl.Type.Results = &ast.FieldList{
		List: []*ast.Field{{
			Type: ast.NewIdent("error"),
		}},
	}
	ret := notImpl.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
	ret.Fun.(*ast.Ident).Name = notImpl.Recv.List[0].Names[0].Name + ".serv." + m.Name.Name
	var args []ast.Expr
	for _, v := range m.Params {
		args = append(args, &ast.BasicLit{
			Kind:  token.TYPE,
			Value: "arg." + strings.Title(v.Name.Name),
		})
	}
	dbCall := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: "db.TagConnect",
		},
		Args: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.TYPE,
				Value: "arg.Dbtag",
			},
		},
	}
	db := ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: "q",
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{dbCall},
	}
	ret.Args = args
	responses := m.ResponseStructFields()
	r := ast.AssignStmt{
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ret},
	}
	for _, v := range responses.List {
		if len(v.Names) > 0 {
			v.Names[0].Name = "resp." + v.Names[0].Name
			r.Lhs = append(r.Lhs, v.Names[0])
		}
	}
	ret.Args[len(ret.Args)-1] = &ast.BasicLit{
		Kind:  token.TYPE,
		Value: "q",
	}
	notImpl.Body.List[0] = &db
	notImpl.Body.List = append(notImpl.Body.List, &r)
	rnil := ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	}
	notImpl.Body.List = append(notImpl.Body.List, &rnil)
	root.Decls = append(root.Decls, notImpl)
}
