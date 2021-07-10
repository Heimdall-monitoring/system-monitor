package probes

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
	logt "github.com/sirupsen/logrus/hooks/test"
)

func newRAMMockRunnerOK() *mockRunner {
	command := `		  16316868 K total memory
		   5157872 K used memory
		   5985484 K active memory
		   2019764 K inactive memory
		   7319832 K free memory
		    662180 K buffer memory
		   3176984 K swap cache
		   7811068 K total swap
		         0 K used swap
		   7811068 K free swap
		    713033 non-nice user cpu ticks
		      1389 nice user cpu ticks
		    159011 system cpu ticks
		  12295837 idle cpu ticks
		     64540 IO-wait cpu ticks
		     26705 IRQ cpu ticks
		     12786 softirq cpu ticks
		         0 stolen cpu ticks
		   2268170 pages paged in
		   3317156 pages paged out
		         0 pages swapped in
		         0 pages swapped out
		  30644254 interrupts
		  87750333 CPU context switches
		1591874687 boot time
		    118575 forks`
	return newMockRunnerOK(command)
}

func newRAMMockRunnerNOK(lineWithError string) *mockRunner {
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
	if lineWithError == "available" {
		lines[0] = "  error K total memory"
	}
	if lineWithError == "used" {
		lines[1] = "   error K used memory"
	}
	if lineWithError == "free" {
		lines[4] = "   error K free memory"
	}
	if lineWithError == "shared" {
		lines[5] = "    error K buffer memory"
	}
	return newMockRunnerOK(strings.Join(lines, "\n"))
}

func TestParsing(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		res, err := getRAMUsage(newRAMMockRunnerNOK(""))
		if err != nil {
			t.Errorf("Unexpected error %s", err.Error())
		}
		export, err := json.Marshal(res)
		if err != nil {
			t.Errorf("Unexpected error exporting result %s", err.Error())
		}
		if string(export) != `{"available":16316868,"used":5157872,"free":7319832,"shared":662180}` {
			t.Errorf("Unexpected result %s", string(export))
		}
	})

	t.Run("Error available", func(t *testing.T) {
		t.Parallel()
		res, err := getRAMUsage(newRAMMockRunnerNOK("available"))
		if err == nil {
			t.Fatal("Error should not be nil")
		}
		if err.Error() != `Error parsing available RAM: strconv.ParseInt: parsing "error": invalid syntax` {
			t.Errorf("Unexpected error %s", err.Error())
		}
		export, err := json.Marshal(res)
		if err != nil {
			t.Errorf("Unexpected error exporting result %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected result %s", string(export))
		}
	})

	t.Run("Error used", func(t *testing.T) {
		t.Parallel()
		res, err := getRAMUsage(newRAMMockRunnerNOK("used"))
		if err == nil {
			t.Fatal("Error should not be nil")
		}
		if err.Error() != `Error parsing used RAM: strconv.ParseInt: parsing "error": invalid syntax` {
			t.Errorf("Unexpected error %s", err.Error())
		}
		export, err := json.Marshal(res)
		if err != nil {
			t.Errorf("Unexpected error exporting result %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected result %s", string(export))
		}
	})

	t.Run("Error free", func(t *testing.T) {
		t.Parallel()
		res, err := getRAMUsage(newRAMMockRunnerNOK("free"))
		if err == nil {
			t.Fatal("Error should not be nil")
		}
		if err.Error() != `Error parsing free RAM: strconv.ParseInt: parsing "error": invalid syntax` {
			t.Errorf("Unexpected error %s", err.Error())
		}
		export, err := json.Marshal(res)
		if err != nil {
			t.Errorf("Unexpected error exporting result %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected result %s", string(export))
		}
	})

	t.Run("Error shared", func(t *testing.T) {
		t.Parallel()
		res, err := getRAMUsage(newRAMMockRunnerNOK("shared"))
		if err == nil {
			t.Fatal("Error should not be nil")
		}
		if err.Error() != `Error parsing shared RAM: strconv.ParseInt: parsing "error": invalid syntax` {
			t.Errorf("Unexpected error %s", err.Error())
		}
		export, err := json.Marshal(res)
		if err != nil {
			t.Errorf("Unexpected error exporting result %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected result %s", string(export))
		}
	})
}

func TestRAMOk(t *testing.T) {
	t.Parallel()
	mockRunner := newRAMMockRunnerOK()

	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)

	testDaemon := &RAMUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           mockRunner,
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             RAMStats{},
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "ram-usage" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		export, err := json.Marshal(testDaemon.Export())
		if err != nil {
			t.Errorf("Unexpected error during export: %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected export %q", string(export))
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedCommands := []string{
			"/usr/bin/vmstat -s",
			"/usr/bin/vmstat -s",
		}
		expectedResult := `{"available":16316868,"used":5157872,"free":7319832,"shared":662180}`
		if strings.Join(mockRunner.commands, " | ") != strings.Join(expectedCommands, " | ") {
			t.Errorf("Unexpected commands %s", strings.Join(mockRunner.commands, " | "))
		}

		export, err := json.Marshal(testDaemon.Export())
		if err != nil {
			t.Errorf("Unexpected error during export %s", err.Error())
		}
		if string(export) != expectedResult {
			t.Errorf("Unexpected export %s", string(export))
		}

		if !testDaemon.Status() {
			t.Error("Status should be true")
		}

		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Starting daemon ram-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon ram-usage",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})

	t.Run("After stop", func(t *testing.T) {
		testDaemon.Stop()

		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Stopping daemon ram-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon ram-usage",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})
}

func TestRAMNOk(t *testing.T) {
	t.Parallel()
	mockRunner := newMockRunnerNOK("sample error")

	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)

	testDaemon := &RAMUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           mockRunner,
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             RAMStats{},
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "ram-usage" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		export, err := json.Marshal(testDaemon.Export())
		if err != nil {
			t.Errorf("Unexpected error during export: %s", err.Error())
		}
		if string(export) != `{"available":0,"used":0,"free":0,"shared":0}` {
			t.Errorf("Unexpected export %q", string(export))
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedCommands := []string{
			"/usr/bin/vmstat -s",
			"/usr/bin/vmstat -s",
		}
		expectedResult := `{"available":0,"used":0,"free":0,"shared":0}`
		if strings.Join(mockRunner.commands, " | ") != strings.Join(expectedCommands, " | ") {
			t.Errorf("Unexpected commands %s", strings.Join(mockRunner.commands, " | "))
		}

		export, err := json.Marshal(testDaemon.Export())
		if err != nil {
			t.Errorf("Unexpected error during export %s", err.Error())
		}
		if string(export) != expectedResult {
			t.Errorf("Unexpected export %s", string(export))
		}

		if !testDaemon.Status() {
			t.Error("Status should be true")
		}

		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Starting daemon ram-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon ram-usage",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.ErrorLevel,
				Message: "Error getting RAM usage data: An error occured during vmstat command execution: \"sample error\"",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.ErrorLevel,
				Message: "Error getting RAM usage data: An error occured during vmstat command execution: \"sample error\"",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})

	t.Run("After stop", func(t *testing.T) {
		testDaemon.Stop()

		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Stopping daemon ram-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon ram-usage",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})
}
