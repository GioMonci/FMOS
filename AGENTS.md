# AGENTS.md

## Project Overview

FMOS stands for **Fix My OS**.

FMOS is a terminal app (TUI) that puts common OS repair, maintenance, and
troubleshooting tools in one place. A lot of useful options are hidden in
settings pages, built-in admin tools, or separate apps. FMOS gathers them into
one simple, beginner-friendly menu. The goal is to make common repair tasks
easier, safer, and less annoying for people who do not want to remember
commands like `sfc`, `dism`, `chkdsk`, or app setup steps.

**Windows is the main target and is built first.** macOS and Ubuntu Linux get
support later, for features that make sense there.

**FMOS does not reinvent existing tools.** It runs tools that already exist
(`sfc`, `DISM`, `winget`, Microsoft Defender, `fastfetch`, and so on), shows
their results clearly, logs what happened, and explains it in plain English.

The roadmap and the order of work live in `todo.md`.

### Current state

The repository contains a minimal Go command-line foundation
(`cmd/fmos/main.go`). It only prints help and version information. It does not
run repair commands or change the system. The TUI is next (step 2 in `todo.md`).

### Planned interface

```text
┌─ System ─────────────────────────┐ ┌─ Menu ───────────┐
│ OS:     Windows 11 Pro 24H2      │ │ > OS Repair      │
│ CPU:    ...                      │ │   Tools          │
│ GPU:    ...                      │ │   Network        │
│ RAM:    12.1 GB / 32 GB          │ │   Logs           │
│ Disk:   C: 210 GB free of 931 GB │ │   Quit           │
│ Admin:  Yes                      │ │                  │
└──────────────────────────────────┘ └──────────────────┘
 ↑/↓ move   Enter select   Esc back   q quit
```

* **System panel (left):** basic, read-only system information. It uses
  `fastfetch --format json` when it is installed and a small built-in fallback
  when it is not.
* **Menu (right):** feature sections. Selecting one shows its actions.
  * **OS Repair:** on Windows, SFC, DISM, and read-only disk checks.
  * **Tools:** installing software (winget), driver checks and updates,
    Microsoft Store updates, and Microsoft Defender scans.
  * **Network:** connection type and link speed, speed and latency checks,
    devices on the local network, and which apps are using the network.
  * **Logs:** view this session's logs and earlier ones.

Start simple. Show basic information first. Deeper features come later.

## General Agent Instructions

When working on this project:

* Prioritize safety over convenience.
* Never run destructive commands automatically.
* Do not hide risky behavior behind simple buttons or menu items.
* Make the code readable, boring, and dependable.
* Prefer clear step-by-step logic over clever abstractions.
* Keep functions small and focused.
* Explain dangerous or system-changing commands clearly.
* Preserve existing behavior unless the requested task says otherwise.
* Do not rewrite large parts of the project unless necessary.
* Do not delete files, rename major folders, or change public APIs without a clear reason.
* Keep the user experience simple and beginner-friendly.
* Do not over-engineer.
* Prefer simple solutions.
* Create easy, readable, and maintainable code.
* Follow the order in `todo.md`, and check items off only after verifying them.

This project is meant to help people fix their computers, not make their problems worse.

## Project Philosophy

FMOS should be:

* Helpful
* Safe
* Clear
* Easy to use
* Honest about what it can and cannot fix

FMOS should not be:

* Overly aggressive
* Full of unexplained system changes
* A one-click disaster button

A good rule:

> If a command could make the system worse, FMOS should explain it before running it.

## Use Existing Tools

FMOS is a central place to reach tools that already exist. It is not a
replacement for them.

* Prefer tools built into the operating system.
* When a built-in app already does the job well (Resource Monitor, Windows
  Update settings, Microsoft Store), opening it is fine. Tell the user what was opened.
* If a feature needs an extra tool (fastfetch, Ookla Speedtest CLI, nethogs),
  check whether it is installed. If it is not, explain what the tool is and
  offer to install it through the official package manager. Install only after
  the user confirms.
* Parse a tool's output only when needed for a clear summary. Always keep the
  raw output in the log.
* Do not reimplement what an existing tool already does reliably.

## Cross-Platform Rules

* **Windows first.** Build and verify a feature on Windows before porting it.
* "Linux" means **Ubuntu** (apt, systemd, NetworkManager). Do not assume other
  distributions unless support for them is added on purpose.
* Keep OS-specific code in files with Go OS suffixes: `_windows.go`,
  `_darwin.go`, `_linux.go`. Shared code (the TUI, logging, the command runner)
  must not contain OS-specific commands or paths.
* Each OS provides its own list of actions. A menu section with nothing
  available on the current OS says so ("Not available on macOS yet") instead
  of disappearing silently.
* Do not translate commands between operating systems blindly. `chkdsk` and
  `fsck` are not equal in risk, and macOS has no SFC equivalent.

## Action Definitions

Describe every action FMOS can run in one place, as a plain Go struct, with:

* The name shown in the menu
* A plain-English explanation
* The exact program and arguments
* Whether admin or root is required
* Whether it changes the system or is read-only
* Whether a restart may be needed
* Whether internet access is needed

The TUI shows this information before running anything that changes the
system. Keep it a simple list of structs. Do not build a plugin system.

## Safety Rules

Never automatically run commands that may cause data loss, system instability, or major configuration changes.

Be extra careful with:

* `chkdsk /f`
* `chkdsk /r`
* `fsck`, especially on mounted filesystems
* Disk formatting
* Partition changes (`diskpart`, `diskutil`, `parted`)
* Registry edits
* Driver removal
* User account changes
* Windows Update resets
* Service disabling
* Firewall changes
* Boot configuration changes
* App uninstallers
* Package removal (`apt remove`, `apt autoremove`)
* Bulk file deletion
* Anything run with `sudo`
* Network scanning (only on networks the user owns)

Risky actions should have:

* A clear explanation
* A confirmation step
* Logging
* A safer alternative when possible

## Running Commands

* Use `os/exec` with a fixed program and a separate argument list. Never build
  a shell command string, and never pass user input through `cmd /c` or `sh -c`.
* Prefer trusted system paths for system tools (for example, the System32
  folder under `%SystemRoot%`) over whatever is first on `PATH`.
* Use PowerShell only when no simpler tool exists. Run it with `-NoProfile
  -NonInteractive` and fixed script text.
* Never freeze the TUI. Run long work in the background (a Bubble Tea
  `tea.Cmd`) and stream output to the screen and the log.
* Show elapsed time for long tasks. Do not show fake percentage progress.
* Commands that need an interactive terminal, such as `sudo` password prompts,
  hand the terminal over with `tea.ExecProcess`. FMOS never asks for or stores passwords.
* If a tool is missing, say which tool it is and how to get it.

## Windows Repair Commands

FMOS may use commands such as:

```powershell
sfc /scannow
```

```powershell
DISM /Online /Cleanup-Image /CheckHealth
```

```powershell
DISM /Online /Cleanup-Image /ScanHealth
```

```powershell
DISM /Online /Cleanup-Image /RestoreHealth
```

```powershell
chkdsk
```

Agents should not assume these commands are harmless.

When adding support for repair commands (on any OS):

* Explain what the command does.
* Explain whether admin permissions are required.
* Explain whether a restart may be needed.
* Capture command output.
* Save logs.
* Report success, failure, or partial completion clearly.
* Do not claim the OS is fixed unless there is evidence.

## Software Installs and Updates

FMOS helps install common apps and keep them updated.

Examples:

* Browsers
* Developer tools
* Archive tools
* Media tools
* Hardware utilities
* Communication apps
* Game launchers

When adding app installation features:

* Use the official package manager for each OS: `winget` on Windows, Homebrew
  on macOS, and `apt` on Ubuntu (`snap` only when that is the official route).
* Never install random executables from unknown sources.
* Show what will be installed, and the exact command, before installing.
* Allow users to choose apps. Do not update everything by default.
* Log successful and failed installs.
* Do not bundle unwanted apps.

Example `winget` command:

```powershell
winget install --id Google.Chrome -e
```

## Logging Requirements

FMOS should keep logs whenever it runs repair, setup, or network actions.

Logs should include:

* Date and time
* Operating system
* Command or action performed
* Whether admin mode was used
* Exit code if available
* Important output
* Error messages
* Whether a restart is recommended

Logs should be saved in:

```text
logs/
```

Do not store logs in random locations.

## Folder Structure

Current project structure:

```text
.
├── cmd/fmos/             main package (help and version only for now)
├── logs/                 Runtime logs and log documentation
├── go.mod                Go module definition
├── README.md             User instructions
├── CONTRIBUTING.md       Contributor instructions
├── AGENTS.md             Instructions for AI coding agents
├── CLAUDE.md             Points Claude Code to this file
└── todo.md               Roadmap
```

Planned Go structure:

```text
.
├── cmd/fmos/             main package: starts the TUI, handles --help and --version
├── internal/
│   ├── tui/              Bubble Tea screens: system panel, menu, action views
│   ├── sysinfo/          System panel data (fastfetch JSON or built-in fallback)
│   ├── actions/          Action definitions per OS (actions_windows.go, ...)
│   ├── runner/           Runs external commands, streams and captures output
│   ├── platform/         Admin/root detection and other small OS helpers
│   └── logging/          Session logs in logs/
├── logs/
├── go.mod
└── go.sum
```

This is the target shape, not something to scaffold all at once. Create a
package only when a feature needs it.

Keep generated logs in `logs/`. Do not add a `docs/` or `tests/` directory
unless the user asks for one. Tests and GitHub Actions are intentionally
deferred during this early rebuild. When tests are added, Go tests go next to
the code as `_test.go` files.

## Build System

FMOS uses Go (1.27 or newer) as a single module, `github.com/GioMonci/FMOS`,
with `go.mod` at the repository root.

Build and run from PowerShell:

```powershell
go build -o build/fmos.exe ./cmd/fmos
.\build\fmos.exe
```

Run without building a binary:

```powershell
go run ./cmd/fmos
```

Cross-compile a Windows build from macOS or Linux:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/fmos.exe ./cmd/fmos
```

Before committing, run:

```powershell
gofmt -l .
go vet ./...
```

Dependencies:

* Bubble Tea, Bubbles, and Lip Gloss (Charm) for the TUI
* `golang.org/x/sys` for Windows admin detection

Keep the build pure Go (`CGO_ENABLED=0`) so cross-compiling stays easy. Do not
add other dependencies, build tools, test frameworks, or CI configuration
unless a requested feature needs them, and explain why when you add one.

## Code Style

Use clear names and simple logic.

Preferred style:

* Code formatted with `gofmt` and clean under `go vet`
* Standard library before third-party dependencies
* Small files and small packages when possible
* Small functions with one purpose
* Minimal global state
* Errors returned with context (`fmt.Errorf("running sfc: %w", err)`)
* `defer` for cleanup (closing files, stopping processes)
* Helpful error messages
* Beginner-friendly output
* No unnecessary dependencies
* No hidden side effects

Avoid:

* Giant source files with unrelated responsibilities
* Ignored errors (`_ = err`) or silent failures
* `panic` for normal errors
* Interfaces before there is a second implementation
* Goroutines without a clear owner and a way to stop
* OS-specific code outside `_windows.go`, `_darwin.go`, or `_linux.go` files
* Obscure abbreviations
* Unexplained admin commands
* Hardcoded machine-specific paths
* Hardcoded usernames
* Hardcoded drive letters unless clearly intentional

## TUI Behavior

Running `fmos` with no arguments opens the TUI. A tiny non-interactive CLI stays available:

```powershell
fmos --help
```

```powershell
fmos --version
```

Scripting subcommands such as `fmos repair --sfc` may come later. Do not add
them until they are needed.

In the TUI:

* Keyboard first: arrow keys or `j`/`k` to move, Enter to select, Esc to go
  back, `q` to quit. Show key hints at the bottom.
* Stay readable in Windows Terminal and the older console window. Do not rely
  on color alone to show meaning.
* Handle small terminal windows. Stack the panels when there is not enough width.
* For system-changing actions, use this flow: explanation, confirmation, run
  with progress, summary.
* Print useful progress updates without spamming the user.

Example tone:

```text
Checking Windows system files...
This may take a while. Do not close this window.
```

When an action finishes, summarize what happened.

Example:

```text
SFC completed successfully.
Windows did not report any integrity violations.
Log saved to logs/sfc-2026-07-06.txt
```

## User Experience Guidelines

FMOS should explain things like a helpful tech friend.

Use plain English.

Instead of:

```text
DISM exited with code 0.
```

Prefer:

```text
DISM finished successfully. Windows did not report a repair failure.
```

Instead of:

```text
Access denied.
```

Prefer:

```text
This action needs administrator permissions. Try running FMOS as Administrator.
```

## Admin Permissions

Many repair commands require Administrator (Windows) or root (macOS/Linux) access.

When admin permission is required:

* Detect whether the process is running as Administrator or root.
* If not, clearly tell the user.
* Do not fail with a confusing error.
* Do not repeatedly request elevation.
* Avoid surprise privilege escalation.

On Windows, a "Relaunch as Administrator" option may trigger the normal UAC
prompt once, but only when the user chooses it. On macOS and Ubuntu, run FMOS
as a normal user, and run individual root-only commands through `sudo` with the
terminal handed over to the user.

Example message:

```text
This repair needs Administrator permissions.
Close this window, right-click FMOS, and choose "Run as administrator."
```

## Error Handling

Errors should be useful.

When something fails, explain:

* What failed
* Why it may have failed
* What the user can try next
* Where the log was saved

Avoid vague messages like:

```text
Something went wrong.
```

Prefer:

```text
DISM could not complete the repair. This may happen if Windows Update is unavailable or the component store is damaged. The full output was saved to logs/dism-restorehealth.txt.
```

## Configuration

If the project uses config files, prefer readable formats such as:

```text
config.json
```

or:

```text
config.yaml
```

Do not commit private machine-specific config files.

Use examples instead:

```text
config.example.json
```

Config files should be documented in the README. Do not add configuration files
until a feature actually needs them.

## Secrets and Private Data

Never commit:

* Passwords
* API keys
* Private tokens
* Personal files
* User-specific logs with sensitive data
* `.env` files
* Machine-specific credentials

If secrets are needed, use environment variables or ignored local config files.

## Documentation Rules

Documentation should be beginner-friendly.

Keep project documentation in the root `README.md` for now. Do not recreate the
`docs/` directory unless the user requests it or the README becomes too large.

When adding a feature, update documentation with:

* What it does
* When to use it
* Which operating systems support it
* Whether it requires admin permissions
* Example command
* Expected output
* Possible risks

Keep commands copy-pasteable.

Good:

```powershell
fmos --help
```

Bad:

```text
Run the repair thing with the system file option.
```

## Git Rules

Keep changes focused.

Do not mix unrelated changes in the same update.

Avoid:

* Large formatting-only changes
* Renaming files for no reason
* Changing unrelated code
* Adding dependencies without explaining why
* Committing generated logs or local outputs

Before adding a dependency, consider whether the standard library or existing tools are enough.

## Generated Files

Generated runtime logs should go into:

```text
logs/
```

Examples:

```text
logs/
├── sfc-2026-07-06.txt
├── dism-checkhealth-2026-07-06.txt
└── setup-apps-2026-07-06.txt
```

Build output belongs in the ignored `build/` directory. Do not place generated
files inside source folders.

## Platform Assumptions

FMOS supports Windows first, then macOS and Ubuntu Linux.

* Do not assume a command works the same way on every OS.
* Do not assume Linux distributions other than Ubuntu.
* Keep OS-specific behavior behind the `_windows.go`, `_darwin.go`, and
  `_linux.go` file boundaries.

When writing user-facing commands, prefer Windows-friendly examples:

```powershell
.\build\fmos.exe
```

Use Bash examples only when the content is specifically about macOS, Linux, or
a Unix-like development environment.

## Good Final Response Format

When finishing a coding task, summarize:

* What changed
* Why it changed
* How to run it
* How to build and verify it
* Any important warnings or limitations

Example:

```text
Changed:
- Added an SFC repair action to the OS Repair menu.
- Added admin permission detection.
- Added log saving to logs/.

Run:
fmos  (then OS Repair > System File Checker)

Note:
This action requires Administrator permissions.
```

## When Unsure

If the safest path is obvious, make a reasonable assumption and explain it.

If the task could cause damage or data loss, ask before making the change.

Default to:

* Safety
* Clear logs
* Reversible changes
* Beginner-friendly explanations
* No surprise system modifications
