package rpcx

import (
	"ASID/generate"
	"fmt"
	"go/ast"
)

type lay struct{}

func (l *lay) TransformAST(ctx *generate.SourceContext, filedir ...string) error {
	//遍历所有接口
	for _, v := range ctx.Interfaces {
		service := generate.NewAstFile("serv")
		ctx.ImportDecls(service) //import
		v.StubStructDecl(service)
		for _, meth := range v.Methods {
			meth.Prefix = ctx.Prefix
			//生成请求结构
			addRequestStruct(service, &meth)
			//生成返回结果结构
			addResponseStruct(service, &meth)
			addMethodService(service, &v, &meth)
			//addRequestStruct(endpoints, meth)
			//addResponseStruct(endpoints, meth)
			//addEndpointMaker(endpoints, iface, meth)
		}
		filedata, err := generate.FormatNode("", service)
		if err != nil {
			panic("err")
		}
		fmt.Printf("%s", filedata)
	}

	return nil
}

func addRequestStruct(root *ast.File, meth *generate.Method) {
	root.Decls = append(root.Decls, meth.RequestStruct())
}

func addResponseStruct(root *ast.File, meth *generate.Method) {
	root.Decls = append(root.Decls, meth.ResponseStruct())
}
