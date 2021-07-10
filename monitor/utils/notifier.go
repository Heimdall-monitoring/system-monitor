package utils

// GoroutineNotifier handles notifying goroutines that they should stop, and waiting for them to be stopped
type GoroutineNotifier struct {
	stopRoutineChan    chan struct{}
	routineStoppedChan chan struct{}
}

// NewGoroutineNotifier creates a new GoroutineNotifier object
func NewGoroutineNotifier() *GoroutineNotifier {
	return &GoroutineNotifier{
		stopRoutineChan:    make(chan struct{}, 1),
		routineStoppedChan: make(chan struct{}, 1),
	}
}

// StopRoutine sends a signal to the notifier that its attached routine should be stopped
func (notifier *GoroutineNotifier) StopRoutine() {
	notifier.stopRoutineChan <- struct{}{}
}

// ConfirmRoutineStopped sends a signal to the notifier that its attached routine has been stopped
func (notifier *GoroutineNotifier) ConfirmRoutineStopped() {
	notifier.routineStoppedChan <- struct{}{}
}

// StopSignalChan returns a unidirectional channel to listen to for a stop signal
func (notifier *GoroutineNotifier) StopSignalChan() <-chan struct{} {
	return notifier.stopRoutineChan
}

// RoutineStoppedChan returns a unidirectional channel to listen to to confirm the goroutine has stopped
func (notifier *GoroutineNotifier) RoutineStoppedChan() <-chan struct{} {
	return notifier.routineStoppedChan
}
