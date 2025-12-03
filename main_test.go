package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if !strings.HasPrefix(Version, "v") {
		t.Errorf("Version should start with 'v', got: %s", Version)
	}
}

func TestPrintVersion(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "goncurrently") {
		t.Errorf("printVersion output should contain 'goncurrently', got: %s", output)
	}
	if !strings.Contains(output, Version) {
		t.Errorf("printVersion output should contain version '%s', got: %s", Version, output)
	}
}

func TestPrintHelp(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expectedStrings := []string{
		"goncurrently",
		"Usage:",
		"--help",
		"--version",
		"commands",
		"setupCommands",
		"shutdownCommands",
		"killOthers",
		"enableTUI",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printHelp output should contain '%s'", expected)
		}
	}
}

func TestCLIHelpFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "goncurrently_test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("goncurrently_test")

	tests := []struct {
		name string
		arg  string
	}{
		{"help flag", "--help"},
		{"short help flag", "-h"},
		{"help command", "help"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./goncurrently_test", tt.arg) // #nosec G204 -- test code with controlled input
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			if !strings.Contains(string(output), "Usage:") {
				t.Errorf("Expected help output, got: %s", string(output))
			}
		})
	}
}

func TestCLIVersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "goncurrently_test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("goncurrently_test")

	tests := []struct {
		name string
		arg  string
	}{
		{"version flag", "--version"},
		{"short version flag", "-v"},
		{"version command", "version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./goncurrently_test", tt.arg) // #nosec G204 -- test code with controlled input
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			if !strings.Contains(string(output), Version) {
				t.Errorf("Expected version in output, got: %s", string(output))
			}
		})
	}
}

func TestCLIUnknownFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "goncurrently_test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("goncurrently_test")

	cmd = exec.Command("./goncurrently_test", "--unknown")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected command to fail with unknown flag")
	}

	if !strings.Contains(string(output), "Unknown option") {
		t.Errorf("Expected 'Unknown option' in output, got: %s", string(output))
	}
}
