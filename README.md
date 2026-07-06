# FMOS

FMOS (Fix My OS) is a program for repairing, cleaning, and setting up
Windows PCs.

## Requirements

- Windows 10 or Windows 11
- Python 3.14
- An administrator account
- An internet connection

## Usage

1. Download or clone this repository.
2. ...
3. ...

## Features

### System Repair

Runs the following built-in Windows commands:

```text
sfc /scannow
DISM /Online /Cleanup-Image /CheckHealth
DISM /Online /Cleanup-Image /ScanHealth
DISM /Online /Cleanup-Image /RestoreHealth
DISM /Online /Cleanup-Image /StartComponentCleanup
```

These commands check and repair the active Windows image, verify protected
system files, and remove replaced Windows components. Depending on the system,
the process may take a significant amount of time.

## Important

FMOS is provided without a warranty and cannot guarantee that every Windows problem
will be repaired.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
