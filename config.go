package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// CommandConfig describes an individual command to run either during setup or main execution.
type CommandConfig struct {
	Name         string            `yaml:"name"`
	Cmd          string            `yaml:"cmd" validate:"required"`
	Args         []string          `yaml:"args"`
	RestartTries int               `yaml:"restartTries"`
	RestartAfter string            `yaml:"restartAfter"`
	Env          map[string]string `yaml:"env"`
	StartAfter   string            `yaml:"startAfter"`
	Silent       bool              `yaml:"silent"`
	Duration     string            `yaml:"duration"`
}

// Config aggregates the complete execution plan for the tool.
type Config struct {
	Commands      []CommandConfig `yaml:"commands" validate:"required,dive,required"`
	KillOthers    bool            `yaml:"killOthers"`
	KillTimeout   int             `yaml:"killTimeout"`
	NoColors      bool            `yaml:"noColors"`
	SetupCommands []CommandConfig `yaml:"setupCommands"`
	EnableTUI     bool            `yaml:"enableTUI"`
}

// loadConfig fully reads configuration data from the provided reader.
func loadConfig(r io.Reader) (Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// assignNames fills missing command names with the executable basename.
func assignNames(cmds []CommandConfig) {
	for i := range cmds {
		if cmds[i].Name != "" {
			continue
		}
		cmd := cmds[i].Cmd
		if idx := strings.LastIndex(cmd, "/"); idx >= 0 && idx+1 < len(cmd) {
			cmds[i].Name = cmd[idx+1:]
			continue
		}
		cmds[i].Name = cmd
	}
}

// mustParseDurationField parses duration fields, terminating the process on invalid values.
func mustParseDurationField(field string, value string, commandName string) time.Duration {
	if value == "" {
		return 0
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		fmt.Fprintf(errorOutput, "invalid duration for %s in command '%s': %v\n", field, commandName, err) //nolint:errcheck
		os.Exit(1)
	}
	return d
}
