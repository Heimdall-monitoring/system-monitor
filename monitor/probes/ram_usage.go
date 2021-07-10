package probes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
)

// RAMStats represent the stats for RAM usage
type RAMStats struct {
	Available int64 `json:"available"`
	Used      int64 `json:"used"`
	Free      int64 `json:"free"`
	Shared    int64 `json:"shared"`
}

// getRAMUsage gets current details on system ram usage
func getRAMUsage(runner commandRunner) (RAMStats, error) {
	ramCommand := []string{"/usr/bin/vmstat", "-s"}
	commandResult := runner.runCommand(ramCommand)

	if commandResult.StatusCode != 0 {
		return RAMStats{}, fmt.Errorf("An error occured during vmstat command execution: %q", commandResult.Stderr)
	}
	lines := strings.Split(commandResult.Stdout, "\n")

	available, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[0]), " ")[0], 10, 64)
	if err != nil {
		return RAMStats{}, fmt.Errorf("Error parsing available RAM: %w", err)
	}

	used, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[1]), " ")[0], 10, 64)
	if err != nil {
		return RAMStats{}, fmt.Errorf("Error parsing used RAM: %w", err)
	}

	free, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[4]), " ")[0], 10, 64)
	if err != nil {
		return RAMStats{}, fmt.Errorf("Error parsing free RAM: %w", err)
	}

	shared, err := strconv.ParseInt(strings.Split(strings.TrimSpace(lines[5]), " ")[0], 10, 64)
	if err != nil {
		return RAMStats{}, fmt.Errorf("Error parsing shared RAM: %w", err)
	}

	return RAMStats{available, used, free, shared}, nil
}

// RAMUsageDaemon gets the RAM usage for the probed machine
type RAMUsageDaemon struct {
	notifier         *utils.GoroutineNotifier
	runner           commandRunner
	refreshFrequency time.Duration
	logger           *log.Logger
	data             RAMStats
	status           bool
}

// NewRAMUsageDaemon returns a RAMUsage daemon for the application
func NewRAMUsageDaemon(logger *log.Logger, refreshFrequency time.Duration) *RAMUsageDaemon {
	return &RAMUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           LinuxCommandRunner{},
		refreshFrequency: refreshFrequency,
		logger:           logger,
		status:           false,
		data:             RAMStats{},
	}
}

// Name returns 'ram-usage'
func (d *RAMUsageDaemon) Name() string {
	return "ram-usage"
}

func (d *RAMUsageDaemon) mainRoutine() {
	d.logger.Debugf("Starting main goroutine for daemon %s", d.Name())
	for {
		select {
		case <-d.notifier.StopSignalChan():
			d.logger.Debugf("Stopping main goroutine for daemon %s", d.Name())
			defer d.notifier.ConfirmRoutineStopped()
			return
		case <-time.After(d.refreshFrequency):
			d.logger.Trace("Refreshing data")
			data, err := getRAMUsage(d.runner)
			if err != nil {
				d.logger.Errorf("Error getting RAM usage data: %s", err.Error())
				continue
			}
			d.data = data
		}
	}
}

// Start handles starting the RAM usage daemon
func (d *RAMUsageDaemon) Start() {
	d.logger.Debugf("Starting daemon %s", d.Name())
	go d.mainRoutine()
	d.status = true
}

// Stop handles stopping the RAM usage daemon
func (d *RAMUsageDaemon) Stop() {
	d.logger.Debugf("Stopping daemon %s", d.Name())
	d.notifier.StopRoutine()
	<-d.notifier.RoutineStoppedChan()
	d.status = false
}

// Status returns true while the daemon is running
func (d *RAMUsageDaemon) Status() bool {
	return d.status
}

// Export returns the currentmy stored data in the daemon as a JSON formattable interface
func (d *RAMUsageDaemon) Export() interface{} {
	return d.data
}
