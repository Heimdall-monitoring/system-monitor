package webserver

import (
	"github.com/aHugues/system-monitor/monitor/probes"
)

// fullStats represent the complete stats returned to the user
type fullStats struct {
	DiskUsage  []probes.DeviceStat `json:"disk-usage"`
	RAMUsage   probes.RAMStats     `json:"ram-usage"`
	SystemInfo probes.SystemInfo   `json:"system-info"`
	Services   map[string]bool     `json:"services-status"`
}

func getFullStats() fullStats {
	services := []string{"sshd"}
	RAMUsage := probes.GetRAMUsage(probes.LinuxCommandRunner{})
	diskStats := probes.GetUsageStats(probes.LinuxCommandRunner{})
	systemInfo := probes.GetSystemInfo(probes.LinuxCommandRunner{})
	serviceStatuses := probes.GetServicesStatuses(probes.LinuxCommandRunner{}, services)
	fullStats := fullStats{
		DiskUsage:  diskStats,
		RAMUsage:   RAMUsage,
		SystemInfo: systemInfo,
		Services:   serviceStatuses,
	}
	return fullStats
}
