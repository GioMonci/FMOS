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

- **Who it's for:** me, my friends, and people whose PCs I build. Most of them
  are not technical. FMOS should let them fix common problems, or collect logs
  to send me, without me being there in person.
- **Language:** Go. Replaces the old C++ foundation and the earlier C#/.NET plan.
- **TUI:** Bubble Tea v2 and Lip Gloss v2 from Charm (`charm.land/...`).
- **Layout:** system info panel on the left, menu on the right.
- **Menu sections:** OS Repair, Tools, Network, Logs.
- **Platforms:** Windows (main), then macOS and Ubuntu.
- **System info:** on Windows, one read-only PowerShell query of built-in
  Windows data, because fastfetch is usually not installed there. On macOS and
  Linux, `fastfetch --format json`, parsed into FMOS's own panel.
- **Distribution:** a single portable `fmos.exe` on GitHub Releases. No installer.
- **Updates:** FMOS checks GitHub for a newer release and updates itself only
  after the user confirms.
- **Website:** one plain HTML/CSS page on GitHub Pages, served from `docs/`.
- **Logs:** the user's own app-data folder (`%LOCALAPPDATA%\FMOS\logs` on
  Windows), wherever the FMOS program is, even Downloads. It is always
  writable, not synced to OneDrive, and keeps every log in one place across
  updates. People reach it through Logs > Open logs folder.

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

- [ ] Code signing: unsigned downloads trigger a Windows SmartScreen warning.
  That's acceptable for friends at first, because the install guide explains
  it. Decide later whether paid signing is worth it.

## Phase 1: Go foundation and TUI

### 1. [x] Replace the C++ foundation with Go

Verified on macOS. The Windows build cross-compiles, but it has not been run on
Windows yet.

- Remove `src/main.cpp` and `CMakeLists.txt`.
- Add `go.mod` and `cmd/fmos/main.go`.
- `fmos --help` and `fmos --version` work without opening the TUI.
- Update `.gitignore` for Go output (keep `build/` ignored).
- Update `README.md`, `CONTRIBUTING.md`, and the "Current state" and folder
  sections of `AGENTS.md` to match.
- Verify it builds and runs on Windows 10 and 11 (x64), and cross-compiles from macOS.

### 2. [x] Build the TUI skeleton

Verified on macOS. Double-clicking `fmos.exe` and the older Windows console
window have not been tested on Windows yet.

- Bubble Tea app with the system panel on the left and the menu on the right.
- Menu: OS Repair, Tools, Network, Logs, Quit.
- Keyboard navigation with key hints at the bottom: arrows or `j`/`k`, Enter, Esc, `q`.
- Sections without features yet say "Coming soon" instead of being hidden.
- Stack the panels vertically when the terminal is too narrow.
- Stay readable in Windows Terminal and the older console window. Do not rely
  on color alone to show meaning.
- Double-clicking `fmos.exe` opens the menu, and the window stays open until
  the user quits.

### 3. [x] Add the system info panel

Verified on macOS, both with and without fastfetch. The Windows query compiles,
but it has not been run on Windows yet.

- Show OS and version, CPU, GPU, RAM used/total, disk space, uptime, current
  user, and whether FMOS is running as admin.
- On Windows, use one read-only PowerShell query. On macOS and Linux, use
  `fastfetch --format json`, and explain how to install it when it is missing.
- Load the details in the background so the menu opens right away.
- If a value cannot be read, show "Unknown" and explain why. Never crash the panel.
- Keep this panel read-only.

### 4. [x] Add session logging

Verified on macOS: logs land in `~/Library/Logs/FMOS`, the warning appears when
the folder cannot be created, two sessions in the same second get separate
files, quitting early is noted, and "Open logs folder" opens Finder. Not yet
run on Windows.

- Create a timestamped session log (`fmos-YYYY-MM-DD_HH-MM-SS.txt`) in the
  user's app-data folder: `%LOCALAPPDATA%\FMOS\logs` on Windows,
  `~/Library/Logs/FMOS` on macOS, `~/.local/state/fmos/logs` on Linux.
- Show the log's full path in the Logs section and after quitting, with an
  "Open logs folder" button (File Explorer, Finder, or `xdg-open`).
- At the start of each session, log the date and time, FMOS version, and the
  System panel details (including OS and admin status) with their `Problems`.
- Check that logs can be saved when FMOS starts. If not, keep running, show a
  warning, and suggest checking disk space and reopening FMOS.
- End each session with a short summary: when it ended, how long it ran, and
  what was run.
- Generated logs stay out of Git.

### 5. [x] Add the command runner

Verified on macOS with Network > Show network details (`ifconfig`), plus test
runs of every outcome. The Windows version (`ipconfig /all`) compiles, but it
has not been run on Windows yet.

- Run known programs by full path with explicit arguments through `os/exec`
  (`internal/runner`). Never build shell command strings or search PATH.
- Show output live on a full-width task screen with the elapsed time. No fake
  percentage progress. After it finishes, the output can be scrolled.
- Report each outcome clearly: finished, failed (with exit code), program not
  found, could not start, or stopped.
- Log each action: the exact command, admin status, start time, duration, exit
  code, outcome, errors, and full output. The end-of-session summary lists
  every action and how it ended.
- Refuse to run system-changing actions when logs cannot be saved.
- Actions are listed in `internal/actions` (one file per OS). The first one is
  Network > Show network details, which is read-only.

**Stopping (decided):** read-only actions can be stopped with Esc.
System-changing actions cannot be stopped from FMOS, so a repair is never cut
off halfway. FMOS cannot quit while an action runs, so every action finishes
and is logged. On Windows, programs run in their own hidden console, so
closing FMOS's window or pressing Ctrl+C does not interrupt them.

**Left for later steps:**

- "Restart recommended" and "unclear result" depend on reading each tool's
  output, so they are added with each repair (steps 12-15).
- Windows tools print in the console's code page, so non-English characters
  may look wrong. SFC prints UTF-16 when its output is captured. Decode both
  properly when adding SFC (step 14).

### 6. [ ] Add admin detection and the action review screen

- Admin detection already exists (`platform.IsAdmin`, added in step 3).
- Describe each action with its name, plain-English explanation, exact command,
  admin requirement, whether it changes the system, whether a restart may be
  needed, and whether it needs internet.
- Before any system-changing action, show that description and ask for confirmation.
- Never restart the computer automatically.

**Getting Administrator rights on Windows (decided):**

- Windows cannot raise a running program to Administrator, because admin
  rights are fixed when a program starts. FMOS can only start a new copy of
  itself as Administrator.
- When an action needs admin and FMOS does not have it, explain why and offer
  **"Restart FMOS as Administrator"**. Choosing it shows the normal Windows
  "Yes/No" (UAC) prompt, opens a new FMOS as Administrator, and closes the old
  window. Use `ShellExecute` with the `runas` verb from `golang.org/x/sys/windows`.
- The new FMOS opens where the user was. The old one passes the section and
  action at startup (for example `--open os-repair/sfc`), and the new one opens
  straight on that action's explanation and confirmation screen. It never runs
  the action automatically, so the user still confirms. Only the section and
  action carry over; anything else simply starts fresh.
- If the user picks "No" on the prompt (Windows error 1223, "cancelled"), keep
  running without admin, say what will not work, and never ask again on its own.
- On a standard (non-admin) Windows account, the Administrator copy runs as the
  admin account, so its log goes to that account's `%LOCALAPPDATA%` folder.
  Pass the current log folder to the new copy so all logs stay in one place.
- **Not chosen: always starting as Administrator** (a `requireAdministrator`
  manifest). Most features do not need admin. A prompt on every launch teaches
  people to click "Yes" without reading. On a standard user account, the
  Administrator copy runs as the admin account instead of the person's own
  account, which is the wrong user for browser settings and their apps.
- App installs do not need an Administrator FMOS. winget shows the Yes/No
  prompt itself when an installer needs it.
- FMOS is not code-signed, so the prompt says "Unknown publisher" (yellow).
  Explain this in the install guide (step 9).
- macOS and Linux do not relaunch. They run single root-only commands through
  `sudo`, with the terminal handed over to the user (`tea.ExecProcess`).

## Phase 2: First release for friends

Set up downloading, installing, and updating early, while FMOS is still
harmless (menu, system info, logs). It is much easier to test downloads,
SmartScreen warnings, and updates with friends before any repair features exist.

### 7. [ ] Add update checking and self-update

- This must be in the very first version you share. A copy without it never
  finds out that a new version exists.
- When FMOS starts, check GitHub's latest release in the background
  (`https://api.github.com/repos/GioMonci/FMOS/releases/latest`), and never
  block the menu while checking. Also add a "Check for updates" menu item.
- If a newer version exists, show its version number and release notes, and
  ask before updating. Never update automatically.
- To update:
  1. Download the new exe next to the current one.
  2. Check its SHA-256 against the checksum file in the release.
  3. Rename the running `fmos.exe` out of the way, then move the new one into place.
  4. Ask the user to restart FMOS.
  5. Delete the old copy the next time FMOS starts.
- If GitHub cannot be reached, show "Could not check for updates" and carry on normally.
- Log every check and every update. Use only the standard library (`net/http`,
  `crypto/sha256`).

### 8. [ ] Set up the release process (manual for now)

- Tag versions `v0.1.0`, `v0.2.0`, and so on. Build with
  `-ldflags "-X main.version=0.1.0"` so FMOS knows its own version.
- Attach `fmos-x64.exe`, `fmos-arm64.exe`, and `SHA256SUMS.txt` to each
  GitHub Release. The file names are listed in the README's Build section.
- Create releases on the GitHub website or with `gh release create`.
- Write short, plain-English release notes, because FMOS shows them on the update screen.
- Check that releases contain no logs, private files, or local settings.
- Automate releases only when doing them by hand becomes a chore.

### 9. [ ] Write the install guide in the README

- Write it for people who are not technical.
- Cover these points:
  - Which file to download (`fmos-x64.exe` or `fmos-arm64.exe`). Settings >
    System > About shows "x64-based processor" or "ARM-based processor".
  - Where to keep the file: anywhere they own, such as Downloads or the
    Desktop. Not `C:\Program Files`, where FMOS cannot update itself.
  - The browser download warning, and SmartScreen's "More info > Run anyway",
    including why each one appears.
  - Opening FMOS, and what the Administrator (Yes/No) prompt looks like when
    FMOS asks to restart as Administrator, including why it says "Unknown publisher".
  - How updates work.
  - How to send logs for help: Logs > Open logs folder, then send the newest file.
  - How to remove FMOS (delete the folder).
- Add screenshots once the TUI exists.

### 10. [ ] Build the GitHub Pages website

- Plain `docs/index.html` and `docs/style.css`. No frameworks, no build step.
- One page: what FMOS is, a download button, install steps, updating, what
  FMOS does and does not do, and a link to the GitHub repo.
- Point the download button at
  `https://github.com/GioMonci/FMOS/releases/latest/download/fmos-x64.exe`,
  with a smaller link for `fmos-arm64.exe`.
  That link always serves the newest release, so the site rarely needs editing.
- Make it work on phones, and support light and dark mode.
- Turn on Pages in the repo: Settings > Pages > Deploy from a branch > `main`,
  `/docs`. The site will be at `https://giomonci.github.io/FMOS/`.
- Keep its steps in sync with the README install guide.

### 11. [ ] Test the first release with friends

- Have one or two people install FMOS from the website without help.
- Publish a small new version and confirm that their copy updates itself.
- Fix anything confusing in the guide before sharing FMOS more widely.

## Phase 3: Windows OS Repair

### 12. [ ] DISM health checks (read-only)

- Offer `DISM /Online /Cleanup-Image /CheckHealth` and `/ScanHealth` as separate choices.
- Explain the difference: CheckHealth reads already-recorded corruption, while
  ScanHealth does a full scan that takes longer.
- Show and log the result. Never start a repair automatically.

### 13. [ ] DISM RestoreHealth

- Explain it, confirm it, and say that it may need Windows Update or internet access.
- Report success, failure, or unresolved results honestly.
- When the repair source is unavailable, tell the user what to try next.

### 14. [ ] SFC check and repair

- Offer `sfc /verifyonly` (check only) and `sfc /scannow` (repair).
- Explain that `/scannow` can replace damaged Windows system files.
- When both are selected, run DISM RestoreHealth before SFC `/scannow`.
- When the output supports it, distinguish three results: no integrity
  violations, files repaired, and corruption that could not be fixed.

### 15. [ ] Disk health check (read-only)

- Show basic drive health (for example `Get-PhysicalDisk` HealthStatus).
- Offer a read-only `chkdsk <drive>:` scan with no `/f` or `/r`.
- If problems are found, explain what `chkdsk /f` would do (it may schedule a
  check at restart). Do not run it in this step.

## Phase 4: Windows Tools

### 16. [ ] Install and update software with winget

- Detect winget (App Installer). If it is missing, explain how to get it.
- Show a short curated app list by category: browsers, developer tools, archive
  tools, media, communication, and game launchers.
- The user picks apps and sees the exact package IDs and commands before
  anything installs. Log each install as succeeded or failed.
- List available updates (`winget upgrade`) and let the user choose which apps
  to update. Do not use `--all` by default.

### 17. [ ] Microsoft Store updates

- Open the Store's Downloads and Updates page
  (`ms-windows-store://downloadsandupdates`; check it works on Windows 10 and 11).
- Optionally list Store app updates through winget's `msstore` source.

### 18. [ ] Driver checks and updates

- List devices with problems (for example `pnputil /enum-devices /problem`). Read-only.
- Show the installed GPU and its driver version.
- Open Windows Update's optional updates page for driver updates
  (`ms-settings:windowsupdate-optionalupdates`).
- Link to the GPU vendor's official app or download page (NVIDIA, AMD, Intel).
- Never uninstall or roll back drivers automatically.

### 19. [ ] Microsoft Defender scans

- Check whether Defender is available and active. If another antivirus is in
  charge, say so, and do not try to override it.
- Offer a signature update, then a Quick or Full scan through a trusted
  `MpCmdRun.exe` path.
- Explain that Defender may quarantine or remove threats, depending on its settings.
- Show elapsed time instead of a percentage. Report completion, detections,
  and unresolved threats separately.

## Phase 5: Windows Network

The goal is to help answer "why am I lagging?"

### 20. [ ] Connection overview (read-only)

- Network > Show network details (raw `ipconfig /all`) already exists from
  step 5. This step adds a plain-English overview on top of it.
- Show active adapters, connection type (Wi-Fi or Ethernet), link speed, IP
  address, gateway, and DNS (`Get-NetAdapter`, `Get-NetIPConfiguration`).
- For Wi-Fi, show the network name, signal strength, and band
  (`netsh wlan show interfaces`).
- Explain that link speed is not the same as internet speed.

### 21. [ ] Speed and latency check

- Ping the gateway and a public host. Explain what the result means: slow to
  the gateway points to a local Wi-Fi or cable problem, while slow only past
  the gateway points to the ISP or the internet.
- Windows has no built-in speed test. If the user agrees, offer the official
  Ookla Speedtest CLI through winget, and mention that it has its own license
  to accept. Otherwise, open a speed test website.

### 22. [ ] Devices on the local network

- Show devices the PC already knows about (`Get-NetNeighbor` / `arp -a`).
- Explain the limitation: this only lists devices the PC has recently talked to.
- Active network scanning comes later, only as an explicit choice, and only on
  networks the user owns.

### 23. [ ] What is using the network

- Group open connections by app (`Get-NetTCPConnection` plus process names).
- Live per-app bandwidth is hard to get from a simple command, so open
  Resource Monitor's Network tab (`resmon`) and point the user to it.
- Read-only. FMOS does not end processes in the first version.

## Phase 6: macOS and Ubuntu

### 24. [ ] Shared groundwork

- System panel from fastfetch. Offer to install it with Homebrew or apt when it
  is missing. Older Ubuntu LTS releases may need the official fastfetch PPA.
- Root-only commands run through `sudo` and hand the terminal over to the user
  (`tea.ExecProcess`). FMOS never asks for or stores passwords.

### 25. [ ] macOS equivalents

- OS Repair: `diskutil verifyVolume` (read-only) and `softwareupdate --list`.
  macOS has no SFC equivalent because the system volume is sealed. Say so
  instead of pretending.
- Tools: Homebrew installs and updates (`brew install`, `brew outdated`).
- Network: `networkQuality` (built-in speed test), `nettop` (per-app network
  usage), and `arp -a`.

### 26. [ ] Ubuntu equivalents

- OS Repair: `systemctl --failed`, `journalctl -p err -b`, and drive health
  through `smartctl` (all read-only). Package fixes such as
  `dpkg --configure -a` and `apt --fix-broken install` change the system and
  need confirmation.
- Tools: apt installs and updates, `ubuntu-drivers list`, and `fwupdmgr get-updates`.
- Network: `ip addr`, `nmcli`, `ip neigh`, `ss -tunp`, and `nethogs` (per-app
  bandwidth, needs sudo).
- Never run `fsck` on a mounted filesystem.

## Phase 7: Wider distribution

### 27. [ ] More platforms and package managers

- Add `fmos-macos` and `fmos-linux` to releases once the macOS and Ubuntu
  features exist. Build with `CGO_ENABLED=0`.
- Publish to winget so people can run `winget install` and `winget upgrade`.
  This needs a manifest pull request to `microsoft/winget-pkgs` for each release.
- Add a Homebrew tap for macOS.
- Add Windows code signing (SmartScreen) and macOS notarization (Gatekeeper)
  if the warnings become a real problem.
- Automate releases with GitHub Actions or GoReleaser once they are frequent.
- Add an installer only if a portable exe stops being enough.

### 28. [ ] Verify the full workflow

- Test on disposable Windows VMs before running repair features on a PC you care about.
- Cover these cases: standard and admin users, declined UAC, offline use, a log
  folder that cannot be written, Defender turned off or replaced by another
  antivirus, missing winget, command failures, failed updates, and small
  terminal windows.
- Update `README.md` and the website with the final download, usage, risk, and
  recovery guidance.
- Tests and CI stay deferred unless they are requested separately.

## Later ideas (parked)

- A "Get help" menu item for helping people remotely. It opens the logs folder
  so they can send logs, and it opens Windows Quick Assist, which is built into
  Windows, for remote help.
- Non-interactive subcommands for scripting (`fmos repair --sfc`).
- Guided browser notification cleanup for Edge and Chrome, which helps with
  fake "virus" pop-ups. Guided first. Automate only what can be safely verified
  and restored.
- GPU driver help: Game Ready vs Studio vs professional drivers, and verified
  official downloads.
- Maintenance presets (office, engineering, gaming) and client-ready reports.
- An optional active network scan of the user's own network.
- View logs inside FMOS.
- Clean up old logs automatically (for example, keep the newest 50). Only ever
  delete FMOS's own `fmos-*.txt` files, never anything else in the folder.

## Scope boundaries

- No automatic registry cleaning, bulk file deletion, browser profile resets,
  driver removal, forced restarts, or disk and partition changes.
- No `chkdsk /f`, `chkdsk /r`, or `fsck` without an explicit explanation and confirmation.
- No silent updates. FMOS may check for updates, but it installs one only
  after the user confirms.
- No background service, telemetry, or remote control built into FMOS.
- No scanning of networks the user does not own.
- Add dependencies, config files, and abstractions only when a feature needs them.
