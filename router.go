package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
)

type outputRouter interface {
	BaseWriter() io.Writer
	LineWriter(name string, col *color.Color, prefix string) func(string)
	Stop()
	Wait()
	Add()
	Done()
}

type panelAppearance struct {
	BorderColor     tcell.Color
	TitleColor      tcell.Color
	BackgroundColor tcell.Color
}

func newOutputRouter(enableTUI bool, commands []CommandConfig, styles map[string]panelAppearance) (outputRouter, error) {
	if !enableTUI {
		return &consoleRouter{
			wg: sync.WaitGroup{},
		}, nil
	}
	commandNames := make([]string, 0, len(commands))
	for _, c := range commands {
		commandNames = append(commandNames, c.Name)
	}
	return newTUIRouter(basePanelName, commandNames, styles)
}

type consoleRouter struct {
	wg sync.WaitGroup
}

func (c *consoleRouter) Add() {
	c.wg.Add(1)
}

func (c *consoleRouter) Done() {
	c.wg.Done()
}

func (c *consoleRouter) Wait() {
	c.wg.Wait()
}

func (c *consoleRouter) BaseWriter() io.Writer {
	return os.Stderr
}

func (c *consoleRouter) LineWriter(_ string, col *color.Color, prefix string) func(string) {
	coloredPrefix := prefix
	if col != nil {
		coloredPrefix = col.Sprint(prefix)
	}
	return func(line string) {
		fmt.Printf(lineJoinFormat, coloredPrefix, line) //nolint:forbidigo
	}
}

// Stop implements outputRouter for console routing and requires no cleanup.
func (c *consoleRouter) Stop() {
	// No cleanup required when writing directly to the console.
}
