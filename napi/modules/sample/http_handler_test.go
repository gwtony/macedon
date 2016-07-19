package sample

import (
	//"io"
	"strings"
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/test"
	//"net/http"
)

func TestAddHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	h := InitHandler("/test", log)
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &AddHandler{h: h, mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("add handler server http done")
}

func TestDeleteHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	h := InitHandler("/test", log)
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &DeleteHandler{h: h, mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("delete handler server http done")
}

func TestReadHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &ReadHandler{mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("read handler server http done")
}

func TestAddServerHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &AddServerHandler{mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("add server handler server http done")
}

func TestDeleteServerHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &DeleteServerHandler{mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("delete server handler server http done")
}

func TestReadServerHandler1(t *testing.T) {
	body := strings.NewReader("{}")
	w, r := test.Test_generate_rr("GET", "/test", body)
	log := test.TestInitlog()
	mc := InitMysqlContext("localhost", "test", "test", "test", log)
	ah := &ReadServerHandler{mc: mc, log: log}
	ah.ServeHTTP(w, r)
	t.Log("read server handler server http done")
}

