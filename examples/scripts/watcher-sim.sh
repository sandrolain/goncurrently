#!/bin/bash
# Watcher simulator - simulates a file watcher

DIR=${1:-.}

echo "Watching directory: $DIR"
echo "Watcher started..."

while true; do
  RANDOM_FILE="file-$((RANDOM % 10)).txt"
  echo "Detected change in: $RANDOM_FILE"
  echo "Rebuilding..."
  sleep 2
  echo "Build complete!"
  sleep $((RANDOM % 5 + 3))
done
