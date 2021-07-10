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

func TestDeviceStatFromDf(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		device, err := DeviceStatFromDf("/dev/sdb6        41G   32G  6.7G  83% /")
		if err != nil {
			t.Fatal("Unexpected error")
		}
		if device.Filesystem != "/dev/sdb6" {
			t.Errorf("Unexpected filesystem %s", device.Filesystem)
		}
		if device.MountPoint != "/" {
			t.Errorf("Unexpected mount point %s", device.MountPoint)
		}
		if device.Size != 41 {
			t.Errorf("Unexpected size %d", device.Size)
		}
		if device.Used != 32 {
			t.Errorf("Unexpected used size %d", device.Used)
		}
		if device.ToString() != "Device /dev/sdb6 - mountpoint /	Total size 41GB - Used 32GB" {
			t.Errorf("Unexpected result of ToString %q", device.ToString())
		}
	})

	t.Run("NOK", func(t *testing.T) {
		t.Parallel()
		_, err := DeviceStatFromDf("it's a trap")
		if err == nil {
			t.Fatal("Error should not be nil")
		}
		if err.Error() != "Line \"it's a trap\" does not match a df line" {
			t.Errorf("Unexpected error %q", err.Error())
		}
	})
}

// Test running the whole command with a sample output
func TestDiskOk(t *testing.T) {
	t.Parallel()
	res := `Filesystem      Size  Used Avail Use% Mounted on
/dev/sdb6        41G   34G  5.1G  87% /
/dev/sda5        96G   43G   48G  48% /home`
	mockRunner := newMockRunnerOK(res)
	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)

	testDaemon := &DiskUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           mockRunner,
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             []DeviceStat{},
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "disk-usage" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		if len(testDaemon.data) != 0 {
			t.Errorf("Invalid length for data %d", len(testDaemon.data))
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedCommands := []string{
			"/usr/bin/df -h -x tmpfs -x devtmpfs -x squashfs",
			"/usr/bin/df -h -x tmpfs -x devtmpfs -x squashfs",
		}
		expectedResult := `[{"filesystem":"/dev/sdb6","mountpoint":"/","size":41,"used":34},{"filesystem":"/dev/sda5","mountpoint":"/home","size":96,"used":43}]`
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
				Message: "Starting daemon disk-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon disk-usage",
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
				Message: "Stopping daemon disk-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon disk-usage",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})

}

// Test running the whole command with an error
func TestDiskNOk(t *testing.T) {
	t.Parallel()
	res := "sample error"
	mockRunner := newMockRunnerNOK(res)
	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)

	testDaemon := &DiskUsageDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           mockRunner,
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             []DeviceStat{},
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "disk-usage" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		if len(testDaemon.data) != 0 {
			t.Errorf("Invalid length for data %d", len(testDaemon.data))
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedCommands := []string{
			"/usr/bin/df -h -x tmpfs -x devtmpfs -x squashfs",
			"/usr/bin/df -h -x tmpfs -x devtmpfs -x squashfs",
		}
		expectedResult := `[]`
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
				Message: "Starting daemon disk-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon disk-usage",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.ErrorLevel,
				Message: "Error getting disk usage stats: Error during df command execution: sample error",
			},
			{
				Level:   log.TraceLevel,
				Message: "Refreshing data",
			},
			{
				Level:   log.ErrorLevel,
				Message: "Error getting disk usage stats: Error during df command execution: sample error",
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
				Message: "Stopping daemon disk-usage",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon disk-usage",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})

}
