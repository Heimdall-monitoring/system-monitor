package webserver

import (
	"github.com/aHugues/system-monitor/monitor/probes"
	"github.com/aHugues/system-monitor/monitor/utils"
)

// fullStats represent the complete stats returned to the user
type fullStats struct {
	DiskUsage  []probes.DeviceStat `json:"disk-usage,omitempty"`
	RAMUsage   probes.RAMStats     `json:"ram-usage,omitempty"`
	SystemInfo probes.SystemInfo   `json:"system-info,omitempty"`
	Services   map[string]bool     `json:"services-status,omitempty"`
	Uptime     int64               `json:"uptime,omitempty"`
}

func (obj *fullStats) update(data fullStats) {
	obj.DiskUsage = data.DiskUsage
	obj.RAMUsage = data.RAMUsage
	obj.SystemInfo = data.SystemInfo
	obj.Services = data.Services
	obj.Uptime = data.Uptime
}

func getFullStats(config utils.ProbesConfig) fullStats {
	fullStats := fullStats{}

	// if len(config.SystemdServices) > 0 {
	// 	fullStats.Services = probes.GetServicesStatuses(probes.LinuxCommandRunner{}, config.SystemdServices)
	// }

	// if config.RAMUsage {
	// 	fullStats.RAMUsage = probes.GetRAMUsage(probes.LinuxCommandRunner{})
	// }

	// if config.DiskUsage {
	// 	fullStats.DiskUsage = probes.GetUsageStats(probes.LinuxCommandRunner{})
	// }

	// if config.SystemInfo {
	// 	fullStats.SystemInfo = probes.GetSystemInfo(probes.LinuxCommandRunner{})
	// }

	// if config.Uptime {
	// 	uptime, err := probes.GetUptime()
	// 	if err == nil {
	// 		fullStats.Uptime = uptime
	// 	}
	// }
	return fullStats
}
