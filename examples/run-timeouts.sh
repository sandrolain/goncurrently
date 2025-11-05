#!/bin/bash
# Run the timeouts example

cd "$(dirname "$0")/.." || exit 1

echo "Running timeouts and delays example..."
echo "This demonstrates duration limits and startup delays"
echo "Press Ctrl+C to stop"
echo ""

cat examples/timeouts.yaml | ./goncurrently
