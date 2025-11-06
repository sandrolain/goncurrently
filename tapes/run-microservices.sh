#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
cat examples/microservices.yaml | timeout 10 ./goncurrently-linux || true
