// +build orcale

//package gosql orcale工具包..引用"github.com/mattn/go-oci8"
package orcale

import (
	"database/sql"
	"time"

	"github.com/kinwyb/go/db"
)

//orcale 操作对象
type orcale struct {
	db.Conn
}

//链接orcale数据库
func Connect(host, username, password, db string) (db.SQL, error) {
	linkstring := username + "/" + password + host + "/" + db
	result := &orcale{}
	sqlDB, err := sql.Open("oci8", linkstring)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(1 * time.Hour) //一个小时后重置链接
	result.SetSQLDB(sqlDB)
	result.SetDataBaseName(db) //记录数据库名称,表名格式化会用到
	return result, nil
}
