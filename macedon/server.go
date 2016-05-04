package macedon

import (
	"os"
	"time"
	"os/signal"
	"syscall"
	"io/ioutil"
	"strings"
)

type Server struct {
	addr    string

	hs      *HttpServer

	pc      *PurgeContext
	cc      *ConsulContext
	sc      *SshContext
	domain  string

	log     *Log
}

func InitServer(conf *Config, log *Log) (*Server, error) {
	s := &Server{}

	if conf.addr == "" {
		log.Error("Empty server address")
		return nil, BadConfigError
	}
	addr := conf.addr
	if !strings.Contains(addr, ":") {
		addr = addr + ":" + DEFAULT_HTTPSERVER_PORT
	}
	s.log = log
	s.addr = addr

	hs := InitHttpServer(s.addr, s.log)
	s.hs = hs
	hs.s = s

	s.log.Debug("Init http server done")

	conf.location = strings.TrimSpace(conf.location)
	if conf.location == "/" {
		log.Error("Location invalid")
		return nil, BadConfigError
	}
	hs.AddRouter(conf.location)

	if conf.purgable == 1 {
		pc, err := InitPurgeContext(conf.ips, conf.sport, conf.cmd, s.log)
		if err != nil {
			s.log.Error("Init purge context failed")
			return nil, err
		}
		s.pc = pc

		key, err := ioutil.ReadFile(conf.skey)
		if err != nil {
			s.log.Error("Read private key from %s failed", conf.skey)
			return nil, err
		}
		sc, err := InitSshContext(string(key), conf.suser, time.Duration(conf.sto) * time.Second, s.log)
		if err != nil {
			s.log.Error("Init ssh context failed")
			return nil, err
		}
		s.sc = sc
	} else { /* Do not purge */
		s.pc = nil
		s.sc = nil
	}

	cc, err := InitConsulContext(conf.caddr, conf.reg_loc, conf.dereg_loc, conf.read_loc, s.log)
	if err != nil {
		s.log.Error("Init consul context failed")
		return nil, err
	}
	s.cc = cc
	s.domain = DEFAULT_SUB_ZONE + conf.domain

	return s, nil
}

func (s *Server)InitReload() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)

	for {
		<-ch

		rcconf := new(Config)
		rconf, err:= rcconf.ReadConf(config_file)
		if err != nil {
			s.log.Error("Reload purge ips read conf failed")
			continue
		}

		sips := strings.Split(rconf.ips, ",")
		ips := []string{}
		for _, ip := range sips {
			if ip != "" {
				ips = append(ips, ip)
			}
		}
		iplen := len(ips)
		if iplen <= 0 {
			s.log.Error("Reload purge ips parse ips failed")
			continue
		}

		s.pc.lock.Lock()
		s.pc.ips = ips
		s.pc.iplen = iplen
		s.pc.lock.Unlock()
		s.log.Debug("Reload done", s.pc.ips)
	}
}

func (s *Server) Run() error {
	go s.InitReload()

	err := s.hs.Run()
	if err != nil {
		s.log.Error("Server run failed: ", err)
		return err
	}

	return nil
}

