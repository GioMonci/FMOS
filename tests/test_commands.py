"""
File: test_commands.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Tests for the repair command registry.

    These check the safety promises FMOS makes about its commands, such as
    "the read-only scan never changes anything" and "chkdsk does not use the
    dangerous /f or /r flags automatically."

    Run from the project root with:

        python -m unittest discover -s tests
"""

import sys
import unittest
from pathlib import Path

# Make the modules in src/ importable when running the tests.
SRC = Path(__file__).resolve().parent.parent / "src"
sys.path.insert(0, str(SRC))

import commands  # noqa: E402


class TestCommandRegistry(unittest.TestCase):
    def test_all_commands_have_required_fields(self):
        for command in commands.all_commands():
            self.assertTrue(command.key, "command needs a key")
            self.assertTrue(command.title, "command needs a title")
            self.assertTrue(command.explanation, "command needs an explanation")
            self.assertTrue(command.args, "command needs args to run")

    def test_get_command_returns_none_for_unknown(self):
        self.assertIsNone(commands.get_command("does-not-exist"))

    def test_scan_sequence_is_read_only(self):
        # Every command used by 'fmos scan' must be marked as not risky.
        for key in commands.SCAN_SEQUENCE:
            command = commands.get_command(key)
            self.assertIsNotNone(command, f"scan uses unknown command '{key}'")
            self.assertFalse(command.is_risky, f"scan step '{key}' should not be risky")

    def test_chkdsk_never_uses_f_or_r_automatically(self):
        # The safety rules call out chkdsk /f and /r as dangerous.
        chkdsk = commands.get_command("chkdsk")
        self.assertIsNotNone(chkdsk)
        self.assertNotIn("/f", chkdsk.args)
        self.assertNotIn("/r", chkdsk.args)

    def test_risky_commands_explain_the_risk(self):
        # Any risky command must give the user a reason before running.
        for command in commands.all_commands():
            if command.is_risky:
                self.assertTrue(
                    command.risk_note,
                    f"risky command '{command.key}' should have a risk_note",
                )


if __name__ == "__main__":
    unittest.main()
