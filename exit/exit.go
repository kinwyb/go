package exit

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//退出函数定义
type Func func(args ...interface{})

//创建监听退出chan
var exitSigle chan os.Signal

var waitGroup *sync.WaitGroup
var cancel context.CancelFunc
var ctx context.Context

//增加退出监控
func Listen(fun Func, args ...interface{}) {
	if fun == nil {
		return
	}
	if exitSigle == nil {
		//创建监听退出chan
		exitSigle = make(chan os.Signal)
		waitGroup = &sync.WaitGroup{}
		ctx, cancel = context.WithCancel(context.TODO())
		go func() {
			<-exitSigle
			cancel() //关闭
		}()
		go func() {
			<-ctx.Done()
			waitGroup.Wait()
			os.Exit(0)
		}()
		//监听指定信号 ctrl+c kill
		signal.Notify(exitSigle,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGUSR1,
			syscall.SIGUSR2)
	}
	waitGroup.Add(1)
	go func(fun Func, args ...interface{}) {
		<-ctx.Done()
		fun(args...)
		waitGroup.Done()
	}(fun, args...)
}
