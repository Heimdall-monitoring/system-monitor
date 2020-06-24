package utils

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/cihub/seelog"
)

// ServerConfig handles the configuration for the web service
type ServerConfig struct {
	ListenMode string `json:"listen-mode"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Socket     string `json:"socket"`
}

// LogConfig handles the configuration for the logger
type LogConfig struct {
	Level string `json:"level"`
}

// ProbesConfig handles the configuration of system probes
type ProbesConfig struct {
	DiskUsage  bool `json:"disk-usage"`
	RAMUsage   bool `json:"ram-usage"`
	SystemInfo bool `json:"system-info"`
	Systemd    bool `json:"systemd"`
}

// FullConfiguration handles the entire configuration of the server
type FullConfiguration struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
	Probes ProbesConfig `json:"probes"`
}

// NewConfig creates a new configuration with default values
func NewConfig() FullConfiguration {
	return FullConfiguration{
		Server: ServerConfig{
			ListenMode: "port",
			Host:       "127.0.0.1",
			Port:       5000,
			Socket:     "",
		},
		Log: LogConfig{
			Level: "INFO",
		},
		Probes: ProbesConfig{
			DiskUsage:  true,
			RAMUsage:   true,
			SystemInfo: true,
			Systemd:    true,
		},
	}
}

// ReadConfigJSON reads a JSON config file and returns the parsed configuration
func ReadConfigJSON(configPath string) (FullConfiguration, error) {
	log.Debugf("Reading configuration from %q", configPath)
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Error reading configuration file: %q", err)
		return FullConfiguration{}, err
	}
	parsedConfig := NewConfig()
	parsingError := json.Unmarshal(conf, &parsedConfig)
	if parsingError != nil {
		log.Errorf("Error parsing JSON configuration: %q", parsingError)
	}
	return parsedConfig, parsingError
}
