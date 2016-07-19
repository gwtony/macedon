package macedon

import (
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/config"
)

func InitContext(conf *config.Config, hs *hserver.HttpServer, log log.Log) error {
	cf := &MacedonConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Macedon parse config failed")
		return err
	}

	h := InitHandler(cf.eaddr, cf.loc, log)

	api_loc := cf.api_loc
	domain := cf.domain

	hs.AddRouter(api_loc + MACEDON_ADD_LOC, &AddHandler{h: h, domain: domain, log: log})
	hs.AddRouter(api_loc + MACEDON_DELETE_LOC, &DeleteHandler{h: h, domain: domain, log: log})
	hs.AddRouter(api_loc + MACEDON_READ_LOC, &ReadHandler{domain: domain, log: log})
	hs.AddRouter(api_loc + MACEDON_SERVER_LOC, &AddServerHandler{log: log})
	hs.AddRouter(api_loc + MACEDON_SERVER_DELETE_LOC, &DeleteServerHandler{log: log})
	hs.AddRouter(api_loc + MACEDON_SERVER_READ_LOC, &ReadServerHandler{log: log})

	return nil
}
