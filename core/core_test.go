package core_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core"
)

// TestNewTerminal verifies default terminal creation
func TestNewTerminal(t *testing.T) {
	term := core.NewTerminal()

	if term == nil {
		t.Fatal("expected non-nil terminal")
	}

	// Default should have no capabilities
	caps := term.Capabilities()
	if caps.SupportsANSI() {
		t.Error("default terminal should not support ANSI")
	}
	if caps.SupportsColor() {
		t.Error("default terminal should not support color")
	}

	// Default size should be VT100 (80x24)
	size := term.Size()
	if size.Width != 80 || size.Height != 24 {
		t.Errorf("expected size 80x24, got %dx%d", size.Width, size.Height)
	}

	// Raw mode should be disabled by default
	if term.IsRawModeEnabled() {
		t.Error("raw mode should be disabled by default")
	}
}

// TestAutoDetect verifies automatic capability detection
func TestAutoDetect(t *testing.T) {
	// This will use real environment variables
	term := core.AutoDetect()

	if term == nil {
		t.Fatal("expected non-nil terminal")
	}

	// Should have capabilities based on actual environment
	caps := term.Capabilities()
	if caps == nil {
		t.Fatal("expected non-nil capabilities")
	}

	// Verify capabilities are coherent
	if caps.SupportsMouse() && !caps.SupportsANSI() {
		t.Error("mouse support requires ANSI support")
	}
	if caps.SupportsAltScreen() && !caps.SupportsANSI() {
		t.Error("alt screen requires ANSI support")
	}
	if caps.SupportsCursorControl() && !caps.SupportsANSI() {
		t.Error("cursor control requires ANSI support")
	}
}

// TestNewTerminalWithCapabilities verifies custom capabilities
func TestNewTerminalWithCapabilities(t *testing.T) {
	caps := core.NewCapabilities(true, core.ColorDepth256, true, true, true)
	term := core.NewTerminalWithCapabilities(caps)

	if term == nil {
		t.Fatal("expected non-nil terminal")
	}

	termCaps := term.Capabilities()
	if !termCaps.SupportsANSI() {
		t.Error("expected ANSI support")
	}
	if termCaps.ColorDepth() != core.ColorDepth256 {
		t.Errorf("expected 256 color depth, got %v", termCaps.ColorDepth())
	}
	if !termCaps.SupportsMouse() {
		t.Error("expected mouse support")
	}
}

// TestNewTerminalWithCapabilities_Nil verifies nil handling
func TestNewTerminalWithCapabilities_Nil(t *testing.T) {
	term := core.NewTerminalWithCapabilities(nil)

	if term == nil {
		t.Fatal("expected non-nil terminal")
	}

	// Should behave like NewTerminal()
	caps := term.Capabilities()
	if caps.SupportsANSI() {
		t.Error("nil capabilities should default to no ANSI")
	}
}

// TestTerminal_WithSize verifies immutable size changes
func TestTerminal_WithSize(t *testing.T) {
	original := core.NewTerminal()
	originalSize := original.Size()

	// Change size
	newSize := core.NewSize(120, 40)
	resized := original.WithSize(newSize)

	// Original should be unchanged
	if original.Size().Width != originalSize.Width || original.Size().Height != originalSize.Height {
		t.Error("original terminal was mutated")
	}

	// New terminal should have new size
	if resized.Size().Width != 120 || resized.Size().Height != 40 {
		t.Errorf("expected size 120x40, got %dx%d", resized.Size().Width, resized.Size().Height)
	}
}

// TestTerminal_WithRawMode verifies raw mode management
func TestTerminal_WithRawMode(t *testing.T) {
	original := core.NewTerminal()

	// Create raw mode
	rawMode, err := core.NewRawMode("fake state")
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}

	// Attach raw mode to terminal
	withRaw := original.WithRawMode(rawMode)

	// Original should not be affected
	if original.IsRawModeEnabled() {
		t.Error("original terminal was mutated")
	}

	// Verify new terminal has raw mode attached
	if withRaw == nil {
		t.Fatal("expected non-nil terminal")
	}

	// Note: RawMode itself doesn't change enabled state until Enable() is called
	// WithRawMode just attaches the RawMode object to the terminal
}

// TestTerminal_WithRawMode_Nil verifies nil handling
func TestTerminal_WithRawMode_Nil(t *testing.T) {
	original := core.NewTerminal()
	withNil := original.WithRawMode(nil)

	// Should return same terminal (no-op)
	if withNil == nil {
		t.Fatal("expected non-nil terminal")
	}

	// Should not enable raw mode
	if withNil.IsRawModeEnabled() {
		t.Error("nil raw mode should not enable raw mode")
	}
}

// TestSize verifies Size creation and validation
func TestSize(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		expectWidth  int
		expectHeight int
	}{
		{"valid", 80, 24, 80, 24},
		{"large", 200, 100, 200, 100},
		{"minimum", 1, 1, 1, 1},
		{"zero width clamped", 0, 24, 1, 24},
		{"zero height clamped", 80, 0, 80, 1},
		{"negative width clamped", -10, 24, 1, 24},
		{"negative height clamped", 80, -10, 80, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := core.NewSize(tt.width, tt.height)

			if size.Width != tt.expectWidth {
				t.Errorf("expected width %d, got %d", tt.expectWidth, size.Width)
			}
			if size.Height != tt.expectHeight {
				t.Errorf("expected height %d, got %d", tt.expectHeight, size.Height)
			}
		})
	}
}

// TestPosition verifies Position creation and validation
func TestPosition(t *testing.T) {
	tests := []struct {
		name      string
		row       int
		col       int
		expectRow int
		expectCol int
	}{
		{"valid", 10, 20, 10, 20},
		{"zero", 0, 0, 0, 0},
		{"negative row clamped", -5, 20, 0, 20},
		{"negative col clamped", 10, -5, 10, 0},
		{"both negative clamped", -5, -10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := core.NewPosition(tt.row, tt.col)

			if pos.Row != tt.expectRow {
				t.Errorf("expected row %d, got %d", tt.expectRow, pos.Row)
			}
			if pos.Col != tt.expectCol {
				t.Errorf("expected col %d, got %d", tt.expectCol, pos.Col)
			}
		})
	}
}

// TestPosition_Add verifies position arithmetic
func TestPosition_Add(t *testing.T) {
	tests := []struct {
		name      string
		start     core.Position
		deltaRow  int
		deltaCol  int
		expectRow int
		expectCol int
	}{
		{"positive delta", core.NewPosition(10, 20), 5, 10, 15, 30},
		{"negative delta", core.NewPosition(10, 20), -3, -5, 7, 15},
		{"zero delta", core.NewPosition(10, 20), 0, 0, 10, 20},
		{"negative result clamped", core.NewPosition(5, 5), -10, -10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Add(tt.deltaRow, tt.deltaCol)

			if result.Row != tt.expectRow {
				t.Errorf("expected row %d, got %d", tt.expectRow, result.Row)
			}
			if result.Col != tt.expectCol {
				t.Errorf("expected col %d, got %d", tt.expectCol, result.Col)
			}

			// Verify immutability
			if tt.start.Row != core.NewPosition(tt.start.Row, tt.start.Col).Row {
				t.Error("original position was mutated")
			}
		})
	}
}

// TestCell verifies Cell creation with manual width
func TestCell(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		width       int
		expectWidth int
	}{
		{"ascii", "A", 1, 1},
		{"emoji", "😀", 2, 2},
		{"zero width", "x", 0, 0},
		{"negative width clamped", "x", -1, 0},
		{"wide character", "あ", 2, 2},
		{"manual incorrect", "👋", 3, 3}, // Wrong width, but user controls
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := core.NewCell(tt.content, tt.width)

			if cell.Content != tt.content {
				t.Errorf("expected content %q, got %q", tt.content, cell.Content)
			}
			if cell.Width != tt.expectWidth {
				t.Errorf("expected width %d, got %d", tt.expectWidth, cell.Width)
			}
		})
	}
}

// TestNewCellAuto_API verifies automatic width calculation at API level
func TestNewCellAuto_API(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
	}{
		// ASCII
		{"ASCII single", "A", 1},
		{"ASCII word", "Hello", 5},
		{"ASCII empty", "", 0},
		{"ASCII space", " ", 1},
		{"ASCII numbers", "123", 3},

		// Emoji
		{"Emoji simple", "😀", 2},
		{"Emoji wave", "👋", 2},
		{"Emoji with modifier", "👋🏻", 2},
		{"Emoji multiple", "😀😃", 4},
		{"Emoji flag", "🇺🇸", 2},
		{"Emoji heart", "❤️", 2},

		// CJK
		{"CJK Chinese", "中", 2},
		{"CJK Chinese word", "中文", 4},
		{"CJK Japanese hiragana", "あ", 2},
		{"CJK Korean", "한", 2},

		// Combining
		{"Combining e acute", "é", 1},
		{"Combining a umlaut", "ä", 1},
		{"Combining word", "café", 4},

		// Mixed
		{"Mixed ASCII+Emoji", "Hi 👋", 5},    // H(1) + i(1) + space(1) + 👋(2)
		{"Mixed ASCII+CJK", "Hello中", 7},    // Hello(5) + 中(2)
		{"Mixed Emoji+CJK", "👋中", 4},        // 👋(2) + 中(2)
		{"Mixed All", "Hi👋中", 6},            // H(1) + i(1) + 👋(2) + 中(2)
		{"Mixed Complex", "Hello 👋 世界", 13}, // Hello(5) + space(1) + 👋(2) + space(1) + 世(2) + 界(2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := core.NewCellAuto(tt.content)
			if cell.Width != tt.width {
				t.Errorf("NewCellAuto(%q) width = %d, want %d", tt.content, cell.Width, tt.width)
			}
			if cell.Content != tt.content {
				t.Errorf("NewCellAuto(%q) content = %q, want %q", tt.content, cell.Content, tt.content)
			}
		})
	}
}

// TestNewCellAuto_API_EdgeCases verifies edge cases
func TestNewCellAuto_API_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{"zero-width space", "\u200b", 0},
		{"newline", "\n", 0},
		{"tab", "\t", 0},
		{"multiple spaces", "   ", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := core.NewCellAuto(tt.content)
			if cell.Width != tt.width {
				t.Errorf("NewCellAuto(%q) width = %d, want %d", tt.content, cell.Width, tt.width)
			}
		})
	}
}

// TestNewCellAuto_API_BackwardsCompatibility verifies NewCell still works
func TestNewCellAuto_API_BackwardsCompatibility(t *testing.T) {
	// NewCell should still work for manual width control
	manualCell := core.NewCell("👋", 3) // Wrong width, but user controls
	if manualCell.Width != 3 {
		t.Errorf("NewCell manual width not respected: got %d, want 3", manualCell.Width)
	}

	// NewCellAuto should calculate correctly
	autoCell := core.NewCellAuto("👋")
	if autoCell.Width != 2 {
		t.Errorf("NewCellAuto width incorrect: got %d, want 2", autoCell.Width)
	}

	// Same content, different widths
	if manualCell.Content != autoCell.Content {
		t.Error("Content should be the same")
	}
	if manualCell.Width == autoCell.Width {
		t.Error("Widths should differ (manual vs auto)")
	}
}

// TestRawMode verifies RawMode creation
func TestRawMode(t *testing.T) {
	state := struct {
		mode uint32
	}{mode: 12345}

	rawMode, err := core.NewRawMode(state)
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}

	if rawMode == nil {
		t.Fatal("expected non-nil raw mode")
	}

	// Should not be enabled initially
	if rawMode.IsEnabled() {
		t.Error("raw mode should not be enabled initially")
	}

	// Original state should be preserved
	if rawMode.OriginalState() == nil {
		t.Error("original state should be preserved")
	}
}

// TestRawMode_NilState verifies nil state handling
func TestRawMode_NilState(t *testing.T) {
	_, err := core.NewRawMode(nil)
	if err == nil {
		t.Error("expected error for nil state")
	}
}

// TestCapabilities_NewCapabilities verifies capability creation
func TestCapabilities_NewCapabilities(t *testing.T) {
	tests := []struct {
		name   string
		ansi   bool
		color  core.ColorDepth
		mouse  bool
		alt    bool
		cursor bool
	}{
		{"all enabled", true, core.ColorDepthTrueColor, true, true, true},
		{"no ansi", false, core.ColorDepthNone, false, false, false},
		{"ansi only", true, core.ColorDepth8, false, false, false},
		{"full color", true, core.ColorDepth256, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caps := core.NewCapabilities(tt.ansi, tt.color, tt.mouse, tt.alt, tt.cursor)

			if caps == nil {
				t.Fatal("expected non-nil capabilities")
			}

			if caps.SupportsANSI() != tt.ansi {
				t.Errorf("expected ANSI %v, got %v", tt.ansi, caps.SupportsANSI())
			}
			if caps.ColorDepth() != tt.color {
				t.Errorf("expected color depth %v, got %v", tt.color, caps.ColorDepth())
			}
			if caps.SupportsMouse() != tt.mouse {
				t.Errorf("expected mouse %v, got %v", tt.mouse, caps.SupportsMouse())
			}
			if caps.SupportsAltScreen() != tt.alt {
				t.Errorf("expected alt screen %v, got %v", tt.alt, caps.SupportsAltScreen())
			}
			if caps.SupportsCursorControl() != tt.cursor {
				t.Errorf("expected cursor control %v, got %v", tt.cursor, caps.SupportsCursorControl())
			}
		})
	}
}

// TestCapabilities_SupportsColor verifies color support detection
func TestCapabilities_SupportsColor(t *testing.T) {
	tests := []struct {
		name       string
		colorDepth core.ColorDepth
		wantColor  bool
	}{
		{"none", core.ColorDepthNone, false},
		{"8 colors", core.ColorDepth8, true},
		{"256 colors", core.ColorDepth256, true},
		{"true color", core.ColorDepthTrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caps := core.NewCapabilities(true, tt.colorDepth, false, false, false)

			if caps.SupportsColor() != tt.wantColor {
				t.Errorf("expected color support %v, got %v", tt.wantColor, caps.SupportsColor())
			}
		})
	}
}

// TestCapabilities_SupportsTrueColor verifies true color detection
func TestCapabilities_SupportsTrueColor(t *testing.T) {
	tests := []struct {
		name          string
		colorDepth    core.ColorDepth
		wantTrueColor bool
	}{
		{"none", core.ColorDepthNone, false},
		{"8 colors", core.ColorDepth8, false},
		{"256 colors", core.ColorDepth256, false},
		{"true color", core.ColorDepthTrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caps := core.NewCapabilities(true, tt.colorDepth, false, false, false)

			if caps.SupportsTrueColor() != tt.wantTrueColor {
				t.Errorf("expected true color %v, got %v", tt.wantTrueColor, caps.SupportsTrueColor())
			}
		})
	}
}

// TestColorDepth_String verifies string representation
func TestColorDepth_String(t *testing.T) {
	tests := []struct {
		depth core.ColorDepth
		want  string
	}{
		{core.ColorDepthNone, "None"},
		{core.ColorDepth8, "8 colors"},
		{core.ColorDepth256, "256 colors"},
		{core.ColorDepthTrueColor, "TrueColor (16.7M colors)"},
		{core.ColorDepth(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.depth.String()
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// TestCapabilities_BusinessRules verifies ANSI dependency rules
func TestCapabilities_BusinessRules(t *testing.T) {
	// Business rule: mouse/alt/cursor require ANSI support
	caps := core.NewCapabilities(false, core.ColorDepthNone, true, true, true)

	// All should be disabled since ANSI is false
	if caps.SupportsMouse() {
		t.Error("mouse should be disabled without ANSI")
	}
	if caps.SupportsAltScreen() {
		t.Error("alt screen should be disabled without ANSI")
	}
	if caps.SupportsCursorControl() {
		t.Error("cursor control should be disabled without ANSI")
	}
}
