#!/bin/bash
# Run the environment variables example

cd "$(dirname "$0")/.." || exit 1

echo "Running environment variables example..."
echo "This shows how to set custom env vars per command"
echo "Press Ctrl+C to stop"
echo ""

cat examples/env-vars.yaml | ./goncurrently
