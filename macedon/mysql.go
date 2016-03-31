package macedon

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

type MysqlContext struct {
	addr    string
	dbname  string
	dbuser  string
	dbpwd   string

	login	string

	log     *Log
}

const DNS_CREATE_SQL          = "insert into records (name, type, domain_id, ttl, content) values ('%s', '%s', %d, %d, '%s');"
const DNS_DELETE_SQL          = "delete from records where name = '%s' and type = '%s';"
const DNS_DELETE_SQL_CONTENT  = "delete from records where name = '%s' and type = '%s' and content = '%s';"
const DNS_UPDATE_SQL          = "update records set disabled = %d where name = '%s' and type = '%s' and content = '%s';"
const DNS_READ_SQL            = "select domain_id, content, ttl, disabled from records where name = '%s' AND type = '%s';"


func InitMysqlContext(addr, dbname, dbuser, dbpwd string, log *Log) (*MysqlContext, error) {
	mc := &MysqlContext{}

	mc.log      = log
	mc.addr     = addr
	mc.dbname   = dbname
	mc.dbuser   = dbuser
	mc.dbpwd    = dbpwd
	mc.login    = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbuser, dbpwd, addr, dbname)

	return mc, nil
}

func (mc *MysqlContext) CreateSql(name, type_s, content string, domain_id, ttl int) (string, error) {
	if name == "" ||
		type_s == "" ||
		content == "" ||
		domain_id <= 0 ||
		ttl <= 0 {
		return "", errors.New("Create sql arguments invalid")
	}
	return fmt.Sprintf(DNS_CREATE_SQL, name, type_s, domain_id, ttl, content), nil
}

func (mc *MysqlContext) DeleteSql(name, type_s, content string) (string, error) {
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
		(disabled != 0 && disabled != 1) {
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

func (mc *MysqlContext) opendb() (*sql.DB, error) {
	db, err := sql.Open("mysql", mc.login)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (mc *MysqlContext) closedb(db *sql.DB) error{
	return db.Close()
}

func (mc *MysqlContext) QueryRead(name, type_s string) (*Response, error) {
	var content string
	var domain_id, ttl, disabled int
	ret := &Response{}
	flag := 0

	db, err := mc.opendb()
	if err != nil {
		mc.log.Error("Open db failed")
		return nil, err
	}
	defer mc.closedb(db)

	query, err := mc.ReadSql(name, type_s)
	if err != nil {
		mc.log.Error("Generate read sql failed")
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		mc.log.Error("Execute read query (%s) failed", query)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&domain_id, &content, &ttl, &disabled)
		if err != nil {
			mc.log.Error("Scan read answer failed")
			//TODO: break or not
		}
		if flag == 0 {
			ret.Result.Affected = 1
			ret.Result.Data.Domain_id = domain_id
			ret.Result.Data.Name = name
			ret.Result.Data.Type = type_s
			ret.Result.Data.Ttl = ttl
			flag = 1
		}

		rec := &Record{}
		rec.Content = content
		rec.Disabled = disabled
		ret.Result.Data.Records = append(ret.Result.Data.Records, *rec)
    }

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed")
		return nil, err
	}

	mc.log.Debug(ret)

	return ret, nil
}
