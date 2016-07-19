package hserver

import (
	"git.lianjia.com/lianjia-sysop/napi/test"
	"git.lianjia.com/lianjia-sysop/napi/errors"
	"testing"
)

func TestInitHttpServer(t *testing.T) {
	log := test.TestInitlog()
	_, err := InitHttpServer(":80", log)
	if err != nil {
		t.Fatalf("init http server error")
	}
	t.Log("init http server done")
}

func TestReturnError1(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnError(w, errors.BadConfigError, log)
}

func TestReturnError2(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnError(w, errors.NoContentError, log)
}

func TestReturnError3(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnError(w, errors.BadRequestError, log)
}

func TestReturnError4(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnError(w, errors.InternalServerError, log)
}

func TestReturnError5(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnError(w, errors.BadGatewayError, log)
}

func TestReturnResponseOk(t *testing.T) {
	w, _ := test.Test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	ReturnResponse(w, 1, log)
	t.Log("Return Response done")
}
