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
	socketErr   chan error
	ctx         context.Context
	nctx        context.Context
	ncancelFunc context.CancelFunc
	conn        net.Listener
	protocol    *Protocol
	lg          logs.Logger
	port        int
	client      chan SClient
}

//TCP连接到服务器的客户端
type SClient struct {
	conn        net.Conn
	protocol    *Protocol
	ctx         context.Context
	cancelFunc  context.CancelFunc
	lg          logs.Logger
	socketError chan<- error
}

//新建一个服务端,port为监听端口
//注意处理channel error 否则程序阻塞无法正常运行
func NewServer(ctx context.Context, port int) *TcpServer {
	return &TcpServer{
		socketErr: make(chan error, 100),
		ctx:       ctx,
		protocol:  NewProtocol(1000),
		lg:        logs.NewLogger(),
		port:      port,
		client:    make(chan SClient, 100),
	}
}

//启动监听该方法会阻塞
func (s *TcpServer) Listen() error {
	var err error
	s.nctx, s.ncancelFunc = context.WithCancel(s.ctx)
	s.conn, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.lg.Error("查询Socket监听失败:" + err.Error())
		return err
	}
	for {
		select {
		case <-s.ctx.Done():
		case <-s.nctx.Done():
			s.conn.Close()
			s.conn = nil
			return errors.New("关闭")
		default:
			if conn, err := s.conn.Accept(); err != nil {
				s.lg.Error("监听请求连接失败:" + err.Error())
				s.socketErr <- err
			} else {
				s.handleConn(conn)
			}
		}
	}
}

//获取客户端连接
func (s *TcpServer) Accept() <-chan SClient {
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
func (s *TcpServer) Error() <-chan error {
	return s.socketErr
}

//处理新连接
func (s *TcpServer) handleConn(conn net.Conn) {
	sclient := SClient{
		conn:        conn,
		protocol:    NewProtocol(1000),
		socketError: s.socketErr,
		lg:          s.lg,
	}
	sclient.ctx, sclient.cancelFunc = context.WithCancel(s.ctx)
	go sclient.readData()
	s.client <- sclient
}

//读取客户端数据
func (s *SClient) readData() {
	defer recoverPainc(s.readData)
	data := make([]byte, 1024)
	for {
		select {
		case <-s.ctx.Done():
		case <-s.ctx.Done():
			if s.conn != nil {
				s.conn.Close()
				s.conn = nil
			}
			//这里显示结束
			return
		default:
			i, err := s.conn.Read(data)
			if err != nil {
				s.lg.Error("数据读取错误:" + err.Error())
				s.socketError <- err
				s.conn.Close()
				s.conn = nil
				s.cancelFunc()
				return
			}
			s.protocol.Unpack(data[0:i])
		}
	}
}

//发送数据
func (s *SClient) Write(data []byte) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	_, err := s.conn.Write(s.protocol.Packet(data))
	if err != nil {
		ret := fmt.Errorf("数据发送失败:%s", err.Error())
		s.cancelFunc()
		return ret
	}
	return nil
}

//处理连接请求
func (s *SClient) Read() <-chan []byte {
	return s.protocol.Read()
}

//关闭连接
func (s *SClient) Close() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
