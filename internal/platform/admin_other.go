//go:build !windows

package platform

import "os"

// IsAdmin reports whether FMOS is running as root.
func IsAdmin() bool {
	return os.Geteuid() == 0
}
