package socket

//错误类型
type ErrorType int

const (
	Connect ErrorType = iota //连接错误
	Read                     //消息读取错误
	Send                     //消息发送错误
	Listen                   //服务器监听错误
	Cancel                   //关闭
)

type Error struct {
	t   ErrorType //错误类型
	err error     //错误
}

//错误类型
func (e *Error) T() ErrorType {
	return e.t
}

//错误信息
func (e *Error) Error() error {
	return e.err
}
