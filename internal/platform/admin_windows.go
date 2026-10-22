// Package platform holds small helpers whose behavior differs between
// operating systems.
package platform

import "golang.org/x/sys/windows"

// IsAdmin reports whether FMOS is running as Administrator (elevated by UAC).
func IsAdmin() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}
