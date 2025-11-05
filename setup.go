package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func runSetupSequence(cmds []CommandConfig, colors []*color.Color, sink outputRouter) {
	for i, c := range cmds {
		col := colors[i%len(colors)]
		identifier := fmt.Sprintf("[setup:%s] ", c.Name)
		stdoutWriter := sink.LineWriter(basePanelName, col, identifier)
		stderrWriter := sink.LineWriter(basePanelName, col, fmt.Sprintf("[setup:%s stderr] ", c.Name))
		if d := mustParseDurationField("startAfter", c.StartAfter, c.Name); d > 0 {
			time.Sleep(d)
		}
		baseLog("[setup:%s] starting", c.Name)
		if !runSetupWithRetries(c, identifier, stdoutWriter, stderrWriter) {
			color.New(color.FgRed, color.Bold).Fprintf(errorOutput, "Setup command '%s' failed after retries\n", c.Name) //nolint:errcheck
			os.Exit(1)
		}
		baseLog("[setup:%s] completed", c.Name)
	}
}
