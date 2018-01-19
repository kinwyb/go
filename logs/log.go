package logs

import (
	"fmt"

	"log"

	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kinwyb/go/exit"
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

func NewFileLogger(file string, t time.Duration, level ...Level) Logger {
	if file == "" || file == "/" { //如果文件地址异常，默认当前路径下的logs
		file = "logs"
	}
	var filedir = "./"
	var filename = ""
	if index := strings.LastIndex(file, "/"); index != -1 { //创建文件夹
		filedir = file[0:index] + "/"
		filedir, _ = filepath.Abs(filedir)
		if _, err := os.Stat(filedir); os.IsNotExist(err) { //如果目录不存在创建目录
			os.MkdirAll(file[0:index], os.ModePerm)
		}
		filename = file[index+1:] //取到文件名
	}
	fd, e := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive)
	if nil == e {
		ret := &logger{
			logger:   log.New(fd, "", log.LstdFlags),
			file:     fd,
			filename: filename,
			filedir:  filedir,
			t:        t,
			level:    Debug,
		}
		if len(level) > 0 {
			ret.level = level[0]
		}
		exit.Listen(ret.exit) //监听退出
		go ret.createLogFile()
		return ret
	} else {
		ret := NewBeegoFileLog(1, filename, 10)
		ret.Error("创建基础日志驱动失败:%s", e.Error())
		return ret
	}
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

func (lg *logger) createLogFile() {
	defer func() {
		if err := recover(); err != nil {
			lg.Critical("文件创建线程崩溃:%s", err)
			go lg.createLogFile()
		}
	}()
	t := time.NewTimer(1)
	first := true
	for {
		if first { //第一次等待时间跳过
			<-t.C
		}
		now := time.Now()
		filename := lg.filename
		switch lg.t {
		case time.Minute:
			filename = fmt.Sprintf("%s_%04d%02d%02d_%02d%02d",
				lg.filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
			if first { //第一次调用.调整时间到下一分钟触发
				next := now.Add(time.Minute)
				next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
				waittime := next.Sub(now)
				t.Reset(waittime)
				t.Reset(next.Sub(now))
			} else {
				t.Reset(time.Minute) //一分钟之后触发
			}
		case time.Hour:
			filename = fmt.Sprintf("%s_%04d%02d%02d_%02d",
				lg.filename, now.Year(), now.Month(), now.Day(), now.Hour())
			if first {
				next := now.Add(time.Hour)
				next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
				waittime := next.Sub(now)
				t.Reset(waittime)
				t.Reset(next.Sub(now))
			} else {
				t.Reset(time.Hour)
			}
		default:
			filename = fmt.Sprintf("%s_%04d%02d%02d",
				lg.filename, now.Year(), now.Month(), now.Day())
			if first {
				next := now.Add(time.Hour * 24)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				waittime := next.Sub(now)
				t.Reset(waittime)
			} else {
				t.Reset(24 * time.Hour)
			}
		}
		first = false
		<-t.C
		lg.saveFile(filename, true)
	}
}

//退出时更新日志文件名
func (lg *logger) exit(args ...interface{}) {
	filename := ""
	now := time.Now()
	switch lg.t {
	case time.Minute:
		filename = fmt.Sprintf("%s_%04d%02d%02d_%02d%02d",
			lg.filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	case time.Hour:
		filename = fmt.Sprintf("%s_%04d%02d%02d_%02d",
			lg.filename, now.Year(), now.Month(), now.Day(), now.Hour())
	default:
		filename = fmt.Sprintf("%s_%04d%02d%02d",
			lg.filename, now.Year(), now.Month(), now.Day())
	}
	lg.saveFile(filename, false)
}

//保存文件并是否创建下一日志文件
func (lg *logger) saveFile(filename string, createNext bool) {
	now := time.Now()
	for { //防止日志文件名重复
		if _, err := os.Stat(filepath.Join(lg.filedir, filename) + ".log"); err == nil { //如果文件已经存在,合并2个文件内容
			filename = filename + "_" + strings.Replace(now.Format("05.99999"), ".", "", -1)
			continue
		}
		//文件名不存在退出循环
		break
	}
	filename = filepath.Join(lg.filedir, filename) + ".log"
	if err := os.Rename(filepath.Join(lg.filedir, lg.filename), filename); err != nil {
		lg.Warning("文件命名失败:%s", err.Error())
		//日志文件重命名失败
	} else {
		if createNext {
			os.Chmod(lg.filedir+"/"+filename, os.ModePerm) //修改权限
			if fd, err := os.OpenFile(filepath.Join(lg.filedir, lg.filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
				lg.logger.SetOutput(fd)
				lg.file.Sync()
				lg.file.Close()
				lg.file = fd
			} else {
				lg.Warning("新文件创建失败:%s\n", err.Error())
			}
		} else { //关闭之前的文件
			lg.file.Sync()
			lg.file.Close()
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
