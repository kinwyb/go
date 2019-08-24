package socket_deprecated

import (
	"bytes"
	"encoding/binary"
)

const bitlength = 10 //数据长度占用字节数

//心跳包内容
var heartBeatBytes = []byte{0x11, 0x22, 0x13, 0x24, 0x15, 0x26, 0x17, 0x28, 0x19, 0x11}

//Protocol 自定义协议
//解决TCP粘包问题
type Protocol struct {
	data       chan []byte    //解析成功的数据
	byteBuffer *bytes.Buffer  //数据存储中心
	dataLength int64          //当前消息数据长度
	heartbeat  []byte         //心跳包数据，如果设置而且接收到改数据会被忽略而不被输出
	Handler    PackageHandler //数据编码handler
}

//NewProtocol 初始化一个Protocol
// chanLength为解析成功数据channel缓冲长度
func NewProtocol(chanLength ...int) *Protocol {
	length := 100
	if chanLength != nil && len(chanLength) > 0 {
		length = chanLength[0]
	}
	return &Protocol{
		data:       make(chan []byte, length),
		byteBuffer: bytes.NewBufferString(""),
	}
}

//Packet 封包
func (p *Protocol) Packet(message []byte) []byte {
	if p.Handler != nil {
		message = p.Handler.Package(message)
	}
	return append(intToByte(int64(len(message))), message...)
}

//Read 获取数据读取的channel对象
func (p *Protocol) Read() <-chan []byte {
	return p.data
}

//设置心跳包数据内容，如果接收到的一条消息刚好于设置的心跳包
//内容一致,这条消息将会忽略不会进入读取成功的消息队列中
func (p *Protocol) SetHeartBeat(b []byte) {
	p.heartbeat = b
}

//解析成功的数据请用Read方法获取
func (p *Protocol) Unpack(buffer []byte) {
	p.byteBuffer.Write(buffer)
	for { //多条数据循环处理
		length := p.byteBuffer.Len()
		if length < bitlength { //前面8个字节是长度
			return
		}
		p.dataLength = byteToInt(p.byteBuffer.Bytes()[0:bitlength])
		if int64(length) < p.dataLength+bitlength { //数据长度不够,等待下次读取数据
			return
		}
		data := make([]byte, p.dataLength+bitlength)
		p.byteBuffer.Read(data)
		msg := data[bitlength:]
		if p.Handler != nil { //解包
			msg = p.Handler.UnPackage(msg)
		}
		if p.heartbeat != nil && bytes.Equal(msg, p.heartbeat) {
			//对比接收到的内容如果和设置的内容一致忽略该条消息
			continue
		}
		p.data <- msg
	}
}

//重置
func (p *Protocol) Reset() {
	p.dataLength = 0
	p.byteBuffer.Reset() //清空重新开始
}

func intToByte(num int64) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, uint64(num))
	return ret
}

func byteToInt(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}
