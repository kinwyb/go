package logs

import (
	"testing"
	"time"
)

func TestNewLogFiles(t *testing.T) {
	logfils := NewLogFiles("./", time.Hour)
	logfils.Info("info", "Info测试:%s", "info")
	logfils.Warning("info", "Info测试:%s", "warning")
	logfils.Error("info", "Info测试:%s", "Error")
	logfils.Critical("info", "Info测试:%s", "Critical")
	logfils.Alert("info", "Info测试:%s", "Alert")
	logfils.Emergency("info", "Info测试:%s", "Emergency")
	logfils.Info("Emergency", "Emergency测试:%s", "Info")
	logfils.Error("Emergency", "Emergency测试:%s", "Error")
}
