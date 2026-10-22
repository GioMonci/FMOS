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

## Building every file for a release

Go can build FMOS for every supported computer from one machine. From macOS or Linux:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/fmos-x64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o build/fmos-arm64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/fmos-macos ./cmd/fmos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/fmos-linux ./cmd/fmos
```

From PowerShell on Windows:

```powershell
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o build/fmos-x64.exe ./cmd/fmos
$env:GOOS = "windows"; $env:GOARCH = "arm64"; go build -o build/fmos-arm64.exe ./cmd/fmos
$env:GOOS = "darwin"; $env:GOARCH = "arm64"; go build -o build/fmos-macos ./cmd/fmos
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o build/fmos-linux ./cmd/fmos
Remove-Item Env:CGO_ENABLED, Env:GOOS, Env:GOARCH
```

Keep these file names exactly. The update checker looks for them in each release.
