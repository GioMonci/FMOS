package platform

import "os/exec"

// OpenFolder shows a folder in Finder.
func OpenFolder(path string) error {
	return startDetached(exec.Command("/usr/bin/open", path))
}
