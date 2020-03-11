package logs

import (
	"sync"
	"testing"
	"time"
)

func Test_logConsole(t *testing.T) {
	l := New()
	format := DefaultTextFormatter
	format.ForceColors = true
	format.DisableColors = false
	format.FullTimestamp = true
	format.EnvironmentOverrideColors = true
	l.SetFormatter(&format)
	l.WithField("dd", "vv").Info("123456")
	l.WithField("dd", "vv").Error("123456")
	l.WithField("dd", "vv").Debug("123456")
	l.WithField("dd", "vv").Fatal("123456")
}

func Test_logFile(t *testing.T) {
	format := DefaultJsonFormatter
	l := New()
	l.ToFile("/Users/heldiam/Desktop/logs/log", 7, &format)
	l2 := New()
	l2.ToFile("/Users/heldiam/Desktop/logs/log", 7, &format)
	w := sync.WaitGroup{}
	w.Add(2)
	go logToFile(l, "L1", &w)
	go logToFile(l2, "L2", &w)
	w.Wait()
}

func logToFile(l *Logger, tag string, w *sync.WaitGroup) {
	i := 0
	t := time.NewTicker(50 * time.Microsecond)
	for {
		<-t.C
		i++
		l.Infof("%s : %d", tag, i)
		if i > 2000 {
			break
		}
	}
	w.Done()
}
