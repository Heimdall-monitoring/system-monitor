package probes

import (
	"reflect"
	"strings"
	"testing"
)

var mockRunDisk func(command []string) CommandResult

type mockCommandRunnerDisk struct{}

func (r mockCommandRunnerDisk) runCommand(command []string) CommandResult {
	return mockRunDisk(command)
}

// Basic test when line corresponds to a correct device
func TestDeviceStatFromDfOk(t *testing.T) {
	device, err := DeviceStatFromDf("/dev/sdb6        41G   32G  6.7G  83% /")
	if err != nil {
		t.Fatal("Invalid test case")
	}
	if device.Filesystem != "/dev/sdb6" || device.MountPoint != "/" || device.Size != 41 || device.Used != 32 {
		t.Fatal("Invalid test case")
	}
}

// Test running the whole command with a sample output
func TestDiskOk(t *testing.T) {
	mockRunner := mockCommandRunnerDisk{}
	mockRunDisk = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/df", "-h", "-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"Filesystem      Size  Used Avail Use% Mounted on",
			"/dev/sdb6        41G   34G  5.1G  87% /",
			"/dev/sda5        96G   43G   48G  48% /home",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetUsageStats(mockRunner)

	device1 := stats[0]
	if device1.Filesystem != "/dev/sdb6" || device1.MountPoint != "/" || device1.Size != 41 || device1.Used != 34 {
		t.Fatal("Invalid test case for device1")
	}

	device2 := stats[1]
	if device2.Filesystem != "/dev/sda5" || device2.MountPoint != "/home" || device2.Size != 96 || device2.Used != 43 {
		t.Fatal("Invalid test case for device2")
	}
}

// Test running the whole command with an error
func TestDiskError(t *testing.T) {
	mockRunner := mockCommandRunnerDisk{}
	mockRunDisk = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/df", "-h", "-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"Filesystem      Size  Used Avail Use% Mounted on",
			"/dev/sdb6        41G   34G  5.1G  87% /",
			"/dev/sda5        96G   43G   48G  48% /home",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "toto",
			StatusCode: 1,
		}
	}

	stats := GetUsageStats(mockRunner)

	if len(stats) != 0 {
		t.Fatal("Expecting no result")
	}
}

// Test displaying a string from a Device
func TestToString(t *testing.T) {
	device := DeviceStat{"/dev/sda1", "/", 42, 13}
	if device.ToString() != "Device /dev/sda1 - mountpoint /\tTotal size 42GB - Used 13GB" {
		t.Fatal("Invalid test case")
	}
}
