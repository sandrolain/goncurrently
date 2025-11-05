package main

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStreamOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasFunc  bool
	}{
		{
			name:     "single line",
			input:    "hello world\n",
			expected: []string{"hello world"},
			hasFunc:  true,
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3\n",
			expected: []string{"line1", "line2", "line3"},
			hasFunc:  true,
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{},
			hasFunc:  true,
		},
		{
			name:     "nil function",
			input:    "test\n",
			expected: nil,
			hasFunc:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			var lines []string
			var mu sync.Mutex

			var writeLine func(string)
			if tt.hasFunc {
				writeLine = func(line string) {
					mu.Lock()
					lines = append(lines, line)
					mu.Unlock()
				}
			}

			streamOutput(writeLine, r)

			if tt.expected != nil {
				if len(lines) != len(tt.expected) {
					t.Errorf("expected %d lines, got %d", len(tt.expected), len(lines))
					return
				}
				for i, expected := range tt.expected {
					if lines[i] != expected {
						t.Errorf("line %d: expected '%s', got '%s'", i, expected, lines[i])
					}
				}
			}
		})
	}
}

func TestShouldRestart(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		timedOut     bool
		triesLeft    int
		restartTries int
		want         bool
		expectedLeft int
	}{
		{
			name:         "no error",
			err:          nil,
			timedOut:     false,
			triesLeft:    3,
			restartTries: 3,
			want:         false,
			expectedLeft: 3,
		},
		{
			name:         "timed out",
			err:          context.DeadlineExceeded,
			timedOut:     true,
			triesLeft:    3,
			restartTries: 3,
			want:         false,
			expectedLeft: 3,
		},
		{
			name:         "error with unlimited retries",
			err:          io.EOF,
			timedOut:     false,
			triesLeft:    0,
			restartTries: -1,
			want:         true,
			expectedLeft: 0,
		},
		{
			name:         "error with tries remaining",
			err:          io.EOF,
			timedOut:     false,
			triesLeft:    2,
			restartTries: 3,
			want:         true,
			expectedLeft: 1,
		},
		{
			name:         "error with no tries remaining",
			err:          io.EOF,
			timedOut:     false,
			triesLeft:    0,
			restartTries: 3,
			want:         false,
			expectedLeft: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triesLeft := tt.triesLeft
			got := shouldRestart(tt.err, tt.timedOut, &triesLeft, tt.restartTries)
			if got != tt.want {
				t.Errorf("shouldRestart() = %v, want %v", got, tt.want)
			}
			if triesLeft != tt.expectedLeft {
				t.Errorf("triesLeft = %d, want %d", triesLeft, tt.expectedLeft)
			}
		})
	}
}

func TestWaitStartDelay(t *testing.T) {
	tests := []struct {
		name       string
		startAfter string
		stopSignal bool
		want       bool
	}{
		{
			name:       "no delay",
			startAfter: "",
			stopSignal: false,
			want:       false,
		},
		{
			name:       "delay without stop",
			startAfter: "10ms",
			stopSignal: false,
			want:       false,
		},
		{
			name:       "delay with stop signal",
			startAfter: "1s",
			stopSignal: true,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommandConfig{
				Name:       "test",
				StartAfter: tt.startAfter,
			}

			var stopCh chan struct{}
			if tt.stopSignal {
				stopCh = make(chan struct{})
				close(stopCh)
			}

			got := waitStartDelay(c, stopCh)
			if got != tt.want {
				t.Errorf("waitStartDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWaitRestartDelay(t *testing.T) {
	tests := []struct {
		name         string
		restartAfter string
		stopSignal   bool
		want         bool
	}{
		{
			name:         "no delay",
			restartAfter: "",
			stopSignal:   false,
			want:         false,
		},
		{
			name:         "delay without stop",
			restartAfter: "10ms",
			stopSignal:   false,
			want:         false,
		},
		{
			name:         "delay with stop signal",
			restartAfter: "1s",
			stopSignal:   true,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommandConfig{
				Name:         "test",
				RestartAfter: tt.restartAfter,
			}

			var stopCh chan struct{}
			if tt.stopSignal {
				stopCh = make(chan struct{})
				close(stopCh)
			}

			got := waitRestartDelay(c, stopCh)
			if got != tt.want {
				t.Errorf("waitRestartDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartProcess(t *testing.T) {
	tests := []struct {
		name      string
		config    CommandConfig
		wantErr   bool
		checkStd  bool
		cmdExists bool
	}{
		{
			name: "simple command",
			config: CommandConfig{
				Name: "echo",
				Cmd:  "echo",
				Args: []string{"hello"},
			},
			wantErr:   false,
			checkStd:  true,
			cmdExists: true,
		},
		{
			name: "silent command",
			config: CommandConfig{
				Name:   "silent",
				Cmd:    "echo",
				Args:   []string{"test"},
				Silent: true,
			},
			wantErr:   false,
			checkStd:  false,
			cmdExists: true,
		},
		{
			name: "command with environment",
			config: CommandConfig{
				Name: "env",
				Cmd:  "echo",
				Args: []string{"test"},
				Env: map[string]string{
					"TEST_VAR": "value",
				},
			},
			wantErr:   false,
			checkStd:  true,
			cmdExists: true,
		},
		{
			name: "command with duration",
			config: CommandConfig{
				Name:     "timed",
				Cmd:      "sleep",
				Args:     []string{"0.01"},
				Duration: "100ms",
			},
			wantErr:   false,
			checkStd:  true,
			cmdExists: true,
		},
		{
			name: "nonexistent command",
			config: CommandConfig{
				Name: "invalid",
				Cmd:  "nonexistentcommand12345",
			},
			wantErr:   true,
			checkStd:  false,
			cmdExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, ctx, cancel, stdout, stderr, err := startProcess(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("startProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				defer func() {
					if cmd != nil && cmd.Process != nil {
						_ = cmd.Process.Kill()
					}
					if cancel != nil {
						cancel()
					}
				}()

				if tt.checkStd {
					if stdout == nil {
						t.Error("expected stdout to be non-nil")
					}
					if stderr == nil {
						t.Error("expected stderr to be non-nil")
					}
				} else if tt.config.Silent {
					if stdout != nil {
						t.Error("expected stdout to be nil for silent command")
					}
					if stderr != nil {
						t.Error("expected stderr to be nil for silent command")
					}
				}

				if tt.config.Duration != "" && ctx == nil {
					t.Error("expected context to be non-nil for command with duration")
				}
			}
		})
	}
}

func TestExecuteOnce(t *testing.T) {
	tests := []struct {
		name         string
		config       CommandConfig
		sendStop     bool
		wantTimedOut bool
		wantStopped  bool
	}{
		{
			name: "successful execution",
			config: CommandConfig{
				Name: "success",
				Cmd:  "echo",
				Args: []string{"test"},
			},
			sendStop:     false,
			wantTimedOut: false,
			wantStopped:  false,
		},
		{
			name: "command with timeout",
			config: CommandConfig{
				Name:     "timeout",
				Cmd:      "sleep",
				Args:     []string{"10"},
				Duration: "50ms",
			},
			sendStop:     false,
			wantTimedOut: true,
			wantStopped:  false,
		},
		{
			name: "interrupted command",
			config: CommandConfig{
				Name: "interrupted",
				Cmd:  "sleep",
				Args: []string{"10"},
			},
			sendStop:     true,
			wantTimedOut: false,
			wantStopped:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopCh := make(chan struct{})
			immediateCh := make(chan struct{})
			signals := stopSignals{
				stop:      stopCh,
				immediate: immediateCh,
			}

			var output bytes.Buffer
			writeFunc := func(line string) {
				output.WriteString(line + "\n")
			}

			done := make(chan bool, 1)
			go func() {
				if tt.sendStop {
					time.Sleep(50 * time.Millisecond)
					close(stopCh)
				}
			}()

			go func() {
				timedOut, interrupted, _ := executeOnce(tt.config, "[test] ", writeFunc, writeFunc, signals, 100*time.Millisecond)
				done <- timedOut
				if interrupted != tt.wantStopped {
					t.Errorf("interrupted = %v, wantStopped %v", interrupted, tt.wantStopped)
				}
			}()

			select {
			case timedOut := <-done:
				if timedOut != tt.wantTimedOut {
					t.Errorf("timedOut = %v, wantTimedOut %v", timedOut, tt.wantTimedOut)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("test timed out")
			}
		})
	}
}

func TestLogCommandLine(t *testing.T) {
	tests := []struct {
		name          string
		hasStdout     bool
		hasStderr     bool
		identifier    string
		message       string
		expectedInOut bool
		expectedInErr bool
	}{
		{
			name:          "with stdout writer",
			hasStdout:     true,
			hasStderr:     false,
			identifier:    "[test] ",
			message:       "hello",
			expectedInOut: true,
			expectedInErr: false,
		},
		{
			name:          "with stderr writer",
			hasStdout:     false,
			hasStderr:     true,
			identifier:    "[test] ",
			message:       "error",
			expectedInOut: false,
			expectedInErr: true,
		},
		{
			name:          "with both writers",
			hasStdout:     true,
			hasStderr:     true,
			identifier:    "[test] ",
			message:       "message",
			expectedInOut: false,
			expectedInErr: true,
		},
		{
			name:          "with no writers",
			hasStdout:     false,
			hasStderr:     false,
			identifier:    "[test] ",
			message:       "fallback",
			expectedInOut: false,
			expectedInErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdoutBuf, stderrBuf bytes.Buffer
			var stdoutWriter, stderrWriter func(string)

			if tt.hasStdout {
				stdoutWriter = func(s string) {
					stdoutBuf.WriteString(s)
				}
			}
			if tt.hasStderr {
				stderrWriter = func(s string) {
					stderrBuf.WriteString(s)
				}
			}

			// Save original errorOutput
			origErrorOutput := errorOutput
			var errBuf bytes.Buffer
			errorOutput = &errBuf
			defer func() {
				errorOutput = origErrorOutput
			}()

			logCommandLine(stdoutWriter, stderrWriter, tt.identifier, tt.message)

			if tt.expectedInOut && !strings.Contains(stdoutBuf.String(), tt.message) {
				t.Error("expected message in stdout")
			}
			if tt.expectedInErr && !strings.Contains(stderrBuf.String(), tt.message) {
				t.Error("expected message in stderr")
			}
		})
	}
}

func TestTerminateProcess(t *testing.T) {
	t.Run("immediate kill with zero timeout", func(t *testing.T) {
		cmd := exec.Command("sleep", "10")
		err := cmd.Start()
		if err != nil {
			t.Fatalf("failed to start command: %v", err)
		}

		done := make(chan error, 1)
		immediate := make(chan struct{})

		go func() {
			done <- cmd.Wait()
		}()

		terminateProcess(cmd, 0, done, immediate)

		// Ensure cleanup
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})
}

func TestRunSetupWithRetries(t *testing.T) {
	tests := []struct {
		name       string
		config     CommandConfig
		shouldPass bool
	}{
		{
			name: "successful setup",
			config: CommandConfig{
				Name: "setup",
				Cmd:  "echo",
				Args: []string{"test"},
			},
			shouldPass: true,
		},
		{
			name: "failing setup with no retries",
			config: CommandConfig{
				Name:         "fail",
				Cmd:          "false",
				RestartTries: 0,
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			writeFunc := func(line string) {
				output.WriteString(line + "\n")
			}

			result := runSetupWithRetries(tt.config, "[test] ", writeFunc, writeFunc)
			if result != tt.shouldPass {
				t.Errorf("runSetupWithRetries() = %v, want %v", result, tt.shouldPass)
			}
		})
	}
}
