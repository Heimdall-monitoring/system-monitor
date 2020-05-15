package probes

import (
	"reflect"
	"strings"
	"testing"
)

var mockRunRAM func(command []string) CommandResult

type mockCommandRunnerRAM struct{}

func (r mockCommandRunnerRAM) runCommand(command []string) CommandResult {
	return mockRunRAM(command)
}

func TestRAMOk(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"  16316868 K total memory",
			"   5157872 K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"   7319832 K free memory",
			"    662180 K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 16316868 || stats.Used != 5157872 || stats.Free != 7319832 || stats.Shared != 662180 {
		t.Fatal("Invalid result")
	}
}

func TestRAMError1(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"  16316868 K total memory",
			"   5157872 K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"   7319832 K free memory",
			"    662180 K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "toto",
			StatusCode: 1,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 0 || stats.Used != 0 || stats.Free != 0 || stats.Shared != 0 {
		t.Fatal("Invalid result")
	}
}

func TestRAMError2(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"     error K total memory",
			"   5157872 K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"   7319832 K free memory",
			"    662180 K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 0 || stats.Used != 0 || stats.Free != 0 || stats.Shared != 0 {
		t.Fatal("Invalid result")
	}
}

func TestRAMError3(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"  16316868 K total memory",
			"     error K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"   7319832 K free memory",
			"    662180 K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 0 || stats.Used != 0 || stats.Free != 0 || stats.Shared != 0 {
		t.Fatal("Invalid result")
	}
}

func TestRAMError4(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"  16316868 K total memory",
			"   5157872 K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"     error K free memory",
			"    662180 K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 0 || stats.Used != 0 || stats.Free != 0 || stats.Shared != 0 {
		t.Fatal("Invalid result")
	}
}

func TestRAMError5(t *testing.T) {
	mockRunner := mockCommandRunnerRAM{}
	mockRunRAM = func(command []string) CommandResult {
		expectedCommand := []string{"/usr/bin/vmstat", "-s"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		lines := []string{
			"  16316868 K total memory",
			"   5157872 K used memory",
			"   5985484 K active memory",
			"   2019764 K inactive memory",
			"   7319832 K free memory",
			"     error K buffer memory",
			"   3176984 K swap cache",
			"   7811068 K total swap",
			"         0 K used swap",
			"   7811068 K free swap",
			"    713033 non-nice user cpu ticks",
			"      1389 nice user cpu ticks",
			"    159011 system cpu ticks",
			"  12295837 idle cpu ticks",
			"     64540 IO-wait cpu ticks",
			"     26705 IRQ cpu ticks",
			"     12786 softirq cpu ticks",
			"         0 stolen cpu ticks",
			"   2268170 pages paged in",
			"   3317156 pages paged out",
			"         0 pages swapped in",
			"         0 pages swapped out",
			"  30644254 interrupts",
			"  87750333 CPU context switches",
			"1591874687 boot time",
			"    118575 forks",
		}
		return CommandResult{
			Stdout:     strings.Join(lines, "\n"),
			Stderr:     "",
			StatusCode: 0,
		}
	}

	stats := GetRAMUsage(mockRunner)
	if stats.Available != 0 || stats.Used != 0 || stats.Free != 0 || stats.Shared != 0 {
		t.Fatal("Invalid result")
	}
}
