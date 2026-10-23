# Contributing to FMOS

1. Create a branch from `main`.
2. Keep changes focused, readable, and easy to review.
3. Before opening a pull request, build the project and check it:

   ```powershell
   gofmt -l .
   go vet ./...
   go build ./...
   ```

   `gofmt -l .` should print nothing.
4. Explain what changed and how it was checked.

Windows repair behavior must be tested carefully on a supported Windows system.
Never add an unexplained or automatically destructive system command.

## Development builds

Go can build FMOS for every supported computer from one machine. Development
builds go in `build/dev/` and report the version `0.1.0-dev`. From macOS or Linux:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/dev/fmos-x64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o build/dev/fmos-arm64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/dev/fmos-macos ./cmd/fmos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/dev/fmos-linux ./cmd/fmos
```

From PowerShell on Windows:

```powershell
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o build/dev/fmos-x64.exe ./cmd/fmos
$env:GOOS = "windows"; $env:GOARCH = "arm64"; go build -o build/dev/fmos-arm64.exe ./cmd/fmos
$env:GOOS = "darwin"; $env:GOARCH = "arm64"; go build -o build/dev/fmos-macos ./cmd/fmos
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o build/dev/fmos-linux ./cmd/fmos
Remove-Item Env:CGO_ENABLED, Env:GOOS, Env:GOARCH
```

Never upload development builds. They contain this computer's folder paths and
the wrong version number.

## Making a release

Releases are built with a release script, so every release is built the same
way: `scripts/release.sh` on macOS or Linux, or `scripts/release.ps1` on
Windows. Both do the same steps and checks, and make the same files.

1. **Commit and push everything.** The script refuses to run with uncommitted
   changes, so a release always matches the code on GitHub.

2. **Pick the version number.** Use `MAJOR.MINOR.PATCH`, and make it higher than
   the last release. FMOS only offers an update when the number goes up.
   - Small fixes: `0.1.0` → `0.1.1`
   - New features: `0.1.1` → `0.2.0`

3. **Run the script** from the `FMOS` folder.

   On macOS or Linux:

   ```bash
   scripts/release.sh 0.1.0
   ```

   On Windows (PowerShell):

   ```powershell
   powershell -ExecutionPolicy Bypass -File scripts\release.ps1 0.1.0
   ```

   `-ExecutionPolicy Bypass` lets this one script run, because Windows blocks
   scripts by default. It does not change any Windows settings.

   **Trying it out:** run either script with a test version such as `0.0.1`.
   It only creates `build/release/v0.0.1/`, with no upload and no tag, so you
   can simply delete that folder afterwards.

   It creates `build/release/v0.1.0/` with `fmos-x64.exe`, `fmos-arm64.exe`,
   `fmos-macos`, `SHA256SUMS.txt`, and `notes.md`. Every program is built with
   `-trimpath` and the release version. The script checks that no folder paths
   from this computer are left inside, that the fingerprints match, and that the
   Mac build reports the right version. It stops if anything is wrong.

4. **Write the release notes** in `build/release/v0.1.0/notes.md`. Keep them
   short and plain English, because FMOS shows them on its Updates screen.

5. **Publish**, using either the command the script prints (needs the
   [GitHub CLI](https://cli.github.com/)) or github.com:
   - **Releases > Draft a new release**
   - **Tag:** `v0.1.0` (create it on `main`)
   - **Title:** `FMOS 0.1.0`
   - **Description:** paste `notes.md`
   - **Files:** attach `fmos-x64.exe`, `fmos-arm64.exe`, `fmos-macos`, and
     `SHA256SUMS.txt`. Do not attach `notes.md`.
   - Do not mark it as a pre-release: FMOS only looks at the latest full release.

6. **Check it.** Open a development build's **Updates** screen. It should offer
   the new version, because `0.1.0-dev` counts as older than `0.1.0`.

Keep the file names exactly as the script makes them. The updater and the
website's download links look for those names in every release.
