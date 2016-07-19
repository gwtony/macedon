package goblin

import (
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/test"
)

func TestInitContext(t *testing.T) {
	conf := &config.Config{}
	log := test.TestInitlog()

	hs, _ := hserver.InitHttpServer(":12345", log)

	InitContext(conf, hs, log)
	t.Log("goblin init context done")
}
