# FMOS Documentation

FMOS (Fix My OS) is a beginner-friendly Windows repair, maintenance, and
first-time setup helper. This folder holds developer-facing documentation.
For user instructions, see the [main README](../README.md).

## How the code is organized

All source lives in `src/`, split into small files that each do one thing:

| File | Responsibility |
| --- | --- |
| `main.py` | Command-line entry point. Parses `gui` / `scan` / `repair` / `setup-apps` / `logs`; opens the GUI when given no command. |
| `gui.py` | Tkinter graphical window. A front end over `operations.py` / `apps.py` — adds no new repair behavior. |
| `operations.py` | Shared "what FMOS does" layer (scan/repair sequences and summaries). Used by both `main.py` and `gui.py`. |
| `commands.py` | Describes each Windows repair command (what it does, admin needed, whether it is risky). Runs nothing. |
| `runner.py` | The only place that actually runs a repair command. Does the safety checks, captures output, saves a log. |
| `apps.py` | First-time setup: the `winget` app catalog and installer. |
| `logger.py` | Writes plain-text logs to `logs/`. |
| `platform_check.py` | Read-only checks: are we on Windows, and are we Administrator? |
| `ui.py` | Beginner-friendly messages and yes/no confirmations. Pluggable so the GUI can route them into the window. |

## How the CLI and GUI share code

`ui.py` sends every message through a swappable output function and asks every
yes/no question through a swappable confirm function. The CLI leaves these at
their console defaults; `gui.py` calls `ui.set_output()` and
`ui.set_confirm()` so the same `operations.py` / `runner.py` code writes into
the text box and shows a pop-up dialog. This is why the two front ends behave
identically and the safety rules live in one place.

The GUI runs repairs on a background thread so the window stays responsive,
and marshals confirmation dialogs back to the main thread (Tkinter requires
UI calls on the main thread).

## Design rules

- **Safety over convenience.** Read-only checks never change the system.
  Risky commands are explained and confirmed first, and default to "no".
- **Describe, then run.** `commands.py` holds the metadata; `runner.py`
  does the running. This makes it easy to see exactly what FMOS can do.
- **No surprise system changes.** `chkdsk` is read-only; `/f` and `/r` are
  never added automatically.
- **Always log.** Every action writes a log to `logs/`.
- **Standard library only.** No third-party dependencies are required.

## Running the tests

From the project root:

```powershell
python -m unittest discover -s tests
```

The tests use `unittest` and mocks, so they run on any operating system
without executing real Windows repair commands.
