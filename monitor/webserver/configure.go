package webserver

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/cihub/seelog"
)

// Server handles the configuration for the web service
type Server struct {
	ListenMode string `json:"listen-mode"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Socket     string `json:"socket"`
}

// Log handles the configuration for the logger
type Log struct {
	Level string `json:"level"`
}

// Configuration handles the entire configuration of the server
type Configuration struct {
	Server Server `json:"server"`
	Log    Log    `json:"log"`
}

// ReadConfigJSON reads a JSON config file and returns the parsed configuration
func ReadConfigJSON(configPath string) (Configuration, error) {
	log.Debugf("Reading configuration from %q", configPath)
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Error reading configuration file: %q", err)
		return Configuration{}, err
	}
	parsedConfig := Configuration{}
	parsingError := json.Unmarshal(conf, &parsedConfig)
	if parsingError != nil {
		log.Errorf("Error parsing JSON configuration: %q", parsingError)
	}
	return parsedConfig, parsingError
}
