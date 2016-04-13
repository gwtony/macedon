package macedon

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	goconf "github.com/msbranco/goconfig"
)

const DEFAULT_CONF_PATH         = "../conf"
const DEFAULT_CONF              = "macedon_new.conf"

const DEFAULT_CREATE_LOCATION   = "/create"
const DEFAULT_DELETE_LOCATION   = "/delete"
const DEFAULT_UPDATE_LOCATION   = "/update"
const DEFAULT_READ_LOCATION     = "/read"
const DEFAULT_NOTIFY_LOCATION   = "/notify"

const DEFAULT_SSH_TIMEOUT       = 5


type Config struct {
	addr       string  /* server bind address */

	location   string  /* handler location */

	maddr      string  /* mysql addr */
	dbname     string  /* db name */
	dbuser     string  /* db username */
	dbpwd      string  /* db password */

	sport      string  /* ssh port */
	suser      string  /* ssh user */
	skey       string  /* ssh key */
	sto        int64   /* ssh timeout */

	ips        string  /* ip to purge */
	cmd        string  /* purge command */
	purgable   int     /* do purge or not */

	dns_server string  /* dns server address */
	updatable  int     /* do dns update or not */

	log        string  /* log file */
	level      string  /* log level */
}

func (conf *Config) ReadConf(file string) (*Config, error) {
	if file == "" {
		file = filepath.Join(DEFAULT_CONF_PATH, DEFAULT_CONF)
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf file %s failed", file)
		return nil, err
	}

	//TODO: check
	conf.addr, err = c.GetString("default", "addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No addr")
		return nil, err
	}
	conf.location, err = c.GetString("default", "location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No addr")
		return nil, err
	}

	conf.log, err = c.GetString("default", "log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] log not found, use default log file")
		conf.log = ""
	}
	conf.level, err = c.GetString("default", "level")
	if err != nil {
		conf.level = "error"
		fmt.Fprintln(os.Stderr, "[Info] level not found, use default log level error")
	}

	conf.maddr, err = c.GetString("default", "mysql_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No mysql_addr")
		return nil, err
	}
	conf.dbname, err = c.GetString("default", "mysql_dbname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No mysql_dbname")
		return nil, err
	}
	conf.dbuser, err = c.GetString("default", "mysql_dbuser")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No mysql_dbuser")
		return nil, err
	}
	conf.dbpwd, err = c.GetString("default", "mysql_dbpwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No mysql_dbpwd")
		return nil, err
	}

	conf.sport, err = c.GetString("default", "ssh_port")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No ssh_port")
		return nil, err
	}
	conf.suser, err = c.GetString("default", "ssh_user")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No ssh_user")
		return nil, err
	}
	conf.skey, err = c.GetString("default", "ssh_key")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: No ssh_key")
		return nil, err
	}
	conf.sto, err = c.GetInt64("default", "ssh_timeout")
	if err != nil {
		conf.sto = DEFAULT_SSH_TIMEOUT
		fmt.Fprintln(os.Stderr, "[Info] ssh_timeout not found, use default timeout")
	}

	conf.purgable = 1
	conf.ips, err = c.GetString("default", "purge_ips")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] purge_ips not found, do not purge")
		conf.purgable = 0
	}
	conf.cmd, err = c.GetString("default", "purge_cmd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] purge_cmd not found, do not purge")
		conf.purgable = 0
	}

	conf.updatable = 1
	conf.dns_server, err = c.GetString("default", "dns_server")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] dns_server not found, do not update")
		conf.updatable = 0
	}

	return conf, nil
}

