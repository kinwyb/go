package mysql

import (
	"strings"
)

const PRIMARY = "primary"

const (
	DuplicateErrorCode = 1062 //字段重复
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
