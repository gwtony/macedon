package macedon

import (
	"time"
	//"strings"
)

type Server struct {
	addr    string

	hs      *HttpServer

	pc      *PurgeContext
	mc      *MysqlContext
	sc      *SshContext

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

	mc, err := InitMysqlContext(conf.maddr, conf.dbname, conf.dbuser, conf.dbpwd, s.log)
	if err != nil {
		s.log.Error("Init mysql client faild")
		return nil, err
	}
	s.mc = mc

	if conf.purgable == 1 {
		pc, err := InitPurgeContext(conf.ips, conf.sport, conf.cmd, s.log)
		if err != nil {
			s.log.Error("Init purge context failed")
		}
		s.pc = pc

		sc, err := InitSshContext(conf.skey, conf.suser, time.Duration(conf.sto) * time.Second, s.log)
		if err != nil {
			s.log.Error("Init ssh context failed")
			return nil, err
		}
		s.sc = sc
	} else {
		s.pc = nil
		s.sc = nil
	}

	return s, nil
}

func (s *Server) MysqlContext() (*MysqlContext) {
	return s.mc
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
		s.log.Error("Server run failed")
	}

	return nil
}

