"""
File: apps.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    First-time setup: install common apps with winget.

    FMOS never installs random downloads. It uses winget (the official
    Windows package manager) and a short, hand-picked list of well-known
    apps. The user always sees the list and chooses what to install. Nothing
    is bundled or installed without being shown first.
"""

import subprocess
from dataclasses import dataclass

import logger
import platform_check
import ui


@dataclass(frozen=True)
class SetupApp:
    """One optional app the user can choose to install via winget."""

    name: str          # friendly name shown to the user
    winget_id: str     # the exact winget package id
    category: str      # what kind of app it is


# A short, boring, well-known list. Everything installs from official
# publishers through winget. Keep this list small and trustworthy.
CATALOG: list[SetupApp] = [
    SetupApp("Google Chrome", "Google.Chrome", "Browser"),
    SetupApp("Mozilla Firefox", "Mozilla.Firefox", "Browser"),
    SetupApp("7-Zip", "7zip.7zip", "Archive tool"),
    SetupApp("VLC Media Player", "VideoLAN.VLC", "Media"),
    SetupApp("Visual Studio Code", "Microsoft.VisualStudioCode", "Developer tool"),
    SetupApp("Git", "Git.Git", "Developer tool"),
    SetupApp("PowerToys", "Microsoft.PowerToys", "Hardware utility"),
]


def show_catalog() -> None:
    """Print the list of apps the user can choose from, with numbers."""
    ui.heading("Apps available to install")
    ui.info("These all install from official sources using winget.\n")
    for number, app in enumerate(CATALOG, start=1):
        print(f"  {number:>2}. {app.name}  ({app.category})")
    ui.info("")


def choose_apps() -> list[SetupApp]:
    """
    Ask the user which apps to install and return their choices.

    The user types the numbers they want, separated by spaces or commas, for
    example: 1 3 5. An empty answer means "install nothing", which is safe.
    """
    show_catalog()
    try:
        raw = input("Enter the numbers of the apps to install (or press Enter to cancel): ")
    except EOFError:
        return []

    chosen: list[SetupApp] = []
    for piece in raw.replace(",", " ").split():
        if not piece.isdigit():
            ui.warning(f"Ignoring '{piece}', which is not a number.")
            continue
        index = int(piece) - 1
        if 0 <= index < len(CATALOG):
            app = CATALOG[index]
            if app not in chosen:
                chosen.append(app)
        else:
            ui.warning(f"Ignoring '{piece}', which is not in the list.")
    return chosen


def install_apps(apps: list[SetupApp]) -> None:
    """
    Install the chosen apps one at a time, showing progress and logging results.

    Each install is confirmed as a group before anything runs, and every
    attempt (success or failure) is logged.
    """
    if not apps:
        ui.info("No apps selected. Nothing was installed.")
        return

    ui.heading("About to install")
    for app in apps:
        print(f"  - {app.name}  (winget id: {app.winget_id})")

    if not platform_check.is_windows():
        ui.warning(
            "winget is a Windows tool and FMOS is not running on Windows, so "
            "nothing was installed."
        )
        return

    if not ui.confirm("\nInstall the apps listed above?"):
        ui.info("Okay, nothing was installed.")
        return

    for app in apps:
        _install_one(app)


def _install_one(app: SetupApp) -> None:
    """Install a single app with winget and log the result."""
    args = ["winget", "install", "--id", app.winget_id, "-e"]
    ui.step(f"Installing {app.name}...")

    try:
        completed = subprocess.run(args, capture_output=True, text=True, check=False)
        exit_code: int | None = completed.returncode
        output = (completed.stdout or "") + (completed.stderr or "")
    except FileNotFoundError:
        exit_code = None
        output = (
            "Could not find 'winget'. It comes with recent versions of Windows "
            "10 and 11. Update Windows, or install 'App Installer' from the "
            "Microsoft Store, then try again."
        )
    except OSError as problem:
        exit_code = None
        output = f"Could not start winget: {problem}"

    log_path = logger.save_log(
        action=f"setup-app-{app.name}",
        command=" ".join(args),
        used_admin=platform_check.is_admin(),
        exit_code=exit_code,
        output=output,
        restart_recommended=False,
    )

    if exit_code == 0:
        ui.success(f"{app.name} installed. Log saved to {log_path}")
    else:
        ui.error(
            f"Could not install {app.name}. This can happen if the app is "
            f"already installed or winget is unavailable. The full output was "
            f"saved to {log_path}"
        )
