//package gosql mysql工具包..引用"github.com/go-sql-driver/mysql"
package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kinwyb/go/db"
)

//mysql 操作对象
type mysql struct {
	db.Conn
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
	sqlDB, err := sql.Open("mysql", linkstring)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(1 * time.Hour) //一个小时后重置链接
	result.SetSQLDB(sqlDB)
	result.SetDataBaseName(db) //记录数据库名称,表名格式化会用到
	return result, nil
}
