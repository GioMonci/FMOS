// Package actions lists what FMOS can run from its menu, with everything the
// user should know before running it. Each operating system's list is in its
// own actions_*.go file.
package actions

import (
	"strings"

	"github.com/GioMonci/FMOS/internal/runner"
)

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
	Findings       []Finding // what the program's output can mean, checked in order
}

// Finding is a sentence a program prints to say what it found, and what that
// means in plain English. Some tools, such as DISM, report "The operation
// completed successfully" whether or not they found damage, so the exit code
// alone does not tell the user what happened.
type Finding struct {
	Phrase  string // printed by the program, such as "The component store is repairable"
	Meaning string // shown to the user and written to the log
}

// Command returns the program and arguments for the runner.
func (a Action) Command() runner.Command {
	return runner.Command{Program: a.Program, Args: a.Args}
}

// Meaning explains in plain English what a run found, using the action's
// Findings. It returns "" when the action has no Findings, or when the
// program did not run to the end.
func (a Action) Meaning(result runner.Result) string {
	if len(a.Findings) == 0 {
		return ""
	}
	if result.Outcome != runner.Finished && result.Outcome != runner.Failed {
		return ""
	}
	for _, finding := range a.Findings {
		if outputContains(result.Output, finding.Phrase) {
			return finding.Meaning
		}
	}
	return "FMOS could not tell from the output what was found. Read the output below."
}

// outputContains reports whether any line contains phrase, ignoring case.
func outputContains(output []string, phrase string) bool {
	phrase = strings.ToLower(phrase)
	for _, line := range output {
		if strings.Contains(strings.ToLower(line), phrase) {
			return true
		}
	}
	return false
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
