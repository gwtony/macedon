package macedon

import (
	"fmt"
	"time"
)

func Run() {
	cconf := new(Config)
	conf, err:= cconf.ReadConf("")
	//conf, _:= cconf.ReadConf("../conf/macedon.conf")
	if err != nil {
		time.Sleep(time.Millisecond * 200)
		return
	}

	if conf == nil {
        fmt.Println("No conf")
        return
    }

    log := GetLogger(conf.log, conf.level)

    server, err := InitServer(conf, log)
    if err != nil {
        log.Error("Init server failed")
		time.Sleep(time.Second)
        return
    }

    server.Run()
}
