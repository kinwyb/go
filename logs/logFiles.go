package logs

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
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
	if filepath != "" { //创建文件文件夹
		os.MkdirAll(filepath, 0777)
	}
	if len(level) > 0 {
		ret.level = level[0]
	}
	return ret
}

//获取指定日志
func (lf *LogFiles) GetLog(filename string) Logger {
	if v, ok := lf.logmap.Load(filename); ok {
		return v.(Logger)
	} else {
		l := logs.NewLogger(3000)
		if lf.filepath != "" {
			l.SetLogger(logs.AdapterFile, `{"filename":"`+lf.filepath+"/"+filename+`","level":`+
				fmt.Sprintf("%d", lf.level)+`,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
		}
		lf.logmap.Store(filename, l)
		return l
		//var l Logger
		//if lf.filepath == "" {
		//	l = NewLogger()
		//} else {
		//	l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, lf.level)
		//}
		//lf.logmap.Store(filename, l)
		//return l
	}
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
			//l = NewFileLogger(filepath.Join(lf.filepath, filename), lf.t, level)
			l := logs.NewLogger(3000)
			if lf.filepath != "" {
				l.SetLogger(logs.AdapterFile, `{"filename":"`+lf.filepath+"/"+filename+`","level":`+
					fmt.Sprintf("%d", lf.level)+`,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
			}
		}
		lf.logmap.Store(filename, l)
	}
}

//输出
func (lf *LogFiles) Debug(filename, format string, args ...interface{}) {
	if lf.level >= Debug {
		l := lf.GetLog(filename)
		l.Debug(format, args...)
	}
}

//输出
func (lf *LogFiles) Info(filename, format string, args ...interface{}) {
	if lf.level >= Info {
		l := lf.GetLog(filename)
		l.Info(format, args...)
	}
}

//警告
func (lf *LogFiles) Warning(filename, format string, args ...interface{}) {
	if lf.level >= Warn {
		l := lf.GetLog(filename)
		l.Warning(format, args...)
	}
}

//错误
func (lf *LogFiles) Error(filename, format string, args ...interface{}) {
	if lf.level >= Error {
		l := lf.GetLog(filename)
		l.Error(format, args...)
	}
}

//关键
func (lf *LogFiles) Critical(filename, format string, args ...interface{}) {
	if lf.level >= Critical {
		l := lf.GetLog(filename)
		l.Critical(format, args...)
	}
}

//警报
func (lf *LogFiles) Alert(filename, format string, args ...interface{}) {
	if lf.level >= Alert {
		l := lf.GetLog(filename)
		l.Alert(format, args...)
	}
}

//紧急
func (lf *LogFiles) Emergency(filename, format string, args ...interface{}) {
	if lf.level >= Emergency {
		l := lf.GetLog(filename)
		l.Emergency(format, args...)
	}
}
