package macedon

import (
	"io"
	"net/http"
)

type CreateHandler struct {
	hs *HttpServer
}
type DeleteHandler struct {
	hs *HttpServer
}
type UpdateHandler struct {
	hs *HttpServer
}
type ReadHandler struct {
	hs *HttpServer
}

func (h* CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
