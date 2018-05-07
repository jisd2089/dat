package log

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"drcs/log/logs"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// 默认日志
var defaultLogger = &LoggerOnLogrus{logrus.StandardLogger(), CallerDepth + 1}

// Logger 日志接口
type Logger interface {
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

//IsError 判断对象是否实现了error接口
func IsError(obj interface{}) bool {
	objType := reflect.TypeOf(obj)
	return objType.Implements(errorType)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// GetLogger 获取指定日志实例
func GetLogger(loggerName string) Logger {
	logrusLogger := logs.GetLogger(loggerName)
	return NewLoggerOnLogrus(logrusLogger)
}
