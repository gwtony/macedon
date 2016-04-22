package macedon

import (
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
