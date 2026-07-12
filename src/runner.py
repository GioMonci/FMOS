"""
File: runner.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Runs a single repair command safely and records what happened.

    The runner is the one place that actually executes system commands. It:

        1. Checks whether the command needs Administrator permissions.
        2. Asks the user to confirm risky commands first.
        3. Runs the command and captures its output.
        4. Saves a log.
        5. Returns a clear result the caller can summarize.

    Nothing here decides on its own to run a risky command. That choice is
    always the user's, made through ui.confirm().
"""

import subprocess
from dataclasses import dataclass

from commands import RepairCommand
import logger
import platform_check
import ui


# What happened when we tried to run a command. This keeps the reasons
# separate so the caller can explain the outcome in plain English.
RESULT_SUCCESS = "success"
RESULT_FAILED = "failed"
RESULT_SKIPPED_NO_ADMIN = "skipped_no_admin"
RESULT_SKIPPED_BY_USER = "skipped_by_user"
RESULT_NOT_WINDOWS = "not_windows"


@dataclass
class RunResult:
    """The outcome of trying to run one repair command."""

    command: RepairCommand
    status: str
    exit_code: int | None = None
    output: str = ""
    log_path: str = ""


def run_command(command: RepairCommand, assume_yes: bool = False) -> RunResult:
    """
    Run one repair command, with all the safety checks in place.

    Set assume_yes=True to skip the confirmation prompt for risky commands.
    This should only be used when the user has already clearly agreed.
    """
    ui.heading(command.title)
    ui.info(command.explanation)

    # 1. Windows-only guard. We do not pretend a Windows repair worked on
    #    another operating system.
    if not platform_check.is_windows():
        ui.warning(
            "This is a Windows repair command, and FMOS is not running on "
            "Windows right now, so it was not run."
        )
        return RunResult(command=command, status=RESULT_NOT_WINDOWS)

    # 2. Administrator check. Tell the user clearly instead of failing with a
    #    confusing "Access denied" later.
    if command.requires_admin and not platform_check.is_admin():
        ui.warning(
            "This repair needs Administrator permissions.\n"
            'Close this window, right-click FMOS, and choose "Run as administrator."'
        )
        return RunResult(command=command, status=RESULT_SKIPPED_NO_ADMIN)

    # 3. Confirmation for risky commands.
    if command.is_risky and not assume_yes:
        if command.risk_note:
            ui.warning(command.risk_note)
        if command.safer_alternative:
            ui.info(f"Safer option: {command.safer_alternative}")
        if not ui.confirm(f"Do you want to run {command.title} now?"):
            ui.info("Okay, skipping this one. Nothing was changed.")
            return RunResult(command=command, status=RESULT_SKIPPED_BY_USER)

    # 4. Run it and capture the output.
    ui.step("Running now. This may take a while. Do not close this window.")
    exit_code, output = _execute(command.args)

    # 5. Save a log no matter how it went.
    used_admin = platform_check.is_admin()
    log_path = logger.save_log(
        action=command.key,
        command=" ".join(command.args),
        used_admin=used_admin,
        exit_code=exit_code,
        output=output,
        restart_recommended=command.may_require_restart and exit_code == 0,
    )

    status = RESULT_SUCCESS if exit_code == 0 else RESULT_FAILED
    return RunResult(
        command=command,
        status=status,
        exit_code=exit_code,
        output=output,
        log_path=str(log_path),
    )


def _execute(args: list[str]) -> tuple[int | None, str]:
    """
    Actually run the command and return (exit_code, combined_output).

    Output from both stdout and stderr is captured together so the log shows
    everything Windows reported. If the program cannot be started at all, we
    return a None exit code and an explanation instead of crashing.
    """
    try:
        completed = subprocess.run(
            args,
            capture_output=True,
            text=True,
            check=False,
        )
    except FileNotFoundError:
        return None, (
            f"Could not find the program '{args[0]}'. On Windows this command "
            f"should be available by default. Make sure you are on Windows and "
            f"that the command name is spelled correctly."
        )
    except OSError as problem:
        return None, f"Could not start the command: {problem}"

    combined = (completed.stdout or "") + (completed.stderr or "")
    return completed.returncode, combined
