#!/bin/bash
# Run the watch and build example

cd "$(dirname "$0")/.." || exit 1

echo "Running watch and build example..."
echo "This simulates a file watcher and development server"
echo "Press Ctrl+C to stop"
echo ""

cat examples/watch-build.yaml | ./goncurrently
