# FMOS

FMOS (Fix My OS) is a terminal app that puts common OS repair, maintenance, and
troubleshooting tools in one place. Windows is the main target. macOS and
Ubuntu support will come later.

The project is being rebuilt in Go. Repair, tools, and network features have
not been added yet. See [todo.md](todo.md) for the roadmap.

## Requirements

- Windows 10 or Windows 11 (64-bit)
- [Go](https://go.dev/dl/) 1.27 or newer, to build FMOS

## Build

FMOS builds into the `build/` folder. Each file is named for the computer it runs on:

| File | Runs on |
| --- | --- |
| `fmos-x64.exe` | Windows on Intel or AMD (most PCs) |
| `fmos-arm64.exe` | Windows on ARM (for example, Snapdragon laptops) |
| `fmos-macos` | macOS on Apple Silicon (M1 and newer) |
| `fmos-linux` | Linux on Intel or AMD |

Not sure which Windows file you need? Open **Settings > System > About** and
look at **System type**. It says either "x64-based processor" or "ARM-based processor".

Build all four from macOS or Linux:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/fmos-x64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o build/fmos-arm64.exe ./cmd/fmos
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/fmos-macos ./cmd/fmos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/fmos-linux ./cmd/fmos
```

Build all four from PowerShell on Windows:

```powershell
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o build/fmos-x64.exe ./cmd/fmos
$env:GOOS = "windows"; $env:GOARCH = "arm64"; go build -o build/fmos-arm64.exe ./cmd/fmos
$env:GOOS = "darwin"; $env:GOARCH = "arm64"; go build -o build/fmos-macos ./cmd/fmos
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o build/fmos-linux ./cmd/fmos
Remove-Item Env:CGO_ENABLED, Env:GOOS, Env:GOARCH
```

## Run

Open the FMOS menu on Windows:

```powershell
.\build\fmos-x64.exe
```

On a Mac:

```bash
./build/fmos-macos
```

| Key | What it does |
| --- | --- |
| `↑` / `↓` (or `k` / `j`) | Move through the menu |
| `Enter` | Open the selected section |
| `Esc` | Go back to the menu |
| `q` | Quit FMOS |

Show the help without opening the menu:

```powershell
.\build\fmos-x64.exe --help
```

Expected output:

```text
FMOS (Fix My OS)

Usage:
  fmos              Open the FMOS menu.
  fmos --help       Show this help message.
  fmos --version    Show the FMOS version.
```

You can also run FMOS without building it first:

```powershell
go run ./cmd/fmos
```

## Current status

FMOS currently shows its menu (OS Repair, Tools, Network, and Logs). Every
section is a placeholder that says "Coming soon". The system information panel
comes next. Repair commands will be added carefully, with clear explanations,
confirmations, and logging before any system-changing behavior is introduced.

Runtime logs will be stored in `logs/`. Generated logs are ignored by Git.

## Safety

FMOS does not currently run repair commands or change your system.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
