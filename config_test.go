package main

import (
	"bytes"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid basic config",
			yaml: `
commands:
  - cmd: echo
    args: ["hello"]
`,
			wantErr: false,
		},
		{
			name: "valid config with setup commands",
			yaml: `
commands:
  - cmd: echo
    args: ["test"]
setupCommands:
  - cmd: echo
    args: ["setup"]
killOthers: true
killTimeout: 5000
noColors: true
enableTUI: false
`,
			wantErr: false,
		},
		{
			name: "valid config with all command fields",
			yaml: `
commands:
  - name: mycommand
    cmd: sleep
    args: ["1"]
    restartTries: 3
    restartAfter: "1s"
    env:
      FOO: bar
    startAfter: "500ms"
    silent: true
    duration: "10s"
`,
			wantErr: false,
		},
		{
			name:    "empty config",
			yaml:    ``,
			wantErr: true,
		},
		{
			name: "invalid yaml",
			yaml: `
commands:
  - cmd: echo
	args: invalid
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.yaml)
			cfg, err := loadConfig(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(cfg.Commands) == 0 {
					t.Error("loadConfig() returned config with no commands")
				}
			}
		})
	}
}

func TestLoadConfigFields(t *testing.T) {
	yaml := `
commands:
  - name: test
    cmd: echo
    args: ["hello", "world"]
    restartTries: 5
    restartAfter: "2s"
    env:
      KEY1: value1
      KEY2: value2
    startAfter: "1s"
    silent: true
    duration: "30s"
setupCommands:
  - cmd: setup
    args: ["init"]
shutdownCommands:
  - cmd: cleanup
    args: ["--all"]
killOthers: true
killTimeout: 3000
noColors: true
enableTUI: true
`
	r := bytes.NewBufferString(yaml)
	cfg, err := loadConfig(r)
	if err != nil {
		t.Fatalf("loadConfig() unexpected error: %v", err)
	}

	if len(cfg.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(cfg.Commands))
	}

	cmd := cfg.Commands[0]
	if cmd.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", cmd.Name)
	}
	if cmd.Cmd != "echo" {
		t.Errorf("expected cmd 'echo', got '%s'", cmd.Cmd)
	}
	if len(cmd.Args) != 2 || cmd.Args[0] != "hello" || cmd.Args[1] != "world" {
		t.Errorf("expected args ['hello', 'world'], got %v", cmd.Args)
	}
	if cmd.RestartTries != 5 {
		t.Errorf("expected restartTries 5, got %d", cmd.RestartTries)
	}
	if cmd.RestartAfter != "2s" {
		t.Errorf("expected restartAfter '2s', got '%s'", cmd.RestartAfter)
	}
	if len(cmd.Env) != 2 || cmd.Env["KEY1"] != "value1" {
		t.Errorf("expected env with KEY1=value1, got %v", cmd.Env)
	}
	if cmd.StartAfter != "1s" {
		t.Errorf("expected startAfter '1s', got '%s'", cmd.StartAfter)
	}
	if !cmd.Silent {
		t.Error("expected silent to be true")
	}
	if cmd.Duration != "30s" {
		t.Errorf("expected duration '30s', got '%s'", cmd.Duration)
	}

	if len(cfg.SetupCommands) != 1 {
		t.Errorf("expected 1 setup command, got %d", len(cfg.SetupCommands))
	}
	if len(cfg.ShutdownCommands) != 1 {
		t.Errorf("expected 1 shutdown command, got %d", len(cfg.ShutdownCommands))
	}
	if cfg.ShutdownCommands[0].Cmd != "cleanup" {
		t.Errorf("expected shutdown command 'cleanup', got '%s'", cfg.ShutdownCommands[0].Cmd)
	}
	if !cfg.KillOthers {
		t.Error("expected killOthers to be true")
	}
	if cfg.KillTimeout != 3000 {
		t.Errorf("expected killTimeout 3000, got %d", cfg.KillTimeout)
	}
	if !cfg.NoColors {
		t.Error("expected noColors to be true")
	}
	if !cfg.EnableTUI {
		t.Error("expected enableTUI to be true")
	}
}

func TestAssignNames(t *testing.T) {
	tests := []struct {
		name     string
		commands []CommandConfig
		expected []string
	}{
		{
			name: "commands with names",
			commands: []CommandConfig{
				{Name: "cmd1", Cmd: "echo"},
				{Name: "cmd2", Cmd: "ls"},
			},
			expected: []string{"cmd1", "cmd2"},
		},
		{
			name: "commands without names",
			commands: []CommandConfig{
				{Cmd: "echo"},
				{Cmd: "ls"},
			},
			expected: []string{"echo", "ls"},
		},
		{
			name: "commands with paths",
			commands: []CommandConfig{
				{Cmd: "/usr/bin/echo"},
				{Cmd: "/bin/ls"},
			},
			expected: []string{"echo", "ls"},
		},
		{
			name: "mixed commands",
			commands: []CommandConfig{
				{Name: "named", Cmd: "/usr/bin/echo"},
				{Cmd: "/bin/ls"},
				{Cmd: "pwd"},
			},
			expected: []string{"named", "ls", "pwd"},
		},
		{
			name:     "empty slice",
			commands: []CommandConfig{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignNames(tt.commands)
			for i, cmd := range tt.commands {
				if i < len(tt.expected) && cmd.Name != tt.expected[i] {
					t.Errorf("command %d: expected name '%s', got '%s'", i, tt.expected[i], cmd.Name)
				}
			}
		})
	}
}

func TestMustParseDurationField(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		value       string
		commandName string
		want        time.Duration
		shouldPanic bool
	}{
		{
			name:        "valid duration seconds",
			field:       "startAfter",
			value:       "5s",
			commandName: "test",
			want:        5 * time.Second,
			shouldPanic: false,
		},
		{
			name:        "valid duration milliseconds",
			field:       "restartAfter",
			value:       "500ms",
			commandName: "test",
			want:        500 * time.Millisecond,
			shouldPanic: false,
		},
		{
			name:        "valid duration minutes",
			field:       "duration",
			value:       "2m",
			commandName: "test",
			want:        2 * time.Minute,
			shouldPanic: false,
		},
		{
			name:        "empty value",
			field:       "startAfter",
			value:       "",
			commandName: "test",
			want:        0,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParseDurationField(tt.field, tt.value, tt.commandName)
			if got != tt.want {
				t.Errorf("mustParseDurationField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustParseDurationFieldInvalid(t *testing.T) {
	// This test is skipped because mustParseDurationField calls os.Exit(1)
	// which cannot be easily tested in unit tests
	t.Skip("Skipping test that calls os.Exit")
}

func TestAssignNamesPreservesExisting(t *testing.T) {
	commands := []CommandConfig{
		{Name: "custom", Cmd: "/usr/bin/foo"},
	}
	assignNames(commands)
	if commands[0].Name != "custom" {
		t.Errorf("expected name to remain 'custom', got '%s'", commands[0].Name)
	}
}

func TestLoadConfigWithMultipleCommands(t *testing.T) {
	yaml := `
commands:
  - cmd: echo
    args: ["first"]
  - cmd: echo
    args: ["second"]
  - cmd: echo
    args: ["third"]
`
	r := bytes.NewBufferString(yaml)
	cfg, err := loadConfig(r)
	if err != nil {
		t.Fatalf("loadConfig() unexpected error: %v", err)
	}

	if len(cfg.Commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(cfg.Commands))
	}
}

func TestConfigDefaults(t *testing.T) {
	yaml := `
commands:
  - cmd: echo
`
	r := bytes.NewBufferString(yaml)
	cfg, err := loadConfig(r)
	if err != nil {
		t.Fatalf("loadConfig() unexpected error: %v", err)
	}

	if cfg.KillOthers != false {
		t.Error("expected killOthers default to be false")
	}
	if cfg.KillTimeout != 0 {
		t.Error("expected killTimeout default to be 0")
	}
	if cfg.NoColors != false {
		t.Error("expected noColors default to be false")
	}
	if cfg.EnableTUI != false {
		t.Error("expected enableTUI default to be false")
	}
}
