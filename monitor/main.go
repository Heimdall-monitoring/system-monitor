package main

import (
	"log"

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
	// config := webserver.Configuration{
	// 	Server: webserver.Server{
	// 		Host: "127.0.0.1",
	// 		Port: 5000,
	// 	},
	// }

	config, err := webserver.ReadConfigJSON("config.json")
	if err != nil {
		log.Fatalf("Error reading config file %q", err)
	}
	webserver.RunServer(config)
}
