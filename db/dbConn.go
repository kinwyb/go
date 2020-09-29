package db

import (
	"context"
	"database/sql"
	"errors"
	"regexp"

	"strconv"
)

//Conn 操作对象
type Conn struct {
	db     *sql.DB
	dbname string
}

func (c *Conn) connect() error {
	if c.db == nil {
		return ErrorNotOpen
	}
	if err := c.db.Ping(); err != nil {
		return ErrorPing
	}
	return nil
}

//设置数据库链接
func (c *Conn) SetSQLDB(dbSQL *sql.DB) {
	c.db = dbSQL
}

//设置数据库名称
func (c *Conn) SetDataBaseName(dbname string) {
	c.dbname = dbname
}

//Close 关闭数据库连接
func (c *Conn) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

//RowsCallbackResult 查询多条数据,结果以回调函数处理
//
//@param sql string SQL
//
//@param callback func(*sql.Rows) 回调函数指针
//
//@param args... interface{} SQL参数
func (c *Conn) QueryRows(sql string, args ...interface{}) QueryResult {
	if err := c.connect(); err != nil {
		return ErrQueryResult(err, sql, args)
	}
	rows, err := c.db.Query(sql, args...)
	if err != nil {
		return ErrQueryResult(err, sql, args)
	}
	return NewQueryResult(rows, sql, args)
}

//Row 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (c *Conn) QueryRow(sql string, args ...interface{}) QueryResult {
	if ok, _ := regexp.MatchString("(?i)(.*?) LIMIT (.*?)\\s?(.*)?", sql); ok {
		sql = regexp.MustCompile("(?i)(.*?) LIMIT (.*?)\\s?(.*)?").ReplaceAllString(sql, "$1")
	} else {
		sql += " LIMIT 1 "
	}
	return c.QueryRows(sql, args...)
}

//Exec 执行一条SQL
//@param sql string SQL
//@param args... interface{} SQL参数
func (c *Conn) Exec(sql string, args ...interface{}) ExecResult {
	if err := c.connect(); err != nil {
		return ErrExecResult(err, sql, args)
	}
	result, err := c.db.Exec(sql, args...)
	if err != nil {
		return ErrExecResult(err, sql, args)
	}
	return NewExecResult(result)
}

//Count SQL语句条数统计
//@param sql string SQL
//@param args... interface{} SQL参数
func (c *Conn) Count(sql string, args ...interface{}) (int64, error) {
	if ok, _ := regexp.MatchString("(?i)(.*?) LIMIT (.*?)\\s?(.*)?", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	if ok, _ := regexp.MatchString("(?i).* group by .*", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	sql = regexp.MustCompile("^(?i)select .*? from (.*) order by (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	sql = regexp.MustCompile("^(?i)select .*? from (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	if err := c.connect(); err != nil {
		return 0, err
	}
	result := c.db.QueryRow(sql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

//ParseSQL 解析SQL
//@param sql string SQL
//@param args map[string]interface{} 参数映射
func (c *Conn) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, error) {
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

//Transaction 事务处理
//@param t TransactionFunc 事务处理函数
func (c *Conn) Transaction(t TransactionFunc, option ...*TxOption) error {
	if err := c.connect(); err != nil {
		return err
	}
	var txOption *sql.TxOptions
	if len(option) > 0 && option[0] != nil {
		txOption = option[0].Option
	}
	tx, err := c.db.BeginTx(context.Background(), txOption)
	if err == nil {
		defer func() {
			if err := recover(); err != nil {
				//发生异常,先回滚事务再继续抛出异常
				tx.Rollback() //回滚
				panic(err)
			}
		}()
		if t != nil {
			e := t(&ConnTx{tx: tx, db: c})
			if e != nil {
				tx.Rollback()
				return e
			}
			err = tx.Commit()
			if err != nil { //事务提交失败,回滚事务,返回错误
				tx.Rollback()
			}
		}
	}
	return err
}

//GetDb 获取数据库对象
func (c *Conn) GetDb() (*sql.DB, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c.db, nil
}

//RowsPage 分页查询
func (c *Conn) QueryWithPage(sql string, page *PageObj, args ...interface{}) QueryResult {
	if page == nil {
		return c.QueryRows(sql, args...)
	}
	countsql := "select count(0) from (" + sql + ") as total"
	if err := c.connect(); err != nil {
		return ErrQueryResult(err, sql, args)
	}
	result := c.db.QueryRow(countsql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return ErrQueryResult(err, sql, args)
	}
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	if count < 1 {
		return NewQueryResult(nil, sql, args)
	}
	sql = sql + " LIMIT " + strconv.FormatInt(int64(currentpage*page.Rows), 10) + "," + strconv.FormatInt(int64(page.Rows), 10)
	return c.QueryRows(sql, args...)
}

func (c *Conn) Prepare(query string) (*sql.Stmt, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	stmt, e := c.db.Prepare(query)
	return stmt, e
}

//格式化表名称,不做处理直接返回
func (c *Conn) Table(tbname string) string {
	if c == nil || c.dbname == "" {
		return tbname
	}
	return "`" + c.dbname + "`." + tbname
}

//数据库名称
func (c *Conn) DataBaseName() string {
	if c == nil || c.dbname == "" {
		return ""
	}
	return c.dbname
}

// 设置最大连接数
func (c *Conn) SetMaxOpenConns(n int) {
	c.db.SetMaxOpenConns(n)
}
