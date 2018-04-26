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
	resultTp := m.ResponseStructName()
	_, paramTp := m.RequestStructName()
	parms.List = []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("ctx")},
		Type:  generate.Sel(ast.NewIdent("context"), ast.NewIdent("Context")),
	}, {
		Names: []*ast.Ident{ast.NewIdent("arg")},
		Type:  paramTp,
	}, {
		Names: []*ast.Ident{ast.NewIdent("resp")},
		Type: &ast.StarExpr{
			X: resultTp,
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
	r := ast.AssignStmt{
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ret},
	}
	for _, v := range responses.List {
		if len(v.Names) > 0 {
			r.Lhs = append(r.Lhs, ast.NewIdent("resp."+v.Names[0].Name))
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

//生成客户代码
func addMethodClient(root *ast.File, ifc *generate.Iface, m *generate.Method, structname string) {
	notImpl := generate.FetchFuncDecl("Example")
	notImpl.Name = m.Name
	notImpl.Recv = generate.FieldList(ifc.Reciever())
	notImpl.Recv.List[0].Type = &ast.StarExpr{
		X: &ast.Ident{Name: structname},
	}
	//生成请求参数
	notImpl.Type.Params = generate.MappedFieldList(func(a generate.Arg) *ast.Field {
		return a.Exported()
	}, m.Params...)
	result := generate.MappedFieldList(func(a generate.Arg) *ast.Field {
		return a.Exported()
	}, m.Results...)
	var returnValue []ast.Expr
	for _, v := range result.List {
		if len(v.Names) > 0 {
			returnValue = append(returnValue, ast.NewIdent("reply."+v.Names[0].Name))
		}
	}
	for _, v := range result.List {
		v.Names = nil
	}
	notImpl.Type.Results = result
	ret := notImpl.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
	ret.Fun.(*ast.Ident).Name = notImpl.Recv.List[0].Names[0].Name + ".client.Call"
	var args []ast.Expr
	args = append(args, &ast.BasicLit{
		Kind:  token.TYPE,
		Value: "context.Background()",
	}, &ast.BasicLit{
		Kind:  token.STRING,
		Value: "\"" + m.Name.Name + "\"",
	}, &ast.BasicLit{
		Kind:  token.TYPE,
		Value: "arg",
	}, &ast.BasicLit{
		Kind:  token.TYPE,
		Value: "reply",
	})
	ret.Args = args
	r := ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{ret},
		Lhs: []ast.Expr{ast.NewIdent("err")},
	}
	//请求参数封装
	argr := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{ast.NewIdent("arg")},
	}
	if len(m.Params) == 1 {
		argr.Rhs = []ast.Expr{
			ast.NewIdent(strings.Title(m.Params[0].Name.Name)),
		}
	} else {
		argr.Rhs = []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X:  m.WrapRequest(),
			},
		}
	}
	notImpl.Body.List[0] = argr
	reply := &ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: m.ResponseStructName(),
				},
			},
		},
		Lhs: []ast.Expr{ast.NewIdent("reply")},
	}
	notImpl.Body.List = append(notImpl.Body.List, reply)
	//返回参数封装
	notImpl.Body.List = append(notImpl.Body.List, &r)
	s := &ast.IfStmt{
		Cond: &ast.BasicLit{
			Kind:  token.TYPE,
			Value: "err != nil",
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "log.Error",
						},
						Args: []ast.Expr{
							&ast.Ident{Name: `"RPCX调用错误:%s",err.Error()`},
						},
					},
				},
			},
		},
	}
	notImpl.Body.List = append(notImpl.Body.List, s)
	rnil := ast.ReturnStmt{
		Results: returnValue,
	}
	notImpl.Body.List = append(notImpl.Body.List, &rnil)
	root.Decls = append(root.Decls, notImpl)
}
