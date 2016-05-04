package macedon

import (
	"sync"
	"strings"
	"errors"
)

type PurgeContext struct {
	ips     []string
	lock    *sync.RWMutex
	iplen   int
	port    string
	cmd     string
	log     *Log
}

func InitPurgeContext(ips string, port string, cmd string, log *Log) (*PurgeContext, error) {
	pc := &PurgeContext{}

	pc.log = log

	sips := strings.Split(ips, ",")
	for _, ip := range sips {
		if ip != "" {
			pc.ips = append(pc.ips, ip)
		}
	}
	pc.iplen = len(pc.ips)
	if pc.iplen <= 0 {
		pc.log.Error("Parse ips failed")
		return nil, errors.New("Parse ips failed")
	}

	pc.lock = &sync.RWMutex{}

	pc.port = port
	pc.cmd = cmd

	return pc, nil
}

/* Do not return any error */
func (pc *PurgeContext) DoPurge(sc *SshContext, name string) error {
	pc.log.Debug("Do purge")

	ch := make(chan int, pc.iplen)

	pc.lock.RLock()
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
			_, err = sconn.SshExec(pc.cmd + " " + name)
			if err != nil {
				pc.log.Error("Execute %s in %s failed", pc.cmd + " " + name, ip)
				return
			}
		}(host)
	}
	pc.lock.RUnlock()

	for i := 0; i < pc.iplen; i++ {
		<-ch
	}

	pc.log.Debug("All purge done")

	return nil
}
