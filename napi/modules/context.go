package modules
import (
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/modules/goblin"
	"git.lianjia.com/lianjia-sysop/napi/modules/sample"
	"git.lianjia.com/lianjia-sysop/napi/modules/macedon"
)

func InitModules(conf *config.Config, hs *hserver.HttpServer, log log.Log) {
	if err := goblin.InitContext(conf, hs, log); err != nil {
		log.Error("goblin module will not start")
	}

	if err := sample.InitContext(conf, hs, log); err != nil {
		log.Error("sample module will not start")
	}
	if err := macedon.InitContext(conf, hs, log); err != nil {
		log.Error("macedon module will not start")
	}
}
