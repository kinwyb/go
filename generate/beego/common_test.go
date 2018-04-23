package beego

import (
	"fmt"
	"testing"
)

func Test_Common(t *testing.T) {
	common := []string{
		"// @Title 添加颜色档案,code编码,name名称",
		"// @Description 添加颜色档案,code编码,name名称",
		"// @Param code query string true 颜色编码",
		"// @Param name query string true 颜色名称",
		"// @router /add",
	}
	fmt.Println(ParseCommon(common))
}
