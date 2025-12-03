# Quick Reference - goncurrently Examples

## Running Examples

All examples assume you've built the tool first:

```bash
go build -o goncurrently .
```

## Examples Overview

| Example | Description | Command | Features Demonstrated |
|---------|-------------|---------|----------------------|
| **Basic** | Simple concurrent commands | `./examples/run-basic.sh` | Basic concurrency |
| **Dev Servers** | Frontend + Backend | `./examples/run-dev-servers.sh` | Restart, env vars, delays |
| **Watch Build** | File watching workflow | `./examples/run-watch-build.sh` | Long-running processes |
| **Microservices** | Multiple services | `./examples/run-microservices.sh` | Setup commands, orchestration |
| **TUI Mode** | Terminal UI | `./examples/run-tui.sh` | TUI visualization |
| **Timeouts** | Duration limits | `./examples/run-timeouts.sh` | Timeouts, delays |
| **Env Vars** | Environment variables | `./examples/run-env-vars.sh` | Custom env per command |
| **Kill Others** | Stop all on exit | `./examples/run-kill-others.sh` | killOthers behavior |
| **Restart** | Auto-restart on failure | `./examples/run-restart.sh` | Retry logic |
| **Shutdown** | Cleanup after completion | `./examples/run-shutdown.sh` | Shutdown commands |

## Direct YAML Usage

You can also pipe any YAML file directly:

```bash
cat examples/basic.yaml | ./goncurrently
cat examples/dev-servers.yaml | ./goncurrently
cat examples/tui-example.yaml | ./goncurrently
```

## Helper Scripts

Use these in your own configurations:

```yaml
commands:
  - cmd: bash
    args: ["examples/scripts/counter.sh", "10", "1"]
  
  - cmd: bash
    args: ["examples/scripts/server-sim.sh", "MyServer", "8080", "60"]
```

### Available Scripts

- `counter.sh [count] [delay]` - Count from 1 to N
- `random-fail.sh [fail_rate]` - Randomly fail (default 50%)
- `server-sim.sh [name] [port] [duration]` - Simulate server
- `watcher-sim.sh [directory]` - Simulate file watcher

## Testing

Run automated tests on all examples:

```bash
./examples/test-all.sh
```

## Tips

- Press `Ctrl+C` once to gracefully stop all processes
- Press `Ctrl+C` twice to force immediate termination
- Use `enableTUI: true` for better visualization of multiple processes
- Use `killOthers: true` to stop all processes when any one exits
- Set `restartTries: -1` for unlimited restarts
