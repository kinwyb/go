package beego

import (
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kinwyb/go/generate"
)

type lay struct{}

func init() {
	generate.RegisterLayouter("beego", &lay{})
}

func (l *lay) TransformAST(ctx *generate.SourceContext, filedir ...string) error {
	packagename := "controllers"
	if len(filedir) > 0 && filedir[0] != "" {
		_, packagename = filepath.Split(filedir[0])
	}
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile(packagename)
		ctx.ImportDecls(service) //import
		name := v.StubName.Name
		name = strings.TrimPrefix(name, "I")
		name = strings.TrimSuffix(name, "EndPoint")
		name = name + "Controller"
		ds := generate.StructDecl(ast.NewIdent(name), &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("Controller"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("Serv")},
					Type:  ast.NewIdent(ctx.Pkg.Name + "." + v.StubName.Name),
				},
			},
		}, v.Comments)
		service.Decls = append(service.Decls, ds)
		for _, meth := range v.Methods {
			if meth.Comments == nil || len(meth.Comments) < 2 {
				//注解不全的忽略
				continue
			} else if !strings.HasPrefix(meth.Comments[0], "// @Title") {
				// 不是@Title开头的忽略
				continue
			}
			meth.Prefix = ctx.Prefix
			addMethodController(service, &v, &meth)
		}
		filedata, err := generate.FormatNode("", service)
		if err != nil {
			panic("err:" + err.Error())
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

//生成服务代码
func addMethodController(root *ast.File, ifc *generate.Iface, m *generate.Method) {
	notImpl := &ast.FuncDecl{
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}
	notImpl.Name = m.Name
	notImpl.Recv = generate.FieldList(ifc.Reciever())
	name := ifc.StubName.Name
	name = strings.TrimPrefix(name, "I")
	name = strings.TrimSuffix(name, "EndPoint")
	name = name + "Controller"
	notImpl.Recv.List[0].Type = &ast.StarExpr{
		X: &ast.Ident{Name: name},
	}
	//生成请求参数
	notImpl.Type.Params = nil
	notImpl.Type.Results = nil
	//解析参数
	params := ParseCommon(m.Comments)
	stmts, hasPage := ParseParam(params, m, notImpl.Recv.List[0].Names[0].Name)
	notImpl.Doc = &ast.CommentGroup{}
	for _, cm := range m.Comments {
		notImpl.Doc.List = append(notImpl.Doc.List, &ast.Comment{
			Text: cm,
		})
	}
	//增加调用代码
	ret := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: notImpl.Recv.List[0].Names[0].Name + ".Serv." + m.Name.Name,
		},
	}
	for _, v := range m.Params {
		if v.Name.Name == "token" {
			ret.Args = append(ret.Args, &ast.BasicLit{
				Kind:  token.TYPE,
				Value: "i.Token",
			})
			continue
		} else if v.IsStar && v.Name.Name != "pg" {
			ret.Args = append(ret.Args, &ast.BasicLit{
				Kind:  token.TYPE,
				Value: "&" + v.Name.Name,
			})
			continue
		}
		ret.Args = append(ret.Args, &ast.BasicLit{
			Kind:  token.TYPE,
			Value: v.Name.Name,
		})
	}
	returnStmt := &ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{ret},
	}
	hasError := false
	hasRet := false
	for _, v := range m.Results {
		if v.IsError {
			returnStmt.Lhs = append(returnStmt.Lhs, &ast.Ident{
				Name: "err",
			})
			hasError = true
		} else {
			returnStmt.Lhs = append(returnStmt.Lhs, &ast.Ident{
				Name: "ret",
			})
			hasRet = true
		}
	}
	stmts = append(stmts, returnStmt)
	//返回结果语句
	if hasError { //如果存在错误.判断错误是否为空返回错误信息
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
								Name: notImpl.Recv.List[0].Names[0].Name + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "err"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		stmts = append(stmts, s)
	}
	if hasPage { //增加分页信息
		pgstmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: notImpl.Recv.List[0].Names[0].Name + ".Page",
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "pg"},
				},
			},
		}
		stmts = append(stmts, pgstmt)
	}
	if hasRet { //如果存在返回内容,返回内容
		retstmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: notImpl.Recv.List[0].Names[0].Name + ".ResponseSUCC",
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "ret"},
				},
			},
		}
		stmts = append(stmts, retstmt)
	} else {
		retstmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: notImpl.Recv.List[0].Names[0].Name + ".ResponseSUCC",
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "Success"},
				},
			},
		}
		stmts = append(stmts, retstmt)
	}
	//合并语句
	notImpl.Body.List = stmts
	root.Decls = append(root.Decls, notImpl)
}
