package main

import (
	"testing"

	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
)

func TestDefaultCommandColors(t *testing.T) {
	colors := defaultCommandColors()

	if len(colors) == 0 {
		t.Error("expected non-empty color slice")
	}

	expectedCount := 6
	if len(colors) != expectedCount {
		t.Errorf("expected %d colors, got %d", expectedCount, len(colors))
	}

	for i, c := range colors {
		if c == nil {
			t.Errorf("color at index %d is nil", i)
		}
	}
}

func TestDefaultPanelStyles(t *testing.T) {
	tests := []struct {
		name     string
		commands []CommandConfig
		expected int
	}{
		{
			name: "single command",
			commands: []CommandConfig{
				{Name: "cmd1"},
			},
			expected: 2, // cmd1 + basePanelName
		},
		{
			name: "multiple commands",
			commands: []CommandConfig{
				{Name: "cmd1"},
				{Name: "cmd2"},
				{Name: "cmd3"},
			},
			expected: 4, // 3 commands + basePanelName
		},
		{
			name:     "no commands",
			commands: []CommandConfig{},
			expected: 1, // just basePanelName
		},
		{
			name: "command with basePanelName",
			commands: []CommandConfig{
				{Name: basePanelName},
			},
			expected: 1, // basePanelName only once
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			styles := defaultPanelStyles(tt.commands)

			if len(styles) != tt.expected {
				t.Errorf("expected %d styles, got %d", tt.expected, len(styles))
			}

			// Check that basePanelName is always present
			if _, ok := styles[basePanelName]; !ok {
				t.Error("basePanelName style not found")
			}

			// Check that each command has a style
			for _, cmd := range tt.commands {
				if _, ok := styles[cmd.Name]; !ok {
					t.Errorf("style for command '%s' not found", cmd.Name)
				}
			}
		})
	}
}

func TestDefaultPanelStylesColors(t *testing.T) {
	commands := []CommandConfig{
		{Name: "cmd1"},
		{Name: "cmd2"},
	}

	styles := defaultPanelStyles(commands)

	for name, style := range styles {
		if name == basePanelName {
			// Base panel should have specific colors
			if style.BorderColor != tcell.ColorDarkGray {
				t.Errorf("basePanelName border color incorrect")
			}
			if style.TitleColor != tcell.ColorWhite {
				t.Errorf("basePanelName title color incorrect")
			}
		} else {
			// Command panels should have non-default colors
			if style.BorderColor == tcell.ColorDefault {
				t.Errorf("command panel '%s' has default border color", name)
			}
			if style.TitleColor == tcell.ColorDefault {
				t.Errorf("command panel '%s' has default title color", name)
			}
		}

		// All panels should have default background
		if style.BackgroundColor != tcell.ColorDefault {
			t.Errorf("panel '%s' background color should be default", name)
		}
	}
}

func TestDefaultPanelStylesCycling(t *testing.T) {
	// Test that colors cycle correctly when there are more commands than colors
	commands := make([]CommandConfig, 8)
	for i := range commands {
		commands[i] = CommandConfig{Name: "cmd" + string(rune('1'+i))}
	}

	styles := defaultPanelStyles(commands)

	// Should have 8 commands + basePanelName
	if len(styles) != 9 {
		t.Errorf("expected 9 styles, got %d", len(styles))
	}

	// Verify all commands have styles
	for _, cmd := range commands {
		if _, ok := styles[cmd.Name]; !ok {
			t.Errorf("style for command '%s' not found", cmd.Name)
		}
	}
}

func TestCommandColorsAreValid(t *testing.T) {
	colors := defaultCommandColors()

	expectedColors := []*color.Color{
		color.New(color.FgCyan),
		color.New(color.FgGreen),
		color.New(color.FgMagenta),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgRed),
	}

	if len(colors) != len(expectedColors) {
		t.Errorf("expected %d colors, got %d", len(expectedColors), len(colors))
	}

	for i := range colors {
		if colors[i] == nil {
			t.Errorf("color at index %d is nil", i)
		}
	}
}
