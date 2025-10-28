// Package ansi provides ANSI escape sequences for mouse tracking.
// This is an infrastructure concern that generates terminal control codes.
package ansi

// MouseMode represents a mouse tracking mode.
type MouseMode int

const (
	// ModeX10 enables X10 mouse tracking (1000).
	// Reports button press only, limited to 223x223 terminal.
	ModeX10 MouseMode = 1000

	// ModeVT200 enables VT200 mouse tracking (1002).
	// Reports button press and release, drag events.
	ModeVT200 MouseMode = 1002

	// ModeButtonMotion enables button-event mouse tracking (1002).
	// Reports motion events when button is pressed.
	ModeButtonMotion MouseMode = 1002

	// ModeAnyMotion enables any-event mouse tracking (1003).
	// Reports all motion events (even without button press).
	ModeAnyMotion MouseMode = 1003

	// ModeSGR enables SGR extended mouse tracking (1006).
	// Modern protocol supporting large terminals, press/release distinction.
	ModeSGR MouseMode = 1006

	// ModeURxvt enables URxvt mouse tracking (1015).
	// Alternative extended protocol.
	ModeURxvt MouseMode = 1015

	// ModeSGRPixels enables SGR pixel mouse tracking (1016).
	// Reports pixel coordinates instead of cell coordinates.
	ModeSGRPixels MouseMode = 1016
)

// EnableMouseTracking returns the ANSI sequence to enable mouse tracking.
func EnableMouseTracking(mode MouseMode) string {
	return "\x1b[?" + itoa(int(mode)) + "h"
}

// DisableMouseTracking returns the ANSI sequence to disable mouse tracking.
func DisableMouseTracking(mode MouseMode) string {
	return "\x1b[?" + itoa(int(mode)) + "l"
}

// EnableFocusEvents returns the ANSI sequence to enable focus events (1004).
func EnableFocusEvents() string {
	return "\x1b[?1004h"
}

// DisableFocusEvents returns the ANSI sequence to disable focus events (1004).
func DisableFocusEvents() string {
	return "\x1b[?1004l"
}

// EnableMouseAll enables comprehensive mouse tracking:
// - SGR extended mode (1006) for modern protocol
// - Button-event tracking (1002) for drag support.
func EnableMouseAll() string {
	return EnableMouseTracking(ModeSGR) + EnableMouseTracking(ModeButtonMotion)
}

// DisableMouseAll disables comprehensive mouse tracking.
func DisableMouseAll() string {
	return DisableMouseTracking(ModeSGR) + DisableMouseTracking(ModeButtonMotion)
}

// itoa converts int to string without importing strconv (performance optimization).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var buf [10]byte
	i := len(buf) - 1

	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}

	return string(buf[i+1:])
}
