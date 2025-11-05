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
