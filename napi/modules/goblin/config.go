package goblin

import (
	"fmt"
	"os"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type GoblinConfig struct {
	maddr   string /* mysql addr */
	dbname  string /* db name */
	dbuser  string /* db username */
	dbpwd   string /* db password */

	api_loc string /* goblin api location */
	loc     string /* goblin location */
}


func (conf *GoblinConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	conf.maddr, err = cf.C.GetString("goblin", "mysql_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_addr")
		return err
	}
	conf.dbname, err = cf.C.GetString("goblin", "mysql_dbname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbname")
		return err
	}
	conf.dbuser, err = cf.C.GetString("goblin", "mysql_dbuser")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbuser")
		return err
	}
	conf.dbpwd, err = cf.C.GetString("goblin", "mysql_dbpwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbpwd")
		return err
	}

	conf.loc, err = cf.C.GetString("goblin", "location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Goblin] Read conf: No goblin_location, use default", GOBLIN_LOCATION)
		conf.loc = GOBLIN_LOCATION
	}

	conf.api_loc, err = cf.C.GetString("goblin", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Goblin] Read conf: No api_location, use default", GOBLIN_API_LOCATION)
		conf.api_loc = GOBLIN_API_LOCATION
	}

	return nil
}
