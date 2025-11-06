#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
# Use script to create a PTY for proper TUI rendering
# -q = quiet mode, -c = command to run, /dev/null = don't save output to file
script -q -c "cat examples/tui-example.yaml | ./goncurrently-linux" /dev/null
