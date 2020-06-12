package probes

import log "github.com/cihub/seelog"

// GetServiceStatus returns True if the probed status is running
func GetServiceStatus(runner commandRunner, service string) bool {
	systemdCommand := []string{"/bin/systemctl", "is-active", "--quiet", service}
	commandResult := runner.runCommand(systemdCommand)

	log.Debugf("Systemd service %q statuscode: %d", service, commandResult.StatusCode)
	return commandResult.StatusCode == 0
}

// GetServicesStatuses computes the statuses for a list of services
func GetServicesStatuses(runner commandRunner, services []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range services {
		result[s] = GetServiceStatus(runner, s)
	}
	return result
}
