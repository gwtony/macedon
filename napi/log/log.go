package log

import (
	"code.google.com/p/log4go"
	"git.lianjia.com/lianjia-sysop/napi/variable"
)
type Log interface {
	Info(arg0 interface{}, args ...interface{})
	Error(arg0 interface{}, args ...interface{})
	Debug(arg0 interface{}, args ...interface{})
}

type Logger struct {
	l log4go.Logger
}

func GetLogger(path string, level string) *Logger {
	var log Logger

	if path == "" {
		path = variable.DEFAULT_LOG_PATH
	}

	lv := log4go.ERROR
	switch level {
	case "debug":
		lv = log4go.DEBUG
	case "error":
		lv = log4go.ERROR
	case "info":
		lv = log4go.INFO
	}

	l := log4go.NewDefaultLogger(lv)
	flw := log4go.NewFileLogWriter(path, false)
	if flw == nil {
		return nil
	}
	flw.SetFormat("[%D %T] [%L] %M")
	//flw.SetRotate(true)
	//flw.SetRotateLines(50)
	flw.SetRotateDaily(true)
	l.AddFilter("log", lv, flw)

	log.l = l

	return &log
}

func (l *Logger) Info(arg0 interface{}, args ...interface{}) {
	l.l.Info(arg0, args...)
}

func (l *Logger) Error(arg0 interface{}, args ...interface{}) {
	l.l.Error(arg0, args...)
}

func (l *Logger) Debug(arg0 interface{}, args ...interface{}) {
	l.l.Debug(arg0, args...)
}
