package sysinfo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/GioMonci/FMOS/internal/platform"
)

// windowsInfoScript asks Windows for system details with read-only queries
// and prints them as JSON. Errors are silenced so that one failed query only
// leaves its own values empty. The lines are joined with "; " so the script
// is a single readable command line, which also separates the fields of the
// final @{ ... } table.
var windowsInfoScript = strings.Join([]string{
	"$ErrorActionPreference = 'SilentlyContinue'",
	"$ProgressPreference = 'SilentlyContinue'",
	"[Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $false",
	"$os = Get-CimInstance Win32_OperatingSystem",
	"$cpu = Get-CimInstance Win32_Processor | Select-Object -First 1",
	"$gpu = (Get-CimInstance Win32_VideoController | ForEach-Object { $_.Name }) -join ', '",
	"$disk = Get-CimInstance Win32_LogicalDisk | Where-Object DeviceID -eq $env:SystemDrive",
	`$cv = Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion'`,
	"$version = $cv.DisplayVersion",
	"if (-not $version) { $version = $cv.ReleaseId }",
	"$uptime = if ($os) { ((Get-Date) - $os.LastBootUpTime).TotalSeconds }",
	"[pscustomobject]@{ osName = $os.Caption",
	"osVersion = $version",
	"osBuild = $os.BuildNumber",
	"cpu = $cpu.Name",
	"gpu = $gpu",
	"memoryTotalKB = $os.TotalVisibleMemorySize",
	"memoryFreeKB = $os.FreePhysicalMemory",
	"drive = $env:SystemDrive",
	"diskTotal = $disk.Size",
	"diskFree = $disk.FreeSpace",
	"uptimeSeconds = $uptime } | ConvertTo-Json -Compress",
}, "; ")

// windowsDetails matches the JSON printed by windowsInfoScript.
type windowsDetails struct {
	OSName        string  `json:"osName"`
	OSVersion     string  `json:"osVersion"`
	OSBuild       string  `json:"osBuild"`
	CPU           string  `json:"cpu"`
	GPU           string  `json:"gpu"`
	MemoryTotalKB uint64  `json:"memoryTotalKB"`
	MemoryFreeKB  uint64  `json:"memoryFreeKB"`
	Drive         string  `json:"drive"`
	DiskTotal     uint64  `json:"diskTotal"`
	DiskFree      uint64  `json:"diskFree"`
	UptimeSeconds float64 `json:"uptimeSeconds"`
}

func gatherForOS(ctx context.Context) Info {
	details, err := readWindowsDetails(ctx)
	if err != nil {
		return Info{
			OS:       "Windows",
			Problems: []string{fmt.Sprintf("Could not read Windows system details: %v.", err)},
		}
	}

	var usedKB uint64
	if details.MemoryTotalKB >= details.MemoryFreeKB {
		usedKB = details.MemoryTotalKB - details.MemoryFreeKB
	}
	return Info{
		OS:     windowsOSName(details),
		CPU:    strings.TrimSpace(details.CPU),
		GPU:    details.GPU,
		RAM:    formatRAM(usedKB*1024, details.MemoryTotalKB*1024),
		Disk:   windowsDisk(details),
		Uptime: formatUptime(time.Duration(details.UptimeSeconds * float64(time.Second))),
	}
}

// readWindowsDetails runs windowsInfoScript with the built-in Windows
// PowerShell and decodes its answer.
func readWindowsDetails(ctx context.Context) (windowsDetails, error) {
	powerShell, err := powerShellPath()
	if err != nil {
		return windowsDetails{}, err
	}

	cmd := exec.CommandContext(ctx, powerShell, "-NoProfile", "-NonInteractive", "-Command", windowsInfoScript)
	output, err := cmd.Output()
	if ctx.Err() != nil {
		return windowsDetails{}, errors.New("Windows took too long to answer")
	}
	if err != nil {
		return windowsDetails{}, fmt.Errorf("PowerShell failed: %w", err)
	}

	var details windowsDetails
	output = bytes.TrimPrefix(output, []byte("\xef\xbb\xbf")) // UTF-8 byte order mark
	if err := json.Unmarshal(output, &details); err != nil {
		return windowsDetails{}, fmt.Errorf("could not understand PowerShell's answer: %w", err)
	}
	return details, nil
}

// powerShellPath finds the built-in Windows PowerShell through the Windows
// folder instead of trusting whatever is first on PATH.
func powerShellPath() (string, error) {
	windowsFolder, err := platform.WindowsFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(windowsFolder, "System32", "WindowsPowerShell", "v1.0", "powershell.exe"), nil
}

// windowsOSName returns text like "Windows 11 Pro 24H2 (build 26100)".
func windowsOSName(details windowsDetails) string {
	name := strings.TrimPrefix(strings.TrimSpace(details.OSName), "Microsoft ")
	if name == "" {
		return ""
	}
	if details.OSVersion != "" {
		name += " " + details.OSVersion
	}
	if details.OSBuild != "" {
		name += " (build " + details.OSBuild + ")"
	}
	return name
}

// windowsDisk returns text like "210.0 GB free of 931.5 GB (C:)".
func windowsDisk(details windowsDetails) string {
	disk := formatDisk(details.DiskFree, details.DiskTotal)
	if disk == "" || details.Drive == "" {
		return disk
	}
	return disk + " (" + details.Drive + ")"
}
