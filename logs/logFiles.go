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
	level    Level
}

//初始化一个日志集合,filepath日志保存路径,如果为空直接输出到屏幕
//t日志文件分割时间,比如每天一个文件或每小时一个文件等
func NewLogFiles(filepath string, t time.Duration, level ...Level) *LogFiles {
	ret := &LogFiles{
		filepath: filepath,
		t:        t,
		logmap:   &sync.Map{},
		level:    Debug,
	}
	if len(level) > 0 {
		ret.level = level[0]
	}
	return ret
}

//设置输出日志等级
func (lf *LogFiles) Level(filename string, level Level) {
	if v, ok := lf.logmap.Load(filename); ok {
		if x, ok := v.(logger); ok {
			x.level = level
		}
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, level)
		}
		lf.logmap.Store(filename, l)
	}
}

//输出
func (lf *LogFiles) Debug(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Debug(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Debug(format, args...)
	}
}

//输出
func (lf *LogFiles) Info(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Info(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Info(format, args...)
	}
}

//警告
func (lf *LogFiles) Warning(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Warning(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Warning(format, args...)
	}
}

//错误
func (lf *LogFiles) Error(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Error(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Error(format, args...)
	}
}

//关键
func (lf *LogFiles) Critical(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Critical(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Critical(format, args...)
	}
}

//警报
func (lf *LogFiles) Alert(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Alert(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Alert(format, args...)
	}
}

//紧急
func (lf *LogFiles) Emergency(filename, format string, args ...interface{}) {
	if v, ok := lf.logmap.Load(filename); ok {
		v.(Logger).Emergency(format, args...)
	} else {
		var l Logger
		if lf.filepath == "" {
			l = NewLogger(lf.level)
		} else {
			l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		}
		lf.logmap.Store(filename, l)
		l.Emergency(format, args...)
	}
}
