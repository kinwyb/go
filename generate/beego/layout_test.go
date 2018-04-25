package beego

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kinwyb/go/generate"
)

func TestLay_TransformAST(t *testing.T) {
	generate.ParseFile(bytes.NewReader([]byte(tp)), &lay{})
}

func TestAST(t *testing.T) {
	src := `
package main
//备注内容
type ier interface {}

func main() {
    println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

const tp = `
package controller

//测试接口备注
type IColorEndPoint interface {

	// @Title 添加颜色档案,code编码,name名称
	// @Description 添加颜色档案,code编码,name名称
	// @Param token header string true 登入状态
	// @Param code query string true 颜色编码
	// @Param name query string true 颜色名称
	// @router /add
	Add(code, name string, token string) err1.Error

	// @Title 更新颜色档案,code要更新的颜色档案编码,name名称,state 更新的状态
	// @Description 更新颜色档案,code要更新的颜色档案编码,name名称,state 更新的状态,名称和状态必填一个
	// @Param token header string true 登入状态
	// @Param code query string true 要更新的颜色编码
	// @Param name query string false 新颜色名称
	// @Param state query int false 状态[2=启用,3=停用]
	// @router /update
	Update(code, name string, state objs.EnableState, token string) err1.Error

	// @Title 更新颜色档案名称,code要更新的颜色档案编码,name名称
	// @Description 更新颜色名称
	// @Param token header string true 登入状态
	// @Param code query string true 要更新的颜色编码
	// @Param name query string true 新颜色名称
	// @router /update/name
	UpdateName(code, name string, token string) err1.Error

	// @Title 更新颜色档案状态,code要更新的颜色档案编码,name名称
	// @Description 更新颜色档案状态
	// @Param token header string true 登入状态
	// @Param code query string true 要更新的颜色编码
	// @Param state query int true 颜色状态[2=启用,3=停用]
	// @router /update/state
	UpdateEnableState(code string, state objs.EnableState, token string) err1.Error

	// @Title 根据编码查询颜色档案 code要查询的编码
	// @Description 根据颜色编码查询颜色信息
	// @Param token header string true 登入状态
	// @Param code query string true 颜色编码
	// @Success 200 {object} objs.ArchivesColor
	// @router /query/code
	QueryByCode(code string, token string) (*objs.ArchivesColor, err1.Error)

	// @Title 根据名称模糊查询颜色档案 name要查询的名称
	// @Description 根据名称查询颜色信息
	// @Param token header string true 登入状态
	// @Param name query string true 颜色名称
	// @Success 200 {array} objs.ArchivesColor
	// @router /query/name
	QueryByLikeName(name string, token string) ([]*objs.ArchivesColor, err1.Error)

	// @Title 获取颜色档案列表
	// @Description 查询颜色列表
	// @Param token header string true 登入状态
	// @Param codeOrName query string false 指定颜色名称或编码
	// @Param state query int false 颜色状态[2=启用,3=停用]
	// @Param page query int false 当前页数
	// @Param pageSize query int false 每页条数
	// @Success 200 {array} objs.ArchivesColor
	// @router /query
	QueryList(codeOrName string, state objs.EnableState,
		pg *db2.PageObj, token string) ([]*objs.ArchivesColor, err1.Error)

	// @Title 获取颜色档案列表
	// @Description 查询颜色列表
	// @Param token header string true 登入状态
	// @Param req body objs.ArchivesColor true 指定颜色名称或编码
	// @Success 200 {array} objs.ArchivesColor
	// @router /query
	QueryList1(req *objs.ArchivesColor) ([]*objs.ArchivesColor, err1.Error)
}
`
