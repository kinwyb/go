//+build !release

package logs

import (
	"log"
)

//调试追踪
func TraceCaller(format string, args ...interface{}) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	log.Printf("[T] "+format, args...)
}

//调试追踪带文件信息
func Trace(format string, args ...interface{}) {
	log.SetFlags(log.LstdFlags)
	log.Printf("[T] "+format, args...)
}
