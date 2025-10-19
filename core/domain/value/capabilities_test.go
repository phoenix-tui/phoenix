package value_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/domain/value"
)

func TestColorDepth_String(t *testing.T) {
	tests := []struct {
		name     string
		depth    value.ColorDepth
		expected string
	}{
		{
			name:     "no color",
			depth:    value.ColorDepthNone,
			expected: "no-color",
		},
		{
			name:     "8 colors",
			depth:    value.ColorDepth8,
			expected: "8-color",
		},
		{
			name:     "256 colors",
			depth:    value.ColorDepth256,
			expected: "256-color",
		},
		{
			name:     "truecolor",
			depth:    value.ColorDepthTrueColor,
			expected: "truecolor",
		},
		{
			name:     "unknown",
			depth:    value.ColorDepth(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.depth.String()

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewCapabilities(t *testing.T) {
	tests := []struct {
		name           string
		ansi           bool
		colors         value.ColorDepth
		mouse          bool
		alt            bool
		cursor         bool
		expectedANSI   bool
		expectedColors value.ColorDepth
		expectedMouse  bool
		expectedAlt    bool
		expectedCursor bool
	}{
		{
			name:           "full capabilities",
			ansi:           true,
			colors:         value.ColorDepthTrueColor,
			mouse:          true,
			alt:            true,
			cursor:         true,
			expectedANSI:   true,
			expectedColors: value.ColorDepthTrueColor,
			expectedMouse:  true,
			expectedAlt:    true,
			expectedCursor: true,
		},
		{
			name:           "no ANSI - all features disabled",
			ansi:           false,
			colors:         value.ColorDepthTrueColor,
			mouse:          true,
			alt:            true,
			cursor:         true,
			expectedANSI:   false,
			expectedColors: value.ColorDepthTrueColor,
			expectedMouse:  false, // Forced false
			expectedAlt:    false, // Forced false
			expectedCursor: false, // Forced false
		},
		{
			name:           "ANSI with no mouse/alt/cursor",
			ansi:           true,
			colors:         value.ColorDepth256,
			mouse:          false,
			alt:            false,
			cursor:         false,
			expectedANSI:   true,
			expectedColors: value.ColorDepth256,
			expectedMouse:  false,
			expectedAlt:    false,
			expectedCursor: false,
		},
		{
			name:           "invalid color depth",
			ansi:           true,
			colors:         value.ColorDepth(999),
			mouse:          false,
			alt:            false,
			cursor:         false,
			expectedANSI:   true,
			expectedColors: value.ColorDepthNone, // Clamped to none
			expectedMouse:  false,
			expectedAlt:    false,
			expectedCursor: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caps := value.NewCapabilities(tt.ansi, tt.colors, tt.mouse, tt.alt, tt.cursor)

			if caps.SupportsANSI() != tt.expectedANSI {
				t.Errorf("expected ANSI %v, got %v", tt.expectedANSI, caps.SupportsANSI())
			}
			if caps.ColorDepth() != tt.expectedColors {
				t.Errorf("expected color depth %v, got %v", tt.expectedColors, caps.ColorDepth())
			}
			if caps.SupportsMouse() != tt.expectedMouse {
				t.Errorf("expected mouse %v, got %v", tt.expectedMouse, caps.SupportsMouse())
			}
			if caps.SupportsAltScreen() != tt.expectedAlt {
				t.Errorf("expected alt screen %v, got %v", tt.expectedAlt, caps.SupportsAltScreen())
			}
			if caps.SupportsCursorControl() != tt.expectedCursor {
				t.Errorf("expected cursor control %v, got %v", tt.expectedCursor, caps.SupportsCursorControl())
			}
		})
	}
}

func TestCapabilities_SupportsColor(t *testing.T) {
	tests := []struct {
		name     string
		caps     *value.Capabilities
		expected bool
	}{
		{
			name:     "no color",
			caps:     value.NewCapabilities(true, value.ColorDepthNone, false, false, false),
			expected: false,
		},
		{
			name:     "8 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth8, false, false, false),
			expected: true,
		},
		{
			name:     "256 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth256, false, false, false),
			expected: true,
		},
		{
			name:     "truecolor",
			caps:     value.NewCapabilities(true, value.ColorDepthTrueColor, false, false, false),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps.SupportsColor()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCapabilities_SupportsTrueColor(t *testing.T) {
	tests := []struct {
		name     string
		caps     *value.Capabilities
		expected bool
	}{
		{
			name:     "no color",
			caps:     value.NewCapabilities(true, value.ColorDepthNone, false, false, false),
			expected: false,
		},
		{
			name:     "8 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth8, false, false, false),
			expected: false,
		},
		{
			name:     "256 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth256, false, false, false),
			expected: false,
		},
		{
			name:     "truecolor",
			caps:     value.NewCapabilities(true, value.ColorDepthTrueColor, false, false, false),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps.SupportsTrueColor()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCapabilities_Supports256Color(t *testing.T) {
	tests := []struct {
		name     string
		caps     *value.Capabilities
		expected bool
	}{
		{
			name:     "no color",
			caps:     value.NewCapabilities(true, value.ColorDepthNone, false, false, false),
			expected: false,
		},
		{
			name:     "8 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth8, false, false, false),
			expected: false,
		},
		{
			name:     "256 colors",
			caps:     value.NewCapabilities(true, value.ColorDepth256, false, false, false),
			expected: true,
		},
		{
			name:     "truecolor (also 256+)",
			caps:     value.NewCapabilities(true, value.ColorDepthTrueColor, false, false, false),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps.Supports256Color()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCapabilities_IsDumbTerminal(t *testing.T) {
	tests := []struct {
		name     string
		caps     *value.Capabilities
		expected bool
	}{
		{
			name:     "dumb terminal",
			caps:     value.NewCapabilities(false, value.ColorDepthNone, false, false, false),
			expected: true,
		},
		{
			name:     "ANSI but no color",
			caps:     value.NewCapabilities(true, value.ColorDepthNone, false, false, false),
			expected: false,
		},
		{
			name:     "no ANSI but color",
			caps:     value.NewCapabilities(false, value.ColorDepth8, false, false, false),
			expected: false,
		},
		{
			name:     "full capabilities",
			caps:     value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps.IsDumbTerminal()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCapabilities_Equal(t *testing.T) {
	caps1 := value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	caps2 := value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)
	caps3 := value.NewCapabilities(true, value.ColorDepth256, true, true, true)

	tests := []struct {
		name     string
		caps1    *value.Capabilities
		caps2    *value.Capabilities
		expected bool
	}{
		{
			name:     "equal capabilities",
			caps1:    caps1,
			caps2:    caps2,
			expected: true,
		},
		{
			name:     "different color depth",
			caps1:    caps1,
			caps2:    caps3,
			expected: false,
		},
		{
			name:     "nil comparison",
			caps1:    caps1,
			caps2:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps1.Equal(tt.caps2)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCapabilities_Immutability(t *testing.T) {
	// Verify that Capabilities is immutable
	caps1 := value.NewCapabilities(true, value.ColorDepthTrueColor, true, true, true)

	// Store original values
	originalANSI := caps1.SupportsANSI()
	originalColors := caps1.ColorDepth()

	// Create another capability - original should be unchanged
	_ = value.NewCapabilities(false, value.ColorDepthNone, false, false, false)

	if caps1.SupportsANSI() != originalANSI {
		t.Error("ANSI support was mutated")
	}
	if caps1.ColorDepth() != originalColors {
		t.Error("color depth was mutated")
	}
}
