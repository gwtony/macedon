package goblin

import (
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/config"
)

func InitContext(conf *config.Config, hs *hserver.HttpServer, log log.Log) error {
	cf := &GoblinConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("goblin parse config failed")
		return err
	}

	h := InitHandler(cf.loc, log)
	mc := InitMysqlContext(cf.maddr, cf.dbname, cf.dbuser, cf.dbpwd, log)

	api_loc := cf.api_loc

	hs.AddRouter(api_loc + GOBLIN_RULE_ADD_LOCATION, &AddHandler{h: h, mc: mc, log: log})
	hs.AddRouter(api_loc + GOBLIN_RULE_DELETE_LOCATION, &DeleteHandler{h: h, mc: mc, log: log})
	hs.AddRouter(api_loc + GOBLIN_RULE_READ_LOCATION, &ReadHandler{mc: mc, log: log})
	hs.AddRouter(api_loc + GOBLIN_SERVER_ADD_LOCATION, &AddServerHandler{mc: mc, log: log})
	hs.AddRouter(api_loc + GOBLIN_SERVER_DELETE_LOCATION, &DeleteServerHandler{mc: mc, log: log})
	hs.AddRouter(api_loc + GOBLIN_SERVER_READ_LOCATION, &ReadServerHandler{mc: mc, log: log})

	return nil
}
