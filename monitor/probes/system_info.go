package probes

import (
	"errors"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
)

// SystemInfo represent the global information for the operating system
type SystemInfo struct {
	OperatingSystem string `json:"operating-system"`
	Kernel          string `json:"kernel"`
	Distro          string `json:"distro"`
	Machine         string `json:"machine"`
	Name            string `json:"name"`
}

// systemInfoFromUname builds the system infos by running the uname command
func systemInfoFromUname(unameLine string) (SystemInfo, error) {
	r := regexp.MustCompile(`^(?P<os>[^ ]+) (?P<name>[^ ]+) (?P<kernel>[^ ]+) (?P<machine>[^ ]+)$`)
	match := r.FindStringSubmatch(unameLine)
	if r.MatchString(unameLine) == false {
		return SystemInfo{}, errors.New("No match")
	}
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	os := result["os"]
	kernel := result["kernel"]
	name := result["name"]
	machine := result["machine"]
	return SystemInfo{
		OperatingSystem: os,
		Kernel:          kernel,
		Machine:         machine,
		Name:            name,
	}, nil
}

// getDistro get the current distro name
func getDistro() (string, error) {
	log.Debug("Reading os-release for distro information")
	dat, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		log.Errorf("Error when reading os-release: %q", err)
		return "", err
	}
	r := regexp.MustCompile(`(?m)^NAME=\"(?P<distro>[^\"]+)\"$`)
	match := r.FindStringSubmatch(string(dat))
	if r.MatchString(string(dat)) == false {
		log.Error("No match for os-release format")
		return "", errors.New("No match")
	}

	log.Debugf("Distro found: %q", match[1])
	return match[1], nil
}

// getSystemInfo return the entire system info
func getSystemInfo(runner commandRunner) SystemInfo {
	unameCommand := []string{"/usr/bin/uname", "-s", "-n", "-m", "-r"}
	distro, _ := getDistro()

	CommandResult := runner.runCommand(unameCommand)
	unameResult := CommandResult.Stdout
	otherInfo, _ := systemInfoFromUname(unameResult)
	otherInfo.Distro = distro
	return otherInfo
}

// SystemInfoDaemon gets the system info for the probed machine
type SystemInfoDaemon struct {
	notifier         *utils.GoroutineNotifier
	runner           commandRunner
	refreshFrequency time.Duration
	logger           *log.Logger
	data             SystemInfo
	status           bool
}

// NewSystemInfoDaemon creates a SystemInfo daemon for the application
func NewSystemInfoDaemon(logger *log.Logger, refreshFrequency time.Duration) *SystemInfoDaemon {
	return &SystemInfoDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           LinuxCommandRunner{},
		refreshFrequency: refreshFrequency,
		logger:           logger,
		status:           false,
		data:             SystemInfo{},
	}
}

// Name returns 'system-info'
func (d *SystemInfoDaemon) Name() string {
	return "system-info"
}

func (d *SystemInfoDaemon) mainRoutine() {
	d.logger.Debugf("Starting main goroutine for daemon %s", d.Name())
	for {
		select {
		case <-d.notifier.StopSignalChan():
			d.logger.Debugf("Stopping main goroutine for daemon %s", d.Name())
			defer d.notifier.ConfirmRoutineStopped()
			return
		case <-time.After(d.refreshFrequency):
			d.logger.Trace("Refreshing data")
			d.data = getSystemInfo(d.runner)
		}
	}
}

// Start handles starting the system info daemon
func (d *SystemInfoDaemon) Start() {
	d.logger.Debugf("Starting daemon %s", d.Name())
	go d.mainRoutine()
	d.status = true
}

// Stop handles stopping the system info daemon
func (d *SystemInfoDaemon) Stop() {
	d.logger.Debugf("Stopping daemon %s", d.Name())
	d.notifier.StopRoutine()
	<-d.notifier.RoutineStoppedChan()
	d.status = false
}

// Status returns true while the daemon is running
func (d *SystemInfoDaemon) Status() bool {
	return d.status
}

// Export returns the currently stored data in the daemon as a JSON formattable interface
func (d *SystemInfoDaemon) Export() interface{} {
	return d.data
}
