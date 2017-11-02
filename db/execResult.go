package db

import (
	"database/sql"

	"github.com/kinwyb/go/err"
)

type ExecResult interface {
	sql.Result
	//出错时回调参数方法
	Error(func(err.Error)) ExecResult
	//是否出错
	HasError() err.Error
}

//获取一个操作结果对象
func NewExecResult(rs sql.Result) ExecResult {
	return &rus{
		err:    nil,
		Result: rs,
	}
}

//查询错误结果
func ErrExecResult(err err.Error) ExecResult {
	return &rus{
		err: err,
	}
}

type rus struct {
	sql.Result
	err err.Error //查询错误
}

func (r *rus) Error(f func(err.Error)) ExecResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rus) HasError() err.Error {
	return r.err
}
