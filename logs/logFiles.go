package logs

import (
	"path/filepath"
	"sync"
	"time"
)

//一个文件日志集合

type LogFiles struct {
	filepath string
	t        time.Duration //每个日志文件记录是时间间隔，比如每天一个文件或每小时一个文件等
	logmap   *sync.Map     //日志集合，key代表日志文件名，Logger日志记录器

}

//初始化一个日志集合,filepath日志保存路径,t日志文件分割时间
//比如每天一个文件或每小时一个文件等
func NewLogFiles(filepath string, t time.Duration) *LogFiles {
	return &LogFiles{
		filepath: filepath,
		t:        t,
		logmap:   &sync.Map{},
	}
}

//输出
func (lf *LogFiles) Info(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Info(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Info(format, args...)
	}
}

//警告
func (lf *LogFiles) Warning(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Warning(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Warning(format, args...)
	}
}

//错误
func (lf *LogFiles) Error(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Error(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Error(format, args...)
	}
}

//关键
func (lf *LogFiles) Critical(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Critical(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Critical(format, args...)
	}
}

//警报
func (lf *LogFiles) Alert(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Alert(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Alert(format, args...)
	}
}

//紧急
func (lf *LogFiles) Emergency(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Emergency(format, args...)
	} else {
		l := NewFileLogger(filepath.Join(lf.filepath, filename), lf.t)
		lf.logmap.Store(filename, l)
		l.Emergency(format, args...)
	}
}
