// Package sysinfo collects the basic, read-only system information shown in
// the FMOS System panel.
package sysinfo

import (
	"context"
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/GioMonci/FMOS/internal/platform"
)

// Info is the system information shown in the System panel. An empty string
// means the value could not be read.
type Info struct {
	OS     string
	CPU    string
	GPU    string
	RAM    string
	Disk   string
	Uptime string
	User   string
	Admin  bool

	// Problems explains in plain English why details are missing.
	Problems []string
}

// Gather collects system information. It never fails as a whole: anything
// that cannot be read is left empty and explained in Problems.
func Gather(ctx context.Context) Info {
	info := gatherForOS(ctx)
	info.Admin = platform.IsAdmin()

	name, err := currentUserName()
	if err != nil {
		info.Problems = append(info.Problems, fmt.Sprintf("Could not read the current user: %v.", err))
	}
	info.User = name
	return info
}

// currentUserName returns the signed-in user's name without the computer or
// domain prefix that Windows adds ("PC-NAME\name").
func currentUserName() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	name := current.Username
	if i := strings.LastIndex(name, `\`); i >= 0 {
		name = name[i+1:]
	}
	return name, nil
}

// formatRAM returns text like "16.0 GB used of 24.0 GB".
func formatRAM(used, total uint64) string {
	if total == 0 {
		return ""
	}
	return formatSize(used) + " used of " + formatSize(total)
}

// formatDisk returns text like "756.4 GB free of 926.4 GB".
func formatDisk(free, total uint64) string {
	if total == 0 {
		return ""
	}
	return formatSize(free) + " free of " + formatSize(total)
}

// formatSize shows a size the way Windows does, where 1 GB is 1024³ bytes.
func formatSize(bytes uint64) string {
	const gigabyte = 1 << 30
	size := float64(bytes) / gigabyte
	if size >= 1024 {
		return fmt.Sprintf("%.1f TB", size/1024)
	}
	return fmt.Sprintf("%.1f GB", size)
}

// formatUptime returns text like "3 days, 4 hours" or "25 minutes".
func formatUptime(uptime time.Duration) string {
	if uptime <= 0 {
		return ""
	}
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	switch {
	case days > 0:
		return plural(days, "day") + ", " + plural(hours, "hour")
	case hours > 0:
		return plural(hours, "hour") + ", " + plural(minutes, "minute")
	case minutes > 0:
		return plural(minutes, "minute")
	default:
		return "less than a minute"
	}
}

// plural returns text like "1 day" or "2 days".
func plural(count int, unit string) string {
	if count == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", count, unit)
}
