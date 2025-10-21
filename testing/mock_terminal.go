package testing

import (
	"fmt"
	"sync"

	"github.com/phoenix-tui/phoenix/terminal/api"
)

// MockTerminal is a recording implementation of the Terminal interface.
//
// All methods are no-ops but record the method name and arguments in the
// Calls slice. Perfect for tests that need to verify terminal operations.
//
// Thread-safe: Can be called from multiple goroutines.
//
// Example:
//
//	mock := testing.NewMockTerminal()
//	m := &MyModel{terminal: mock}
//
//	m.Render()
//
//	// Verify operations
//	assert.Contains(t, mock.Calls, "ClearLine")
//	assert.Equal(t, 1, mock.CallCount("ClearLine"))
//	assert.Equal(t, "SetCursorPosition(10, 5)", mock.Calls[0])
type MockTerminal struct {
	inAltScreen bool // Tracks alternate screen state
	inRawMode   bool // Tracks raw mode state
	mu          sync.Mutex
	Calls       []string // All recorded method calls with arguments
}

// NewMockTerminal creates a new mock terminal.
func NewMockTerminal() *MockTerminal {
	return &MockTerminal{
		Calls: make([]string, 0),
	}
}

// record adds a method call to the Calls slice (thread-safe).
func (m *MockTerminal) record(call string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, call)
}

// CallCount returns the number of times a method was called.
//
// Example:
//
//	count := mock.CallCount("ClearLine")
func (m *MockTerminal) CallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for _, call := range m.Calls {
		if call == method || len(call) > len(method) && call[:len(method)] == method && call[len(method)] == '(' {
			count++
		}
	}
	return count
}

// Reset clears all recorded calls.
//
// Useful when you want to reuse the same mock in multiple test phases.
func (m *MockTerminal) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = make([]string, 0)
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Operations                                           │
// └─────────────────────────────────────────────────────────────┘

// SetCursorPosition sets the cursor position (mock implementation).
func (m *MockTerminal) SetCursorPosition(x, y int) error {
	m.record(fmt.Sprintf("SetCursorPosition(%d, %d)", x, y))
	return nil
}

// GetCursorPosition returns the current cursor position (mock implementation).
func (m *MockTerminal) GetCursorPosition() (x, y int, err error) {
	m.record("GetCursorPosition")
	return 0, 0, nil
}

// MoveCursorUp moves the cursor up (mock implementation).
func (m *MockTerminal) MoveCursorUp(n int) error {
	m.record(fmt.Sprintf("MoveCursorUp(%d)", n))
	return nil
}

// MoveCursorDown moves the cursor down (mock implementation).
func (m *MockTerminal) MoveCursorDown(n int) error {
	m.record(fmt.Sprintf("MoveCursorDown(%d)", n))
	return nil
}

// MoveCursorLeft moves the cursor left (mock implementation).
func (m *MockTerminal) MoveCursorLeft(n int) error {
	m.record(fmt.Sprintf("MoveCursorLeft(%d)", n))
	return nil
}

// MoveCursorRight moves the cursor right (mock implementation).
func (m *MockTerminal) MoveCursorRight(n int) error {
	m.record(fmt.Sprintf("MoveCursorRight(%d)", n))
	return nil
}

// SaveCursorPosition saves the current cursor position (mock implementation).
func (m *MockTerminal) SaveCursorPosition() error {
	m.record("SaveCursorPosition")
	return nil
}

// RestoreCursorPosition restores the saved cursor position (mock implementation).
func (m *MockTerminal) RestoreCursorPosition() error {
	m.record("RestoreCursorPosition")
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Visibility & Style                                   │
// └─────────────────────────────────────────────────────────────┘

// HideCursor hides the cursor (mock implementation).
func (m *MockTerminal) HideCursor() error {
	m.record("HideCursor")
	return nil
}

// ShowCursor shows the cursor (mock implementation).
func (m *MockTerminal) ShowCursor() error {
	m.record("ShowCursor")
	return nil
}

// SetCursorStyle sets the cursor style (mock implementation).
func (m *MockTerminal) SetCursorStyle(style api.CursorStyle) error {
	m.record(fmt.Sprintf("SetCursorStyle(%s)", style))
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Operations                                           │
// └─────────────────────────────────────────────────────────────┘

// Clear clears the screen (mock implementation).
func (m *MockTerminal) Clear() error {
	m.record("Clear")
	return nil
}

// ClearLine clears the current line (mock implementation).
func (m *MockTerminal) ClearLine() error {
	m.record("ClearLine")
	return nil
}

// ClearFromCursor clears from cursor to end of screen (mock implementation).
func (m *MockTerminal) ClearFromCursor() error {
	m.record("ClearFromCursor")
	return nil
}

// ClearLines clears specified number of lines (mock implementation).
func (m *MockTerminal) ClearLines(count int) error {
	m.record(fmt.Sprintf("ClearLines(%d)", count))
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Output                                                      │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) Write(s string) error {
	m.record(fmt.Sprintf("Write(%q)", s))
	return nil
}

// WriteAt writes text at specified position (mock implementation).
func (m *MockTerminal) WriteAt(x, y int, s string) error {
	m.record(fmt.Sprintf("WriteAt(%d, %d, %q)", x, y, s))
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Buffer (Windows Console API only)                    │
// └─────────────────────────────────────────────────────────────┘

// ReadScreenBuffer reads the screen buffer (mock implementation).
func (m *MockTerminal) ReadScreenBuffer() ([][]rune, error) {
	m.record("ReadScreenBuffer")
	return nil, nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Info                                               │
// └─────────────────────────────────────────────────────────────┘

// Size returns the terminal size (mock implementation).
func (m *MockTerminal) Size() (width, height int, err error) {
	m.record("Size")
	return 80, 24, nil
}

// ColorDepth returns the color depth (mock implementation).
func (m *MockTerminal) ColorDepth() int {
	m.record("ColorDepth")
	return 256
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Capabilities Discovery                                      │
// └─────────────────────────────────────────────────────────────┘

// SupportsDirectPositioning returns whether direct positioning is supported (mock implementation).
func (m *MockTerminal) SupportsDirectPositioning() bool {
	m.record("SupportsDirectPositioning")
	return false
}

// SupportsReadback returns whether readback is supported (mock implementation).
func (m *MockTerminal) SupportsReadback() bool {
	m.record("SupportsReadback")
	return false
}

// SupportsTrueColor returns whether true color is supported (mock implementation).
func (m *MockTerminal) SupportsTrueColor() bool {
	m.record("SupportsTrueColor")
	return true
}

// Platform returns the platform type (mock implementation).
func (m *MockTerminal) Platform() api.Platform {
	m.record("Platform")
	return api.PlatformUnknown
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Alternate Screen Buffer                                     │
// └─────────────────────────────────────────────────────────────┘

// EnterAltScreen enters alternate screen buffer (mock implementation).
func (m *MockTerminal) EnterAltScreen() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "EnterAltScreen")

	if m.inAltScreen {
		return fmt.Errorf("already in alternate screen")
	}

	m.inAltScreen = true
	return nil
}

// ExitAltScreen exits alternate screen buffer (mock implementation).
func (m *MockTerminal) ExitAltScreen() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "ExitAltScreen")

	if !m.inAltScreen {
		return fmt.Errorf("not in alternate screen")
	}

	m.inAltScreen = false
	return nil
}

// IsInAltScreen returns whether in alternate screen buffer (mock implementation).
func (m *MockTerminal) IsInAltScreen() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "IsInAltScreen")
	return m.inAltScreen
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Mode (Raw vs Cooked)                               │
// └─────────────────────────────────────────────────────────────┘

// IsInRawMode returns whether in raw mode (mock implementation).
func (m *MockTerminal) IsInRawMode() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "IsInRawMode")
	return m.inRawMode
}

// EnterRawMode enters raw mode (mock implementation).
func (m *MockTerminal) EnterRawMode() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "EnterRawMode")

	if m.inRawMode {
		return fmt.Errorf("terminal: already in raw mode")
	}

	m.inRawMode = true
	return nil
}

// ExitRawMode exits raw mode (mock implementation).
func (m *MockTerminal) ExitRawMode() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, "ExitRawMode")

	if !m.inRawMode {
		return fmt.Errorf("terminal: not in raw mode")
	}

	m.inRawMode = false
	return nil
}
