#!/bin/bash
# Run the development servers example

cd "$(dirname "$0")/.." || exit 1

echo "Running development servers example..."
echo "This simulates a frontend (port 3000) and backend (port 8080) server"
echo "Press Ctrl+C to stop"
echo ""

cat examples/dev-servers.yaml | ./goncurrently
