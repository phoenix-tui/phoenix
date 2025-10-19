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

func (n *NullTerminal) SetCursorPosition(x, y int) error {
	return nil
}

func (n *NullTerminal) GetCursorPosition() (x, y int, err error) {
	return 0, 0, nil
}

func (n *NullTerminal) MoveCursorUp(count int) error {
	return nil
}

func (n *NullTerminal) MoveCursorDown(count int) error {
	return nil
}

func (n *NullTerminal) MoveCursorLeft(count int) error {
	return nil
}

func (n *NullTerminal) MoveCursorRight(count int) error {
	return nil
}

func (n *NullTerminal) SaveCursorPosition() error {
	return nil
}

func (n *NullTerminal) RestoreCursorPosition() error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Visibility & Style                                   │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) HideCursor() error {
	return nil
}

func (n *NullTerminal) ShowCursor() error {
	return nil
}

func (n *NullTerminal) SetCursorStyle(style api.CursorStyle) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Operations                                           │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) Clear() error {
	return nil
}

func (n *NullTerminal) ClearLine() error {
	return nil
}

func (n *NullTerminal) ClearFromCursor() error {
	return nil
}

func (n *NullTerminal) ClearLines(count int) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Output                                                      │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) Write(s string) error {
	return nil
}

func (n *NullTerminal) WriteAt(x, y int, s string) error {
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Buffer (Windows Console API only)                    │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) ReadScreenBuffer() ([][]rune, error) {
	return nil, nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Info                                               │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) Size() (width, height int, err error) {
	return 80, 24, nil // Default terminal size
}

func (n *NullTerminal) ColorDepth() int {
	return 256 // Assume 256 colors
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Capabilities Discovery                                      │
// └─────────────────────────────────────────────────────────────┘

func (n *NullTerminal) SupportsDirectPositioning() bool {
	return false // Conservative default
}

func (n *NullTerminal) SupportsReadback() bool {
	return false // Conservative default
}

func (n *NullTerminal) SupportsTrueColor() bool {
	return true // Optimistic default
}

func (n *NullTerminal) Platform() api.Platform {
	return api.PlatformUnknown
}
