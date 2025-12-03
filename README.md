# goncurrently

A powerful, flexible command-line tool for running multiple commands concurrently in Go. Perfect for development workflows, build processes, and managing multiple services.

*Partially inspired by the [concurrently](https://www.npmjs.com/package/concurrently) npm package.*

[![Go Report Card](https://goreportcard.com/badge/github.com/sandrolain/goncurrently)](https://goreportcard.com/report/github.com/sandrolain/goncurrently)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Demo

### Terminal User Interface (TUI) Mode

![TUI Mode Demo](tapes/tui-example.gif)

The TUI mode provides a beautiful terminal interface with split panels for each command, showing real-time output and status.

### Microservices Architecture

![Microservices Demo](tapes/microservices.gif)

Run multiple services concurrently with setup commands, auto-restart, and graceful shutdown.

### Development Servers

![Dev Servers Demo](tapes/dev-servers.gif)

Perfect for running multiple development servers (frontend, backend, watchers) with colored output and auto-restart on failure.

### Basic Concurrent Execution

![Basic Demo](tapes/basic.gif)

Simple concurrent command execution with color-coded output for easy identification.

---

See the [`tapes/`](tapes/) directory for VHS tape files used to generate these demos.

## Features

- 🚀 **Concurrent Execution**: Run multiple commands simultaneously
- 🔄 **Auto-restart**: Automatically restart failed processes with configurable retry logic
- 🎨 **Color-coded Output**: Distinguish between different processes with colored output
- 📊 **TUI Mode**: Optional terminal user interface for better visualization
- ⏱️ **Timing Control**: Set delays before starting or restarting commands
- 🔧 **Setup Commands**: Run initialization commands before starting main processes
- 🛑 **Graceful Shutdown**: Handle SIGINT/SIGTERM signals with graceful termination
- ⚙️ **Environment Variables**: Set custom environment variables per command
- 🔇 **Silent Mode**: Suppress output from specific commands
- ⏰ **Command Timeouts**: Set maximum execution time for commands

## Installation

### Using go install

```bash
go install github.com/sandrolain/goncurrently@latest
```

This will install the `goncurrently` binary in your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your `PATH`.

### From Source

```bash
git clone https://github.com/sandrolain/goncurrently.git
cd goncurrently
go build -o goncurrently .
```

## Usage

> 💡 **Tip**: Check out the [Demo section](#demo) above to see goncurrently in action with animated examples!

goncurrently reads its configuration from YAML via stdin:

```bash
cat config.yaml | goncurrently
```

Or using a heredoc:

```bash
goncurrently << EOF
commands:
  - cmd: npm
    args: ["run", "dev"]
  - cmd: go
    args: ["run", "main.go"]
EOF
```

## Configuration

### Basic Configuration

```yaml
commands:
  - cmd: echo
    args: ["Hello, World!"]
```

### Full Configuration Example

```yaml
# Setup commands run sequentially before main commands
setupCommands:
  - name: database-migration
    cmd: ./scripts/migrate.sh
    args: ["up"]
    restartTries: 3
    restartAfter: "1s"

# Main commands run concurrently
commands:
  - name: frontend
    cmd: npm
    args: ["run", "dev"]
    restartTries: -1  # Unlimited restarts
    restartAfter: "2s"
    startAfter: "1s"  # Wait 1s before starting
    env:
      PORT: "3000"
      NODE_ENV: "development"

  - name: backend
    cmd: go
    args: ["run", "cmd/server/main.go"]
    restartTries: 5
    restartAfter: "1s"
    env:
      PORT: "8080"
      LOG_LEVEL: "debug"

  - name: watcher
    cmd: nodemon
    args: ["--watch", "src", "--exec", "npm run build"]
    silent: true  # Don't show output

  - name: cleanup
    cmd: ./cleanup.sh
    duration: "5m"  # Auto-stop after 5 minutes

# Shutdown commands run sequentially after all main commands complete
shutdownCommands:
  - name: stop-services
    cmd: docker
    args: ["compose", "down"]

  - name: cleanup-temp
    cmd: rm
    args: ["-rf", "/tmp/app-cache"]

# Global settings
killOthers: true      # Stop all commands if one exits
killTimeout: 5000     # Wait 5s before force-killing (milliseconds)
noColors: false       # Enable colored output
enableTUI: false      # Enable terminal UI mode
```

### Configuration Fields

#### Command Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `name` | string | Name of the command (auto-generated from cmd if not provided) | - |
| `cmd` | string | **Required**. Command to execute | - |
| `args` | []string | Command arguments | `[]` |
| `restartTries` | int | Number of restart attempts (-1 for unlimited) | `0` |
| `restartAfter` | string | Delay before restarting (e.g., "1s", "500ms") | `0` |
| `startAfter` | string | Delay before initial start | `0` |
| `env` | map[string]string | Environment variables | `{}` |
| `silent` | bool | Suppress command output | `false` |
| `duration` | string | Maximum execution time | - |

#### Global Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `commands` | []CommandConfig | **Required**. Commands to run concurrently | - |
| `setupCommands` | []CommandConfig | Commands to run sequentially before main commands | `[]` |
| `shutdownCommands` | []CommandConfig | Commands to run sequentially after all main commands complete | `[]` |
| `killOthers` | bool | Stop all commands if any one exits | `false` |
| `killTimeout` | int | Timeout in milliseconds before force kill | `0` |
| `noColors` | bool | Disable colored output | `false` |
| `enableTUI` | bool | Enable terminal UI mode | `false` |

## Examples

For ready-to-run examples, see the [`examples/`](examples/) directory. You can try them immediately:

```bash
# Build the tool first
go build -o goncurrently .

# Run any example
./examples/run-basic.sh
./examples/run-dev-servers.sh
./examples/run-tui.sh
```

See the [examples README](examples/README.md) for a complete list of available examples.

### Development Server

Run a frontend and backend server concurrently:

```yaml
commands:
  - name: frontend
    cmd: npm
    args: ["start"]
    env:
      PORT: "3000"
  
  - name: backend
    cmd: go
    args: ["run", "main.go"]
    env:
      PORT: "8080"

killOthers: true  # Stop both if either crashes
```

### Microservices with Setup

```yaml
setupCommands:
  - name: postgres
    cmd: docker
    args: ["start", "postgres"]
  
  - name: redis
    cmd: docker
    args: ["start", "redis"]

commands:
  - name: auth-service
    cmd: ./bin/auth-service
    restartTries: -1
    restartAfter: "2s"
  
  - name: api-gateway
    cmd: ./bin/api-gateway
    restartTries: -1
    restartAfter: "2s"
    startAfter: "3s"  # Wait for auth service
```

### Watch and Build

```yaml
commands:
  - name: watch-ts
    cmd: tsc
    args: ["--watch"]
  
  - name: watch-sass
    cmd: sass
    args: ["--watch", "src:dist"]
  
  - name: dev-server
    cmd: node
    args: ["server.js"]
    restartTries: -1
    restartAfter: "1s"
```

### Testing with Timeout

```yaml
commands:
  - name: unit-tests
    cmd: go
    args: ["test", "./..."]
    duration: "5m"  # Timeout after 5 minutes
  
  - name: integration-tests
    cmd: npm
    args: ["run", "test:integration"]
    duration: "10m"
    startAfter: "2s"
```

## TUI Mode

Enable the Terminal User Interface for a better visualization of multiple processes:

```yaml
enableTUI: true
commands:
  - cmd: npm
    args: ["start"]
  - cmd: go
    args: ["run", "main.go"]
```

In TUI mode, each command gets its own panel with colored borders and dedicated output area. The layout automatically adjusts based on the number of commands.

## Signal Handling

goncurrently handles interrupt signals gracefully:

- **First SIGINT/SIGTERM**: Initiates graceful shutdown, sends SIGTERM to all processes
- **Second SIGINT/SIGTERM**: Forces immediate termination

You can configure the grace period with `killTimeout` (in milliseconds).

## Duration Format

Duration strings support the following units:

- `ns` - nanoseconds
- `us` or `µs` - microseconds
- `ms` - milliseconds
- `s` - seconds
- `m` - minutes
- `h` - hours

Examples: `"100ms"`, `"2s"`, `"1.5m"`, `"1h30m"`

## Environment Variables

Commands inherit all environment variables from the parent process. Additional variables can be set per command:

```yaml
commands:
  - cmd: node
    args: ["server.js"]
    env:
      NODE_ENV: "production"
      PORT: "3000"
      DEBUG: "*"
```

## Exit Behavior

### Default Behavior

By default, goncurrently continues running even if individual commands exit.

### Kill Others Mode

With `killOthers: true`, goncurrently stops all commands when any command exits (excluding successful timeouts and restart-eligible failures).

## Restart Logic

Commands can be configured to restart automatically:

- `restartTries: -1` - Unlimited restarts
- `restartTries: 0` - No restarts (default)
- `restartTries: N` - Restart up to N times

Successful completions and timeout exits don't trigger restarts. Only error exits trigger the restart logic.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Your License Here]

## Related Projects

- [concurrently](https://github.com/open-cli-tools/concurrently) - Node.js version
- [hivemind](https://github.com/DarthSim/hivemind) - Ruby-based process manager
- [overmind](https://github.com/DarthSim/overmind) - Tmux-based process manager

## Author

Sandro Lain

## Support

For issues and questions, please use the [GitHub issue tracker](https://github.com/sandrolain/goncurrently/issues).
