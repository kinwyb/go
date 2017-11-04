package logs

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	logger.Info("1.Info log ...")
	logger.Warning("2.Warning log ...")
	logger.Error("3.Error log ...")
	logger.Critical("4.Critical ...")
	logger.Alert("6.Alert log ...")
	logger.Emergency("7.Emergency log ...")
}

func TestNewBeegoLogger(t *testing.T) {
	logger := NewBeegoLogger()
	logger.Info("1.Info log ...")
	logger.Warning("2.Warning log ...")
	logger.Error("3.Error log ...")
	logger.Critical("4.Critical ...")
	logger.Alert("6.Alert log ...")
	logger.Emergency("7.Emergency log ...")
}

func TestNewFileLogger(t *testing.T) {
	logger := NewFileLogger("/Users/heldiam/Desktop/test.log", time.Minute)
	ter := time.NewTicker(10 * time.Second)
	i := 1
	for {
		<-ter.C
		fmt.Printf("输出....\n")
		logger.Info("%d.Info log ...", i)
		logger.Warning("%d.Warning log ...", i)
		logger.Error("%d.Error log ...", i)
		logger.Critical("%d.Critical ...", i)
		logger.Alert("%d.Alert log ...", i)
		logger.Emergency("%d.Emergency log ...", i)
		i++
		if i > 30 {
			break
		}
	}
}

func BenchmarkNewLogger(b *testing.B) {
	//10000	     44039 ns/op
	//20000	     39594 ns/op
	//30000	     41878 ns/op
	logger := NewLogger()
	b.N = 30000
	//b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("1.Info log ...")
			logger.Warning("2.Warning log ...")
			logger.Error("3.Error log ...")
			logger.Critical("4.Critical ...")
			logger.Alert("6.Alert log ...")
			logger.Emergency("7.Emergency log ...")
		}
	})
}

func BenchmarkNewBeegoLogger(b *testing.B) {
	//10000	    179570 ns/op
	//20000	     96171 ns/op
	//30000	     94612 ns/op
	logger := NewBeegoLogger()
	b.N = 30000
	//b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("1.Info log ...")
			logger.Warning("2.Warning log ...")
			logger.Error("3.Error log ...")
			logger.Critical("4.Critical ...")
			logger.Alert("6.Alert log ...")
			logger.Emergency("7.Emergency log ...")
		}
	})
}
