package generate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"runtime"
	"sort"
	"strings"

	"reflect"

	"golang.org/x/tools/imports"
)

var ostype = runtime.GOOS

type sortableDecls []ast.Decl

func (sd sortableDecls) Len() int {
	return len(sd)
}

func (sd sortableDecls) Less(i int, j int) bool {
	switch left := sd[i].(type) {
	case *ast.GenDecl:
		switch right := sd[j].(type) {
		default:
			return left.Tok == token.IMPORT
		case *ast.GenDecl:
			return left.Tok == token.IMPORT && right.Tok != token.IMPORT
		}
	}
	return false
}

func (sd sortableDecls) Swap(i int, j int) {
	sd[i], sd[j] = sd[j], sd[i]
}

func FormatNode(fname string, node ast.Node) (*bytes.Buffer, error) {
	if file, is := node.(*ast.File); is {
		sort.Stable(sortableDecls(file.Decls))
	}
	outfset := token.NewFileSet()
	buf := &bytes.Buffer{}
	err := format.Node(buf, outfset, node)
	if err != nil {
		return nil, err
	}
	if ostype == "windows" && fname == "" {
		fname = "\\tmp.go"
	}
	imps, err := imports.Process(fname, buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(imps), nil
}

func StructDecl(name *ast.Ident, fields *ast.FieldList, comment ...*ast.CommentGroup) ast.Decl {
	return typeDecl(&ast.TypeSpec{
		Name: name,
		Type: &ast.StructType{
			Fields: fields,
		},
	}, comment...)
}

func typeDecl(ts *ast.TypeSpec, comment ...*ast.CommentGroup) ast.Decl {
	ret := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{ts},
	}
	if len(comment) > 0 {
		ret.Doc = comment[0]
	}
	return ret
}

func Field(n *ast.Ident, t ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{n},
		Type:  t,
	}
}

func FieldList(list ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: list}
}

func MappedFieldList(fn func(Arg) *ast.Field, args ...Arg) *ast.FieldList {
	fl := &ast.FieldList{List: []*ast.Field{}}
	for _, a := range args {
		fl.List = append(fl.List, fn(a))
	}
	return fl
}

func ScopeWith(names ...string) *ast.Scope {
	scope := ast.NewScope(nil)
	for _, name := range names {
		scope.Insert(ast.NewObj(ast.Var, name))
	}
	return scope
}

func TypeOf(tp interface{}) string {
	return reflect.TypeOf(tp).Name()
}

func export(s string) string {
	return strings.Title(s)
}

func InventName(t ast.Expr, scope *ast.Scope) *ast.Ident {
	n := baseName(t)
	for try := 0; ; try++ {
		nstr := pickName(n, try)
		obj := ast.NewObj(ast.Var, nstr)
		if alt := scope.Insert(obj); alt == nil {
			return ast.NewIdent(nstr)
		}
	}
}

func baseName(t ast.Expr) string {
	switch tt := t.(type) {
	default:
		panic(fmt.Sprintf("don't know how to choose a base name for %T (%[1]v)", tt))
	case *ast.MapType:
		return "map"
	case *ast.InterfaceType:
		return "inf"
	case *ast.ArrayType:
		return "slice"
	case *ast.Ident:
		return tt.Name
	case *ast.SelectorExpr:
		return tt.Sel.Name
	case *ast.StarExpr:
		return baseName(tt.X)
	}
}

func pickName(base string, idx int) string {
	if idx == 0 {
		switch base {
		default:
			return strings.Split(base, "")[0]
		case "Context":
			return "ctx"
		case "error":
			return "err"
		}
	}
	return fmt.Sprintf("%s%d", base, idx)
}

func FetchFuncDecl(name string) *ast.FuncDecl {
	root := templateAST()
	for _, decl := range root.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name == name {
				return f
			}
		}
	}
	panic(fmt.Errorf("No function called %q in 'templates/full.go'", name))
}

func templateAST() *ast.File {
	tpASTs, err := parser.ParseFile(token.NewFileSet(), "",
		bytes.NewReader([]byte(template)), parser.DeclarationErrors)
	if err != nil {
		panic(err)
	}
	return tpASTs
}

var template = `
package generate

import (
	"context"
	"errors"
)

type ExampleService struct{}

type ExampleRequest struct {
	I int
	S string
}
type ExampleResponse struct {
	S   string
	Err error
}

func (f ExampleService) Example(ctx context.Context, i int, s string) (string, error) {
	panic(errors.New("not implemented"))
}

`
