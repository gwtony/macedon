package macedon

import (
	"io"
	//"fmt"
	"net/http"
)

type CreateHandler struct {
	hs  *HttpServer
	log *Log
}
type DeleteHandler struct {
	hs  *HttpServer
	log *Log
}
type UpdateHandler struct {
	hs  *HttpServer
	log *Log
}
type ReadHandler struct {
	hs  *HttpServer
	log *Log
}
type NotifyHandler struct {
	hs  *HttpServer
	log *Log
}

func (h* CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mc := h.hs.Server().MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		//
	}
	defer mc.Close(db)

	mc.QueryCreate(db, "qqq.example.org", "A", "192.168.8.8", 5, 86400)
	io.WriteString(w, "create\n")
}
func (h* DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mc := h.hs.Server().MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		//
	}
	defer mc.Close(db)

	mc.QueryDelete(db, "qqq.example.org", "A", "")
	io.WriteString(w, "delete\n")
}
func (h* UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mc := h.hs.Server().MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		//TODO: return 502
	}
	defer mc.Close(db)
	mc.QueryUpdate(db, "qqq.example.org", "A", "192.168.8.8", 1)
	io.WriteString(w, "update\n")
}
func (h* ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mc := h.hs.Server().MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		//TODO: return 502
	}
	defer mc.Close(db)

	mc.QueryRead(db, "test.example.org", "A")

	io.WriteString(w, "read\n")
}
func (h* NotifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "notify\n")
}
