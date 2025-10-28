package value

import "testing"

// TestTerminalCapabilityString tests the String() method.
func TestTerminalCapabilityString(t *testing.T) {
	tests := []struct {
		name string
		tc   TerminalCapability
		want string
	}{
		{"NoColor", NoColor, "NoColor"},
		{"ANSI16", ANSI16, "ANSI16"},
		{"ANSI256", ANSI256, "ANSI256"},
		{"TrueColor", TrueColor, "TrueColor"},
		{"Invalid", TerminalCapability(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc.String()
			if got != tt.want {
				t.Errorf("TerminalCapability.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestSupportsColor tests the SupportsColor() method.
func TestSupportsColor(t *testing.T) {
	tests := []struct {
		name string
		tc   TerminalCapability
		want bool
	}{
		{"NoColor", NoColor, false},
		{"ANSI16", ANSI16, true},
		{"ANSI256", ANSI256, true},
		{"TrueColor", TrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc.SupportsColor()
			if got != tt.want {
				t.Errorf("TerminalCapability.SupportsColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSupportsTrueColor tests the SupportsTrueColor() method.
func TestSupportsTrueColor(t *testing.T) {
	tests := []struct {
		name string
		tc   TerminalCapability
		want bool
	}{
		{"NoColor", NoColor, false},
		{"ANSI16", ANSI16, false},
		{"ANSI256", ANSI256, false},
		{"TrueColor", TrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc.SupportsTrueColor()
			if got != tt.want {
				t.Errorf("TerminalCapability.SupportsTrueColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSupports256Color tests the Supports256Color() method.
func TestSupports256Color(t *testing.T) {
	tests := []struct {
		name string
		tc   TerminalCapability
		want bool
	}{
		{"NoColor", NoColor, false},
		{"ANSI16", ANSI16, false},
		{"ANSI256", ANSI256, true},
		{"TrueColor", TrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc.Supports256Color()
			if got != tt.want {
				t.Errorf("TerminalCapability.Supports256Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSupports16Color tests the Supports16Color() method.
func TestSupports16Color(t *testing.T) {
	tests := []struct {
		name string
		tc   TerminalCapability
		want bool
	}{
		{"NoColor", NoColor, false},
		{"ANSI16", ANSI16, true},
		{"ANSI256", ANSI256, true},
		{"TrueColor", TrueColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc.Supports16Color()
			if got != tt.want {
				t.Errorf("TerminalCapability.Supports16Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTerminalCapabilityConstants tests that constants have expected values.
func TestTerminalCapabilityConstants(t *testing.T) {
	if NoColor != 0 {
		t.Errorf("NoColor = %d, want 0", NoColor)
	}
	if ANSI16 != 1 {
		t.Errorf("ANSI16 = %d, want 1", ANSI16)
	}
	if ANSI256 != 2 {
		t.Errorf("ANSI256 = %d, want 2", ANSI256)
	}
	if TrueColor != 3 {
		t.Errorf("TrueColor = %d, want 3", TrueColor)
	}
}
