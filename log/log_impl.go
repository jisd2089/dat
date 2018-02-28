package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	CallerDepth = 3 // 函数调用栈深度
)

type LoggerOnLogrus struct {
	logrusLogger *logrus.Logger
	callerDepth  int
}

func NewLoggerOnLogrus(logrusLogger *logrus.Logger) *LoggerOnLogrus {
	return &LoggerOnLogrus{logrusLogger: logrusLogger, callerDepth: CallerDepth}
}

func (logger *LoggerOnLogrus) Error(format string, args ...interface{}) {
	if len(args) > 0 && IsError(args[len(args)-1]) {
		err := args[len(args)-1].(error)
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).WithError(err).Errorf(format, args[:len(args)-1]...)
	} else {
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).Errorf(format, args...)
	}
}

func (logger *LoggerOnLogrus) Warn(format string, args ...interface{}) {
	if len(args) > 0 && IsError(args[len(args)-1]) {
		err := args[len(args)-1].(error)
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).WithError(err).Warnf(format, args[:len(args)-1]...)
	} else {
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).Warnf(format, args...)
	}
}

func (logger *LoggerOnLogrus) Info(format string, args ...interface{}) {
	if len(args) > 0 && IsError(args[len(args)-1]) {
		err := args[len(args)-1].(error)
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).WithError(err).Infof(format, args[:len(args)-1]...)
	} else {
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).Infof(format, args...)
	}
}

func (logger *LoggerOnLogrus) Debug(format string, args ...interface{}) {
	if len(args) > 0 && IsError(args[len(args)-1]) {
		err := args[len(args)-1].(error)
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).WithError(err).Debugf(format, args[:len(args)-1]...)
	} else {
		withFileAndLineAndFunc(logger.logrusLogger, logger.callerDepth).Debugf(format, args...)
	}
}

func withFileAndLineAndFunc(logrusLogger *logrus.Logger, callerDepth int) *logrus.Entry {
	pc := make([]uintptr, 1)
	runtime.Callers(callerDepth, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return logrusLogger.WithFields(logrus.Fields{
		"file": frame.File,
		"line": frame.Line,
		"func": frame.Function,
	})
}
