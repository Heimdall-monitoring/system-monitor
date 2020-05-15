package webserver

import (
	"github.com/aHugues/system-monitor/monitor/probes"
)

// fullStats represent the complete stats returned to the user
type fullStats struct {
	DiskUsage  []probes.DeviceStat `json:"disk-usage"`
	RAMUsage   probes.RAMStats     `json:"ram-usage"`
	SystemInfo probes.SystemInfo   `json:"system-info"`
}

func getFullStats() fullStats {
	RAMUsage := probes.GetRAMUsage(probes.LinuxCommandRunner{})
	diskStats := probes.GetUsageStats(probes.LinuxCommandRunner{})
	systemInfo := probes.GetSystemInfo(probes.LinuxCommandRunner{})
	fullStats := fullStats{
		DiskUsage:  diskStats,
		RAMUsage:   RAMUsage,
		SystemInfo: systemInfo,
	}
	return fullStats
}
