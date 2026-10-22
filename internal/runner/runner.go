// Package runner runs external programs for FMOS. It passes their output on
// line by line as it arrives and reports clearly how each run ended.
package runner

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// waitDelay bounds how long Run waits for output after the program exits,
// in case a helper program it started keeps the output pipe open.
const waitDelay = 10 * time.Second

// Command is one program to run, with arguments passed exactly as written.
// Program must be a full path, so FMOS never runs whatever happens to be
// first on PATH.
type Command struct {
	Program string
	Args    []string
}

// String returns the command the way it would be typed, for screens and logs.
func (c Command) String() string {
	parts := []string{quoteIfNeeded(c.Program)}
	for _, arg := range c.Args {
		parts = append(parts, quoteIfNeeded(arg))
	}
	return strings.Join(parts, " ")
}

func quoteIfNeeded(text string) string {
	if strings.ContainsAny(text, " \t") {
		return `"` + text + `"`
	}
	return text
}

// Outcome says how a run ended.
type Outcome int

const (
	Finished      Outcome = iota // ran and exited with code 0
	Failed                       // ran and exited with another code
	NotFound                     // the program does not exist
	FailedToStart                // the program exists but could not start
	Stopped                      // the user stopped it before it finished
)

// String returns a short plain-English name for the outcome.
func (o Outcome) String() string {
	switch o {
	case Finished:
		return "Finished"
	case Failed:
		return "Failed"
	case NotFound:
		return "Program not found"
	case FailedToStart:
		return "Could not start"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

// Result describes one run.
type Result struct {
	Command  Command
	Started  time.Time
	Duration time.Duration
	ExitCode int      // -1 when the program did not run to completion
	Output   []string // every line printed (normal output and errors), in order
	Outcome  Outcome
	Err      error // what went wrong, when Outcome is not Finished
}

// Run runs cmd and waits for it to finish. Each line the program prints is
// passed to onLine as soon as it arrives. Cancelling ctx stops the program,
// so only pass a cancellable ctx for actions that are safe to stop.
func Run(ctx context.Context, cmd Command, onLine func(string)) Result {
	result := Result{Command: cmd, Started: time.Now(), ExitCode: -1}
	name := filepath.Base(cmd.Program)
	if !filepath.IsAbs(cmd.Program) {
		result.Outcome = NotFound
		result.Err = fmt.Errorf("%s was not found (FMOS needs its full path)", name)
		return finish(result)
	}
	if _, err := os.Stat(cmd.Program); err != nil {
		result.Outcome = NotFound
		result.Err = fmt.Errorf("%s was not found at %s", name, cmd.Program)
		return finish(result)
	}

	process := exec.CommandContext(ctx, cmd.Program, cmd.Args...)
	detachFromConsole(process)
	process.WaitDelay = waitDelay

	// Normal output and error output share one pipe, so lines stay in the
	// order the program printed them.
	reader, writer := io.Pipe()
	process.Stdout = writer
	process.Stderr = writer

	if err := process.Start(); err != nil {
		result.Outcome = FailedToStart
		result.Err = fmt.Errorf("could not start %s: %w", name, err)
		if ctx.Err() != nil { // the user stopped it before it began
			result.Outcome = Stopped
			result.Err = errors.New("stopped before it started")
		}
		return finish(result)
	}

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		scanner.Split(scanLines)
		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), " \t")
			result.Output = append(result.Output, line)
			onLine(line)
		}
	}()

	waitErr := process.Wait()
	writer.Close() // lets the reader above finish
	<-readDone

	if process.ProcessState != nil {
		result.ExitCode = process.ProcessState.ExitCode()
	}
	switch {
	case ctx.Err() != nil:
		result.Outcome = Stopped
		result.Err = errors.New("stopped before it finished")
	case waitErr != nil && !(errors.Is(waitErr, exec.ErrWaitDelay) && result.ExitCode == 0):
		result.Outcome = Failed
		result.Err = waitErr
	default:
		result.Outcome = Finished
	}
	return finish(result)
}

// finish records how long the run took. Run calls it on every return.
func finish(result Result) Result {
	result.Duration = time.Since(result.Started)
	return result
}

// scanLines splits output into lines at "\n", "\r\n", or a lone "\r".
// Windows tools such as SFC and DISM rewrite progress lines with a lone "\r".
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			return i + 1, data[:i], nil
		}
		// A "\r": check whether a "\n" follows before deciding.
		if i+1 < len(data) {
			if data[i+1] == '\n' {
				return i + 2, data[:i], nil
			}
			return i + 1, data[:i], nil
		}
		if atEOF {
			return i + 1, data[:i], nil
		}
		return 0, nil, nil // need more data
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
