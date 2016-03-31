package macedon

import (
	"strings"
	"errors"
)

type PurgeContext struct {
	ips     []string
	iplen   int
	port    string
	cmd     string
	log     *Log
}

func InitPurgeContext(ips string, port string, cmd string, log *Log) (*PurgeContext, error) {
	pc := &PurgeContext{}

	pc.log = log

	pc.ips = strings.Split(ips, ",")
	pc.iplen = len(pc.ips)
	if pc.iplen <= 0 {
		pc.log.Error("Parse ips failed")
		return nil, errors.New("Parse ips failed")
	}

	pc.port = port
	pc.cmd = cmd

	return pc, nil
}

func (pc *PurgeContext) DoPurge(sc *SshContext) (error){
	pc.log.Debug("Do purge")

	ch := make(chan int, pc.iplen)

	for _, host := range pc.ips {
		pc.log.Debug("Purge ip: %s", host)

		go func(ip string) {
			defer func () { ch <- 1 }()

			sconn, err := sc.InitSshConn(ip + ":" + pc.port)
			if err != nil {
				pc.log.Error("Init ssh to %s failed", ip)
				return
			}

			defer sconn.SshClose()

			/* Do not care purge result ? */
			_, err = sconn.SshExec(pc.cmd)
			if err != nil {
				pc.log.Error("Execute %s in %s failed", pc.cmd, ip)
				return
			}
		}(host)
	}

	for i := 0; i < pc.iplen; i++ {
		<-ch
	}

	pc.log.Debug("All purge done")

	return nil
}
