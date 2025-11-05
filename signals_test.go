package main

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewTerminationManager(t *testing.T) {
	handler := func(sig os.Signal, immediate bool) {
		// Handler for testing
	}

	tm := newTerminationManager(handler)
	if tm == nil {
		t.Fatal("newTerminationManager() returned nil")
	}
	defer tm.Shutdown()

	if tm.stopCh == nil {
		t.Error("stopCh is nil")
	}
	if tm.immediateCh == nil {
		t.Error("immediateCh is nil")
	}
	if tm.shutdownCh == nil {
		t.Error("shutdownCh is nil")
	}
	if tm.done == nil {
		t.Error("done channel is nil")
	}
}

func TestTerminationManager_StopSignals(t *testing.T) {
	tm := newTerminationManager(nil)
	defer tm.Shutdown()

	signals := tm.StopSignals()
	if signals.stop == nil {
		t.Error("stop channel is nil")
	}
	if signals.immediate == nil {
		t.Error("immediate channel is nil")
	}

	// Channels should be open initially
	select {
	case <-signals.stop:
		t.Error("stop channel should not be closed initially")
	default:
		// Expected
	}

	select {
	case <-signals.immediate:
		t.Error("immediate channel should not be closed initially")
	default:
		// Expected
	}
}

func TestTerminationManager_RequestStop(t *testing.T) {
	tm := newTerminationManager(nil)
	defer tm.Shutdown()

	signals := tm.StopSignals()

	// Stop channel should be open
	select {
	case <-signals.stop:
		t.Error("stop channel should not be closed before RequestStop")
	default:
		// Expected
	}

	tm.RequestStop()

	// Stop channel should now be closed
	select {
	case <-signals.stop:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("stop channel should be closed after RequestStop")
	}

	// Multiple calls should be safe
	tm.RequestStop()
	tm.RequestStop()
}

func TestTerminationManager_Shutdown(t *testing.T) {
	tm := newTerminationManager(nil)

	// Shutdown should close stop channel
	tm.Shutdown()

	signals := tm.StopSignals()
	select {
	case <-signals.stop:
		// Expected
	default:
		t.Error("stop channel should be closed after Shutdown")
	}

	// Should wait for listener goroutine to finish
	// The done channel should be closed
	select {
	case <-tm.done:
		// Expected
	case <-time.After(500 * time.Millisecond):
		t.Error("done channel should be closed after Shutdown")
	}

	// Multiple shutdowns should be safe
	tm.Shutdown()
}

func TestTerminationManager_FirstSignal(t *testing.T) {
	t.Skip("Skipping test that sends OS signals which interfere with test execution")
}

func TestTerminationManager_SecondSignal(t *testing.T) {
	t.Skip("Skipping test that sends OS signals which interfere with test execution")
}

func TestTerminationManager_ConcurrentRequestStop(t *testing.T) {
	tm := newTerminationManager(nil)
	defer tm.Shutdown()

	// Multiple concurrent RequestStop calls should be safe
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tm.RequestStop()
		}()
	}

	wg.Wait()

	// Stop channel should be closed exactly once
	signals := tm.StopSignals()
	select {
	case <-signals.stop:
		// Expected
	default:
		t.Error("stop channel should be closed")
	}
}

func TestTerminationManager_TriggerImmediate(t *testing.T) {
	tm := newTerminationManager(nil)
	defer tm.Shutdown()

	signals := tm.StopSignals()

	// Initially immediate should be open
	select {
	case <-signals.immediate:
		t.Error("immediate channel should not be closed initially")
	default:
		// Expected
	}

	tm.triggerImmediate()

	// Immediate channel should now be closed
	select {
	case <-signals.immediate:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("immediate channel should be closed after triggerImmediate")
	}

	// Stop should also be closed
	select {
	case <-signals.stop:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("stop channel should be closed after triggerImmediate")
	}

	// Multiple calls should be safe
	tm.triggerImmediate()
}

func TestTerminationManager_WithNilHandler(t *testing.T) {
	tm := newTerminationManager(nil)
	defer tm.Shutdown()

	// Should work fine with nil handler
	tm.RequestStop()
	tm.triggerImmediate()

	signals := tm.StopSignals()
	select {
	case <-signals.stop:
		// Expected
	default:
		t.Error("stop channel should be closed")
	}
}
