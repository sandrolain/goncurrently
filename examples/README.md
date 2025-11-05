# Examples

This directory contains practical examples of `goncurrently` configurations that you can try immediately.

## Interactive Selector

For an interactive menu to select and run examples:

```bash
./examples/run-interactive.sh
```

## Quick Start

Each example can be run with:

```bash
cat examples/<example-name>.yaml | goncurrently
```

Or using the provided run scripts:

```bash
./examples/run-<example-name>.sh
```

## Available Examples

### 1. Basic Example (`basic.yaml`)

Simple example running multiple echo commands concurrently.

```bash
cat examples/basic.yaml | goncurrently
```

### 2. Development Servers (`dev-servers.yaml`)

Simulates a frontend and backend development environment.

```bash
cat examples/dev-servers.yaml | goncurrently
```

### 3. Watch and Build (`watch-build.yaml`)

Example of file watching and building processes.

```bash
cat examples/watch-build.yaml | goncurrently
```

### 4. Microservices (`microservices.yaml`)

Example with setup commands and multiple services with restart logic.

```bash
cat examples/microservices.yaml | goncurrently
```

### 5. TUI Mode (`tui-example.yaml`)

Example with Terminal UI enabled for better visualization.

```bash
cat examples/tui-example.yaml | goncurrently
```

### 6. Timeouts and Delays (`timeouts.yaml`)

Example demonstrating duration limits and startup delays.

```bash
cat examples/timeouts.yaml | goncurrently
```

### 7. Environment Variables (`env-vars.yaml`)

Example showing how to set environment variables per command.

```bash
cat examples/env-vars.yaml | goncurrently
```

### 8. Kill Others (`kill-others.yaml`)

Example demonstrating the killOthers behavior.

```bash
cat examples/kill-others.yaml | goncurrently
```

## Helper Scripts

The `scripts/` subdirectory contains helper scripts used by the examples:

- `counter.sh` - Counts to a number with delays
- `random-fail.sh` - Randomly succeeds or fails (for testing restarts)
- `server-sim.sh` - Simulates a server that runs for a duration
- `watcher-sim.sh` - Simulates a file watcher

## Testing Examples

To test that all examples work:

```bash
./examples/test-all.sh
```
