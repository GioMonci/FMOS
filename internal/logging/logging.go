// Package logging saves a plain-text log for each FMOS session in the user's
// own app-data folder, so people can see what FMOS did and send it to someone
// for help. Each operating system's folder is defined in its dir_*.go file.
package logging

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GioMonci/FMOS/internal/actions"
	"github.com/GioMonci/FMOS/internal/platform"
	"github.com/GioMonci/FMOS/internal/runner"
	"github.com/GioMonci/FMOS/internal/sysinfo"
)

const timeLayout = "2006-01-02 15:04:05 -07:00"

// Session is the log file for one run of FMOS.
type Session struct {
	file         *os.File
	started      time.Time
	systemLogged bool
	actions      []string // one summary line per action, for the end of the log
}

// Start creates a new session log and writes its header. An error means
// logs cannot be saved. FMOS can still show information, but it must not run
// actions that change the system without a log.
func Start(version string) (*Session, error) {
	dir, err := logsDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("could not create the logs folder %s: %w", dir, err)
	}

	started := time.Now()
	file, err := createLogFile(dir, started)
	if err != nil {
		return nil, fmt.Errorf("could not create a log file in %s: %w", dir, err)
	}

	header := fmt.Sprintf("FMOS session log\n\nStarted:  %s\nVersion:  %s\nLog file: %s\n",
		started.Format(timeLayout), version, file.Name())
	if _, err := file.WriteString(header); err != nil {
		file.Close()
		return nil, fmt.Errorf("could not write to %s: %w", file.Name(), err)
	}
	return &Session{file: file, started: started}, nil
}

// Path returns the full path of the session's log file.
func (s *Session) Path() string {
	return s.file.Name()
}

// Dir returns the folder that holds all FMOS logs.
func (s *Session) Dir() string {
	return filepath.Dir(s.file.Name())
}

// LogSystemInfo writes the System panel details, including why any of them
// could not be read.
func (s *Session) LogSystemInfo(info sysinfo.Info) error {
	var text strings.Builder
	text.WriteString("\nSystem\n")
	for _, field := range info.Fields() {
		value := field.Value
		if value == "" {
			value = "Unknown"
		}
		fmt.Fprintf(&text, "  %-8s %s\n", field.Label+":", value)
	}
	for _, problem := range info.Problems {
		fmt.Fprintf(&text, "  Problem: %s\n", strings.ReplaceAll(problem, "\n", " "))
	}

	s.systemLogged = true
	_, err := s.file.WriteString(text.String())
	return err
}

// LogAction writes everything about one finished action: the exact command,
// whether FMOS ran as admin, timing, exit code, outcome, and full output. The
// file is flushed to disk right away, in case the computer crashes later.
func (s *Session) LogAction(action actions.Action, result runner.Result) error {
	name := action.Name
	exitCode := "none (the program did not run to completion)"
	if result.ExitCode >= 0 {
		exitCode = fmt.Sprint(result.ExitCode)
	}
	var text strings.Builder
	fmt.Fprintf(&text, "\nAction: %s\n", name)
	fmt.Fprintf(&text, "  Command:   %s\n", result.Command)
	fmt.Fprintf(&text, "  Admin:     %s\n", yesNo(platform.IsAdmin()))
	fmt.Fprintf(&text, "  Started:   %s\n", result.Started.Format(timeLayout))
	fmt.Fprintf(&text, "  Duration:  %s\n", logDuration(result.Duration))
	fmt.Fprintf(&text, "  Exit code: %s\n", exitCode)
	fmt.Fprintf(&text, "  Outcome:   %s\n", result.Outcome)
	fmt.Fprintf(&text, "  Changes the system:  %s\n", yesNo(action.ChangesSystem))
	fmt.Fprintf(&text, "  May need a restart:  %s\n", yesNo(action.MayNeedRestart))
	if result.Err != nil {
		fmt.Fprintf(&text, "  Error:     %v\n", result.Err)
	}
	if meaning := action.Meaning(result); meaning != "" {
		fmt.Fprintf(&text, "  Result:    %s\n", meaning)
	}
	fmt.Fprintf(&text, "  Output (%d lines):\n", len(result.Output))
	for _, line := range result.Output {
		fmt.Fprintf(&text, "    | %s\n", line)
	}

	summary := fmt.Sprintf("%s: %s", name, result.Outcome)
	if result.Outcome == runner.Failed {
		summary += fmt.Sprintf(" (exit code %d)", result.ExitCode)
	}
	s.actions = append(s.actions, summary)

	if _, err := s.file.WriteString(text.String()); err != nil {
		return err
	}
	return s.file.Sync()
}

func yesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// Note writes a short timestamped line, such as the result of an update check.
func (s *Session) Note(text string) error {
	_, err := fmt.Fprintf(s.file, "\n%s  %s\n", time.Now().Format("15:04:05"), text)
	return err
}

// logDuration shows short runs in milliseconds and longer ones in seconds.
func logDuration(duration time.Duration) time.Duration {
	if duration < time.Second {
		return duration.Round(time.Millisecond)
	}
	return duration.Round(time.Second)
}

// Close writes the end-of-session summary and closes the log file.
func (s *Session) Close() error {
	ended := time.Now()
	var text strings.Builder
	text.WriteString("\nSummary\n")
	fmt.Fprintf(&text, "  Ended:   %s\n", ended.Format(timeLayout))
	fmt.Fprintf(&text, "  Ran for: %s\n", ended.Sub(s.started).Round(time.Second))
	if !s.systemLogged {
		text.WriteString("  System details were not gathered before FMOS closed.\n")
	}
	if len(s.actions) == 0 {
		text.WriteString("  No actions were run in this session.\n")
	} else {
		fmt.Fprintf(&text, "  Actions run: %d\n", len(s.actions))
		for _, action := range s.actions {
			fmt.Fprintf(&text, "    - %s\n", action)
		}
	}

	_, writeErr := s.file.WriteString(text.String())
	closeErr := s.file.Close()
	return errors.Join(writeErr, closeErr)
}

// createLogFile creates a new file named after the start time, such as
// fmos-2026-10-22_13-45-10.txt. If two sessions start in the same second (for
// example when FMOS restarts as Administrator), later ones get "-2", "-3", ...
func createLogFile(dir string, started time.Time) (*os.File, error) {
	base := "fmos-" + started.Format("2006-01-02_15-04-05")
	for attempt := 1; ; attempt++ {
		name := base + ".txt"
		if attempt > 1 {
			name = fmt.Sprintf("%s-%d.txt", base, attempt)
		}
		file, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, fs.ErrExist) && attempt < 100 {
			continue
		}
		return file, err
	}
}
