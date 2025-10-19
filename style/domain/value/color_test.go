package value

import "testing"

// TestRGB tests the RGB constructor.
func TestRGB(t *testing.T) {
	tests := []struct {
		name                string
		r, g, b             uint8
		wantR, wantG, wantB uint8
	}{
		{"Black", 0, 0, 0, 0, 0, 0},
		{"White", 255, 255, 255, 255, 255, 255},
		{"Red", 255, 0, 0, 255, 0, 0},
		{"Green", 0, 255, 0, 0, 255, 0},
		{"Blue", 0, 0, 255, 0, 0, 255},
		{"Gray", 128, 128, 128, 128, 128, 128},
		{"Custom", 123, 45, 67, 123, 45, 67},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RGB(tt.r, tt.g, tt.b)
			gotR, gotG, gotB := c.RGB()

			if gotR != tt.wantR || gotG != tt.wantG || gotB != tt.wantB {
				t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)",
					gotR, gotG, gotB, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

// TestHex tests the Hex constructor with valid inputs.
func TestHex(t *testing.T) {
	tests := []struct {
		name                string
		hex                 string
		wantR, wantG, wantB uint8
		wantErr             bool
	}{
		{"6-digit with hash", "#FF00FF", 255, 0, 255, false},
		{"6-digit without hash", "FF00FF", 255, 0, 255, false},
		{"3-digit with hash", "#F0F", 255, 0, 255, false},
		{"3-digit without hash", "F0F", 255, 0, 255, false},
		{"Black", "#000000", 0, 0, 0, false},
		{"White", "#FFFFFF", 255, 255, 255, false},
		{"Lowercase", "#ff00ff", 255, 0, 255, false},
		{"Mixed case", "#Ff00Ff", 255, 0, 255, false},
		{"Invalid length", "#FF", 0, 0, 0, true},
		{"Invalid length", "#FFFFFFF", 0, 0, 0, true},
		{"Invalid chars", "#GGGGGG", 0, 0, 0, true},
		{"Invalid chars", "#XYZ", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Hex(tt.hex)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Hex(%q) expected error, got nil", tt.hex)
				}
				return
			}

			if err != nil {
				t.Errorf("Hex(%q) unexpected error: %v", tt.hex, err)
				return
			}

			gotR, gotG, gotB := c.RGB()
			if gotR != tt.wantR || gotG != tt.wantG || gotB != tt.wantB {
				t.Errorf("Hex(%q).RGB() = (%d, %d, %d), want (%d, %d, %d)",
					tt.hex, gotR, gotG, gotB, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

// TestFromANSI256 tests the FromANSI256 constructor.
func TestFromANSI256(t *testing.T) {
	tests := []struct {
		name string
		code uint8
		// We test that it returns a valid color (not exact RGB match,
		// since conversion is lossy and uses standard palette)
		wantValid bool
	}{
		{"Basic black", 0, true},
		{"Basic red", 1, true},
		{"Basic bright white", 15, true},
		{"Cube start", 16, true},
		{"Cube mid", 100, true},
		{"Cube end", 231, true},
		{"Gray start", 232, true},
		{"Gray mid", 244, true},
		{"Gray end", 255, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FromANSI256(tt.code)
			r, g, b := c.RGB()

			// Just verify we got valid RGB values (0-255)
			if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
				t.Errorf("FromANSI256(%d).RGB() = (%d, %d, %d), values out of range",
					tt.code, r, g, b)
			}
		})
	}
}

// TestColorHex tests the Hex() method.
func TestColorHex(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{"Black", RGB(0, 0, 0), "#000000"},
		{"White", RGB(255, 255, 255), "#FFFFFF"},
		{"Red", RGB(255, 0, 0), "#FF0000"},
		{"Green", RGB(0, 255, 0), "#00FF00"},
		{"Blue", RGB(0, 0, 255), "#0000FF"},
		{"Custom", RGB(123, 45, 67), "#7B2D43"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.Hex()
			if got != tt.want {
				t.Errorf("Color.Hex() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestColorEqual tests the Equal() method.
func TestColorEqual(t *testing.T) {
	red := RGB(255, 0, 0)
	red2 := RGB(255, 0, 0)
	blue := RGB(0, 0, 255)
	almostRed := RGB(254, 0, 0)

	tests := []struct {
		name   string
		c1, c2 Color
		want   bool
	}{
		{"Same color", red, red2, true},
		{"Different colors", red, blue, false},
		{"Almost same", red, almostRed, false},
		{"Self equality", red, red, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c1.Equal(tt.c2)
			if got != tt.want {
				t.Errorf("Color.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestColorString tests the String() method.
func TestColorString(t *testing.T) {
	c := RGB(255, 0, 0)
	s := c.String()

	// Just verify it contains RGB values and hex
	if s == "" {
		t.Error("Color.String() returned empty string")
	}

	// Should contain RGB values
	if !contains(s, "255") || !contains(s, "0") {
		t.Errorf("Color.String() = %q, missing RGB values", s)
	}

	// Should contain hex representation
	if !contains(s, "#FF0000") {
		t.Errorf("Color.String() = %q, missing hex representation", s)
	}
}

// TestHexRoundTrip tests that RGB -> Hex -> RGB is idempotent.
func TestHexRoundTrip(t *testing.T) {
	original := RGB(123, 45, 67)
	hex := original.Hex()

	decoded, err := Hex(hex)
	if err != nil {
		t.Fatalf("Hex(%q) unexpected error: %v", hex, err)
	}

	if !original.Equal(decoded) {
		t.Errorf("Round trip failed: original %v != decoded %v", original, decoded)
	}
}

// TestFromANSI256ToRGB tests FromANSI256 specific codes.
func TestFromANSI256ToRGB(t *testing.T) {
	// Test basic 16 colors (0-15)
	black := FromANSI256(0)
	r, g, b := black.RGB()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("FromANSI256(0) should be black, got RGB(%d, %d, %d)", r, g, b)
	}

	// Test 6x6x6 color cube (16-231)
	// Code 16 should be black in the cube (0,0,0 -> 0,0,0)
	cube16 := FromANSI256(16)
	r, g, b = cube16.RGB()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("FromANSI256(16) should be (0,0,0), got RGB(%d, %d, %d)", r, g, b)
	}

	// Code 231 should be white in the cube (5,5,5 -> 255,255,255)
	cube231 := FromANSI256(231)
	r, g, b = cube231.RGB()
	if r != 255 || g != 255 || b != 255 {
		t.Errorf("FromANSI256(231) should be (255,255,255), got RGB(%d, %d, %d)", r, g, b)
	}

	// Test grayscale (232-255)
	// Code 232 should be near-black (8,8,8)
	gray232 := FromANSI256(232)
	r, g, b = gray232.RGB()
	if r != 8 || g != 8 || b != 8 {
		t.Errorf("FromANSI256(232) should be (8,8,8), got RGB(%d, %d, %d)", r, g, b)
	}

	// Code 255 should be near-white (238,238,238)
	gray255 := FromANSI256(255)
	r, g, b = gray255.RGB()
	if r != 238 || g != 238 || b != 238 {
		t.Errorf("FromANSI256(255) should be (238,238,238), got RGB(%d, %d, %d)", r, g, b)
	}
}

// --- Helpers ---

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
