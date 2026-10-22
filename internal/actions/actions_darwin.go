package actions

func all() []Action {
	return []Action{
		{
			Section: "network",
			Name:    "Show network details",
			Description: "Lists your network interfaces with their addresses. " +
				"It only reads information and changes nothing.",
			Program: "/sbin/ifconfig",
		},
	}
}
