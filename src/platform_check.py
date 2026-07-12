"""
File: platform_check.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Small helpers for checking the environment FMOS is running in.

    FMOS is built for Windows. These helpers let the rest of the program
    answer two simple questions safely:

        1. Are we actually on Windows?
        2. Are we running with Administrator permissions?

    Everything here is read-only. Nothing in this file changes the system.
"""

import ctypes
import os
import platform


def is_windows() -> bool:
    """Return True if FMOS is running on Windows."""
    return platform.system() == "Windows"


def is_admin() -> bool:
    """
    Return True if the current process has Administrator permissions.

    On Windows we ask the operating system directly. On other systems we
    fall back to checking for a user id of 0 (root), which lets the code
    be tested on Linux/macOS without pretending to be Windows.
    """
    if is_windows():
        try:
            return bool(ctypes.windll.shell32.IsUserAnAdmin())
        except Exception:
            # If the check itself fails, assume we are NOT admin. It is
            # safer to under-claim permissions than to over-claim them.
            return False

    # Non-Windows fallback, mainly useful for local development and tests.
    if hasattr(os, "geteuid"):
        return os.geteuid() == 0
    return False


def describe_environment() -> str:
    """Return a short, plain-English description of where FMOS is running."""
    system = platform.system() or "Unknown"
    release = platform.release() or ""
    admin = "with Administrator permissions" if is_admin() else "without Administrator permissions"
    return f"{system} {release} ({admin})".strip()
