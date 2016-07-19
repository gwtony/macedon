package hserver

import (
	//"fmt"
	"io"
	"time"
	"net/http"
	//"io/ioutil"
	"encoding/json"
	"git.lianjia.com/lianjia-sysop/napi/variable"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
	"git.lianjia.com/lianjia-sysop/napi/router"
)

type HttpServer struct {
	addr        string
	location    string

	router      *router.Router

	log         log.Log
}

func InitHttpServer(addr string, log log.Log) (*HttpServer, error) {
	hs := &HttpServer{}

	hs.addr = addr
	hs.log  = log

	hs.router = router.InitRouter(log)

	return hs, nil
}

func (hs *HttpServer) AddRouter(url string, h http.Handler) error {
	return hs.router.AddRouter(url, h)
}


func (hs *HttpServer) Run() error {
	s := &http.Server{
		Addr:           hs.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        hs.router,
	}

	return s.ListenAndServe()
}

func ReturnError(w http.ResponseWriter, err error, log log.Log) {
	if err == errors.NoContentError {
		log.Debug("Request no content")
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err == errors.BadRequestError {
		log.Debug("Return bad request")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err == errors.BadGatewayError {
		log.Debug("Return bad gateway")
		http.Error(w, "", http.StatusBadGateway)
		return
	}

	log.Debug("Return internal server error")
	http.Error(w, "", http.StatusInternalServerError)
}

func ReturnResponse(w http.ResponseWriter, resp interface{}, log log.Log) {
	if resp == nil {
		log.Debug("Return OK")
		w.WriteHeader(http.StatusOK)
		return
	}
	respj, err := json.Marshal(resp)
	if err != nil {
		log.Error("Encode json failed: ", resp)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
	w.WriteHeader(http.StatusOK)

	log.Debug("Return OK")

	io.WriteString(w, string(respj))
}

