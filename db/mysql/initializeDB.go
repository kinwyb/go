package mysql

import (
	"database/sql"
	"strings"

	"github.com/kinwyb/go/db"

	"github.com/kinwyb/go/logs"

	"github.com/kinwyb/go/err1"
)

var dbhost string
var dbusername string
var dbpassword string
var dbname string
var conn *Conn
var lg logs.Logger

//设置数据库基础连接
func InitializeDB(host, username, password, name string, log ...logs.Logger) {
	dbhost = host
	dbusername = username
	dbpassword = password
	dbname = name
	if len(log) > 0 {
		lg = log[0]
	}
}

//获取数据库连接
func GetDBConnect() *Conn {
	if conn == nil {
		conn = InitializeConnect(dbhost, dbusername, dbpassword, dbname, lg)
	}
	return conn
}

var notInitializeQueryResult = db.ErrQueryResult(db.DatabaseNotInitialize)
var connectFailQueryResult = db.ErrQueryResult(db.DatabaseConnectFail)
var notInitializeExecResult = db.ErrExecResult(db.DatabaseNotInitialize)
var connectFailExecResult = db.ErrExecResult(db.DatabaseConnectFail)

//数据库连接
type Conn struct {
	conn        db.SQL      //数据库连接
	dbname      string      //数据库名称
	host        string      //数据库地址
	username    string      //数据库用户名
	password    string      //数据库密码
	connectSucc bool        //数据库连接是否成功
	lg          logs.Logger //日志
}

func (d *Conn) FormatError(e error) err1.Error {
	if m, ok := d.conn.(*mysql); ok {
		return m.FormatError(e)
	}
	return err1.NewError(-1, e.Error())
}

//获取完整表名[附带数据库名称]
func (d *Conn) Table(tbname string) string {
	if d == nil || d.conn == nil {
		return tbname
	}
	return d.conn.Table(tbname)
}

//Rows 查询多条数据,结果以[]map[string]interface{}方式返回
//返回结果,使用本package中的类型函数进行数据解析
//eg:
//		result := QueryRow(...)
//		result.Error(func(error.Error){
//			这里处理错误
// 		}).Rows(func(map[string]interface{}) bool {
//			return true //返回true，继续循环读取下一条
// 		})
//param sql string SQL
//param args... interface{} SQL参数
func (d *Conn) QueryRows(sql string, args ...interface{}) db.QueryResult {
	if d == nil || d.conn == nil {
		return notInitializeQueryResult
	} else if !d.connectSucc {
		return connectFailQueryResult
	}
	return d.conn.QueryRows(sql, args...)
}

//Rows 查询多条数据,结果以[]map[string]interface{}方式返回
//返回结果,使用本package中的类型函数进行数据解析
//eg:
//		result := QueryRow(...)
//		result.Error(func(error.Error){
//			这里处理错误
// 		}).Rows(func(map[string]interface{}) bool {
//			return true //返回true，继续循环读取下一条
// 		})
//param sql string SQL
//param args... interface{} SQL参数
func (d *Conn) QueryRow(sql string, args ...interface{}) db.QueryResult {
	if d == nil || d.conn == nil {
		return notInitializeQueryResult
	} else if !d.connectSucc {
		return connectFailQueryResult
	}
	return d.conn.QueryRow(sql, args...)
}

//Exec 执行一条SQL
//param sql string SQL
//param args... interface{} SQL参数
func (d *Conn) Exec(sql string, args ...interface{}) db.ExecResult {
	if d == nil || d.conn == nil {
		return notInitializeExecResult
	} else if !d.connectSucc {
		return connectFailExecResult
	}
	return d.conn.Exec(sql, args...)
}

//Count SQL语句条数统计
//param sql string SQL
//param args... interface{} SQL参数
func (d *Conn) Count(sql string, args ...interface{}) (int64, err1.Error) {
	if d == nil || d.conn == nil {
		return 0, db.DatabaseNotInitialize
	} else if !d.connectSucc {
		return 0, db.DatabaseConnectFail
	}
	return d.conn.Count(sql, args...)
}

//RowsPage 分页查询SQL
//返回结果,使用本package中的类型函数进行数据解析
//eg:
//		result := QueryRow(...)
//		result.Error(func(error.Error){
//			这里处理错误
// 		}).Rows(func(map[string]interface{}) bool {
//			return true //返回true，继续循环读取下一条
// 		})
//param sql string SQL
//param page *PageObj 分页数据
//param args... interface{} SQL参数
func (d *Conn) QueryWithPage(sql string, page *db.PageObj, args ...interface{}) db.QueryResult {
	if d == nil || d.conn == nil {
		return notInitializeQueryResult
	} else if !d.connectSucc {
		return connectFailQueryResult
	}
	return d.conn.QueryWithPage(sql, page, args...)
}

//ParseSQL 解析SQL
//param sql string SQL
//param args map[string]interface{} 参数映射
func (d *Conn) ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, err1.Error) {
	if d == nil || d.conn == nil {
		return "", nil, db.DatabaseNotInitialize
	} else if !d.connectSucc {
		return "", nil, db.DatabaseConnectFail
	}
	return d.conn.ParseSQL(sql, args)
}

func (d *Conn) Prepare(query string) (*sql.Stmt, err1.Error) {
	if d == nil || d.conn == nil {
		return nil, db.DatabaseNotInitialize
	} else if !d.connectSucc {
		return nil, db.DatabaseConnectFail
	}
	return d.conn.Prepare(query)
}

//Transaction 事务处理
//param t TransactionFunc 事务处理函数
func (d *Conn) Transaction(t db.TransactionFunc, new ...bool) err1.Error {
	if d == nil || d.conn == nil {
		return db.DatabaseNotInitialize
	} else if !d.connectSucc {
		return db.DatabaseConnectFail
	}
	return d.conn.Transaction(t)
}

func (d *Conn) DataBaseName() string {
	if d == nil || d.conn == nil || !d.connectSucc {
		return ""
	}
	return d.conn.DataBaseName()
}

//GetDb 获取数据库对象
func (d *Conn) GetDb() (*sql.DB, err1.Error) {
	if d == nil || d.conn == nil {
		return nil, db.DatabaseNotInitialize
	} else if !d.connectSucc {
		return nil, db.DatabaseConnectFail
	}
	return d.conn.GetDb()
}

//获取db.SQL
func (d *Conn) GetConn() db.SQL {
	if d == nil || d.conn == nil {
		return nil
	} else if !d.connectSucc {
		d.ReConnect()
	}
	return d.conn
}

//Close 关闭数据库
func (d *Conn) Close() {
	if d == nil || d.conn == nil {
		return
	}
	d.conn.Close()
	d.conn = nil
	d.connectSucc = false
}

//再次连接数据库
func (d *Conn) ReConnect() {
	if d == nil {
		return
	}
	d.connectSucc = false
	var err error
	d.conn, err = Connect(d.host, d.username, d.password, d.dbname)
	if err != nil && d.lg != nil {
		d.lg.Error("数据库连接失败:%s", err.Error())
	} else {
		d.connectSucc = true
	}
}

//初始化数据库
func InitializeConnect(host, username, password, dbname string, log ...logs.Logger) *Conn {
	if !strings.Contains(host, ":") { //不带端口给加上默认端口
		host = host + ":3306"
	}
	ret := &Conn{
		dbname:   dbname,
		host:     host,
		username: username,
		password: password,
	}
	if len(log) > 0 {
		ret.lg = log[0]
	}
	ret.ReConnect()
	return ret
}
