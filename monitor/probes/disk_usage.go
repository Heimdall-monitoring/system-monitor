package probes

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

// DeviceStat represents the usage stats for a device
type DeviceStat struct {
	Filesystem string `json:"filesystem"`
	MountPoint string `json:"mountpoint"`
	Size       int64  `json:"size"`
	Used       int64  `json:"used"`
}

// ToString creates a string from a given DeviceStats
func (dev *DeviceStat) ToString() string {
	return fmt.Sprintf("Device %s - mountpoint %s\tTotal size %dGB - Used %dGB", dev.Filesystem, dev.MountPoint, dev.Size, dev.Used)
}

// DeviceStatFromDf builds the device stats from a DF output
func DeviceStatFromDf(dfLine string) (DeviceStat, error) {
	r := regexp.MustCompile(`^(?P<device_path>(/(\w+))+) +(?P<size>\d+)\D +(?P<used>\d+)\D.+(?P<mountpoint>(/(\w*))+)$`)
	match := r.FindStringSubmatch(dfLine)
	if r.MatchString(dfLine) == false {
		return DeviceStat{}, errors.New("No match")
	}
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	filesystem := result["device_path"]
	size, _ := strconv.ParseInt(result["size"], 10, 64)
	used, _ := strconv.ParseInt(result["used"], 10, 64)
	mountpoint := result["mountpoint"]
	return DeviceStat{filesystem, mountpoint, size, used}, nil
}

// GetUsageStats compute the disk usage on the machine and return it
func GetUsageStats(runner commandRunner) []DeviceStat {
	dfCommand := []string{"/usr/bin/df", "-h", "-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs"}
	commandResult := runner.runCommand(dfCommand)

	if commandResult.StatusCode != 0 {
		log.Errorf("An error occured during df command execution: %q", commandResult.Stderr)
		return []DeviceStat{}
	}

	lines := strings.Split(commandResult.Stdout, "\n")
	result := []DeviceStat{}
	for _, line := range lines {
		device, err := DeviceStatFromDf(line)
		if err == nil {
			result = append(result, device)
		}
	}
	return result
}
