package macedon

import (
	goconf "github.com/msbranco/goconfig"
	"path/filepath"
	//"time"
)

const DEFAULT_CONF_PATH         = "../conf"
const DEFAULT_CONF              = "macedon.conf"

const DEFAULT_CREATE_LOCATION   = "/create"
const DEFAULT_DELETE_LOCATION   = "/delete"
const DEFAULT_UPDATE_LOCATION   = "/update"
const DEFAULT_READ_LOCATION     = "/read"
const DEFAULT_NOTIFY_LOCATION   = "/notify"

const DEFAULT_SSH_TIMEOUT       = 5


type Config struct {
	addr       string  /* server bind address */

	location   string	/* handler location */

	maddr      string	/* mysql addr */
	dbname     string	/* db name */
	dbuser     string	/* db username */
	dbpwd      string	/* db password */

	sport      string  /* ssh port */
	suser      string  /* ssh user */
	skey       string  /* ssh key */
	sto        int64   /* ssh timeout */

	ips        string	/* ip to purge */
	cmd        string	/* purge command */

	log        string	/* log file */
	level      string	/* log level */
}

func (conf *Config) ReadConf(file string) (*Config, error) {
	if file == "" {
		file = filepath.Join(DEFAULT_CONF_PATH, DEFAULT_CONF)
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		return nil, err
	}

	//TODO: check
	conf.addr, _        = c.GetString("default", "addr")
	conf.location, _    = c.GetString("default", "location")

	conf.log, _         = c.GetString("default", "log")
	conf.level, _       = c.GetString("default", "level")

	conf.maddr, _       = c.GetString("default", "mysql_addr")
	conf.dbname, _      = c.GetString("default", "mysql_dbname")
	conf.dbuser, _      = c.GetString("default", "mysql_dbuser")
	conf.dbpwd, _       = c.GetString("default", "mysql_dbpwd")

	conf.sport, _       = c.GetString("default", "ssh_port")
	conf.suser, _       = c.GetString("default", "ssh_user")
	conf.skey, _        = c.GetString("default", "ssh_key")
	conf.sto, err       = c.GetInt64("default", "ssh_timeout")
	if err != nil {
		conf.sto = DEFAULT_SSH_TIMEOUT
	}

	conf.ips, _         = c.GetString("default", "purge_ips")
	conf.cmd, _         = c.GetString("default", "purge_cmd")

	return conf, nil
}

