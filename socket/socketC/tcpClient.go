package main

/*
#include "socketCCode.h"
// 连接成功触发函数
void TcpClientConnectCallback(SocketTcpClientConfig *config);

// 收到消息触发函数
void TcpClientMessageCallback(SocketTcpClientConfig *config,char *data,GoInt len);

// 关闭触发函数
void TcpClientCloseCallback(SocketTcpClientConfig *config);

// 错误触发函数
void TcpClientErrorCallback(SocketTcpClientConfig *config,GoInt errCode,char *msg,GoInt msgLen);

*/
import "C"
import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"

	"github.com/kinwyb/go/socket"
)

var decoder *encoding.Decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
var encoder *encoding.Encoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()

var cMap = map[string]*cTcpClient{}

type cTcpClient struct {
	s       unsafe.Pointer
	client  *socket.TcpClient
	config  *socket.TcpClientConfig
	handler *cTcpClientHandler
}

type cTcpClientHandler struct {
	s unsafe.Pointer
}

func (c *cTcpClientHandler) errHandler(code socket.ErrorType, err error) {
	fmt.Printf("错误:%d %s\n", code, err.Error())
	e := err.Error()
	msg := C.CString(e)
	msgLen := C.GoInt(len(e))
	C.TcpClientErrorCallback((*C.SocketTcpClientConfig)(c.s), C.GoInt(code), msg, msgLen)
}

func (c *cTcpClientHandler) msgHandler(msg []byte) {
	v := C.CBytes(msg)
	C.TcpClientMessageCallback((*C.SocketTcpClientConfig)(c.s), (*C.char)(v), C.GoInt(len(msg)))
}

func (c *cTcpClientHandler) connectHandler() {
	fmt.Printf("连接成功\n")
	C.TcpClientConnectCallback((*C.SocketTcpClientConfig)(c.s))
}

func (c *cTcpClientHandler) closeHandler() {
	fmt.Printf("连接关闭\n")
	C.TcpClientCloseCallback((*C.SocketTcpClientConfig)(c.s))
}

//export GoNewTcpClient
func GoNewTcpClient(s unsafe.Pointer, id string, address string, autoReConnect int,
	ReConnectWaitTime int, HeartBeatTime int, ConnectTimeOut int, HeartBeatData []byte) {

	fmt.Printf("初始化连接 %d %d %d %d\n", autoReConnect, ReConnectWaitTime, HeartBeatTime, ConnectTimeOut)
	address, _ = decoder.String(address)
	client := &cTcpClient{
		s:       s,
		config:  nil,
		handler: &cTcpClientHandler{s: s},
	}
	client.config = &socket.TcpClientConfig{
		ServerAddress:     address,
		AutoReConnect:     autoReConnect == 1,
		ReConnectWaitTime: time.Second * time.Duration(ReConnectWaitTime),
		Protocol:          nil,
		Log:               nil,
		ErrorHandler:      client.handler.errHandler,
		ConnectHandler:    client.handler.connectHandler,
		MessageHandler:    client.handler.msgHandler,
		CloseHandler:      client.handler.closeHandler,
		HeartBeatData:     HeartBeatData,
		HeartBeatTime:     time.Second * time.Duration(HeartBeatTime),
		ConnectTimeOut:    time.Second * time.Duration(ConnectTimeOut),
	}
	fmt.Printf("开始连接:%s\n", client.config.ServerAddress)
	c, err := socket.NewTcpClient(context.Background(), client.config)
	if err != nil {
		client.handler.errHandler(-1, fmt.Errorf("客户端创建错误:%s", err.Error()))
		return
	}
	client.client = c
	client.client.Connect()
	cMap[id] = client
}

//export GoTcpClientClose
func GoTcpClientClose(id string) {
	if v, ok := cMap[id]; ok {
		v.client.Close()
	}
}

//export TcpClientWrite
func TcpClientWrite(id string, data []byte) {
	if v, ok := cMap[id]; ok {
		err := v.client.Write(data)
		if err != nil {
			v.handler.errHandler(socket.SendErr, err)
		}
	}
}
