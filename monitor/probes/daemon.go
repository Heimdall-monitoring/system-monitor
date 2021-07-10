package probes

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Daemon interface is implemented by all probing daemons and allow a high level interface to handle goroutines
type Daemon interface {
	Start()
	Stop()
	Status() bool
	Name() string
	Export() interface{}
}

// MainDaemon is the main daemon that will be used when the application is starting
type MainDaemon struct {
	subDaemons []Daemon
	logger     *log.Logger
}

// NewMainDaemon returns a new MainDaemon to handle the main application
func NewMainDaemon(logger *log.Logger) *MainDaemon {
	return &MainDaemon{
		subDaemons: []Daemon{
			NewDiskUsageDaemon(logger, 10*time.Second),
			NewRAMUsageDaemon(logger, 1*time.Second),
			NewSystemInfoDaemon(logger, 30*time.Second),
		},
		logger: logger,
	}
}

// Start stats all subdaemons attached to the main daemon
func (d *MainDaemon) Start() {
	d.logger.Info("Starting main daemon")
	for _, daemon := range d.subDaemons {
		d.logger.Infof("Starting daemon %s", daemon.Name())
		daemon.Start()
	}
	d.logger.Info("Main daemon started")
}

// Stop stops all subdaemons attached to the main daemon
func (d *MainDaemon) Stop() {
	d.logger.Info("Stopping main daemon")
	for _, daemon := range d.subDaemons {
		if !daemon.Status() {
			d.logger.Warnf("Daemon %s is already stopped, skipping...", daemon.Name())
			continue
		}
		d.logger.Infof("Stopping daemon %s", daemon.Name())
		daemon.Stop()
	}
	d.logger.Info("Main daemon stopped")
}

// Status returns true while all subdaemons are running
func (d *MainDaemon) Status() bool {
	for _, daemon := range d.subDaemons {
		if !daemon.Status() {
			return false
		}
	}
	return true
}

// Name returns the name of the Main Daemon
func (d *MainDaemon) Name() string {
	return "main-daemon"
}
