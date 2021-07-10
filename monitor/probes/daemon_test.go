package probes

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"

	logt "github.com/sirupsen/logrus/hooks/test"
)

type mockData struct {
	Int1 int
	Str1 string
}

type testDaemon struct {
	status bool
	id     int
}

func newTestDaemon(id int) *testDaemon {
	return &testDaemon{status: false, id: id}
}

func (d *testDaemon) Start() {
	d.status = true
}

func (d *testDaemon) Stop() {
	d.status = false
}

func (d *testDaemon) Status() bool {
	return d.status
}

func (d *testDaemon) Name() string {
	return fmt.Sprintf("test-daemon-%d", d.id)
}

func (d *testDaemon) Export() interface{} {
	return mockData{0, "test"}
}

func TestMainDaemon(t *testing.T) {
	t.Parallel()

	logger, caplog := logt.NewNullLogger()
	logger.SetLevel(log.DebugLevel)

	testDaemon0 := newTestDaemon(0)
	testDaemon1 := newTestDaemon(1)
	testDaemon2 := newTestDaemon(2)
	testSubDaemons := []*testDaemon{testDaemon0, testDaemon1, testDaemon2}

	testMainDaemon := MainDaemon{
		logger:     logger,
		subDaemons: []Daemon{testDaemon0, testDaemon1, testDaemon2},
	}

	if testMainDaemon.Name() != "main-daemon" {
		t.Errorf("Invalid name for daemon %q, expecting 'main-daemon'", testMainDaemon.Name())
	}

	t.Run("all stopped", func(t *testing.T) {
		for _, daemon := range testSubDaemons {
			if daemon.status {
				t.Errorf("daemon %s should be stopped", daemon.Name())
			}
		}
		if testMainDaemon.Status() {
			t.Error("Main daemon should be stopped")
		}
		if len(caplog.Entries) != 0 {
			t.Errorf("Unexpected number of logs %d", len(caplog.Entries))
		}
		caplog.Reset()
	})

	t.Run("start all", func(t *testing.T) {
		testMainDaemon.Start()
		for _, daemon := range testSubDaemons {
			if !daemon.status {
				t.Errorf("daemon %s should be started", daemon.Name())
			}
		}
		if !testMainDaemon.Status() {
			t.Error("Main daemon should be started")
		}

		expectedEntries := []expectedLogEntry{
			{
				Level:   log.InfoLevel,
				Message: "Starting main daemon",
			},
			{
				Level:   log.InfoLevel,
				Message: "Starting daemon test-daemon-0",
			},
			{
				Level:   log.InfoLevel,
				Message: "Starting daemon test-daemon-1",
			},
			{
				Level:   log.InfoLevel,
				Message: "Starting daemon test-daemon-2",
			},
			{
				Level:   log.InfoLevel,
				Message: "Main daemon started",
			},
		}

		compareLogEntries(t, caplog.Entries, expectedEntries)
		caplog.Reset()
	})

	t.Run("stop all", func(t *testing.T) {
		if !testMainDaemon.Status() {
			t.Error("Main daemon should be started")
		}
		testMainDaemon.Stop()
		for _, daemon := range testSubDaemons {
			if daemon.status {
				t.Errorf("daemon %s should be stopped", daemon.Name())
			}
		}
		if testMainDaemon.Status() {
			t.Error("Main daemon should be stopped")
		}

		expectedEntries := []expectedLogEntry{
			{
				Level:   log.InfoLevel,
				Message: "Stopping main daemon",
			},
			{
				Level:   log.InfoLevel,
				Message: "Stopping daemon test-daemon-0",
			},
			{
				Level:   log.InfoLevel,
				Message: "Stopping daemon test-daemon-1",
			},
			{
				Level:   log.InfoLevel,
				Message: "Stopping daemon test-daemon-2",
			},
			{
				Level:   log.InfoLevel,
				Message: "Main daemon stopped",
			},
		}

		compareLogEntries(t, caplog.Entries, expectedEntries)
		caplog.Reset()
	})

	t.Run("stop already stopped daemons", func(t *testing.T) {
		if testMainDaemon.Status() {
			t.Error("Main daemon should already be stopped")
		}
		testMainDaemon.Stop()
		for _, daemon := range testSubDaemons {
			if daemon.status {
				t.Errorf("daemon %s should be stopped", daemon.Name())
			}
		}
		if testMainDaemon.Status() {
			t.Error("Main daemon should still be stopped")
		}

		expectedEntries := []expectedLogEntry{
			{
				Level:   log.InfoLevel,
				Message: "Stopping main daemon",
			},
			{
				Level:   log.WarnLevel,
				Message: "Daemon test-daemon-0 is already stopped, skipping...",
			},
			{
				Level:   log.WarnLevel,
				Message: "Daemon test-daemon-1 is already stopped, skipping...",
			},
			{
				Level:   log.WarnLevel,
				Message: "Daemon test-daemon-2 is already stopped, skipping...",
			},
			{
				Level:   log.InfoLevel,
				Message: "Main daemon stopped",
			},
		}

		compareLogEntries(t, caplog.Entries, expectedEntries)
		caplog.Reset()
	})

}
