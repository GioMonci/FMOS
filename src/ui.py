"""
File: ui.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Beginner-friendly messaging helpers for FMOS.

    The goal is to talk to the user like a helpful tech friend: plain
    English, clear steps, and honest summaries. All user-facing text should
    go through these helpers so the tone stays consistent.

    By default these helpers print to the console for the command-line tool.
    The graphical interface (gui.py) can redirect them by calling
    set_output() and set_confirm(), so the same repair logic in runner.py and
    apps.py works for both the CLI and the GUI without any changes.
"""

from typing import Callable


# These two hooks decide WHERE messages go and HOW questions are asked.
# The CLI leaves them at the defaults. The GUI replaces them so messages
# appear in a window and confirmations use a pop-up dialog.
_output: Callable[[str], None] = print
_asker: Callable[[str], bool] = None  # set below to _console_confirm


def set_output(func: Callable[[str], None]) -> None:
    """Send all FMOS messages to func(text) instead of printing them."""
    global _output
    _output = func


def set_confirm(func: Callable[[str], bool]) -> None:
    """Ask yes/no questions through func(question) -> bool instead of the console."""
    global _asker
    _asker = func


def reset() -> None:
    """Restore the default console behavior (used mainly by tests)."""
    global _output, _asker
    _output = print
    _asker = _console_confirm


def heading(text: str) -> None:
    """Show a section heading so the user can see where they are."""
    line = "=" * len(text)
    _output(f"\n{text}\n{line}")


def info(text: str) -> None:
    """Show a normal informational message."""
    _output(text)


def step(text: str) -> None:
    """Show a progress step, prefixed so it stands out."""
    _output(f"-> {text}")


def success(text: str) -> None:
    """Show a success message."""
    _output(f"[OK] {text}")


def warning(text: str) -> None:
    """Show a warning the user should notice but that is not fatal."""
    _output(f"[!] {text}")


def error(text: str) -> None:
    """Show an error message in plain English."""
    _output(f"[X] {text}")


def confirm(question: str) -> bool:
    """
    Ask the user a yes/no question and return their answer.

    Risky actions must call this first. The answer defaults to "no" so that a
    system-changing command is never triggered by accident.
    """
    return _asker(question)


def _console_confirm(question: str) -> bool:
    """Default confirmation: ask on the console. Empty answer means "no"."""
    prompt = f"{question} [y/N]: "
    try:
        answer = input(prompt).strip().lower()
    except EOFError:
        # No interactive input available (for example, output is piped).
        # Treat that as "no" to stay on the safe side.
        return False
    return answer in ("y", "yes")


# Start with the safe console default.
_asker = _console_confirm
