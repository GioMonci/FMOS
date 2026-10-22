package platform

import (
	"errors"
	"os"
)

// WindowsFolder returns the Windows folder, usually C:\Windows. FMOS finds
// built-in tools through it instead of trusting whatever is first on PATH.
func WindowsFolder() (string, error) {
	folder := os.Getenv("SystemRoot")
	if folder == "" {
		return "", errors.New("the Windows folder (SystemRoot) is not set")
	}
	return folder, nil
}
