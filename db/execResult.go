package db

import (
	"database/sql"

	"github.com/gogo/protobuf/proto"

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

type rusMsg struct {
	lastInsertId int64
	rowsAffected int64
	err          err1.Error //查询错误
}

func (r *rusMsg) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

func (r *rusMsg) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func (r *rusMsg) Error(f func(err1.Error)) ExecResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rusMsg) HasError(reportZeroChange ...bool) err1.Error {
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

func ExecResultToBytes(v ExecResult) []byte {
	msg := &ExecResultMsg{}
	l, err := v.LastInsertId()
	if err != nil {
		msg.ErrMsg = err.Error()
		msg.ErrCode = -1
	}
	msg.LastInsertId = l
	r, err := v.RowsAffected()
	if err != nil {
		msg.ErrMsg = err.Error()
		msg.ErrCode = -1
	}
	msg.RowsAffected = r
	e := v.HasError(true)
	if e != nil {
		msg.ErrCode = e.Code()
		msg.ErrMsg = e.Msg()
	}
	ret, _ := proto.Marshal(msg)
	return ret
}

func BytesToExecResult(v []byte) ExecResult {
	msg := &ExecResultMsg{}
	proto.Unmarshal(v, msg)
	ret := &rusMsg{
		lastInsertId: msg.LastInsertId,
		rowsAffected: msg.RowsAffected,
	}
	if msg.ErrCode != 0 {
		ret.err = err1.NewError(msg.ErrCode, msg.ErrMsg)
	}
	return ret
}
