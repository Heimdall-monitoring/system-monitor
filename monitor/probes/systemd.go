package probes

import (
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
)

// getServiceStatus returns True if the probed status is running
func getServiceStatus(runner commandRunner, service string) bool {
	systemdCommand := []string{"/bin/systemctl", "is-active", "--quiet", service}
	commandResult := runner.runCommand(systemdCommand)
	return commandResult.StatusCode == 0
}

// etServicesStatuses computes the statuses for a list of services
func getServicesStatuses(runner commandRunner, services []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range services {
		result[s] = getServiceStatus(runner, s)
	}
	return result
}

// SystemdDaemon gets the status for the given systemd services
type SystemdDaemon struct {
	notifier         *utils.GoroutineNotifier
	runner           commandRunner
	refreshFrequency time.Duration
	logger           *log.Logger
	data             map[string]bool
	services         []string
	status           bool
}

// NewSystemdDaemon returns a SystemdDaemon for the given services
func NewSystemdDaemon(logger *log.Logger, refreshFrequency time.Duration, services []string) *SystemdDaemon {
	return &SystemdDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           LinuxCommandRunner{},
		refreshFrequency: refreshFrequency,
		logger:           logger,
		data:             make(map[string]bool),
		services:         services,
		status:           false,
	}

}

// Name returns 'systemd'
func (d *SystemdDaemon) Name() string {
	return "systemd"
}

func (d *SystemdDaemon) mainRoutine() {
	d.logger.Debugf("Starting main goroutine for daemon %s", d.Name())
	for {
		select {
		case <-d.notifier.StopSignalChan():
			d.logger.Debugf("Stopping main goroutine for daemon %s", d.Name())
			defer d.notifier.ConfirmRoutineStopped()
			return
		case <-time.After(d.refreshFrequency):
			d.logger.Trace("Refreshing data")
			d.data = getServicesStatuses(d.runner, d.services)
		}
	}
}

// Start handles starting the systemd daemon
func (d *SystemdDaemon) Start() {
	d.logger.Debugf("Starting daemon %s", d.Name())
	go d.mainRoutine()
	d.status = true
}

// Stop handles stopping the systemd daemon
func (d *SystemdDaemon) Stop() {
	d.logger.Debugf("Stopping daemon %s", d.Name())
	d.notifier.StopRoutine()
	<-d.notifier.RoutineStoppedChan()
	d.status = false
}

// Status returns true while the daemon is running
func (d *SystemdDaemon) Status() bool {
	return d.status
}

// Export returns the currently stored data in the daemon as a JSON formattable interface
func (d *SystemdDaemon) Export() interface{} {
	return d.data
}
