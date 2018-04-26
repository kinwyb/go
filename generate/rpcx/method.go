package rpcx

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/kinwyb/go/generate"
)

//生成服务代码
func addMethodService(root *ast.File, ifc *generate.Iface, m *generate.Method, structname string) {
	notImpl := generate.FetchFuncDecl("Example")
	notImpl.Name = m.Name
	notImpl.Recv = generate.FieldList(ifc.Reciever())
	notImpl.Recv.List[0].Type = &ast.StarExpr{
		X: &ast.Ident{Name: structname},
	}
	//生成请求参数
	parms := &ast.FieldList{}
	_, resultTp := m.ResponseStructName()
	_, paramTp := m.RequestStructName()
	parms.List = []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("ctx")},
		Type:  generate.Sel(ast.NewIdent("context"), ast.NewIdent("Context")),
	}, {
		Names: []*ast.Ident{ast.NewIdent("arg")},
		Type:  paramTp,
	}, {
		Names: []*ast.Ident{ast.NewIdent("resp")},
		Type:  resultTp,
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
	if len(m.Params) == 1 {
		args = append(args, &ast.BasicLit{
			Kind:  token.TYPE,
			Value: "arg",
		})
	} else {
		for _, v := range m.Params {
			args = append(args, &ast.BasicLit{
				Kind:  token.TYPE,
				Value: "arg." + strings.Title(v.Name.Name),
			})
		}
	}
	ret.Args = args
	responses := m.ResponseStructFields()
	if len(responses.List) == 1 {
		responses.List[0].Names = []*ast.Ident{
			ast.NewIdent("resp"),
		}
	} else {
		for i, v := range responses.List {
			responses.List[i].Names = []*ast.Ident{
				ast.NewIdent("resp." + v.Names[0].Name),
			}
		}
	}
	r := ast.AssignStmt{
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ret},
	}
	for _, v := range responses.List {
		if len(v.Names) > 0 {
			r.Lhs = append(r.Lhs, v.Names[0])
		}
	}
	notImpl.Body.List[0] = &r
	rnil := ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	}
	notImpl.Body.List = append(notImpl.Body.List, &rnil)
	root.Decls = append(root.Decls, notImpl)
}
