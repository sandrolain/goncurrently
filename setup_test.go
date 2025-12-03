package main

import (
	"bytes"
	"testing"
)

func TestRunSetupSequence(t *testing.T) {
	// Save original errorOutput
	origErrorOutput := errorOutput
	defer func() {
		errorOutput = origErrorOutput
	}()

	var buf bytes.Buffer
	errorOutput = &buf

	tests := []struct {
		name     string
		commands []CommandConfig
		wantErr  bool
	}{
		{
			name:     "empty setup commands",
			commands: []CommandConfig{},
			wantErr:  false,
		},
		{
			name: "single successful command",
			commands: []CommandConfig{
				{
					Name: "test",
					Cmd:  "echo",
					Args: []string{"hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple commands",
			commands: []CommandConfig{
				{
					Name: "first",
					Cmd:  "echo",
					Args: []string{"first"},
				},
				{
					Name: "second",
					Cmd:  "echo",
					Args: []string{"second"},
				},
			},
			wantErr: false,
		},
		{
			name: "command with start delay",
			commands: []CommandConfig{
				{
					Name:       "delayed",
					Cmd:        "echo",
					Args:       []string{"test"},
					StartAfter: "10ms",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			colors := defaultCommandColors()
			router := &consoleRouter{}

			// This will call os.Exit(1) on failure, so we can't easily test failure cases
			if !tt.wantErr {
				runSetupSequence(tt.commands, colors, router)
			}
		})
	}
}

func TestRunShutdownSequence(t *testing.T) {
	// Save original errorOutput
	origErrorOutput := errorOutput
	defer func() {
		errorOutput = origErrorOutput
	}()

	var buf bytes.Buffer
	errorOutput = &buf

	tests := []struct {
		name     string
		commands []CommandConfig
	}{
		{
			name:     "empty shutdown commands",
			commands: []CommandConfig{},
		},
		{
			name: "single successful command",
			commands: []CommandConfig{
				{
					Name: "cleanup",
					Cmd:  "echo",
					Args: []string{"cleaning up"},
				},
			},
		},
		{
			name: "multiple commands",
			commands: []CommandConfig{
				{
					Name: "cleanup-db",
					Cmd:  "echo",
					Args: []string{"cleaning database"},
				},
				{
					Name: "cleanup-cache",
					Cmd:  "echo",
					Args: []string{"cleaning cache"},
				},
			},
		},
		{
			name: "command with start delay",
			commands: []CommandConfig{
				{
					Name:       "delayed-cleanup",
					Cmd:        "echo",
					Args:       []string{"delayed cleanup"},
					StartAfter: "10ms",
				},
			},
		},
		{
			name: "failing command continues to next",
			commands: []CommandConfig{
				{
					Name: "fail-cmd",
					Cmd:  "false",
				},
				{
					Name: "success-cmd",
					Cmd:  "echo",
					Args: []string{"this should run"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			colors := defaultCommandColors()
			router := &consoleRouter{}

			// runShutdownSequence should not panic or exit on failures
			runShutdownSequence(tt.commands, colors, router)
		})
	}
}
