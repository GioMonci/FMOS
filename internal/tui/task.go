package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/GioMonci/FMOS/internal/actions"
	"github.com/GioMonci/FMOS/internal/platform"
	"github.com/GioMonci/FMOS/internal/runner"
)

// task is an action that is running or has finished.
type task struct {
	action   actions.Action
	started  time.Time
	now      time.Time          // refreshed every second while running
	lines    []string           // output so far
	events   chan tea.Msg       // progress from the background run
	stop     context.CancelFunc // stops a read-only action early
	result   *runner.Result     // nil while running
	scrollUp int                // lines scrolled up from the newest output
	notice   string             // short message, such as "Stopping..."
}

// Messages from a running task.
type (
	outputLineMsg string
	taskDoneMsg   runner.Result
	tickMsg       time.Time
)

// startAction checks that an action may run, then runs it in the background
// and switches to its screen.
func (m model) startAction(action actions.Action) (tea.Model, tea.Cmd) {
	if action.ChangesSystem && !m.logsAvailable() {
		m.notice = "This action changes your system, so FMOS only runs it while " +
			"it can save a log. " + logProblem(m.logErr)
		return m, nil
	}
	if action.NeedsAdmin && !platform.IsAdmin() {
		m.notice = "This action needs Administrator permissions. Close FMOS, " +
			`right-click it, and choose "Run as administrator".`
		return m, nil
	}

	runCtx, stop := context.WithCancel(context.Background())
	if action.ChangesSystem {
		runCtx = context.Background() // never cut off a repair halfway
	}
	t := &task{
		action:  action,
		started: time.Now(),
		now:     time.Now(),
		events:  make(chan tea.Msg, 100),
		stop:    stop,
	}
	go func() {
		defer stop()
		result := runner.Run(runCtx, action.Command(), func(line string) {
			t.events <- outputLineMsg(line)
		})
		t.events <- taskDoneMsg(result)
	}()

	m.task = t
	m.screen = taskScreen
	m.notice = ""
	return m, tea.Batch(waitForTaskEvent(t.events), tickEverySecond())
}

// waitForTaskEvent waits for the next message from the running task.
func waitForTaskEvent(events <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg { return <-events }
}

// tickEverySecond refreshes the elapsed time while a task runs.
func tickEverySecond() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg { return tickMsg(now) })
}

// updateTask handles output, the elapsed-time tick, and the final result.
func (m model) updateTask(msg tea.Msg) (tea.Model, tea.Cmd) {
	t := m.task
	if t == nil {
		return m, nil
	}
	switch msg := msg.(type) {
	case outputLineMsg:
		t.lines = append(t.lines, string(msg))
		return m, waitForTaskEvent(t.events)
	case tickMsg:
		if t.result == nil {
			t.now = time.Time(msg)
			return m, tickEverySecond()
		}
	case taskDoneMsg:
		result := runner.Result(msg)
		t.result = &result
		t.notice = ""
		if m.log != nil {
			if err := m.log.LogAction(t.action.Name, result); err != nil {
				m.logErr = err
			}
		}
	}
	return m, nil
}

// handleTaskKey handles keys on the task screen. While an action runs, FMOS
// stays on this screen so every run finishes and is logged. Read-only actions
// can be stopped with Esc; system-changing ones cannot.
func (m model) handleTaskKey(key string) (tea.Model, tea.Cmd) {
	t := m.task
	if t.result == nil {
		switch {
		case t.action.ChangesSystem && (key == "esc" || key == "q" || key == "ctrl+c"):
			t.notice = "Please wait. Stopping this halfway could leave your system in a worse state."
		case key == "esc":
			t.stop()
			t.notice = "Stopping..."
		case key == "q" || key == "ctrl+c":
			t.notice = "Press Esc to stop it before quitting."
		}
		return m, nil
	}

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = sectionScreen
	case "up", "k":
		t.scrollUp++
	case "down", "j":
		t.scrollUp--
	case "pgup":
		t.scrollUp += m.outputHeight()
	case "pgdown":
		t.scrollUp -= m.outputHeight()
	case "home":
		t.scrollUp = len(t.lines)
	case "end":
		t.scrollUp = 0
	}
	t.scrollUp = max(min(t.scrollUp, m.maxScrollUp()), 0)
	return m, nil
}

// renderTask draws the full-width screen for a running or finished action.
func (m model) renderTask(width int) string {
	contentWidth := width - panelStyle.GetHorizontalFrameSize()
	body := m.statusBlock(width) + "\n\n" + m.visibleOutput(contentWidth)
	return panel(m.task.action.Name, body, width, 0)
}

// statusBlock is the task's status text, wrapped to fit inside the panel.
func (m model) statusBlock(width int) string {
	contentWidth := width - panelStyle.GetHorizontalFrameSize()
	return lipgloss.NewStyle().Width(contentWidth).Render(m.taskStatus())
}

// taskStatus says what is happening, or how the run ended, in plain English.
func (m model) taskStatus() string {
	t := m.task
	var status string
	if t.result == nil {
		status = fmt.Sprintf("Running %s (%s so far)...",
			shortCommand(t.action.Command()), formatElapsed(t.now.Sub(t.started)))
		if t.action.ChangesSystem {
			status += "\nThis may take a while. Do not close this window."
		}
	} else {
		status = resultSummary(*t.result)
		if m.logsAvailable() {
			status += "\nThe full output is in this session's log (Logs > Open logs folder)."
		}
	}
	if t.notice != "" {
		status += "\n" + t.notice
	}
	return status
}

// resultSummary explains in plain English how a run ended.
func resultSummary(result runner.Result) string {
	took := formatElapsed(result.Duration)
	switch result.Outcome {
	case runner.Finished:
		return "Finished in " + took + "."
	case runner.Failed:
		return fmt.Sprintf("Did not finish successfully. It stopped with exit code %d after %s.",
			result.ExitCode, took)
	case runner.NotFound:
		return "FMOS could not find the program it needs: " + result.Err.Error() + "."
	case runner.FailedToStart:
		return "The program could not start: " + result.Err.Error() + "."
	case runner.Stopped:
		return "Stopped after " + took + ", as you asked."
	default:
		return "Ended in an unexpected way after " + took + "."
	}
}

// visibleOutput returns the output lines that fit on the screen: the newest
// ones, or older ones after scrolling up. Long lines are cut at the edge.
func (m model) visibleOutput(contentWidth int) string {
	t := m.task
	if len(t.lines) == 0 {
		if t.result == nil {
			return mutedStyle.Render("Waiting for output...")
		}
		return mutedStyle.Render("The program printed nothing.")
	}
	start, end := m.outputRange()
	shown := strings.Join(t.lines[start:end], "\n")
	return lipgloss.NewStyle().MaxWidth(contentWidth).Render(shown)
}

// outputRange returns which output lines are on screen, as a slice range.
func (m model) outputRange() (start, end int) {
	end = len(m.task.lines) - min(m.task.scrollUp, m.maxScrollUp())
	start = max(end-m.outputHeight(), 0)
	return start, end
}

// maxScrollUp is how far the output can scroll up before reaching the top.
func (m model) maxScrollUp() int {
	return max(len(m.task.lines)-m.outputHeight(), 0)
}

// outputHeight is how many output lines fit on the task screen.
func (m model) outputHeight() int {
	width := min(m.width, maxWidth)
	used := 1 + 1 // header and key hints
	used += 2     // panel border
	used += 2     // title and the blank line under it
	used += 1     // blank line between the status and the output
	used += lipgloss.Height(m.statusBlock(width))
	if m.logErr != nil {
		used += lipgloss.Height(m.warning(width))
	}
	return max(m.height-used, 3)
}

func (m model) taskKeyHints() string {
	t := m.task
	switch {
	case t.result == nil && t.action.ChangesSystem:
		return "please wait until it finishes"
	case t.result == nil:
		return "esc stop"
	case len(t.lines) > m.outputHeight():
		start, end := m.outputRange()
		return fmt.Sprintf("↑/↓ scroll (lines %d–%d of %d) • esc back • q quit",
			start+1, end, len(t.lines))
	default:
		return "esc back • q quit"
	}
}

// shortCommand shows a command without its folder, such as "ipconfig.exe /all".
func shortCommand(command runner.Command) string {
	command.Program = filepath.Base(command.Program)
	return command.String()
}

// formatElapsed returns text like "less than a second", "12s", or "3m 05s".
func formatElapsed(elapsed time.Duration) string {
	switch {
	case elapsed < time.Second:
		return "less than a second"
	case elapsed < time.Minute:
		return fmt.Sprintf("%ds", int(elapsed.Seconds()))
	default:
		return fmt.Sprintf("%dm %02ds", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
	}
}
