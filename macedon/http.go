package macedon

import (
	//"fmt"
	"time"
	"net/http"
)

type HttpServer struct {
	s		*Server

	addr	string

	create	string
	remove  string
	update	string
	read	string

	creater	*CreateHandler
	remover *RemoveHandler
	updater *UpdateHandler
	reader	*ReadHandler

	log		*Log
}

func InitHttpServer(addr string, log *Log) (*HttpServer, error) {
	hs := &HttpServer{}

	hs.addr = addr
	hs.log  = log

	hs.creater = &CreateHandler{}
	hs.creater.hs = hs
	hs.remover = &RemoveHandler{}
	hs.remover.hs = hs
	hs.updater = &UpdateHandler{}
	hs.updater.hs = hs
	hs.reader  = &ReadHandler{}
	hs.reader.hs = hs

	return hs, nil
}

func (hs *HttpServer) Run() error {
	s := &http.Server{
		Addr:			hs.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.ListenAndServe()
}

func (hs *HttpServer) AddRouter(create, remove, update, read string) error {
	hs.log.Debug("Add router", create, remove, update, read)
	time.Sleep(time.Second * 5)
	http.Handle(create, hs.creater)
	hs.log.Debug("Add router done")
	time.Sleep(time.Second * 5)

	http.Handle(remove, hs.remover)
	http.Handle(update, hs.updater)
	http.Handle(read, hs.reader)

	return nil
}
