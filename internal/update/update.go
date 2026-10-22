// Package update checks GitHub for a newer FMOS release and, when the user
// agrees, replaces the running program with it.
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

const (
	latestReleaseURL = "https://api.github.com/repos/GioMonci/FMOS/releases/latest"
	checksumsName    = "SHA256SUMS.txt"
)

// assetNames maps each platform to its file in a release. These match the
// file names in the README's Build section.
var assetNames = map[string]string{
	"windows/amd64": "fmos-x64.exe",
	"windows/arm64": "fmos-arm64.exe",
	"darwin/arm64":  "fmos-macos",
	"linux/amd64":   "fmos-linux",
}

// Release is a published version of FMOS.
type Release struct {
	Version      string // such as "0.2.0"
	Notes        string // release notes, written for users
	AssetName    string // this computer's file, such as "fmos-x64.exe"
	AssetURL     string // empty when the release has no file for this computer
	ChecksumsURL string // empty when the release has no checksum file
}

// CanInstall reports whether the release has what an update needs: this
// computer's program file and the checksum file to verify it.
func (r Release) CanInstall() bool {
	return r.AssetURL != "" && r.ChecksumsURL != ""
}

// Check asks GitHub for the latest release. newer is true when that release
// is newer than current. When no release has been published yet, newer is
// false and err is nil.
func Check(ctx context.Context, current string) (release Release, newer bool, err error) {
	return check(ctx, latestReleaseURL, current)
}

func check(ctx context.Context, url, current string) (Release, bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Release{}, false, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "FMOS/"+current)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return Release{}, false, fmt.Errorf("could not reach GitHub: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return Release{}, false, nil // no release has been published yet
	}
	if response.StatusCode != http.StatusOK {
		return Release{}, false, fmt.Errorf("GitHub answered %q", response.Status)
	}

	var latest struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&latest); err != nil {
		return Release{}, false, fmt.Errorf("could not understand GitHub's answer: %w", err)
	}

	release := Release{
		Version:   strings.TrimPrefix(latest.TagName, "v"),
		Notes:     strings.TrimSpace(strings.ReplaceAll(latest.Body, "\r\n", "\n")),
		AssetName: assetNames[runtime.GOOS+"/"+runtime.GOARCH],
	}
	for _, asset := range latest.Assets {
		if asset.Name == release.AssetName && release.AssetName != "" {
			release.AssetURL = asset.URL
		}
		if asset.Name == checksumsName {
			release.ChecksumsURL = asset.URL
		}
	}

	newer, err := isNewer(release.Version, current)
	if err != nil {
		return release, false, err
	}
	return release, newer, nil
}

// isNewer reports whether version a is newer than version b. Versions look
// like "1.2.3", optionally followed by a suffix such as "-dev". A version
// with a suffix comes before the same version without one.
func isNewer(a, b string) (bool, error) {
	aNumbers, aSuffix, err := parseVersion(a)
	if err != nil {
		return false, err
	}
	bNumbers, bSuffix, err := parseVersion(b)
	if err != nil {
		return false, err
	}
	for i := range aNumbers {
		if aNumbers[i] != bNumbers[i] {
			return aNumbers[i] > bNumbers[i], nil
		}
	}
	return aSuffix == "" && bSuffix != "", nil
}

// parseVersion splits a version such as "1.2.3-dev" into [1 2 3] and "dev".
func parseVersion(version string) ([3]int, string, error) {
	var numbers [3]int
	core, suffix, _ := strings.Cut(strings.TrimPrefix(version, "v"), "-")
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return numbers, "", fmt.Errorf("%q is not a version number like 1.2.3", version)
	}
	for i, part := range parts {
		number, err := strconv.Atoi(part)
		if err != nil || number < 0 {
			return numbers, "", fmt.Errorf("%q is not a version number like 1.2.3", version)
		}
		numbers[i] = number
	}
	return numbers, suffix, nil
}
