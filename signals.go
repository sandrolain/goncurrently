package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type stopSignals struct {
	stop      <-chan struct{}
	immediate <-chan struct{}
}

type terminationManager struct {
	stopCh        chan struct{}
	stopOnce      sync.Once
	immediateCh   chan struct{}
	immediateOnce sync.Once
	shutdownCh    chan struct{}
	shutdownOnce  sync.Once
	done          chan struct{}
	handler       func(os.Signal, bool)
}

func newTerminationManager(handler func(os.Signal, bool)) *terminationManager {
	tm := &terminationManager{
		stopCh:      make(chan struct{}),
		immediateCh: make(chan struct{}),
		shutdownCh:  make(chan struct{}),
		done:        make(chan struct{}),
		handler:     handler,
	}
	go tm.listen()
	return tm
}

func (tm *terminationManager) listen() {
	defer close(tm.done)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	firstHandled := false
	for {
		select {
		case sig := <-sigCh:
			if sig == nil {
				continue
			}
			if !firstHandled {
				if tm.handler != nil {
					tm.handler(sig, false)
				}
				tm.RequestStop()
				firstHandled = true
			} else {
				if tm.handler != nil {
					tm.handler(sig, true)
				}
				tm.triggerImmediate()
				return
			}
		case <-tm.shutdownCh:
			return
		}
	}
}

func (tm *terminationManager) triggerImmediate() {
	tm.immediateOnce.Do(func() {
		close(tm.immediateCh)
	})
	tm.RequestStop()
}

func (tm *terminationManager) StopSignals() stopSignals {
	return stopSignals{
		stop:      tm.stopCh,
		immediate: tm.immediateCh,
	}
}

func (tm *terminationManager) RequestStop() {
	tm.stopOnce.Do(func() {
		close(tm.stopCh)
	})
}

func (tm *terminationManager) triggerShutdown() {
	tm.shutdownOnce.Do(func() {
		close(tm.shutdownCh)
	})
}

func (tm *terminationManager) Shutdown() {
	tm.RequestStop()
	tm.triggerShutdown()
	<-tm.done
}
