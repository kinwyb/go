//package gosql mssql工具包..引用"github.com/denisenkom/go-mssqldb"
package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kinwyb/go/conv"
	"github.com/xwb1989/sqlparser"
	"net/url"
	"regexp"
	"strings"
	"time"

	sqlserver "github.com/denisenkom/go-mssqldb"
	"github.com/kinwyb/go/db"
)

//mssql 操作对象
type mssql struct {
	db.Conn
	linkString string
}

//链接mssql数据库
//eg:sqlserver://sa:mypass@localhost?database=master
func Connect(host, username, password, db string, params ...url.Values) (db.SQL, error) {
	query := url.Values{}
	if len(params) > 0 {
		query = params[0]
	}
	query.Set("database", db)
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
	sqlDB.SetMaxOpenConns(30)               //最大连接数
	sqlDB.SetConnMaxLifetime(1 * time.Hour) //一个小时后重置链接
	result.SetSQLDB(sqlDB)
	result.SetDataBaseName(db) //记录数据库名称,表名格式化会用到
	return result, nil
}

//格式化表名称,不做处理直接返回
func (m *mssql) Table(tbname string) string {
	return tbname
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
		return db.ErrQueryResult(fmt.Errorf("参数缺少,目标参数%d个,实际参数%d个", i, len(args)), sql, args)
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
func (m *mssql) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, error) {
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
				return sql, nil, errors.New("缺少参数[" + s[0] + "]的值")
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
	sql = strings.ReplaceAll(sql, "?", "@")
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return db.ErrQueryResult(fmt.Errorf("sql语句解析错误:%w", err), sql, args)
	}
	selectColumn := ""
	from := ""
	where := ""
	orderBy := ""
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.SelectExprs.Format(buf)
		selectColumn = buf.String()
		buf.Reset()
		stmt.From.Format(buf)
		from = buf.String()
		buf.Reset()
		stmt.Where.Format(buf)
		where = buf.String()
		buf.Reset()
		stmt.OrderBy.Format(buf)
		orderBy = buf.String()
	default:
		return db.ErrQueryResult(errors.New("只支持select语句"), sql, args)
	}
	where = strings.ReplaceAll(where, "@", "?")
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("SELECT count(0) num FROM ")
	sqlBuilder.WriteString(from)
	sqlBuilder.WriteString(where)
	result := m.QueryRows(sqlBuilder.String(), args...)
	count := conv.ToInt64(result.Get("num"))
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	if count < 1 {
		return db.NewQueryResult(nil, sql, args)
	}
	sqlBuilder.Reset()
	sqlBuilder.WriteString("SELECT TOP ")
	sqlBuilder.WriteString(conv.ToString(page.Rows))
	sqlBuilder.WriteString(" * FROM (SELECT ROW_NUMBER() OVER (")
	sqlBuilder.WriteString(orderBy)
	sqlBuilder.WriteString(") as RowNumber,")
	sqlBuilder.WriteString(selectColumn)
	sqlBuilder.WriteString(" FROM ")
	sqlBuilder.WriteString(from)
	sqlBuilder.WriteString(where)
	sqlBuilder.WriteString(") as tmp WHERE RowNumber > ")
	sqlBuilder.WriteString(conv.ToString(page.Rows * currentpage))
	sqlBuilder.WriteString(" ORDER BY RowNumber ASC ")
	sql = sqlBuilder.String()
	return m.QueryRows(sql, args...)
}

//Exec 执行一条SQL
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssql) Exec(sql string, args ...interface{}) db.ExecResult {
	i := 0
	sql = regexp.MustCompile("(\\?)").ReplaceAllStringFunc(sql, func(s string) string {
		i++
		return fmt.Sprintf("@p%d", i)
	})
	if len(args) < i {
		return db.ErrExecResult(fmt.Errorf("参数缺少,目标参数%d个,实际参数%d个", i, len(args)), sql, args)
	}
	return m.Conn.Exec(sql, args...)
}

//Transaction 事务处理
//@param t TransactionFunc 事务处理函数
func (m *mssql) Transaction(t db.TransactionFunc, option ...*db.TxOption) error {
	f := func(tx db.TxSQL) error {
		return t(&mssqlTx{
			TxSQL: tx,
			db:    m,
		})
	}
	return m.Conn.Transaction(f, option...)
}

type mssqlTx struct {
	db.TxSQL
	db *mssql
}

//Transaction 事务处理
//@param t TransactionFunc 事务处理函数
func (m *mssqlTx) Transaction(t db.TransactionFunc, options ...*db.TxOption) error {
	if t != nil {
		if len(options) > 0 && options[0] != nil && options[0].New {
			options[0].New = false
			//要求新事物返回新事务
			return m.db.Transaction(t, options...)
		}
		//本身就是事务了，直接调用即可
		return t(m)
	}
	return nil
}

//RowsCallbackResult 查询多条数据,结果以回调函数处理
//
//@param sql string SQL
//
//@param callback func(*sql.Rows) 回调函数指针
//
//@param args... interface{} SQL参数
func (m *mssqlTx) QueryRows(sql string, args ...interface{}) db.QueryResult {
	i := 0
	sql = regexp.MustCompile("(\\?)").ReplaceAllStringFunc(sql, func(s string) string {
		i++
		return fmt.Sprintf("@p%d", i)
	})
	if len(args) < i {
		return db.ErrQueryResult(fmt.Errorf("参数缺少,目标参数%d个,实际参数%d个", i, len(args)), sql, args)
	}
	return m.TxSQL.QueryRows(sql, args...)
}

//Row 查询单条语句,返回结果
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

//ParseSQL 解析SQL
//@param sql string SQL
//@param args map[string]interface{} 参数映射
func (m *mssqlTx) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, error) {
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
				return sql, nil, errors.New("缺少参数[" + s[0] + "]的值")
			}
		}
		return cp.ReplaceAllString(sql, "?"), result, nil
	}
	return sql, nil, nil
}

//格式化表名称,不做处理直接返回
func (m *mssqlTx) Table(tbname string) string {
	return tbname
}

//RowsPage 分页查询
func (m *mssqlTx) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.QueryResult {
	if page == nil {
		return m.QueryRows(sql, args...)
	}
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return db.ErrQueryResult(fmt.Errorf("sql语句解析错误:%w", err), sql, args)
	}
	selectColumn := ""
	from := ""
	where := ""
	orderBy := ""
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.SelectExprs.Format(buf)
		selectColumn = buf.String()
		buf.Reset()
		stmt.From.Format(buf)
		from = buf.String()
		buf.Reset()
		stmt.Where.Format(buf)
		where = buf.String()
		buf.Reset()
		stmt.OrderBy.Format(buf)
		orderBy = buf.String()
	default:
		return db.ErrQueryResult(errors.New("只支持select语句"), sql, args)
	}
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("SELECT count(0) num FROM ")
	sqlBuilder.WriteString(from)
	sqlBuilder.WriteString(where)
	result := m.QueryRows(sqlBuilder.String(), args...)
	count := conv.ToInt64(result.Get("num"))
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	if count < 1 {
		return db.NewQueryResult(nil, sql, args)
	}
	sqlBuilder.Reset()
	sqlBuilder.WriteString("SELECT TOP ")
	sqlBuilder.WriteString(conv.ToString(page.Rows))
	sqlBuilder.WriteString(" * FROM (SELECT ROW_NUMBER() OVER (")
	sqlBuilder.WriteString(orderBy)
	sqlBuilder.WriteString(") as RowNumber,")
	sqlBuilder.WriteString(selectColumn)
	sqlBuilder.WriteString(" FROM ")
	sqlBuilder.WriteString(from)
	sqlBuilder.WriteString(where)
	sqlBuilder.WriteString(") as tmp WHERE RowNumber > ")
	sqlBuilder.WriteString(conv.ToString(page.Rows * currentpage))
	sqlBuilder.WriteString(" ORDER BY RowNumber ASC ")
	sql = sqlBuilder.String()
	return m.QueryRows(sql, args...)
}

//Exec 执行一条SQL
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mssqlTx) Exec(sql string, args ...interface{}) db.ExecResult {
	i := 0
	sql = regexp.MustCompile("(\\?)").ReplaceAllStringFunc(sql, func(s string) string {
		i++
		return fmt.Sprintf("@p%d", i)
	})
	if len(args) < i {
		return db.ErrExecResult(fmt.Errorf("参数缺少,目标参数%d个,实际参数%d个", i, len(args)), sql, args)
	}
	return m.TxSQL.Exec(sql, args...)
}

// sqlserver解码
func UniqueIdentifierToString(v interface{}) string {
	if v == nil {
		return ""
	}
	i := sqlserver.UniqueIdentifier{}
	_ = i.Scan(v)
	return i.String()
}
