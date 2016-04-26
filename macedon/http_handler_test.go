package macedon

import (
	"io"
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"
)

func test_generate_rr(method, uri string, body io.Reader) (*httptest.ResponseRecorder, *http.Request){
	r, _ := http.NewRequest(method, uri, body)
	w := httptest.NewRecorder()
	return w, r
}

var testMockCs *httptest.Server = nil

func testGetMockHs(t *testing.T, log *Log) (* HttpServer) {
	if testMockCs == nil {
		testMockCs = testCreateConsulServer("[{\"ServiceName\": \"test\", \"ServiceAddress\": \"192.168.0.1\", \"ServiceId\": \"1\"}]")
	}

	testMockCc := testCreateConsulContext(t, strings.Trim(testMockCs.URL, "http://"), log)

	testMockHs := &HttpServer{}
	s := &Server{}
	testMockHs.s = s
	s.cc = testMockCc

	return testMockHs
}

func testFreeMockHs() {
	testMockCs.Close()
	testMockCs = nil
}

func testCheckReturn(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 200 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func testCheckReturn_400(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 400 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func testCheckReturn_204(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 204 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func testCheckReturn_500(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 500 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func testCheckReturn_502(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 502 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}

func TestCreateHandlerBadrequest1(t *testing.T) {
	w, r := test_generate_rr("POST", "/dns/create", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerBadrequest2(t *testing.T) {
	w, r := test_generate_rr("GET", "/dns/create", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerBadrequest3(t *testing.T) {
	str := ""
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/create", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerBadrequest4(t *testing.T) {
	str := "{\"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/create", body)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerBadrequest5(t *testing.T) {
	str := "}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/create", body)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerBadrequest6(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/create", body)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	handler.hs.s.domain = "test.com"
	defer testFreeMockHs()

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestCreateHandlerOk(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/create", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &CreateHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	defer testFreeMockHs()

	handler.ServeHTTP(w, r)

	testCheckReturn(w, t)
}

func TestDeleteHandlerBadrequest1(t *testing.T) {

	w, r := test_generate_rr("POST", "/dns/delete", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerBadrequest2(t *testing.T) {

	w, r := test_generate_rr("GET", "/dns/delete", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerBadrequest3(t *testing.T) {
	str := ""
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/delete", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerBadrequest4(t *testing.T) {
	str := "{\"Address\":\"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/delete", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerBadrequest5(t *testing.T) {
	str := "}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/delete", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerBadrequest6(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\":\"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/delete", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	defer testFreeMockHs()

	handler.hs.s.domain = "test.com"

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestDeleteHandlerOk(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/delete", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &DeleteHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	defer testFreeMockHs()
	handler.ServeHTTP(w, r)

	testCheckReturn(w, t)
}

func TestReadHandlerBadrequest1(t *testing.T) {

	w, r := test_generate_rr("POST", "/dns/read", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestReadHandlerBadrequest2(t *testing.T) {

	w, r := test_generate_rr("GET", "/dns/read", nil)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestReadHandlerBadrequest3(t *testing.T) {
	str := "{\"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)

	w, r := test_generate_rr("POST", "/dns/read", body)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestReadHandlerBadrequest4(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)
	w, r := test_generate_rr("POST", "/dns/read", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	handler.hs.s.domain = "test.com"

	defer testFreeMockHs()

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestReadHandlerBadrequest5(t *testing.T) {
	str := "\"Name\": \"test\"}"
	body := strings.NewReader(str)
	w, r := test_generate_rr("POST", "/dns/read", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	testCheckReturn_400(w, t)
}

func TestReadHandlerOk(t *testing.T) {
	str := "{\"Name\": \"test\", \"Address\": \"192.168.0.1\"}"
	body := strings.NewReader(str)
	w, r := test_generate_rr("POST", "/dns/read", body)
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	handler := &ReadHandler{}
	handler.log = log

	handler.hs = testGetMockHs(t, log)
	defer testFreeMockHs()

	handler.ServeHTTP(w, r)

	testCheckReturn(w, t)
}

func TestReturnError1(t *testing.T) {
	w, _ := test_generate_rr("GET", "/dns/read", nil)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	returnError(w, BadConfigError, log)
}

func TestReturnError2(t *testing.T) {
	w, _ := test_generate_rr("GET", "/dns/read", nil)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	returnError(w, NoContentError, log)
}

func TestReturnError3(t *testing.T) {
	w, _ := test_generate_rr("GET", "/dns/read", nil)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	returnError(w, BadRequestError, log)
}

func TestReturnError4(t *testing.T) {
	w, _ := test_generate_rr("GET", "/dns/read", nil)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	returnError(w, InternalServerError, log)
}

func TestReturnError5(t *testing.T) {
	w, _ := test_generate_rr("GET", "/dns/read", nil)

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	returnError(w, BadGatewayError, log)
}
