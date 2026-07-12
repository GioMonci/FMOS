"""
File: operations.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    The shared "what FMOS actually does" layer.

    Both the command line (main.py) and the graphical interface (gui.py) call
    into this module, so the two front ends always behave the same way and the
    safety rules live in exactly one place. Everything here talks to the user
    through ui.*, which means the CLI and GUI can each decide where those
    messages go.
"""

from commands import REPAIR_SEQUENCE, SCAN_SEQUENCE, get_command
import platform_check
import runner
import ui


def run_sequence(keys: list[str], assume_yes: bool = False) -> list[runner.RunResult]:
    """Run a list of command keys in order and collect their results."""
    results: list[runner.RunResult] = []
    for key in keys:
        command = get_command(key)
        if command is None:
            ui.error(f"Internal problem: unknown command '{key}'. Skipping it.")
            continue
        results.append(runner.run_command(command, assume_yes=assume_yes))
    return results


def summarize(results: list[runner.RunResult]) -> None:
    """Print a short, honest summary of what happened."""
    ui.heading("Summary")
    if not results:
        ui.info("Nothing was run.")
        return

    restart_needed = False
    for result in results:
        title = result.command.title
        if result.status == runner.RESULT_SUCCESS:
            ui.success(f"{title}: finished successfully.")
            if result.command.may_require_restart:
                restart_needed = True
        elif result.status == runner.RESULT_FAILED:
            ui.error(
                f"{title}: did not finish successfully "
                f"(exit code {result.exit_code}). Log: {result.log_path}"
            )
        elif result.status == runner.RESULT_SKIPPED_NO_ADMIN:
            ui.warning(f"{title}: skipped because Administrator permissions are needed.")
        elif result.status == runner.RESULT_SKIPPED_BY_USER:
            ui.info(f"{title}: skipped by you. Nothing was changed.")
        elif result.status == runner.RESULT_NOT_WINDOWS:
            ui.warning(f"{title}: skipped because FMOS is not running on Windows.")

    if restart_needed:
        ui.info("\nA restart may help finish the repairs. Save your work and restart when you can.")

    ui.info("\nFMOS cannot promise every problem is fixed. Review the logs above for details.")


def scan() -> list[runner.RunResult]:
    """Run the read-only health check. Never changes the system."""
    ui.heading("FMOS: Checking your system")
    ui.info(f"Environment: {platform_check.describe_environment()}")
    ui.info("These checks are read-only. They look for problems but do not change anything.")
    results = run_sequence(SCAN_SEQUENCE)
    summarize(results)
    return results


def repair(only_sfc: bool = False, only_dism: bool = False, assume_yes: bool = False) -> list[runner.RunResult]:
    """
    Run repair commands. Risky steps ask for confirmation first (via ui.confirm).

    Pass only_sfc or only_dism to run a single step. With neither set, the
    standard repair sequence runs.
    """
    ui.heading("FMOS: Repairing your system")
    ui.info(f"Environment: {platform_check.describe_environment()}")

    if only_sfc and not only_dism:
        keys = ["sfc"]
    elif only_dism and not only_sfc:
        keys = ["dism-restore"]
    elif only_sfc and only_dism:
        keys = ["sfc", "dism-restore"]
    else:
        keys = REPAIR_SEQUENCE

    results = run_sequence(keys, assume_yes=assume_yes)
    summarize(results)
    return results
