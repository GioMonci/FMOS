"""
File: test_ui.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Tests for the pluggable output and confirmation hooks in ui.py.

    These hooks are what let the GUI reuse the command-line repair logic:
    the GUI redirects messages into its window and confirmations into a
    pop-up dialog. If these break, the GUI and CLI could drift apart.

    Run from the project root with:

        python -m unittest discover -s tests
"""

import sys
import unittest
from pathlib import Path

# Make the modules in src/ importable when running the tests.
SRC = Path(__file__).resolve().parent.parent / "src"
sys.path.insert(0, str(SRC))

import ui  # noqa: E402


class TestUiHooks(unittest.TestCase):
    def tearDown(self):
        # Always put ui back to its console defaults after each test.
        ui.reset()

    def test_output_can_be_redirected(self):
        captured: list[str] = []
        ui.set_output(captured.append)

        ui.info("hello")
        ui.success("done")

        self.assertIn("hello", captured)
        self.assertIn("[OK] done", captured)

    def test_confirm_can_be_redirected(self):
        ui.set_confirm(lambda question: True)
        self.assertTrue(ui.confirm("Proceed?"))

        ui.set_confirm(lambda question: False)
        self.assertFalse(ui.confirm("Proceed?"))

    def test_reset_restores_console_default(self):
        ui.set_output(lambda text: None)
        ui.reset()
        # After reset the default output is the built-in print function.
        self.assertIs(ui._output, print)


if __name__ == "__main__":
    unittest.main()
