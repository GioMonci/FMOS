"""
File: logger.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Saves FMOS logs to a predictable place.

    Every repair or setup action should leave a log behind so the user (or
    someone helping them) can review exactly what happened. Logs are written
    to the project's logs/ folder as plain text files.

    A log records:
        - date and time
        - the action performed
        - the exact command
        - whether Administrator mode was used
        - the exit code
        - captured output
        - whether a restart is recommended
"""

import datetime
from pathlib import Path


# logs/ lives at the project root, one level up from this src/ folder.
LOG_DIR = Path(__file__).resolve().parent.parent / "logs"


def _timestamp_for_filename() -> str:
    """Return a filesystem-safe timestamp like 2026-07-12_140355."""
    return datetime.datetime.now().strftime("%Y-%m-%d_%H%M%S")


def _readable_timestamp() -> str:
    """Return a human-readable timestamp like 2026-07-12 14:03:55."""
    return datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def save_log(
    action: str,
    command: str,
    used_admin: bool,
    exit_code: int | None,
    output: str,
    restart_recommended: bool,
) -> Path:
    """
    Write a single log file describing one action and return its path.

    The file name includes the action name and a timestamp so logs never
    overwrite each other, for example: logs/sfc-scan-2026-07-12_140355.txt
    """
    LOG_DIR.mkdir(parents=True, exist_ok=True)

    safe_action = action.replace(" ", "-").replace("/", "-").lower()
    file_path = LOG_DIR / f"{safe_action}-{_timestamp_for_filename()}.txt"

    exit_code_text = "unknown" if exit_code is None else str(exit_code)
    restart_text = "Yes" if restart_recommended else "No"

    contents = (
        f"FMOS log\n"
        f"========\n"
        f"Date and time:       {_readable_timestamp()}\n"
        f"Action:              {action}\n"
        f"Command:             {command}\n"
        f"Ran as Administrator: {'Yes' if used_admin else 'No'}\n"
        f"Exit code:           {exit_code_text}\n"
        f"Restart recommended: {restart_text}\n"
        f"\n"
        f"--- Output ---\n"
        f"{output.strip()}\n"
    )

    file_path.write_text(contents, encoding="utf-8")
    return file_path


def list_logs() -> list[Path]:
    """Return saved log files, newest first."""
    if not LOG_DIR.exists():
        return []
    logs = [p for p in LOG_DIR.glob("*.txt")]
    return sorted(logs, key=lambda p: p.stat().st_mtime, reverse=True)
