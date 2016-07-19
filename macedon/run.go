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

func ShowHelp() {
	fmt.Printf("%s [-f config_file | -v | -h]\n", os.Args[0])
}

func ShowVersion() {
	fmt.Println("Version:", VERSION)
}

func ParseOption() bool {
	var c int
	OptErr = 0
	OptInd = 1
	for {
		if c = Getopt("f:hv"); c == EOF {
			break
		}
		switch c {
		case 'f':
			config_file = OptArg
			break
		case 'h':
			ShowHelp()
			return false
		case 'v':
			ShowVersion()
			return false
		default:
			fmt.Println("[Error] Invalid arguments")
			return false
		}
	}

	return true
}

var config_file string


func Run() {
	res := ParseOption()
	if !res {
		return
	}

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
		log.Error("Server run failed: ", err)
		time.Sleep(DEFAULT_QUIT_WAIT_TIME)
		return
	}
}
