package exit

import (
	"fmt"
	"testing"
	"time"
)

func TestExit(t *testing.T) {
	exitFun1 := func(args ...interface{}) {
		fmt.Println("退出函数1")
	}
	exitFun2 := func(args ...interface{}) {
		fmt.Println("退出函数2")
	}
	exitFun3 := func(args ...interface{}) {
		fmt.Println("退出函数3")
	}
	exitFun4 := func(args ...interface{}) {
		fmt.Println("退出函数4")
	}
	exitFun5 := func(args ...interface{}) {
		fmt.Println("退出函数5")
	}
	Listen(exitFun1)
	Listen(exitFun2)
	Listen(exitFun3)
	Listen(exitFun4)
	Listen(exitFun5)
	time.Sleep(10 * time.Second)
	t.Logf("开启关闭")
}
