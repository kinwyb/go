package socket

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/kinwyb/go/logs"
)

// tcp客户端配置
type TcpClientConfig struct {
	ServerAddress     string                 //服务地址
	AutoReConnect     bool                   //自动重连
	ReConnectWaitTime time.Duration          //重连等待时间
	Protocol          TcpProtocol            //连接处理协议,如果没有则收到什么数据返回什么数据,写入什么数据发送什么数据
	Log               logs.ILogger           //日志
	ErrorHandler      func(ErrorType, error) //错误处理
	ConnectHandler    func()                 //连接成功调用
	MessageHandler    func([]byte)           //消息处理
	CloseHandler      func()                 //连接关闭处理
	HeartBeatData     []byte                 //心跳内容,如果为空不发送心跳信息
	HeartBeatTime     time.Duration          //心跳间隔
	ConnectTimeOut    time.Duration          //连接超时事件
}

// tcp客户端
type TcpClient struct {
	ctx           context.Context  //上下文context
	conn          net.Conn         //socket连接
	config        *TcpClientConfig //配置
	isConnect     chan bool        //是否连接
	reConnectChan chan bool        //重连出发器
	connectSucc   bool             //连接成功
	doClose       bool             //关闭操作
}

//创建Tcp客户端
func NewTcpClient(ctx context.Context, config *TcpClientConfig) (*TcpClient, error) {
	if config == nil {
		return nil, errors.New("配置信息不能为空")
	}
	ret := &TcpClient{
		config:        config,
		ctx:           ctx,
		isConnect:     make(chan bool, 1),
		reConnectChan: make(chan bool, 1),
	}
	go ret.processMsg() //解析数据
	go ret.stateCheckGoroutine()
	return ret, nil
}

//连接服务器,该方法会阻塞知道连接异常或ctx关闭
func (c *TcpClient) connectServer() {
	defer recoverPainc(c.config.Log, c.connectServer)
	if c.connectSucc {
		return
	}
	c.doClose = false
	var err error
	c.conn, err = net.DialTimeout("tcp",
		c.config.ServerAddress, c.config.ConnectTimeOut)
	if err != nil {
		c.connectSucc = false
		c.isConnect <- false
		if c.config.ErrorHandler != nil {
			c.config.ErrorHandler(ConnectErr, err)
		}
		c.reConnectChan <- true //发起重连接
		return
	}
	c.connectSucc = true
	go c.heartbeat() //心跳...
	go c.readData()  //接受数据
	if c.config.ConnectHandler != nil {
		c.config.ConnectHandler()
	}
	c.isConnect <- true
	if c.config.Log != nil {
		c.config.Log.Info("服务器连接成功")
	}
}

//连接服务器返回连接结果是否成功
func (c *TcpClient) ConnectAsync() <-chan bool {
	go c.connectServer()
	return c.isConnect
}

//连接服务器返回连接结果是否成功
func (c *TcpClient) Connect() bool {
	c.connectServer()
	return <-c.isConnect
}

//读取数据
func (c *TcpClient) readData() {
	data := make([]byte, 1024)
	for {
		if c.conn == nil || c.doClose {
			return
		}
		i, err := c.conn.Read(data)
		if c.doClose {
			return
		}
		if err != nil {
			if err != io.EOF && c.config.ErrorHandler != nil {
				c.config.ErrorHandler(ReadErr, err)
			}
			c.connectClose()
			c.reConnectChan <- true //发起重连接
			return
		}
		if c.config.Protocol != nil {
			c.config.Protocol.UnPacking(data[0:i])
		} else if c.config.MessageHandler != nil {
			c.config.MessageHandler(data[0:i])
		}
	}
}

// 状态检测线程
func (c *TcpClient) stateCheckGoroutine() {
	defer recoverPainc(c.config.Log, c.stateCheckGoroutine)
	for {
		select {
		case <-c.reConnectChan:
			if c.connectSucc || !c.config.AutoReConnect { //连接成功或者不需要重连的直接返回
				continue
			}
			time.Sleep(c.config.ReConnectWaitTime)
			if c.doClose {
				return
			}
			c.Connect()
		case <-c.ctx.Done(): //如果上下文结束,关闭整个连接
			c.Close()
			return
		}

	}
}

//发送数据
func (c *TcpClient) Write(data []byte) error {
	if c.conn == nil {
		return errors.New("未连接服务器")
	}
	_, err := c.conn.Write(c.packingData(data))
	return err
}

//编码数据
func (c *TcpClient) packingData(msg []byte) []byte {
	if c.config.Protocol != nil {
		return c.config.Protocol.Packing(msg)
	}
	return msg
}

//解析返回收到的有效内容
func (c *TcpClient) processMsg() {
	defer recoverPainc(c.config.Log, c.processMsg)
	if c.config.Protocol != nil {
		for {
			v := <-c.config.Protocol.ReadMsg()
			if c.config.MessageHandler != nil {
				c.config.MessageHandler(v)
			}
		}
	}
}

//关闭连接
func (c *TcpClient) Close() {
	c.doClose = true
	c.connectClose()
}

func (c *TcpClient) connectClose() {
	if c.connectSucc {
		c.connectSucc = false
		if c.config.CloseHandler != nil {
			c.config.CloseHandler()
		}
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

//心跳包
func (c *TcpClient) heartbeat() {
	defer recoverPainc(c.config.Log, c.heartbeat)
	if c.config.HeartBeatData == nil { //没有心跳需求
		return
	}
	t := time.NewTicker(c.config.HeartBeatTime)
	for {
		select {
		case <-t.C:
			if c.conn == nil {
				continue
			}
			_, err := c.conn.Write(c.packingData(c.config.HeartBeatData))
			if err != nil {
				if c.config.ErrorHandler != nil {
					c.config.ErrorHandler(SendErr, errors.New("心跳发送失败"))
				}
				c.Close()   //关闭之前连接
				c.Connect() //重新连接
			}
		case <-c.ctx.Done():
			return
		}
	}
}
