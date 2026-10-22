// Package tui draws the FMOS terminal menu.
package tui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/GioMonci/FMOS/internal/sysinfo"
)

// Run opens the FMOS menu and returns when the user quits.
func Run(version string) error {
	_, err := tea.NewProgram(newModel(version)).Run()
	return err
}

// model holds everything the menu needs to draw itself.
type model struct {
	version    string
	width      int          // terminal width in columns
	cursor     int          // index of the selected menu item
	open       bool         // true while the selected section is showing
	info       sysinfo.Info // details shown in the System panel
	infoLoaded bool         // false until the system details arrive
}

func newModel(version string) model {
	// 80 columns is a safe guess until the terminal reports its real size.
	return model{version: version, width: 80}
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
		m.width = msg.Width
	case systemInfoMsg:
		m.info = sysinfo.Info(msg)
		m.infoLoaded = true
	case tea.KeyPressMsg:
		return m.handleKey(msg.String())
	}
	return m, nil
}

// handleKey reacts to a single key press.
func (m model) handleKey(key string) (tea.Model, tea.Cmd) {
	if key == "q" || key == "ctrl+c" {
		return m, tea.Quit
	}

	if m.open {
		if key == "esc" {
			m.open = false
		}
		return m, nil
	}

	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(menuItems)-1 {
			m.cursor++
		}
	case "enter":
		if menuItems[m.cursor].quits {
			return m, tea.Quit
		}
		m.open = true
	}
	return m, nil
}
