package server

import (
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/test"

)

func TestInitServerOK(t *testing.T) {
	conf := &config.Config{}
	conf.Addr = "localhost:80"
	log := test.TestInitlog()
	InitServer(conf, log)
	t.Log("init server done")
}
