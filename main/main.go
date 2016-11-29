package main

import (
	"fmt"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/macedon"
)

func main() {
	err := api.Init()
	if err != nil {
		fmt.Println("Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	err = macedon.InitContext(config, log)
	if err != nil {
		fmt.Println("Init Macedon failed")
	}

	api.Run()
}
