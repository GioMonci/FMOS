package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// maxDownloadSize stops a broken or unexpected download from filling the disk.
const maxDownloadSize = 200 << 20 // 200 MB

// releaseDownloads is the only place FMOS downloads updates from: this
// project's own GitHub Releases. A release whose links point anywhere else is
// refused.
var releaseDownloads = "https://github.com/GioMonci/FMOS/releases/download/"

// Install downloads the release's program file next to the running FMOS,
// checks it against the release's SHA-256 checksum, and swaps it in. If any
// step fails, the running version stays in place. The new version starts the
// next time FMOS is opened.
func Install(ctx context.Context, release Release) error {
	program, err := currentProgram()
	if err != nil {
		return err
	}
	return installTo(ctx, release, program)
}

// installTo does the work of Install for the program at the given path.
func installTo(ctx context.Context, release Release, program string) error {
	if release.AssetURL == "" {
		return errors.New("this release has no download for this computer")
	}
	if release.ChecksumsURL == "" {
		return errors.New("this release has no checksum file, so the download could not be verified")
	}
	for _, url := range []string{release.AssetURL, release.ChecksumsURL} {
		if !strings.HasPrefix(url, releaseDownloads) {
			return fmt.Errorf("the download link does not point to FMOS's GitHub Releases page: %s", url)
		}
	}

	want, err := fetchChecksum(ctx, release.ChecksumsURL, release.AssetName)
	if err != nil {
		return err
	}

	newPath := program + ".new"
	got, err := download(ctx, release.AssetURL, newPath)
	if err != nil {
		os.Remove(newPath)
		return err
	}
	if got != want {
		os.Remove(newPath)
		return errors.New("the downloaded file did not match its checksum, so it was not used")
	}
	return swap(program, newPath)
}

// CleanUp removes files left by an earlier update: the previous version
// ("<program>.old") and any unfinished download ("<program>.new"). Missing
// files are fine. On Windows the previous version cannot be removed while it
// is still running, so the next start tries again.
func CleanUp() error {
	program, err := currentProgram()
	if err != nil {
		return err
	}
	var problems []error
	for _, leftover := range []string{program + ".old", program + ".new"} {
		if err := os.Remove(leftover); err != nil && !errors.Is(err, fs.ErrNotExist) {
			problems = append(problems, err)
		}
	}
	return errors.Join(problems...)
}

// currentProgram returns the full path of the running FMOS program.
func currentProgram() (string, error) {
	program, err := os.Executable()
	if err == nil {
		program, err = filepath.EvalSymlinks(program)
	}
	if err != nil {
		return "", fmt.Errorf("could not find where FMOS is running from: %w", err)
	}
	return program, nil
}

// fetchChecksum reads the expected SHA-256 of one file from the release's
// checksum file, which has lines like "<hash>  fmos-x64.exe".
func fetchChecksum(ctx context.Context, url, name string) (string, error) {
	body, err := get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("could not download the checksum file: %w", err)
	}
	defer body.Close()

	data, err := io.ReadAll(io.LimitReader(body, 64<<10))
	if err != nil {
		return "", fmt.Errorf("could not download the checksum file: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == name {
			return strings.ToLower(fields[0]), nil
		}
	}
	return "", fmt.Errorf("the checksum file has no entry for %s", name)
}

// download saves the file at url to path and returns its SHA-256 as hex.
func download(ctx context.Context, url, path string) (string, error) {
	body, err := get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("could not download the new version: %w", err)
	}
	defer body.Close()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return "", fmt.Errorf("could not save the new version next to FMOS: %w", err)
	}
	hash := sha256.New()
	_, copyErr := io.Copy(io.MultiWriter(file, hash), io.LimitReader(body, maxDownloadSize))
	closeErr := file.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return "", fmt.Errorf("could not download the new version: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// get starts a download and returns its body. Anything but "200 OK" is an error.
func get(ctx context.Context, url string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "FMOS")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("the server answered %q", response.Status)
	}
	return response.Body, nil
}

// swap moves the new program into place. The running program is renamed to
// "<program>.old" first, which works even while it runs (Windows allows
// renaming a running program, just not deleting it). If moving the new file
// in fails, the old one is put back.
func swap(program, newPath string) error {
	oldPath := program + ".old"
	if err := os.Remove(oldPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		os.Remove(newPath)
		return fmt.Errorf("could not remove the copy left by an earlier update: %w", err)
	}
	if err := os.Rename(program, oldPath); err != nil {
		os.Remove(newPath)
		return fmt.Errorf("could not move the current version aside: %w", err)
	}
	if err := os.Rename(newPath, program); err != nil {
		if restoreErr := os.Rename(oldPath, program); restoreErr != nil {
			return fmt.Errorf("could not put the new version in place (%v) or restore the "+
				"current one (%v); the current version is saved as %s", err, restoreErr, oldPath)
		}
		return fmt.Errorf("could not put the new version in place: %w", err)
	}
	return nil
}
