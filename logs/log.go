package logs

import (
	"fmt"

	"log"

	"os"
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

func NewFileLogger(filename string, t time.Duration) Logger {
	if fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
		ret := &logger{
			logger:   log.New(fd, "", log.LstdFlags),
			file:     fd,
			filename: filename,
			filedir:  "./",
			t:        t,
		}
		if index := strings.LastIndex(filename, "/"); index != -1 {
			ret.filedir = filename[0:index] + "/"
			ret.filename = filename[index:]
			os.MkdirAll(filename[0:index], os.ModePerm)
		}
		go createLogFile(ret)
		return ret
	}
	return NewBeegoFileLog(1, filename, 10)
}

func NewLogger() Logger {
	return &logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

//输出
func (lg *logger) Info(format string, args ...interface{}) {
	lg.logger.Printf("[I] "+format, args...)
}

//警告
func (lg *logger) Warning(format string, args ...interface{}) {
	lg.logger.Printf("[W] "+format, args...)
}

//错误
func (lg *logger) Error(format string, args ...interface{}) {
	lg.logger.Printf("[E] "+format, args...)
}

//关键
func (lg *logger) Critical(format string, args ...interface{}) {
	lg.logger.Printf("[C] "+format, args...)
}

//警报
func (lg *logger) Alert(format string, args ...interface{}) {
	lg.logger.Printf("[A] "+format, args...)
}

//紧急
func (lg *logger) Emergency(format string, args ...interface{}) {
	lg.logger.Printf("[E] "+format, args...)
}

func createLogFile(lg *logger) {
	t := time.NewTimer(1)
	now := time.Now()
	for {
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
		fmt.Printf("开始创建新文件....\n")
		now = time.Now()
		filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d%02d.log", lg.filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		if err := os.Rename(lg.filedir+"/"+lg.filename, lg.filedir+"/"+filename); err != nil {
			fmt.Printf("文件命名失败:%s", err.Error())
			//日志文件重命名失败
		} else {
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
