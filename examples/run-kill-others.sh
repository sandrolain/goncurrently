#!/bin/bash
# Run the kill others example

cd "$(dirname "$0")/.." || exit 1

echo "Running kill others example..."
echo "All processes will stop when the first one completes"
echo "Press Ctrl+C to stop"
echo ""

cat examples/kill-others.yaml | ./goncurrently
