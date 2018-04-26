package rpcx

import (
	"bytes"
	"testing"

	"github.com/kinwyb/go/generate"
)

func TestLay_TransformAST(t *testing.T) {
	generate.ParseFile(bytes.NewReader([]byte(tp)), &lay{})
}

func TestLayClient_TransformAST(t *testing.T) {
	generate.ParseFile(bytes.NewReader([]byte(tp)), &layclient{})
}

var tp = `
package gokit

import (
	"ASID/models/objs"

	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err1"
)

// ColorService describes the service.
type ColorService interface {

	//添加颜色档案,code编码,name名称
	Add(code, name, creator string, dbtag string) err1.Error

	//更新颜色档案,code要更新的颜色档案编码,name名称,modifier 修改人
	//state 更新的状态
	Update(code, name, modifier string,
		state objs.EnableState, dbtag string) (string,err1.Error)

	//更新颜色档案名称,code要更新的颜色档案编码,name名称,modifier 修改人
	UpdateName(code, name, modifier string, dbtag string) err1.Error

	//更新颜色档案状态,code要更新的颜色档案编码,name名称,modifier 修改人
	UpdateEnableState(code, modifier string,
		state objs.EnableState, dbtag string) err1.Error

	//根据编码查询颜色档案 code要查询的编码
	QueryByCode(code string) *objs.ArchivesColor

	//根据名称模糊查询颜色档案 name要查询的名称
	QueryByLikeName(name string, dbtag string) []*objs.ArchivesColor

	//获取颜色档案列表
	QueryList(codeOrName string, state objs.EnableState, pg *db.PageObj, dbtag string) []*objs.ArchivesColor
}
`
