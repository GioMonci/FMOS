package runner

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// detachFromConsole gives the program its own hidden console instead of
// sharing FMOS's window. Anything it writes straight to a console then cannot
// draw over FMOS's screen, and closing FMOS's window or pressing Ctrl+C does
// not cut off a repair halfway.
func detachFromConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NO_WINDOW}
}
