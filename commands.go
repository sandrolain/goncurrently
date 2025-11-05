package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func streamOutput(writeLine func(string), r io.Reader) {
	if writeLine == nil {
		_, _ = io.Copy(io.Discard, r) //nolint:errcheck
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		writeLine(scanner.Text())
	}
}

func logCommandLine(stdoutWriter, stderrWriter func(string), identifier, message string) {
	switch {
	case stderrWriter != nil:
		stderrWriter(message)
	case stdoutWriter != nil:
		stdoutWriter(message)
	default:
		fmt.Fprintf(errorOutput, "%s%s\n", identifier, message) //nolint:errcheck
	}
}

func terminateProcess(cmd *exec.Cmd, killTimeout time.Duration, done <-chan error, immediate <-chan struct{}) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Signal(syscall.SIGTERM) //nolint:errcheck
	if killTimeout <= 0 {
		_ = cmd.Process.Kill() //nolint:errcheck
		return
	}
	timer := time.NewTimer(killTimeout)
	defer timer.Stop()
	for {
		select {
		case <-done:
			return
		case <-immediate:
			_ = cmd.Process.Kill() //nolint:errcheck
			return
		case <-timer.C:
			_ = cmd.Process.Kill() //nolint:errcheck
			return
		}
	}
}

func startProcess(c CommandConfig) (cmd *exec.Cmd, ctx context.Context, cancel context.CancelFunc, stdout, stderr io.ReadCloser, err error) {
	if dur := mustParseDurationField("duration", c.Duration, c.Name); dur > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), dur)
		cmd = exec.CommandContext(ctx, c.Cmd, c.Args...) // #nosec G204 -- test tool with controlled config
	} else {
		cmd = exec.Command(c.Cmd, c.Args...) // #nosec G204 -- test tool with controlled config
	}
	if c.Env != nil {
		env := os.Environ()
		for k, v := range c.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}
	if c.Silent {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		if stdout, err = cmd.StdoutPipe(); err != nil {
			return nil, ctx, cancel, nil, nil, err
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			return nil, ctx, cancel, nil, nil, err
		}
	}
	if err = cmd.Start(); err != nil {
		if cancel != nil {
			cancel()
		}
		return nil, ctx, nil, nil, nil, err
	}
	return cmd, ctx, cancel, stdout, stderr, nil
}

func executeOnce(c CommandConfig, identifier string, stdoutWriter, stderrWriter func(string), signals stopSignals, killTimeout time.Duration) (bool, bool, error) {
	cmd, ctx, cancel, stdout, stderr, err := startProcess(c)
	if err != nil {
		logCommandLine(stdoutWriter, stderrWriter, identifier, fmt.Sprintf("failed to start: %v", err))
		return false, false, err
	}
	if !c.Silent {
		go streamOutput(stdoutWriter, stdout)
		go streamOutput(stderrWriter, stderr)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	select {
	case err := <-done:
		timedOut := ctx != nil && ctx.Err() == context.DeadlineExceeded
		return timedOut, false, err
	case <-signals.stop:
		logCommandLine(stdoutWriter, stderrWriter, identifier, "interrupted")
		terminateProcess(cmd, killTimeout, done, signals.immediate)
		return false, true, nil
	}
}

func logCommandOutcome(name string, err error, timedOut bool) {
	switch {
	case err == nil:
		baseLog("[%s] completed successfully", name)
	case timedOut:
		baseLog("[%s] timed out: %v", name, err)
	default:
		baseLog("[%s] exited with error: %v", name, err)
	}
}

func handleNoRestart(name string, killOthers bool, alert *color.Color, requestStop func()) {
	if killOthers {
		alert.Fprintf(errorOutput, "Stopping all processes due to killOnExit triggered by '%s'\n", name) //nolint:errcheck
		if requestStop != nil {
			requestStop()
		}
	}
	baseLog("[%s] will not restart (killOthers=%t)", name, killOthers)
}

func logRestartSchedule(name string, attempt int, restartTries int, triesLeft int) {
	if restartTries >= 0 {
		baseLog("[%s] scheduling restart (attempt %d of %d, remaining retries %d)", name, attempt, restartTries+1, triesLeft)
		return
	}
	baseLog("[%s] scheduling restart (attempt %d)", name, attempt)
}

func runManagedCommand(c CommandConfig, col *color.Color, sink outputRouter, signals stopSignals, killTimeout time.Duration, killOthers bool, requestStop func()) {
	identifier := fmt.Sprintf("[%s] ", c.Name)
	stdoutPrefix := identifier
	stderrPrefix := fmt.Sprintf("[%s stderr] ", c.Name)
	if _, isTUI := sink.(*tuiRouter); isTUI {
		stdoutPrefix = ""
		stderrPrefix = "[stderr] "
	}
	stdoutWriter := sink.LineWriter(c.Name, col, stdoutPrefix)
	stderrWriter := sink.LineWriter(c.Name, col, stderrPrefix)
	alert := color.New(color.FgRed, color.Bold)
	triesLeft := c.RestartTries
	if waitStartDelay(c, signals.stop) {
		baseLog("[%s] start aborted before launch", c.Name)
		return
	}
	baseLog("[%s] starting", c.Name)
	attempt := 1
	for {
		timedOut, interrupted, err := executeOnce(c, identifier, stdoutWriter, stderrWriter, signals, killTimeout)
		if interrupted {
			baseLog("[%s] interrupted", c.Name)
			return
		}
		logCommandOutcome(c.Name, err, timedOut)
		if !shouldRestart(err, timedOut, &triesLeft, c.RestartTries) {
			handleNoRestart(c.Name, killOthers, alert, requestStop)
			return
		}
		attempt++
		logRestartSchedule(c.Name, attempt, c.RestartTries, triesLeft)
		if waitRestartDelay(c, signals.stop) {
			baseLog("[%s] restart aborted due to stop signal", c.Name)
			return
		}
		baseLog("[%s] restarting now", c.Name)
	}
}

func waitStartDelay(c CommandConfig, stop <-chan struct{}) bool {
	if d := mustParseDurationField("startAfter", c.StartAfter, c.Name); d > 0 {
		if stop == nil {
			time.Sleep(d)
			return false
		}
		select {
		case <-stop:
			return true
		case <-time.After(d):
			return false
		}
	}
	return false
}

func waitRestartDelay(c CommandConfig, stop <-chan struct{}) bool {
	if d := mustParseDurationField("restartAfter", c.RestartAfter, c.Name); d > 0 {
		if stop == nil {
			time.Sleep(d)
			return false
		}
		select {
		case <-stop:
			return true
		case <-time.After(d):
			return false
		}
	}
	return false
}

func shouldRestart(err error, timedOut bool, triesLeft *int, restartTries int) bool {
	if err == nil || timedOut {
		return false
	}
	if restartTries < 0 {
		return true
	}
	if *triesLeft > 0 {
		*triesLeft--
		return true
	}
	return false
}

func runSetupWithRetries(c CommandConfig, identifier string, stdoutWriter, stderrWriter func(string)) bool {
	triesLeft := c.RestartTries
	for {
		timedOut, _, err := executeOnce(c, identifier, stdoutWriter, stderrWriter, stopSignals{}, 0)
		if err == nil || timedOut {
			return true
		}
		if !shouldRestart(err, timedOut, &triesLeft, c.RestartTries) {
			return false
		}
		_ = waitRestartDelay(c, nil)
	}
}
