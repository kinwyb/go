package socket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kinwyb/go/logs"
)

func Test_NewTcpClient(t *testing.T) {
	log := logs.GetDefaultLogger()
	ctx, cancel := context.WithCancel(context.Background())
	client, err := NewTcpClient(ctx, &TcpClientConfig{
		ServerAddress: "127.0.0.1:1222",
		AutoReConnect: true,
		Log:           log,
		//Protocol:          NewProtocol(100),
		ReConnectWaitTime: 5 * time.Second,
		ErrorHandler: func(errorType ErrorType, e error) {
			log.Errorf("%s => %s", errorType, e.Error())
		},
		MessageHandler: func(msg []byte) {
			log.Infof("收到消息: %s", string(msg))
		},
		CloseHandler: func() {
			log.Info("连接关闭")
			//cancel()
		},
		ConnectTimeOut: 30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	go func(f func()) {
		time.Sleep(2 * time.Minute)
		f()
	}(cancel)
	connect := client.Connect()
	log.Infof("服务器连接状态: %v", connect)
	if !connect {
		t.Fatalf("服务器连接失败")
	}
	ticker := time.NewTicker(1 * time.Second)
	i := 1
	for {
		select {
		case <-ticker.C:
			client.Write([]byte(fmt.Sprintf("消息:%d", i)))
			i++
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
