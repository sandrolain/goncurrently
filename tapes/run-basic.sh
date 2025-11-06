#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
cat examples/basic.yaml | timeout 3 ./goncurrently-linux || true
