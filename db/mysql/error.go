package mysql

import (
	"strings"

	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err1"
)

var SQLEmptyChange = err1.NewError(101, "数据无变化")

const (
	DuplicateErrorCode = 1062 //字段重复
	PRIMARY            = "primary"
)

//获取重复字段
func GetDuplicateField(errmsg string) string {
	strs := strings.Split(errmsg, " key ")
	if len(strs) < 2 {
		return PRIMARY
	}
	field := strings.Trim(strs[1], " ")
	field = strings.Trim(strs[1], "'")
	return strings.ToLower(field)
}

//执行结果是否有错误
func ExecResultHasError(execresult db.ExecResult, reportZeroChange bool, param ...map[string]string) err1.Error {
	retError := execresult.HasError(reportZeroChange)
	if retError != nil {
		if retError.Code() == DuplicateErrorCode { //字段重复
			field := GetDuplicateField(retError.Msg())
			if len(param) < 1 {
				param = []map[string]string{}
			}
			if v, ok := param[0][field]; ok {
				return err1.NewError(retError.Code(), "["+v+"]重复")
			} else if field == PRIMARY {
				return err1.NewError(retError.Code(), "[主键]重复")
			}
			return err1.NewError(retError.Code(), "唯一数据重复")
		}
		return retError
	}
	return nil
}
