package sample

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type MysqlContext struct {
	addr    string
	dbname  string
	dbuser  string
	dbpwd   string

	login   string

	log     log.Log
}

type ServersResult []string

const (
	SAMPLE_CLEAR_EXPIRE_SQL   = "delete from rules where expire < ?"
	SAMPLE_GET_SERVER_SQL     = "select addr from servers where product = ? and state = 1"

	SAMPLE_READ_HOST_SQL      = "select type, host, band, action_type, action_value, rule_id, product, expire from rules where host = ? and expire > ? and deleted = 0"
	SAMPLE_READ_HOST_LIKE_SQL = "select type, host, band, action_type, action_value, rule_id, product, expire from rules where host like ? and expire > ? and deleted = 0"
	SAMPLE_READ_RULEID_SQL    = "select type, host, band, action_type, action_value, rule_id, product, expire from rules where rule_id = ? and expire > ?"
	SAMPLE_ADD_SQL            = "insert into rules (type, host, band, action_type, action_value, rule_id, product, expire) values (?, ?, ?, ?, ?, ?, ?, ?)"
	SAMPLE_UPDATE_DELETED_SQL = "update rules set deleted = 1 where rule_id = ?"
	SAMPLE_UPDATE_RESULT_SQL  = "update rules set result = 1 where rule_id = ?"

	SAMPLE_ADD_SERVER_SQL     = "insert into servers (addr, product, state) values (?, ?, 1)"
	SAMPLE_DELETE_SERVER_SQL  = "delete from servers where addr = ? and product = ?"
	SAMPLE_READ_SERVER_SQL    = "select addr, product from servers where addr = ?"
)

func InitMysqlContext(addr, dbname, dbuser, dbpwd string, log log.Log) *MysqlContext {
	mc := &MysqlContext{}

	mc.log      = log
	mc.addr     = addr
	mc.dbname   = dbname
	mc.dbuser   = dbuser
	mc.dbpwd    = dbpwd
	mc.login    = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbuser, dbpwd, addr, dbname)

	return mc
}

func (mc *MysqlContext) Open() (*sql.DB, error) {
	db, err := sql.Open("mysql", mc.login)
	if err != nil {
		mc.log.Error("open failed: %s", err)
		return nil, err
	}

	return db, nil
}

func (mc *MysqlContext) Close(db *sql.DB) error{
	return db.Close()
}

func (mc *MysqlContext) QueryGetServer(db *sql.DB, product string) (*ServersResult, error) {
	var addr string
	flag := 0
	sr := &ServersResult{}
	rows, err := db.Query(SAMPLE_GET_SERVER_SQL, product)
	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		mc.log.Error("Execute get server for product %s failed: %s", product, err)
		return nil, errors.BadGatewayError
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&addr)
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		if err != nil {
			mc.log.Error("Scan read answer failed: %s", err)
			return nil, errors.InternalServerError
		}
		if flag == 0 {
			flag = 1
		}
		*sr = append(*sr, addr)
	}

	if flag == 0 {
		mc.log.Error("Scan no answer")
		return nil, errors.NoContentError
	}

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed: %s", err)
		return nil, errors.InternalServerError
	}

	return sr, nil
}

func (mc *MysqlContext) QueryInsert(db *sql.DB, rtype, host, band, action_type, action_value, rule_id, product string, expire int) error {
	return mc.QueryWrite(db, SAMPLE_ADD_SQL, rtype, host, band, action_type, action_value, rule_id, product, expire)
}

func (mc *MysqlContext) QueryReadByRuleid(db *sql.DB, rule_id string) (*SampleRule, error) {
	var err error
	var expire, iband int
	var host, rtype, band, action_type, action_value, product string
	var resp *SampleRule

	now := int(time.Now().Unix())

	mc.log.Info("rule_id is ", rule_id)
	rows, err := db.Query(SAMPLE_READ_RULEID_SQL, rule_id, now)

	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		mc.log.Error("Execute read rule_id: %s failed: %s", rule_id, err)
		return nil, errors.BadGatewayError
	}

	defer rows.Close()

	flag := 0
	for rows.Next() {
		err := rows.Scan(&rtype, &host, &band, &action_type, &action_value, &rule_id, &product, &expire)
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		if err != nil {
			mc.log.Error("Scan read answer failed: %s", err)
			return nil, errors.InternalServerError
		}

		if flag == 0 {
			flag = 1
		}
		mc.log.Debug("sql result is: ", rtype, host, band, action_type, action_value, rule_id, product, expire)

		resp = &SampleRule{}
		resp.Type = rtype
		resp.Match.Host = strings.Split(host, ",")
		for _, v := range strings.Split(band, ",") {
			iband, _ = strconv.Atoi(v)
			resp.Match.Band = append(resp.Match.Band, iband)
		}
		resp.Match.Expire = expire
		resp.Action.Type = action_type
		resp.Action.Value = action_value
		resp.Product = product

		break
	}

	if flag == 0 {
		mc.log.Error("Scan no answer")
		return nil, errors.NoContentError
	}

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed: %s", err)
		return nil, errors.InternalServerError
	}

	mc.log.Debug("Query response is ", resp)

	return resp, nil
}

func (mc *MysqlContext) QueryReadByHost(db *sql.DB, host, htype string) (*SampleRuleResponse, error) {
	var err error
	var sqlstr string
	var expire, iband int
	var rule_id, rtype, band, action_type, action_value, product string
	var resp *SampleRule
	var resps SampleRuleResponse

	now := int(time.Now().Unix())

	if htype == EQUAL_HOST {
		sqlstr = SAMPLE_READ_HOST_SQL
	} else if htype == LIKE_HOST {
		sqlstr = SAMPLE_READ_HOST_LIKE_SQL
		host = "%" + host + "%"
	}

	mc.log.Debug("ReadByHost is", sqlstr, host)
	rows, err := db.Query(sqlstr, host, now)

	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		mc.log.Error("Execute read host: %s failed: %s", host, err)
		return nil, errors.BadGatewayError
	}

	defer rows.Close()

	flag := 0
	for rows.Next() {
		err := rows.Scan(&rtype, &host, &band, &action_type, &action_value, &rule_id, &product, &expire)
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		if err != nil {
			mc.log.Error("Scan read answer failed: %s", err)
			return nil, errors.InternalServerError
		}

		if flag == 0 {
			flag = 1
		}

		resp = &SampleRule{}
		resp.Type = rtype
		resp.Match.Host = strings.Split(host, ",")
		for _, v := range strings.Split(band, ",") {
			iband, _ = strconv.Atoi(v)
			resp.Match.Band = append(resp.Match.Band, iband)
		}
		resp.Match.Expire = expire
		resp.Action.Type = action_type
		resp.Action.Value = action_value
		resp.Product = product
		resps = append(resps, *resp)
	}

	if flag == 0 {
		mc.log.Error("Scan no answer")
		return nil, errors.NoContentError
	}

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed: %s", err)
		return nil, errors.InternalServerError
	}

	mc.log.Debug(resps)

	return &resps, nil
}

func (mc *MysqlContext) QueryUpdateResult(db *sql.DB, ruleid string) error {
	res, err := db.Exec(SAMPLE_UPDATE_RESULT_SQL, ruleid)

	if err != nil {
		mc.log.Error("Execute update sql ", SAMPLE_UPDATE_RESULT_SQL, " for ruleid: ", ruleid, " failed: ", err)
		return errors.BadGatewayError
	}

	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}

	if int(affected) > 0 {
		return nil
	}

	mc.log.Info("No such ruleid: ", ruleid)
	return errors.NoContentError
}

func (mc *MysqlContext) QueryUpdateDeleted(db *sql.DB, ruleid string) error {
	res, err := db.Exec(SAMPLE_UPDATE_DELETED_SQL, ruleid)

	if err != nil {
		mc.log.Error("Execute update sql %s for ruleid: %s failed: %s", SAMPLE_UPDATE_DELETED_SQL, ruleid, err)
		return errors.BadGatewayError
	}

	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}

	if int(affected) > 0 {
		return nil
	}

	mc.log.Info("No such ruleid: ", ruleid)
	return errors.NoContentError
}

func (mc *MysqlContext) QueryWrite(db *sql.DB, query string, args ...interface{}) error {
	res, err := db.Exec(query, args...)

	if err != nil {
		mc.log.Error("Execute write sql: ", query, args, " failed: ", err)
		return errors.BadGatewayError
	}
	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}
	if int(affected) <= 0 {
		return errors.BadGatewayError
	}

	return nil
}

func (mc *MysqlContext) QueryReadServer(db *sql.DB, addr string) ([]ServerResponse, error) {
	var product string
	flag := 0

	rows, err := db.Query(SAMPLE_READ_SERVER_SQL, addr)
	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		mc.log.Error("Execute get server for addr %s failed: %s", addr, err)
		return nil, errors.BadGatewayError
	}
	defer rows.Close()

	sr := []ServerResponse{}
	for rows.Next() {
		err := rows.Scan(&addr, &product)
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		if err != nil {
			mc.log.Error("Scan read answer failed: %s", err)
			return nil, errors.InternalServerError
		}
		if flag == 0 {
			flag = 1
		}
		tsr := &ServerResponse{}
		tsr.Addr = addr
		tsr.Product = product
		sr = append(sr, *tsr)
	}

	if flag == 0 {
		mc.log.Error("Scan no answer")
		return nil, errors.NoContentError
	}

	err = rows.Err()
	if err != nil {
		mc.log.Error("Iterate row failed: %s", err)
		return nil, errors.InternalServerError
	}

	return sr, nil
}

func (mc *MysqlContext) QueryAddServer(db *sql.DB, addr, product string) error {
	res, err := db.Exec(SAMPLE_ADD_SERVER_SQL, addr, product)

	if err != nil {
		mc.log.Error("Execute add server sql for %s, %s failed: %s", addr, product, err)
		return errors.BadGatewayError
	}

	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}

	if int(affected) > 0 {
		return nil
	}

	mc.log.Error("Add server failed")
	return errors.BadGatewayError
}

func (mc *MysqlContext) QueryDeleteServer(db *sql.DB, addr, product string) error {
	res, err := db.Exec(SAMPLE_DELETE_SERVER_SQL, addr, product)

	if err != nil {
		mc.log.Error("Execute delete server sql for addr: %s, product: %s failed: %s", addr, product, err)
		return errors.BadGatewayError
	}

	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}

	if int(affected) > 0 {
		return nil
	}

	mc.log.Error("Delete server no content")
	return errors.NoContentError
}

func (mc *MysqlContext) QueryClearExpire(db *sql.DB, expire int) error {
	res, err := db.Exec(SAMPLE_CLEAR_EXPIRE_SQL, expire)

	if err != nil {
		mc.log.Error("Execute clear expire rules before %d failed", expire)
		return errors.BadGatewayError
	}

	affected, err := res.RowsAffected()
	if err != nil {
		mc.log.Error("Get rows affected failed: %s", err)
		return errors.InternalServerError
	}

	if int(affected) > 0 {
		return nil
	}
	mc.log.Error("Clear expire no content")
	return errors.NoContentError
}
