package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	labelStyle    = lipgloss.NewStyle().Width(labelWidth)
)

// labelWidth fits the longest System panel label ("Uptime:") plus a space.
const labelWidth = 8

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.WindowTitle = "FMOS"
	return v
}

// render draws the whole screen: header, panels, and key hints.
func (m model) render() string {
	width := min(m.width, maxWidth)
	header := titleStyle.Render("FMOS (Fix My OS)") + mutedStyle.Render("  version "+m.version)
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.renderPanels(width),
		mutedStyle.Render(m.keyHints()),
	)
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

// rightPanelContent returns the menu, or the screen of the open section.
func (m model) rightPanelContent() (title, body string) {
	selected := menuItems[m.cursor]
	if m.open {
		return selected.title, selected.description + "\n\n" +
			mutedStyle.Render("Coming soon. This part of FMOS has not been built yet.")
	}

	lines := make([]string, 0, len(menuItems))
	for i, item := range menuItems {
		if i == m.cursor {
			lines = append(lines, selectedStyle.Render("> "+item.title))
		} else {
			lines = append(lines, "  "+item.title)
		}
	}
	return "Menu", strings.Join(lines, "\n") + "\n\n" + mutedStyle.Render(selected.description)
}

func (m model) keyHints() string {
	if m.open {
		return "esc back • q quit"
	}
	return "↑/↓ move • enter select • q quit"
}

// systemPanelBody shows the system details as label and value rows, fitted
// to the given width. Long values wrap inside their own column so the labels
// stay lined up.
func (m model) systemPanelBody(width int) string {
	if !m.infoLoaded {
		return mutedStyle.Render("Loading system info...")
	}

	info := m.info
	admin := "No"
	if info.Admin {
		admin = "Yes"
	}
	rows := []struct{ label, value string }{
		{"OS", info.OS},
		{"CPU", info.CPU},
		{"GPU", info.GPU},
		{"RAM", info.RAM},
		{"Disk", info.Disk},
		{"Uptime", info.Uptime},
		{"User", info.User},
		{"Admin", admin},
	}

	valueStyle := lipgloss.NewStyle().Width(max(width-labelWidth, 10))
	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		value := row.value
		if value == "" {
			value = mutedStyle.Render("Unknown")
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render(row.label+":"),
			valueStyle.Render(value),
		))
	}
	for _, problem := range info.Problems {
		lines = append(lines, "", mutedStyle.Width(width).Render(problem))
	}
	return strings.Join(lines, "\n")
}
