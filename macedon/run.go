package macedon

import (
	"fmt"
)

func Run() {
    cconf := new(Config)
    conf, _:= cconf.ReadConf("../conf/macedon.conf")

	if conf == nil {
        fmt.Println("no conf")
        return
    }

    log := GetLogger(conf.log, conf.level)

    server, err := InitServer(conf, log)
    if err != nil {
        log.Error("init server failed")
        return
    }

    server.Run()
}
