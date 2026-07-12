"""
File: main.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Command-line entry point for FMOS (Fix My OS).

    FMOS is a Windows repair, maintenance, and first-time setup helper. This
    file wires up the simple commands the user types and hands the real work
    off to the small, focused modules in this folder.

    Commands:
        fmos gui           Open the graphical window (also the default)
        fmos scan          Read-only health check (changes nothing)
        fmos repair        Run SFC, then DISM RestoreHealth (asks before risky steps)
        fmos repair --sfc  Run only System File Check
        fmos repair --dism Run only DISM RestoreHealth
        fmos setup-apps    Install common apps with winget
        fmos logs          Show where logs are saved and list recent ones

    Running with no command opens the graphical window.
"""

import argparse

import apps
import logger
import operations
import platform_check
import ui


def command_gui(_args: argparse.Namespace) -> int:
    """Open the Tkinter graphical window."""
    import gui  # imported here so the CLI works even without a display
    gui.launch()
    return 0


def command_scan(_args: argparse.Namespace) -> int:
    """Read-only health check. Never changes the system."""
    operations.scan()
    return 0


def command_repair(args: argparse.Namespace) -> int:
    """Run repair commands. Risky steps ask for confirmation first."""
    operations.repair(only_sfc=args.sfc, only_dism=args.dism, assume_yes=args.yes)
    return 0


def command_setup_apps(_args: argparse.Namespace) -> int:
    """Let the user pick common apps and install them with winget."""
    ui.heading("FMOS: First-time app setup")
    ui.info(f"Environment: {platform_check.describe_environment()}")
    chosen = apps.choose_apps()
    apps.install_apps(chosen)
    return 0


def command_logs(_args: argparse.Namespace) -> int:
    """Show where logs live and list the most recent ones."""
    ui.heading("FMOS logs")
    ui.info(f"Logs are saved in: {logger.LOG_DIR}")
    saved = logger.list_logs()
    if not saved:
        ui.info("No logs have been saved yet.")
        return 0
    ui.info("\nMost recent logs:")
    for path in saved[:15]:
        print(f"  - {path.name}")
    return 0


def build_parser() -> argparse.ArgumentParser:
    """Build the command-line parser with beginner-friendly help text."""
    parser = argparse.ArgumentParser(
        prog="fmos",
        description="FMOS (Fix My OS): a safe, beginner-friendly Windows repair helper.",
    )
    subparsers = parser.add_subparsers(dest="command")

    gui_cmd = subparsers.add_parser("gui", help="Open the graphical window (also the default).")
    gui_cmd.set_defaults(func=command_gui)

    scan = subparsers.add_parser("scan", help="Read-only health check (changes nothing).")
    scan.set_defaults(func=command_scan)

    repair = subparsers.add_parser(
        "repair",
        help="Repair Windows system files. Asks before risky steps.",
    )
    repair.add_argument("--sfc", action="store_true", help="Run only System File Check.")
    repair.add_argument("--dism", action="store_true", help="Run only DISM RestoreHealth.")
    repair.add_argument(
        "--yes",
        action="store_true",
        help="Skip confirmation prompts. Only use this if you understand the risks.",
    )
    repair.set_defaults(func=command_repair)

    setup = subparsers.add_parser("setup-apps", help="Install common apps with winget.")
    setup.set_defaults(func=command_setup_apps)

    logs = subparsers.add_parser("logs", help="Show where logs are saved and list recent ones.")
    logs.set_defaults(func=command_logs)

    return parser


def main() -> None:
    """Main entry point. Parse the command line and run the chosen command."""
    parser = build_parser()
    args = parser.parse_args()

    if not getattr(args, "command", None):
        # No command given: open the graphical window, since FMOS is a
        # GUI-first app. CLI users can still pass a command like "scan".
        raise SystemExit(command_gui(args))

    exit_code = args.func(args)
    raise SystemExit(exit_code)


if __name__ == "__main__":
    main()
