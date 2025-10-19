package model_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/domain/model"
	"github.com/phoenix-tui/phoenix/core/domain/value"
)

func TestNewTerminal(t *testing.T) {
	tests := []struct {
		name     string
		caps     *value.Capabilities
		wantANSI bool
	}{
		{
			name:     "with full capabilities",
			caps:     value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true),
			wantANSI: true,
		},
		{
			name:     "with no capabilities",
			caps:     value.NewCapabilities(false, value.ColorDepthNone, false, false, false),
			wantANSI: false,
		},
		{
			name:     "nil capabilities (defensive)",
			caps:     nil,
			wantANSI: false, // Should create minimal capabilities
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term := model.NewTerminal(tt.caps)

			if term == nil {
				t.Fatal("expected terminal, got nil")
			}

			if term.Capabilities() == nil {
				t.Error("capabilities should never be nil")
			}

			if term.IsRawMode() {
				t.Error("new terminal should not be in raw mode")
			}

			// Check default size
			size := term.Size()
			if size.Width != 80 || size.Height != 24 {
				t.Errorf("expected default size 80x24, got %dx%d", size.Width, size.Height)
			}

			if term.SupportsANSI() != tt.wantANSI {
				t.Errorf("expected ANSI %v, got %v", tt.wantANSI, term.SupportsANSI())
			}
		})
	}
}

func TestTerminal_WithSize(t *testing.T) {
	caps := value.NewCapabilities(true, value.ColorDepth256, true, true, true)
	term := model.NewTerminal(caps)

	tests := []struct {
		name    string
		newSize value.Size
	}{
		{
			name:    "different size",
			newSize: value.NewSize(90, 30),
		},
		{
			name:    "large terminal",
			newSize: value.NewSize(200, 60),
		},
		{
			name:    "small terminal",
			newSize: value.NewSize(40, 12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTerm := term.WithSize(tt.newSize)

			// Check new terminal has updated size
			if !newTerm.Size().Equal(tt.newSize) {
				t.Errorf("expected size %v, got %v", tt.newSize, newTerm.Size())
			}

			// Check original unchanged (immutability)
			if newTerm == term {
				t.Error("WithSize should return new instance")
			}

			if term.Size().Equal(tt.newSize) {
				t.Error("original terminal was mutated")
			}
		})
	}
}

func TestTerminal_WithRawMode(t *testing.T) {
	caps := value.NewCapabilities(true, value.ColorDepth256, true, true, true)
	term := model.NewTerminal(caps)

	// Create raw mode entity
	rawMode, err := model.NewRawMode("original-state")
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}
	rawMode = rawMode.Enable()

	// Apply raw mode to terminal
	termRaw := term.WithRawMode(rawMode)

	// Check raw mode is set
	if !termRaw.IsRawMode() {
		t.Error("terminal should be in raw mode")
	}

	if termRaw.RawMode() == nil {
		t.Error("raw mode entity should not be nil")
	}

	// Check original unchanged (immutability)
	if term.IsRawMode() {
		t.Error("original terminal was mutated")
	}

	if termRaw == term {
		t.Error("WithRawMode should return new instance")
	}
}

func TestTerminal_WithoutRawMode(t *testing.T) {
	caps := value.NewCapabilities(true, value.ColorDepth256, true, true, true)
	term := model.NewTerminal(caps)

	// Put into raw mode first
	rawMode, _ := model.NewRawMode("original-state")
	rawMode = rawMode.Enable()
	termRaw := term.WithRawMode(rawMode)

	// Verify in raw mode
	if !termRaw.IsRawMode() {
		t.Fatal("terminal should be in raw mode")
	}

	// Disable raw mode
	termNormal := termRaw.WithoutRawMode()

	// Check raw mode is cleared
	if termNormal.IsRawMode() {
		t.Error("terminal should not be in raw mode")
	}

	if termNormal.RawMode() != nil {
		t.Error("raw mode entity should be nil")
	}

	// Check original unchanged (immutability)
	if !termRaw.IsRawMode() {
		t.Error("original terminal was mutated")
	}

	if termNormal == termRaw {
		t.Error("WithoutRawMode should return new instance")
	}
}

func TestTerminal_ConvenienceMethods(t *testing.T) {
	tests := []struct {
		name           string
		caps           *value.Capabilities
		wantANSI       bool
		wantColor      bool
		wantColorDepth value.ColorDepth
	}{
		{
			name:           "full capabilities",
			caps:           value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true),
			wantANSI:       true,
			wantColor:      true,
			wantColorDepth: value.ColorDepthTrueColor,
		},
		{
			name:           "no capabilities",
			caps:           value.NewCapabilities(false, value.ColorDepthNone, false, false, false),
			wantANSI:       false,
			wantColor:      false,
			wantColorDepth: value.ColorDepthNone,
		},
		{
			name:           "ANSI but no color",
			caps:           value.NewCapabilities(true, value.ColorDepthNone, false, false, false),
			wantANSI:       true,
			wantColor:      false,
			wantColorDepth: value.ColorDepthNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term := model.NewTerminal(tt.caps)

			if term.SupportsANSI() != tt.wantANSI {
				t.Errorf("SupportsANSI: expected %v, got %v", tt.wantANSI, term.SupportsANSI())
			}

			if term.SupportsColor() != tt.wantColor {
				t.Errorf("SupportsColor: expected %v, got %v", tt.wantColor, term.SupportsColor())
			}

			if term.ColorDepth() != tt.wantColorDepth {
				t.Errorf("ColorDepth: expected %v, got %v", tt.wantColorDepth, term.ColorDepth())
			}
		})
	}
}

func TestTerminal_Immutability(t *testing.T) {
	// Verify Terminal is truly immutable across operations
	caps := value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	term1 := model.NewTerminal(caps)

	// Store original state
	originalSize := term1.Size()
	originalRawMode := term1.IsRawMode()
	originalCaps := term1.Capabilities()

	// Perform multiple operations
	term2 := term1.WithSize(value.NewSize(100, 50))
	rawMode, _ := model.NewRawMode("state")
	term3 := term2.WithRawMode(rawMode.Enable())
	term4 := term3.WithoutRawMode()

	// Verify term1 unchanged
	if !term1.Size().Equal(originalSize) {
		t.Error("term1 size was mutated")
	}

	if term1.IsRawMode() != originalRawMode {
		t.Error("term1 raw mode was mutated")
	}

	if term1.Capabilities() != originalCaps {
		t.Error("term1 capabilities were mutated")
	}

	// Verify all instances are different
	if term1 == term2 || term1 == term3 || term1 == term4 {
		t.Error("operations returned same instance (not immutable)")
	}
}
