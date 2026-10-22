package platform

import "os/exec"

// startDetached starts a program that opens its own window, such as a file
// manager, without waiting for that window to close.
func startDetached(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	// Wait in the background only to clean up the finished process. Its exit
	// code is ignored: File Explorer reports 1 even when it opened fine.
	go func() { _ = cmd.Wait() }()
	return nil
}
