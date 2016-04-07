package macedon

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type MysqlContext struct {
	addr    string
	dbname  string
	dbuser  string
	dbpwd   string

	login   string

	log     *Log
}

const DNS_CREATE_SQL          = "insert into records (name, type, domain_id, ttl, content) values (?, ?, ?, ?, ?)"
const DNS_DELETE_SQL          = "delete from records where name = ? and type = ?"
const DNS_DELETE_SQL_CONTENT  = "delete from records where name = ? and type = ? and content = ?"
const DNS_UPDATE_SQL          = "update records set disabled = ? where name = ? and type = ? and content = ?"
const DNS_READ_SQL            = "select domain_id, content, ttl, disabled from records where name = ? AND type = ?"
const DNS_NOTIFY_SQL          = "update domains set notified_serial = notified_serial + 1 where name = ?"


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

func (mc *MysqlContext) Open() (*sql.DB, error) {
	db, err := sql.Open("mysql", mc.login)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (mc *MysqlContext) Close(db *sql.DB) error{
	return db.Close()
}

func (mc *MysqlContext) QueryRead(db *sql.DB, name, type_s string) (*Response, error) {
	var content string
	var domain_id, ttl, disabled int
	ret := &Response{}
	flag := 0

	if name == "" || type_s == "" {
		mc.log.Error("Read sql arguments invalid")
		return nil, BadRequestError
	}

	rows, err := db.Query(DNS_READ_SQL, name, type_s)
	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, NoContentError
		}
		mc.log.Error("Execute read (name: %s, type: %s) failed: ", name, type_s, err)
		return nil, BadGatewayError
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&domain_id, &content, &ttl, &disabled)
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, NoContentError
		}
		if err != nil {
			mc.log.Error("Scan read answer failed")
			return nil, InternalServerError
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

	if flag == 0 {
		mc.log.Error("Scan no answer")
		return nil, NoContentError
	}

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed")
		return nil, InternalServerError
	}

	mc.log.Debug(ret)

	return ret, nil
}

func (mc *MysqlContext) QueryWrite(db *sql.DB, query string, args ...interface{}) (*Response, error) {
	ret := &Response{}

	res, err := db.Exec(query, args...)

	if err != nil {
		mc.log.Error("Execute write sql: ", query, args, " failed: ", err)
		return nil, BadGatewayError
	}
	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed")
		return nil, InternalServerError
	}

	ret.Result.Affected = int(affected)

	return ret, nil
}

func (mc *MysqlContext) QueryCreate(db *sql.DB, name, type_s, content string, domain_id, ttl int) (*Response, error) {
	if name == "" ||
		type_s == "" ||
		content == "" ||
		domain_id <= 0 ||
		ttl <= 0 {
		mc.log.Error("Create sql arguments invalid")
		return nil, BadRequestError
	}

	return mc.QueryWrite(db, DNS_CREATE_SQL, name, type_s, domain_id, ttl, content)
}

func (mc *MysqlContext) QueryDelete(db *sql.DB, name, type_s, content string) (*Response, error) {
	if name == "" || type_s == "" {
		mc.log.Error("Delete sql arguments invalid")
		return nil, BadRequestError
	}

	if content == "" {
		return mc.QueryWrite(db, DNS_DELETE_SQL, name, type_s)
	} else {
		return mc.QueryWrite(db, DNS_DELETE_SQL_CONTENT, name, type_s, content)
	}
}

func (mc *MysqlContext) QueryUpdate(db *sql.DB, name, type_s, content string, disabled int) (*Response, error) {
	if name == "" ||
		type_s == "" ||
		content == "" ||
		(disabled != 0 && disabled != 1) {
		mc.log.Error("Update sql arguments invalid")
		return nil, BadRequestError
	}

	return mc.QueryWrite(db, DNS_UPDATE_SQL, disabled, name, type_s, content)
}

func (mc *MysqlContext ) QueryNotify(db *sql.DB, name string) (*Response, error) {
	if name == "" {
		mc.log.Error("Read sql arguments invalid")
		return nil, BadRequestError
	}

	ret := &Response{}

	res, err := db.Exec(DNS_NOTIFY_SQL, name)

	if err != nil {
		mc.log.Error("Execute notify sql: ", DNS_NOTIFY_SQL, name, " failed: ", err)
		return nil, BadGatewayError
	}
	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed")
		return nil, InternalServerError
	}

	ret.Result.Affected = int(affected)

	return ret, nil
}
