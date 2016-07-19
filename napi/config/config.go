package config

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	goconf "github.com/msbranco/goconfig"
	"git.lianjia.com/lianjia-sysop/napi/errors"
	"git.lianjia.com/lianjia-sysop/napi/variable"
)

type Config struct {
	Addr       string  /* server bind address */

	Location   string  /* handler location */

	Log        string  /* log file */
	Level      string  /* log level */

	C          *goconf.ConfigFile /* goconfig struct */
}

func (conf *Config) ReadConf(file string) error {
	if file == "" {
		file = filepath.Join(variable.DEFAULT_CONFIG_PATH, variable.DEFAULT_CONFIG_FILE)
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf file %s failed", file)
		return err
	}
	conf.C = c
	return nil
}

func (conf *Config) ParseConf() error {
	var err error

	if conf.C == nil {
		fmt.Fprintln(os.Stderr, "[Error] Must read config first")
		return errors.BadConfigError
	}

	conf.Addr, err = conf.C.GetString("default", "addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Default] Read conf: No addr")
		return err
	}

	conf.Log, err = conf.C.GetString("default", "log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] Log not found, use default log file")
		conf.Log = ""
	}
	conf.Level, err = conf.C.GetString("default", "level")
	if err != nil {
		conf.Level = "error"
		fmt.Fprintln(os.Stderr, "[Info] [Default] Level not found, use default log level error")
	}

	return nil
}

