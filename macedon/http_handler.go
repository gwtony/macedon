package macedon

import (
	"io"
	"fmt"
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
	pc := h.hs.Server().PurgeContext()
	sc := h.hs.Server().SshContext()
	fmt.Println(mc, pc, sc)
	pc.DoPurge(sc)

	io.WriteString(w, "create\n")
}
func (h* DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "delete\n")
}
func (h* UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "update\n")
}
func (h* ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "read\n")
}
func (h* NotifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "notify\n")
}
