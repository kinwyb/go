package rpcx

import (
	"fmt"
	"go/ast"

	"io/ioutil"
	"os"
	"path/filepath"

	"strings"

	"github.com/kinwyb/go/generate"
)

type lay struct{}

type layclient struct{}

//服务端
func (l *lay) TransformAST(ctx *generate.SourceContext, filedir ...string) error {
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile("rpcx")
		ctx.ImportDecls(service) //import
		name := v.StubName.Name
		name = strings.TrimPrefix(name, "I")
		name = strings.TrimSuffix(name, "EndPoint")
		name = name + "Rpcx"
		ds := generate.StructDecl(ast.NewIdent(name), &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{ast.NewIdent("serv")},
				Type:  ast.NewIdent("endPoints." + v.StubName.Name),
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
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile("rpcx")
		ctx.ImportDecls(service) //import
		name := v.StubName.Name
		name = strings.TrimPrefix(name, "I")
		name = strings.TrimSuffix(name, "EndPoint")
		name = name + "RpcxClient"
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
