package probes

import (
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

// GetUptime return a string con
func GetUptime() (int64, error) {
	log.Debug("Reading uptime from /proc/uptime")
	dat, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		log.Errorf("Error reading uptime: %q", err)
		return 0, err
	}
	rawUptime, err := strconv.ParseFloat(strings.Split(string(dat), " ")[0], 64)
	if err != nil {
		log.Errorf("Impossible to parse uptime: %q", err)
		return 0, err
	}
	return int64(rawUptime), nil
}
