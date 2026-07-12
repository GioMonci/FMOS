# FMOS

FMOS (Fix My OS) is a program for repairing, cleaning, and setting up
Windows PCs.

## Requirements

- Windows 10 or Windows 11
- Python 3.14
- An administrator account
- An internet connection

## Usage

FMOS has a simple graphical window (built with Tkinter) **and** a command line.
Both use the exact same safe repair logic.

1. Download or clone this repository.
2. Open a terminal in the project folder. For repairs, open it as
   **Administrator** (right-click, "Run as administrator").

### Graphical window (recommended)

```powershell
python src/main.py          # opens the FMOS window
python src/main.py gui      # same thing, explicitly
```

The window has buttons for **Scan**, **Repair** (with checkboxes for the SFC
and DISM steps), **Set up apps**, and **Show logs**. Progress appears in the
text area, and risky steps pop up a yes/no confirmation first.

Helpful touches:

- A banner at the top shows whether FMOS has Administrator permissions.
- Output is colour-coded (green success, amber warnings, red errors).
- A progress bar shows when a long repair is still working.
- **Clear**, **Save log…**, and **Open logs folder** buttons manage the output.
- Closing the window mid-repair asks for confirmation first.

### Command line

```powershell
python src/main.py scan          # Read-only health check (changes nothing)
python src/main.py repair        # Repair system files (asks before risky steps)
python src/main.py repair --sfc  # Run only System File Check
python src/main.py repair --dism # Run only DISM RestoreHealth
python src/main.py setup-apps    # Install common apps with winget
python src/main.py logs          # Show where logs are saved
```

## Features

### System check (`scan`)

A read-only health check. It runs `DISM /CheckHealth` and a read-only
`chkdsk`, reports what it finds, and **never changes anything**. Start here.

### System repair (`repair`)

Runs the following built-in Windows commands to check and repair the active
Windows image and protected system files:

```text
sfc /scannow
DISM /Online /Cleanup-Image /RestoreHealth
```

`RestoreHealth` changes system files, so FMOS explains it and asks you to
confirm before running it. Depending on the system, repairs may take a
significant amount of time. Use `--sfc` or `--dism` to run just one step.

### First-time app setup (`setup-apps`)

Shows a short list of well-known apps (browsers, 7-Zip, VLC, VS Code, and
more) and installs the ones you choose using `winget`, the official Windows
package manager. Nothing is installed without being shown and confirmed first.

### Logs (`logs`)

Every repair and install writes a plain-text log to the `logs/` folder,
including the command, exit code, output, and whether a restart is
recommended. Use the `logs` command to see where they are saved.

## Safety

FMOS prioritizes safety over convenience:

- Read-only checks (`scan`) never change your system.
- Risky steps are explained and require confirmation.
- `chkdsk` runs in read-only mode; it never uses `/f` or `/r` automatically.
- Commands needing Administrator permissions say so clearly instead of
  failing with a confusing error.
- Everything that runs is logged.

## Important

FMOS is provided without a warranty and cannot guarantee that every Windows problem
will be repaired.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
