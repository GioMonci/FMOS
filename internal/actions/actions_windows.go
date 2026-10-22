package actions

import (
	"path/filepath"

	"github.com/GioMonci/FMOS/internal/platform"
)

func all() []Action {
	return []Action{
		{
			Section: "network",
			Name:    "Show network details",
			Description: "Lists your network adapters with their IP addresses, " +
				"gateway, and DNS servers. It only reads information and changes nothing.",
			Program: system32("ipconfig.exe"),
			Args:    []string{"/all"},
		},
	}
}

// system32 returns the full path of a built-in Windows tool. If the Windows
// folder is unknown, it returns only the name, which the runner reports as
// "not found" instead of searching PATH.
func system32(name string) string {
	windowsFolder, err := platform.WindowsFolder()
	if err != nil {
		return name
	}
	return filepath.Join(windowsFolder, "System32", name)
}
