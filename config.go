package macedon

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
)

// MacedonConfig  Macedon config
type MacedonConfig struct {
	eaddr      []string /* etcd addr */

	apiLoc     string   /* macedon api location */
	loc        string   /* macedon location */

	domain     string

	purgeCmd   string
	purgeTo    time.Duration

	token      string   /* access token */
}

// ParseConfig parses config
func (conf *MacedonConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	eaddrStr, err := cf.C.GetString("macedon", "etcd_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [macedon] Read conf: No etcd_addr")
		return err
	}
	if eaddrStr == "" {
		fmt.Fprintln(os.Stderr, "[Error] [macedon] Empty etcd server address")
		return errors.BadConfigError
	}
	eaddr := strings.Split(eaddrStr, ",")
	for i := 0; i < len(eaddr); i++ {
		if eaddr[i] != "" {
			if !strings.Contains(eaddr[i], ":") {
				conf.eaddr = append(conf.eaddr, eaddr[i] + ":" + DEFAULT_ETCD_PORT)
			} else {
				conf.eaddr = append(conf.eaddr, eaddr[i])
			}
		}
	}

	conf.loc, err = cf.C.GetString("macedon", "location")
	if err != nil {
		conf.loc = DEFAULT_SKYDNS_LOC
	}

	conf.apiLoc, err = cf.C.GetString("macedon", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No api_location, use default location", MACEDON_LOC)
		conf.apiLoc = MACEDON_LOC
	}

	conf.domain, err = cf.C.GetString("macedon", "domain")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No domain")
		return err
	}

	conf.purgeCmd, err = cf.C.GetString("macedon", "purge_cmd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: use default purge_cmd")
		conf.purgeCmd = DEFAULT_PURGE_CMD
	}

	purgeTo, err := cf.C.GetInt64("macedon", "purge_timeout")
	if err != nil || purgeTo <= 0 {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: use default purge_timeout: ", DEFAULT_PURGE_TIMEOUT)
		purgeTo = DEFAULT_PURGE_TIMEOUT
	}
	conf.purgeTo =  time.Duration(purgeTo) * time.Second

	conf.token, err = cf.C.GetString("macedon", "token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No token")
		conf.token = ""
	}

	return nil
}
