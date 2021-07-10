package probes

import (
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

type expectedLogEntry struct {
	Level   log.Level
	Message string
	Fields  map[string]interface{}
}

func compareLogEntries(t *testing.T, logLst []log.Entry, expectedLogLst []expectedLogEntry) {
	if logLst == nil || expectedLogLst == nil {
		t.Fatalf("unexpected nil value while comparing logs")
	}
	if len(logLst) != len(expectedLogLst) {
		t.Fatalf("missing or spurious log entries; expected %d entries, found %d entries", len(expectedLogLst), len(logLst))
	}
	for i := range logLst {
		if logLst[i].Message != expectedLogLst[i].Message || logLst[i].Level != expectedLogLst[i].Level {
			t.Errorf("unexpected log entry [%d]: expected (%d) %s [%v], found (%d) %s [%v]", i, expectedLogLst[i].Level, expectedLogLst[i].Message, expectedLogLst[i].Fields, logLst[i].Level, logLst[i].Message, logLst[i].Data)
		}
		if len(logLst[i].Data) != len(expectedLogLst[i].Fields) {
			t.Fatalf("missing or spurious fields: expected %d; found %d", len(expectedLogLst[i].Fields), len(logLst[i].Data))
		}
		for key, expectedValue := range expectedLogLst[i].Fields {
			actualValue, ok := logLst[i].Data[key]
			if !ok || actualValue != expectedValue {
				t.Errorf("unexpected log entry [%d]: expected (%d) %s [%v], found (%d) %s [%v]", i, expectedLogLst[i].Level, expectedLogLst[i].Message, expectedLogLst[i].Fields, logLst[i].Level, logLst[i].Message, logLst[i].Data)
			}
		}
	}
}

type mockRunner struct {
	stdout     string
	stderr     string
	returncode int
	commands   []string
}

func (r *mockRunner) runCommand(command []string) CommandResult {
	r.commands = append(r.commands, strings.Join(command, " "))
	return CommandResult{
		Stderr:     r.stderr,
		Stdout:     r.stdout,
		StatusCode: r.returncode,
	}
}

// newMockRunnerOK returns a mockRunner for a command that will suceed
func newMockRunnerOK(result string) *mockRunner {
	return &mockRunner{
		stdout:     result,
		stderr:     "",
		returncode: 0,
	}
}

// newMockRunnerNOK returns a mockRunner for a command that will fail
func newMockRunnerNOK(err string) *mockRunner {
	return &mockRunner{
		stdout:     "",
		stderr:     err,
		returncode: 1,
	}
}
