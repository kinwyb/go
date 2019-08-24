package socket

import (
	"context"
	"testing"

	"github.com/kinwyb/go/logs"
)

func Test_NewTcpServer(t *testing.T) {
	log := logs.NewLogger()
	ctx, _ := context.WithCancel(context.Background())
	config := &TcpServerConfig{
		Port:          1222,
		ServerAddress: "",
		Log:           log,
		ErrorHandler: func(errType ErrorType, err error, clientID ...string) {
			log.Error("%s => %s", errType.String(), err.Error())
		},
		NewClientHandler: func(clientID string) TcpProtocol {
			log.Info("新客户端:%s", clientID)
			return nil
		},
		CloseHandler: func() {
			log.Info("服务端关闭")
		},
		ClientCloseHandler: func(clientID string) {
			log.Info("客户端[%s]连接关闭", clientID)
		},
	}
	server, err := NewTcpServer(ctx, config)
	config.MessageHandler = func(clientID string, msg []byte) {
		log.Info("[%s]收到消息:%s", clientID, string(msg))
		server.Write(clientID, msg)
	}
	if err != nil {
		t.Fatal(err)
	}
	//go func(f func()) {
	//	time.Sleep(2 * time.Minute)
	//	f()
	//}(cancel)
	server.Listen()
}
