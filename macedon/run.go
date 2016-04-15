package macedon

import (
	"os"
	"fmt"
	"time"
	. "github.com/mattn/go-getopt"
)

const (
	DEFAULT_QUIT_WAIT_TIME = time.Millisecond * 200
)

func show_help() {
	fmt.Println("%s [-f config_file | -v | -h]", os.Args[0])
}

func show_version() {
	fmt.Println("Version: %s", VERSION)
}

func parse_option() {
	var c int
	OptErr = 0
	for {
		if c = Getopt("f:h:v"); c == EOF {
			break
		}
		switch c {
		case 'f':
			config_file = OptArg
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
			fmt.Println("[Error] Invalid arguments")
		}
	}
}

var config_file string


func Run() {
	parse_option()

	cconf := new(Config)
	conf, err:= cconf.ReadConf(config_file)
	if err != nil {
		time.Sleep(DEFAULT_QUIT_WAIT_TIME)
		return
	}

	if conf == nil {
        fmt.Println("[Error] No conf")
        return
    }

    log := GetLogger(conf.log, conf.level)

    server, err := InitServer(conf, log)
    if err != nil {
        log.Error("Init server failed")
		time.Sleep(DEFAULT_QUIT_WAIT_TIME)
        return
    }

    err = server.Run()
	if err != nil {
		time.Sleep(DEFAULT_QUIT_WAIT_TIME)
		return
	}
}
