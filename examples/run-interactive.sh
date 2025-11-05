#!/bin/bash
# Interactive example selector

cd "$(dirname "$0")/.." || exit 1

echo "========================================="
echo "  goncurrently - Example Selector"
echo "========================================="
echo ""
echo "Select an example to run:"
echo ""
echo "1) Basic - Simple concurrent commands"
echo "2) Development Servers - Frontend + Backend"
echo "3) Watch & Build - File watching workflow"
echo "4) Microservices - Multiple services with setup"
echo "5) TUI Mode - Terminal User Interface"
echo "6) Timeouts - Duration limits and delays"
echo "7) Environment Variables - Custom env vars"
echo "8) Kill Others - Stop all when one exits"
echo "9) Restart - Auto-restart on failure"
echo "0) Exit"
echo ""
read -p "Enter choice [0-9]: " choice

case $choice in
  1)
    echo "Running: Basic Example"
    cat examples/basic.yaml | ./goncurrently
    ;;
  2)
    echo "Running: Development Servers"
    cat examples/dev-servers.yaml | ./goncurrently
    ;;
  3)
    echo "Running: Watch & Build"
    cat examples/watch-build.yaml | ./goncurrently
    ;;
  4)
    echo "Running: Microservices"
    cat examples/microservices.yaml | ./goncurrently
    ;;
  5)
    echo "Running: TUI Mode"
    cat examples/tui-example.yaml | ./goncurrently
    ;;
  6)
    echo "Running: Timeouts"
    cat examples/timeouts.yaml | ./goncurrently
    ;;
  7)
    echo "Running: Environment Variables"
    cat examples/env-vars.yaml | ./goncurrently
    ;;
  8)
    echo "Running: Kill Others"
    cat examples/kill-others.yaml | ./goncurrently
    ;;
  9)
    echo "Running: Restart"
    cat examples/restart.yaml | ./goncurrently
    ;;
  0)
    echo "Exiting..."
    exit 0
    ;;
  *)
    echo "Invalid choice!"
    exit 1
    ;;
esac
