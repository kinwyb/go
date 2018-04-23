package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"sync"
)

var NoLayouterError = errors.New("Layouter空")
var layMap = &sync.Map{}

type Layouter interface {
	TransformAST(ctx *SourceContext, filedir ...string) error
}

//注册解析对象
func RegisterLayouter(name string, lay Layouter) {
	if _, ok := layMap.Load(name); ok {
		fmt.Printf("%s对象已经存在\n", name)
		return
	}
	layMap.Store(name, lay)
}

//根据解析名称解析文件
func ParseFileByLayoutName(source io.Reader, layname string, outfiledir ...string) (*SourceContext, error) {
	lay, ok := layMap.Load(layname)
	if !ok {
		return nil, NoLayouterError
	}
	return ParseFile(source, lay.(Layouter), outfiledir...)
}

func ParseFile(source io.Reader, lay Layouter, outfiledir ...string) (*SourceContext, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	//ast.Print(fset, f)
	context := &SourceContext{}
	visitor := &parseVisitor{src: context}
	ast.Walk(visitor, f)
	if context.validate() != nil {
		return nil, context.validate()
	} else if lay != nil { //调用具体的
		return context, lay.TransformAST(context, outfiledir...)
	}
	return context, NoLayouterError
}

func NewAstFile(pkgname string) *ast.File {
	file := &ast.File{
		Name:  ast.NewIdent(pkgname),
		Decls: []ast.Decl{},
	}
	return file
}

func AddImport(root *ast.File, path string) {
	for _, d := range root.Decls {
		if imp, is := d.(*ast.GenDecl); is && imp.Tok == token.IMPORT {
			for _, s := range imp.Specs {
				if s.(*ast.ImportSpec).Path.Value == `"`+path+`"` {
					return // already have one
					// xxx aliased imports?
				}
			}
		}
	}
	root.Decls = append(root.Decls, importFor(importSpec(path)))
}

func importFor(is *ast.ImportSpec) *ast.GenDecl {
	return &ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{is}}
}

func importSpec(path string) *ast.ImportSpec {
	return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}}
}
