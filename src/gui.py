"""
File: gui.py
Author: Gio Monci
Updated: 07/12/2026

Description:
    A simple Tkinter graphical window for FMOS (Fix My OS).

    This window is a friendly front end over the exact same safe logic the
    command line uses (operations.py, runner.py, apps.py). It does not add any
    new repair behavior of its own, so all the safety rules still apply:

        - Read-only checks never change the system.
        - Risky repairs pop up a confirmation dialog first.
        - Everything that runs is logged to the logs/ folder.

    Friendly touches in this window:
        - Messages are colour-coded (green success, amber warnings, red
          errors, bold headings) so long logs are easy to skim.
        - A moving progress bar shows when a long repair is still working.
        - A banner up top says whether FMOS has Administrator permissions.
        - Buttons to clear the log, save it to a file, or open the logs folder.

    How it stays responsive and safe:
        - Repairs and installs run on a background thread so the window never
          freezes. There is no "stop" button on purpose: interrupting a repair
          like SFC or DISM could leave Windows in a worse state.
        - All FMOS messages are routed into the text box on screen.
        - Confirmation questions are marshalled back to the main thread and
          shown as a normal yes/no dialog, then the worker waits for the answer.
"""

import os
import subprocess
import sys
import threading
import tkinter as tk
from tkinter import filedialog, messagebox, scrolledtext, ttk

import apps
import logger
import operations
import platform_check
import ui


# Colours used for the different kinds of messages in the output box.
# Chosen to stay readable on a normal light background.
TAG_COLORS = {
    "heading": "#0b3d66",
    "success": "#1a7f37",
    "warning": "#b35900",
    "error": "#c01c28",
    "step": "#5a5a5a",
}


def tag_for_line(line: str, next_line: str = "") -> str:
    """
    Decide which colour tag a single output line should use.

    This looks at the plain-text prefixes FMOS already emits (see ui.py), so
    the GUI can colour messages without the rest of the program knowing or
    caring that a GUI exists. Kept as a plain function so it can be tested
    without opening a window.

    A heading is a title line immediately followed by a line of "=" signs,
    plus that underline line itself.
    """
    stripped = line.strip()

    if line.startswith("[OK]"):
        return "success"
    if line.startswith("[!]"):
        return "warning"
    if line.startswith("[X]"):
        return "error"
    if line.startswith("-> "):
        return "step"

    # Heading underline, e.g. "======".
    if stripped and set(stripped) == {"="}:
        return "heading"
    # Heading title: the next line is an underline of "=" signs.
    next_stripped = next_line.strip()
    if next_stripped and set(next_stripped) == {"="}:
        return "heading"

    return ""  # normal text, no special colour


class Tooltip:
    """A small hover label that explains what a button does."""

    def __init__(self, widget: tk.Widget, text: str) -> None:
        self.widget = widget
        self.text = text
        self.tip: tk.Toplevel | None = None
        widget.bind("<Enter>", self._show)
        widget.bind("<Leave>", self._hide)

    def _show(self, _event: object) -> None:
        if self.tip is not None:
            return
        x = self.widget.winfo_rootx() + 16
        y = self.widget.winfo_rooty() + self.widget.winfo_height() + 4
        self.tip = tk.Toplevel(self.widget)
        self.tip.wm_overrideredirect(True)
        self.tip.wm_geometry(f"+{x}+{y}")
        tk.Label(
            self.tip,
            text=self.text,
            background="#ffffe0",
            relief="solid",
            borderwidth=1,
            justify="left",
            padx=6,
            pady=3,
            wraplength=260,
        ).pack()

    def _hide(self, _event: object) -> None:
        if self.tip is not None:
            self.tip.destroy()
            self.tip = None


class FmosApp:
    """The main FMOS window and everything it can do."""

    def __init__(self, root: tk.Tk) -> None:
        self.root = root
        self.worker: threading.Thread | None = None

        root.title("FMOS - Fix My OS")
        root.geometry("760x600")
        root.minsize(620, 480)
        root.protocol("WM_DELETE_WINDOW", self._on_close)

        # A consistent, modern-looking theme across operating systems.
        style = ttk.Style()
        if "clam" in style.theme_names():
            style.theme_use("clam")

        self._build_header()
        self._build_banner()
        self._build_buttons()
        self._build_output()
        self._build_statusbar()

        # Route all FMOS messages and confirmations into this window.
        ui.set_output(self._emit)
        ui.set_confirm(self._confirm)

        self._emit("Welcome to FMOS (Fix My OS).")
        self._emit(f"Environment: {platform_check.describe_environment()}")
        self._emit('Start with "Scan" - it only checks and changes nothing.\n')

    # ---- Window building -------------------------------------------------

    def _build_header(self) -> None:
        header = ttk.Frame(self.root, padding=(12, 10))
        header.pack(fill="x")

        ttk.Label(header, text="Fix My OS", font=("Segoe UI", 18, "bold")).pack(anchor="w")
        ttk.Label(
            header,
            text="A safe, beginner-friendly Windows repair helper.",
            font=("Segoe UI", 10),
            foreground="#555555",
        ).pack(anchor="w")

    def _build_banner(self) -> None:
        """A coloured strip that shows the Administrator / Windows status."""
        on_windows = platform_check.is_windows()
        is_admin = platform_check.is_admin()

        if not on_windows:
            text = "Not running on Windows - repairs will not actually run. You can still explore FMOS safely."
            bg, fg = "#fff3cd", "#664d03"
        elif not is_admin:
            text = 'Not running as Administrator. Repairs need it: close FMOS, right-click it, and choose "Run as administrator".'
            bg, fg = "#fff3cd", "#664d03"
        else:
            text = "Running as Administrator. FMOS can perform repairs."
            bg, fg = "#d1e7dd", "#0f5132"

        banner = tk.Label(
            self.root,
            text=text,
            background=bg,
            foreground=fg,
            anchor="w",
            padx=12,
            pady=6,
            wraplength=740,
            justify="left",
        )
        banner.pack(fill="x")

    def _build_buttons(self) -> None:
        bar = ttk.Frame(self.root, padding=(12, 6))
        bar.pack(fill="x")

        # Repair options let the user pick which repair steps to run.
        self.run_sfc = tk.BooleanVar(value=True)
        self.run_dism = tk.BooleanVar(value=True)

        actions = ttk.Frame(bar)
        actions.pack(fill="x")

        self.scan_button = ttk.Button(actions, text="Scan (safe check)", width=18, command=self.on_scan)
        self.scan_button.grid(row=0, column=0, padx=4, pady=4)
        Tooltip(self.scan_button, "Read-only health check. Looks for problems but changes nothing.")

        self.repair_button = ttk.Button(actions, text="Repair", width=18, command=self.on_repair)
        self.repair_button.grid(row=0, column=1, padx=4, pady=4)
        Tooltip(self.repair_button, "Repair Windows system files. Asks before any risky step.")

        self.apps_button = ttk.Button(actions, text="Set up apps", width=18, command=self.on_setup_apps)
        self.apps_button.grid(row=0, column=2, padx=4, pady=4)
        Tooltip(self.apps_button, "Install common apps with winget. You choose what to install.")

        self.logs_button = ttk.Button(actions, text="Show logs", width=18, command=self.on_logs)
        self.logs_button.grid(row=0, column=3, padx=4, pady=4)
        Tooltip(self.logs_button, "List the log files FMOS has saved.")

        # Checkboxes that control what "Repair" does.
        options = ttk.Frame(bar)
        options.pack(fill="x", pady=(4, 0))
        ttk.Label(options, text="Repair steps:").pack(side="left", padx=(0, 6))
        ttk.Checkbutton(options, text="System File Check (SFC)", variable=self.run_sfc).pack(side="left")
        ttk.Checkbutton(options, text="DISM RestoreHealth", variable=self.run_dism).pack(side="left")

        # These buttons act on the log view itself and are always safe to use.
        tools = ttk.Frame(bar)
        tools.pack(fill="x", pady=(6, 0))
        ttk.Button(tools, text="Clear", command=self.on_clear).pack(side="left", padx=(0, 4))
        ttk.Button(tools, text="Save log...", command=self.on_save_log).pack(side="left", padx=4)
        ttk.Button(tools, text="Open logs folder", command=self.on_open_logs_folder).pack(side="left", padx=4)

        # Only the repair/action buttons get disabled while work is running.
        self._buttons = [self.scan_button, self.repair_button, self.apps_button, self.logs_button]

    def _build_output(self) -> None:
        wrapper = ttk.Frame(self.root, padding=(12, 8))
        wrapper.pack(fill="both", expand=True)

        self.output = scrolledtext.ScrolledText(
            wrapper, wrap="word", font=("Consolas", 10), state="disabled"
        )
        self.output.pack(fill="both", expand=True)

        # Set up the colours for each kind of message.
        for name, colour in TAG_COLORS.items():
            self.output.tag_configure(name, foreground=colour)
        self.output.tag_configure("heading", font=("Consolas", 10, "bold"))

    def _build_statusbar(self) -> None:
        bar = ttk.Frame(self.root, padding=(12, 4))
        bar.pack(fill="x", side="bottom")

        self.status = ttk.Label(bar, text="Ready.", foreground="#555555")
        self.status.pack(side="left")

        # An indeterminate bar that just shows "still working" during long jobs.
        self.progress = ttk.Progressbar(bar, mode="indeterminate", length=160)
        # It is packed only while busy (see _set_busy).

    # ---- Output and confirmation (thread-safe) ---------------------------

    def _emit(self, text: str) -> None:
        """Append a line of FMOS output to the text box, from any thread."""
        self.root.after(0, self._append, text + "\n")

    def _append(self, text: str) -> None:
        """
        Insert text into the output box, colouring each line. Main thread only.
        """
        self.output.configure(state="normal")
        lines = text.split("\n")
        for index, line in enumerate(lines):
            is_last = index == len(lines) - 1
            next_line = lines[index + 1] if index + 1 < len(lines) else ""
            tag = tag_for_line(line, next_line)
            content = line if is_last else line + "\n"
            if content:
                self.output.insert("end", content, tag)
        self.output.see("end")
        self.output.configure(state="disabled")

    def _confirm(self, question: str) -> bool:
        """
        Ask a yes/no question safely from a background thread.

        Tkinter dialogs must run on the main thread, so we schedule the dialog
        there and make the worker thread wait for the answer.
        """
        answer: dict[str, bool] = {}
        answered = threading.Event()

        def ask() -> None:
            answer["value"] = messagebox.askyesno("Please confirm", question)
            answered.set()

        self.root.after(0, ask)
        answered.wait()
        return answer.get("value", False)

    # ---- Button handlers -------------------------------------------------

    def on_scan(self) -> None:
        self._run_in_background("Checking your system...", operations.scan)

    def on_repair(self) -> None:
        only_sfc = self.run_sfc.get()
        only_dism = self.run_dism.get()
        if not only_sfc and not only_dism:
            messagebox.showinfo("Nothing selected", "Pick at least one repair step first.")
            return

        # If both are checked, that is the full standard repair, so we pass
        # neither "only" flag. If exactly one is checked, run just that step.
        if only_sfc and only_dism:
            self._run_in_background("Repairing your system...", lambda: operations.repair())
        else:
            self._run_in_background(
                "Repairing your system...",
                lambda: operations.repair(only_sfc=only_sfc, only_dism=only_dism),
            )

    def on_setup_apps(self) -> None:
        AppPicker(self.root, on_install=self._install_apps)

    def on_logs(self) -> None:
        self._emit(f"\nLogs are saved in: {logger.LOG_DIR}")
        saved = logger.list_logs()
        if not saved:
            self._emit("No logs have been saved yet.")
            return
        self._emit("Most recent logs:")
        for path in saved[:15]:
            self._emit(f"  - {path.name}")

    def on_clear(self) -> None:
        """Empty the output box."""
        self.output.configure(state="normal")
        self.output.delete("1.0", "end")
        self.output.configure(state="disabled")

    def on_save_log(self) -> None:
        """Save everything currently shown in the output box to a text file."""
        text = self.output.get("1.0", "end").strip()
        if not text:
            messagebox.showinfo("Nothing to save", "The log is empty right now.")
            return
        path = filedialog.asksaveasfilename(
            title="Save FMOS output",
            defaultextension=".txt",
            filetypes=[("Text files", "*.txt"), ("All files", "*.*")],
        )
        if not path:
            return
        try:
            with open(path, "w", encoding="utf-8") as handle:
                handle.write(text + "\n")
        except OSError as problem:
            messagebox.showerror("Could not save", f"FMOS could not save the file:\n{problem}")
            return
        self._emit(f"[OK] Saved this output to {path}")

    def on_open_logs_folder(self) -> None:
        """Open the logs/ folder in the system file browser."""
        logger.LOG_DIR.mkdir(parents=True, exist_ok=True)
        path = str(logger.LOG_DIR)
        try:
            if platform_check.is_windows():
                os.startfile(path)  # type: ignore[attr-defined]
            elif sys.platform == "darwin":
                subprocess.run(["open", path], check=False)
            else:
                subprocess.run(["xdg-open", path], check=False)
        except OSError as problem:
            messagebox.showerror("Could not open folder", f"FMOS could not open the folder:\n{problem}")

    # ---- Background work -------------------------------------------------

    def _install_apps(self, selected: list) -> None:
        self._run_in_background("Installing apps...", lambda: apps.install_apps(selected))

    def _run_in_background(self, status_text: str, work) -> None:
        """
        Run work() on a background thread so the window stays responsive.

        Only one job runs at a time. The action buttons are disabled while it
        runs, and the progress bar shows that FMOS is still working.
        """
        if self.worker is not None and self.worker.is_alive():
            messagebox.showinfo("Please wait", "FMOS is still working on the last task.")
            return

        self._set_busy(True, status_text)

        def wrapped() -> None:
            try:
                work()
            except Exception as problem:  # keep the window alive on any error
                self._emit(f"[X] Something went wrong: {problem}")
            finally:
                self.root.after(0, self._set_busy, False, "Ready.")

        self.worker = threading.Thread(target=wrapped, daemon=True)
        self.worker.start()

    def _set_busy(self, busy: bool, status_text: str) -> None:
        """Enable or disable the action buttons and the progress bar."""
        state = "disabled" if busy else "normal"
        for button in self._buttons:
            button.configure(state=state)
        self.status.configure(text=status_text)

        if busy:
            self.progress.pack(side="right")
            self.progress.start(12)
        else:
            self.progress.stop()
            self.progress.pack_forget()

    def _on_close(self) -> None:
        """Ask before closing if a repair is still running."""
        if self.worker is not None and self.worker.is_alive():
            leave = messagebox.askyesno(
                "FMOS is still working",
                "A task is still running. Closing now could interrupt a repair.\n\n"
                "Close FMOS anyway?",
            )
            if not leave:
                return
        self.root.destroy()


class AppPicker(tk.Toplevel):
    """A small pop-up window that lets the user choose apps to install."""

    def __init__(self, parent: tk.Misc, on_install) -> None:
        super().__init__(parent)
        self.on_install = on_install
        self.title("Set up apps")
        self.geometry("440x460")
        self.transient(parent)

        ttk.Label(
            self,
            text="Choose apps to install (with winget):",
            font=("Segoe UI", 11, "bold"),
            padding=(12, 10),
        ).pack(anchor="w")

        self.choices: list[tuple[apps.SetupApp, tk.BooleanVar]] = []
        body = ttk.Frame(self, padding=(12, 0))
        body.pack(fill="both", expand=True)
        for app in apps.CATALOG:
            var = tk.BooleanVar(value=False)
            ttk.Checkbutton(
                body,
                text=f"{app.name}  ({app.category})",
                variable=var,
            ).pack(fill="x", anchor="w")
            self.choices.append((app, var))

        footer = ttk.Frame(self, padding=(12, 10))
        footer.pack(fill="x")
        ttk.Button(footer, text="Cancel", command=self.destroy).pack(side="right", padx=4)
        ttk.Button(footer, text="Install selected", command=self._confirm_and_close).pack(
            side="right", padx=4
        )

    def _confirm_and_close(self) -> None:
        selected = [app for app, var in self.choices if var.get()]
        if not selected:
            messagebox.showinfo("Nothing selected", "Tick at least one app, or press Cancel.")
            return
        self.destroy()
        # apps.install_apps shows its own final confirmation before installing.
        self.on_install(selected)


def launch() -> None:
    """Create the window and start the FMOS graphical interface."""
    root = tk.Tk()
    FmosApp(root)
    root.mainloop()


if __name__ == "__main__":
    launch()
