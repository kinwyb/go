package logs

import (
	"fmt"

	"log"

	"os"
	"path/filepath"
	"strings"
	"time"
)

type Level int

const (
	Emergency Level = iota
	Alert
	Critical
	Error
	Warn
	Info
	Debug
)

//Logger 日志接口
type Logger interface {
	//输出
	Debug(format string, args ...interface{})
	//输出
	Info(format string, args ...interface{})
	//警告
	Warning(format string, args ...interface{})
	//错误
	Error(format string, args ...interface{})
	//关键
	Critical(format string, args ...interface{})
	//警报
	Alert(format string, args ...interface{})
	//紧急
	Emergency(format string, args ...interface{})
}

//日志对象
type logger struct {
	logger   *log.Logger
	filename string
	filedir  string
	file     *os.File
	t        time.Duration
	level    Level
}

func NewFileLogger(filename string, t time.Duration, level ...Level) Logger {
	fd, e := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive)
	if nil == e {
		ret := &logger{
			logger:   log.New(fd, "", log.LstdFlags),
			file:     fd,
			filename: filename,
			filedir:  "./",
			t:        t,
			level:    Debug,
		}
		if len(level) > 0 {
			ret.level = level[0]
		}
		if index := strings.LastIndex(filename, "/"); index != -1 {
			ret.filedir = filename[0:index] + "/"
			ret.filedir, _ = filepath.Abs(ret.filedir)
			ret.filename = filename[index:]
			os.MkdirAll(filename[0:index], os.ModePerm)
		}
		go createLogFile(ret)
		return ret
	} else {
		fmt.Printf("%s", e.Error())
	}
	return NewBeegoFileLog(1, filename, 10)
}

func NewLogger(level ...Level) Logger {
	ret := &logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  Debug,
	}
	if len(level) > 0 {
		ret.level = level[0]
	}
	return ret
}

//输出
func (lg *logger) Debug(format string, args ...interface{}) {
	if lg.level < Debug {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[D] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[34m[D] "+format+"\x1b[0m", args...)
	}
}

//输出
func (lg *logger) Info(format string, args ...interface{}) {
	if lg.level < Info {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[I] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[36m[I] "+format+"\x1b[0m", args...)
	}
}

//警告
func (lg *logger) Warning(format string, args ...interface{}) {
	if lg.level < Warn {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[W] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[33m[W] "+format+"\x1b[0m", args...)
	}
}

//错误
func (lg *logger) Error(format string, args ...interface{}) {
	if lg.level < Error {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[E] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[31m[E] "+format+"\x1b[0m", args...)
	}
}

//关键
func (lg *logger) Critical(format string, args ...interface{}) {
	if lg.level < Critical {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[C] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[1m\x1b[31m[C] "+format+"\x1b[0m\x1b[21m", args...)
	}
}

//警报
func (lg *logger) Alert(format string, args ...interface{}) {
	if lg.level < Alert {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[A] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[35m[A] "+format+"\x1b[0m", args...)
	}
}

//紧急
func (lg *logger) Emergency(format string, args ...interface{}) {
	if lg.level < Emergency {
		return
	}
	if lg.filedir != "" {
		lg.logger.Printf("[EG] "+format, args...)
	} else {
		lg.logger.Printf("\x1b[1m\x1b[35m[EG] "+format+"\x1b[0m\x1b[21m", args...)
	}
}

func createLogFile(lg *logger) {
	t := time.NewTimer(1)
	now := time.Now()
	first := true
	for {
		if first { //第一次等待时间跳过
			<-t.C
			first = false
		}
		switch lg.t {
		case time.Minute:
			next := now.Add(time.Minute)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
			t.Reset(next.Sub(now))
		case time.Hour:
			next := now.Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
			t.Reset(next.Sub(now))
		default:
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t.Reset(next.Sub(now))
		}
		<-t.C
		now = time.Now()
		filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d%02d.log", lg.filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		if err := os.Rename(lg.filedir+"/"+lg.filename, lg.filedir+"/"+filename); err != nil {
			fmt.Printf("文件命名失败:%s", err.Error())
			//日志文件重命名失败
		} else {
			os.Chmod(lg.filedir+"/"+filename, os.ModePerm) //修改权限
			if fd, err := os.OpenFile(lg.filedir+"/"+lg.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
				lg.logger.SetOutput(fd)
				lg.file.Sync()
				lg.file.Close()
				lg.file = fd
			} else {
				fmt.Printf("新文件创建失败:%s\n", err.Error())
			}
		}
	}
}

//WriteLog 写入日志
func WriteLog(log Logger, level Level, format string, args ...interface{}) {
	if log == nil {
		return
	}
	switch level {
	case Emergency:
		log.Emergency(format, args...)
	case Alert:
		log.Alert(format, args...)
	case Critical:
		log.Critical(format, args...)
	case Error:
		log.Error(format, args...)
	case Warn:
		log.Warning(format, args...)
	case Info:
		log.Info(format, args...)
	case Debug:
		log.Debug(format, args...)
	}
}
