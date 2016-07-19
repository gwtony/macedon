package macedon

import (
	"fmt"
	"os"
	"strings"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type MacedonConfig struct {
	eaddr_str string   /* etcd addr string */
	eaddr     []string /* etcd addr */

	api_loc   string   /* macedon api location */
	loc       string   /* macedon location */
}


func (conf *MacedonConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	conf.eaddr_str, err = cf.C.GetString("macedon", "etcd_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [macedon] Read conf: No etcd_addr")
		return err
	}
	if conf.eaddr_str == "" {
		fmt.Fprintln("[Error] [macedon] Empty etcd server address")
		return errors.BadConfigError
	}
	conf.eaddr = strings.Split(conf.eaddr_str, ",")
	for i := 0; i < len(conf.eaddr); i++ {
		if !strings.Contains(conf.eaddr[i], ":") {
			conf.eaddr[i] = conf.eaddr[i] + ":" + DEFAULT_ETCD_PORT
		}
	}

	conf.loc, err = cf.C.GetString("macedon", "location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No macedon_location, use default location", MACEDON_LOCATION)
		conf.loc = MACEDON_LOCATION
	}

	conf.api_loc, err = cf.C.GetString("macedon", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No api_location, use default location", MACEDON_API_LOCATION)
		conf.api_loc = MACEDON_API_LOCATION
	}

	conf.loc, err = cf.C.GetString("macedon", "domain")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No domain")
		return err
	}

	return nil
}
