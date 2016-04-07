package macedon

import (
	//"fmt"
	"time"
	"net/http"
)

type HttpServer struct {
	s           *Server

	addr        string
	location    string

	creater     *CreateHandler
	deleter     *DeleteHandler
	updater     *UpdateHandler
	reader      *ReadHandler
	notifyer    *NotifyHandler

	log         *Log
}

func InitHttpServer(addr string, log *Log) (*HttpServer, error) {
	hs := &HttpServer{}

	hs.addr = addr
	hs.log  = log

	hs.creater = &CreateHandler{}
	hs.creater.hs = hs
	hs.creater.log = log

	hs.deleter = &DeleteHandler{}
	hs.deleter.hs = hs
	hs.deleter.log = log

	hs.updater = &UpdateHandler{}
	hs.updater.hs = hs
	hs.updater.log = log

	hs.reader  = &ReadHandler{}
	hs.reader.hs = hs
	hs.reader.log = log

	hs.notifyer = &NotifyHandler{}
	hs.notifyer.hs = hs
	hs.notifyer.log = log

	return hs, nil
}

func (hs *HttpServer) Run() error {
	s := &http.Server{
		Addr:           hs.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.ListenAndServe()
}

func (hs *HttpServer) AddRouter(location string) error {
	hs.log.Debug("Add router %s", location)

	//TODO: deal invalid location
	hs.location = location

	http.Handle(location + DEFAULT_CREATE_LOCATION, hs.creater)
	http.Handle(location + DEFAULT_DELETE_LOCATION, hs.deleter)
	http.Handle(location + DEFAULT_UPDATE_LOCATION, hs.updater)
	http.Handle(location + DEFAULT_READ_LOCATION, hs.reader)
	http.Handle(location + DEFAULT_NOTIFY_LOCATION, hs.notifyer)

	return nil
}

func (hs *HttpServer) Server() (*Server) {
	return hs.s
}
