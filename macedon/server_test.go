package macedon

import (
	"testing"
)

func Test_Server_1(t *testing.T) {
	conf := &Config{
				"172.30.19.33",
				"/dns",
				"22", "root", "./id_isa.pub", 10,
				"127.0.0.1", "echo", 0,
				"172.30.19.33", "", "", "", "test.com",
				"",
				"",
	}

	log := test_init_log()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer test_destroy_log()

	_, err := InitServer(conf, log)
	if err != nil {
		t.Error("Test server 1 failed!")
	}
	t.Log("Test server 1 ok.")
}
