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
		ret := &logger{
			lg:    l,
			level: lf.level,
		}
		lf.logmap.Store(filename, ret)
		return ret
	}
}

//设置输出日志等级
func (lf *LogFiles) Level(level Level) {
	lf.level = level
	lf.logmap.Range(func(key, value interface{}) bool {
		if v, ok := value.(*logger); ok {
			v.level = level
		}
		return true
	})
}

func (lf *LogFiles) Notice(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Notice(format, args...)
}

//输出
func (lf *LogFiles) Debug(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Debug(format, args...)
}

//输出
func (lf *LogFiles) Info(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Info(format, args...)
}

//警告
func (lf *LogFiles) Warning(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Warning(format, args...)
}

//错误
func (lf *LogFiles) Error(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Error(format, args...)
}

//关键
func (lf *LogFiles) Critical(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Critical(format, args...)
}

//警报
func (lf *LogFiles) Alert(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Alert(format, args...)
}

//紧急
func (lf *LogFiles) Emergency(filename, format string, args ...interface{}) {
	lf.GetLog(filename).Emergency(format, args...)
}
