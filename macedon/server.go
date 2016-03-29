package macedon

import (
	//"time"
	//"strings"
)

type Server struct {
	addr	string
	mysql	*MysqlClient
	http	*HttpServer
	log		*Log
}

func InitServer(conf *Config, log *Log) (*Server, error) {
	s := &Server{}

	s.log = log
	s.addr = conf.addr

	http, err := InitHttpServer(conf.addr, s.log)
	if err != nil {
		s.log.Error("Init http server failed")
		return nil, err
	}
	s.http = http
	http.s = s

	s.log.Debug("Init http server done")
	//time.Sleep(time.Second * 5)

	err = http.AddRouter(conf.location)
	if err != nil {
		s.log.Error("Server add router failed")
		return nil, err
	}

	mysql, err := InitMysqlClient(conf.maddr, conf.dbname, conf.dbuser, conf.dbpwd, s.log)
	if err != nil {
		s.log.Error("Init mysql client faild")
		return nil, err
	}

	mysql.s = s
	s.mysql = mysql

	return s, nil
}


func (s *Server) Run() error {
	err := s.http.Run()
	if err != nil {
		s.log.Error("Server run failed")
	}

	return nil
}

