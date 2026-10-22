//go:build !windows

package sysinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// fastfetchModule is one entry in fastfetch's JSON output. A module that
// failed has an Error instead of a Result.
type fastfetchModule struct {
	Type   string          `json:"type"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

// parseFastfetch turns fastfetch's JSON output into Info. A module that
// cannot be read leaves its value empty and adds a problem.
func parseFastfetch(data []byte) (Info, error) {
	var modules []fastfetchModule
	if err := json.Unmarshal(data, &modules); err != nil {
		return Info{}, err
	}

	var info Info
	for _, module := range modules {
		if err := info.addFastfetchModule(module); err != nil {
			info.Problems = append(info.Problems,
				fmt.Sprintf("fastfetch could not read %s details: %v.", module.Type, err))
		}
	}
	return info, nil
}

// addFastfetchModule fills in the Info field that one module provides.
func (info *Info) addFastfetchModule(module fastfetchModule) error {
	if module.Error != "" {
		return errors.New(module.Error)
	}

	switch module.Type {
	case "OS":
		var result struct {
			PrettyName string `json:"prettyName"`
			Name       string `json:"name"`
			Version    string `json:"version"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		info.OS = result.PrettyName
		if info.OS == "" {
			info.OS = strings.TrimSpace(result.Name + " " + result.Version)
		}

	case "CPU":
		var result struct {
			Name string `json:"cpu"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		info.CPU = result.Name

	case "GPU":
		var result []struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		names := make([]string, 0, len(result))
		for _, gpu := range result {
			names = append(names, gpu.Name)
		}
		info.GPU = strings.Join(names, ", ")

	case "Memory":
		var result struct {
			Total uint64 `json:"total"`
			Used  uint64 `json:"used"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		info.RAM = formatRAM(result.Used, result.Total)

	case "Disk":
		var result []struct {
			Mountpoint string `json:"mountpoint"`
			Bytes      struct {
				Free  uint64 `json:"free"`
				Total uint64 `json:"total"`
			} `json:"bytes"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		// Show the main disk, the one mounted at "/".
		for _, disk := range result {
			if disk.Mountpoint == "/" {
				info.Disk = formatDisk(disk.Bytes.Free, disk.Bytes.Total)
			}
		}

	case "Uptime":
		var result struct {
			Milliseconds uint64 `json:"uptime"`
		}
		if err := json.Unmarshal(module.Result, &result); err != nil {
			return err
		}
		info.Uptime = formatUptime(time.Duration(result.Milliseconds) * time.Millisecond)
	}
	return nil
}
