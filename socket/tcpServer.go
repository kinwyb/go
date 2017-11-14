package socket

import (
	"context"
	"fmt"
	"net"

	"errors"

	"github.com/kinwyb/go/logs"
)

//TCP服务器对象
type TcpServer struct {
	socketErr   chan Error
	ctx         context.Context
	nctx        context.Context
	ncancelFunc context.CancelFunc
	conn        net.Listener
	lg          logs.Logger
	port        int
	client      chan *SClient
	IsClose     <-chan struct{} //关闭channel
}

//TCP连接到服务器的客户端
type SClient struct {
	conn        net.Conn
	protocol    *Protocol
	ctx         context.Context
	cancelFunc  context.CancelFunc
	lg          logs.Logger
	socketError chan<- Error
	IsClose     <-chan struct{} //关闭channel
	ID          string          //客户端唯一标示
}

//新建一个服务端,port为监听端口
//注意处理channel error 否则程序阻塞无法正常运行
func NewTcpServer(ctx context.Context, port int) *TcpServer {
	return &TcpServer{
		socketErr: make(chan Error, 100),
		ctx:       ctx,
		lg:        logs.NewLogger(),
		port:      port,
		client:    make(chan *SClient, 100),
	}
}

//启动监听该方法会阻塞
func (s *TcpServer) Listen() {
	var err error
	s.nctx, s.ncancelFunc = context.WithCancel(s.ctx)
	s.IsClose = s.nctx.Done()
	s.conn, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.lg.Error("查询Socket监听失败:" + err.Error())
		s.socketErr <- Error{
			t:   Listen,
			err: err,
		}
	}
	go s.handleConn()
	select {
	case <-s.ctx.Done():
	case <-s.nctx.Done():
	}
	s.conn.Close()
}

//获取客户端连接
func (s *TcpServer) Accept() <-chan *SClient {
	return s.client
}

//关闭服务
func (s *TcpServer) Close() {
	if s.ncancelFunc != nil {
		s.ncancelFunc()
	}
}

//设置日志
func (s *TcpServer) SetLogger(lg logs.Logger) {
	s.lg = lg
}

//获取错误channel
func (s *TcpServer) Error() <-chan Error {
	return s.socketErr
}

//处理新连接
func (s *TcpServer) handleConn() {
	for {
		if conn, err := s.conn.Accept(); err != nil {
			s.lg.Error("监听请求连接失败:" + err.Error())
			s.socketErr <- Error{
				t:   Listen,
				err: err,
			}
		} else {
			sclient := &SClient{
				conn:        conn,
				protocol:    NewProtocol(1000),
				socketError: s.socketErr,
				lg:          s.lg,
			}
			sclient.protocol.SetHeartBeat(heartBeatBytes)
			sclient.ctx, sclient.cancelFunc = context.WithCancel(s.ctx)
			sclient.IsClose = sclient.ctx.Done()
			go sclient.readData()
			s.client <- sclient
		}
	}
}

//读取客户端数据
func (s *SClient) readData() {
	defer recoverPainc(s.readData)
	data := make([]byte, 1024)
	for {
		i, err := s.conn.Read(data)
		if err != nil {
			s.lg.Error("%s=>数据读取错误:%s", s.ID, err.Error())
			s.socketError <- Error{
				t:   Read,
				err: fmt.Errorf("%s=>%s", s.ID, err.Error()),
			}
			s.Close()
			return
		}
		s.protocol.Unpack(data[0:i])
	}
}

//发送数据
func (s *SClient) Write(data []byte) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	_, err := s.conn.Write(s.protocol.Packet(data))
	if err != nil {
		return fmt.Errorf("数据发送失败:%s", err.Error())
	}
	return nil
}

//读取消息
func (s *SClient) Read() <-chan []byte {
	return s.protocol.data
}

//关闭连接
func (s *SClient) Close() {
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.conn.Close()
		s.socketError <- Error{
			t:   Cancel,
			err: fmt.Errorf("%s", s.ID),
		}
	}
}
