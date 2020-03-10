package log

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// line number hook for log the call context,
type lineHook struct {
	Field string
	// skip为遍历调用栈开始的索引位置
	Skip   int
	levels []log.Level
}

// Levels implement levels
func (hook lineHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire implement fire
func (hook lineHook) Fire(entry *log.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	var pc uintptr
	// 遍历调用栈的最大索引为第11层.
	for i := 0; i < 11; i++ {
		file, line, pc = getCaller(skip + i)
		// 过滤掉所有logrus包，即可得到生成代码信息
		if !strings.HasPrefix(file, "logrus") &&
			!strings.HasPrefix(file, "log") {
			break
		}
	}

	fullFnName := runtime.FuncForPC(pc)

	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		// 取得函数名
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d:%s()", file, line, fnName)
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0

	// 获取包名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}
