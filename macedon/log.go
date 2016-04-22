package macedon

import (
	"fmt"
	"code.google.com/p/log4go"
)

type Log struct {
	l log4go.Logger
}

func GetLogger(path string, level string) *Log {
	var log Log


	if path == "" {
		fmt.Println("path is ", path)
		path = DEFAULT_LOG_PATH
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

	/* for test */
	l := log4go.NewConsoleLogger(lv)

	//l := make(log4go.Logger)
	flw := log4go.NewFileLogWriter(path, false)
	flw.SetFormat("[%D %T] [%L] %M")
	//flw.SetRotate(true)
	//flw.SetRotateLines(50)
	flw.SetRotateDaily(true)
	l.AddFilter("log", lv, flw)

	log.l = l

	return &log
}

func (l *Log) Info(arg0 interface{}, args ...interface{}) {
	l.l.Info(arg0, args...)
}

func (l *Log) Error(arg0 interface{}, args ...interface{}) {
	l.l.Error(arg0, args...)
}

func (l *Log) Debug(arg0 interface{}, args ...interface{}) {
	l.l.Debug(arg0, args...)
}
