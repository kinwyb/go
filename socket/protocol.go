package socket

import (
	"bytes"
	"encoding/binary"
)

const bitlength = 10 //数据长度占用字节数

//Protocol 自定义协议解析
type Protocol struct {
	data       chan []byte   //解析成功的数据
	byteBuffer *bytes.Buffer //数据存储中心
	dataLength int64         //当前消息数据长度
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
func (p *Protocol) Packing(message []byte) []byte {
	return append(intToByte(int64(len(message))), message...)
}

//Read 获取数据读取的channel对象
func (p *Protocol) ReadMsg() <-chan []byte {
	return p.data
}

//解析成功的数据请用Read方法获取
func (p *Protocol) UnPacking(buffer []byte) {
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
