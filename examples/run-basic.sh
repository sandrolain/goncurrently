#!/bin/bash
# Run the basic example

cd "$(dirname "$0")/.." || exit 1

echo "Running basic example..."
echo "Press Ctrl+C to stop"
echo ""

cat examples/basic.yaml | ./goncurrently
