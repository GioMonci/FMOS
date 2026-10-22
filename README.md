# FMOS: Fix My OS

FMOS is a simple terminal app that puts common computer repair, maintenance,
and troubleshooting tools in one menu, with plain-English explanations.

Lots of useful Windows tools are hidden in settings pages or need commands
nobody remembers. FMOS gathers them in one place. It doesn't reinvent them: it
runs the tools your computer already has, shows what they found, and keeps a
log you can send to someone helping you.

- **Windows first.** Some features also work on macOS and Ubuntu Linux.
- **Safe by design.** FMOS explains anything that would change your computer
  before it runs, and never does it without asking.
- **Free and open source** under the GPL-3.0 license.

> **Early development.** FMOS can already show system details and network
> details, save logs, and update itself. Repairs such as SFC and DISM are
> coming next. See the [roadmap](todo.md).

## Contents

- [What FMOS can do](#what-fmos-can-do)
- [Getting FMOS](#getting-fmos)
- [Using FMOS](#using-fmos)
- [Features](#features)
- [Safety and privacy](#safety-and-privacy)
- [Building FMOS yourself](#building-fmos-yourself)
- [Contributing](#contributing) and [License](#license)

## What FMOS can do

| Feature | What it does | Windows | macOS | Linux |
| --- | --- | --- | --- | --- |
| [System panel](#system-panel) | Shows your OS, CPU, GPU, RAM, disk space, and more | Yes | Yes, with fastfetch | Yes, with fastfetch |
| [Show network details](#network--show-network-details) | Lists your network adapters and addresses | Yes | Yes | Yes |
| [Logs](#logs) | Saves a log of every session, one key to open them | Yes | Yes | Yes |
| [Updates](#updates) | Checks for a new version and installs it when you choose | Yes | Checks only | Checks only |

**Coming next:** system file and disk repairs (SFC, DISM, CHKDSK), installing
and updating apps (winget), driver and Microsoft Store updates, Microsoft
Defender scans, and network speed checks.

## Getting FMOS

FMOS has not been released yet. When it is, you will download a single file
from the [Releases page](https://github.com/GioMonci/FMOS/releases). There is
nothing to install. Until then, you can [build it yourself](#building-fmos-yourself).

Pick the file for your computer:

| File | For |
| --- | --- |
| `fmos-x64.exe` | Windows on Intel or AMD (almost all desktops and laptops) |
| `fmos-arm64.exe` | Windows on ARM (for example, Snapdragon laptops) |

Not sure which one? Open **Settings > System > About** and look at
**System type**. It says either "x64-based processor" or "ARM-based processor".

## Using FMOS

On Windows, double-click the FMOS file. On macOS or Linux, run it from a
terminal. You will see something like this:

```text
FMOS (Fix My OS)  version 0.1.0
┌────────────────────────────────────────┐ ┌─────────────────────────┐
│ System                                 │ │ Menu                    │
│                                        │ │                         │
│ OS:     Windows 11 Pro 24H2            │ │ > OS Repair             │
│ CPU:    AMD Ryzen 7 7800X3D            │ │   Tools                 │
│ GPU:    NVIDIA GeForce RTX 4070        │ │   Network               │
│ RAM:    12.1 GB used of 32.0 GB        │ │   Logs                  │
│ Disk:   210.0 GB free of 931.5 GB (C:) │ │   Updates               │
│ Uptime: 2 days, 4 hours                │ │   Quit                  │
│ User:   you                            │ │                         │
│ Admin:  No                             │ │ Check and repair system │
│                                        │ │ files and disks.        │
└────────────────────────────────────────┘ └─────────────────────────┘
↑/↓ move • enter select • q quit
```

Everything works with the keyboard. The bottom line always shows the keys you
can use right now.

| Key | What it does |
| --- | --- |
| `↑` `↓` (or `k` `j`) | Move through a list, or scroll an action's output |
| `Enter` | Open a section, or run the selected action |
| `Esc` | Go back. While a read-only action runs, stop it |
| `PgUp` `PgDn` `Home` `End` | Jump through an action's output |
| `q` | Quit (not while an action or update is running) |

FMOS also has two command-line options:

```powershell
.\fmos-x64.exe --help       # show the options
.\fmos-x64.exe --version    # show the version number
```

## Features

### System panel

The left side of the screen shows your operating system, processor, graphics,
memory, free disk space, how long the computer has been on, who is signed in,
and whether FMOS is running as Administrator. It only reads information.

- **Windows:** the details come from Windows itself. Nothing extra is needed.
- **macOS and Linux:** the details come from
  [fastfetch](https://github.com/fastfetch-cli/fastfetch). Without it, most
  details show as "Unknown", and FMOS tells you how to install it:

  ```bash
  brew install fastfetch        # macOS
  sudo apt install fastfetch    # Ubuntu
  ```

### Network > Show network details

Lists your network adapters with their IP addresses, gateway, and DNS servers.
Use it to check whether your computer is connected to your network, or to send
the details to someone helping you.

Choose **Network**, then press Enter on **Show network details**. The output
appears live, and when it is done FMOS tells you how it went:

```text
Finished in less than a second.
The full output is in this session's log (Logs > Open logs folder).
```

| | |
| --- | --- |
| **Runs** | `ipconfig /all` on Windows, `ifconfig` on macOS, `ip address` on Linux |
| **Administrator** | Not needed |
| **Risks** | None. It only reads information and changes nothing |

### Logs

Every time FMOS opens, it saves a log of that session, for example
`fmos-2026-10-22_13-45-10.txt`. A log contains your system details and
everything FMOS ran, with the full output and how each one ended. It is plain
text, so Notepad or any text editor can open it.

**To send a log to someone helping you:** choose **Logs**, press Enter on
**Open logs folder**, and send the newest file.

Logs are saved in your own user folder, wherever the FMOS file is:

| System | Logs folder |
| --- | --- |
| Windows | `%LOCALAPPDATA%\FMOS\logs` (for example `C:\Users\you\AppData\Local\FMOS\logs`) |
| macOS | `~/Library/Logs/FMOS` |
| Linux | `~/.local/state/fmos/logs` |

If FMOS cannot save logs, it shows a red warning and keeps working. Check that
your disk has free space, then close FMOS and open it again.

### Updates

When FMOS opens, it quietly checks whether a newer version has been released.
If there is one, the top line says **Update available**. Choose **Updates** to
see what's new, then press Enter on **Update to ...**.

FMOS then:

1. Downloads the new version from this project's GitHub Releases page.
2. Checks it against the release's checksum, and refuses to use a damaged or
   altered file.
3. Replaces itself. The new version starts the next time you open FMOS.

If anything goes wrong, your current version keeps working. Being offline is
fine: FMOS says it could not check, and carries on.

| | |
| --- | --- |
| **Administrator** | Not needed, as long as FMOS is in a folder you own, such as Downloads or the Desktop. It cannot update itself inside `C:\Program Files` |
| **Internet** | Needed to check and to download |
| **Risks** | Low. It only replaces the FMOS file, and only after you choose to update |

## Safety and privacy

- **FMOS does not change your system yet.** The only file it ever replaces is
  itself, and only when you choose to update.
- **No Administrator rights needed** for anything FMOS does so far.
- **Nothing runs behind your back.** FMOS runs fixed, built-in tools by their
  full path, never through a command shell, and shows each command before you run it.
- **It only goes online to check for updates.** The check is a normal web
  request to GitHub that includes FMOS's version number. FMOS sends no
  information about your computer, and has no tracking.
- **Your logs stay on your computer.** They include details such as your user
  name, so only share them with people you trust.

## Building FMOS yourself

You only need this if you want to build FMOS from its source code. It takes
four steps: install the tools, download the code, build, and run.

### 1. Install Git and Go

FMOS needs [Git](https://git-scm.com/) to download the code, and
[Go](https://go.dev/) 1.27 or newer to build it.

**Windows (PowerShell):**

```powershell
winget install --id Git.Git -e
winget install --id GoLang.Go -e
```

Then close PowerShell and open it again, so it finds the new commands.

**macOS:** install [Homebrew](https://brew.sh/) first if you do not have it.

```bash
xcode-select --install    # installs Git; skip it if "git --version" already works
brew install go
```

**Ubuntu:**

```bash
sudo apt install git
sudo snap install go --classic    # Ubuntu's own "golang" package is too old
```

Check that it worked. This should print `go1.27` or newer:

```bash
go version
```

### 2. Download the code

This puts the code in a new `FMOS` folder inside your user folder, for example
`C:\Users\you\FMOS` or `/Users/you/FMOS`.

**Windows (PowerShell):**

```powershell
cd $HOME
git clone https://github.com/GioMonci/FMOS.git
cd FMOS
```

**macOS or Ubuntu:**

```bash
cd ~
git clone https://github.com/GioMonci/FMOS.git
cd FMOS
```

No Git? On the [GitHub page](https://github.com/GioMonci/FMOS), choose
**Code > Download ZIP**, extract it, and open a terminal in the extracted folder.

### 3. Build

Run these from the `FMOS` folder. The program is saved in its `build` folder.

**Windows (PowerShell):**

```powershell
$env:CGO_ENABLED = "0"
$env:GOARCH = "amd64"; go build -o build/fmos-x64.exe ./cmd/fmos
$env:GOARCH = "arm64"; go build -o build/fmos-arm64.exe ./cmd/fmos
Remove-Item Env:CGO_ENABLED, Env:GOARCH
```

**macOS or Ubuntu:**

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/fmos-macos ./cmd/fmos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/fmos-linux ./cmd/fmos
```

The macOS file is for Apple Silicon Macs (M1 and newer).

### 4. Run

| System | Command |
| --- | --- |
| Windows on Intel or AMD | `.\build\fmos-x64.exe` |
| Windows on ARM | `.\build\fmos-arm64.exe` |
| macOS | `./build/fmos-macos` |
| Ubuntu | `./build/fmos-linux` |

On Windows you can also double-click the file in the `build` folder.

### Getting the newest code later

From the `FMOS` folder, download the latest changes, then build again (step 3):

```bash
git pull
```

### For developers

You can skip building and run FMOS straight from the source:

```bash
go run ./cmd/fmos
```

To build every file at once for a release, including the Windows files from a
Mac, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
