package probes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
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
		return DeviceStat{}, fmt.Errorf("Line %q does not match a df line", dfLine)
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
func getUsageStats(runner commandRunner) ([]DeviceStat, error) {
	dfCommand := []string{"/usr/bin/df", "-h", "-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs"}
	commandResult := runner.runCommand(dfCommand)

	if commandResult.StatusCode != 0 {
		return []DeviceStat{}, fmt.Errorf("Error during df command execution: %s", commandResult.Stderr)
	}

	lines := strings.Split(commandResult.Stdout, "\n")
	result := []DeviceStat{}
	for _, line := range lines {
		device, err := DeviceStatFromDf(line)
		if err == nil {
			result = append(result, device)
		}
	}
	return result, nil
}

// DiskUsageDaemon gets the disk usage for the probed machine
type DiskUsageDaemon struct {
	notifier         *utils.GoroutineNotifier
	runner           commandRunner
	refreshFrequency time.Duration
	logger           *log.Logger
	data             []DeviceStat
	status           bool
}

// NewDiskUsageDaemon returns a DiskUsage daemon for the application
func NewDiskUsageDaemon(logger *log.Logger, refreshFrequency time.Duration) *DiskUsageDaemon {
	return &DiskUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           LinuxCommandRunner{},
		refreshFrequency: refreshFrequency,
		logger:           logger,
		status:           false,
		data:             []DeviceStat{},
	}
}

// Name returns 'disk-usage'
func (d *DiskUsageDaemon) Name() string {
	return "disk-usage"
}

func (d *DiskUsageDaemon) mainRoutine() {
	d.logger.Debugf("Starting main goroutine for daemon %s", d.Name())
	for {
		select {
		case <-d.notifier.StopSignalChan():
			d.logger.Debugf("Stopping main goroutine for daemon %s", d.Name())
			defer d.notifier.ConfirmRoutineStopped()
			return
		case <-time.After(d.refreshFrequency):
			d.logger.Trace("Refreshing data")
			data, err := getUsageStats(d.runner)
			if err != nil {
				d.logger.Errorf("Error getting disk usage stats: %s", err.Error())
				continue
			}
			d.data = data
		}
	}
}

// Start handles starting the disk usage daemon
func (d *DiskUsageDaemon) Start() {
	d.logger.Debugf("Starting daemon %s", d.Name())
	go d.mainRoutine()
	d.status = true
}

// Stop handles stopping the disk usage daemon
func (d *DiskUsageDaemon) Stop() {
	d.logger.Debugf("Stopping daemon %s", d.Name())
	d.notifier.StopRoutine()
	<-d.notifier.RoutineStoppedChan()
	d.status = false
}

// Status returns true while the daemon is running
func (d *DiskUsageDaemon) Status() bool {
	return d.status
}

// Export returns the currently stored data in the daemon as a JSON formattable interface
func (d *DiskUsageDaemon) Export() interface{} {
	return d.data
}
