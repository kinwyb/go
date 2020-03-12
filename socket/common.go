package socket

import "github.com/kinwyb/go/logs"

// tcp连接协议
type TcpProtocol interface {
	Packing(data []byte) []byte //打包数据
	UnPacking(data []byte)      //解析数据
	ReadMsg() <-chan []byte     //读取解析出来的有效数据
}

//异常处理
func recoverPainc(lg logs.ILogger, f ...func()) {
	if r := recover(); r != nil {
		if lg != nil {
			lg.Fatalf("异常:%s", r)
		}
		if f != nil {
			for _, fu := range f {
				go fu()
			}
		}
	}
}
