package exit

import (
	"os"
	"os/signal"
	"syscall"
)

//退出函数定义
type Func func(args ...interface{})

//创建监听退出chan
var exitSigle chan os.Signal
var signalType []os.Signal
var exitFuncs []*exitFun

// 初始化基础信号
func initSignal() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL}
}

//增加退出监控
func Listen(fun Func, args ...interface{}) {
	if fun == nil {
		return
	} else if len(signalType) < 1 {
		signalType = initSignal()
	}
	if exitSigle == nil {
		//创建监听退出chan
		exitSigle = make(chan os.Signal)
		//监听指定信号 ctrl+c kill
		signal.Notify(exitSigle, signalType...)
		go func() {
			<-exitSigle
			for _, v := range exitFuncs {
				if v.fun != nil {
					v.fun(v.args...)
				}
			}
			os.Exit(0)
		}()
	}
	exitFuncs = append(exitFuncs, &exitFun{
		fun:  fun,
		args: args,
	})
}

func Exit() {
	if exitSigle != nil {
		// 发送停止信号
		exitSigle <- syscall.SIGQUIT
		// 等待停止函数完成
		<-exitSigle
	}
}

// 退出函数结构
type exitFun struct {
	fun  Func
	args []interface{}
}
