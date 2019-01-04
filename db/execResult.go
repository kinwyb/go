package db

import (
	"database/sql"

	"github.com/kinwyb/go/err1"
)

type ExecResult interface {
	sql.Result
	//出错时回调参数方法
	Error(func(err1.Error)) ExecResult
	//是否出错
	HasError(reportZeroChange ...bool) err1.Error
}

//获取一个操作结果对象
func NewExecResult(rs sql.Result) ExecResult {
	return &rus{
		err:    nil,
		Result: rs,
	}
}

//查询错误结果
func ErrExecResult(err err1.Error) ExecResult {
	return &rus{
		err: err,
	}
}

type rus struct {
	sql.Result
	err err1.Error //查询错误
}

func (r *rus) Error(f func(err1.Error)) ExecResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rus) HasError(reportZeroChange ...bool) err1.Error {
	if r.err != nil {
		return r.err
	} else if len(reportZeroChange) < 1 {
		reportZeroChange = []bool{false}
	}
	changrow, _ := r.RowsAffected()
	if changrow == 0 && reportZeroChange[0] {
		return SQLEmptyChange
	}
	return nil
}
