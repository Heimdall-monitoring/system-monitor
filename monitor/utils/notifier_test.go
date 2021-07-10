package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestNotifier(t *testing.T) {
	notifier := NewGoroutineNotifier()
	errsChan := make(chan error, 1)

	go func(notifier *GoroutineNotifier, errs chan<- error) {
		select {
		case <-notifier.StopSignalChan():
			notifier.ConfirmRoutineStopped()
		case <-time.After(50 * time.Millisecond):
			errs <- fmt.Errorf("No stop signal received in 50ms")
		}
	}(notifier, errsChan)

	notifier.StopRoutine()

	select {
	case err := <-errsChan:
		t.Fatalf("Unexpected error: %s", err)
	case <-notifier.RoutineStoppedChan():
		return
	case <-time.After(100 * time.Millisecond):
		t.Fatal("No confirmation received in 100ms")
	}
}
