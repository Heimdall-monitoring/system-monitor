package probes

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// RAMStats represent the stats for RAM usage
type RAMStats struct {
	Available int64
	Used      int64
	Free      int64
	Shared    int64
}

// GetRAMUsage gets current details on system ram usage
func GetRAMUsage() RAMStats {
	cmd := exec.Command("/usr/bin/vmstat", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(out.String(), "\n")
	available, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[0]), " ")[0], 10, 64)
	used, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[1]), " ")[0], 10, 64)
	free, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[4]), " ")[0], 10, 64)
	shared, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[5]), " ")[0], 10, 64)
	return RAMStats{available, used, free, shared}
}
