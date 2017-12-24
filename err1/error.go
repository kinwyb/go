package err1

import (
	"encoding/json"
	"fmt"
)

//Error 错误接口
type Error interface {
	json.Marshaler
	Code() int64    //自定义错误编码
	Msg() string    //自定义错误消息
	Err() error     //具体的错误
	Caller() string //返回调用堆栈信息
	Error() string  //继承全局的error接口
}

//err 公用错误对象
type err struct {
	code int64
	msg  string
	e    error
}

//Error 错误
func (e *err) Error() string {
	if e.msg != "" {
		return e.msg
	} else if e.e != nil {
		return e.e.Error()
	}
	return " none"
}

func (e *err) Caller() string {
	return CallInfo(3)
}

//Code 错误代码
func (e *err) Code() int64 {
	return e.code
}

//Msg 错误消息
func (e *err) Msg() string {
	return e.msg
}

//Err 原始错误
func (e *err) Err() error {
	if e.e == nil {
		return nil
	}
	return e.e
}

func (e *err) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte{}, nil
	}
	return []byte(fmt.Sprintf("{\"code\":%d,\"msg\":\"%s\",\"errmsg\":\"%s\"}", e.code, e.msg, e.Error())), nil
}

//NewError 新建错误
func NewError(code int64, msg string, errs ...error) Error {
	ret := &err{code: code, msg: msg}
	if errs != nil && len(errs) > 0 {
		ret.e = errs[0]
	}
	return ret
}
