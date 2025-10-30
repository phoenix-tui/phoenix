package value

import (
	"strings"
	"testing"
)

func TestNewTextStyles(t *testing.T) {
	styles := NewTextStyles()

	if styles.Bold {
		t.Error("Expected Bold to be false")
	}
	if styles.Italic {
		t.Error("Expected Italic to be false")
	}
	if styles.Underline {
		t.Error("Expected Underline to be false")
	}
	if styles.Color != "" {
		t.Errorf("Expected Color to be empty, got %s", styles.Color)
	}
	if !styles.IsPlain() {
		t.Error("Expected IsPlain to be true for default styles")
	}
}

func TestNewTextStylesWithColor(t *testing.T) {
	tests := []struct {
		name      string
		color     string
		wantError bool
	}{
		{"Valid color", "#FF0000", false},
		{"Valid color lowercase", "#ff0000", false},
		{"Empty color", "", false},
		{"Invalid format - no hash", "FF0000", true},
		{"Invalid format - short", "#FFF", true},
		{"Invalid format - long", "#FF00000", true},
		{"Invalid format - non-hex", "#GGGGGG", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			styles, err := NewTextStylesWithColor(tt.color)
			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.color != "" && styles.Color != tt.color {
					t.Errorf("Expected color %s, got %s", tt.color, styles.Color)
				}
			}
		})
	}
}

func TestTextStylesWithBold(t *testing.T) {
	styles := NewTextStyles()

	// Enable bold
	styled := styles.WithBold(true)
	if !styled.Bold {
		t.Error("Expected Bold to be true")
	}

	// Original should be unchanged
	if styles.Bold {
		t.Error("Original TextStyles was mutated")
	}

	// Disable bold
	plain := styled.WithBold(false)
	if plain.Bold {
		t.Error("Expected Bold to be false")
	}
}

func TestTextStylesWithItalic(t *testing.T) {
	styles := NewTextStyles()

	styled := styles.WithItalic(true)
	if !styled.Italic {
		t.Error("Expected Italic to be true")
	}

	if styles.Italic {
		t.Error("Original TextStyles was mutated")
	}
}

func TestTextStylesWithUnderline(t *testing.T) {
	styles := NewTextStyles()

	styled := styles.WithUnderline(true)
	if !styled.Underline {
		t.Error("Expected Underline to be true")
	}

	if styles.Underline {
		t.Error("Original TextStyles was mutated")
	}
}

func TestTextStylesWithColor(t *testing.T) {
	styles := NewTextStyles()

	// Valid color
	styled, err := styles.WithColor("#FF0000")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if styled.Color != "#FF0000" {
		t.Errorf("Expected color #FF0000, got %s", styled.Color)
	}

	// Original should be unchanged
	if styles.Color != "" {
		t.Error("Original TextStyles was mutated")
	}

	// Invalid color
	_, err = styles.WithColor("invalid")
	if err == nil {
		t.Error("Expected error for invalid color")
	}

	// Clear color
	cleared, err := styled.WithColor("")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cleared.Color != "" {
		t.Error("Expected color to be cleared")
	}
}

func TestTextStylesIsPlain(t *testing.T) {
	tests := []struct {
		name      string
		styles    TextStyles
		wantPlain bool
	}{
		{
			name:      "Plain styles",
			styles:    NewTextStyles(),
			wantPlain: true,
		},
		{
			name:      "Bold only",
			styles:    NewTextStyles().WithBold(true),
			wantPlain: false,
		},
		{
			name:      "Italic only",
			styles:    NewTextStyles().WithItalic(true),
			wantPlain: false,
		},
		{
			name:      "Underline only",
			styles:    NewTextStyles().WithUnderline(true),
			wantPlain: false,
		},
		{
			name: "Color only",
			styles: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#FF0000")
				return s
			}(),
			wantPlain: false,
		},
		{
			name: "All styles",
			styles: func() TextStyles {
				s := NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
				s, _ = s.WithColor("#FF0000")
				return s
			}(),
			wantPlain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.styles.IsPlain() != tt.wantPlain {
				t.Errorf("IsPlain() = %v, want %v", tt.styles.IsPlain(), tt.wantPlain)
			}
		})
	}
}

func TestTextStylesEquals(t *testing.T) {
	tests := []struct {
		name   string
		s1     TextStyles
		s2     TextStyles
		equals bool
	}{
		{
			name:   "Both plain",
			s1:     NewTextStyles(),
			s2:     NewTextStyles(),
			equals: true,
		},
		{
			name:   "Same bold",
			s1:     NewTextStyles().WithBold(true),
			s2:     NewTextStyles().WithBold(true),
			equals: true,
		},
		{
			name:   "Different bold",
			s1:     NewTextStyles().WithBold(true),
			s2:     NewTextStyles(),
			equals: false,
		},
		{
			name: "Same color (different case)",
			s1: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#FF0000")
				return s
			}(),
			s2: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#ff0000")
				return s
			}(),
			equals: true,
		},
		{
			name: "Different color",
			s1: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#FF0000")
				return s
			}(),
			s2: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#00FF00")
				return s
			}(),
			equals: false,
		},
		{
			name: "All same",
			s1: func() TextStyles {
				s := NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
				s, _ = s.WithColor("#FF0000")
				return s
			}(),
			s2: func() TextStyles {
				s := NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
				s, _ = s.WithColor("#ff0000")
				return s
			}(),
			equals: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.s1.Equals(tt.s2) != tt.equals {
				t.Errorf("Equals() = %v, want %v", tt.s1.Equals(tt.s2), tt.equals)
			}
		})
	}
}

func TestTextStylesString(t *testing.T) {
	tests := []struct {
		name   string
		styles TextStyles
		want   string
	}{
		{
			name:   "Plain",
			styles: NewTextStyles(),
			want:   "plain",
		},
		{
			name:   "Bold only",
			styles: NewTextStyles().WithBold(true),
			want:   "bold",
		},
		{
			name:   "Italic only",
			styles: NewTextStyles().WithItalic(true),
			want:   "italic",
		},
		{
			name:   "Underline only",
			styles: NewTextStyles().WithUnderline(true),
			want:   "underline",
		},
		{
			name: "Color only",
			styles: func() TextStyles {
				s, _ := NewTextStyles().WithColor("#FF0000")
				return s
			}(),
			want: "color:#FF0000",
		},
		{
			name:   "Bold and italic",
			styles: NewTextStyles().WithBold(true).WithItalic(true),
			want:   "bold,italic",
		},
		{
			name: "All styles",
			styles: func() TextStyles {
				s := NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
				s, _ = s.WithColor("#FF0000")
				return s
			}(),
			want: "bold,italic,underline,color:#FF0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.styles.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseColorRGB(t *testing.T) {
	tests := []struct {
		name  string
		color string
		wantR int
		wantG int
		wantB int
		err   bool
	}{
		{"Red", "#FF0000", 255, 0, 0, false},
		{"Green", "#00FF00", 0, 255, 0, false},
		{"Blue", "#0000FF", 0, 0, 255, false},
		{"White", "#FFFFFF", 255, 255, 255, false},
		{"Black", "#000000", 0, 0, 0, false},
		{"Gray", "#808080", 128, 128, 128, false},
		{"Custom", "#123456", 18, 52, 86, false},
		{"Lowercase", "#abcdef", 171, 205, 239, false},
		{"Invalid format", "FF0000", 0, 0, 0, true},
		{"Invalid hex", "#GGGGGG", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, err := ParseColorRGB(tt.color)
			if tt.err {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if r != tt.wantR || g != tt.wantG || b != tt.wantB {
					t.Errorf("ParseColorRGB() = (%d, %d, %d), want (%d, %d, %d)",
						r, g, b, tt.wantR, tt.wantG, tt.wantB)
				}
			}
		})
	}
}

func TestFormatColorRGB(t *testing.T) {
	tests := []struct {
		name string
		r    int
		g    int
		b    int
		want string
		err  bool
	}{
		{"Red", 255, 0, 0, "#FF0000", false},
		{"Green", 0, 255, 0, "#00FF00", false},
		{"Blue", 0, 0, 255, "#0000FF", false},
		{"White", 255, 255, 255, "#FFFFFF", false},
		{"Black", 0, 0, 0, "#000000", false},
		{"Gray", 128, 128, 128, "#808080", false},
		{"Custom", 18, 52, 86, "#123456", false},
		{"Invalid R negative", -1, 0, 0, "", true},
		{"Invalid R overflow", 256, 0, 0, "", true},
		{"Invalid G negative", 0, -1, 0, "", true},
		{"Invalid G overflow", 0, 256, 0, "", true},
		{"Invalid B negative", 0, 0, -1, "", true},
		{"Invalid B overflow", 0, 0, 256, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatColorRGB(tt.r, tt.g, tt.b)
			if tt.err {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("FormatColorRGB() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestColorRoundTrip(t *testing.T) {
	colors := []string{
		"#FF0000", "#00FF00", "#0000FF",
		"#FFFFFF", "#000000", "#808080",
		"#123456", "#ABCDEF",
	}

	for _, color := range colors {
		t.Run(color, func(t *testing.T) {
			// Parse to RGB
			r, g, b, err := ParseColorRGB(color)
			if err != nil {
				t.Fatalf("Failed to parse color: %v", err)
			}

			// Format back to hex
			result, err := FormatColorRGB(r, g, b)
			if err != nil {
				t.Fatalf("Failed to format color: %v", err)
			}

			// Should match (case-insensitive)
			if !strings.EqualFold(result, color) {
				t.Errorf("Round trip failed: %s -> (%d,%d,%d) -> %s",
					color, r, g, b, result)
			}
		})
	}
}
