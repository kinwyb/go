package logs

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

// 输出到文件
func newFileWriter(logPath string, maxRemainCnt uint, log *logrus.Logger) {
	writer, err := rotatelogs.New(
		logPath+"%Y%m%d",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logPath),
		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	log.SetOutput(writer)
	log.ExitFunc = func(code int) {
		writer.Close()
	}
}

// 输出到文件hook
func newFileHook(logPath string, maxRemainCnt uint,
	log *logrus.Logger, format logrus.Formatter) logrus.Hook {
	writer, err := rotatelogs.New(
		logPath+"%Y_%m_%d",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logPath),
		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	if format == nil { //默认json格式
		format = &DefaultJsonFormatter
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, format)
	return lfsHook
}
