package err1

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

//获取代码调用堆栈
func Caller() []CallerInfo {
	var infos []CallerInfo
	skip := 2
	for ; ; skip++ {
		name, file, line, ok := callerInfo(skip + 1)
		if !ok {
			break
		}
		if strings.HasPrefix(name, "runtime.") {
			break
		}
		infos = append(infos, CallerInfo{
			FuncName: name,
			FileName: file,
			FileLine: line,
		})
	}
	return infos
}

// 获取堆栈信息字符串=>文件名:行数 方法名称
// 如果没有结果返回 ???:- ???
func CallInfo(skip int) string {
	name, file, line, ok := callerInfo(skip)
	if !ok {
		return "???:- ???"
	}
	return fmt.Sprintf("%s:%d %s", file, line, name)
}

//打印堆栈信息
func PrintCaller() string {
	bf := &bytes.Buffer{}
	info := Caller()
	for i, v := range info {
		bf.WriteString(fmt.Sprintf("%3d => %s[%d] %s\n", i, v.FileName, v.FileLine, v.FuncName))
	}
	return bf.String()
}

//获取堆栈信息
func callerInfo(skip int) (name, file string, line int, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		name = "???"
		file = "???"
		line = 1
		return
	}
	name = runtime.FuncForPC(pc).Name()
	// Truncate file name at last file name separator.
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	} else if idx = strings.LastIndex(name, "\\"); idx >= 0 {
		name = name[idx+1:]
	}
	// Truncate file name at last file name separator.
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	} else if idx = strings.LastIndex(file, "\\"); idx >= 0 {
		file = file[idx+1:]
	}
	return
}
