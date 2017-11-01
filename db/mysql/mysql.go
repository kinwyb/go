//package gosql mysql工具包..引用"github.com/go-sql-driver/mysql"
package mysql

import (
	"database/sql"
	"errors"
	"regexp"

	"strconv"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err"
)

var rep *regexp.Regexp

func init() {
	rep, _ = regexp.Compile("\\s?Error (\\d+):(.*)")
}

//mysql 操作对象
type mysql struct {
	db *sql.DB
}

func (m *mysql) connect() err.Error {
	if m.db == nil {
		return m.FormatError(db.ErrorNotOpen)
	}
	if err := m.db.Ping(); err != nil {
		return m.FormatError(err)
	}
	return nil
}

func (m *mysql) FormatError(e error) err.Error {
	if e == nil {
		return nil
	}
	code := int64(1)
	msg := e.Error()
	if rep.MatchString(e.Error()) {
		d := rep.FindAllStringSubmatch(e.Error(), -1)
		msg = d[0][2]
		cod, err := strconv.ParseInt(d[0][1], 10, 64)
		if err == nil {
			code = cod
		}
	}
	return err.NewError(code, msg, e)
}

//链接mysql数据库，其中other参数代表链接字符串附加的配置信息
//eg:mysql://lcfgly:wang93426@tcp(api.zhifangw.cn:3306)/rfid?loc=Local&multiStatements=true
//其中other="loc=Local&multiStatements=true"
func Connect(host, username, password, db string, other ...string) (db.SQL, error) {
	linkstring := username + ":" + password + "@tcp(" + host + ")/" + db
	if len(other) > 0 {
		linkstring += "?" + other[0]
	}
	result := &mysql{}
	var err error
	result.db, err = sql.Open("mysql", linkstring)
	if err != nil {
		return nil, err
	}
	result.db.SetConnMaxLifetime(1 * time.Hour) //一个小时后重置链接
	return result, nil
}

//Close 关闭数据库连接
func (m *mysql) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

//RowsCallbackResult 查询多条数据,结果以回调函数处理
//
//@param sql string SQL
//
//@param callback func(*sql.Rows) 回调函数指针
//
//@param args... interface{} SQL参数
func (m *mysql) QueryRows(sql string, args ...interface{}) db.Row {
	if err := m.connect(); err != nil {
		return db.DbErr(err)
	}
	rows, err := m.db.Query(sql, args...)
	if err != nil {
		return db.DbErr(m.FormatError(err))
	}
	return db.DbRows(rows, m)
}

//Row 查询单条语句,返回结果
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mysql) QueryRow(sql string, args ...interface{}) db.Row {
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
func (m *mysql) Exec(sql string, args ...interface{}) db.Result {
	if err := m.connect(); err != nil {
		return db.DbErrResult(err)
	}
	result, err := m.db.Exec(sql, args...)
	if err != nil {
		return db.DbErrResult(m.FormatError(err))
	}
	return db.DbResult(result)
}

//Count SQL语句条数统计
//@param sql string SQL
//@param args... interface{} SQL参数
func (m *mysql) Count(sql string, args ...interface{}) (int64, err.Error) {
	if ok, _ := regexp.MatchString("(?i)(.*?) LIMIT (.*?)\\s?(.*)?", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	if ok, _ := regexp.MatchString("(?i).* group by .*", sql); ok {
		sql = "SELECT COUNT(1) FROM (" + sql + ") as tmp"
	}
	sql = regexp.MustCompile("^(?i)select .*? from (.*) order by (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	sql = regexp.MustCompile("^(?i)select .*? from (.*)").ReplaceAllString(sql, "SELECT count(1) FROM $1")
	if err := m.connect(); err != nil {
		return 0, err
	}
	result := m.db.QueryRow(sql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return 0, m.FormatError(err)
	}
	return count, nil
}

//ParseSQL 解析SQL
//@param sql string SQL
//@param args map[string]interface{} 参数映射
func (m *mysql) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, err.Error) {
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

//Transaction 事务处理
//@param t TransactionFunc 事务处理函数
func (m *mysql) Transaction(t db.TransactionFunc) err.Error {
	if err := m.connect(); err != nil {
		return err
	}
	tx, err := m.db.Begin()
	if err == nil {
		if t != nil {
			e := t(&mysqlTx{tx: tx, fmterr: m})
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
	return m.FormatError(err)
}

//GetDb 获取数据库对象
func (m *mysql) GetDb() (*sql.DB, err.Error) {
	if err := m.connect(); err != nil {
		return nil, err
	}
	return m.db, nil
}

//RowsPage 分页查询
func (m *mysql) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.Row {
	if page == nil {
		return m.QueryRows(sql, args...)
	}
	countsql := "select count(0) from (" + sql + ") as total"
	if err := m.connect(); err != nil {
		return db.DbErr(err)
	}
	result := m.db.QueryRow(countsql, args...)
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return db.DbErr(m.FormatError(err))
	}
	page.SetTotal(count)
	currentpage := 0
	if page.Page-1 > 0 {
		currentpage = page.Page - 1
	}
	sql = sql + " LIMIT " + strconv.FormatInt(int64(currentpage*page.Rows), 10) + "," + strconv.FormatInt(int64(page.Rows), 10)
	return m.QueryRows(sql, args...)
}
