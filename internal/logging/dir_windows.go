package logging

import (
	"errors"
	"os"
	"path/filepath"
)

// logsDir returns %LOCALAPPDATA%\FMOS\logs, the standard place for an app's
// own files on Windows. The user can always write there, wherever the FMOS
// program itself is, and it is not synced to OneDrive.
func logsDir() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", errors.New("the LOCALAPPDATA folder is not set")
	}
	return filepath.Join(localAppData, "FMOS", "logs"), nil
}
