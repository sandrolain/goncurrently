#!/bin/bash
# Run the shutdown commands example

cd "$(dirname "$0")/.." || exit 1

echo "Running shutdown commands example..."
echo "This demonstrates cleanup tasks after main commands complete"
echo "Press Ctrl+C to stop"
echo ""

cat examples/shutdown.yaml | ./goncurrently
