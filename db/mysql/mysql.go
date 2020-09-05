//package gosql mysql工具包..引用"github.com/go-sql-driver/mysql"
package mysql

import (
	"database/sql"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kinwyb/go/db"
)

var rep *regexp.Regexp

func init() {
	rep, _ = regexp.Compile("\\s?Error (\\d+):(.*)")
}

//mysql 操作对象
type mysql struct {
	db.Conn
	linkString string
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
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) //一个小时后重置链接
	result.SetSQLDB(sqlDB)
	result.SetDataBaseName(db) //记录数据库名称,表名格式化会用到
	result.linkString = linkstring
	return result, nil
}

// 解析错误
func formatError(e error) (int64, error) {
	if e == nil {
		return 0, nil
	}
	code := int64(1)
	msg := e.Error()
	if rep.MatchString(msg) {
		d := rep.FindAllStringSubmatch(msg, -1)
		msg = d[0][2]
		cod, err := strconv.ParseInt(d[0][1], 10, 64)
		if err == nil {
			code = cod
		}
	}
	return code, e
}
