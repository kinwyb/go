package beego

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/kinwyb/go/generate"
)

type Param struct {
	// @Param pageSize query int false 每页条数
	Name    string //参数名称
	Tl      string //查询类型
	Tp      string //参数类型
	Request bool   //是否必要
	Desc    string //描述
}

//解析参数信息
func ParseCommon(common []string) map[string]*Param {
	ret := map[string]*Param{}
	//解析参数信息
	for _, v := range common {
		if strings.HasPrefix(v, "// @Param") {
			fmt.Printf("解析参数：%s\n", v)
			v = strings.Trim(v, " ")
			v = strings.TrimPrefix(v, "// @Param")
			v = strings.Trim(v, " ")
			vs := strings.Split(v, " ")
			p := &Param{}
			p.Name = vs[0]
			p.Tl = vs[1]
			p.Tp = vs[2]
			if len(vs) > 3 {
				p.Request = "true" == vs[3]
			}
			if len(vs) > 4 {
				p.Desc = vs[4]
			}
			ret[p.Name] = p
		}
	}
	return ret
}

//解析参数到代码,参数代码，bool表示是否存在分页参数
func ParseParam(params map[string]*Param, m *generate.Method, rev string) ([]ast.Stmt, bool) {
	var ret []ast.Stmt
	hasPage := false
	for i, v := range m.Params {
		paramName := v.Name.Name
		if paramName == "pg" {
			ret = append(ret, queryPage(rev)...)
			hasPage = true
			continue
		}
		p := params[paramName]
		if p == nil {
			fmt.Printf("[%s]参数[%s]不存在\n", m.Name.Name, paramName)
			continue
		} else if p.Name == "token" { //token值
			continue
		}
		switch p.Tl {
		case "query":
			ret = append(ret, queryParam(paramName, p, rev)...)
		case "body":
			m.Params[i].IsBody = true
			ret = append(ret, body(paramName, p, rev)...)
		}
	}
	return ret, hasPage
}

//query参数代码
func queryParam(paramName string, p *Param, rev string) []ast.Stmt {
	switch p.Tp {
	case "int":
		return queryInt(rev, paramName, p.Request)
	case "int64":
		return queryInt64(rev, paramName, p.Request)
	case "float", "float64":
		return queryFloat(rev, paramName, p.Request)
	case "bool":
		return queryBool(rev, paramName, p.Request)
	default:
		return queryString(rev, paramName, p.Request)
	}
}

func queryString(rev string, param string, request bool) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: param,
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetString",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"" + param + "\"",
						},
					},
				},
			},
		},
	}
	if request {
		s := &ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: param + " == \"\"",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamMissing"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		ret = append(ret, s)
	}
	return ret
}

func queryInt(rev string, param string, request bool) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: param,
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetInt",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"" + param + "\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "-1",
						},
					},
				},
			},
		},
	}
	if request {
		s := &ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: param + " == -1",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamMissing"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		ret = append(ret, s)
	}
	return ret
}

func queryInt64(rev string, param string, request bool) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: param,
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetInt64",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"" + param + "\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "-1",
						},
					},
				},
			},
		},
	}
	if request {
		s := &ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: param + " == -1",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamMissing"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		ret = append(ret, s)
	}
	return ret
}

func queryFloat(rev string, param string, request bool) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: param,
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetFloat",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"" + param + "\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "-1",
						},
					},
				},
			},
		},
	}
	if request {
		s := &ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: param + " == -1",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamMissing"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		ret = append(ret, s)
	}
	return ret
}

func queryBool(rev string, param string, request bool) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: param,
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetBool",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"" + param + "\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "false",
						},
					},
				},
			},
		},
	}
	if request {
		s := &ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: param + " == -1",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamMissing"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}
		ret = append(ret, s)
	}
	return ret
}

func queryPage(rev string) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "page",
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetInt",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"page\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "1",
						},
					},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "pageSize",
				},
				&ast.Ident{
					Name: "_",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: rev + ".GetInt",
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "\"pageSize\"",
						},
						&ast.BasicLit{
							Kind:  token.TYPE,
							Value: "20",
						},
					},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "pg",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: ast.NewIdent("db.PageObj"),
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Page"),
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "page",
								},
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Rows"),
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "pageSize",
								},
							},
						},
					},
				},
			},
		},
	}
}

//body类型代码
func body(paramName string, p *Param, rev string) []ast.Stmt {
	ret := []ast.Stmt{
		&ast.DeclStmt{ // var req xxxxx
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("req")},
						Type:  ast.NewIdent(p.Tp),
					},
				},
			},
		},
		&ast.AssignStmt{ // err = json.Unmarshal(i.Ctx.Input.RequestBody,req)
			Lhs: []ast.Expr{
				ast.NewIdent("e"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.TYPE,
					Value: "json.Unmarshal(" + rev + ".Ctx.Input.RequestBody,&req)",
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BasicLit{
				Kind:  token.TYPE,
				Value: "e != nil",
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: rev + ".RespError",
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "ParamDecodeFail"},
								&ast.Ident{Name: "e"},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		},
	}
	return ret
}
