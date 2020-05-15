package main

import (
	"github.com/aHugues/system-monitor/monitor/webserver"
)

func main() {
	// ramUsage := probes.GetRAMUsage()
	// fmt.Printf("Ram usage: %d total, %d used, %d free, %d shared\n", ramUsage.Available, ramUsage.Used, ramUsage.Free, ramUsage.Shared)

	// devices := probes.GetUsageStats(probes.CommandRunner)
	// if len(devices) == 0 {
	// 	fmt.Println("No storage devices found.")
	// } else {
	// 	for _, dev := range devices {
	// 		fmt.Println(dev.ToString())
	// 	}
	// }

	webserver.RunServer()
}
