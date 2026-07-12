"""
File: test_runner.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Tests for the command runner's safety behavior.

    These focus on the checks that protect the user, without actually running
    any real Windows repair commands. They confirm that the runner refuses to
    run commands when it should, and that it never claims success on a system
    where the command cannot run.

    Run from the project root with:

        python -m unittest discover -s tests
"""

import sys
import unittest
from pathlib import Path
from unittest import mock

# Make the modules in src/ importable when running the tests.
SRC = Path(__file__).resolve().parent.parent / "src"
sys.path.insert(0, str(SRC))

import commands  # noqa: E402
import runner  # noqa: E402


class TestRunnerSafety(unittest.TestCase):
    def test_skips_when_not_windows(self):
        command = commands.get_command("dism-check")
        with mock.patch("platform_check.is_windows", return_value=False):
            result = runner.run_command(command)
        self.assertEqual(result.status, runner.RESULT_NOT_WINDOWS)
        self.assertIsNone(result.exit_code)

    def test_skips_when_admin_required_but_missing(self):
        command = commands.get_command("sfc")  # requires admin
        with mock.patch("platform_check.is_windows", return_value=True), \
             mock.patch("platform_check.is_admin", return_value=False):
            result = runner.run_command(command)
        self.assertEqual(result.status, runner.RESULT_SKIPPED_NO_ADMIN)

    def test_risky_command_skipped_when_user_declines(self):
        command = commands.get_command("dism-restore")  # risky
        with mock.patch("platform_check.is_windows", return_value=True), \
             mock.patch("platform_check.is_admin", return_value=True), \
             mock.patch("ui.confirm", return_value=False):
            result = runner.run_command(command)
        self.assertEqual(result.status, runner.RESULT_SKIPPED_BY_USER)

    def test_successful_run_saves_a_log(self):
        command = commands.get_command("dism-check")

        fake_completed = mock.Mock()
        fake_completed.returncode = 0
        fake_completed.stdout = "No component store corruption detected."
        fake_completed.stderr = ""

        with mock.patch("platform_check.is_windows", return_value=True), \
             mock.patch("platform_check.is_admin", return_value=True), \
             mock.patch("subprocess.run", return_value=fake_completed), \
             mock.patch("logger.save_log", return_value=Path("logs/fake.txt")) as save_log:
            result = runner.run_command(command)

        self.assertEqual(result.status, runner.RESULT_SUCCESS)
        self.assertEqual(result.exit_code, 0)
        save_log.assert_called_once()


if __name__ == "__main__":
    unittest.main()
