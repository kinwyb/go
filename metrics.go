package heldiamgo

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kinwyb/go/logs"

	"github.com/shirou/gopsutil/process"

	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/vrischmann/go-metrics-influxdb"
)

var EnableMetrics = false //是否开启性能统计
var MetricsRegistry metrics.Registry
var MetricsCancel context.CancelFunc

var log logs.ILogger = logs.GetDefaultLogger()

//开启性能统计
func StartMetrics(influxdbURL string, influxdbName string,
	influxdbUser string, influxdbPwd string) {
	if influxdbURL == "" || influxdbName == "" { //时序数据库地址错误不启动
		return
	}
	EnableMetrics = true
	MetricsRegistry = metrics.NewRegistry()
	metrics.RegisterDebugGCStats(MetricsRegistry)
	metrics.RegisterRuntimeMemStats(MetricsRegistry)
	go metrics.CaptureDebugGCStats(MetricsRegistry, time.Second*5)
	go metrics.CaptureRuntimeMemStats(MetricsRegistry, time.Second*5)
	go influxdb.InfluxDB(MetricsRegistry, time.Second, influxdbURL,
		influxdbName, "go-metrics", influxdbUser, influxdbPwd, true)
	ctx := context.Background()
	ctx, MetricsCancel = context.WithCancel(ctx)
	go runInfo(MetricsRegistry, time.Second*5, ctx) //运行时执行的数据
}

// 爆露性能统计http接口
func MetricExpHandler() http.Handler {
	if MetricsRegistry == nil {
		MetricsRegistry = metrics.NewRegistry()
	}
	return exp.ExpHandler(MetricsRegistry)
}

// 设置日志
func SetLogger(lg logs.ILogger) {
	log = lg
}

//运行时的数据
func runInfo(r metrics.Registry, d time.Duration, ctx context.Context) {
	ps, err := process.Processes()
	if err != nil {
		log.Errorf("运行数据获取失败:%s", err.Error())
		return
	}
	var proc *process.Process
	for _, v := range ps {
		if cmd, err := v.Cmdline(); err == nil {
			if cmd == os.Args[0] {
				proc = v
				break
			}
		}
	}
	if proc == nil {
		log.Error("运行进程ID获取失败")
		return
	}
	log.Infof("获取到运行进程ID:%d", proc.Pid)
	memRSS := metrics.NewGauge()
	memVMS := metrics.NewGauge()
	r.Register("process.Mem.RSS", memRSS)
	r.Register("process.Mem.VMS", memVMS)
	infs, err := net.Interfaces()
	if err != nil {
		log.Error("系统网络获取失败")
	}
	netRecvGauge := map[string]metrics.Gauge{}
	netSentGauge := map[string]metrics.Gauge{}
	netLastRecvBytes := map[string]uint64{}
	netLastSentBytes := map[string]uint64{}
	for _, v := range infs {
		addr, _ := v.Addrs()
		if len(addr) < 1 {
			continue
		}
		addrString := addr[0].String()
		netRecvGauge[v.Name] = metrics.NewGauge()
		netSentGauge[v.Name] = metrics.NewGauge()
		r.Register(fmt.Sprintf("process.net.recv.%s", addrString), netRecvGauge[v.Name])
		r.Register(fmt.Sprintf("process.net.sent.%s", addrString), netSentGauge[v.Name])
	}
	dBits := uint64(d / time.Second)
	for {
		select {
		case <-time.Tick(d):
			//内存
			mem, err := proc.MemoryInfo()
			if err != nil {
				log.Errorf("进程运行内存获取失败:%s", err.Error())
			} else {
				memRSS.Update(int64(mem.RSS))
				memVMS.Update(int64(mem.VMS))
			}
			//网络
			netInfo, err := proc.NetIOCounters(true)
			if err == nil {
				for _, v := range netInfo {
					if gauge, ok := netRecvGauge[v.Name]; ok {
						recv := netLastRecvBytes[v.Name]
						if recv < 1 {
							netLastRecvBytes[v.Name] = v.BytesRecv
						} else {
							bits := v.BytesRecv - recv
							bits = bits / dBits
							gauge.Update(int64(bits))
							netLastRecvBytes[v.Name] = v.BytesRecv
						}
					}
					if gauge, ok := netSentGauge[v.Name]; ok {
						sent := netLastSentBytes[v.Name]
						if sent < 1 {
							netLastSentBytes[v.Name] = v.BytesSent
						} else {
							bits := v.BytesSent - sent
							bits = bits / dBits
							gauge.Update(int64(bits))
							netLastSentBytes[v.Name] = v.BytesSent
						}
					}
				}
			}
		case <-ctx.Done():
			goto end
		}
	}
end:
	return
}

//结束性能统计
func StopMetrics() {
	EnableMetrics = false
	MetricsRegistry.UnregisterAll()
	MetricsRegistry = nil
	MetricsCancel() //关闭
}
