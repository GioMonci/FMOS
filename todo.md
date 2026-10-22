# FMOS To-Do

FMOS (Fix My OS) is a Go terminal app (TUI) that puts common OS repair,
maintenance, and troubleshooting tools in one place. A lot of useful options are
hidden in settings pages, built-in admin tools, or separate apps. FMOS gathers
them into one simple menu.

**We are not reinventing the wheel.** FMOS runs tools that already exist,
shows their results clearly, logs what happened, and explains it in plain
English. Fancier features can come later.

**Windows first.** Every feature is built and verified on Windows before it is
ported. macOS and Ubuntu come after that, for features that make sense there.

Work through the phases in order. Check off an item only when its behavior has
been verified on a real machine or VM. Follow the rules in `AGENTS.md`.

## Decisions

- **Language:** Go. Replaces the current C++ foundation and the earlier C#/.NET plan.
- **TUI:** Bubble Tea, Bubbles, and Lip Gloss from Charm.
- **Layout:** system info panel on the left, menu on the right.
- **Menu sections:** OS Repair, Tools, Network, Logs.
- **Platforms:** Windows (main), then macOS and Ubuntu.
- **System info:** `fastfetch --format json` when installed, parsed into
  FMOS's own panel, with a small built-in fallback when it is missing.
  fastfetch is usually not installed on Windows, so the Windows fallback needs
  to be good on its own.

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

## Open questions

- [ ] Where should logs go for a downloaded binary: `logs/` next to the
  executable, or under the user's profile with a `logs/` subfolder when the
  executable folder is not writable?
- [ ] Use Bubble Tea v1 or v2? Use whichever is the current stable release
  when step 1 starts.
- [ ] Is FMOS mainly for your own machines, or also for servicing other
  people's PCs? The earlier plan was technician-focused (client reports,
  office/engineering/gaming presets). Those ideas are parked under
  "Later ideas" until this is decided.

## Phase 1: Go foundation and TUI

### 1. [ ] Replace the C++ foundation with Go

- Remove `src/main.cpp` and `CMakeLists.txt`.
- Add `go.mod` and `cmd/fmos/main.go`.
- `fmos --help` and `fmos --version` work without opening the TUI.
- Update `.gitignore` for Go output (keep `build/` ignored).
- Update `README.md`, `CONTRIBUTING.md`, and the "Current state" and folder
  sections of `AGENTS.md` to match.
- Verify it builds and runs on Windows 10 and 11 (x64), and cross-compiles from macOS.

### 2. [ ] Build the TUI skeleton

- Bubble Tea app with the system panel on the left and the menu on the right.
- Menu: OS Repair, Tools, Network, Logs, Quit.
- Keyboard navigation with key hints at the bottom: arrows or `j`/`k`, Enter, Esc, `q`.
- Sections without features yet say "Coming soon" instead of being hidden.
- Stack the panels vertically when the terminal is too narrow.
- Stay readable in Windows Terminal and the older console window. Do not rely
  on color alone to show meaning.

### 3. [ ] Add the system info panel

- Show OS and version, CPU, GPU, RAM used/total, disk space, uptime, current
  user, and whether FMOS is running as admin.
- Use `fastfetch --format json` when it is installed. Otherwise use a built-in
  fallback, starting with Windows.
- If a value cannot be read, show "Unknown" and log why. Never crash the panel.
- Keep this panel read-only.

### 4. [ ] Add session logging

- Create a timestamped session log under `logs/` and show its full path.
- Record date and time, OS, admin status, the exact command, exit code,
  important output, errors, outcome, and whether a restart is recommended.
- Check that the log folder is writable before running any action.
- End each session with a short plain-English summary: what ran, what passed,
  what failed, and what needs follow-up.
- Generated logs stay out of Git.

### 5. [ ] Add the command runner

- Run known programs with explicit arguments through `os/exec`. Never build
  shell command strings.
- Stream output to the TUI and the log while showing elapsed time. Do not show
  fake percentage progress.
- Report each outcome clearly: tool not found, failed to start, failed,
  finished, or unclear result.
- Decide how cancellation works for each tool. Stopping a queue prevents the
  next task from starting, but does not abruptly kill a repair that is running.
- Test with harmless commands (`whoami`, `ipconfig`) before wiring up repairs.

### 6. [ ] Add admin detection and the action review screen

- Detect Administrator (Windows) or root (macOS/Linux).
- Describe each action with its name, plain-English explanation, exact command,
  admin requirement, whether it changes the system, whether a restart may be
  needed, and whether it needs internet.
- Before any system-changing action, show that description and ask for confirmation.
- If admin is required and missing, explain how to rerun FMOS as Administrator.
  A "Relaunch as Administrator" option may trigger one normal UAC prompt, only
  when the user chooses it, and must never loop.
- Never restart the computer automatically.

## Phase 2: Windows OS Repair

### 7. [ ] DISM health checks (read-only)

- Offer `DISM /Online /Cleanup-Image /CheckHealth` and `/ScanHealth` as separate choices.
- Explain the difference: CheckHealth reads already-recorded corruption, while
  ScanHealth does a full scan that takes longer.
- Show and log the result. Never start a repair automatically.

### 8. [ ] DISM RestoreHealth

- Explain it, confirm it, and say that it may need Windows Update or internet access.
- Report success, failure, or unresolved results honestly.
- When the repair source is unavailable, tell the user what to try next.

### 9. [ ] SFC check and repair

- Offer `sfc /verifyonly` (check only) and `sfc /scannow` (repair).
- Explain that `/scannow` can replace damaged Windows system files.
- When both are selected, run DISM RestoreHealth before SFC `/scannow`.
- When the output supports it, distinguish three results: no integrity
  violations, files repaired, and corruption that could not be fixed.

### 10. [ ] Disk health check (read-only)

- Show basic drive health (for example `Get-PhysicalDisk` HealthStatus).
- Offer a read-only `chkdsk <drive>:` scan with no `/f` or `/r`.
- If problems are found, explain what `chkdsk /f` would do (it may schedule a
  check at restart). Do not run it in this step.

## Phase 3: Windows Tools

### 11. [ ] Install and update software with winget

- Detect winget (App Installer). If it is missing, explain how to get it.
- Show a short curated app list by category: browsers, developer tools, archive
  tools, media, communication, and game launchers.
- The user picks apps and sees the exact package IDs and commands before
  anything installs. Log each install as succeeded or failed.
- List available updates (`winget upgrade`) and let the user choose which apps
  to update. Do not use `--all` by default.

### 12. [ ] Microsoft Store updates

- Open the Store's Downloads and Updates page
  (`ms-windows-store://downloadsandupdates`; check it works on Windows 10 and 11).
- Optionally list Store app updates through winget's `msstore` source.

### 13. [ ] Driver checks and updates

- List devices with problems (for example `pnputil /enum-devices /problem`). Read-only.
- Show the installed GPU and its driver version.
- Open Windows Update's optional updates page for driver updates
  (`ms-settings:windowsupdate-optionalupdates`).
- Link to the GPU vendor's official app or download page (NVIDIA, AMD, Intel).
- Never uninstall or roll back drivers automatically.

### 14. [ ] Microsoft Defender scans

- Check whether Defender is available and active. If another antivirus is in
  charge, say so, and do not try to override it.
- Offer a signature update, then a Quick or Full scan through a trusted
  `MpCmdRun.exe` path.
- Explain that Defender may quarantine or remove threats, depending on its settings.
- Show elapsed time instead of a percentage. Report completion, detections,
  and unresolved threats separately.

## Phase 4: Windows Network

The goal is to help answer "why am I lagging?"

### 15. [ ] Connection overview (read-only)

- Show active adapters, connection type (Wi-Fi or Ethernet), link speed, IP
  address, gateway, and DNS (`Get-NetAdapter`, `Get-NetIPConfiguration`).
- For Wi-Fi, show the network name, signal strength, and band
  (`netsh wlan show interfaces`).
- Explain that link speed is not the same as internet speed.

### 16. [ ] Speed and latency check

- Ping the gateway and a public host. Explain what the result means: slow to
  the gateway points to a local Wi-Fi or cable problem, while slow only past
  the gateway points to the ISP or the internet.
- Windows has no built-in speed test. If the user agrees, offer the official
  Ookla Speedtest CLI through winget, and mention that it has its own license
  to accept. Otherwise, open a speed test website.

### 17. [ ] Devices on the local network

- Show devices the PC already knows about (`Get-NetNeighbor` / `arp -a`).
- Explain the limitation: this only lists devices the PC has recently talked to.
- Active network scanning comes later, only as an explicit choice, and only on
  networks the user owns.

### 18. [ ] What is using the network

- Group open connections by app (`Get-NetTCPConnection` plus process names).
- Live per-app bandwidth is hard to get from a simple command, so open
  Resource Monitor's Network tab (`resmon`) and point the user to it.
- Read-only. FMOS does not end processes in the first version.

## Phase 5: macOS and Ubuntu

### 19. [ ] Shared groundwork

- System panel from fastfetch. Offer to install it with Homebrew or apt when it
  is missing. Older Ubuntu LTS releases may need the official fastfetch PPA.
- Root-only commands run through `sudo` and hand the terminal over to the user
  (`tea.ExecProcess`). FMOS never asks for or stores passwords.

### 20. [ ] macOS equivalents

- OS Repair: `diskutil verifyVolume` (read-only) and `softwareupdate --list`.
  macOS has no SFC equivalent because the system volume is sealed. Say so
  instead of pretending.
- Tools: Homebrew installs and updates (`brew install`, `brew outdated`).
- Network: `networkQuality` (built-in speed test), `nettop` (per-app network
  usage), and `arp -a`.

### 21. [ ] Ubuntu equivalents

- OS Repair: `systemctl --failed`, `journalctl -p err -b`, and drive health
  through `smartctl` (all read-only). Package fixes such as
  `dpkg --configure -a` and `apt --fix-broken install` change the system and
  need confirmation.
- Tools: apt installs and updates, `ubuntu-drivers list`, and `fwupdmgr get-updates`.
- Network: `ip addr`, `nmcli`, `ip neigh`, `ss -tunp`, and `nethogs` (per-app
  bandwidth, needs sudo).
- Never run `fsck` on a mounted filesystem.

## Phase 6: Packaging and release

### 22. [ ] Release builds

- One binary per platform: windows/amd64 and windows/arm64, darwin/arm64, and
  linux/amd64. Build with `CGO_ENABLED=0` so no runtime install is needed.
- Publish to GitHub Releases with release notes, checksums, and the license.
- Plan Windows code signing (SmartScreen) and macOS notarization (Gatekeeper)
  before wider distribution.
- Confirm release packages contain no logs, private files, or local settings.

### 23. [ ] Verify the full workflow

- Test on disposable Windows VMs before running repair features on a PC you care about.
- Cover these cases: standard and admin users, declined UAC, offline use, a log
  folder that cannot be written, Defender turned off or replaced by another
  antivirus, missing winget, command failures, and small terminal windows.
- Update `README.md` with the final build, download, usage, risk, and recovery guidance.
- Tests and CI stay deferred unless they are requested separately.

## Later ideas (parked)

- Non-interactive subcommands for scripting (`fmos repair --sfc`).
- Distribution through winget and a Homebrew tap. Add an installer (the earlier
  plan used Inno Setup) only if one is needed.
- A small GitHub Pages download site.
- Guided browser notification cleanup for Edge and Chrome, which helps with
  fake "virus" pop-ups. Guided first. Automate only what can be safely verified
  and restored.
- GPU driver help: Game Ready vs Studio vs professional drivers, and verified
  official downloads.
- Maintenance presets (office, engineering, gaming) and client-ready reports,
  if FMOS becomes technician-focused.
- An optional active network scan of the user's own network.

## Scope boundaries

- No automatic registry cleaning, bulk file deletion, browser profile resets,
  driver removal, forced restarts, or disk and partition changes.
- No `chkdsk /f`, `chkdsk /r`, or `fsck` without an explicit explanation and confirmation.
- No background service, telemetry, remote maintenance, or automatic updates in
  the first release.
- No scanning of networks the user does not own.
- Add dependencies, config files, and abstractions only when a feature needs them.
