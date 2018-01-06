package mysql

import (
	"database/sql"
	"regexp"
	"strconv"

	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err1"
)

//MySQLTx 事务操作
type mysqlTx struct {
	tx     *sql.Tx
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
func (m *mysqlTx) QueryRows(sql string, args ...interface{}) db.QueryResult {
	rows, err := m.tx.Query(sql, args...)
	if err != nil {
		return db.ErrQueryResult(m.fmterr.FormatError(err))
	}
	return db.NewQueryResult(rows, m.fmterr)
}

func (m *mysqlTx) Prepare(query string) (*sql.Stmt, err1.Error) {
	stmt, err := m.tx.Prepare(query)
	return stmt, FormatError(err)
}

//QueryResult 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mysqlTx) QueryRow(sql string, args ...interface{}) db.QueryResult {
	if ok, _ := regexp.MatchString("(?i)(.*?) LIMIT (.*?)\\s?(.*)?", sql); ok {
		sql = regexp.MustCompile("(?i)(.*?) LIMIT (.*?)\\s?(.*)?").ReplaceAllString(sql, "$1")
	} else {
		sql += " LIMIT 1 "
	}
	return m.QueryRows(sql, args...)
}

//Exec 执行一条SQL
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mysqlTx) Exec(sql string, args ...interface{}) db.ExecResult {
	result, err := m.tx.Exec(sql, args...)
	if err != nil {
		return db.ErrExecResult(m.fmterr.FormatError(err))
	}
	return db.NewExecResult(result)
}

//Count SQL语句条数统计
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mysqlTx) Count(sql string, args ...interface{}) (int64, err1.Error) {
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
func (m *mysqlTx) GetTx() *sql.Tx {
	return m.tx
}

//RowsPage 分页查询
func (m *mysqlTx) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.QueryResult {
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
	sql = sql + " LIMIT " + strconv.FormatInt(int64(currentpage*page.Rows), 10) + "," + strconv.FormatInt(int64(page.Rows), 10)
	return m.QueryRows(sql, args...)
}

//格式化表名称,不做处理直接返回
func (m *mysqlTx) Table(tbname string) string {
	return tbname
}
