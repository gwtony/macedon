package goblin

import (
	"fmt"
	"time"
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

type GoblinMysqlResponse GoblinRequest
type ServersResult []string

const (
	GOBLIN_GET_SERVER_SQL     = "select addr from servers where product = ? and state = 1"

	GOBLIN_READ_SQL           = "select ip, uid, uuid, product, expire, action, rule_id from rules where ip = ? and uid = ? and uuid = ? and product = ? and expire > ?"
	GOBLIN_ADD_SQL            = "insert into rules (ip, uid, uuid, product, expire, action, rule_id) values (?, ?, ?, ?, ?, ?, ?)"
	GOBLIN_DELETE_SQL         = "delete from rules where ip = ? and uid = ? and uuid = ?"
	GOBLIN_UPDATE_DELETED_SQL = "update rules set deleted = 1 where rule_id = ?"
	GOBLIN_UPDATE_RESULT_SQL  = "update rules set result = 1 where rule_id = ?"

	GOBLIN_ADD_SERVER_SQL     = "insert into servers (addr, product, state) values (?, ?, 1)"
	GOBLIN_DELETE_SERVER_SQL  = "delete from servers where addr = ? and product = ?"
	GOBLIN_READ_SERVER_SQL    = "select addr, product from servers where addr = ?"
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
	rows, err := db.Query(GOBLIN_GET_SERVER_SQL, product)
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

func (mc *MysqlContext) QueryInsert(db *sql.DB, ip, uid, uuid, product string,
		expire int, action, ruleid string) error {
	return mc.QueryWrite(db, GOBLIN_ADD_SQL, ip, uid, uuid, product, expire, action, ruleid)
}

func (mc *MysqlContext) QueryRead(db *sql.DB, ip, uid, uuid, product string) (*GoblinMysqlResponse, error) {
	var expire int
	var action, ruleid string

	now := int(time.Now().Unix())
	rows, err := db.Query(GOBLIN_READ_SQL, ip, uid, uuid, product, now)
	if err != nil {
		if err == sql.ErrNoRows {
			mc.log.Error("Scan no answer")
			return nil, errors.NoContentError
		}
		mc.log.Error("Execute read (ip: %s, uid: %s, uuid: %s, product: %s) failed: %s", ip, uid, uuid, product, err)
		return nil, errors.BadGatewayError
	}

	defer rows.Close()

	gr := &GoblinMysqlResponse{}
	flag := 0
	for rows.Next() {
		err := rows.Scan(&ip, &uid, &uuid, &product, &expire, &action, &ruleid)
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

		gr.Ip = ip
		gr.Uid = uid
		gr.Uuid = uuid
		gr.Ruleid = ruleid
		gr.Product = product
		gr.Expire = expire
		gr.Action = action

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

	mc.log.Debug(gr)

	return gr, nil
}

func (mc *MysqlContext) QueryUpdateResult(db *sql.DB, ruleid string) error {
	if ruleid == "" {
		return errors.BadRequestError
	}

	res, err := db.Exec(GOBLIN_UPDATE_RESULT_SQL, ruleid)

	if err != nil {
		mc.log.Error("Execute update sql ", GOBLIN_UPDATE_RESULT_SQL, " for ruleid: ", ruleid, " failed: ", err)
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
	if ruleid == "" {
		return errors.BadRequestError
	}

	res, err := db.Exec(GOBLIN_UPDATE_DELETED_SQL, ruleid)

	if err != nil {
		mc.log.Error("Execute update sql %s for ruleid: %s failed: %s", GOBLIN_UPDATE_DELETED_SQL, ruleid, err)
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

	rows, err := db.Query(GOBLIN_READ_SERVER_SQL, addr)
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
	res, err := db.Exec(GOBLIN_ADD_SERVER_SQL, addr, product)

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
	res, err := db.Exec(GOBLIN_DELETE_SERVER_SQL, addr, product)

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

	mc.log.Error("Delete server failed")
	return errors.BadGatewayError
}
