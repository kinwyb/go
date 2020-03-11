// +build !windows,!nacl,!plan9

package logs

import (
	"github.com/sirupsen/logrus/hooks/syslog"
	slog "log/syslog"
)

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
	l.EnableSource()
	l.SetOutput(hook.Writer)
	l.ExitFunc = func(code int) {
		hook.Writer.Close()
	}
	return nil
}
