package macedon

import (
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	s		*Server

	addr	string
	dbname	string
	dbuser	string
	dbpwd	string

	log		*Log
}

func InitMysqlClient(addr, dbname, dbuser, dbpwd string, log *Log) (*MysqlClient, error) {
	mc := &MysqlClient{}

	mc.addr		= addr
	mc.dbname	= dbname
	mc.dbuser	= dbuser
	mc.dbpwd	= dbpwd
	mc.log		= log

	return mc, nil
}


