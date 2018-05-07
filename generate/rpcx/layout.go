package rpcx

import (
	"fmt"
	"go/ast"
	"go/token"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kinwyb/go/generate"
)

type lay struct{}

type layclient struct{}

//服务端
func (l *lay) TransformAST(ctx *generate.SourceContext, filedir ...string) error {
	packagename := "rpcx"
	if len(filedir) > 0 {
		_, packagename = filepath.Split(filedir[0])
	}
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile(packagename)
		ctx.ImportDecls(service) //import
		name := v.StubName.Name
		//name = strings.TrimPrefix(name, "I")
		//name = strings.TrimSuffix(name, "EndPoint")
		name = name + "Rpcx"
		ds := generate.StructDecl(ast.NewIdent(name), &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{ast.NewIdent("serv")},
				Type:  ast.NewIdent(ctx.Pkg.Name + "." + v.StubName.Name),
			}},
		})
		service.Decls = append(service.Decls, ds)
		for _, meth := range v.Methods {
			meth.Prefix = ctx.Prefix + name
			//生成请求结构
			addRequestStruct(service, &meth)
			//生成返回结果结构
			addResponseStruct(service, &meth)
			addMethodService(service, &v, &meth, name)
		}
		filedata, err := generate.FormatNode("", service)
		if err != nil {
			panic(err)
		}
		if len(filedir) < 1 || filedir[0] == "" {
			fmt.Printf("%s", filedata)
		} else {
			path := filepath.Join(filedir[0], name+".go")
			err := ioutil.WriteFile(path, filedata.Bytes(), os.ModePerm)
			if err != nil {
				fmt.Printf("[%s]文件保存错误:%s\n", name, err.Error())
			} else {
				fmt.Printf("[%s]文件保存成功\n", name)
			}
		}
	}

	return nil
}

//客户端
func (l *layclient) TransformAST(ctx *generate.SourceContext, filedir ...string) error {
	packagename := "rpcxclient"
	if len(filedir) > 0 {
		_, packagename = filepath.Split(filedir[0])
	}
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile(packagename)
		ctx.ImportDecls(service) //import
		name := v.StubName.Name
		//name = strings.TrimPrefix(name, "I")
		//name = strings.TrimSuffix(name, "EndPoint")
		name = name + "RpcxClient"
		addClientNewStruct(service, name)
		ds := generate.StructDecl(ast.NewIdent(name), &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{ast.NewIdent("client")},
				Type:  ast.NewIdent("client.XClient"),
			}},
		})
		service.Decls = append(service.Decls, ds)
		for _, meth := range v.Methods {
			meth.Prefix = ctx.Prefix + name
			//生成请求结构
			addRequestStruct(service, &meth)
			//生成返回结果结构
			addResponseStruct(service, &meth)
			addMethodClient(service, &v, &meth, name)
		}
		filedata, err := generate.FormatNode("", service)
		if err != nil {
			panic(err)
		}
		if len(filedir) < 1 || filedir[0] == "" {
			fmt.Printf("%s", filedata)
		} else {
			path := filepath.Join(filedir[0], name+".go")
			err := ioutil.WriteFile(path, filedata.Bytes(), os.ModePerm)
			if err != nil {
				fmt.Printf("[%s]文件保存错误:%s\n", name, err.Error())
			} else {
				fmt.Printf("[%s]文件保存成功\n", name)
			}
		}
	}

	return nil
}

func addRequestStruct(root *ast.File, meth *generate.Method) {
	result := meth.RequestStruct()
	if result == nil {
		return
	}
	root.Decls = append(root.Decls, result)
}

func addResponseStruct(root *ast.File, meth *generate.Method) {
	result := meth.ResponseStruct()
	if result == nil {
		return
	}
	root.Decls = append(root.Decls, result)
}

func addClientNewStruct(root *ast.File, name string) {
	newfunc := &ast.FuncDecl{
		Name: ast.NewIdent("New" + name),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("discovery"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "discovery.Clone", //表达式内容
							},
							Args: []ast.Expr{ //表达式参数集合
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "servicePath",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("xclient"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "client.NewXClient", //表达式内容
							},
							Args: []ast.Expr{ //表达式参数集合
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "servicePath",
								},
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "client.Failover",
								},
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "client.RoundRobin",
								},
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "discovery",
								},
								&ast.BasicLit{
									Kind:  token.TYPE,
									Value: "client.DefaultOption",
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(name),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("client"),
										Value: ast.NewIdent("xclient"),
									},
								},
							},
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
						Names: []*ast.Ident{
							ast.NewIdent("servicePath"),
						},
					},
					{
						Type: ast.NewIdent("client.ServiceDiscovery"),
						Names: []*ast.Ident{
							ast.NewIdent("discovery"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent(name),
						},
					},
				},
			},
		},
	}
	if newfunc != nil {
		root.Decls = append(root.Decls, newfunc)
	}
}
