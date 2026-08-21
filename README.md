# FMOS

FMOS (Fix My OS) is a Windows repair, maintenance, and first-time setup helper.

The project has been reset and is being rebuilt in C++. Repair and setup
features have not been added yet.

## Requirements

- Windows 10 or Windows 11
- A C++17-compatible compiler
- CMake 3.20 or newer

## Build

From PowerShell in the project folder:

```powershell
cmake -S . -B build
cmake --build build --config Release
```

The executable will be in `build/Release/` when using Visual Studio, or directly
in `build/` with a single-configuration generator.

## Run

```powershell
.\build\Release\fmos.exe help
```

## Current status

FMOS currently provides only a small command-line foundation. Windows repair
commands will be added carefully, with clear confirmations and logging before
any system-changing behavior is introduced.

Runtime logs will be stored in `logs/`. Generated logs are ignored by Git.

## Safety

FMOS does not currently run repair commands or modify Windows.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
