#!/bin/bash
# Server simulator - simulates a server running for a duration

NAME=${1:-Server}
PORT=${2:-8080}
DURATION=${3:-30}

echo "[$NAME] Starting on port $PORT..."
echo "[$NAME] Ready to accept connections"

for i in $(seq 1 $DURATION); do
  echo "[$NAME] Handling request #$i"
  sleep 1
done

echo "[$NAME] Shutting down..."
