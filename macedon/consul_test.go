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

func testCreateConsulServer500(data string) (* httptest.Server) {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func TestRegisterServiceConnectConsulFailed(t *testing.T) {
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

func TestDeRegisterServiceConnectConsulFailed(t *testing.T) {
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

func TestListServiceConnectConsulFailed(t *testing.T) {
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

func TestListServiceEmptyResult(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	_, err := cc.ListService("test", "192.168.0.1")
	if err != nil && err == NoContentError {
		t.Log("register service error NoContentError")
	}
}

func TestListServiceFailedBadResponseJson(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("\"Name\": {]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	_, err := cc.ListService("test", "192.168.0.1")
	if err != nil && err == InternalServerError {
		t.Log("register service error InternalServerError")
	}
}

func TestDeRegisterServiceEmptyResult(t *testing.T) {
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

func TestDeRegisterServiceAddressNotMatch(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[{\"ServiceName\": \"test\", \"ServiceAddress\": \"192.168.0.1\", \"ServiceId\": \"1\"}]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	err := cc.DeRegisterService("test", "192.168.0.2")
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
func TestDeRegisterServiceNameOnly(t *testing.T) {
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

func TestOperateServiceRegisterWithNoId(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[{\"ServiceName\": \"test\", \"ServiceAddress\": \"192.168.0.1\", \"ServiceId\": \"1\"}]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)

	_, err := cc.OperateService("test", "192.168.0.1", "", REGISTER)
	if err != InternalServerError {
		t.Fatal("operate service check error")
	}
	t.Log("operate service internal server error done")
}


func TestOperateServiceEmptyResult(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)

	_, err := cc.OperateService("test", "192.168.0.1", "", READ)
	if err == nil || err != NoContentError {
		t.Fatal("operate service check error", err)
	}
	t.Log("operate service internal server error done")
}

func TestOperateServiceBadHttpCode(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer500("")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	_, err := cc.OperateService("test", "192.168.0.1", "1", REGISTER)
	if err != BadGatewayError {
		t.Fatal("operate service check error", err)
	}
	t.Log("operate service bad gateway error done")
}

func TestOperateServiceBadResponseJson(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	ts := testCreateConsulServer("[{\"ServiceName\": \"test\", [ \"ServiceAddress\": \"192.168.0.1\"; \"ServiceId\": \"1\"}]")
	defer ts.Close()

	cc := testCreateConsulContext(t, strings.Trim(ts.URL, "http://"), log)
	_, err := cc.OperateService("test", "192.168.0.1", "1", READ)
	if err != InternalServerError {
		t.Fatal("operate service check error", err)
	}
	t.Log("operate service internal server error done")
}

