package logging

import (
	"fmt"
	"os"
	"path/filepath"
)

// logsDir returns ~/Library/Logs/FMOS, the standard place for app logs on macOS.
func logsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find your home folder: %w", err)
	}
	return filepath.Join(home, "Library", "Logs", "FMOS"), nil
}
