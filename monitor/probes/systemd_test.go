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

type mockRunSystemd struct {
	t *testing.T
}

func (r *mockRunSystemd) runCommand(command []string) CommandResult {
	commandBegin := strings.Join(command[:len(command)-1], " ")
	if commandBegin != "/bin/systemctl is-active --quiet" {
		r.t.Fatalf("Unexpected command %s", strings.Join(command, " "))
	}
	service := command[len(command)-1]
	if service == "service-ok.service" {
		return CommandResult{StatusCode: 0}
	} else if service == "service-error.service" {
		return CommandResult{StatusCode: 1}
	}
	r.t.Fatalf("Unexpected service %s", service)
	return CommandResult{}
}

func TestDaemon(t *testing.T) {
	t.Parallel()
	mockRunner := &mockRunSystemd{t}

	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.TraceLevel)

	testDaemon := &SystemdDaemon{
		notifier:         utils.NewGoroutineNotifier(),
		runner:           mockRunner,
		refreshFrequency: 30 * time.Millisecond,
		logger:           logger,
		status:           false,
		data:             make(map[string]bool),
		services:         []string{"service-ok.service", "service-error.service"},
	}

	t.Run("Before start", func(t *testing.T) {
		if testDaemon.Name() != "systemd" {
			t.Errorf("Unexpected name %s", testDaemon.Name())
		}
		if testDaemon.Status() {
			t.Error("Status should be false")
		}
		export, err := json.Marshal(testDaemon.Export())
		if err != nil {
			t.Errorf("Unexpected error during export: %s", err.Error())
		}
		if string(export) != `{}` {
			t.Errorf("Unexpected export %q", string(export))
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
	})

	t.Run("After start", func(t *testing.T) {
		testDaemon.Start()
		time.Sleep(75 * time.Millisecond)

		expectedResult := `{"service-error.service":false,"service-ok.service":true}`

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
				Message: "Starting daemon systemd",
			},
			{
				Level:   log.DebugLevel,
				Message: "Starting main goroutine for daemon systemd",
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
				Message: "Stopping daemon systemd",
			},
			{
				Level:   log.DebugLevel,
				Message: "Stopping main goroutine for daemon systemd",
			},
		}
		compareLogEntries(t, caplog.Entries, expectedLogs)
		caplog.Reset()
	})
}

// // Test a single up service
// func TestServiceUp(t *testing.T) {
// 	mockRunner := mockCommandRunnerSystemd{}
// 	mockRunSystemd = func(command []string) CommandResult {
// 		expectedCommand := []string{"/bin/systemctl", "is-active", "--quiet", "test-service.service"}
// 		if !reflect.DeepEqual(command, expectedCommand) {
// 			t.Fatal("Command is invalid")
// 		}
// 		return CommandResult{StatusCode: 0}
// 	}

// 	serviceStatus := getServiceStatus(mockRunner, "test-service.service")
// 	if !serviceStatus {
// 		t.Fatal("Invalid service status")
// 	}
// }

// // Test a single down service
// func TestServiceDown(t *testing.T) {
// 	mockRunner := mockCommandRunnerSystemd{}
// 	mockRunSystemd = func(command []string) CommandResult {
// 		expectedCommand := []string{"/bin/systemctl", "is-active", "--quiet", "test-service.service"}
// 		if !reflect.DeepEqual(command, expectedCommand) {
// 			t.Fatal("Command is invalid")
// 		}
// 		return CommandResult{StatusCode: 1}
// 	}

// 	serviceStatus := getServiceStatus(mockRunner, "test-service.service")
// 	if serviceStatus {
// 		t.Fatal("Invalid service status")
// 	}
// }

// // Test several services
// func TestSeveralServices(t *testing.T) {
// 	callCount := 0
// 	mockRunner := mockCommandRunnerSystemd{}
// 	mockRunSystemd = func(command []string) CommandResult {
// 		var expectedCommand []string
// 		var result int
// 		switch callCount {
// 		case 0:
// 			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-test-1.service"}
// 			result = 0
// 		case 1:
// 			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-error"}
// 			result = 1
// 		default:
// 			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-ok"}
// 			result = 0
// 		}
// 		if !reflect.DeepEqual(command, expectedCommand) {
// 			t.Fatalf("Command is invalid, expected %q, got %q", expectedCommand, command)
// 		}
// 		callCount++
// 		return CommandResult{StatusCode: result}
// 	}

// 	statuses := getServicesStatuses(mockRunner, []string{"service-test-1.service", "service-error", "service-ok"})
// 	if !statuses["service-test-1.service"] || statuses["service-error"] || !statuses["service-ok"] {
// 		t.Fatal("Invalid service status")
// 	}
// }
