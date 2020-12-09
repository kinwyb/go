package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kinwyb/go/socket"

	"github.com/kinwyb/go/logs"
	"github.com/spf13/cobra"
)

var (
	serverPort       int64  //服务器端口
	serverListenAddr string //服务器监听地址
	serverAddr       string //服务器地址
	aliveTime        int64  //存货时间
	interval         int64  //间隔时间
	userInput        bool   //用户输入
	log              = logs.GetDefaultLogger()
)

var server = &cobra.Command{
	Use:   "server",
	Short: "tcp服务端测试工具",
	Long:  `tcp服务端测试工具,收到的消息会原样返回`,
	//Args:  cobra.MinimumNArgs(1),
	Run: serverRun,
}

var client = &cobra.Command{
	Use:   "client",
	Short: "tcp客户端测试工具",
	Long:  `tcp客户端测试工具,按interval设置的时间间隔往服务器发送数据,interval默认值1s`,
	//Args:  cobra.MinimumNArgs(1),
	Run: clientRun,
}

func serverRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	config := &socket.TcpServerConfig{
		Port:          serverPort,
		ServerAddress: serverListenAddr,
		Log:           log,
		ErrorHandler: func(errType socket.ErrorType, err error, clientID ...string) {
			log.Errorf("%s => %s", errType.String(), err.Error())
		},
		NewClientHandler: func(clientID string) socket.TcpProtocol {
			log.Infof("新客户端:%s", clientID)
			return nil
		},
		CloseHandler: func() {
			log.Infof("服务端关闭")
		},
		ClientCloseHandler: func(clientID string) {
			log.Infof("客户端[%s]连接关闭", clientID)
		},
	}
	server, err := socket.NewTcpServer(ctx, config)
	config.MessageHandler = func(clientID string, msg []byte) {
		log.Infof("[%s]收到消息:%s", clientID, string(msg))
		server.Write(clientID, msg)
	}
	if err != nil {
		log.Errorf("服务端创建失败:%s", err.Error())
		return
	}
	if aliveTime > 0 {
		go func(f func()) {
			time.Sleep(time.Duration(aliveTime) * time.Minute)
			f()
		}(cancel)
	}
	server.Listen()
}

func clientRun(cmd *cobra.Command, args []string) {
	if serverAddr == "" {
		log.Error("服务器地址不能为空")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	client, err := socket.NewTcpClient(ctx, &socket.TcpClientConfig{
		ServerAddress:     serverAddr,
		AutoReConnect:     true,
		Log:               log,
		ReConnectWaitTime: 5 * time.Second,
		ErrorHandler: func(errorType socket.ErrorType, e error) {
			log.Errorf("%s => %s", errorType, e.Error())
		},
		MessageHandler: func(msg []byte) {
			if userInput {
				fmt.Printf("收到消息: %s", string(msg))
				fmt.Print("> ")
			} else {
				log.Infof("收到消息: %s", string(msg))
			}
		},
		CloseHandler: func() {
			log.Info("连接关闭")
		},
		ConnectTimeOut: 30 * time.Second,
	})
	if err != nil {
		log.Errorf("客户端创建失败:%s", err.Error())
		return
	}
	if aliveTime > 0 {
		go func(f func()) {
			time.Sleep(time.Duration(aliveTime) * time.Second)
			f()
		}(cancel)
	}
	connect := client.Connect()
	log.Infof("服务器连接状态: %v", connect)
	if !connect {
		log.Error("服务器连接失败")
		return
	}
	if interval < 1 {
		interval = 1
	}
	if userInput {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var str = ""
				fmt.Print("> ")
				fmt.Scanln(&str)
				if str == "\\exit" {
					os.Exit(0)
				}
				client.Write([]byte(fmt.Sprintf("%s", str)))
			}
		}
	} else {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
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
}

func main() {
	// 服务端命令解析
	server.Flags().Int64VarP(&serverPort, "port", "p", 1222, "服务监听端口号")
	server.Flags().StringVarP(&serverListenAddr, "address", "a", "127.0.0.1", "服务器监听地址")
	server.Flags().Int64Var(&aliveTime, "alive", 0, "持续时间(s)默认0表示一直运行")
	// 客户端命令解析
	client.Flags().StringVarP(&serverAddr, "server", "s", "127.0.0.1:1222", "服务器地址")
	client.Flags().Int64VarP(&interval, "interval", "i", 1, "消息发送间隔(s),默认:1s")
	client.Flags().Int64Var(&aliveTime, "alive", 0, "持续时间(s)默认0表示一直运行")
	client.Flags().BoolVar(&userInput, "input", false, "手工输入")

	var rootCmd = &cobra.Command{Use: "tcpTool"}
	rootCmd.AddCommand(server, client)
	rootCmd.Execute()
}
