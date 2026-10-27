package actions

import (
	"path/filepath"

	"github.com/GioMonci/FMOS/internal/platform"
)

func all() []Action {
	return []Action{
		{
			ID:      "sfc-verifyonly",
			Section: "os-repair",
			Name:    "Check system files (SFC)",
			Description: "Checks the Windows system files your computer is running for " +
				"damage. Takes several minutes. It does not repair or change anything.",
			Program:    system32("sfc.exe"),
			Args:       []string{"/verifyonly"},
			NeedsAdmin: true,
			Findings:   sfcFindings(),
		},
		{
			ID:      "dism-checkhealth",
			Section: "os-repair",
			Name:    "Quick health check (DISM)",
			Description: "Asks Windows whether it has already found damage in its own " +
				"repair files (the component store, which Windows uses to fix and update " +
				"itself). Takes a few seconds. It does not repair or change anything.",
			Program:    system32("Dism.exe"),
			Args:       dismArgs("/CheckHealth"),
			NeedsAdmin: true,
			Findings: dismFindings("Windows has not recorded any damage. This quick check " +
				"only looks for damage Windows already knows about. For a thorough check, " +
				"run the full health scan."),
		},
		{
			ID:      "dism-scanhealth",
			Section: "os-repair",
			Name:    "Full health scan (DISM)",
			Description: "Scans all of Windows' own repair files (the component store, " +
				"which Windows uses to fix and update itself) for damage. More thorough " +
				"than the quick check, but takes several minutes. It does not repair or " +
				"change anything.",
			Program:    system32("Dism.exe"),
			Args:       dismArgs("/ScanHealth"),
			NeedsAdmin: true,
			Findings:   dismFindings("DISM scanned Windows' repair files and found no damage."),
		},
		{
			ID:      "network-details",
			Section: "network",
			Name:    "Show network details",
			Description: "Lists your network adapters with their IP addresses, " +
				"gateway, and DNS servers. It only reads information and changes nothing.",
			Program: system32("ipconfig.exe"),
			Args:    []string{"/all"},
		},
	}
}

// sfcFindings lists what "sfc /verifyonly" can report. SFC has no option for
// English output, so on other Windows languages none of these match, and
// FMOS says it could not tell what was found.
func sfcFindings() []Finding {
	return []Finding{
		{
			Phrase:  "did not find any integrity violations",
			Meaning: "SFC checked Windows' system files and found no damage.",
		},
		{
			Phrase: "found integrity violations",
			Meaning: "SFC found damaged Windows system files. FMOS did not change " +
				"anything. Run the DISM health checks next, to see whether Windows can " +
				"repair them, and send this session's log to the person helping you. " +
				"Repairs are coming to FMOS soon.",
		},
		{
			Phrase: "system repair pending which requires reboot",
			Meaning: "Windows is waiting to finish an earlier repair. Restart your " +
				"computer, then run this check again.",
		},
		{
			Phrase: "could not perform the requested operation",
			Meaning: "SFC could not run the check. Restart your computer and try " +
				"again. If it keeps happening, run the DISM health checks and send " +
				"this session's log to the person helping you.",
		},
		{
			Phrase: "could not start the repair service",
			Meaning: "SFC could not start the Windows service it needs. Restart your " +
				"computer and try again.",
		},
	}
}

// dismArgs returns the arguments for a DISM component store check, such as
// "/CheckHealth". "/English" keeps DISM's output in English on every Windows
// language, so FMOS can read what it found.
func dismArgs(check string) []string {
	return []string{"/Online", "/English", "/Cleanup-Image", check}
}

// dismFindings lists what DISM's health checks can report. Only the healthy
// result differs between the quick check and the full scan. FMOS never
// starts a repair from these results.
func dismFindings(healthy string) []Finding {
	return []Finding{
		{
			Phrase:  "No component store corruption detected",
			Meaning: healthy,
		},
		{
			Phrase: "The component store is repairable",
			Meaning: "DISM found damage that Windows can repair. FMOS did not change " +
				"anything. A repair (DISM RestoreHealth) is coming to FMOS soon. Until " +
				"then, send this session's log to the person helping you.",
		},
		{
			Phrase: "The component store cannot be repaired",
			Meaning: "DISM found damage that Windows cannot repair by itself. FMOS did " +
				"not change anything. Send this session's log to the person helping you.",
		},
	}
}

// system32 returns the full path of a built-in Windows tool. If the Windows
// folder is unknown, it returns only the name, which the runner reports as
// "not found" instead of searching PATH.
func system32(name string) string {
	windowsFolder, err := platform.WindowsFolder()
	if err != nil {
		return name
	}
	return filepath.Join(windowsFolder, "System32", name)
}
