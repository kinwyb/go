package db

import (
	"database/sql"

	"github.com/kinwyb/go/err"
)

//查询结果返回接口
type Row interface {
	//获取结果
	GetRows() *QueryResult
	//逐条获取结果
	//如果参数func返回true，并且还有下一条结果则再次调用func返回下一条
	Rows(func(map[string]interface{}) bool) Row
	//对象ID编码
	Id() int64
	//出错时回调参数方法
	Error(func(err.Error)) Row
	//是否出错
	HasError() err.Error
	//是否为空
	IsEmpty() bool
	//结果空是回调参数方法
	Empty(func()) Row
	//关闭查询结果.
	//如果读取了结果内容查询会自动关闭,只有不需要获取查询结果的时候才需要手动调用关闭查询结果
	Close()
}

type res struct {
	id     int64
	data   *QueryResult //查询结果
	err    err.Error  //查询错误
	norow  bool         //空数据
	errFmt FormatError  //错误格式化
}

func (r *res) GetRows() *QueryResult {
	return r.data
}

func (r *res) Rows(f func(map[string]interface{}) bool) Row {
	if f != nil && r.data != nil {
		r.data.ForEach(f)
	}
	return r
}

func (r *res) Error(f func(err.Error)) Row {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *res) Close() {
	if r.data != nil && r.data.rows != nil {
		r.data.rows.Close()
	}
}

func (r *res) Id() int64 {
	return r.id
}

func (r *res) HasError() err.Error {
	return r.err
}

func (r *res) IsEmpty() bool {
	return r.norow
}

func (r *res) Empty(f func()) Row {
	if r.norow && f != nil {
		f()
	}
	return r
}

type Result interface {
	sql.Result
	//出错时回调参数方法
	Error(func(err.Error)) Result
	//是否出错
	HasError() err.Error
}

type rus struct {
	sql.Result
	err err.Error //查询错误
}

func (r *rus) Error(f func(err.Error)) Result {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *rus) HasError() err.Error {
	return r.err
}
