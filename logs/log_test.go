package logs

import "testing"

func Test_logConsole(t *testing.T) {
	l := New()
	l.EnableTextFormat(true)
	l.WithField("dd", "vv").Info("123456")
	l.WithField("dd", "vv").Error("123456")
	l.WithField("dd", "vv").Debug("123456")
	l.WithField("dd", "vv").Fatal("123456")
}

func Test_logFile(t *testing.T) {
	l := New()
	l.EnableJsonFormat()
	l.ToFile("/Users/heldiam/Desktop/logs/log", 7)
	l.WithField("dd", "vv").Info("123456")
	l.WithField("dd", "vv").Error("123456")
	l.WithField("dd", "vv").Debug("123456")
	l.WithField("dd", "vv").Fatal("123456")
}
