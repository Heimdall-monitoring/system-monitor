package main

import (
	"fmt"

	"github.com/aHugues/system-monitor/monitor/probes"
)

func main() {
	ramUsage := probes.GetRAMUsage()
	fmt.Printf("Ram usage: %d total, %d used, %d free, %d shared", ramUsage.Available, ramUsage.Used, ramUsage.Free, ramUsage.Shared)
}
