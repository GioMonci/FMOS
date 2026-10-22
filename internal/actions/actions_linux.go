package actions

func all() []Action {
	return []Action{
		{
			ID:      "network-details",
			Section: "network",
			Name:    "Show network details",
			Description: "Lists your network interfaces with their addresses. " +
				"It only reads information and changes nothing.",
			Program: "/usr/sbin/ip",
			Args:    []string{"address"},
		},
	}
}
