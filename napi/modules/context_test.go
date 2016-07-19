package modules
import (
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/test"
)

func TestInitModules(t *testing.T) {
	conf := &config.Config{}
	log := test.TestInitlog()

	hs, _ := hserver.InitHttpServer(":12345", log)
	InitModules(conf, hs, log)
	t.Log("init modules done")
}
