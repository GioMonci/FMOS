package tui

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/GioMonci/FMOS/internal/actions"
	"github.com/GioMonci/FMOS/internal/platform"
)

// reviewLabelWidth fits the longest review label ("Needs Administrator:").
const reviewLabelWidth = 22

// chooseAction runs a read-only action that needs nothing special right
// away. Anything else goes to the review screen first.
func (m model) chooseAction(action actions.Action) (tea.Model, tea.Cmd) {
	if action.NeedsReview(platform.IsAdmin()) {
		return m.openReview(action), nil
	}
	return m.startAction(action)
}

// openReview shows the review screen for an action, with the menu behind it
// pointing at the action's section.
func (m model) openReview(action actions.Action) model {
	for i, item := range menuItems {
		if item.id == action.Section {
			m.cursor = i
		}
	}
	for i, sectionAction := range actions.InSection(action.Section) {
		if sectionAction.ID == action.ID {
			m.actionCursor = i
		}
	}
	m.review = action
	m.screen = reviewScreen
	m.notice = ""
	return m
}

// handleReviewKey handles keys on the review screen. Only "y" goes ahead, so
// pressing Enter twice by accident never starts anything.
func (m model) handleReviewKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = sectionScreen
		m.notice = ""
	case "y":
		if m.needsAdminRestart() {
			return m.restartAsAdmin()
		}
		return m.startAction(m.review)
	}
	return m, nil
}

// needsAdminRestart reports whether the reviewed action needs Administrator
// rights that this copy of FMOS does not have.
func (m model) needsAdminRestart() bool {
	return m.review.NeedsAdmin && !platform.IsAdmin()
}

// restartAsAdmin opens a new copy of FMOS as Administrator on this review
// screen, then closes this one. The new copy never runs the action on its
// own; the user confirms again there.
func (m model) restartAsAdmin() (tea.Model, tea.Cmd) {
	err := platform.RestartAsAdmin([]string{"--open", m.review.ID})
	switch {
	case err == nil:
		m.writeNote("Restarted FMOS as Administrator to run " + m.review.Name +
			". The new window keeps its own log.")
		m.restartedAsAdmin = true
		return m, tea.Quit
	case errors.Is(err, platform.ErrElevationDeclined):
		m.writeNote("The Administrator request was declined, so " + m.review.Name + " was not run.")
		m.notice = "Windows did not give FMOS Administrator permissions, so " + m.review.Name +
			" was not run. Press Y to ask again, or Esc to go back."
	default:
		m.notice = "Could not restart FMOS as Administrator: " + err.Error() + "."
	}
	return m, nil
}

// renderReview draws the full-width review screen: what the action does,
// what it will change, and what the user needs to press to go ahead.
func (m model) renderReview(width int) string {
	action := m.review
	contentWidth := width - panelStyle.GetHorizontalFrameSize()

	admin := "No."
	switch {
	case action.NeedsAdmin && platform.IsAdmin():
		admin = "Yes. FMOS is running as Administrator."
	case action.NeedsAdmin:
		admin = "Yes. FMOS is not running as Administrator yet."
	}
	changes := "No. It only reads information."
	if action.ChangesSystem {
		changes = "Yes."
	}
	restart := "No."
	if action.MayNeedRestart {
		restart = "Possibly. FMOS never restarts your computer by itself."
	}
	internet := "No."
	if action.NeedsInternet {
		internet = "Yes."
	}

	facts := labeledRows([]row{
		{"Runs:", shortCommand(action.Command())},
		{"Changes your system:", changes},
		{"Needs Administrator:", admin},
		{"May need a restart:", restart},
		{"Needs internet:", internet},
	}, reviewLabelWidth, contentWidth)

	prompt := "Press Y to run it, or Esc to go back."
	if m.needsAdminRestart() {
		prompt = "Press Y to restart FMOS as Administrator. Windows will ask you to " +
			"confirm, then FMOS opens again on this screen. Nothing runs until you " +
			"press Y there."
	}

	body := action.Description + "\n\n" + facts + "\n\n" +
		lipgloss.NewStyle().Width(contentWidth).Render(selectedStyle.Render("> ")+prompt)
	if m.notice != "" {
		body += "\n\n" + lipgloss.NewStyle().Width(contentWidth).Render(m.notice)
	}
	return panel(action.Name, body, width, 0)
}

func (m model) reviewKeyHints() string {
	if m.needsAdminRestart() {
		return "y restart as Administrator • esc back • q quit"
	}
	return "y run • esc back • q quit"
}
