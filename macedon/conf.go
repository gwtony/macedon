package macedon

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	goconf "github.com/msbranco/goconfig"
)

type Config struct {
	addr       string  /* server bind address */

	location   string  /* handler location */

	sport      string  /* ssh port */
	suser      string  /* ssh user */
	skey       string  /* ssh key */
	sto        int64   /* ssh timeout */

	ips        string  /* ip to purge */
	cmd        string  /* purge command */
	purgable   int     /* do purge or not */

	caddr      string  /* consul server address */
	reg_loc    string  /* register location */
	dereg_loc  string  /* deregister location */
	read_loc   string  /* read location */
	domain     string  /* consul domain name */

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
	conf.caddr, err = c.GetString("default", "consul_addrs")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: no consul_addrs")
		return nil, err
	}
	conf.reg_loc, err = c.GetString("default", "register_location")
	if err != nil {
		conf.reg_loc = DEFAULT_REGISTER_LOC
		fmt.Fprintln(os.Stderr, "[Info] register_location not found, use default")
	}
	conf.dereg_loc, err = c.GetString("default", "deregister_location")
	if err != nil {
		conf.dereg_loc = DEFAULT_DEREGISTER_LOC
		fmt.Fprintln(os.Stderr, "[Info] deregister_location not found, use default")
	}
	conf.read_loc, err = c.GetString("default", "read_location")
	if err != nil {
		conf.read_loc = DEFAULT_READ_LOC
		fmt.Fprintln(os.Stderr, "[Info] read_location not found, use default")
	}
	conf.domain, err = c.GetString("default", "domain")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf: no domain")
		return nil, err
	}
	return conf, nil
}

