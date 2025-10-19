package value

import "testing"

// --- Constructor Tests ---

func TestNewMargin(t *testing.T) {
	tests := []struct {
		name                                     string
		top, right, bottom, left                 int
		wantTop, wantRight, wantBottom, wantLeft int
	}{
		{
			name: "all positive values",
			top:  1, right: 2, bottom: 3, left: 4,
			wantTop: 1, wantRight: 2, wantBottom: 3, wantLeft: 4,
		},
		{
			name: "all zeros",
			top:  0, right: 0, bottom: 0, left: 0,
			wantTop: 0, wantRight: 0, wantBottom: 0, wantLeft: 0,
		},
		{
			name: "negative values clamped to zero",
			top:  -1, right: -5, bottom: -10, left: -3,
			wantTop: 0, wantRight: 0, wantBottom: 0, wantLeft: 0,
		},
		{
			name: "mixed positive and negative",
			top:  5, right: -2, bottom: 3, left: -1,
			wantTop: 5, wantRight: 0, wantBottom: 3, wantLeft: 0,
		},
		{
			name: "large values",
			top:  100, right: 200, bottom: 150, left: 75,
			wantTop: 100, wantRight: 200, wantBottom: 150, wantLeft: 75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMargin(tt.top, tt.right, tt.bottom, tt.left)

			if m.Top() != tt.wantTop {
				t.Errorf("Top() = %d, want %d", m.Top(), tt.wantTop)
			}
			if m.Right() != tt.wantRight {
				t.Errorf("Right() = %d, want %d", m.Right(), tt.wantRight)
			}
			if m.Bottom() != tt.wantBottom {
				t.Errorf("Bottom() = %d, want %d", m.Bottom(), tt.wantBottom)
			}
			if m.Left() != tt.wantLeft {
				t.Errorf("Left() = %d, want %d", m.Left(), tt.wantLeft)
			}
		})
	}
}

func TestUniformMargin(t *testing.T) {
	tests := []struct {
		name string
		all  int
		want int
	}{
		{name: "zero margin", all: 0, want: 0},
		{name: "positive margin", all: 5, want: 5},
		{name: "negative margin clamped", all: -3, want: 0},
		{name: "large margin", all: 100, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := UniformMargin(tt.all)

			if m.Top() != tt.want {
				t.Errorf("Top() = %d, want %d", m.Top(), tt.want)
			}
			if m.Right() != tt.want {
				t.Errorf("Right() = %d, want %d", m.Right(), tt.want)
			}
			if m.Bottom() != tt.want {
				t.Errorf("Bottom() = %d, want %d", m.Bottom(), tt.want)
			}
			if m.Left() != tt.want {
				t.Errorf("Left() = %d, want %d", m.Left(), tt.want)
			}
		})
	}
}

func TestVerticalHorizontalMargin(t *testing.T) {
	tests := []struct {
		name       string
		vertical   int
		horizontal int
		wantV      int
		wantH      int
	}{
		{
			name:     "both positive",
			vertical: 2, horizontal: 4,
			wantV: 2, wantH: 4,
		},
		{
			name:     "both zero",
			vertical: 0, horizontal: 0,
			wantV: 0, wantH: 0,
		},
		{
			name:     "negative vertical clamped",
			vertical: -3, horizontal: 5,
			wantV: 0, wantH: 5,
		},
		{
			name:     "negative horizontal clamped",
			vertical: 5, horizontal: -2,
			wantV: 5, wantH: 0,
		},
		{
			name:     "both negative clamped",
			vertical: -1, horizontal: -1,
			wantV: 0, wantH: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := VerticalHorizontalMargin(tt.vertical, tt.horizontal)

			if m.Top() != tt.wantV {
				t.Errorf("Top() = %d, want %d", m.Top(), tt.wantV)
			}
			if m.Bottom() != tt.wantV {
				t.Errorf("Bottom() = %d, want %d", m.Bottom(), tt.wantV)
			}
			if m.Left() != tt.wantH {
				t.Errorf("Left() = %d, want %d", m.Left(), tt.wantH)
			}
			if m.Right() != tt.wantH {
				t.Errorf("Right() = %d, want %d", m.Right(), tt.wantH)
			}
		})
	}
}

// --- Getter Tests ---

func TestMarginGetters(t *testing.T) {
	m := NewMargin(1, 2, 3, 4)

	if got := m.Top(); got != 1 {
		t.Errorf("Top() = %d, want 1", got)
	}
	if got := m.Right(); got != 2 {
		t.Errorf("Right() = %d, want 2", got)
	}
	if got := m.Bottom(); got != 3 {
		t.Errorf("Bottom() = %d, want 3", got)
	}
	if got := m.Left(); got != 4 {
		t.Errorf("Left() = %d, want 4", got)
	}
}

// --- Calculation Tests ---

func TestMarginHorizontal(t *testing.T) {
	tests := []struct {
		name string
		m    Margin
		want int
	}{
		{name: "zero margin", m: UniformMargin(0), want: 0},
		{name: "uniform margin", m: UniformMargin(5), want: 10},
		{name: "asymmetric margin", m: NewMargin(1, 2, 3, 4), want: 6}, // left(4) + right(2)
		{name: "only horizontal", m: VerticalHorizontalMargin(0, 5), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Horizontal(); got != tt.want {
				t.Errorf("Horizontal() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMarginVertical(t *testing.T) {
	tests := []struct {
		name string
		m    Margin
		want int
	}{
		{name: "zero margin", m: UniformMargin(0), want: 0},
		{name: "uniform margin", m: UniformMargin(5), want: 10},
		{name: "asymmetric margin", m: NewMargin(1, 2, 3, 4), want: 4}, // top(1) + bottom(3)
		{name: "only vertical", m: VerticalHorizontalMargin(5, 0), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Vertical(); got != tt.want {
				t.Errorf("Vertical() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMarginTotal(t *testing.T) {
	tests := []struct {
		name     string
		m        Margin
		wantVert int
		wantHorz int
	}{
		{
			name:     "zero margin",
			m:        UniformMargin(0),
			wantVert: 0, wantHorz: 0,
		},
		{
			name:     "uniform margin",
			m:        UniformMargin(5),
			wantVert: 10, wantHorz: 10,
		},
		{
			name:     "asymmetric margin",
			m:        NewMargin(1, 2, 3, 4),
			wantVert: 4, wantHorz: 6, // top(1)+bottom(3)=4, left(4)+right(2)=6
		},
		{
			name:     "vertical horizontal constructor",
			m:        VerticalHorizontalMargin(3, 7),
			wantVert: 6, wantHorz: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVert, gotHorz := tt.m.Total()

			if gotVert != tt.wantVert {
				t.Errorf("Total() vertical = %d, want %d", gotVert, tt.wantVert)
			}
			if gotHorz != tt.wantHorz {
				t.Errorf("Total() horizontal = %d, want %d", gotHorz, tt.wantHorz)
			}
		})
	}
}

// --- Equality Tests ---

func TestMarginEqual(t *testing.T) {
	tests := []struct {
		name string
		m1   Margin
		m2   Margin
		want bool
	}{
		{
			name: "equal uniform",
			m1:   UniformMargin(5),
			m2:   UniformMargin(5),
			want: true,
		},
		{
			name: "equal asymmetric",
			m1:   NewMargin(1, 2, 3, 4),
			m2:   NewMargin(1, 2, 3, 4),
			want: true,
		},
		{
			name: "not equal - different top",
			m1:   NewMargin(1, 2, 3, 4),
			m2:   NewMargin(5, 2, 3, 4),
			want: false,
		},
		{
			name: "not equal - different right",
			m1:   NewMargin(1, 2, 3, 4),
			m2:   NewMargin(1, 5, 3, 4),
			want: false,
		},
		{
			name: "not equal - different bottom",
			m1:   NewMargin(1, 2, 3, 4),
			m2:   NewMargin(1, 2, 5, 4),
			want: false,
		},
		{
			name: "not equal - different left",
			m1:   NewMargin(1, 2, 3, 4),
			m2:   NewMargin(1, 2, 3, 5),
			want: false,
		},
		{
			name: "zero vs non-zero",
			m1:   UniformMargin(0),
			m2:   UniformMargin(5),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Equal(tt.m2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- String Tests ---

func TestMarginString(t *testing.T) {
	tests := []struct {
		name string
		m    Margin
		want string
	}{
		{
			name: "zero margin",
			m:    UniformMargin(0),
			want: "Margin(top=0, right=0, bottom=0, left=0)",
		},
		{
			name: "uniform margin",
			m:    UniformMargin(5),
			want: "Margin(top=5, right=5, bottom=5, left=5)",
		},
		{
			name: "asymmetric margin",
			m:    NewMargin(1, 2, 3, 4),
			want: "Margin(top=1, right=2, bottom=3, left=4)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Edge Case Tests ---

func TestMarginEdgeCases(t *testing.T) {
	t.Run("large values don't overflow", func(t *testing.T) {
		m := NewMargin(1000000, 2000000, 3000000, 4000000)
		vert, horz := m.Total()

		if vert != 4000000 {
			t.Errorf("Total() vertical = %d, want 4000000", vert)
		}
		if horz != 6000000 {
			t.Errorf("Total() horizontal = %d, want 6000000", horz)
		}
	})

	t.Run("all negative becomes all zero", func(t *testing.T) {
		m := NewMargin(-10, -20, -30, -40)
		vert, horz := m.Total()

		if vert != 0 {
			t.Errorf("Total() vertical = %d, want 0", vert)
		}
		if horz != 0 {
			t.Errorf("Total() horizontal = %d, want 0", horz)
		}
	})
}

// --- Type Safety Tests ---

func TestMarginPaddingTypeSafety(t *testing.T) {
	// This test verifies that Margin and Padding are distinct types
	// and cannot be accidentally used interchangeably.

	m := NewMargin(1, 2, 3, 4)
	p := NewPadding(1, 2, 3, 4)

	// These should compile (same type)
	_ = m.Equal(NewMargin(1, 2, 3, 4))
	_ = p.Equal(NewPadding(1, 2, 3, 4))

	// These should NOT compile (different types):
	// _ = m.Equal(p) // compile error: cannot use p (Padding) as Margin
	// _ = p.Equal(m) // compile error: cannot use m (Margin) as Padding

	// This test passes if it compiles, demonstrating type safety
	if m.Top() != p.Top() {
		t.Errorf("Values should be equal even though types differ")
	}
}
