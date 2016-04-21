package macedon

import (
	"time"
	"io/ioutil"
	//"strings"
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

	s.log = log
	s.addr = conf.addr

	hs, err := InitHttpServer(conf.addr, s.log)
	if err != nil {
		s.log.Error("Init http server failed")
		return nil, err
	}
	s.hs = hs
	hs.s = s

	s.log.Debug("Init http server done")

	err = hs.AddRouter(conf.location)
	if err != nil {
		s.log.Error("Server add router failed")
		return nil, err
	}

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

func (s *Server) HttpServer() (*HttpServer) {
	return s.hs
}

func (s *Server) PurgeContext() (*PurgeContext) {
	return s.pc
}
func (s *Server) SshContext() (*SshContext) {
	return s.sc
}

func (s *Server) Run() error {
	err := s.hs.Run()
	if err != nil {
		s.log.Error("Server run failed: ", err)
		return err
	}

	return nil
}

