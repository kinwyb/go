package logs

import (
	"github.com/sirupsen/logrus"
)

var DefaultJsonFormatter = logrus.JSONFormatter{
	TimestampFormat: "2006-01-02 15:04:05.999",
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyLevel: "@level",
		logrus.FieldKeyMsg:   "@message",
		logrus.FieldKeyFunc:  "@caller",
	},
}
var DefaultTextFormatter = logrus.TextFormatter{
	DisableColors:   true,
	TimestampFormat: "2006-01-02 15:04:05.999",
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
	enableSource bool
}

// 新建一个日志文件
func New() *Logger {
	ret := &Logger{
		Logger:       logrus.New(),
		enableSource: false,
	}
	format := DefaultTextFormatter
	ret.Logger.SetFormatter(&format)
	ret.Logger.SetLevel(logrus.TraceLevel)
	return ret
}

func (l *Logger) EnableSource() {
	if l.enableSource {
		return
	}
	l.Logger.SetReportCaller(true)
}

// 日志输入到文件
func (l *Logger) ToFile(logPath string, maxDay uint, format logrus.Formatter) {
	l.Logger.SetFormatter(format)
	l.EnableSource() //只写入文件时候
	newFileWriter(logPath, maxDay, l.Logger)
}

// 附加输出到文件
func (l *Logger) HookToFile(logPath string, maxDay uint, format logrus.Formatter) {
	l.Logger.AddHook(newFileHook(logPath, maxDay, l.Logger, format))
}

func (l *Logger) RemoveHook(h logrus.Hook) {
	newHooks := logrus.LevelHooks{}
	// 因为PanicLevel最低级的,这个等级是包含所有hook的,遍历这个等级即可
	for _, hooks := range l.Hooks[logrus.PanicLevel] {
		if hooks != h {
			newHooks.Add(hooks)
		}
	}
	l.ReplaceHooks(newHooks)
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

type ILogger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}
