package db

import (
	"database/sql"
	"errors"
)

//一个简便的数据库操作封装
var ErrorNotOpen = errors.New("数据库连接失败")
var ErrorPing = errors.New("数据库网络异常")

// 事务参数
type TxOption struct {
	New    bool           //新事务
	Option *sql.TxOptions //事务参数
}

//TransactionFunc 事务回调函数
type TransactionFunc func(tx TxSQL) error

//数据库操作接口
type SQL interface {
	Query
	//ParseSQL 解析SQL
	//param sql string SQL
	//param args map[string]interface{} 参数映射
	ParseSQL(sql string, args map[string]interface{}) (string, []interface{}, error)
	//GetDb 获取数据库对象
	GetDb() (*sql.DB, error)
	//Close 关闭数据库
	Close()
}

//事物数据操作接口
type TxSQL interface {
	Query
	//GetTx 获取事务对象
	GetTx() *sql.Tx
}

//查询操作集合
type Query interface {
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
	QueryRows(sql string, args ...interface{}) QueryResult
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
	QueryRow(sql string, args ...interface{}) QueryResult
	//Exec 执行一条SQL
	//param sql string SQL
	//param args... interface{} SQL参数
	Exec(sql string, args ...interface{}) ExecResult
	//Count SQL语句条数统计
	//param sql string SQL
	//param args... interface{} SQL参数
	Count(sql string, args ...interface{}) (int64, error)
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
	QueryWithPage(sql string, page *PageObj, args ...interface{}) QueryResult
	//Prepare 预处理
	Prepare(query string) (*sql.Stmt, error)
	//格式化表名称
	Table(tbname string) string
	//数据库名称
	DataBaseName() string
	// 设置最大连接数
	SetMaxOpenConns(n int)
	//Transaction 事务处理
	//param t TransactionFunc 事务处理函数
	//param new bool 是否创建新事物,默认false,如果设置true不管事务是否存在都会创建新事物
	Transaction(t TransactionFunc, option ...*TxOption) error
}
