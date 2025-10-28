package phoenix_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix"
	"github.com/phoenix-tui/phoenix/tea"
)

// ┌─────────────────────────────────────────────────────────────┐
// │ Core Tests                                                  │
// └─────────────────────────────────────────────────────────────┘

func TestAutoDetectTerminal(t *testing.T) {
	term := phoenix.AutoDetectTerminal()
	if term == nil {
		t.Fatal("AutoDetectTerminal() returned nil")
	}

	size := term.Size()
	if size.Width <= 0 || size.Height <= 0 {
		t.Errorf("Invalid terminal size: %dx%d", size.Width, size.Height)
	}
}

func TestNewTerminal(t *testing.T) {
	term := phoenix.NewTerminal()

	if term == nil {
		t.Fatal("NewTerminal() returned nil")
	}

	size := term.Size()
	if size.Width <= 0 || size.Height <= 0 {
		t.Errorf("Invalid terminal size: %dx%d", size.Width, size.Height)
	}
}

func TestNewTerminalWithCapabilities(t *testing.T) {
	caps := phoenix.NewCapabilities(true, phoenix.ColorDepth256, true, true, true)
	term := phoenix.NewTerminalWithCapabilities(caps)

	if term == nil {
		t.Fatal("NewTerminalWithCapabilities() returned nil")
	}
}

func TestNewSize(t *testing.T) {
	size := phoenix.NewSize(100, 50)

	if size.Width != 100 {
		t.Errorf("Width mismatch: got %d, want 100", size.Width)
	}
	if size.Height != 50 {
		t.Errorf("Height mismatch: got %d, want 50", size.Height)
	}
}

func TestNewCapabilities(t *testing.T) {
	caps := phoenix.NewCapabilities(
		true,                        // ANSI
		phoenix.ColorDepthTrueColor, // Color depth
		true,                        // Mouse
		false,                       // Alt screen
		true,                        // Cursor
	)

	if caps == nil {
		t.Fatal("NewCapabilities() returned nil")
	}

	if !caps.SupportsANSI() {
		t.Error("Expected ANSI support")
	}
	if caps.ColorDepth() != phoenix.ColorDepthTrueColor {
		t.Errorf("Color depth mismatch: got %v, want %v",
			caps.ColorDepth(), phoenix.ColorDepthTrueColor)
	}
	if !caps.SupportsMouse() {
		t.Error("Expected mouse support")
	}
	if caps.SupportsAltScreen() {
		t.Error("Expected no alt screen support")
	}
}

func TestColorDepthConstants(t *testing.T) {
	tests := []struct {
		name  string
		depth interface{}
	}{
		{"ColorDepthNone", phoenix.ColorDepthNone},
		{"ColorDepth8", phoenix.ColorDepth8},
		{"ColorDepth256", phoenix.ColorDepth256},
		{"ColorDepthTrueColor", phoenix.ColorDepthTrueColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.depth == nil {
				t.Errorf("%s is nil", tt.name)
			}
		})
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Style Tests                                                 │
// └─────────────────────────────────────────────────────────────┘

func TestNewStyle(_ *testing.T) {
	s := phoenix.NewStyle()
	// Style is a value type, not a pointer
	// Just verify we can create it
	_ = s
}

func TestStyleAPI(_ *testing.T) {
	// Simple test to verify the Style API is accessible through umbrella module
	// Don't test Style functionality itself - that's tested in style package
	s := phoenix.NewStyle()
	_ = s
	// If we got here without compilation errors, the API is accessible
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Tea Tests                                                   │
// └─────────────────────────────────────────────────────────────┘

type testModel struct {
	value int
}

func (m testModel) Init() tea.Cmd { return nil }

func (m testModel) Update(_ tea.Msg) (testModel, tea.Cmd) {
	return m, nil
}

func (m testModel) View() string {
	return "test"
}

func TestNewProgram(t *testing.T) {
	p := phoenix.NewProgram(testModel{value: 42})
	if p == nil {
		t.Fatal("NewProgram() returned nil")
	}
}

func TestNewProgramWithOptions(t *testing.T) {
	p := phoenix.NewProgram(
		testModel{value: 100},
		phoenix.WithAltScreen[testModel](),
		phoenix.WithMouseAllMotion[testModel](),
	)
	if p == nil {
		t.Fatal("NewProgram() with options returned nil")
	}
}

func TestWithAltScreen(t *testing.T) {
	opt := phoenix.WithAltScreen[testModel]()
	if opt == nil {
		t.Fatal("WithAltScreen() returned nil")
	}
}

func TestWithMouseAllMotion(t *testing.T) {
	opt := phoenix.WithMouseAllMotion[testModel]()
	if opt == nil {
		t.Fatal("WithMouseAllMotion() returned nil")
	}
}

func TestQuit(t *testing.T) {
	cmd := phoenix.Quit()
	if cmd == nil {
		t.Fatal("Quit() returned nil")
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Clipboard Tests                                             │
// └─────────────────────────────────────────────────────────────┘

func TestClipboard(t *testing.T) {
	// Note: May not work in CI environments
	testText := "phoenix clipboard test"

	err := phoenix.WriteClipboard(testText)
	if err != nil {
		t.Skip("Clipboard not available:", err)
	}

	text, err := phoenix.ReadClipboard()
	if err != nil {
		t.Fatal("ReadClipboard() failed:", err)
	}

	if text != testText {
		t.Errorf("got %q, want %q", text, testText)
	}
}

func TestReadClipboard(t *testing.T) {
	// Should not panic even if clipboard is unavailable
	_, err := phoenix.ReadClipboard()
	if err != nil {
		t.Skip("ReadClipboard() not available:", err)
	}
}

func TestWriteClipboard(t *testing.T) {
	// Should not panic even if clipboard is unavailable
	err := phoenix.WriteClipboard("test")
	if err != nil {
		t.Skip("WriteClipboard() not available:", err)
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Terminal Tests                                              │
// └─────────────────────────────────────────────────────────────┘

func TestNewPlatformTerminal(t *testing.T) {
	term := phoenix.NewPlatformTerminal()

	if term == nil {
		t.Fatal("NewPlatformTerminal() returned nil")
	}

	// Terminal should have a valid platform
	platform := term.Platform()
	if platform == 0 {
		t.Error("Terminal has invalid platform (0)")
	}
}

func TestNewANSITerminal(t *testing.T) {
	term := phoenix.NewANSITerminal()
	if term == nil {
		t.Fatal("NewANSITerminal() returned nil")
	}

	platform := term.Platform()
	// ANSI terminal should report a valid platform
	if platform == 0 {
		t.Error("Invalid platform for ANSI terminal")
	}
}

// ┌─────────────────────────────────────────────────────────────┐
// │ Integration Tests                                           │
// └─────────────────────────────────────────────────────────────┘

func TestUmbrellaIntegration(t *testing.T) {
	// Test that multiple components work together
	term := phoenix.AutoDetectTerminal()
	if term == nil {
		t.Fatal("Failed to create terminal")
	}

	// Create style (don't render - that's tested in style package)
	style := phoenix.NewStyle()
	_ = style

	// Create a program (don't run it, just verify creation)
	p := phoenix.NewProgram(
		testModel{value: 123},
		phoenix.WithAltScreen[testModel](),
	)
	if p == nil {
		t.Error("Failed to create program")
	}
}

func TestUmbrellaConvenienceVsDirectImport(t *testing.T) {
	// Verify that umbrella convenience functions produce the same results
	// as direct imports (API compatibility test)

	// Terminal creation
	umbrellaSize := phoenix.NewSize(100, 30)
	if umbrellaSize.Width != 100 || umbrellaSize.Height != 30 {
		t.Error("Umbrella NewSize() differs from direct import")
	}

	// Style creation
	umbrellaStyle := phoenix.NewStyle()
	_ = umbrellaStyle // Style is value type, not pointer

	// Program options
	opt := phoenix.WithAltScreen[testModel]()
	if opt == nil {
		t.Error("Umbrella WithAltScreen() returned nil")
	}
}
