package probes

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RAMStats represent the stats for RAM usage
type RAMStats struct {
	Available int64 `json:"available"`
	Used      int64 `json:"used"`
	Free      int64 `json:"free"`
	Shared    int64 `json:"shared"`
}

// GetRAMUsage gets current details on system ram usage
func GetRAMUsage(runner commandRunner) RAMStats {
	ramCommand := []string{"/usr/bin/vmstat", "-s"}
	commandResult := runner.runCommand(ramCommand)

	if commandResult.StatusCode != 0 {
		log.Errorf("An error occured during vmstat command execution: %q", commandResult.Stderr)
		return RAMStats{}
	}

	lines := strings.Split(commandResult.Stdout, "\n")

	available, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[0]), " ")[0], 10, 64)
	if err != nil {
		log.Errorf("Error parsing available RAM: %q", err)
		return RAMStats{}
	}

	used, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[1]), " ")[0], 10, 64)
	if err != nil {
		log.Errorf("Error parsing used RAM: %q", err)
		return RAMStats{}
	}

	free, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[4]), " ")[0], 10, 64)
	if err != nil {
		log.Errorf("Error parsing free RAM: %q", err)
		return RAMStats{}
	}

	shared, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[5]), " ")[0], 10, 64)
	if err != nil {
		log.Errorf("Error parsing shared RAM: %q", err)
		return RAMStats{}
	}

	return RAMStats{available, used, free, shared}
}
