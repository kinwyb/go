package db

import (
	"database/sql"
	"regexp"
	"strconv"
)

//ConnTx 事务操作
type ConnTx struct {
	tx *sql.Tx
	db *Conn
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
func (m *ConnTx) QueryRows(sql string, args ...interface{}) QueryResult {
	rows, err := m.tx.Query(sql, args...)
	if err != nil {
		return ErrQueryResult(err, sql, args)
	}
	return NewQueryResult(rows, sql, args)
}

func (m *ConnTx) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := m.tx.Prepare(query)
	return stmt, err
}

//QueryResult 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *ConnTx) QueryRow(sql string, args ...interface{}) QueryResult {
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
func (m *ConnTx) Exec(sql string, args ...interface{}) ExecResult {
	result, err := m.tx.Exec(sql, args...)
	if err != nil {
		return ErrExecResult(err, sql, args)
	}
	return NewExecResult(result)
}

//Count SQL语句条数统计
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *ConnTx) Count(sql string, args ...interface{}) (int64, error) {
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
		return 0, err
	}
	return count, nil
}

//GetTx 获取事务对象
func (m *ConnTx) GetTx() *sql.Tx {
	return m.tx
}

//RowsPage 分页查询
func (m *ConnTx) QueryWithPage(sql string, page *PageObj, args ...interface{}) QueryResult {
	if page == nil {
		return m.QueryRows(sql, args...)
	}
	countsql := "select count(0) from (" + sql + ") as total"
	result := m.tx.QueryRow(countsql, args...)
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
	return m.QueryRows(sql, args...)
}

//格式化表名称,不做处理直接返回
func (m *ConnTx) Table(tbname string) string {
	if m == nil || m.db == nil || m.db.dbname == "" {
		return tbname
	}
	return "`" + m.db.dbname + "`." + tbname
}

//Transaction 事务处理
//param t TransactionFunc 事务处理函数
func (m *ConnTx) Transaction(t TransactionFunc, option ...*TxOption) error {
	if t != nil {
		if len(option) > 0 && option[0] != nil {
			if option[0].New {
				option[0].New = false
				return m.db.Transaction(t, option[0])
			}
		}
		//本身就是事务了，直接调用即可
		return t(m)
	}
	return nil
}

//数据库名称
func (m *ConnTx) DataBaseName() string {
	if m.db != nil {
		return m.db.dbname
	}
	return ""
}

// 设置最大连接数
func (m *ConnTx) SetMaxOpenConns(n int) {
	// 无效设置
}
