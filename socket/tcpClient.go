package socket

import (
	"context"
	"net"
	"time"

	"errors"

	"github.com/kinwyb/go/logs"
)

type TcpClient struct {
	socketError chan error         //连接错误队列
	addr        string             //服务器地址
	ctx         context.Context    //
	nctx        context.Context    //内部的context
	ncancelFunc context.CancelFunc //关闭函数
	protocol    *Protocol          //沾包处理
	conn        net.Conn           //socket连接
	lg          logs.Logger
}

//新客户端对象,注意输出channel error
//否则错误阻塞可能会到这无法正常运行
func NewTcpClient(ctx context.Context, addr string) *TcpClient {
	return &TcpClient{
		socketError: make(chan error, 100),
		ctx:         ctx,
		addr:        addr,
		protocol:    NewProtocol(1000),
		lg:          logs.NewLogger(),
	}
}

//连接服务器,该方法会阻塞知道连接异常或ctx关闭
//建议使用go ConnectServer()方式运行
func (c *TcpClient) ConnectServer() {
	defer recoverPainc(c.ConnectServer)
	var err error
	c.protocol.SetHeartBeat(heartBeatBytes) //设置心跳包内容
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		c.lg.Error("服务器连接失败:%s", err.Error())
		c.socketError <- err
		return
	}
	c.nctx, c.ncancelFunc = context.WithCancel(c.ctx)
	c.lg.Info("服务器连接成功...")
	c.protocol = NewProtocol(1000)
	go c.readData()
	go c.heartbeat() //心跳线程...
	select {
	case <-c.ctx.Done():
	case <-c.nctx.Done():
	}
	c.conn.Close()
}

//读取数据
func (c *TcpClient) readData() {
	defer recoverPainc(c.readData)
	data := make([]byte, 1024)
	for {
		select {
		case <-c.ctx.Done():
		case <-c.nctx.Done():
			//这里显示结束
			return
		default:
			i, err := c.conn.Read(data)
			if err != nil {
				c.lg.Error("数据读取错误:" + err.Error())
				c.socketError <- err
				c.ncancelFunc()
				go c.ConnectServer()
				return
			}
			c.protocol.Unpack(data[0:i])
		}
	}
}

//发送数据
func (c *TcpClient) Write(data []byte) {
	_, err := c.conn.Write(c.protocol.Packet(data))
	if err != nil {
		c.socketError <- errors.New("心跳发送失败")
		c.ncancelFunc()      //关闭连接
		go c.ConnectServer() //重新开启连接
	}
}

//获取接收到的消息channel
func (c *TcpClient) Msg() <-chan []byte {
	return c.protocol.Read()
}

//获取错误channel
func (c *TcpClient) Error() <-chan error {
	return c.socketError
}

//关闭连接
func (c *TcpClient) Close() {
	if c.ncancelFunc != nil {
		c.ncancelFunc()
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
			_, err := c.conn.Write(c.protocol.Packet([]byte("heartbeat")))
			if err != nil {
				c.socketError <- errors.New("心跳发送失败")
				c.ncancelFunc()      //关闭连接
				go c.ConnectServer() //重新开启连接
			}
		case <-c.ctx.Done():
		case <-c.nctx.Done():
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
