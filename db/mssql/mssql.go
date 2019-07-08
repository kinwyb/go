//package gosql mssql工具包..引用"github.com/denisenkom/go-mssqldb"
package mssql

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

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
