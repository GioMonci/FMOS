package platform

import (
	"errors"
	"os/exec"
)

// OpenFolder shows a folder in the desktop's file manager.
func OpenFolder(path string) error {
	xdgOpen, err := exec.LookPath("xdg-open")
	if err != nil {
		return errors.New("xdg-open is not installed, so FMOS cannot open folders")
	}
	return startDetached(exec.Command(xdgOpen, path))
}
