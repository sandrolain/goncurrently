#!/bin/bash
# Run the TUI example

cd "$(dirname "$0")/.." || exit 1

echo "Running TUI mode example..."
echo "This demonstrates the Terminal User Interface"
echo "Press Ctrl+C to stop"
echo ""

cat examples/tui-example.yaml | ./goncurrently
