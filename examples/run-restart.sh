#!/bin/bash
# Run the restart example

cd "$(dirname "$0")/.." || exit 1

echo "Running restart example..."
echo "This demonstrates automatic restart on failure"
echo "Press Ctrl+C to stop"
echo ""

cat examples/restart.yaml | ./goncurrently
