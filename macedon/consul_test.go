package macedon

import (
	"fmt"
	"strings"
	"net/http"
	"testing"
	"net/http/httptest"
)

func testCreateConsulContext(t *testing.T, server string, log *Log) (* ConsulContext) {
	cc, err := InitConsulContext(server, DEFAULT_REGISTER_LOC, DEFAULT_DEREGISTER_LOC, DEFAULT_READ_LOC, log)
	if err != nil {
		t.Fatal("init consul context failed")
	}

	return cc
}

func testCreateConsulServer(data string) (* httptest.Server) {
	 return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, data)
	}))
}

func TestRegisterServiceFailed(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	cc := testCreateConsulContext(t, "127.0.0.1", log)
	err := cc.RegisterService("test", "192.168.0.1")
	if err != nil && err == BadGatewayError {
		t.Log("register service error badgateway")
	}
}

func TestDeRegisterServiceFailed(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	cc := testCreateConsulContext(t, "127.0.0.1", log)
	err := cc.DeRegisterService("test", "192.168.0.1")
	if err != nil && err == BadGatewayError {
		t.Log("register service error badgateway")
	}
}

func TestListServiceFailed(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	cc := testCreateConsulContext(t, "127.0.0.1", log)
	_, err := cc.ListService("test", "192.168.0.1")
	if err != nil && err == BadGatewayError {
		t.Log("register service error badgateway")
	}
}

func TestDeRegisterServiceNotFound(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	err := cc.DeRegisterService("test", "192.168.0.1")
	if err != nil && err != NoContentError {
		t.Fatal("deregister service error: ", err)
	}
	t.Log("Deregister service not found")
}
func TestDeRegisterServiceOk(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[{\"ServiceName\": \"test\", \"ServiceAddress\": \"192.168.0.1\", \"ServiceId\": \"1\"}]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	err := cc.DeRegisterService("test", "192.168.0.1")
	if err != nil {
		t.Fatal("deregister service error: ", err)
	}
	t.Log("deregister service done")
}
func TestDeRegisterServiceOk_2(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[{\"ServiceName\": \"test\", \"ServiceAddress\": \"192.168.0.1\", \"ServiceId\": \"1\"}]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	err := cc.DeRegisterService("test", "")
	if err != nil {
		t.Fatal("deregister service error: ", err)
	}
	t.Log("deregister service done")
}
