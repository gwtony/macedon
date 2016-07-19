package router

import (
	//"net/http"
	//"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/test"
	"testing"

)

func TestInitRouterOk(t *testing.T) {
//	w, _ := test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	r := InitRouter(log)
	if r == nil {
		t.Fatal("init router failed")
	}
	t.Log("init router done")
}

func TestAddRouterOk(t *testing.T) {
	log := test.TestInitlog()
	r := InitRouter(log)
	r.AddRouter("/test", &test.Thandler{})
	t.Log("add router done")
}
func TestServeHTTPOk(t *testing.T) {
	w, req := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	r := InitRouter(log)
	r.AddRouter("/test", &test.Thandler{})
	r.ServeHTTP(w, req)
	t.Log("serve http done")
}
