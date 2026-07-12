"""
File: test_gui.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    Tests for the GUI's message-colouring logic.

    The colour of each line is decided by tag_for_line(), a plain function
    that reads the same text prefixes FMOS already emits. Testing it here
    keeps the colours correct without needing to open a real window.

    Run from the project root with:

        python -m unittest discover -s tests
"""

import sys
import unittest
from pathlib import Path

# Make the modules in src/ importable when running the tests.
SRC = Path(__file__).resolve().parent.parent / "src"
sys.path.insert(0, str(SRC))

from gui import tag_for_line  # noqa: E402


class TestTagForLine(unittest.TestCase):
    def test_prefixes_map_to_colours(self):
        self.assertEqual(tag_for_line("[OK] done"), "success")
        self.assertEqual(tag_for_line("[!] careful"), "warning")
        self.assertEqual(tag_for_line("[X] failed"), "error")
        self.assertEqual(tag_for_line("-> running"), "step")

    def test_plain_text_has_no_tag(self):
        self.assertEqual(tag_for_line("just some text"), "")

    def test_underline_line_is_a_heading(self):
        self.assertEqual(tag_for_line("======"), "heading")

    def test_title_followed_by_underline_is_a_heading(self):
        # A title line is a heading when the next line underlines it.
        self.assertEqual(tag_for_line("Summary", "======="), "heading")
        # But not when the next line is normal text.
        self.assertEqual(tag_for_line("Summary", "all good"), "")


if __name__ == "__main__":
    unittest.main()
