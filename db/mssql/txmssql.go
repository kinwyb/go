package mssql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err1"
)

//mssqlTx 事务操作
type mssqlTx struct {
	tx     *sql.Tx
	mssql  *mssql
	fmterr db.FormatError
}

//Rows 查询多条数据,结果以[]map[string]interface{}方式返回
//
//返回结果,使用本package中的类型函数进行数据解析
//eg:
//		result := Rows(...)
//		for _,mp := range result {
//			Int(mp["colum"])
//			String(mp["colum"])
//			.......
//		}
//@param sql string SQL
//
//@param args... interface{} SQL参数
func (m *mssqlTx) QueryRows(sql string, args ...interface{}) db.QueryResult {
	i := 0
	sql = regexp.MustCompile("(\\?)").ReplaceAllStringFunc(sql, func(s string) string {
		i++
		return fmt.Sprintf("@p%d", i)
	})
	if len(args) < i {
		return db.ErrQueryResult(
			err1.NewError(-1, "参数缺少,目标参数%d个,实际参数%d个").Format(i, len(args)))
	}
	rows, err := m.tx.Query(sql, args...)
	if err != nil {
		return db.ErrQueryResult(m.fmterr.FormatError(err))
	}
	return db.NewQueryResult(rows, m.fmterr)
}

func (m *mssqlTx) Prepare(query string) (*sql.Stmt, err1.Error) {
	stmt, err := m.tx.Prepare(query)
	return stmt, FormatError(err)
}

//QueryResult 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssqlTx) QueryRow(sql string, args ...interface{}) db.QueryResult {
	if ok, _ := regexp.MatchString("(?i)(.*?) TOP (.*?)\\s?(.*)?", sql); ok {
		sql = regexp.MustCompile("(?i)(.*?) TOP (.*?)\\s?(.*)?").ReplaceAllString(sql, "$1")
	} else {
		sql = strings.Replace(sql, "SELECT ", "SELECT TOP 1 ", 1)
	}
	return m.QueryRows(sql, args...)
}

//Exec 执行一条SQL
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssqlTx) Exec(sql string, args ...interface{}) db.ExecResult {
	result, err := m.tx.Exec(sql, args...)
	if err != nil {
		return db.ErrExecResult(m.fmterr.FormatError(err))
	}
	return db.NewExecResult(result)
}

//Count SQL语句条数统计
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssqlTx) Count(sql string, args ...interface{}) (int64, err1.Error) {
	if ok, _ := regexp.MatchString("(?i)(.*?) LIMIT (.*?)\\s?(.*)?", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	if ok, _ := regexp.MatchString("(?i).* group by .*", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	sql = regexp.MustCompile("^(?i)select .*? from (.*) order by (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	sql = regexp.MustCompile("^(?i)select .*? from (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	result := m.tx.QueryRow(sql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return 0, m.fmterr.FormatError(err)
	}
	return count, nil
}

//GetTx 获取事务对象
func (m *mssqlTx) GetTx() *sql.Tx {
	return m.tx
}

//RowsPage 分页查询
func (m *mssqlTx) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.QueryResult {
	if page == nil {
		return m.QueryRows(sql, args...)
	}
	countsql := "select count(0) from (" + sql + ") as total"
	result := m.tx.QueryRow(countsql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return db.ErrQueryResult(m.fmterr.FormatError(err))
	}
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	if count < 1 {
		return db.NewQueryResult(nil, nil)
	}
	sql = sql + " LIMIT " + strconv.FormatInt(int64(currentpage*page.Rows), 10) + "," + strconv.FormatInt(int64(page.Rows), 10)
	return m.QueryRows(sql, args...)
}

//格式化表名称,不做处理直接返回
func (m *mssqlTx) Table(tbname string) string {
	if m == nil || m.mssql == nil || m.mssql.dbname == "" {
		return tbname
	}
	return "`" + m.mssql.dbname + "`." + tbname
}

//Transaction 事务处理
//param t TransactionFunc 事务处理函数
func (m *mssqlTx) Transaction(t db.TransactionFunc, new ...bool) err1.Error {
	if t != nil {
		if len(new) > 0 && new[0] && m.mssql != nil {
			//要求新事物返回新事务
			return m.mssql.Transaction(t)
		}
		//本身就是事务了，直接调用即可
		return t(m)
	}
	return nil
}

//数据库名称
func (m *mssqlTx) DataBaseName() string {
	if m.mssql != nil {
		return m.mssql.dbname
	}
	return ""
}
