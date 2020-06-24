package main

import (
	"flag"

	log "github.com/cihub/seelog"

	"github.com/aHugues/system-monitor/monitor/utils"
	"github.com/aHugues/system-monitor/monitor/webserver"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the JSON configuration file")
	flag.Parse()

	config, err := utils.ReadConfigJSON(*configPath)
	if err != nil {
		log.Warn("Impossible to read from JSON configuration, using default config")
		config = utils.NewConfig()
	}
	webserver.RunServer(config)
}
