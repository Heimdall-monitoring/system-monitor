package probes

import (
	"reflect"
	"testing"
)

var mockRunSystemd func(command []string) CommandResult

type mockCommandRunnerSystemd struct{}

func (r mockCommandRunnerSystemd) runCommand(command []string) CommandResult {
	return mockRunSystemd(command)
}

// Test a single up service
func TestServiceUp(t *testing.T) {
	mockRunner := mockCommandRunnerSystemd{}
	mockRunSystemd = func(command []string) CommandResult {
		expectedCommand := []string{"/bin/systemctl", "is-active", "--quiet", "test-service.service"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		return CommandResult{StatusCode: 0}
	}

	serviceStatus := GetServiceStatus(mockRunner, "test-service.service")
	if !serviceStatus {
		t.Fatal("Invalid service status")
	}
}

// Test a single down service
func TestServiceDown(t *testing.T) {
	mockRunner := mockCommandRunnerSystemd{}
	mockRunSystemd = func(command []string) CommandResult {
		expectedCommand := []string{"/bin/systemctl", "is-active", "--quiet", "test-service.service"}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatal("Command is invalid")
		}
		return CommandResult{StatusCode: 1}
	}

	serviceStatus := GetServiceStatus(mockRunner, "test-service.service")
	if serviceStatus {
		t.Fatal("Invalid service status")
	}
}

// Test several services
func TestSeveralServices(t *testing.T) {
	callCount := 0
	mockRunner := mockCommandRunnerSystemd{}
	mockRunSystemd = func(command []string) CommandResult {
		var expectedCommand []string
		var result int
		switch callCount {
		case 0:
			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-test-1.service"}
			result = 0
		case 1:
			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-error"}
			result = 1
		default:
			expectedCommand = []string{"/bin/systemctl", "is-active", "--quiet", "service-ok"}
			result = 0
		}
		if !reflect.DeepEqual(command, expectedCommand) {
			t.Fatalf("Command is invalid, expected %q, got %q", expectedCommand, command)
		}
		callCount++
		return CommandResult{StatusCode: result}
	}

	statuses := GetServicesStatuses(mockRunner, []string{"service-test-1.service", "service-error", "service-ok"})
	if !statuses["service-test-1.service"] || statuses["service-error"] || !statuses["service-ok"] {
		t.Fatal("Invalid service status")
	}
}
