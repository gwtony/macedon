package macedon

import (
	//"fmt"
	"time"
	"testing"
)

func TestPurgeOk(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	sc, err := InitSshContext(test_key, "root", time.Second * 5, log)
	if err != nil {
		t.Fatal("init ssh context failed")
	}
	pc, err := InitPurgeContext("127.0.0.1", "22", "echo", log)
	pc.DoPurge(sc)

	t.Log("purge ok")
}

func TestPurgeEmptyIp(t *testing.T) {
	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}

	defer testDestroylog()

	_, err := InitPurgeContext("", "22", "echo", log)
	if err != nil {
		t.Log("purge empty ip done")
	} else {
		t.Fatal("purge check ip failed")
	}
}
