package db

import (
	"errors"
)

var (
	DatabaseNotInitialize = errors.New("数据库连接尚未初始化")
	DatabaseConnectFail   = errors.New("数据库连接失败")
	SQLError              = errors.New("数据库操作异常")
	SQLEmptyChange        = errors.New("数据无变化")
)
