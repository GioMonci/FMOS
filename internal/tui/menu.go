package tui

// menuItem is one entry in the main menu.
type menuItem struct {
	id          string // stable name used in code, such as "logs"
	title       string
	description string
}

// menuItems lists the main menu in display order. Sections without actions
// yet show a "Coming soon" screen.
var menuItems = []menuItem{
	{
		id:          "os-repair",
		title:       "OS Repair",
		description: "Check and repair system files and disks.",
	},
	{
		id:          "tools",
		title:       "Tools",
		description: "Install apps, check drivers, get app updates, and scan for viruses.",
	},
	{
		id:          "network",
		title:       "Network",
		description: "Check your connection and see what is using your network.",
	},
	{
		id:          "logs",
		title:       "Logs",
		description: "See what FMOS did and how it went.",
	},
	{
		id:          "updates",
		title:       "Updates",
		description: "Check for a newer version of FMOS and update it.",
	},
	{
		id:          "quit",
		title:       "Quit",
		description: "Close FMOS.",
	},
}
