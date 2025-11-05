#!/bin/bash
# Counter script - counts from 1 to N with delays

COUNT=${1:-5}
DELAY=${2:-1}

echo "Starting counter (counting to $COUNT with ${DELAY}s delay)..."

for i in $(seq 1 $COUNT); do
  echo "Count: $i"
  sleep $DELAY
done

echo "Counter finished!"
