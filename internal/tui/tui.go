// Package tui draws the FMOS terminal menu.
package tui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/GioMonci/FMOS/internal/actions"
	"github.com/GioMonci/FMOS/internal/logging"
	"github.com/GioMonci/FMOS/internal/platform"
	"github.com/GioMonci/FMOS/internal/sysinfo"
)

// Run opens the FMOS menu and returns when the user quits. session is this
// run's log, or nil when logs cannot be saved; logErr then explains why.
func Run(version string, session *logging.Session, logErr error) error {
	_, err := tea.NewProgram(newModel(version, session, logErr)).Run()
	return err
}

// screen is the part of FMOS that is showing.
type screen int

const (
	menuScreen    screen = iota // the main menu
	sectionScreen               // one section: its actions, or a "Coming soon" note
	taskScreen                  // an action running, or its result
)

// model holds everything FMOS needs to draw itself.
type model struct {
	version      string
	width        int // terminal width in columns
	height       int // terminal height in rows
	screen       screen
	cursor       int              // selected menu item
	actionCursor int              // selected action in the open section
	notice       string           // short message shown in the open section
	info         sysinfo.Info     // details shown in the System panel
	infoLoaded   bool             // false until the system details arrive
	log          *logging.Session // nil when logs cannot be saved
	logErr       error            // why logs cannot be saved, if they cannot
	task         *task            // the action running or last run, if any
}

func newModel(version string, session *logging.Session, logErr error) model {
	// 80x24 is a safe guess until the terminal reports its real size.
	return model{version: version, width: 80, height: 24, log: session, logErr: logErr}
}

// systemInfoMsg delivers the system details once they have been gathered.
type systemInfoMsg sysinfo.Info

// loadSystemInfo gathers system details in the background, so the menu
// opens right away even when Windows is slow to answer.
func loadSystemInfo() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return systemInfoMsg(sysinfo.Gather(ctx))
}

func (m model) Init() tea.Cmd {
	return loadSystemInfo
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case systemInfoMsg:
		m.info = sysinfo.Info(msg)
		m.infoLoaded = true
		if m.log != nil {
			if err := m.log.LogSystemInfo(m.info); err != nil {
				m.logErr = err
			}
		}
	case outputLineMsg, taskDoneMsg, tickMsg:
		return m.updateTask(msg)
	case tea.KeyPressMsg:
		return m.handleKey(msg.String())
	}
	return m, nil
}

// handleKey passes a key press to the screen that is showing.
func (m model) handleKey(key string) (tea.Model, tea.Cmd) {
	switch m.screen {
	case taskScreen:
		return m.handleTaskKey(key)
	case sectionScreen:
		return m.handleSectionKey(key)
	default:
		return m.handleMenuKey(key)
	}
}

func (m model) handleMenuKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		m.cursor = max(m.cursor-1, 0)
	case "down", "j":
		m.cursor = min(m.cursor+1, len(menuItems)-1)
	case "enter":
		if m.section().id == "quit" {
			return m, tea.Quit
		}
		m.screen = sectionScreen
		m.actionCursor = 0
		m.notice = ""
	}
	return m, nil
}

func (m model) handleSectionKey(key string) (tea.Model, tea.Cmd) {
	sectionActions := actions.InSection(m.section().id)
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = menuScreen
		m.notice = ""
	case "up", "k":
		m.actionCursor = max(m.actionCursor-1, 0)
	case "down", "j":
		m.actionCursor = min(m.actionCursor+1, max(len(sectionActions)-1, 0))
	case "enter":
		if m.section().id == "logs" && m.logsAvailable() {
			m.notice = m.openLogsFolder()
		} else if len(sectionActions) > 0 {
			return m.startAction(sectionActions[m.actionCursor])
		}
	}
	return m, nil
}

// section returns the selected menu item.
func (m model) section() menuItem {
	return menuItems[m.cursor]
}

// logsAvailable reports whether this session's log is being saved.
func (m model) logsAvailable() bool {
	return m.log != nil && m.logErr == nil
}

// openLogsFolder opens the logs folder in the system's file manager and
// returns a short message saying how it went.
func (m model) openLogsFolder() string {
	if err := platform.OpenFolder(m.log.Dir()); err != nil {
		return "Could not open the logs folder: " + err.Error()
	}
	return "Opened the logs folder."
}
