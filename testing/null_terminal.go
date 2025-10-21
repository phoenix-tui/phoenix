package testing

import (
	"github.com/phoenix-tui/phoenix/terminal/api"
)

// NullTerminal is a no-op implementation of the Terminal interface.
//
// All methods do nothing and return success. Perfect for tests that don't
// care about terminal operations but need to avoid nil pointer panics.
//
// Example:
//
//	m := &MyModel{
//	    terminal: testing.NewNullTerminal(),
//	}
//	m.Render() // All terminal calls succeed silently
type NullTerminal struct{}

// NewNullTerminal creates a new no-op terminal.
func NewNullTerminal() api.Terminal {
	return &NullTerminal{}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Operations                                           │
// └─────────────────────────────────────────────────────────────┘

// SetCursorPosition does nothing (null implementation).
func (n *NullTerminal) SetCursorPosition(_, _ int) error {
	return nil
}

// GetCursorPosition returns zero position (null implementation).
func (n *NullTerminal) GetCursorPosition() (x, y int, err error) {
	return 0, 0, nil
}

// MoveCursorUp does nothing (null implementation).
func (n *NullTerminal) MoveCursorUp(_ int) error {
	return nil
}

// MoveCursorDown does nothing (null implementation).
func (n *NullTerminal) MoveCursorDown(_ int) error {
	return nil
}

// MoveCursorLeft does nothing (null implementation).
func (n *NullTerminal) MoveCursorLeft(_ int) error {
	return nil
}

// MoveCursorRight does nothing (null implementation).
func (n *NullTerminal) MoveCursorRight(_ int) error {
	return nil
}

// SaveCursorPosition does nothing (null implementation).
func (n *NullTerminal) SaveCursorPosition() error {
	return nil
}

// RestoreCursorPosition does nothing (null implementation).
func (n *NullTerminal) RestoreCursorPosition() error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Visibility & Style                                   │
// └─────────────────────────────────────────────────────────────┘

// HideCursor does nothing (null implementation).
func (n *NullTerminal) HideCursor() error {
	return nil
}

// ShowCursor does nothing (null implementation).
func (n *NullTerminal) ShowCursor() error {
	return nil
}

// SetCursorStyle does nothing (null implementation).
func (n *NullTerminal) SetCursorStyle(_ api.CursorStyle) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Operations                                           │
// └─────────────────────────────────────────────────────────────┘

// Clear does nothing (null implementation).
func (n *NullTerminal) Clear() error {
	return nil
}

// ClearLine does nothing (null implementation).
func (n *NullTerminal) ClearLine() error {
	return nil
}

// ClearFromCursor does nothing (null implementation).
func (n *NullTerminal) ClearFromCursor() error {
	return nil
}

// ClearLines does nothing (null implementation).
func (n *NullTerminal) ClearLines(_ int) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Output                                                      │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) Write(_ string) error {
	return nil
}

// WriteAt does nothing (null implementation).
func (n *NullTerminal) WriteAt(_, _ int, _ string) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Buffer (Windows Console API only)                    │
// └─────────────────────────────────────────────────────────────┘

// ReadScreenBuffer returns nil (null implementation).
func (n *NullTerminal) ReadScreenBuffer() ([][]rune, error) {
	return nil, nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Info                                               │
// └─────────────────────────────────────────────────────────────┘

// Size returns zero size (null implementation).
func (n *NullTerminal) Size() (width, height int, err error) {
	return 80, 24, nil // Default terminal size
}

// ColorDepth returns zero (null implementation).
func (n *NullTerminal) ColorDepth() int {
	return 256 // Assume 256 colors
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Capabilities Discovery                                      │
// └─────────────────────────────────────────────────────────────┘

// SupportsDirectPositioning returns false (null implementation).
func (n *NullTerminal) SupportsDirectPositioning() bool {
	return false // Conservative default
}

// SupportsReadback returns false (null implementation).
func (n *NullTerminal) SupportsReadback() bool {
	return false // Conservative default
}

// SupportsTrueColor returns false (null implementation).
func (n *NullTerminal) SupportsTrueColor() bool {
	return true // Optimistic default
}

// Platform returns null platform (null implementation).
func (n *NullTerminal) Platform() api.Platform {
	return api.PlatformUnknown
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Alternate Screen Buffer                                     │
// └─────────────────────────────────────────────────────────────┘

// EnterAltScreen does nothing (null implementation).
func (n *NullTerminal) EnterAltScreen() error {
	return nil
}

// ExitAltScreen does nothing (null implementation).
func (n *NullTerminal) ExitAltScreen() error {
	return nil
}

// IsInAltScreen returns false (null implementation).
func (n *NullTerminal) IsInAltScreen() bool {
	return false
}
