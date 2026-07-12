"""
File: commands.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    The list of Windows repair commands FMOS knows about.

    This file only DESCRIBES commands. It does not run anything. Keeping the
    description separate from the running code makes it easy to see, at a
    glance, exactly what each command does, whether it needs Administrator
    permissions, and whether it is risky.

    Each command is a RepairCommand with plain-English fields so the rest of
    the program (and any human reading this file) can understand it without
    guessing.
"""

from dataclasses import dataclass, field


@dataclass(frozen=True)
class RepairCommand:
    """A single Windows repair command and everything FMOS needs to know about it."""

    key: str                       # short name used on the command line, e.g. "sfc"
    title: str                     # friendly title shown to the user
    explanation: str               # plain-English description of what it does
    args: list[str]                # the actual command, as a list of arguments
    requires_admin: bool           # does it need "Run as administrator"?
    is_risky: bool = False         # should FMOS confirm before running it?
    may_require_restart: bool = False
    safer_alternative: str = ""    # optional hint at a safer option
    risk_note: str = field(default="")  # shown before running a risky command


# NOTE ON /f AND /r:
# The AGENTS.md safety rules call out `chkdsk /f` and `chkdsk /r` as risky.
# FMOS deliberately offers only a READ-ONLY disk check by default. It reports
# problems and tells the user how to schedule a full repair themselves, rather
# than scheduling a disk-locking, restart-requiring repair automatically.

_COMMANDS: dict[str, RepairCommand] = {
    "sfc": RepairCommand(
        key="sfc",
        title="System File Check (SFC)",
        explanation=(
            "Scans your protected Windows system files and repairs any that "
            "are missing or damaged, using a known-good copy from Windows."
        ),
        args=["sfc", "/scannow"],
        requires_admin=True,
        may_require_restart=True,
    ),
    "dism-check": RepairCommand(
        key="dism-check",
        title="DISM Check Health",
        explanation=(
            "Quickly checks whether the Windows image has been flagged as "
            "damaged. This is a fast, read-only check that changes nothing."
        ),
        args=["DISM", "/Online", "/Cleanup-Image", "/CheckHealth"],
        requires_admin=True,
    ),
    "dism-scan": RepairCommand(
        key="dism-scan",
        title="DISM Scan Health",
        explanation=(
            "Does a deeper, read-only scan of the Windows image for damage. "
            "It does not repair anything, but it can take a while."
        ),
        args=["DISM", "/Online", "/Cleanup-Image", "/ScanHealth"],
        requires_admin=True,
    ),
    "dism-restore": RepairCommand(
        key="dism-restore",
        title="DISM Restore Health",
        explanation=(
            "Repairs the Windows image by downloading and replacing damaged "
            "files. This can change system files and needs a working internet "
            "connection or Windows Update to be available."
        ),
        args=["DISM", "/Online", "/Cleanup-Image", "/RestoreHealth"],
        requires_admin=True,
        is_risky=True,
        may_require_restart=True,
        risk_note=(
            "This repairs core Windows files. It is usually safe, but it does "
            "change the system and can take a long time. Do not turn off your "
            "PC while it runs."
        ),
        safer_alternative="Run 'fmos scan' first to check whether a repair is even needed.",
    ),
    "chkdsk": RepairCommand(
        key="chkdsk",
        title="Check Disk (read-only)",
        explanation=(
            "Checks your main drive for errors WITHOUT trying to fix them. "
            "This is the safe version: it only reports problems."
        ),
        # No /f and no /r on purpose. Those lock the drive and require a
        # restart. FMOS reports issues and lets the user decide.
        args=["chkdsk"],
        requires_admin=True,
        safer_alternative=(
            "If chkdsk reports errors, you can schedule a full repair yourself "
            "with 'chkdsk /f' from an Administrator terminal, then restart."
        ),
    ),
}


def all_commands() -> list[RepairCommand]:
    """Return every known repair command in a stable order."""
    return list(_COMMANDS.values())


def get_command(key: str) -> RepairCommand | None:
    """Return the command with this key, or None if it is not known."""
    return _COMMANDS.get(key)


# The commands used by a plain 'fmos scan' (all read-only, no changes).
SCAN_SEQUENCE = ["dism-check", "chkdsk"]

# The commands used by a plain 'fmos repair' (in this order).
REPAIR_SEQUENCE = ["sfc", "dism-restore"]
