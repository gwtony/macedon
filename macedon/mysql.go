package macedon

import (
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

type MysqlContext struct {
	addr    string
	dbname  string
	dbuser  string
	dbpwd   string

	log     *Log
}

const DNS_CREATE_SQL "insert into records (name, type, domain_id, ttl, content) values ('%s', '%s', %d, %d, '%s');"
const DNS_DELETE_SQL "delete from records where name = '%s' AND type = '%s';"
const DNS_DELETE_SQL_CONTENT "delete from records where name = '%s' and type = '%s' and content = '%s';"
const DNS_UPDATE_SQL "update records set disabled = %d where name = '%s' and type = '%s' and content = '%s';"
const DNS_READ_SQL "select domain_id, name, type, content, ttl, disabled from records where name = '%s' AND type = '%s';"


func InitMysqlContext(addr, dbname, dbuser, dbpwd string, log *Log) (*MysqlContext, error) {
	mc := &MysqlContext{}

	mc.addr     = addr
	mc.dbname   = dbname
	mc.dbuser   = dbuser
	mc.dbpwd    = dbpwd
	mc.log      = log

	return mc, nil
}

func (mc *MysqlContext) CreateSql(name, type_s, content string, domain_id, ttl int) (string, error) {
	if name == "" ||
		type_s = "" ||
		content = "" ||
		domain_id <= 0 ||
		ttl <= 0 {
		return "", errors.New("Create sql arguments invalid")
	}
	return fmt.Sprintf(DNS_CREATE_SQL, name, type_s, domain_id, ttl, content), nil
}

func (mc *MysqlContext) DeleteSql(name, type_s, content) (string, error) {
	if name == "" || type_s == "" {
		return "", errors.New("Delete sql arguments invalid")
	}

	if content == "" {
		return fmt.Sprintf(DNS_DELETE_SQL, name, type_s), nil
	}

	return fmt.Sprintf(DNS_DELETE_SQL_CONTENT, name, type_s, content), nil
}

func (mc *MysqlContext) UpdateSql(name, type_s, content string, disabled int) (string, error) {
	if name == "" ||
		type_s == "" ||
		content == "" ||
		(disables != 0 && disabled != 1) {
		return "", errors.New("Update sql arguments invalid")
	}

	return fmt.Sprintf(DNS_UPDATE_SQL, disabled, name, type_s, content), nil
}

func (mc *MysqlContext) ReadSql(name, type_s string) (string, error) {
	if name == "" || type_s == "" {
		return "", errors.New("Read sql arguments invalid")
	}

	return fmt.Sprintf(DNS_READ_SQL, name, type_s), nil
}
