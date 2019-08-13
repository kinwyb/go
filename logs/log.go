package logs

import (
	"time"

	"github.com/astaxie/beego/logs"

	"io"
)

type Level = int

const (
	Emergency Level = iota
	Alert
	Critical
	Error
	Warn
	Notice
	Info
	Debug
)

//Logger 日志接口
type Logger interface {
	io.Writer
	//输出
	Debug(format string, args ...interface{})
	//输出
	Info(format string, args ...interface{})
	//警告
	Warning(format string, args ...interface{})
	//注意
	Notice(format string, args ...interface{})
	//错误
	Error(format string, args ...interface{})
	//关键
	Critical(format string, args ...interface{})
	//警报
	Alert(format string, args ...interface{})
	//紧急
	Emergency(format string, args ...interface{})
}

func NewLogger(level ...Level) Logger {
	ret := logs.NewLogger()
	ret.SetLevel(Debug)
	if len(level) > 0 {
		ret.SetLevel(level[0])
	}
	return ret
}

//WriteLog 写入日志
func WriteLog(log Logger, level Level, format string, args ...interface{}) {
	if log == nil {
		return
	}
	switch level {
	case Emergency:
		log.Emergency(format, args...)
	case Alert:
		log.Alert(format, args...)
	case Critical:
		log.Critical(format, args...)
	case Error:
		log.Error(format, args...)
	case Warn:
		log.Warning(format, args...)
	case Info:
		log.Info(format, args...)
	case Debug:
		log.Debug(format, args...)
	}
}

//注册一个日志获取函数
type RegisterLogFunc func(log *LogFiles)

var logFactory = NewLogFiles("", 24*time.Hour)

var logmap []RegisterLogFunc

//设置日志路径
func SetLogPath(filepath string, level ...Level) {
	if filepath == "" {
		return
	} else if len(level) < 1 {
		level = []Level{Debug}
	}
	logFactory = NewLogFiles(filepath, 24*time.Hour, level[0])
	for _, v := range logmap {
		if v != nil {
			v(logFactory)
		}
	}
}

func SetLevel(level Level) {
	logFactory.Level(level)
}

func RegisterLog(fun RegisterLogFunc) {
	if fun != nil {
		logmap = append(logmap, fun)
	}
}

//获取一个日志
func GetLogger(logname string) Logger {
	return logFactory.GetLog(logname)
}
