//go:build !windows

package platform

// RestartAsAdmin is not available on macOS and Linux. There, root-only
// commands will run through sudo instead (todo.md step 24).
func RestartAsAdmin(args []string) error {
	return ErrElevationNotSupported
}
