// Package actions lists what FMOS can run from its menu, with everything the
// user should know before running it. Each operating system's list is in its
// own actions_*.go file.
package actions

import "github.com/GioMonci/FMOS/internal/runner"

// Action is one thing FMOS can run from a menu section.
type Action struct {
	ID             string   // stable name, such as "network-details" (used by --open)
	Section        string   // menu section ID, such as "network"
	Name           string   // shown in the menu
	Description    string   // plain-English explanation of what it does
	Program        string   // full path to the program
	Args           []string // arguments, passed exactly as written
	NeedsAdmin     bool     // needs Administrator (Windows) or root
	ChangesSystem  bool     // false for read-only actions
	MayNeedRestart bool
	NeedsInternet  bool
}

// Command returns the program and arguments for the runner.
func (a Action) Command() runner.Command {
	return runner.Command{Program: a.Program, Args: a.Args}
}

// NeedsReview reports whether FMOS must show the review screen and wait for
// the user to confirm before running the action. Read-only actions that need
// nothing special run right away.
func (a Action) NeedsReview(isAdmin bool) bool {
	return a.ChangesSystem || (a.NeedsAdmin && !isAdmin)
}

// Find returns the action with the given ID.
func Find(id string) (Action, bool) {
	for _, action := range all() {
		if action.ID == id {
			return action, true
		}
	}
	return Action{}, false
}

// InSection returns the actions in one menu section, in display order.
func InSection(section string) []Action {
	var found []Action
	for _, action := range all() {
		if action.Section == section {
			found = append(found, action)
		}
	}
	return found
}
