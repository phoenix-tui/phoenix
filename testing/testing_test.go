package testing

import (
	"sync"
	"testing"

	"github.com/phoenix-tui/phoenix/terminal/api"
)

// ┌─────────────────────────────────────────────────────────────┐
// │ NullTerminal Tests                                          │
// └─────────────────────────────────────────────────────────────┘

func TestNullTerminal_ImplementsTerminalInterface(_ *testing.T) {
	var _ api.Terminal = (*NullTerminal)(nil)
}

func TestNullTerminal_AllMethodsReturnNil(t *testing.T) {
	term := NewNullTerminal()

	// Cursor operations
	if err := term.SetCursorPosition(10, 5); err != nil {
		t.Errorf("SetCursorPosition() = %v, want nil", err)
	}
	if _, _, err := term.GetCursorPosition(); err != nil {
		t.Errorf("GetCursorPosition() = %v, want nil", err)
	}
	if err := term.MoveCursorUp(1); err != nil {
		t.Errorf("MoveCursorUp() = %v, want nil", err)
	}
	if err := term.MoveCursorDown(1); err != nil {
		t.Errorf("MoveCursorDown() = %v, want nil", err)
	}
	if err := term.MoveCursorLeft(1); err != nil {
		t.Errorf("MoveCursorLeft() = %v, want nil", err)
	}
	if err := term.MoveCursorRight(1); err != nil {
		t.Errorf("MoveCursorRight() = %v, want nil", err)
	}
	if err := term.SaveCursorPosition(); err != nil {
		t.Errorf("SaveCursorPosition() = %v, want nil", err)
	}
	if err := term.RestoreCursorPosition(); err != nil {
		t.Errorf("RestoreCursorPosition() = %v, want nil", err)
	}

	// Cursor visibility
	if err := term.HideCursor(); err != nil {
		t.Errorf("HideCursor() = %v, want nil", err)
	}
	if err := term.ShowCursor(); err != nil {
		t.Errorf("ShowCursor() = %v, want nil", err)
	}
	if err := term.SetCursorStyle(api.CursorBlock); err != nil {
		t.Errorf("SetCursorStyle() = %v, want nil", err)
	}

	// Screen operations
	if err := term.Clear(); err != nil {
		t.Errorf("Clear() = %v, want nil", err)
	}
	if err := term.ClearLine(); err != nil {
		t.Errorf("ClearLine() = %v, want nil", err)
	}
	if err := term.ClearFromCursor(); err != nil {
		t.Errorf("ClearFromCursor() = %v, want nil", err)
	}
	if err := term.ClearLines(5); err != nil {
		t.Errorf("ClearLines() = %v, want nil", err)
	}

	// Output
	if err := term.Write("test"); err != nil {
		t.Errorf("Write() = %v, want nil", err)
	}
	if err := term.WriteAt(0, 0, "test"); err != nil {
		t.Errorf("WriteAt() = %v, want nil", err)
	}

	// Screen buffer
	if _, err := term.ReadScreenBuffer(); err != nil {
		t.Errorf("ReadScreenBuffer() = %v, want nil", err)
	}

	// Terminal info
	if _, _, err := term.Size(); err != nil {
		t.Errorf("Size() = %v, want nil", err)
	}
}

func TestNullTerminal_ReasonableDefaults(t *testing.T) {
	term := NewNullTerminal()

	// Size should return reasonable default
	width, height, err := term.Size()
	if err != nil {
		t.Errorf("Size() error = %v, want nil", err)
	}
	if width != 80 || height != 24 {
		t.Errorf("Size() = (%d, %d), want (80, 24)", width, height)
	}

	// ColorDepth should return reasonable value
	depth := term.ColorDepth()
	if depth != 256 {
		t.Errorf("ColorDepth() = %d, want 256", depth)
	}

	// Capabilities should be conservative
	if term.SupportsDirectPositioning() {
		t.Error("SupportsDirectPositioning() = true, want false (conservative)")
	}
	if term.SupportsReadback() {
		t.Error("SupportsReadback() = true, want false (conservative)")
	}
	if !term.SupportsTrueColor() {
		t.Error("SupportsTrueColor() = false, want true (optimistic)")
	}

	// Platform should be unknown
	if platform := term.Platform(); platform != api.PlatformUnknown {
		t.Errorf("Platform() = %v, want PlatformUnknown", platform)
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ MockTerminal Tests                                          │
// └─────────────────────────────────────────────────────────────┘

func TestMockTerminal_ImplementsTerminalInterface(_ *testing.T) {
	var _ api.Terminal = (*MockTerminal)(nil)
}

func TestMockTerminal_RecordsCalls(t *testing.T) {
	mock := NewMockTerminal()

	// Perform some operations
	_ = mock.SetCursorPosition(10, 5)
	_ = mock.ClearLine()
	_ = mock.Write("Hello")
	_ = mock.ClearLines(3)

	// Verify calls were recorded
	expectedCalls := []string{
		"SetCursorPosition(10, 5)",
		"ClearLine",
		`Write("Hello")`,
		"ClearLines(3)",
	}

	if len(mock.Calls) != len(expectedCalls) {
		t.Fatalf("len(Calls) = %d, want %d", len(mock.Calls), len(expectedCalls))
	}

	for i, expected := range expectedCalls {
		if mock.Calls[i] != expected {
			t.Errorf("Calls[%d] = %q, want %q", i, mock.Calls[i], expected)
		}
	}
}

func TestMockTerminal_CallCount(t *testing.T) {
	mock := NewMockTerminal()

	// Call same method multiple times
	_ = mock.ClearLine()
	_ = mock.ClearLine()
	_ = mock.ClearLine()
	_ = mock.Clear()

	if count := mock.CallCount("ClearLine"); count != 3 {
		t.Errorf("CallCount(ClearLine) = %d, want 3", count)
	}

	if count := mock.CallCount("Clear"); count != 1 {
		t.Errorf("CallCount(Clear) = %d, want 1", count)
	}

	if count := mock.CallCount("HideCursor"); count != 0 {
		t.Errorf("CallCount(HideCursor) = %d, want 0", count)
	}
}

func TestMockTerminal_CallCountWithArgs(t *testing.T) {
	mock := NewMockTerminal()

	_ = mock.SetCursorPosition(10, 5)
	_ = mock.SetCursorPosition(20, 10)
	_ = mock.MoveCursorUp(3)

	// CallCount should match method name prefix
	if count := mock.CallCount("SetCursorPosition"); count != 2 {
		t.Errorf("CallCount(SetCursorPosition) = %d, want 2", count)
	}

	if count := mock.CallCount("MoveCursorUp"); count != 1 {
		t.Errorf("CallCount(MoveCursorUp) = %d, want 1", count)
	}
}

func TestMockTerminal_Reset(t *testing.T) {
	mock := NewMockTerminal()

	// Perform some operations
	_ = mock.ClearLine()
	_ = mock.Write("test")

	if len(mock.Calls) != 2 {
		t.Fatalf("len(Calls) before reset = %d, want 2", len(mock.Calls))
	}

	// Reset
	mock.Reset()

	if len(mock.Calls) != 0 {
		t.Errorf("len(Calls) after reset = %d, want 0", len(mock.Calls))
	}

	// Verify new calls are recorded after reset
	_ = mock.ShowCursor()
	if len(mock.Calls) != 1 {
		t.Errorf("len(Calls) after reset and new call = %d, want 1", len(mock.Calls))
	}
}

func TestMockTerminal_ThreadSafety(t *testing.T) {
	mock := NewMockTerminal()

	// Simulate concurrent access
	var wg sync.WaitGroup
	numGoroutines := 100
	numCallsPerGoroutine := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numCallsPerGoroutine; j++ {
				_ = mock.ClearLine()
			}
		}()
	}

	wg.Wait()

	// Verify all calls were recorded
	expectedCalls := numGoroutines * numCallsPerGoroutine
	if len(mock.Calls) != expectedCalls {
		t.Errorf("len(Calls) = %d, want %d", len(mock.Calls), expectedCalls)
	}

	if count := mock.CallCount("ClearLine"); count != expectedCalls {
		t.Errorf("CallCount(ClearLine) = %d, want %d", count, expectedCalls)
	}
}

func TestMockTerminal_AllMethodsRecord(t *testing.T) {
	mock := NewMockTerminal()

	// Test all methods record correctly
	tests := []struct {
		name     string
		call     func()
		expected string
	}{
		{"SetCursorPosition", func() { _ = mock.SetCursorPosition(1, 2) }, "SetCursorPosition(1, 2)"},
		{"GetCursorPosition", func() { _, _, _ = mock.GetCursorPosition() }, "GetCursorPosition"},
		{"MoveCursorUp", func() { _ = mock.MoveCursorUp(3) }, "MoveCursorUp(3)"},
		{"MoveCursorDown", func() { _ = mock.MoveCursorDown(4) }, "MoveCursorDown(4)"},
		{"MoveCursorLeft", func() { _ = mock.MoveCursorLeft(5) }, "MoveCursorLeft(5)"},
		{"MoveCursorRight", func() { _ = mock.MoveCursorRight(6) }, "MoveCursorRight(6)"},
		{"SaveCursorPosition", func() { _ = mock.SaveCursorPosition() }, "SaveCursorPosition"},
		{"RestoreCursorPosition", func() { _ = mock.RestoreCursorPosition() }, "RestoreCursorPosition"},
		{"HideCursor", func() { _ = mock.HideCursor() }, "HideCursor"},
		{"ShowCursor", func() { _ = mock.ShowCursor() }, "ShowCursor"},
		{"SetCursorStyle", func() { _ = mock.SetCursorStyle(api.CursorBlock) }, "SetCursorStyle(Block)"},
		{"Clear", func() { _ = mock.Clear() }, "Clear"},
		{"ClearLine", func() { _ = mock.ClearLine() }, "ClearLine"},
		{"ClearFromCursor", func() { _ = mock.ClearFromCursor() }, "ClearFromCursor"},
		{"ClearLines", func() { _ = mock.ClearLines(7) }, "ClearLines(7)"},
		{"Write", func() { _ = mock.Write("test") }, `Write("test")`},
		{"WriteAt", func() { _ = mock.WriteAt(8, 9, "hello") }, `WriteAt(8, 9, "hello")`},
		{"ReadScreenBuffer", func() { _, _ = mock.ReadScreenBuffer() }, "ReadScreenBuffer"},
		{"Size", func() { _, _, _ = mock.Size() }, "Size"},
		{"ColorDepth", func() { _ = mock.ColorDepth() }, "ColorDepth"},
		{"SupportsDirectPositioning", func() { _ = mock.SupportsDirectPositioning() }, "SupportsDirectPositioning"},
		{"SupportsReadback", func() { _ = mock.SupportsReadback() }, "SupportsReadback"},
		{"SupportsTrueColor", func() { _ = mock.SupportsTrueColor() }, "SupportsTrueColor"},
		{"Platform", func() { _ = mock.Platform() }, "Platform"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Reset()
			tt.call()

			if len(mock.Calls) != 1 {
				t.Fatalf("len(Calls) = %d, want 1", len(mock.Calls))
			}

			if mock.Calls[0] != tt.expected {
				t.Errorf("Calls[0] = %q, want %q", mock.Calls[0], tt.expected)
			}
		})
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Integration Tests (Realistic Usage)                        │
// └─────────────────────────────────────────────────────────────┘

// TestNullTerminal_InRealModel demonstrates using NullTerminal in a model test.
func TestNullTerminal_InRealModel(_ *testing.T) {
	type Model struct {
		terminal api.Terminal
	}

	render := func(m *Model) {
		_ = m.terminal.HideCursor()
		_ = m.terminal.SetCursorPosition(0, 0)
		_ = m.terminal.Write("Hello, World!")
		_ = m.terminal.ShowCursor()
	}

	// Test with NullTerminal - no panics!
	m := &Model{terminal: NewNullTerminal()}
	render(m) // Should not panic or fail
}

// TestMockTerminal_InRealModel demonstrates using MockTerminal to verify calls.
func TestMockTerminal_InRealModel(t *testing.T) {
	type Model struct {
		terminal api.Terminal
	}

	render := func(m *Model) {
		_ = m.terminal.ClearLine()
		_ = m.terminal.Write("Status: Ready")
	}

	mock := NewMockTerminal()
	m := &Model{terminal: mock}
	render(m)

	// Verify expected operations
	if count := mock.CallCount("ClearLine"); count != 1 {
		t.Errorf("Expected ClearLine to be called once, got %d", count)
	}

	if count := mock.CallCount("Write"); count != 1 {
		t.Errorf("Expected Write to be called once, got %d", count)
	}

	// Verify Write was called with correct argument
	foundWrite := false
	for _, call := range mock.Calls {
		if call == `Write("Status: Ready")` {
			foundWrite = true
			break
		}
	}
	if !foundWrite {
		t.Error("Expected Write to be called with 'Status: Ready'")
	}
}
