package platform

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// RestartAsAdmin starts a new copy of FMOS as Administrator, with the given
// arguments. Windows shows its normal Administrator (UAC) prompt first. A
// running program cannot gain Administrator rights, so this is the only way.
// The caller should quit after this returns nil.
func RestartAsAdmin(args []string) error {
	program, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find where FMOS is running from: %w", err)
	}
	folder, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not find the current folder: %w", err)
	}

	verb, verbErr := windows.UTF16PtrFromString("runas") // "run as Administrator"
	file, fileErr := windows.UTF16PtrFromString(program)
	params, paramsErr := windows.UTF16PtrFromString(windows.ComposeCommandLine(args))
	dir, dirErr := windows.UTF16PtrFromString(folder)
	if err := errors.Join(verbErr, fileErr, paramsErr, dirErr); err != nil {
		return fmt.Errorf("could not prepare the restart: %w", err)
	}

	err = windows.ShellExecute(0, verb, file, params, dir, windows.SW_SHOWNORMAL)
	if errors.Is(err, windows.ERROR_CANCELLED) {
		return ErrElevationDeclined
	}
	if err != nil {
		return fmt.Errorf("Windows could not start FMOS as Administrator: %w", err)
	}
	return nil
}
