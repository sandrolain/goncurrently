#!/bin/bash
# Run the microservices example

cd "$(dirname "$0")/.." || exit 1

echo "Running microservices example..."
echo "This simulates multiple microservices with setup commands"
echo "Press Ctrl+C to stop"
echo ""

cat examples/microservices.yaml | ./goncurrently
