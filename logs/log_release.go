//+build release

package logs

//调试追踪
func Trace(format string, args ...interface{}) {}

//调试追踪带文件信息
func TraceCaller(format string, args ...interface{}) {}
