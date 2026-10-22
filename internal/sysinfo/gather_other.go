//go:build !windows

package sysinfo

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// fastfetchArgs asks fastfetch for only the modules FMOS shows, as JSON.
// "--config none" ignores the user's own fastfetch settings.
var fastfetchArgs = []string{
	"--config", "none",
	"--format", "json",
	"--structure", "OS:CPU:GPU:Memory:Disk:Uptime",
}

func gatherForOS(ctx context.Context) Info {
	fastfetch, err := exec.LookPath("fastfetch")
	if err != nil {
		return Info{OS: osName(), Problems: []string{fastfetchMissingMessage()}}
	}

	output, err := exec.CommandContext(ctx, fastfetch, fastfetchArgs...).Output()
	if err != nil {
		return Info{
			OS:       osName(),
			Problems: []string{fmt.Sprintf("fastfetch could not read system details: %v.", err)},
		}
	}

	info, err := parseFastfetch(output)
	if err != nil {
		return Info{
			OS:       osName(),
			Problems: []string{fmt.Sprintf("Could not understand fastfetch's answer: %v.", err)},
		}
	}
	return info
}

// osName returns a basic name for the operating system, used when fastfetch
// cannot provide the full one.
func osName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}

func fastfetchMissingMessage() string {
	switch runtime.GOOS {
	case "darwin":
		return "Install fastfetch to see more details:\nbrew install fastfetch"
	case "linux":
		return "Install fastfetch to see more details:\nsudo apt install fastfetch"
	default:
		return "Install fastfetch to see more details."
	}
}
