//+build !release

package logs

import (
	"fmt"
	"log"
	"runtime"
)

//调试追踪
func TraceCaller(format string, args ...interface{}) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	log.Printf("\x1b[46m[T] "+format+"\x1b[0m", args...)
}

//调试追踪带文件信息
func Trace(format string, args ...interface{}) {
	log.SetFlags(log.LstdFlags)
	_, file, line, ok := runtime.Caller(1)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Printf(fmt.Sprintf("\x1b[34m%s:%d\x1b[0m \x1b[32m[T] ", file, line)+format+"\x1b[0m", args...)
	} else {
		log.Printf("\x1b[33m[T] "+format, args...)
	}
}
