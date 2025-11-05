package main

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestConsoleRouter_Add_Done_Wait(t *testing.T) {
	router := &consoleRouter{
		wg: sync.WaitGroup{},
	}

	// Test Add increments wait group
	router.Add()
	router.Add()

	done := make(chan bool)
	go func() {
		router.Wait()
		done <- true
	}()

	// Give goroutine time to start waiting
	time.Sleep(50 * time.Millisecond)

	// Wait should block until Done is called twice
	select {
	case <-done:
		t.Error("Wait() should not return before Done() is called")
	default:
		// Expected
	}

	router.Done()

	// Give time for potential completion
	time.Sleep(50 * time.Millisecond)

	select {
	case <-done:
		t.Error("Wait() should not return before all Done() calls")
	default:
		// Expected
	}

	router.Done()

	// Should complete now
	select {
	case <-done:
		// Expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Wait() should return after all Done() calls")
	}
}

func TestConsoleRouter_BaseWriter(t *testing.T) {
	router := &consoleRouter{
		wg: sync.WaitGroup{},
	}

	writer := router.BaseWriter()
	if writer == nil {
		t.Error("BaseWriter() returned nil")
	}

	// Should be able to write to it
	_, err := io.WriteString(writer, "test message\n")
	if err != nil {
		t.Errorf("Write to BaseWriter() failed: %v", err)
	}
}

func TestConsoleRouter_LineWriter(t *testing.T) {
	tests := []struct {
		name    string
		cmdName string
		prefix  string
		color   *color.Color
		line    string
	}{
		{
			name:    "with color",
			cmdName: "test",
			prefix:  "[test] ",
			color:   color.New(color.FgGreen),
			line:    "hello",
		},
		{
			name:    "without color",
			cmdName: "test",
			prefix:  "[test] ",
			color:   nil,
			line:    "hello",
		},
		{
			name:    "empty prefix",
			cmdName: "test",
			prefix:  "",
			color:   color.New(color.FgBlue),
			line:    "message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := &consoleRouter{
				wg: sync.WaitGroup{},
			}

			lineWriter := router.LineWriter(tt.cmdName, tt.color, tt.prefix)
			if lineWriter == nil {
				t.Fatal("LineWriter() returned nil")
			}

			// Just call it to ensure no panic
			lineWriter(tt.line)
		})
	}
}

func TestConsoleRouter_Stop(t *testing.T) {
	router := &consoleRouter{
		wg: sync.WaitGroup{},
	}

	// Stop should not panic and should be safe to call
	router.Stop()

	// Should be safe to call multiple times
	router.Stop()
}

func TestNewOutputRouter_Console(t *testing.T) {
	commands := []CommandConfig{
		{Name: "cmd1", Cmd: "echo"},
		{Name: "cmd2", Cmd: "ls"},
	}
	styles := defaultPanelStyles(commands)

	router, err := newOutputRouter(false, commands, styles)
	if err != nil {
		t.Fatalf("newOutputRouter() error = %v", err)
	}

	if router == nil {
		t.Fatal("newOutputRouter() returned nil router")
	}

	// Should be a console router
	if _, ok := router.(*consoleRouter); !ok {
		t.Error("expected consoleRouter type when enableTUI is false")
	}

	// Test basic operations
	router.Add()
	defer router.Done()

	writer := router.BaseWriter()
	if writer == nil {
		t.Error("BaseWriter() returned nil")
	}

	lineWriter := router.LineWriter("cmd1", color.New(color.FgCyan), "[test] ")
	if lineWriter == nil {
		t.Error("LineWriter() returned nil")
	}

	router.Stop()
}

func TestNewOutputRouter_TUI(t *testing.T) {
	t.Skip("TUI router requires interactive terminal, skipping in unit tests")

	commands := []CommandConfig{
		{Name: "cmd1", Cmd: "echo"},
	}
	styles := defaultPanelStyles(commands)

	router, err := newOutputRouter(true, commands, styles)
	if err != nil {
		t.Fatalf("newOutputRouter() error = %v", err)
	}

	if router == nil {
		t.Fatal("newOutputRouter() returned nil router")
	}

	// Should be a TUI router
	if _, ok := router.(*tuiRouter); !ok {
		t.Error("expected tuiRouter type when enableTUI is true")
	}

	router.Stop()
}

func TestConsoleRouterConcurrency(t *testing.T) {
	router := &consoleRouter{
		wg: sync.WaitGroup{},
	}

	// Test concurrent Add/Done operations
	concurrency := 10
	for i := 0; i < concurrency; i++ {
		router.Add()
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			router.Done()
		}()
	}

	done := make(chan bool)
	go func() {
		router.Wait()
		done <- true
	}()

	wg.Wait()

	select {
	case <-done:
		// Expected
	case <-time.After(500 * time.Millisecond):
		t.Error("Wait() should return after all concurrent Done() calls")
	}
}
