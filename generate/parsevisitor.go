package generate

import (
	"fmt"
	"go/ast"
)

type (
	parseVisitor struct {
		src *SourceContext
	}

	genDelVisitor struct {
		comments *ast.CommentGroup
		ps       *parseVisitor
		src      *SourceContext
	}

	typeSpecVisitor struct {
		src   *SourceContext
		gen   *genDelVisitor
		node  *ast.TypeSpec
		iface *Iface
		name  *ast.Ident
	}

	interfaceTypeVisitor struct {
		node    *ast.TypeSpec
		ts      *typeSpecVisitor
		methods []Method
	}

	methodVisitor struct {
		depth           int
		node            *ast.TypeSpec
		list            *[]Method
		name            *ast.Ident
		params, results *[]Arg
		isMethod        bool
		comments        []string
	}

	argListVisitor struct {
		list *[]Arg
	}

	argVisitor struct {
		node  *ast.TypeSpec
		parts []ast.Expr
		list  *[]Arg
	}
)

func (v *parseVisitor) Visit(n ast.Node) ast.Visitor {
	switch rn := n.(type) {
	default:
		return v
	case *ast.GenDecl:
		return &genDelVisitor{comments: rn.Doc, ps: v, src: v.src}
	case *ast.File:
		v.src.Pkg = rn.Name
		return v
	case *ast.ImportSpec:
		v.src.Imports = append(v.src.Imports, rn)
		return nil

	case *ast.TypeSpec:
		switch rn.Type.(type) {
		default:
			v.src.Types = append(v.src.Types, rn)
		case *ast.InterfaceType:
			// can't output interfaces
			// because they'd conflict with our implementations
		}
		return &typeSpecVisitor{src: v.src, node: rn}
	}
}

func (v *genDelVisitor) Visit(n ast.Node) ast.Visitor {
	switch rn := n.(type) {
	default:
		return v
	case *ast.File:
		v.src.Pkg = rn.Name
		return v
	case *ast.ImportSpec:
		v.src.Imports = append(v.src.Imports, rn)
		return nil

	case *ast.TypeSpec:
		switch rn.Type.(type) {
		default:
			v.src.Types = append(v.src.Types, rn)
		case *ast.InterfaceType:
			// can't output interfaces
			// because they'd conflict with our implementations
		}
		return &typeSpecVisitor{src: v.src, node: rn, gen: v}
	}
}

/*
package foo

type FooService interface {
	Bar(ctx context.Context, i int, s string) (string, error)
}
*/

func (v *typeSpecVisitor) Visit(n ast.Node) ast.Visitor {
	switch rn := n.(type) {
	default:
		return v
	case *ast.Ident:
		if v.name == nil {
			v.name = rn
		}
		return v
	case *ast.InterfaceType:
		return &interfaceTypeVisitor{ts: v, methods: []Method{}}
	case nil:
		if v.iface != nil {
			v.iface.Name = v.name
			sn := *v.name
			v.iface.StubName = &sn
			v.iface.StubName.Name = v.name.String()
			if v.gen != nil {
				v.iface.Comments = v.gen.comments
			}
			v.src.Interfaces = append(v.src.Interfaces, *v.iface)
		}
		return nil
	}
}

func (v *interfaceTypeVisitor) Visit(n ast.Node) ast.Visitor {
	switch n.(type) {
	default:
		return v
	case *ast.Field:
		return &methodVisitor{list: &v.methods}
	case nil:
		v.ts.iface = &Iface{Methods: v.methods}
		if v.ts.gen != nil {
			v.ts.iface.Comments = v.ts.gen.comments
		}
		return nil
	}
}

func (v *methodVisitor) Visit(n ast.Node) ast.Visitor {
	switch rn := n.(type) {
	default:
		v.depth++
		return v
	case *ast.Comment:
		v.comments = append(v.comments, rn.Text)
		v.depth++
		return v
	case *ast.Ident:
		if rn.IsExported() {
			v.name = rn
		}
		v.depth++
		return v
	case *ast.FuncType:
		v.depth++
		v.isMethod = true
		return v
	case *ast.FieldList:
		if v.params == nil {
			v.params = &[]Arg{}
			return &argListVisitor{list: v.params}
		}
		if v.results == nil {
			v.results = &[]Arg{}
		}
		return &argListVisitor{list: v.results}
	case nil:
		v.depth--
		if v.depth == 0 && v.isMethod && v.name != nil {
			method := Method{Name: v.name, Comments: v.comments}
			if v.results != nil {
				method.Results = *v.results
			}
			if v.params != nil {
				method.Params = *v.params
			}
			*v.list = append(*v.list, method)
			v.comments = v.comments[:0]
		}
		return nil
	}
}

func (v *argListVisitor) Visit(n ast.Node) ast.Visitor {
	switch n.(type) {
	default:
		return nil
	case *ast.Field:
		return &argVisitor{list: v.list}
	}
}

func (v *argVisitor) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.CommentGroup, *ast.BasicLit:
		return nil
	case *ast.Ident: //Expr -> everything, but clarity
		if t.Name != "_" {
			v.parts = append(v.parts, t)
		}
	case ast.Expr:
		v.parts = append(v.parts, t)
	case nil:
		names := v.parts[:len(v.parts)-1]
		tp := v.parts[len(v.parts)-1]
		iserr, isstar := v.parseTpIsError(tp)
		if len(names) == 0 {
			*v.list = append(*v.list, Arg{Typ: tp, IsError: iserr, IsStar: isstar})
			return nil
		}
		for _, n := range names {
			*v.list = append(*v.list, Arg{
				Name:    n.(*ast.Ident),
				Typ:     tp,
				IsError: iserr,
				IsStar:  isstar,
			})
		}
	}
	return nil
}

//是否是错误类型
func (v *argVisitor) parseTpIsError(tp ast.Expr) (bool, bool) {
	if v, ok := tp.(*ast.SelectorExpr); ok {
		return fmt.Sprintf("%s", v.Sel.Name) == "Error", false
	} else if _, ok := tp.(*ast.StarExpr); ok {
		return false, true
	}
	return false, false
}
