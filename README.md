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

From PowerShell in the project folder:

```powershell
go build -o build/fmos.exe ./cmd/fmos
```

The program is saved as `build\fmos.exe`.

To build the Windows version from macOS or Linux:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/fmos.exe ./cmd/fmos
```

## Run

```powershell
.\build\fmos.exe --help
```

Expected output:

```text
FMOS (Fix My OS)

Usage:
  fmos --help       Show this help message.
  fmos --version    Show the FMOS version.

The FMOS menu is not available yet. For now, FMOS only shows this help.
```

You can also run FMOS without building it first:

```powershell
go run ./cmd/fmos --help
```

## Current status

FMOS currently provides only a small command-line foundation. The menu-based
terminal interface comes next. Repair commands will be added carefully, with
clear explanations, confirmations, and logging before any system-changing
behavior is introduced.

Runtime logs will be stored in `logs/`. Generated logs are ignored by Git.

## Safety

FMOS does not currently run repair commands or change your system.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
