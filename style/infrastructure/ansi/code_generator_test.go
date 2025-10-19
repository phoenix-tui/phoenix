package ansi

import "testing"

// TestNewANSICodeGenerator tests the constructor.
func TestNewANSICodeGenerator(t *testing.T) {
	gen := NewANSICodeGenerator()
	if gen == nil {
		t.Error("NewANSICodeGenerator() returned nil")
	}
}

// TestForeground tests TrueColor foreground code generation.
func TestForeground(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name    string
		r, g, b uint8
		want    string
	}{
		{"Red", 255, 0, 0, "\x1b[38;2;255;0;0m"},
		{"Green", 0, 255, 0, "\x1b[38;2;0;255;0m"},
		{"Blue", 0, 0, 255, "\x1b[38;2;0;0;255m"},
		{"White", 255, 255, 255, "\x1b[38;2;255;255;255m"},
		{"Black", 0, 0, 0, "\x1b[38;2;0;0;0m"},
		{"Custom", 123, 45, 67, "\x1b[38;2;123;45;67m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Foreground(tt.r, tt.g, tt.b)
			if got != tt.want {
				t.Errorf("Foreground(%d, %d, %d) = %q, want %q",
					tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

// TestBackground tests TrueColor background code generation.
func TestBackground(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name    string
		r, g, b uint8
		want    string
	}{
		{"Red", 255, 0, 0, "\x1b[48;2;255;0;0m"},
		{"Green", 0, 255, 0, "\x1b[48;2;0;255;0m"},
		{"Blue", 0, 0, 255, "\x1b[48;2;0;0;255m"},
		{"White", 255, 255, 255, "\x1b[48;2;255;255;255m"},
		{"Black", 0, 0, 0, "\x1b[48;2;0;0;0m"},
		{"Custom", 123, 45, 67, "\x1b[48;2;123;45;67m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Background(tt.r, tt.g, tt.b)
			if got != tt.want {
				t.Errorf("Background(%d, %d, %d) = %q, want %q",
					tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

// TestForeground256 tests 256-color foreground code generation.
func TestForeground256(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name string
		code uint8
		want string
	}{
		{"Code 0", 0, "\x1b[38;5;0m"},
		{"Code 1", 1, "\x1b[38;5;1m"},
		{"Code 16", 16, "\x1b[38;5;16m"},
		{"Code 128", 128, "\x1b[38;5;128m"},
		{"Code 255", 255, "\x1b[38;5;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Foreground256(tt.code)
			if got != tt.want {
				t.Errorf("Foreground256(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// TestBackground256 tests 256-color background code generation.
func TestBackground256(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name string
		code uint8
		want string
	}{
		{"Code 0", 0, "\x1b[48;5;0m"},
		{"Code 1", 1, "\x1b[48;5;1m"},
		{"Code 16", 16, "\x1b[48;5;16m"},
		{"Code 128", 128, "\x1b[48;5;128m"},
		{"Code 255", 255, "\x1b[48;5;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Background256(tt.code)
			if got != tt.want {
				t.Errorf("Background256(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// TestForeground16 tests 16-color foreground code generation.
func TestForeground16(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name string
		code uint8
		want string
	}{
		{"Black", 0, "\x1b[30m"},
		{"Red", 1, "\x1b[31m"},
		{"Green", 2, "\x1b[32m"},
		{"Yellow", 3, "\x1b[33m"},
		{"Blue", 4, "\x1b[34m"},
		{"Magenta", 5, "\x1b[35m"},
		{"Cyan", 6, "\x1b[36m"},
		{"White", 7, "\x1b[37m"},
		{"Bright Black", 8, "\x1b[90m"},
		{"Bright Red", 9, "\x1b[91m"},
		{"Bright Green", 10, "\x1b[92m"},
		{"Bright Yellow", 11, "\x1b[93m"},
		{"Bright Blue", 12, "\x1b[94m"},
		{"Bright Magenta", 13, "\x1b[95m"},
		{"Bright Cyan", 14, "\x1b[96m"},
		{"Bright White", 15, "\x1b[97m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Foreground16(tt.code)
			if got != tt.want {
				t.Errorf("Foreground16(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// TestBackground16 tests 16-color background code generation.
func TestBackground16(t *testing.T) {
	gen := NewANSICodeGenerator()

	tests := []struct {
		name string
		code uint8
		want string
	}{
		{"Black", 0, "\x1b[40m"},
		{"Red", 1, "\x1b[41m"},
		{"Green", 2, "\x1b[42m"},
		{"Yellow", 3, "\x1b[43m"},
		{"Blue", 4, "\x1b[44m"},
		{"Magenta", 5, "\x1b[45m"},
		{"Cyan", 6, "\x1b[46m"},
		{"White", 7, "\x1b[47m"},
		{"Bright Black", 8, "\x1b[100m"},
		{"Bright Red", 9, "\x1b[101m"},
		{"Bright Green", 10, "\x1b[102m"},
		{"Bright Yellow", 11, "\x1b[103m"},
		{"Bright Blue", 12, "\x1b[104m"},
		{"Bright Magenta", 13, "\x1b[105m"},
		{"Bright Cyan", 14, "\x1b[106m"},
		{"Bright White", 15, "\x1b[107m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Background16(tt.code)
			if got != tt.want {
				t.Errorf("Background16(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

// TestReset tests reset code generation.
func TestReset(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[0m"
	got := gen.Reset()

	if got != want {
		t.Errorf("Reset() = %q, want %q", got, want)
	}
}

// TestBold tests bold code generation.
func TestBold(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[1m"
	got := gen.Bold()

	if got != want {
		t.Errorf("Bold() = %q, want %q", got, want)
	}
}

// TestItalic tests italic code generation.
func TestItalic(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[3m"
	got := gen.Italic()

	if got != want {
		t.Errorf("Italic() = %q, want %q", got, want)
	}
}

// TestUnderline tests underline code generation.
func TestUnderline(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[4m"
	got := gen.Underline()

	if got != want {
		t.Errorf("Underline() = %q, want %q", got, want)
	}
}

// TestStrikethrough tests strikethrough code generation.
func TestStrikethrough(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[9m"
	got := gen.Strikethrough()

	if got != want {
		t.Errorf("Strikethrough() = %q, want %q", got, want)
	}
}

// TestBoldOff tests bold-off code generation.
func TestBoldOff(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[22m"
	got := gen.BoldOff()

	if got != want {
		t.Errorf("BoldOff() = %q, want %q", got, want)
	}
}

// TestItalicOff tests italic-off code generation.
func TestItalicOff(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[23m"
	got := gen.ItalicOff()

	if got != want {
		t.Errorf("ItalicOff() = %q, want %q", got, want)
	}
}

// TestUnderlineOff tests underline-off code generation.
func TestUnderlineOff(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[24m"
	got := gen.UnderlineOff()

	if got != want {
		t.Errorf("UnderlineOff() = %q, want %q", got, want)
	}
}

// TestStrikethroughOff tests strikethrough-off code generation.
func TestStrikethroughOff(t *testing.T) {
	gen := NewANSICodeGenerator()
	want := "\x1b[29m"
	got := gen.StrikethroughOff()

	if got != want {
		t.Errorf("StrikethroughOff() = %q, want %q", got, want)
	}
}

// TestAllCodesAreValid tests that all generated codes are valid ANSI sequences.
func TestAllCodesAreValid(t *testing.T) {
	gen := NewANSICodeGenerator()

	// All codes should start with ESC (0x1b) and end with 'm'
	codes := []string{
		gen.Foreground(255, 0, 0),
		gen.Background(255, 0, 0),
		gen.Foreground256(128),
		gen.Background256(128),
		gen.Foreground16(9),
		gen.Background16(9),
		gen.Reset(),
		gen.Bold(),
		gen.Italic(),
		gen.Underline(),
		gen.Strikethrough(),
		gen.BoldOff(),
		gen.ItalicOff(),
		gen.UnderlineOff(),
		gen.StrikethroughOff(),
	}

	for i, code := range codes {
		if len(code) == 0 {
			t.Errorf("Code %d is empty", i)
			continue
		}

		if code[0] != '\x1b' {
			t.Errorf("Code %d (%q) doesn't start with ESC", i, code)
		}

		if code[len(code)-1] != 'm' {
			t.Errorf("Code %d (%q) doesn't end with 'm'", i, code)
		}
	}
}
