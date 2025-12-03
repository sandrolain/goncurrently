package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
)

// Version is the current version of goncurrently.
const Version = "v1.2.0"

func printVersion() {
	fmt.Printf("goncurrently %s\n", Version) //nolint:forbidigo
}

func printHelp() {
	help := `goncurrently - Run multiple commands concurrently

Usage:
  cat config.yaml | goncurrently
  goncurrently < config.yaml
  goncurrently --help
  goncurrently --version

Commands:
  --help, -h       Show this help message
  --version, -v    Show version information

Configuration (via YAML on stdin):
  commands           List of commands to run concurrently (required)
  setupCommands      Commands to run sequentially before main commands
  shutdownCommands   Commands to run sequentially after all main commands complete
  killOthers         Stop all commands if any one exits (default: false)
  killTimeout        Timeout in milliseconds before force kill (default: 0)
  noColors           Disable colored output (default: false)
  enableTUI          Enable terminal UI mode (default: false)

Command Configuration:
  name               Name of the command (auto-generated if not provided)
  cmd                Command to execute (required)
  args               Command arguments
  restartTries       Number of restart attempts (-1 for unlimited)
  restartAfter       Delay before restarting (e.g., "1s", "500ms")
  startAfter         Delay before initial start
  env                Environment variables (map)
  silent             Suppress command output (default: false)
  duration           Maximum execution time

Examples:
  # Run a simple configuration
  cat <<EOF | goncurrently
  commands:
    - cmd: npm
      args: ["run", "dev"]
    - cmd: go
      args: ["run", "main.go"]
  EOF

For more information, visit: https://github.com/sandrolain/goncurrently
`
	fmt.Print(help) //nolint:forbidigo
}

func main() {
	// Handle command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			printHelp()
			return
		case "--version", "-v", "version":
			printVersion()
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %s\nUse --help for usage information.\n", os.Args[1])
			os.Exit(1)
		}
	}

	cfg, err := loadConfig(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	color.NoColor = cfg.NoColors
	assignNames(cfg.Commands)
	assignNames(cfg.SetupCommands)
	assignNames(cfg.ShutdownCommands)

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "invalid config: %v\n", err)
		os.Exit(1)
	}

	colors := defaultCommandColors()
	panelStyles := defaultPanelStyles(cfg.Commands)

	router, err := newOutputRouter(cfg.EnableTUI, cfg.Commands, panelStyles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize output routing: %v\n", err)
		os.Exit(1)
	}
	defer router.Stop()

	errorOutput = router.BaseWriter()
	baseLog("Initialized goncurrently | commands=%d setup=%d shutdown=%d killOthers=%t", len(cfg.Commands), len(cfg.SetupCommands), len(cfg.ShutdownCommands), cfg.KillOthers)

	termination := newTerminationManager(func(sig os.Signal, immediate bool) {
		if immediate {
			color.New(color.FgRed, color.Bold).Fprintf(errorOutput, "Second interrupt (%s) received, forcing termination...\n", sig) //nolint:errcheck
			router.Stop()
			return
		}
		color.New(color.FgRed, color.Bold).Fprintf(errorOutput, "Interrupt (%s) received, stopping all processes...\n", sig) //nolint:errcheck
	})
	defer termination.Shutdown()

	signals := termination.StopSignals()
	requestStop := termination.RequestStop

	runSetupSequence(cfg.SetupCommands, colors, router)
	if len(cfg.SetupCommands) > 0 {
		baseLog("Setup phase completed")
	}
	for i, c := range cfg.Commands {
		router.Add()
		go func(idx int, cc CommandConfig) {
			defer router.Done()
			baseLog("[%s] worker initialized", cc.Name)
			runManagedCommand(
				cc,
				colors[idx%len(colors)],
				router,
				signals,
				time.Duration(cfg.KillTimeout)*time.Millisecond,
				cfg.KillOthers,
				requestStop,
			)
		}(i, c)
	}

	router.Wait()

	if len(cfg.ShutdownCommands) > 0 {
		baseLog("Running shutdown commands...")
		runShutdownSequence(cfg.ShutdownCommands, colors, router)
		baseLog("Shutdown phase completed")
	}
}
