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
	mu    sync.Mutex
	Calls []string // All recorded method calls with arguments
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

func (m *MockTerminal) SetCursorPosition(x, y int) error {
	m.record(fmt.Sprintf("SetCursorPosition(%d, %d)", x, y))
	return nil
}

func (m *MockTerminal) GetCursorPosition() (x, y int, err error) {
	m.record("GetCursorPosition")
	return 0, 0, nil
}

func (m *MockTerminal) MoveCursorUp(n int) error {
	m.record(fmt.Sprintf("MoveCursorUp(%d)", n))
	return nil
}

func (m *MockTerminal) MoveCursorDown(n int) error {
	m.record(fmt.Sprintf("MoveCursorDown(%d)", n))
	return nil
}

func (m *MockTerminal) MoveCursorLeft(n int) error {
	m.record(fmt.Sprintf("MoveCursorLeft(%d)", n))
	return nil
}

func (m *MockTerminal) MoveCursorRight(n int) error {
	m.record(fmt.Sprintf("MoveCursorRight(%d)", n))
	return nil
}

func (m *MockTerminal) SaveCursorPosition() error {
	m.record("SaveCursorPosition")
	return nil
}

func (m *MockTerminal) RestoreCursorPosition() error {
	m.record("RestoreCursorPosition")
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Cursor Visibility & Style                                   │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) HideCursor() error {
	m.record("HideCursor")
	return nil
}

func (m *MockTerminal) ShowCursor() error {
	m.record("ShowCursor")
	return nil
}

func (m *MockTerminal) SetCursorStyle(style api.CursorStyle) error {
	m.record(fmt.Sprintf("SetCursorStyle(%s)", style))
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Operations                                           │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) Clear() error {
	m.record("Clear")
	return nil
}

func (m *MockTerminal) ClearLine() error {
	m.record("ClearLine")
	return nil
}

func (m *MockTerminal) ClearFromCursor() error {
	m.record("ClearFromCursor")
	return nil
}

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

func (m *MockTerminal) WriteAt(x, y int, s string) error {
	m.record(fmt.Sprintf("WriteAt(%d, %d, %q)", x, y, s))
	return nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Screen Buffer (Windows Console API only)                    │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) ReadScreenBuffer() ([][]rune, error) {
	m.record("ReadScreenBuffer")
	return nil, nil
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Info                                               │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) Size() (width, height int, err error) {
	m.record("Size")
	return 80, 24, nil
}

func (m *MockTerminal) ColorDepth() int {
	m.record("ColorDepth")
	return 256
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Capabilities Discovery                                      │
// └─────────────────────────────────────────────────────────────┘

func (m *MockTerminal) SupportsDirectPositioning() bool {
	m.record("SupportsDirectPositioning")
	return false
}

func (m *MockTerminal) SupportsReadback() bool {
	m.record("SupportsReadback")
	return false
}

func (m *MockTerminal) SupportsTrueColor() bool {
	m.record("SupportsTrueColor")
	return true
}

func (m *MockTerminal) Platform() api.Platform {
	m.record("Platform")
	return api.PlatformUnknown
}
