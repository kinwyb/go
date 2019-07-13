//package gosql mssql工具包..引用"github.com/denisenkom/go-mssqldb"
package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kinwyb/go/err1"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/kinwyb/go/db"
)

//mssql 操作对象
type mssql struct {
	db.Conn
	linkString string
}

//链接mssql数据库
//eg:sqlserver://sa:mypass@localhost?database=master
func Connect(host, username, password, db string) (db.SQL, error) {
	query := url.Values{}
	query.Add("database", db)
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s", host),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	result := &mssql{
		linkString: u.String(),
	}
	sqlDB, err := sql.Open("sqlserver", result.linkString)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(1 * time.Hour) //一个小时后重置链接
	result.SetSQLDB(sqlDB)
	result.SetDataBaseName(db) //记录数据库名称,表名格式化会用到
	result.SetReconnectFunc(result.reconnect)
	return result, nil
}

// 重新连接
func (c *mssql) reconnect() (*sql.DB, error) {
	return sql.Open("sqlserver", c.linkString)
}

//格式化表名称,不做处理直接返回
func (c *mssql) Table(tbname string) string {
	if c == nil || c.Conn.DataBaseName() == "" {
		return tbname
	}
	return "[" + c.Conn.DataBaseName() + "].[dbo].[" + tbname + "]"
}

//RowsCallbackResult 查询多条数据,结果以回调函数处理
//
//@param sql string SQL
//
//@param callback func(*sql.Rows) 回调函数指针
//
//@param args... interface{} SQL参数
func (m *mssql) QueryRows(sql string, args ...interface{}) db.QueryResult {
	i := 0
	sql = regexp.MustCompile("(\\?)").ReplaceAllStringFunc(sql, func(s string) string {
		i++
		return fmt.Sprintf("@p%d", i)
	})
	if len(args) < i {
		return db.ErrQueryResult(
			err1.NewError(-1, "参数缺少,目标参数%d个,实际参数%d个").Format(i, len(args)))
	}
	return m.Conn.QueryRows(sql, args...)
}

//Row 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssql) QueryRow(sql string, args ...interface{}) db.QueryResult {
	if ok, _ := regexp.MatchString("(?i)(.*?) TOP (.*?)\\s?(.*)?", sql); ok {
		sql = regexp.MustCompile("(?i)(.*?) TOP (.*?)\\s?(.*)?").ReplaceAllString(sql, "$1")
	} else {
		sql = strings.Replace(sql, "SELECT ", "SELECT TOP 1 ", 1)
	}
	return m.QueryRows(sql, args...)
}

//ParseSQL 解析SQL
//@param sql string SQL
//@param args map[string]interface{} 参数映射
func (m *mssql) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, err1.Error) {
	cp, err := regexp.Compile("@([^\\s|,|\\)]*)")
	if err != nil {
		return sql, nil, nil
	}
	pts := cp.FindAllStringSubmatch(sql, -1)
	if pts != nil && args != nil { //匹配到数据
		result := make([]interface{}, len(pts))
		for index, s := range pts {
			if v, ok := args[s[1]]; ok { //存在参数
				result[index] = v
			} else {
				return sql, nil, m.FormatError(errors.New("缺少参数[" + s[0] + "]的值"))
			}
		}
		return cp.ReplaceAllString(sql, "?"), result, nil
	}
	return sql, nil, nil
}

//RowsPage 分页查询
func (m *mssql) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.QueryResult {
	if page == nil {
		return m.QueryRows(sql, args...)
	}
	countsql := "select count(0) num from (" + sql + ") as total"
	result := m.QueryRow(countsql, args...)
	count := db.Int64Default(result.Get("num"))
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	if count < 1 {
		return db.NewQueryResult(nil, nil)
	}
	sql = strings.Replace(sql, "SELECT ",
		"SELECT TOP "+strconv.FormatInt(int64(currentpage*page.Rows), 10)+","+strconv.FormatInt(int64(page.Rows), 10),
		1)
	return m.QueryRows(sql, args...)
}
