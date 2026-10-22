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

	"github.com/GioMonci/FMOS/internal/sysinfo"
)

const timeLayout = "2006-01-02 15:04:05 -07:00"

// Session is the log file for one run of FMOS.
type Session struct {
	file         *os.File
	started      time.Time
	systemLogged bool
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
	text.WriteString("  No actions were run in this session.\n")

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
