package logs

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

//新文件日志
func NewBeegoFileLog(level Level, fileName string, maxday int) Logger {
	ret := logs.NewLogger()
	ret.SetLevel(int(level))
	ret.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","level":%d,"maxlines":0,"maxsize":0,"daily":true,"maxdays":%d}`, fileName, level, maxday))
	return ret
}

//新日志
func NewBeegoLogger() Logger {
	return logs.NewLogger()
}
