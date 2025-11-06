#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
cat examples/dev-servers.yaml | timeout 8 ./goncurrently-linux || true
