package probes

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
)

// getUptime return a string con
func (d *UptimeDaemon) getUptime() (int64, error) {
	d.logger.Debug("Reading uptime from /proc/uptime")
	dat, err := ioutil.ReadFile(d.uptimeFile)
	if err != nil {
		return 0, fmt.Errorf("Impossible to read uptime: %w", err)
	}
	rawUptime, err := strconv.ParseFloat(strings.Split(string(dat), " ")[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Impossible to parse uptime value: %w", err)
	}
	return int64(rawUptime), nil
}

// UptimeDaemon gets the current machine uptime
type UptimeDaemon struct {
	notifier         *utils.GoroutineNotifier
	refreshFrequency time.Duration
	logger           *log.Logger
	data             int64
	status           bool
	uptimeFile       string
}

// NewUptimeDaemon returns an Uptime daemon for the application
func NewUptimeDaemon(logger *log.Logger, refreshFrequency time.Duration) *UptimeDaemon {
	return &UptimeDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		refreshFrequency: refreshFrequency,
		logger:           logger,
		data:             -1,
		status:           false,
		uptimeFile:       "/proc/uptime",
	}
}

// Name returns 'uptime'
func (d *UptimeDaemon) Name() string {
	return "uptime"
}

func (d *UptimeDaemon) mainRoutine() {
	d.logger.Debugf("Starting main goroutine for daemon %s", d.Name())
	for {
		select {
		case <-d.notifier.StopSignalChan():
			d.logger.Debugf("Stopping main goroutine for daemon %s", d.Name())
			defer d.notifier.ConfirmRoutineStopped()
			return
		case <-time.After(d.refreshFrequency):
			d.logger.Trace("Refreshing data")
			data, err := d.getUptime()
			if err != nil {
				d.logger.Errorf("Error refreshing data: %s", err.Error())
				continue
			}
			d.data = data
		}
	}
}

// Start handles starting the uptime daemon
func (d *UptimeDaemon) Start() {
	d.logger.Debugf("Starting daemon %s", d.Name())
	go d.mainRoutine()
	d.status = true
}

// Stop handles stopping the uptime daemon
func (d *UptimeDaemon) Stop() {
	d.logger.Debugf("Stopping daemon %s", d.Name())
	d.notifier.StopRoutine()
	<-d.notifier.RoutineStoppedChan()
	d.status = false
}

// Status returns true while the daemon is running
func (d *UptimeDaemon) Status() bool {
	return d.status
}

// Export returns the currently stored data in the daemon as a JSON formattable interface
func (d *UptimeDaemon) Export() interface{} {
	return d.data
}
