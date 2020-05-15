package probes

import (
	"errors"
	"io/ioutil"
	"regexp"

	log "github.com/cihub/seelog"
)

// SystemInfo represent the global information for the operating system
type SystemInfo struct {
	OperatingSystem string `json:"operating-system"`
	Kernel          string `json:"kernel"`
	Distro          string `json:"distro"`
	Machine         string `json:"machine"`
	Name            string `json:"name"`
}

// SystemInfoFromUname builds the system infos by running the uname command
func SystemInfoFromUname(unameLine string) (SystemInfo, error) {
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

// GetDistro get the current distro name
func GetDistro() (string, error) {
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

// GetSystemInfo return the entire system info
func GetSystemInfo(runner commandRunner) SystemInfo {
	unameCommand := []string{"/usr/bin/uname", "-s", "-n", "-m", "-r"}
	distro, _ := GetDistro()

	CommandResult := runner.runCommand(unameCommand)
	unameResult := CommandResult.Stdout
	otherInfo, _ := SystemInfoFromUname(unameResult)
	otherInfo.Distro = distro
	return otherInfo
}
