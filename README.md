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
| [Updates](#updates) | Checks for a new version and installs it when you choose | Yes | Yes | Checks only |

**Coming next:** system file and disk repairs (SFC, DISM, CHKDSK), installing
and updating apps (winget), driver and Microsoft Store updates, Microsoft
Defender scans, and network speed checks.

## Getting FMOS

FMOS is a single program file. There is no installer: you download the file
and double-click it. These steps are for Windows. For a Mac, see
[On a Mac](#on-a-mac). The same steps are on the
[FMOS website](https://giomonci.github.io/FMOS/).

### 1. Download the right file

Most computers need the first file. Click its name to download it.

| Download | For |
| --- | --- |
| [fmos-x64.exe](https://github.com/GioMonci/FMOS/releases/latest/download/fmos-x64.exe) | Windows on Intel or AMD (almost all desktops and laptops) |
| [fmos-arm64.exe](https://github.com/GioMonci/FMOS/releases/latest/download/fmos-arm64.exe) | Windows on ARM (for example, Snapdragon laptops) |
| [fmos-macos](https://github.com/GioMonci/FMOS/releases/latest/download/fmos-macos) | macOS on Apple Silicon (M1 and newer) |

Not sure which Windows file you need? Open **Settings > System > About** and
look at **System type**. It says either "x64-based processor" or
"ARM-based processor".

These links always give you the newest version. Every version is also listed
on the [Releases page](https://github.com/GioMonci/FMOS/releases). Each release
there also has a `SHA256SUMS.txt` file, which FMOS uses to check its own
updates. You can ignore it.

**Only download FMOS from these links.** The warnings in the next steps are
fine to get past for this file, because you know where it came from. Do not do
that for a copy of FMOS from anywhere else.

### 2. Keep the file if your browser warns you

Your browser may say that `fmos-x64.exe` "isn't commonly downloaded" and hold
the file back. This happens because FMOS is new and few people have downloaded
it. It does not mean anything harmful was found.

- **Microsoft Edge:** in the downloads list, point at the file, click the
  **...** button, and choose **Keep**. In the next window, click **Show more**,
  then **Keep anyway**.
- **Other browsers:** open the downloads list and look for a **Keep** button.

### 3. Choose where to keep it

Keep FMOS in one of your own folders. Leaving it in **Downloads** is fine, or
you can move it to your **Desktop**.

Do not put it in `C:\Program Files`. FMOS updates itself by replacing its own
file, and only Administrator programs can change files there, so FMOS could
not update itself.

To open FMOS from the Desktop while the file stays in Downloads, right-click
the file and choose **Show more options > Send to > Desktop (create
shortcut)**. On Windows 10, choose **Send to > Desktop (create shortcut)**.

### 4. Open FMOS

Double-click the file. The first time, Windows may show a blue window that
says **Windows protected your PC**. To open FMOS anyway:

1. Click **More info**.
2. Check that it says **App: fmos-x64.exe** (or `fmos-arm64.exe`).
3. Click **Run anyway**.

**Why Windows warns you:** Microsoft Defender SmartScreen warns about programs
that it has not seen many people use and that are not code-signed. Code
signing is a paid certificate that tells Windows who made a program. FMOS does
not have one yet, so Windows calls it an "Unknown publisher". If the warning
appears again, do the same steps.

FMOS opens in a text window and works with the keyboard. See
[Using FMOS](#using-fmos). If the text looks squashed, maximize the window.

**If Windows says Smart App Control blocked FMOS,** there is no Run anyway
button. Some newer Windows 11 computers block every program that is not
code-signed, so FMOS cannot run there yet. Do not turn off Smart App Control
just for FMOS. Tell the person who shared FMOS with you instead.

### 5. If FMOS asks for Administrator permissions

Some repairs need Administrator permissions. None of FMOS's current features
do, so you will not see this yet.

When a repair needs them, its review screen explains why and offers to restart
FMOS as Administrator, but only if you press **Y**. Windows then dims the
screen and asks:

```text
Do you want to allow this app from an unknown publisher to make changes to your device?

fmos-x64.exe
Publisher: Unknown
```

- **Yes:** FMOS closes and opens again as Administrator, on the same review
  screen. Nothing runs until you press Y there.
- **No:** FMOS keeps running without Administrator permissions and tells you
  which repair it could not run.
- If your Windows account is not an administrator, Windows asks for an
  administrator's password instead. If you do not know it, click **No**.

It says "Unknown publisher" for the same reason as the warning in step 4:
FMOS is not code-signed yet. FMOS never asks for Administrator permissions on
its own. It only asks when you press Y on a review screen.

### Updating FMOS

You never need to download FMOS again. When a new version comes out, the top
line of FMOS says **Update available**. Choose **Updates**, read what's new,
and press Enter on **Update to ...**. When it is done, quit FMOS (press **q**)
and open it again. See [Updates](#updates) for how it works.

After an update, you may see a file named `fmos-x64.exe.old` next to FMOS.
That is the previous version. FMOS deletes it the next time it opens.

### Sending your logs for help

If something goes wrong, or someone is helping you from far away, send them
your FMOS log:

1. In FMOS, choose **Logs**, then press Enter on **Open logs folder**.
2. Quit FMOS (press **q**), so the log is complete.
3. Send the newest file. Its name has the date and time, for example
   `fmos-2026-10-22_13-45-10.txt`.

Logs include details such as your Windows user name, so only send them to
people you trust. FMOS never sends logs anywhere by itself.

### Removing FMOS

1. Delete the FMOS file, and `fmos-x64.exe.old` if there is one. Delete your
   Desktop shortcut too, if you made one.
2. To delete the logs as well, press **Windows key + R**, type
   `%LOCALAPPDATA%`, press Enter, and delete the **FMOS** folder there.

That is everything. FMOS adds nothing else to your computer: no Start menu
entries, no settings, and nothing that runs in the background.

### On a Mac

Download `fmos-macos` from the table in step 1. A downloaded program is
blocked until you allow it. In Terminal:

```bash
cd ~/Downloads
chmod +x fmos-macos      # allow it to run as a program
xattr -c fmos-macos      # remove the "downloaded from the internet" block
./fmos-macos
```

Updates work the same way as on Windows. To remove FMOS, delete `fmos-macos`,
and delete the `~/Library/Logs/FMOS` folder for the logs.

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

**Before anything that changes your computer,** FMOS shows a review screen: what
the action does, the exact command, whether it needs Administrator rights,
internet, or a restart. Nothing happens until you press **Y**. Esc goes back.
If an action needs Administrator rights, the same screen offers to restart FMOS
as Administrator. Windows asks you to confirm, and FMOS reopens on that screen.

FMOS also has a few command-line options:

```powershell
.\fmos-x64.exe --help                     # show the options
.\fmos-x64.exe --version                  # show the version number
.\fmos-x64.exe --open network-details     # open one action's review screen
```

FMOS uses `--open` itself when it restarts as Administrator. Even then,
nothing runs until you press Y.

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

1. Downloads the new version from this project's GitHub Releases page, and
   nowhere else.
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
  itself, and only when you choose to update. When repairs arrive, each one
  shows a review screen first and waits for you to press Y.
- **No Administrator rights needed** for anything FMOS does so far. FMOS never
  asks for them on its own. It only restarts as Administrator when you choose to.
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

Run these from the `FMOS` folder. The program is saved in the `build\dev` folder.

**Windows (PowerShell):**

```powershell
$env:CGO_ENABLED = "0"
$env:GOARCH = "amd64"; go build -o build/dev/fmos-x64.exe ./cmd/fmos
$env:GOARCH = "arm64"; go build -o build/dev/fmos-arm64.exe ./cmd/fmos
Remove-Item Env:CGO_ENABLED, Env:GOARCH
```

**macOS or Ubuntu:**

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/dev/fmos-macos ./cmd/fmos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/dev/fmos-linux ./cmd/fmos
```

The macOS file is for Apple Silicon Macs (M1 and newer).

### 4. Run

| System | Command |
| --- | --- |
| Windows on Intel or AMD | `.\build\dev\fmos-x64.exe` |
| Windows on ARM | `.\build\dev\fmos-arm64.exe` |
| macOS | `./build/dev/fmos-macos` |
| Ubuntu | `./build/dev/fmos-linux` |

On Windows you can also double-click the file in the `build\dev` folder.

### Development builds and release builds

The `build` folder has two kinds of builds, kept in separate folders so they
never get mixed up:

| Folder | Made by | Version shown | What it is for |
| --- | --- | --- | --- |
| `build/dev/` | The build commands above | `0.1.0-dev` | Trying out changes on your own computer |
| `build/release/v0.1.0/` | `scripts/release.sh 0.1.0` (Mac) or `scripts\release.ps1 0.1.0` (Windows) | `0.1.0` | Uploading to GitHub Releases, and nothing else |

Release builds have the same file names as development builds, because the
updater and the download links look for those exact names. To tell them apart:

- **The folder.** Only upload files from a `build/release/vX.Y.Z/` folder.
- **The version.** Run the program with `--version`. Release builds say
  `FMOS 0.1.0`; development builds say `FMOS 0.1.0-dev`.
- **The fingerprint.** `SHA256SUMS.txt` in the release folder lists each
  file's fingerprint, and the release script prints the first few characters.

The release script also strips this computer's folder paths out of the files
(`-trimpath`), and stops with an error if any are left. See
[CONTRIBUTING.md](CONTRIBUTING.md) for the full release steps.

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

To make a release, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

FMOS is licensed under the [GNU General Public License v3.0](LICENSE).
