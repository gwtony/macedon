package macedon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func test_generate_rr(method, uri string, body io.Reader) (*httptest.ResponseRecorder, *http.Request){
	r, _ := http.NewRequest(method, uri, body)
	w := httptest.NewRecorder()
	return w, r
}

func test_check_return(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 200 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func test_check_return_400(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 400 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func test_check_return_204(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 204 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func test_check_return_500(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 500 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}
func test_check_return_502(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 502 {
		t.Fatalf("handler run error, code is %d", w.Code)
	}
	t.Log("handler done")
}

func Test_CreateHandler_ServeHTTP_Method(t *testing.T) {

	w, r := test_generate_rr("POST", "/dns/create", nil)
	log := test_init_log()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer test_destroy_log()

	handler := &CreateHandler{}
	handler.log = log

	handler.ServeHTTP(w, r)

	test_check_return_400(w, t)
}
