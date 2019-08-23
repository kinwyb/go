package socket_deprecated

import (
	"context"
	"fmt"
	"net"

	"github.com/kinwyb/go/socket"

	"github.com/kinwyb/go"

	"errors"

	"github.com/kinwyb/go/logs"
)

//数据消息
type TcpMsg struct {
	ClientID string `description:"客户端标记"`
	Data     []byte `description:"内容"`
}

// 新连接回调函数
var NewClientCallBackFun func(client *SClient)

// 客户端断开回调函数
var ClientDoneCallBackFun func(clientID string)

//TCP服务器对象
type TcpServer struct {
	socketErr   chan *socket.Error
	ctx         context.Context
	nctx        context.Context
	ncancelFunc context.CancelFunc
	conn        net.Listener
	lg          logs.Logger
	port        int
	clients     map[string]*SClient
	readData    chan *TcpMsg //读取到的数据
	doClose     bool
	IsClose     <-chan struct{} //关闭channel
}

//新建一个服务端,port为监听端口
//注意处理channel error 否则程序阻塞无法正常运行
func NewTcpServer(ctx context.Context, port int) *TcpServer {
	return &TcpServer{
		socketErr: make(chan *socket.Error, 100),
		ctx:       ctx,
		lg:        logs.NewLogger(),
		port:      port,
		clients:   map[string]*SClient{},
		readData:  make(chan *TcpMsg, 1000),
	}
}

//启动监听该方法会阻塞
func (s *TcpServer) Listen() {
	var err error
	s.doClose = false
	s.nctx, s.ncancelFunc = context.WithCancel(s.ctx)
	s.IsClose = s.nctx.Done()
	s.conn, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.lg.Error("查询Socket监听失败:" + err.Error())
		s.socketErr <- socket.NewError(socket.ListenErr, err)
		return
	}
	go s.handleConn()
	select {
	case <-s.ctx.Done():
	case <-s.nctx.Done():
	}
	s.Close() //关闭
}

//关闭服务
func (s *TcpServer) Close() {
	if s.ncancelFunc != nil {
		s.ncancelFunc()
		s.ncancelFunc = nil
	}
	s.doClose = true
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

//设置日志
func (s *TcpServer) SetLogger(lg logs.Logger) {
	s.lg = lg
}

//获取错误channel
func (s *TcpServer) Error() <-chan *socket.Error {
	return s.socketErr
}

//处理新连接
func (s *TcpServer) handleConn() {
	for {
		if conn, err := s.conn.Accept(); err != nil {
			if s.doClose { //关闭
				return
			}
			s.lg.Error("监听请求连接失败:" + err.Error())
			s.socketErr <- socket.NewError(socket.ListenErr, err)
		} else {
			sclient := &SClient{
				conn:     conn,
				server:   s,
				protocol: NewProtocol(1000),
				lg:       s.lg,
				ID:       heldiamgo.IDGen(),
			}
			sclient.ctx, sclient.cancelFunc = context.WithCancel(s.nctx)
			sclient.IsClose = sclient.ctx.Done()
			sclient.doClose = false
			sclient.protocol.SetHeartBeat(heartBeatBytes) //设置心跳
			go sclient.readData()
			s.clients[sclient.ID] = sclient
			if NewClientCallBackFun != nil {
				NewClientCallBackFun(sclient)
			}
		}
	}
}

//发送数据
func (s *TcpServer) Write(msg *TcpMsg) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	client := s.clients[msg.ClientID]
	if client == nil {
		return errors.New("客户端链接不存在")
	}
	_, err := client.conn.Write(client.protocol.Packet(msg.Data))
	if err != nil {
		client.Close()
		return fmt.Errorf("数据发送失败:%s", err.Error())
	}
	return nil
}

//读取消息
func (s *TcpServer) Read() <-chan *TcpMsg {
	return s.readData
}

//TCP连接到服务器的客户端
type SClient struct {
	conn       net.Conn
	server     *TcpServer
	protocol   *Protocol
	ctx        context.Context
	cancelFunc context.CancelFunc
	lg         logs.Logger
	IsClose    <-chan struct{} //关闭channel
	ID         string          //客户端唯一标示
	doClose    bool            //关闭
}

//读取客户端数据
func (s *SClient) readData() {
	defer func() {
		if err := recover(); err != nil {
			s.lg.Error("客户链接[%s]数据读取异常:%s", s.ID, err)
			s.Close()
		}
	}()
	data := make([]byte, 1024)
	for {
		i, err := s.conn.Read(data)
		if s.doClose {
			return
		} else if err != nil {
			s.lg.Error("%s=>数据读取错误:%s", s.ID, err.Error())
			s.server.socketErr <- socket.NewError(socket.ReadErr, fmt.Errorf("%s=>%s", s.ID, err.Error()))
			s.Close()
			return
		}
		s.protocol.Unpack(data[0:i])
		select {
		case data := <-s.protocol.Read():
			//todo: tcpmsg可以做池
			s.server.readData <- &TcpMsg{
				ClientID: s.ID,
				Data:     data,
			}
		default:
		}
	}
}

//发送数据
func (s *SClient) write(msg []byte) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	_, err := s.conn.Write(s.protocol.Packet(msg))
	if err != nil {
		return fmt.Errorf("数据发送失败:%s", err.Error())
	}
	return nil
}

//关闭连接
func (s *SClient) Close() {
	if ClientDoneCallBackFun != nil {
		ClientDoneCallBackFun(s.ID)
	}
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
	s.doClose = true
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
		s.server.socketErr <- socket.NewError(socket.CancelErr, fmt.Errorf("%s", s.ID))
	}
	s.protocol.Reset()
}
