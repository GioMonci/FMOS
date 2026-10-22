//go:build !windows

package runner

import "os/exec"

// detachFromConsole does nothing on macOS and Linux. FMOS's screen keeps the
// terminal in raw mode, so Ctrl+C reaches FMOS as a key press instead of
// stopping the program.
func detachFromConsole(cmd *exec.Cmd) {}
