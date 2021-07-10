package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/aHugues/system-monitor/monitor/utils"
	"github.com/aHugues/system-monitor/monitor/webserver"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the JSON configuration file")
	debug := flag.Bool("debug", false, "Start application in debug mode")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	config, err := utils.ReadConfigJSON(*configPath)
	if err != nil {
		log.Warn("Impossible to read from JSON configuration, using default config")
		config = utils.NewConfig()
	}
	webserver.RunServer(config)
	log.Info("Stopping server")
}
