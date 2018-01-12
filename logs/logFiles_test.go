package logs

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewLogFiles(t *testing.T) {
	logfils := NewLogFiles("/Users/heldiam/Desktop/", time.Minute)
	for {
		logfils.Info("info", "Info测试:%s", "info")
		logfils.Warning("info", "Info测试:%s", "warning")
		logfils.Error("info", "Info测试:%s", "Error")
		logfils.Critical("info", "Info测试:%s", "Critical")
		logfils.Alert("info", "Info测试:%s", "Alert")
		logfils.Emergency("info", "Info测试:%s", "Emergency")
		logfils.Info("Emergency", "Emergency测试:%s", "Info")
		logfils.Error("Emergency", "Emergency测试:%s", "Error")
		time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
	}
}

func TestTimer(t *testing.T) {
	t1 := time.NewTimer(3 * time.Second)
	for {
		<-t1.C
		println("time down ")
		t1.Reset(1 * time.Second)
	}
}
