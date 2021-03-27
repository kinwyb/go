package socket

//错误类型
type ErrorType int

const (
	ConnectErr ErrorType = iota + 1 //连接错误
	ReadErr                         //消息读取错误
	SendErr                         //消息发送错误
	ListenErr                       //服务器监听错误
)

func (e ErrorType) String() string {
	switch e {
	case ConnectErr:
		return "连接错误"
	case ReadErr:
		return "读取错误"
	case SendErr:
		return "发送错误"
	case ListenErr:
		return "监听错误"
	default:
		return "未知错误"
	}
}

type Error struct {
	t   ErrorType //错误类型
	err error     //错误
}

// 新增错误
func NewError(t ErrorType, err error) *Error {
	return &Error{
		t:   t,
		err: err,
	}
}

//错误类型
func (e *Error) T() ErrorType {
	return e.t
}

//错误信息
func (e *Error) Error() error {
	return e.err
}
