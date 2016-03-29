package macedon

import (
	"code.google.com/p/log4go"
)

type Log struct {
	l log4go.Logger
}

const DEFAULT_LOG_PATH = "../log/macedon.log"

func GetLogger(path string, level string) *Log {
	var log Log

	if path == "" {
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

	l := log4go.NewDefaultLogger(lv)
	l.AddFilter("log", lv, log4go.NewFileLogWriter(path, false))

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
