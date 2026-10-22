package platform

import "errors"

var (
	// ErrElevationDeclined means the user chose "No" on the Windows
	// Administrator (UAC) prompt.
	ErrElevationDeclined = errors.New("the Administrator request was declined")

	// ErrElevationNotSupported means this operating system cannot restart
	// FMOS with more rights from inside FMOS.
	ErrElevationNotSupported = errors.New("restarting FMOS with more rights is only available on Windows")
)
