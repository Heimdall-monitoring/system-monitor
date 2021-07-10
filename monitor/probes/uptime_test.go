package probes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/aHugues/system-monitor/monitor/utils"
	log "github.com/sirupsen/logrus"
	logt "github.com/sirupsen/logrus/hooks/test"
)

func TestUptimeOK(t *testing.T) {
	t.Parallel()
	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)
	uptimeFile := path.Join(t.TempDir(), "uptime")

	ioutil.WriteFile(uptimeFile, []byte("8202.92 129956.14"), os.FileMode(0o666))

	testDaemon := &UptimeDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             -1,
		uptimeFile:       uptimeFile,
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "uptime" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		if testDaemon.data != -1 {
			t.Errorf("Invalid value for data %d", testDaemon.data)
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedResult := `8202`
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
				Message: "Starting daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon uptime",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.DebugLevel,
				Message: "Reading uptime from /proc/uptime",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.DebugLevel,
				Message: "Reading uptime from /proc/uptime",
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
				Message: "Stopping daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon uptime",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})
}

func TestUptimeNOK(t *testing.T) {
	t.Parallel()
	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)
	uptimeFile := path.Join(t.TempDir(), "uptime")

	// ioutil.WriteFile(uptimeFile, []byte("8202.92 129956.14"), os.FileMode(0o666))

	testDaemon := &UptimeDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             -1,
		uptimeFile:       uptimeFile,
	}

	t.Run("No file", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(50 * time.Millisecond)
		testDaemon.Stop()

		if testDaemon.data != -1 {
			t.Errorf("Unexpected data: %d", testDaemon.data)
		}
		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Starting daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon uptime",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.DebugLevel,
				Message: "Reading uptime from /proc/uptime",
			},
			{
				Level:   log.ErrorLevel,
				Message: fmt.Sprintf("Error refreshing data: Impossible to read uptime: open %s: no such file or directory", uptimeFile),
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon uptime",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})

	t.Run("Invalid content", func(t *testing.T) {
		ioutil.WriteFile(uptimeFile, []byte("invalid content"), os.FileMode(0o666))

		testDaemon.Start()
		time.Sleep(50 * time.Millisecond)
		testDaemon.Stop()

		if testDaemon.data != -1 {
			t.Errorf("Unexpected data: %d", testDaemon.data)
		}
		expectedLogs := []expectedLogEntry{
			{
				Level:   log.DebugLevel,
				Message: "Starting daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon uptime",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.DebugLevel,
				Message: "Reading uptime from /proc/uptime",
			},
			{
				Level:   log.ErrorLevel,
				Message: "Error refreshing data: Impossible to parse uptime value: strconv.ParseFloat: parsing \"invalid\": invalid syntax",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping daemon uptime",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon uptime",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})
}
