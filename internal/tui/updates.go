package tui

import (
	"context"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/GioMonci/FMOS/internal/update"
)

// updateState tracks checking for, and installing, a newer FMOS.
type updateState struct {
	checking   bool
	checkErr   error          // why the last check failed
	release    update.Release // the latest release, once checked
	available  bool           // the latest release is newer than this one
	installing bool
	installErr error // why the last update failed
	installed  bool  // the new version is in place and starts next time
}

// Messages from the background update work.
type (
	updateCheckedMsg struct {
		release update.Release
		newer   bool
		err     error
	}
	updateInstalledMsg struct{ err error }
	cleanUpMsg         struct{ err error }
)

// checkForUpdates asks GitHub for the latest release in the background.
func checkForUpdates(version string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		release, newer, err := update.Check(ctx, version)
		return updateCheckedMsg{release: release, newer: newer, err: err}
	}
}

// installUpdate downloads, verifies, and swaps in a release in the background.
func installUpdate(release update.Release) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		return updateInstalledMsg{err: update.Install(ctx, release)}
	}
}

// cleanUpAfterUpdate removes files left over from an earlier update.
func cleanUpAfterUpdate() tea.Msg {
	return cleanUpMsg{err: update.CleanUp()}
}

// updateUpdates handles the results of the background update work and logs them.
func (m model) updateUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateCheckedMsg:
		m.updates.checking = false
		m.updates.checkErr = msg.err
		m.updates.release = msg.release
		m.updates.available = msg.newer
		switch {
		case msg.err != nil:
			m.writeNote("Update check failed: " + msg.err.Error())
		case msg.newer:
			m.writeNote("Update check: FMOS " + msg.release.Version + " is available.")
		default:
			m.writeNote("Update check: FMOS " + m.version + " is the newest version.")
		}
	case updateInstalledMsg:
		m.updates.installing = false
		m.updates.installErr = msg.err
		m.updates.installed = msg.err == nil
		if msg.err != nil {
			m.writeNote("Update to " + m.updates.release.Version + " failed, so the current version " +
				"was kept: " + msg.err.Error())
		} else {
			m.writeNote("Updated FMOS from " + m.version + " to " + m.updates.release.Version +
				" (checksum verified). It starts the next time FMOS opens.")
		}
	case cleanUpMsg:
		if msg.err != nil {
			m.writeNote("Could not remove files left over from an earlier update: " + msg.err.Error())
		}
	}
	return m, nil
}

// handleUpdatesEnter runs the Updates section's button: check again, update,
// or try the update again.
func (m model) handleUpdatesEnter() (tea.Model, tea.Cmd) {
	u := &m.updates
	switch {
	case u.checking || u.installing || u.installed:
		return m, nil
	case u.available && u.release.CanInstall():
		u.installing = true
		u.installErr = nil
		return m, installUpdate(u.release)
	default:
		u.checking = true
		return m, checkForUpdates(m.version)
	}
}

// updatesBody shows the Updates section.
func (m model) updatesBody() string {
	u := m.updates
	body := "You have FMOS " + m.version + ".\n\n"
	switch {
	case u.installing:
		return body + "Downloading and checking FMOS " + u.release.Version + "...\n" +
			"Please wait. Do not close this window."
	case u.installed:
		return body + "FMOS was updated to " + u.release.Version + ". Quit FMOS and open it " +
			"again to start the new version."
	case u.checking:
		return body + mutedStyle.Render("Checking for a newer version...")
	case u.checkErr != nil:
		return body + "Could not check for updates. You may be offline. FMOS keeps working " +
			"normally.\n\n" + mutedStyle.Render(u.checkErr.Error()) + "\n\n" +
			selectedStyle.Render("> Check again")
	case !u.available:
		return body + "This is the newest version.\n\n" + selectedStyle.Render("> Check again")
	}

	body += "FMOS " + u.release.Version + " is available.\n\n"
	if u.release.Notes != "" {
		body += "What's new:\n" + firstLines(u.release.Notes, 8) + "\n\n"
	}
	if !u.release.CanInstall() {
		return body + mutedStyle.Render("There is no download for this computer yet.") + "\n\n" +
			selectedStyle.Render("> Check again")
	}
	body += mutedStyle.Render("Updating downloads the new version from GitHub, checks that "+
		"the file is intact, and replaces this copy. If anything goes wrong, this version "+
		"keeps working.") + "\n\n"
	if u.installErr != nil {
		body += "The last try did not finish: " + u.installErr.Error() + ".\n\n" +
			selectedStyle.Render("> Try again")
		return body
	}
	return body + selectedStyle.Render("> Update to "+u.release.Version)
}

// updatesButton names what Enter does in the Updates section, for the key hints.
func (m model) updatesButton() string {
	u := m.updates
	switch {
	case u.checking || u.installing || u.installed:
		return ""
	case u.available && u.release.CanInstall():
		return "update"
	default:
		return "check again"
	}
}

// updateBadge is a short note for the header about the update status.
func (m model) updateBadge() string {
	u := m.updates
	switch {
	case u.installing:
		return "Updating to " + u.release.Version + ", please wait..."
	case u.installed:
		return "Updated to " + u.release.Version + ": restart FMOS to use it"
	case u.available && u.release.CanInstall():
		return "Update available: " + u.release.Version + " (see Updates)"
	default:
		return ""
	}
}

// writeNote adds a line to the session log. If writing fails, FMOS shows its
// "cannot save logs" warning.
func (m *model) writeNote(text string) {
	if m.log == nil {
		return
	}
	if err := m.log.Note(text); err != nil {
		m.logErr = err
	}
}

// firstLines returns up to count lines of text, noting when more were cut.
func firstLines(text string, count int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= count {
		return text
	}
	return strings.Join(lines[:count], "\n") + "\n" + mutedStyle.Render("(more in the release notes on GitHub)")
}
