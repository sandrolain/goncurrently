package main

import (
	"bytes"
	"testing"
)

func TestBaseLog(t *testing.T) {
	// Save original errorOutput
	origErrorOutput := errorOutput
	defer func() {
		errorOutput = origErrorOutput
	}()

	var buf bytes.Buffer
	errorOutput = &buf

	tests := []struct {
		name   string
		format string
		args   []interface{}
		want   string
	}{
		{
			name:   "simple message",
			format: "test message",
			args:   nil,
			want:   "test message\n",
		},
		{
			name:   "formatted message",
			format: "count: %d",
			args:   []interface{}{42},
			want:   "count: 42\n",
		},
		{
			name:   "message with newline",
			format: "test\n",
			args:   nil,
			want:   "test\n",
		},
		{
			name:   "multiple args",
			format: "%s = %d",
			args:   []interface{}{"value", 123},
			want:   "value = 123\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			baseLog(tt.format, tt.args...)
			got := buf.String()
			if got != tt.want {
				t.Errorf("baseLog() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if basePanelName != "goncurrently" {
		t.Errorf("basePanelName = %q, want %q", basePanelName, "goncurrently")
	}
	if lineJoinFormat != "%s%s\n" {
		t.Errorf("lineJoinFormat = %q, want %q", lineJoinFormat, "%s%s\n")
	}
}
