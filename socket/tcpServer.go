package socket

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/kinwyb/go"

	"errors"

	"github.com/kinwyb/go/logs"
)

// tcp服务端配置
type TcpServerConfig struct {
	Port               int64                                                  //监听端口
	ServerAddress      string                                                 //服务监听地址
	Log                logs.ILogger                                           //日志
	ErrorHandler       func(errType ErrorType, err error, clientID ...string) //错误处理
	NewClientHandler   func(clientID string) TcpProtocol                      //新客户端连接回调,返回该客户端处理协议,可以返回nil
	MessageHandler     func(clientID string, msg []byte)                      //消息处理
	CloseHandler       func()                                                 //服务关闭回调
	ClientCloseHandler func(clientID string)                                  //客户端关闭回调
}

//TCP服务器对象
type TcpServer struct {
	ctx      context.Context
	config   *TcpServerConfig
	conn     net.Listener
	clients  map[string]*SClient
	readData chan []byte //读取到的数据
	doClose  bool        //关闭操作
}

//新建一个服务端
func NewTcpServer(ctx context.Context, config *TcpServerConfig) (*TcpServer, error) {
	if config == nil {
		return nil, errors.New("配置信息不能为空")
	}
	ret := &TcpServer{
		ctx:     ctx,
		config:  config,
		clients: map[string]*SClient{},
	}
	return ret, nil
}

//启动监听该方法会
func (s *TcpServer) Listen() {
	var err error
	s.doClose = false
	s.conn, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.ServerAddress, s.config.Port))
	if err != nil {
		if s.config.Log != nil {
			s.config.Log.Infof("查询Socket监听失败:%s", err.Error())
		}
		if s.config.ErrorHandler != nil {
			s.config.ErrorHandler(ListenErr, err)
		}
		return
	}
	if s.config.Log != nil {
		s.config.Log.Infof("服务器监听开启 => %s:%d", s.config.ServerAddress, s.config.Port)
	}
	go s.handleConn()
	<-s.ctx.Done()
	s.Close()
}

//关闭服务
func (s *TcpServer) Close() {
	s.doClose = true
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	if s.config.CloseHandler != nil {
		s.config.CloseHandler()
	}
	for _, v := range s.clients {
		v.Close()
	}
}

//处理新连接
func (s *TcpServer) handleConn() {
	defer recoverPainc(s.config.Log, s.handleConn)
	for {
		if conn, err := s.conn.Accept(); err != nil {
			if s.doClose { //关闭
				return
			}
			if s.config.Log != nil {
				s.config.Log.Errorf("监听请求连接失败:%s", err.Error())
			}
			if s.config.ErrorHandler != nil {
				s.config.ErrorHandler(ListenErr, err)
			}
			s.Close()
			return
		} else {
			s.newClientAccept(conn)
		}
	}
}

//发送数据
func (s *TcpServer) Write(clientID string, msg []byte) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	client := s.clients[clientID]
	if client == nil {
		return errors.New("客户端链接不存在")
	}
	return client.write(msg)
}

//读取消息
func (s *TcpServer) messageRecv(clientID string, msg []byte) {
	if s.config.MessageHandler != nil {
		s.config.MessageHandler(clientID, msg)
	}
}

// 关闭指定客户端
func (s *TcpServer) CloseClient(clientID string) {
	if c, ok := s.clients[clientID]; ok {
		if c != nil {
			c.Close()
		} else {
			delete(s.clients, clientID)
		}
	}
}

func (s *TcpServer) newClientAccept(conn net.Conn) {
	sclient := &SClient{
		conn:   conn,
		server: s,
		ID:     clientIDGen(),
	}
	var protocol TcpProtocol
	if s.config.NewClientHandler != nil {
		protocol = s.config.NewClientHandler(sclient.ID)
	}
	sclient.protocol = protocol
	sclient.doClose = false
	go sclient.readData()
	if protocol != nil {
		go sclient.recvProtocolMsg()
	}
	s.clients[sclient.ID] = sclient
}

// 客户端id生成器
func clientIDGen() string {
	id := heldiamgo.ID()
	return strings.ToUpper(strconv.FormatUint(id, 32))
}

//TCP连接到服务器的客户端
type SClient struct {
	conn     net.Conn
	server   *TcpServer
	protocol TcpProtocol
	ID       string //客户端唯一标示
	doClose  bool   //关闭
}

//读取客户端数据
func (s *SClient) readData() {
	defer s.Close()
	data := make([]byte, 1024)
	for {
		i, err := s.conn.Read(data)
		if s.doClose {
			return
		} else if err != nil {
			if err == io.EOF { //读取结束
				return
			}
			if s.server.config.Log != nil {
				s.server.config.Log.Errorf("%s=>数据读取错误:%s", s.ID, err.Error())
			}
			if s.server.config.ErrorHandler != nil {
				s.server.config.ErrorHandler(ReadErr, err, s.ID)
			}
			return
		}
		if s.protocol != nil {
			s.protocol.UnPacking(data[0:i])
		} else {
			s.server.messageRecv(s.ID, data[0:i])
		}
	}
}

// 获取通讯协议解析出来的消息
func (s *SClient) recvProtocolMsg() {
	defer recoverPainc(s.server.config.Log, s.recvProtocolMsg)
	if s.protocol != nil {
		for {
			msg := <-s.protocol.ReadMsg()
			s.server.messageRecv(s.ID, msg)
		}
	}
}

//发送数据
func (s *SClient) write(msg []byte) error {
	if s.conn == nil {
		return errors.New("连接已关闭")
	}
	if s.protocol != nil {
		msg = s.protocol.Packing(msg)
	}
	_, err := s.conn.Write(msg)
	return err
}

//关闭连接
func (s *SClient) Close() {
	if !s.doClose && s.server.config.ClientCloseHandler != nil {
		s.server.config.ClientCloseHandler(s.ID)
	}
	s.doClose = true
	delete(s.server.clients, s.ID)
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}
