package logging

import (
	"fmt"
	"os"
	"path/filepath"
)

// logsDir returns $XDG_STATE_HOME/fmos/logs, usually ~/.local/state/fmos/logs,
// the standard place for app logs on Linux.
func logsDir() (string, error) {
	stateFolder := os.Getenv("XDG_STATE_HOME")
	if stateFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not find your home folder: %w", err)
		}
		stateFolder = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(stateFolder, "fmos", "logs"), nil
}
