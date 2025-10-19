// Package model provides domain models for the Phoenix Tea event loop.
// This package defines the message types that flow through the Elm Architecture.
package model

import (
	"fmt"
	"strings"
)

const (
	unknownKeyName = "unknown"
)

// Msg represents any message that can be sent through the event loop.
// This is a marker interface - any type can be a message.
type Msg interface{}

// KeyType represents the type of key pressed.
type KeyType int

// Key type constants define all supported keyboard events.
const (
	KeyRune KeyType = iota // Regular character key
	KeyEnter
	KeyBackspace
	KeyTab
	KeyEsc
	KeySpace
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlC // Ctrl+C (common, so dedicated type)
)

// KeyMsg represents a keyboard event.
type KeyMsg struct {
	Type  KeyType // Type of key pressed
	Rune  rune    // The actual rune (for KeyRune type)
	Alt   bool    // Alt modifier
	Ctrl  bool    // Ctrl modifier
	Shift bool    // Shift modifier
}

// keyTypeName returns the string representation of a KeyType.
//
//nolint:gocyclo,cyclop,funlen // Exhaustive switch for all KeyType constants (29 cases) - complexity is intentional
func keyTypeName(kt KeyType, r rune) string {
	switch kt {
	case KeyRune:
		return string(r)
	case KeyEnter:
		return "enter"
	case KeyBackspace:
		return "backspace"
	case KeyTab:
		return "tab"
	case KeyEsc:
		return "esc"
	case KeySpace:
		return "space"
	case KeyUp:
		return "↑"
	case KeyDown:
		return "↓"
	case KeyLeft:
		return "←"
	case KeyRight:
		return "→"
	case KeyHome:
		return "home"
	case KeyEnd:
		return "end"
	case KeyPgUp:
		return "pgup"
	case KeyPgDown:
		return "pgdown"
	case KeyDelete:
		return "delete"
	case KeyInsert:
		return "insert"
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyCtrlC:
		return "ctrl+c"
	default:
		return unknownKeyName
	}
}

// String returns a human-readable representation of the key.
//
// Examples:
//   - KeyMsg{Type: KeyRune, Rune: 'a'}                    → "a"
//   - KeyMsg{Type: KeyRune, Rune: 'A', Shift: true}       → "A"
//   - KeyMsg{Type: KeyEnter}                              → "enter"
//   - KeyMsg{Type: KeyRune, Rune: 'c', Ctrl: true}        → "ctrl+c"
//   - KeyMsg{Type: KeyCtrlC}                              → "ctrl+c"
//   - KeyMsg{Type: KeyUp}                                 → "↑"
//   - KeyMsg{Type: KeyF1}                                 → "F1"
func (k KeyMsg) String() string {
	// Special case: dedicated Ctrl+C
	if k.Type == KeyCtrlC {
		return "ctrl+c"
	}

	var parts []string

	// Add modifiers
	if k.Alt {
		parts = append(parts, "alt")
	}
	if k.Ctrl {
		parts = append(parts, "ctrl")
	}
	if k.Shift && k.Type != KeyRune {
		// For runes, shift is implicit in the character itself (A vs a)
		parts = append(parts, "shift")
	}

	// Add key name
	keyName := keyTypeName(k.Type, k.Rune)
	parts = append(parts, keyName)

	return strings.Join(parts, "+")
}

// MouseButton represents which mouse button was used.
type MouseButton int

// Mouse button constants define all supported mouse buttons.
const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
)

// MouseAction represents what the mouse did.
type MouseAction int

// Mouse action constants define all supported mouse actions.
const (
	MouseActionPress   MouseAction = iota // Button pressed
	MouseActionRelease                    // Button released
	MouseActionMotion                     // Mouse moved
)

// MouseMsg represents a mouse event.
type MouseMsg struct {
	X      int         // Column (0-based)
	Y      int         // Row (0-based)
	Button MouseButton // Which button
	Action MouseAction // What happened
	Alt    bool        // Alt modifier
	Ctrl   bool        // Ctrl modifier
	Shift  bool        // Shift modifier
}

// mouseButtonName returns the string representation of a MouseButton.
func mouseButtonName(mb MouseButton) string {
	switch mb {
	case MouseButtonNone:
		return "mouse"
	case MouseButtonLeft:
		return "left"
	case MouseButtonMiddle:
		return "middle"
	case MouseButtonRight:
		return "right"
	case MouseButtonWheelUp:
		return "wheel up"
	case MouseButtonWheelDown:
		return "wheel down"
	default:
		return unknownKeyName
	}
}

// mouseActionName returns the string representation of a MouseAction.
func mouseActionName(ma MouseAction) string {
	switch ma {
	case MouseActionPress:
		return "press"
	case MouseActionRelease:
		return "release"
	case MouseActionMotion:
		return "motion"
	default:
		return unknownKeyName
	}
}

// String returns a human-readable representation of the mouse event.
//
// Examples:
//   - MouseMsg{X: 10, Y: 5, Button: MouseButtonLeft, Action: MouseActionPress}    → "left press at (10, 5)"
//   - MouseMsg{X: 20, Y: 10, Button: MouseButtonRight, Action: MouseActionRelease} → "right release at (20, 10)"
//   - MouseMsg{X: 15, Y: 8, Button: MouseButtonNone, Action: MouseActionMotion}   → "mouse motion at (15, 8)"
//   - MouseMsg{X: 0, Y: 0, Button: MouseButtonWheelUp, Action: MouseActionPress}  → "wheel up at (0, 0)"
func (m MouseMsg) String() string {
	buttonName := mouseButtonName(m.Button)
	actionName := mouseActionName(m.Action)

	// For motion, we don't need to specify button (unless it's being held)
	if m.Action == MouseActionMotion && m.Button == MouseButtonNone {
		return fmt.Sprintf("mouse motion at (%d, %d)", m.X, m.Y)
	}

	// For wheel events, the action is implicit
	if m.Button == MouseButtonWheelUp || m.Button == MouseButtonWheelDown {
		return fmt.Sprintf("%s at (%d, %d)", buttonName, m.X, m.Y)
	}

	return fmt.Sprintf("%s %s at (%d, %d)", buttonName, actionName, m.X, m.Y)
}

// WindowSizeMsg represents a terminal resize event.
type WindowSizeMsg struct {
	Width  int // Terminal width in columns
	Height int // Terminal height in rows
}

// String returns a human-readable representation.
//
// Example:
//   - WindowSizeMsg{Width: 80, Height: 24} → "window resize: 80x24"
func (w WindowSizeMsg) String() string {
	return fmt.Sprintf("window resize: %dx%d", w.Width, w.Height)
}

// IsValid checks if the window size is valid (positive dimensions).
func (w WindowSizeMsg) IsValid() bool {
	return w.Width > 0 && w.Height > 0
}

// QuitMsg signals that the program should quit.
// This is a message, not a command. The application can choose to ignore it
// or perform cleanup before actually quitting.
type QuitMsg struct{}

// String returns a human-readable representation.
func (q QuitMsg) String() string {
	return "quit"
}
