# AGENTS.md

## Project Overview

FMOS stands for **Fix My OS**.

This project is a Windows repair, maintenance, and first-time setup helper. The goal is to make common Windows repair tasks easier, safer, and less annoying for users who do not want to manually remember commands like `sfc`, `dism`, `chkdsk`, or app setup steps.

FMOS may eventually support:

* Running common Windows repair commands
* Checking system health
* Running `sfc /scannow`
* Running DISM repair commands
* Checking disk health
* Installing first-time setup apps
* Applying basic Windows settings
* Saving logs for review
* Giving beginner-friendly explanations of what happened

## General Agent Instructions

When working on this project:

* Prioritize safety over convenience.
* Never run destructive commands automatically.
* Do not hide risky behavior behind simple buttons or commands.
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

This project is meant to help people fix Windows, not make their problems worse.

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

> If a command could make Windows worse, FMOS should explain it before running it.

## Safety Rules

Never automatically run commands that may cause data loss, system instability, or major configuration changes.

Be extra careful with:

* `chkdsk /f`
* `chkdsk /r`
* Disk formatting
* Partition changes
* Registry edits
* Driver removal
* User account changes
* Windows Update resets
* Service disabling
* Firewall changes
* Boot configuration changes
* App uninstallers
* Bulk file deletion

Risky actions should have:

* A clear explanation
* A confirmation step
* Logging
* A safer alternative when possible

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

When adding support for repair commands:

* Explain what the command does.
* Explain whether admin permissions are required.
* Explain whether a restart may be needed.
* Capture command output.
* Save logs.
* Report success, failure, or partial completion clearly.
* Do not claim the OS is fixed unless there is evidence.

## First-Time Setup Features

FMOS may eventually help install common first-time setup apps.

Examples:

* Browsers
* Developer tools
* Archive tools
* Media tools
* Hardware utilities
* Communication apps
* Game launchers

When adding app installation features:

* Prefer official package managers when possible.
* Prefer `winget` when available.
* Never install random executables from unknown sources.
* Show what will be installed before installing.
* Allow users to choose apps.
* Log successful and failed installs.
* Do not bundle unwanted apps.

Example `winget` command:

```powershell
winget install --id Google.Chrome -e
```

## Logging Requirements

FMOS should keep logs whenever it runs repair or setup actions.

Logs should include:

* Date and time
* Command or action performed
* Whether admin mode was used
* Exit code if available
* Important output
* Error messages
* Whether a restart is recommended

Logs should be saved somewhere predictable, such as:

```text
logs/
```

or:

```text
output/logs/
```

Avoid storing logs in random locations.

## Folder Structure

Recommended project structure:

```text
.
├── src/                  Main source code
├── logs/                 Runtime logs
├── docs/                 Standard documentation
├── README.md             Instructions for Users
└── AGENTS.md             Instructions for AI coding agents
```

Adjust this if the actual project structure changes.

## Code Style

Use clear names and simple logic.

Preferred style:

* Clear function names
* Small files when possible
* Small functions with one purpose
* Minimal global state
* Helpful error messages
* Beginner-friendly output
* No unnecessary dependencies
* No hidden side effects

Avoid:

* Giant scripts with everything in one file
* Obscure abbreviations
* Silent failures
* Unexplained admin commands
* Hardcoded machine-specific paths
* Hardcoded usernames
* Hardcoded drive letters unless clearly intentional

## CLI Behavior

If FMOS has a command-line interface, it should be easy to understand.

Good examples:

```powershell
fmos scan
```

```powershell
fmos repair
```

```powershell
fmos repair --sfc
```

```powershell
fmos repair --dism
```

```powershell
fmos setup-apps
```

```powershell
fmos logs
```

Commands should print useful progress updates without spamming the user.

Example tone:

```text
Checking Windows system files...
This may take a while. Do not close this window.
```

When a command finishes, summarize what happened.

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

Many Windows repair commands require Administrator access.

When admin permission is required:

* Detect whether the process is running as Administrator.
* If not, clearly tell the user.
* Do not fail with a confusing error.
* Do not repeatedly request elevation.
* Avoid surprise privilege escalation.

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

Config files should be documented in the README.

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

When adding a feature, update documentation with:

* What it does
* When to use it
* Whether it requires admin permissions
* Example command
* Expected output
* Possible risks

Keep commands copy-pasteable.

Good:

```powershell
fmos repair --sfc
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

Generated files should usually go into:

```text
output/
```

or:

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

Do not place generated files inside source folders unless there is a clear reason.

## Platform Assumptions

FMOS is primarily for Windows.

Agents should not assume Linux/macOS behavior unless the project explicitly adds cross-platform support.

When writing commands, prefer Windows-friendly examples:

```powershell
```

instead of Bash, unless the file is specifically for a Unix-like environment.

## Good Final Response Format

When finishing a coding task, summarize:

* What changed
* Why it changed
* How to run it
* How to test it
* Any important warnings or limitations

Example:

```text
Changed:
- Added an SFC repair command.
- Added admin permission detection.
- Added log saving to logs/.

Run:
fmos repair --sfc

Note:
This command requires Administrator permissions.
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
