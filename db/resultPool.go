package db

import (
	"database/sql"

	"time"

	"github.com/kinwyb/go/err"
)

//获取一个查询结果对象
func getRes() *res {
	return &res{
		id: time.Now().UnixNano(),
	}
}

//获取一个操作结果对象
func getResult() *rus {
	return &rus{}
}

func DbRows(rows *sql.Rows, fmterr FormatError) Row {
	res := getRes()
	res.data = &QueryResult{res: res}
	res.errFmt = fmterr
	err := res.data.setResult(rows)
	if err != nil {
		res.err = fmterr.FormatError(err)
	} else {
		res.err = nil
	}
	res.norow = len(res.data.data) < 1
	return res
}

//返回一个查询错误
func DbErr(err err.Error) Row {
	ret := getRes()
	ret.err = err
	return ret
}

func DbErrResult(err err.Error) Result {
	ret := getResult()
	ret.err = err
	return ret
}

func DbResult(rs sql.Result) Result {
	ret := getResult()
	ret.err = nil
	ret.Result = rs
	return ret
}
