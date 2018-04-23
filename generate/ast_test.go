package generate

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	src := bytes.NewReader([]byte(fileData))
	data, err := ParseFile(src, nil)
	if err != nil && err != NoLayouterError {
		t.Fatal(err)
	}
	fmt.Printf("pkg := %+v\n", data.Pkg.Name)
	fmt.Printf("imports := %+v\n", data.Imports[0].Path.Value)
	for _, v := range data.Interfaces {
		fmt.Printf("接口: %+v\n", v.Name.Name)
		for _, x := range v.Methods {
			fmt.Printf("接口方法: %+v\n", x.Name.Name)
			fmt.Printf("接口方法备注: \n\t%+v\n", strings.Join(x.Comments, "\n\t"))
		}
	}
}

var fileData = `
package astf

import (
	"tes/ssd"
)

type testinterface interface{

	// 1.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	// 2.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	// 3.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	UpdateName(code, name, modifier string, dbtag string) err1.Error

	// 4.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	// 5.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	// 6.更新部门档案名称,code要更新的部门档案编码,name名称,modifier 修改人
	UpdateName2(code, name, modifier string, dbtag string) err1.Error
}

//dfadfs
//asfsfs
func d() error {
	return nil
}

`
