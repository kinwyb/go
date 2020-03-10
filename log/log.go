package log

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/syslog"
	slog "log/syslog"
)

type Logger struct {
	*logrus.Logger
}

// 新建一个日志文件
func New() *Logger {
	ret := &Logger{
		Logger: logrus.New(),
	}
	ret.Logger.AddHook(&lineHook{
		Field: "source",
	})
	ret.Logger.SetLevel(logrus.TraceLevel)
	return ret
}

// 日志输入到文件
func (l *Logger) ToFile(logPath string, maxDay uint) {
	newFileWriter(logPath, maxDay, l.Logger)
}

// 附加输出到文件
func (l *Logger) HookToFile(logPath string, maxDay uint) {
	l.Logger.AddHook(newFileHook(logPath, maxDay, l.Logger))
}

func (l *Logger) EnableJsonFormat() {
	l.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyFunc:  "@caller",
		},
	})
}

func (l *Logger) EnableTextFormat(disableColor bool) {
	l.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     !disableColor,
		DisableColors:   disableColor,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// 日志输入到syslog
func (l *Logger) HookToSysLog(network, raddr string, priority slog.Priority, tag string) error {
	hook, err := syslog.NewSyslogHook(network, raddr, priority, tag)
	if err != nil {
		l.Error("Unable to connect to local syslog daemon")
		return err
	}
	l.AddHook(hook)
	return nil
}

func (l *Logger) ToSysLog(network, raddr string, priority slog.Priority, tag string) error {
	hook, err := syslog.NewSyslogHook(network, raddr, priority, tag)
	if err != nil {
		l.Error("Unable to connect to local syslog daemon")
		return err
	}
	l.SetOutput(hook.Writer)
	l.ExitFunc = func(code int) {
		hook.Writer.Close()
	}
	return nil
}

// 输入到elasticsearch https://github.com/sohlich/elogrus
// 输入到logstash https://github.com/bshuster-repo/logrus-logstash-hook
