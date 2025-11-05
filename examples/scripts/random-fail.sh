#!/bin/bash
# Random fail script - randomly succeeds or fails

FAIL_RATE=${1:-50}

echo "Starting with $FAIL_RATE% fail rate..."
sleep 1

RANDOM_NUM=$((RANDOM % 100))

if [ $RANDOM_NUM -lt $FAIL_RATE ]; then
  echo "ERROR: Random failure occurred!"
  exit 1
else
  echo "SUCCESS: Completed successfully"
  exit 0
fi
