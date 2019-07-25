package main

import (
	"getbox/config"
	"getbox/pbx"
	"getbox/utils"
	"log"
)

func main() {
	config := config.Init(".")
	pbx.IsNrCoresSet(config)

	if config.GetBool("log_to_file") {
		logfile := utils.LogOutput(config)
		log.SetOutput(logfile)
		defer logfile.Close()
	}

	pbx.NewGetbox(config)
}
