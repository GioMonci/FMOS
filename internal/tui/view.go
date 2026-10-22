package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/GioMonci/FMOS/internal/actions"
)

const (
	// maxWidth keeps the layout readable on very wide terminals.
	maxWidth = 100
	// sideBySideMinWidth is the narrowest terminal that fits both panels
	// next to each other. Narrower terminals stack them instead.
	sideBySideMinWidth = 72
)

// The styles use basic ANSI colors so they follow the terminal's own color
// scheme. Plain box borders are used because older Windows console fonts
// draw them correctly.
var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.BrightBlack).
			Padding(0, 1)
	titleStyle    = lipgloss.NewStyle().Bold(true)
	selectedStyle = lipgloss.NewStyle().Bold(true).Reverse(true)
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	warningStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Red)
)

// labelWidth fits the longest System panel label ("Uptime:") plus a space.
const labelWidth = 8

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.WindowTitle = "FMOS"
	return v
}

// render draws the whole screen: header, panels (or a running task), a
// warning when logs cannot be saved, and key hints.
func (m model) render() string {
	width := min(m.width, maxWidth)
	header := titleStyle.Render("FMOS (Fix My OS)") + mutedStyle.Render("  version "+m.version)
	if badge := m.updateBadge(); badge != "" {
		header += "  •  " + titleStyle.Render(badge)
	}

	body := m.renderPanels(width)
	switch m.screen {
	case taskScreen:
		body = m.renderTask(width)
	case reviewScreen:
		body = m.renderReview(width)
	}
	parts := []string{header, body}
	if m.logErr != nil {
		parts = append(parts, m.warning(width))
	}
	parts = append(parts, mutedStyle.Render(m.keyHints()))
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// warning explains that logs cannot be saved.
func (m model) warning(width int) string {
	return warningStyle.Width(width).Render("Warning: " + logProblem(m.logErr))
}

// renderPanels places the system panel and the menu side by side, or stacks
// them when the terminal is too narrow.
func (m model) renderPanels(width int) string {
	rightTitle, rightBody := m.rightPanelContent()
	frame := panelStyle.GetHorizontalFrameSize() // border and padding

	if width < sideBySideMinWidth {
		return lipgloss.JoinVertical(lipgloss.Left,
			panel("System", m.systemPanelBody(width-frame), width, 0),
			panel(rightTitle, rightBody, width, 0),
		)
	}

	systemWidth := width * 3 / 5
	menuWidth := width - systemWidth - 1 // one column of space between panels
	systemBody := m.systemPanelBody(systemWidth - frame)

	// Measure both panels first so they can be drawn with the same height.
	height := max(
		lipgloss.Height(panel("System", systemBody, systemWidth, 0)),
		lipgloss.Height(panel(rightTitle, rightBody, menuWidth, 0)),
	)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		panel("System", systemBody, systemWidth, height),
		" ",
		panel(rightTitle, rightBody, menuWidth, height),
	)
}

// panel draws a bordered box with a bold title. The width includes the
// border. A height of 0 lets the box grow to fit its content.
func panel(title, body string, width, height int) string {
	content := titleStyle.Render(title) + "\n\n" + body
	return panelStyle.Width(width).Height(height).Render(content)
}

// rightPanelContent returns the menu, or the open section.
func (m model) rightPanelContent() (title, body string) {
	selected := m.section()
	if m.screen == sectionScreen {
		return selected.title, m.sectionBody()
	}

	titles := make([]string, 0, len(menuItems))
	for _, item := range menuItems {
		titles = append(titles, item.title)
	}
	return "Menu", choiceList(titles, m.cursor) + "\n\n" + mutedStyle.Render(selected.description)
}

// choiceList draws a list where the selected entry is highlighted and marked
// with ">", so it does not rely on color alone.
func choiceList(entries []string, selected int) string {
	lines := make([]string, 0, len(entries))
	for i, entry := range entries {
		if i == selected {
			lines = append(lines, selectedStyle.Render("> "+entry))
		} else {
			lines = append(lines, "  "+entry)
		}
	}
	return strings.Join(lines, "\n")
}

// sectionBody shows an open section: the Logs screen, the section's actions,
// or a "Coming soon" note.
func (m model) sectionBody() string {
	section := m.section()
	switch section.id {
	case "logs":
		return m.logsBody()
	case "updates":
		return m.updatesBody()
	}
	sectionActions := actions.InSection(section.id)
	if len(sectionActions) == 0 {
		return section.description + "\n\n" +
			mutedStyle.Render("Coming soon. This part of FMOS has not been built yet.")
	}

	names := make([]string, 0, len(sectionActions))
	for _, action := range sectionActions {
		names = append(names, action.Name)
	}
	chosen := sectionActions[m.actionCursor]
	body := choiceList(names, m.actionCursor) + "\n\n" + chosen.Description + "\n\n" +
		mutedStyle.Render("Runs: "+shortCommand(chosen.Command()))
	if m.notice != "" {
		body += "\n\n" + m.notice
	}
	return body
}

// logsBody shows where this session's log is being saved, with a button to
// open the logs folder.
func (m model) logsBody() string {
	if !m.logsAvailable() {
		return logProblem(m.logErr)
	}
	body := "This session is being saved to:\n" + m.log.Path() + "\n\n" +
		selectedStyle.Render("> Open logs folder")
	if m.notice != "" {
		body += "\n\n" + mutedStyle.Render(m.notice)
	}
	return body
}

// logProblem explains why logs cannot be saved and what to do about it.
func logProblem(err error) string {
	reason := "the reason is unknown"
	if err != nil {
		reason = err.Error()
	}
	return "FMOS cannot save logs (" + reason + "). Check that your disk has free " +
		"space, then close FMOS and open it again."
}

func (m model) keyHints() string {
	switch {
	case m.screen == taskScreen:
		return m.taskKeyHints()
	case m.screen == reviewScreen:
		return m.reviewKeyHints()
	case m.screen == sectionScreen && m.section().id == "logs" && m.logsAvailable():
		return "enter open logs folder • esc back • q quit"
	case m.screen == sectionScreen && m.section().id == "updates" && m.updates.installing:
		return "esc back • please wait for the update to finish before quitting"
	case m.screen == sectionScreen && m.section().id == "updates" && m.updatesButton() != "":
		return "enter " + m.updatesButton() + " • esc back • q quit"
	case m.screen == sectionScreen && len(actions.InSection(m.section().id)) > 0:
		return "↑/↓ move • enter run • esc back • q quit"
	case m.screen == sectionScreen:
		return "esc back • q quit"
	default:
		return "↑/↓ move • enter select • q quit"
	}
}

// systemPanelBody shows the system details as label and value rows, fitted
// to the given width. Long values wrap inside their own column so the labels
// stay lined up.
func (m model) systemPanelBody(width int) string {
	if !m.infoLoaded {
		return mutedStyle.Render("Loading system info...")
	}

	var rows []row
	for _, field := range m.info.Fields() {
		value := field.Value
		if value == "" {
			value = mutedStyle.Render("Unknown")
		}
		rows = append(rows, row{field.Label + ":", value})
	}
	lines := []string{labeledRows(rows, labelWidth, width)}
	for _, problem := range m.info.Problems {
		lines = append(lines, "", mutedStyle.Width(width).Render(problem))
	}
	return strings.Join(lines, "\n")
}

// row is one label and its value, such as "CPU:" and the processor name.
type row struct{ label, value string }

// labeledRows draws label and value rows fitted to width. Long values wrap
// inside their own column so the labels stay lined up.
func labeledRows(rows []row, labelWidth, width int) string {
	labelStyle := lipgloss.NewStyle().Width(labelWidth)
	valueStyle := lipgloss.NewStyle().Width(max(width-labelWidth, 10))
	lines := make([]string, 0, len(rows))
	for _, r := range rows {
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render(r.label), valueStyle.Render(r.value)))
	}
	return strings.Join(lines, "\n")
}
