package value

import "testing"

// TestRoundedBorder tests the pre-defined RoundedBorder.
func TestRoundedBorder(t *testing.T) {
	b := RoundedBorder

	if b.Top != "─" {
		t.Errorf("RoundedBorder.Top = %q, want %q", b.Top, "─")
	}
	if b.Bottom != "─" {
		t.Errorf("RoundedBorder.Bottom = %q, want %q", b.Bottom, "─")
	}
	if b.Left != "│" {
		t.Errorf("RoundedBorder.Left = %q, want %q", b.Left, "│")
	}
	if b.Right != "│" {
		t.Errorf("RoundedBorder.Right = %q, want %q", b.Right, "│")
	}
	if b.TopLeft != "╭" {
		t.Errorf("RoundedBorder.TopLeft = %q, want %q", b.TopLeft, "╭")
	}
	if b.TopRight != "╮" {
		t.Errorf("RoundedBorder.TopRight = %q, want %q", b.TopRight, "╮")
	}
	if b.BottomLeft != "╰" {
		t.Errorf("RoundedBorder.BottomLeft = %q, want %q", b.BottomLeft, "╰")
	}
	if b.BottomRight != "╯" {
		t.Errorf("RoundedBorder.BottomRight = %q, want %q", b.BottomRight, "╯")
	}
}

// TestThickBorder tests the pre-defined ThickBorder.
func TestThickBorder(t *testing.T) {
	b := ThickBorder

	if b.Top != "━" {
		t.Errorf("ThickBorder.Top = %q, want %q", b.Top, "━")
	}
	if b.Bottom != "━" {
		t.Errorf("ThickBorder.Bottom = %q, want %q", b.Bottom, "━")
	}
	if b.Left != "┃" {
		t.Errorf("ThickBorder.Left = %q, want %q", b.Left, "┃")
	}
	if b.Right != "┃" {
		t.Errorf("ThickBorder.Right = %q, want %q", b.Right, "┃")
	}
	if b.TopLeft != "┏" {
		t.Errorf("ThickBorder.TopLeft = %q, want %q", b.TopLeft, "┏")
	}
	if b.TopRight != "┓" {
		t.Errorf("ThickBorder.TopRight = %q, want %q", b.TopRight, "┓")
	}
	if b.BottomLeft != "┗" {
		t.Errorf("ThickBorder.BottomLeft = %q, want %q", b.BottomLeft, "┗")
	}
	if b.BottomRight != "┛" {
		t.Errorf("ThickBorder.BottomRight = %q, want %q", b.BottomRight, "┛")
	}
}

// TestDoubleBorder tests the pre-defined DoubleBorder.
func TestDoubleBorder(t *testing.T) {
	b := DoubleBorder

	if b.Top != "═" {
		t.Errorf("DoubleBorder.Top = %q, want %q", b.Top, "═")
	}
	if b.Bottom != "═" {
		t.Errorf("DoubleBorder.Bottom = %q, want %q", b.Bottom, "═")
	}
	if b.Left != "║" {
		t.Errorf("DoubleBorder.Left = %q, want %q", b.Left, "║")
	}
	if b.Right != "║" {
		t.Errorf("DoubleBorder.Right = %q, want %q", b.Right, "║")
	}
	if b.TopLeft != "╔" {
		t.Errorf("DoubleBorder.TopLeft = %q, want %q", b.TopLeft, "╔")
	}
	if b.TopRight != "╗" {
		t.Errorf("DoubleBorder.TopRight = %q, want %q", b.TopRight, "╗")
	}
	if b.BottomLeft != "╚" {
		t.Errorf("DoubleBorder.BottomLeft = %q, want %q", b.BottomLeft, "╚")
	}
	if b.BottomRight != "╝" {
		t.Errorf("DoubleBorder.BottomRight = %q, want %q", b.BottomRight, "╝")
	}
}

// TestNormalBorder tests the pre-defined NormalBorder.
func TestNormalBorder(t *testing.T) {
	b := NormalBorder

	if b.Top != "─" {
		t.Errorf("NormalBorder.Top = %q, want %q", b.Top, "─")
	}
	if b.Bottom != "─" {
		t.Errorf("NormalBorder.Bottom = %q, want %q", b.Bottom, "─")
	}
	if b.Left != "│" {
		t.Errorf("NormalBorder.Left = %q, want %q", b.Left, "│")
	}
	if b.Right != "│" {
		t.Errorf("NormalBorder.Right = %q, want %q", b.Right, "│")
	}
	if b.TopLeft != "┌" {
		t.Errorf("NormalBorder.TopLeft = %q, want %q", b.TopLeft, "┌")
	}
	if b.TopRight != "┐" {
		t.Errorf("NormalBorder.TopRight = %q, want %q", b.TopRight, "┐")
	}
	if b.BottomLeft != "└" {
		t.Errorf("NormalBorder.BottomLeft = %q, want %q", b.BottomLeft, "└")
	}
	if b.BottomRight != "┘" {
		t.Errorf("NormalBorder.BottomRight = %q, want %q", b.BottomRight, "┘")
	}
}

// TestHiddenBorder tests the pre-defined HiddenBorder.
func TestHiddenBorder(t *testing.T) {
	b := HiddenBorder

	if b.Top != "" || b.Bottom != "" || b.Left != "" || b.Right != "" ||
		b.TopLeft != "" || b.TopRight != "" || b.BottomLeft != "" || b.BottomRight != "" {
		t.Error("HiddenBorder should have all empty strings")
	}

	if !b.IsHidden() {
		t.Error("HiddenBorder.IsHidden() should return true")
	}
}

// TestASCIIBorder tests the pre-defined ASCIIBorder.
func TestASCIIBorder(t *testing.T) {
	b := ASCIIBorder

	if b.Top != "-" {
		t.Errorf("ASCIIBorder.Top = %q, want %q", b.Top, "-")
	}
	if b.Bottom != "-" {
		t.Errorf("ASCIIBorder.Bottom = %q, want %q", b.Bottom, "-")
	}
	if b.Left != "|" {
		t.Errorf("ASCIIBorder.Left = %q, want %q", b.Left, "|")
	}
	if b.Right != "|" {
		t.Errorf("ASCIIBorder.Right = %q, want %q", b.Right, "|")
	}
	if b.TopLeft != "+" || b.TopRight != "+" || b.BottomLeft != "+" || b.BottomRight != "+" {
		t.Error("ASCIIBorder corners should all be '+'")
	}
}

// TestBorderEqual tests the Equal() method.
func TestBorderEqual(t *testing.T) {
	b1 := RoundedBorder
	b2 := RoundedBorder
	b3 := ThickBorder

	if !b1.Equal(b2) {
		t.Error("Same borders should be equal")
	}

	if b1.Equal(b3) {
		t.Error("Different borders should not be equal")
	}

	// Custom border test
	custom1 := Border{Top: "x", Bottom: "x", Left: "y", Right: "y"}
	custom2 := Border{Top: "x", Bottom: "x", Left: "y", Right: "y"}
	custom3 := Border{Top: "x", Bottom: "x", Left: "z", Right: "y"}

	if !custom1.Equal(custom2) {
		t.Error("Same custom borders should be equal")
	}

	if custom1.Equal(custom3) {
		t.Error("Different custom borders should not be equal")
	}
}

// TestBorderIsHidden tests the IsHidden() method.
func TestBorderIsHidden(t *testing.T) {
	tests := []struct {
		name   string
		border Border
		want   bool
	}{
		{"HiddenBorder", HiddenBorder, true},
		{"Empty border", Border{}, true},
		{"RoundedBorder", RoundedBorder, false},
		{"Partial", Border{Top: "x"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.border.IsHidden()
			if got != tt.want {
				t.Errorf("Border.IsHidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBorderString tests the String() method.
func TestBorderString(t *testing.T) {
	// Hidden border
	hidden := HiddenBorder
	if s := hidden.String(); s != "Border(hidden)" {
		t.Errorf("HiddenBorder.String() = %q, want %q", s, "Border(hidden)")
	}

	// Visible border
	rounded := RoundedBorder
	s := rounded.String()
	if s == "" {
		t.Error("RoundedBorder.String() should not be empty")
	}
	if !contains(s, "╭") || !contains(s, "╮") || !contains(s, "╰") || !contains(s, "╯") {
		t.Errorf("RoundedBorder.String() = %q, missing corner characters", s)
	}
}

// TestCustomBorder tests creating a custom border.
func TestCustomBorder(t *testing.T) {
	custom := Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	if custom.Top != "═" {
		t.Errorf("Custom border Top = %q, want %q", custom.Top, "═")
	}

	if custom.IsHidden() {
		t.Error("Custom border should not be hidden")
	}

	// Should equal DoubleBorder
	if !custom.Equal(DoubleBorder) {
		t.Error("Custom border should equal DoubleBorder")
	}
}

// TestBorderUnicodeCorrectness tests that Unicode characters are valid.
func TestBorderUnicodeCorrectness(t *testing.T) {
	borders := []Border{
		RoundedBorder,
		ThickBorder,
		DoubleBorder,
		NormalBorder,
	}

	for _, b := range borders {
		// Each character should be a valid UTF-8 string
		chars := []string{b.Top, b.Bottom, b.Left, b.Right,
			b.TopLeft, b.TopRight, b.BottomLeft, b.BottomRight}

		for _, ch := range chars {
			if ch == "" {
				t.Errorf("Border has empty character: %v", b)
			}
			// Each character should be exactly 1 rune (single Unicode codepoint)
			runes := []rune(ch)
			if len(runes) != 1 {
				t.Errorf("Border character should be 1 rune, got %d: %q", len(runes), ch)
			}
		}
	}
}
