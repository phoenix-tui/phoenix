// Package ansi provides ANSI escape sequence generation.
package ansi

import "fmt"

// ANSI control sequences and escape codes.
const (
	// CSI - Control Sequence Introducer.
	CSI = "\x1b["

	// ESC - Escape character.
	ESC = "\x1b"

	// Reset - Reset all attributes.
	Reset = CSI + "0m"

	// Text attributes.
	Bold      = CSI + "1m"
	Dim       = CSI + "2m"
	Italic    = CSI + "3m"
	Underline = CSI + "4m"
	Blink     = CSI + "5m"
	Reverse   = CSI + "7m"
	Hidden    = CSI + "8m"
	Strike    = CSI + "9m"

	// Text attribute resets.
	NoBold      = CSI + "22m"
	NoDim       = CSI + "22m"
	NoItalic    = CSI + "23m"
	NoUnderline = CSI + "24m"
	NoBlink     = CSI + "25m"
	NoReverse   = CSI + "27m"
	NoHidden    = CSI + "28m"
	NoStrike    = CSI + "29m"

	// Screen operations.
	ClearScreen      = CSI + "2J"
	ClearLine        = CSI + "2K"
	ClearLineFromPos = CSI + "0K"
	ClearLineToPos   = CSI + "1K"
	ClearScreenBelow = CSI + "0J"
	ClearScreenAbove = CSI + "1J"

	// Cursor operations.
	CursorHome    = CSI + "H"
	CursorHide    = CSI + "?25l"
	CursorShow    = CSI + "?25h"
	CursorSave    = CSI + "s"
	CursorRestore = CSI + "u"

	// Alternative screen buffer.
	AltScreenEnable  = CSI + "?1049h"
	AltScreenDisable = CSI + "?1049l"

	// Mouse tracking.
	MouseTrackingEnable    = CSI + "?1000h"
	MouseTrackingDisable   = CSI + "?1000l"
	MouseCellMotionEnable  = CSI + "?1002h"
	MouseCellMotionDisable = CSI + "?1002l"
	MouseAllMotionEnable   = CSI + "?1003h"
	MouseAllMotionDisable  = CSI + "?1003l"
	MouseSGREnable         = CSI + "?1006h"
	MouseSGRDisable        = CSI + "?1006l"
)

// MoveCursor returns ANSI sequence to move cursor to (x, y).
// Coordinates are 1-based (terminal convention).
func MoveCursor(x, y int) string {
	return fmt.Sprintf("%s%d;%dH", CSI, y+1, x+1)
}

// MoveCursorUp returns ANSI sequence to move cursor up n lines.
func MoveCursorUp(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%dA", CSI, n)
}

// MoveCursorDown returns ANSI sequence to move cursor down n lines.
func MoveCursorDown(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%dB", CSI, n)
}

// MoveCursorRight returns ANSI sequence to move cursor right n columns.
func MoveCursorRight(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%dC", CSI, n)
}

// MoveCursorLeft returns ANSI sequence to move cursor left n columns.
func MoveCursorLeft(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%dD", CSI, n)
}

// SetFg256 returns ANSI sequence to set 256-color foreground.
func SetFg256(color uint8) string {
	return fmt.Sprintf("%s38;5;%dm", CSI, color)
}

// SetBg256 returns ANSI sequence to set 256-color background.
func SetBg256(color uint8) string {
	return fmt.Sprintf("%s48;5;%dm", CSI, color)
}

// SetFgRGB returns ANSI sequence to set RGB foreground color.
func SetFgRGB(r, g, b uint8) string {
	return fmt.Sprintf("%s38;2;%d;%d;%dm", CSI, r, g, b)
}

// SetBgRGB returns ANSI sequence to set RGB background color.
func SetBgRGB(r, g, b uint8) string {
	return fmt.Sprintf("%s48;2;%d;%d;%dm", CSI, r, g, b)
}

// ResetFg returns ANSI sequence to reset foreground color to default.
func ResetFg() string {
	return CSI + "39m"
}

// ResetBg returns ANSI sequence to reset background color to default.
func ResetBg() string {
	return CSI + "49m"
}

// SetCursorShape sets cursor shape (0=default, 1=blinking block, 2=steady block, etc.).
func SetCursorShape(shape int) string {
	return fmt.Sprintf("%s%d q", CSI, shape)
}

// SetScrollRegion sets scrolling region from top to bottom (1-based).
func SetScrollRegion(top, bottom int) string {
	return fmt.Sprintf("%s%d;%dr", CSI, top, bottom)
}

// SaveCursorPosition returns ANSI sequence to save cursor position (DEC).
func SaveCursorPosition() string {
	return ESC + "7"
}

// RestoreCursorPosition returns ANSI sequence to restore cursor position (DEC).
func RestoreCursorPosition() string {
	return ESC + "8"
}

// SetTitle sets the terminal window title.
func SetTitle(title string) string {
	return fmt.Sprintf("%s]0;%s\x07", ESC, title)
}

// Bell returns the bell/beep control character.
func Bell() string {
	return "\x07"
}
