package tui

// menuItem is one entry in the main menu.
type menuItem struct {
	title       string
	description string
	quits       bool // selecting this item closes FMOS
}

// menuItems lists the main menu in display order. No section has actions
// yet, so opening one shows a "Coming soon" screen.
var menuItems = []menuItem{
	{
		title:       "OS Repair",
		description: "Check and repair system files and disks.",
	},
	{
		title:       "Tools",
		description: "Install apps, check drivers, get app updates, and scan for viruses.",
	},
	{
		title:       "Network",
		description: "Check your connection and see what is using your network.",
	},
	{
		title:       "Logs",
		description: "See what FMOS did and how it went.",
	},
	{
		title:       "Quit",
		description: "Close FMOS.",
		quits:       true,
	},
}
