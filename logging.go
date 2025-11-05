package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	basePanelName  = "goncurrently"
	lineJoinFormat = "%s%s\n"
)

var errorOutput io.Writer = os.Stderr

func baseLog(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(errorOutput, format, args...) //nolint:errcheck
}
