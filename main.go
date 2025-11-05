package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
)

func main() {
	cfg, err := loadConfig(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	color.NoColor = cfg.NoColors
	assignNames(cfg.Commands)
	assignNames(cfg.SetupCommands)

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
	baseLog("Initialized goncurrently | commands=%d setup=%d killOthers=%t", len(cfg.Commands), len(cfg.SetupCommands), cfg.KillOthers)

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
}
