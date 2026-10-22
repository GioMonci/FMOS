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
| `↑` / `↓` (or `k` / `j`) | Move through the menu, or scroll a finished action's output |
| `Enter` | Open the selected section, or run the selected action |
| `Esc` | Go back. While a read-only action runs, stop it |
| `Home` / `End`, `PgUp` / `PgDn` | Jump through a finished action's output |
| `q` | Quit FMOS (not while an action is running) |

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

FMOS currently shows a System panel and its menu (OS Repair, Tools, Network,
and Logs), and saves a log of each session (see [Logs](#logs)). The Logs
section shows where the log is saved. Network has its first action (below).
OS Repair and Tools are placeholders that say "Coming soon".

The System panel shows your OS, CPU, GPU, RAM, disk space, uptime, user, and
whether FMOS is running as Administrator. It only reads information and never
changes anything.

- **Windows:** the details come from Windows itself. Nothing extra is needed.
- **macOS and Linux:** the details come from
  [fastfetch](https://github.com/fastfetch-cli/fastfetch). If it is not
  installed, FMOS shows "Unknown" and tells you how to install it:

  ```bash
  brew install fastfetch        # macOS
  sudo apt install fastfetch    # Ubuntu
  ```

Repair commands will be added carefully, with clear explanations,
confirmations, and logging before any system-changing behavior is introduced.

### Network > Show network details

Lists your network adapters with their IP addresses, gateway, and DNS servers.
Use it to check whether your PC has an address on your network, or to send the
details to someone helping you.

- **Runs:** `ipconfig /all` on Windows, `ifconfig` on macOS, `ip address` on Linux.
- **Administrator:** not needed.
- **Risks:** none. It only reads information and changes nothing.

Open FMOS, choose **Network**, and press Enter on **Show network details**.
The output appears as it is printed, with the elapsed time. When it finishes,
FMOS says how it ended, for example:

```text
Finished in less than a second.
The full output is in this session's log (Logs > Open logs folder).
```

Use the arrow keys to scroll through the output, and Esc to go back. The full
output is also saved in the session log.

## Logs

Every time FMOS opens, it saves a log of that session, for example
`fmos-2026-10-22_13-45-10.txt`. The log lists your system details and a
summary of the session. You can open it with Notepad or any text editor, and
send it to someone helping you.

**To find your logs:** open FMOS, choose **Logs**, and press Enter on
**Open logs folder**. FMOS also prints the log's location when you quit.

Logs are saved in your own user folder, no matter where the FMOS program is
(even if it is still in Downloads):

| System | Logs folder |
| --- | --- |
| Windows | `%LOCALAPPDATA%\FMOS\logs`, for example `C:\Users\you\AppData\Local\FMOS\logs` |
| macOS | `~/Library/Logs/FMOS` |
| Linux | `~/.local/state/fmos/logs` |

Logs contain details about your PC, such as your user name, so only share them
with people you trust.

If FMOS cannot save logs, it shows a warning and keeps working. Check that your
disk has free space, then close FMOS and open it again.

## Safety

FMOS does not currently run repair commands or change your system.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
