package socket_deprecated

import (
	"context"
	"net"
	"time"

	"github.com/kinwyb/go/socket"

	"errors"

	"fmt"

	"github.com/kinwyb/go/logs"
)

type TcpClient struct {
	socketError   chan *socket.Error //连接错误队列
	addr          string             //服务器地址
	sctx          context.Context    //context
	ctx           context.Context    //context
	ncancelFunc   context.CancelFunc //关闭函数
	protocol      *Protocol          //沾包处理
	conn          net.Conn           //socket连接
	lg            logs.Logger        //日志
	IsClose       <-chan struct{}    //关闭channel
	isConnect     chan bool          //是否连接
	reConnectChan chan bool          //重连标记
	reConnect     bool               //重连
	reConnectTime time.Duration      //重连时间
	connectSucc   bool               //连接成功
	doClose       bool               //关闭操作
}

//新客户端对象,注意输出channel error
//否则错误阻塞可能会到这无法正常运行
func NewTcpClient(ctx context.Context, addr string) *TcpClient {
	ret := &TcpClient{
		socketError:   make(chan *socket.Error, 100),
		addr:          addr,
		sctx:          ctx,
		protocol:      NewProtocol(1000),
		lg:            logs.NewLogger(),
		isConnect:     make(chan bool),
		reConnectChan: make(chan bool),
		reConnectTime: 3 * time.Second,
	}
	go ret.reConnGoroutine()
	return ret
}

//连接服务器,该方法会阻塞知道连接异常或ctx关闭
func (c *TcpClient) connectServer() {
	defer recoverPainc(c.connectServer)
	c.ctx, c.ncancelFunc = context.WithCancel(c.sctx)
	var err error
	c.doClose = false
	c.protocol.Reset()
	c.protocol.SetHeartBeat(heartBeatBytes) //设置心跳包内容
	c.conn, err = net.DialTimeout("tcp", c.addr, c.reConnectTime)
	if err != nil {
		c.connectSucc = false
		c.isConnect <- false
		c.lg.Error("服务器连接失败:%s", err.Error())
		c.socketError <- socket.NewError(socket.ConnectErr, err)
		c.reConnectChan <- true //发起重连接
		return
	}
	c.connectSucc = true
	c.isConnect <- true
	c.lg.Info("服务器连接成功...")
	c.IsClose = c.ctx.Done()
	go c.heartbeat() //心跳...
	c.readData()     //接受数据
	c.Close()
}

// 设置是否重连
func (c *TcpClient) SetReConn(reconn bool, reConnectTime time.Duration) {
	c.reConnect = reconn
	c.reConnectTime = reConnectTime
}

//连接服务器返回连接结果是否成功
func (c *TcpClient) Connect() <-chan bool {
	go c.connectServer()
	return c.isConnect
}

//读取数据
func (c *TcpClient) readData() {
	data := make([]byte, 1024)
	for {
		i, err := c.conn.Read(data)
		if c.doClose {
			return
		} else if err != nil {
			c.lg.Error("数据读取错误:" + err.Error())
			c.socketError <- socket.NewError(socket.ReadErr, err)
			c.Close()
			c.reConnectChan <- true //发起重连接
			return
		}
		c.protocol.Unpack(data[0:i])
	}
}

func (c *TcpClient) reConnGoroutine() {
	for {
		<-c.reConnectChan
		if c.connectSucc || !c.reConnect { //连接成功或者不需要重连的直接返回
			continue
		}
		t := time.NewTimer(c.reConnectTime)
		<-t.C
		go c.connectServer()
		<-c.isConnect
	}
}

//发送数据
func (c *TcpClient) Write(data []byte) {
	if c.conn == nil {
		c.socketError <- socket.NewError(socket.SendErr, errors.New("连接已经关闭"))
		return
	}
	_, err := c.conn.Write(c.protocol.Packet(data))
	if err != nil {
		c.socketError <- socket.NewError(socket.SendErr, err)
	}
}

//获取接收到的消息channel
func (c *TcpClient) Msg() <-chan []byte {
	return c.protocol.Read()
}

//获取错误channel
func (c *TcpClient) Error() <-chan *socket.Error {
	return c.socketError
}

//关闭连接
func (c *TcpClient) Close() {
	if c.ncancelFunc != nil {
		c.ncancelFunc()
		c.ncancelFunc = nil
	}
	c.connectSucc = false
	c.doClose = true
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.socketError <- socket.NewError(socket.CancelErr, fmt.Errorf("关闭"))
	}
}

//设置日志
func (c *TcpClient) SetLogger(lg logs.Logger) {
	c.lg = lg
}

//心跳包
func (c *TcpClient) heartbeat() {
	defer recoverPainc(c.heartbeat)
	t := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-t.C:
			_, err := c.conn.Write(c.protocol.Packet(heartBeatBytes))
			if err != nil {
				c.socketError <- socket.NewError(socket.SendErr, errors.New("心跳发送失败"))
				c.Close()   //关闭之前连接
				c.Connect() //重新连接
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//异常处理
func recoverPainc(f ...func()) {
	if r := recover(); r != nil {
		if f != nil {
			for _, fu := range f {
				go fu()
			}
		}
	}
}
