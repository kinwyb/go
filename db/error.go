package db

import "github.com/kinwyb/go/err1"

var (
	DatabaseNotInitialize = err1.NewError(101, "数据库连接尚未初始化")
	DatabaseConnectFail   = err1.NewError(102, "数据库连接失败")
	SQLError              = err1.NewError(103, "数据库操作异常")
	SQLEmptyChange        = err1.NewError(104, "数据无变化")
)
