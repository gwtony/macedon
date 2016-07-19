package run

import (
	"os"
	"fmt"
	"time"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/utils"
	"git.lianjia.com/lianjia-sysop/napi/server"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/variable"
)

func show_help() {
	fmt.Println(os.Args[0], "[-f config_file | -v | -h]")
}

func show_version() {
	fmt.Println("Version:", variable.VERSION)
}

func parse_option() {
	var c int
	utils.OptErr = 0
	for {
		if c = utils.Getopt("f:hv"); c == utils.EOF {
			break
		}
		switch c {
		case 'f':
			config_file = utils.OptArg
			break
		case 'h':
			show_help()
			os.Exit(0)
			break
		case 'v':
			show_version()
			os.Exit(0)
			break
		default:
			fmt.Printf("[Error] Invalid arguments %c\n", c)
			os.Exit(0)
			break
		}
	}
}

var config_file string


func Run() {
	parse_option()

	conf := new(config.Config)
	err := conf.ReadConf(config_file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read config file failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return
	}
	err = conf.ParseConf()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Parse config file failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return
	}

    rlog := log.GetLogger(conf.Log, conf.Level)
	if rlog == nil {
		fmt.Fprintln(os.Stderr, "[Error] Init log failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return
	}

    server, err := server.InitServer(conf, rlog)
    if err != nil {
        rlog.Error("[Error] Init server failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
        return
    }

    err = server.Run()
	if err != nil {
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return
	}
}
