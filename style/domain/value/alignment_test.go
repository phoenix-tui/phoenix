package value

import "testing"

// --- HorizontalAlignment Tests ---

func TestHorizontalAlignmentString(t *testing.T) {
	tests := []struct {
		name      string
		alignment HorizontalAlignment
		want      string
	}{
		{name: "left", alignment: AlignLeft, want: "Left"},
		{name: "center", alignment: AlignCenter, want: "Center"},
		{name: "right", alignment: AlignRight, want: "Right"},
		{name: "unknown", alignment: HorizontalAlignment(99), want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- VerticalAlignment Tests ---

func TestVerticalAlignmentString(t *testing.T) {
	tests := []struct {
		name      string
		alignment VerticalAlignment
		want      string
	}{
		{name: "top", alignment: AlignTop, want: "Top"},
		{name: "middle", alignment: AlignMiddle, want: "Middle"},
		{name: "bottom", alignment: AlignBottom, want: "Bottom"},
		{name: "unknown", alignment: VerticalAlignment(99), want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- NewAlignment Tests ---

func TestNewAlignment(t *testing.T) {
	tests := []struct {
		name       string
		horizontal HorizontalAlignment
		vertical   VerticalAlignment
	}{
		{name: "left-top", horizontal: AlignLeft, vertical: AlignTop},
		{name: "center-middle", horizontal: AlignCenter, vertical: AlignMiddle},
		{name: "right-bottom", horizontal: AlignRight, vertical: AlignBottom},
		{name: "left-bottom", horizontal: AlignLeft, vertical: AlignBottom},
		{name: "right-top", horizontal: AlignRight, vertical: AlignTop},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAlignment(tt.horizontal, tt.vertical)

			if a.Horizontal() != tt.horizontal {
				t.Errorf("Horizontal() = %v, want %v", a.Horizontal(), tt.horizontal)
			}
			if a.Vertical() != tt.vertical {
				t.Errorf("Vertical() = %v, want %v", a.Vertical(), tt.vertical)
			}
		})
	}
}

// --- Convenience Constructor Tests ---

func TestLeftTop(t *testing.T) {
	a := LeftTop()
	if a.Horizontal() != AlignLeft {
		t.Errorf("Horizontal() = %v, want AlignLeft", a.Horizontal())
	}
	if a.Vertical() != AlignTop {
		t.Errorf("Vertical() = %v, want AlignTop", a.Vertical())
	}
}

func TestLeftMiddle(t *testing.T) {
	a := LeftMiddle()
	if a.Horizontal() != AlignLeft {
		t.Errorf("Horizontal() = %v, want AlignLeft", a.Horizontal())
	}
	if a.Vertical() != AlignMiddle {
		t.Errorf("Vertical() = %v, want AlignMiddle", a.Vertical())
	}
}

func TestLeftBottom(t *testing.T) {
	a := LeftBottom()
	if a.Horizontal() != AlignLeft {
		t.Errorf("Horizontal() = %v, want AlignLeft", a.Horizontal())
	}
	if a.Vertical() != AlignBottom {
		t.Errorf("Vertical() = %v, want AlignBottom", a.Vertical())
	}
}

func TestCenterTop(t *testing.T) {
	a := CenterTop()
	if a.Horizontal() != AlignCenter {
		t.Errorf("Horizontal() = %v, want AlignCenter", a.Horizontal())
	}
	if a.Vertical() != AlignTop {
		t.Errorf("Vertical() = %v, want AlignTop", a.Vertical())
	}
}

func TestCenterMiddle(t *testing.T) {
	a := CenterMiddle()
	if a.Horizontal() != AlignCenter {
		t.Errorf("Horizontal() = %v, want AlignCenter", a.Horizontal())
	}
	if a.Vertical() != AlignMiddle {
		t.Errorf("Vertical() = %v, want AlignMiddle", a.Vertical())
	}
}

func TestCenterBottom(t *testing.T) {
	a := CenterBottom()
	if a.Horizontal() != AlignCenter {
		t.Errorf("Horizontal() = %v, want AlignCenter", a.Horizontal())
	}
	if a.Vertical() != AlignBottom {
		t.Errorf("Vertical() = %v, want AlignBottom", a.Vertical())
	}
}

func TestRightTop(t *testing.T) {
	a := RightTop()
	if a.Horizontal() != AlignRight {
		t.Errorf("Horizontal() = %v, want AlignRight", a.Horizontal())
	}
	if a.Vertical() != AlignTop {
		t.Errorf("Vertical() = %v, want AlignTop", a.Vertical())
	}
}

func TestRightMiddle(t *testing.T) {
	a := RightMiddle()
	if a.Horizontal() != AlignRight {
		t.Errorf("Horizontal() = %v, want AlignRight", a.Horizontal())
	}
	if a.Vertical() != AlignMiddle {
		t.Errorf("Vertical() = %v, want AlignMiddle", a.Vertical())
	}
}

func TestRightBottom(t *testing.T) {
	a := RightBottom()
	if a.Horizontal() != AlignRight {
		t.Errorf("Horizontal() = %v, want AlignRight", a.Horizontal())
	}
	if a.Vertical() != AlignBottom {
		t.Errorf("Vertical() = %v, want AlignBottom", a.Vertical())
	}
}

// --- Getter Tests ---

func TestAlignmentGetters(t *testing.T) {
	a := NewAlignment(AlignCenter, AlignMiddle)

	if got := a.Horizontal(); got != AlignCenter {
		t.Errorf("Horizontal() = %v, want AlignCenter", got)
	}
	if got := a.Vertical(); got != AlignMiddle {
		t.Errorf("Vertical() = %v, want AlignMiddle", got)
	}
}

// --- Equality Tests ---

func TestAlignmentEqual(t *testing.T) {
	tests := []struct {
		name string
		a1   Alignment
		a2   Alignment
		want bool
	}{
		{
			name: "equal - left top",
			a1:   LeftTop(),
			a2:   LeftTop(),
			want: true,
		},
		{
			name: "equal - center middle",
			a1:   CenterMiddle(),
			a2:   CenterMiddle(),
			want: true,
		},
		{
			name: "equal - via constructor",
			a1:   NewAlignment(AlignRight, AlignBottom),
			a2:   RightBottom(),
			want: true,
		},
		{
			name: "not equal - different horizontal",
			a1:   LeftTop(),
			a2:   RightTop(),
			want: false,
		},
		{
			name: "not equal - different vertical",
			a1:   LeftTop(),
			a2:   LeftBottom(),
			want: false,
		},
		{
			name: "not equal - both different",
			a1:   LeftTop(),
			a2:   RightBottom(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a1.Equal(tt.a2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- String Tests ---

func TestAlignmentString(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      string
	}{
		{
			name:      "left top",
			alignment: LeftTop(),
			want:      "Alignment(Left, Top)",
		},
		{
			name:      "center middle",
			alignment: CenterMiddle(),
			want:      "Alignment(Center, Middle)",
		},
		{
			name:      "right bottom",
			alignment: RightBottom(),
			want:      "Alignment(Right, Bottom)",
		},
		{
			name:      "left middle",
			alignment: LeftMiddle(),
			want:      "Alignment(Left, Middle)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alignment.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Comprehensive Coverage Tests ---

func TestAllAlignmentCombinations(t *testing.T) {
	// Test all 9 combinations of horizontal and vertical alignment
	combinations := []struct {
		constructor func() Alignment
		horizontal  HorizontalAlignment
		vertical    VerticalAlignment
	}{
		{LeftTop, AlignLeft, AlignTop},
		{LeftMiddle, AlignLeft, AlignMiddle},
		{LeftBottom, AlignLeft, AlignBottom},
		{CenterTop, AlignCenter, AlignTop},
		{CenterMiddle, AlignCenter, AlignMiddle},
		{CenterBottom, AlignCenter, AlignBottom},
		{RightTop, AlignRight, AlignTop},
		{RightMiddle, AlignRight, AlignMiddle},
		{RightBottom, AlignRight, AlignBottom},
	}

	for _, combo := range combinations {
		a := combo.constructor()

		if a.Horizontal() != combo.horizontal {
			t.Errorf("%v: Horizontal() = %v, want %v",
				a.String(), a.Horizontal(), combo.horizontal)
		}
		if a.Vertical() != combo.vertical {
			t.Errorf("%v: Vertical() = %v, want %v",
				a.String(), a.Vertical(), combo.vertical)
		}

		// Verify consistency with NewAlignment
		a2 := NewAlignment(combo.horizontal, combo.vertical)
		if !a.Equal(a2) {
			t.Errorf("Constructor %v not equal to NewAlignment(%v, %v)",
				a.String(), combo.horizontal, combo.vertical)
		}
	}
}
