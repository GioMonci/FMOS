package platform

import (
	"os/exec"
	"path/filepath"
)

// OpenFolder shows a folder in File Explorer.
func OpenFolder(path string) error {
	windowsFolder, err := WindowsFolder()
	if err != nil {
		return err
	}
	return startDetached(exec.Command(filepath.Join(windowsFolder, "explorer.exe"), path))
}
