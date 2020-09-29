package db

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/proto"
)

type ExecResult interface {
	sql.Result
	// 出错时回调参数方法
	Error(func(error)) ExecResult
	// 错误保存到日志
	ErrorToLog(log *logrus.Entry, msg string) ExecResult
	//是否出错
	HasError(reportZeroChange ...bool) error
}

//获取一个操作结果对象
func NewExecResult(rs sql.Result) ExecResult {
	return &rus{
		err:    nil,
		Result: rs,
	}
}

//查询错误结果
func ErrExecResult(err error, sql string, args []interface{}) ExecResult {
	return &rus{
		sql:  sql,
		args: args,
		err:  err,
	}
}

type rus struct {
	sql.Result
	sql  string
	args []interface{}
	err  error //查询错误
}

func (r *rus) Error(f func(err error)) ExecResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rus) ErrorToLog(log *logrus.Entry, msg string) ExecResult {
	if r.err != nil && log != nil {
		log.WithField("sql", r.sql).
			WithField("req", r.args).
			WithError(r.err).Errorf("SQL错误:%s", msg)
	}
	return r
}

func (r *rus) HasError(reportZeroChange ...bool) error {
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
	err          error //查询错误
}

func (r *rusMsg) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

func (r *rusMsg) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func (r *rusMsg) Error(f func(err error)) ExecResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rusMsg) ErrorToLog(log *logrus.Entry, msg string) ExecResult {
	if r.err != nil && log != nil {
		log.WithError(r.err).Errorf("SQL错误:%s", msg)
	}
	return r
}

func (r *rusMsg) HasError(reportZeroChange ...bool) error {
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
		msg.ErrCode = -1
		msg.ErrMsg = e.Error()
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
	if msg.ErrMsg != "" {
		ret.err = errors.New(msg.ErrMsg)
	}
	return ret
}
