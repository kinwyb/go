package logs

import (
	"github.com/sirupsen/logrus"
)

var DefaultJsonFormatter = logrus.JSONFormatter{
	TimestampFormat: "2006-01-02 15:04:05",
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyLevel: "@level",
		logrus.FieldKeyMsg:   "@message",
		logrus.FieldKeyFunc:  "@caller",
	},
}
var DefaultTextFormatter = logrus.TextFormatter{
	DisableColors:   true,
	TimestampFormat: "2006-01-02 15:04:05",
}
var DefaultColorTextFormatter = logrus.TextFormatter{
	DisableColors:             false,
	ForceColors:               true,
	FullTimestamp:             true,
	EnvironmentOverrideColors: true,
	TimestampFormat:           "2006-01-02 15:04:05.999",
}

type Logger struct {
	*logrus.Logger
}

// 新建一个日志文件
func New() *Logger {
	ret := &Logger{
		Logger: logrus.New(),
	}
	format := DefaultColorTextFormatter
	ret.Logger.SetFormatter(&format)
	ret.Logger.AddHook(&lineHook{
		Field: "source",
		Skip:  3,
	})
	ret.Logger.SetLevel(logrus.TraceLevel)
	return ret
}

// 日志输入到文件
func (l *Logger) ToFile(logPath string, maxDay uint, format logrus.Formatter) {
	l.Logger.SetFormatter(format)
	newFileWriter(logPath, maxDay, l.Logger)
}

// 附加输出到文件
func (l *Logger) HookToFile(logPath string, maxDay uint, format logrus.Formatter) {
	l.Logger.AddHook(newFileHook(logPath, maxDay, l.Logger, format))
}

// 输入到elasticsearch https://github.com/sohlich/elogrus
// 输入到logstash https://github.com/bshuster-repo/logrus-logstash-hook

// 默认的日志
var defaultLog = New()

// 获取默认日志中
func GetDefaultLogger() *Logger {
	return defaultLog
}

func Tracef(format string, args ...interface{}) {
	defaultLog.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLog.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	defaultLog.Infof(format, args...)
}
func Printf(format string, args ...interface{}) {
	defaultLog.Printf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	defaultLog.Warnf(format, args...)
}
func Warningf(format string, args ...interface{}) {
	defaultLog.Warningf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	defaultLog.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	defaultLog.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	defaultLog.Panicf(format, args...)
}
