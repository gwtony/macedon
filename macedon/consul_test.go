package macedon

import (
	"os"
	"testing"
)

func Test_Consul_RegisterService_Failed(t *testing.T) {

	log := test_log_init()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer os.RemoveAll(tempDir)

	cc, err := InitConsulContext("127.0.0.1", DEFAULT_REGISTER_LOC, DEFAULT_DEREGISTER_LOC, DEFAULT_READ_LOC, log)
	if err != nil {
		t.Fatal("init consul context failed")
	}

	err = cc.RegisterService("test", "192.168.0.1")
	if err != nil && err == BadGatewayError {
		t.Log("register service error badgateway")
	}
}
