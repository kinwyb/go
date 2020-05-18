package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kinwyb/go/db"
)

var SQLEmptyChange = errors.New("数据无变化")
var DuplicateField = errors.New("字段重复")

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
func ExecResultHasError(execresult db.ExecResult, reportZeroChange bool, param ...map[string]string) error {
	retError := execresult.HasError(reportZeroChange)
	if retError != nil {
		errCode, _ := formatError(retError)
		if errCode == DuplicateErrorCode { //字段重复
			field := GetDuplicateField(retError.Error())
			if len(param) > 0 {
				if v, ok := param[0][field]; ok {
					return fmt.Errorf("[%s]%w", v, DuplicateField)
				}
			} else if field == PRIMARY {
				return fmt.Errorf("[主键]%w", DuplicateField)
			}
			return fmt.Errorf("唯一%w", DuplicateField)
		}
		return retError
	}
	return nil
}
