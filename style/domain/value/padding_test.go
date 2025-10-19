package value

import "testing"

// --- Constructor Tests ---

func TestNewPadding(t *testing.T) {
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
			p := NewPadding(tt.top, tt.right, tt.bottom, tt.left)

			if p.Top() != tt.wantTop {
				t.Errorf("Top() = %d, want %d", p.Top(), tt.wantTop)
			}
			if p.Right() != tt.wantRight {
				t.Errorf("Right() = %d, want %d", p.Right(), tt.wantRight)
			}
			if p.Bottom() != tt.wantBottom {
				t.Errorf("Bottom() = %d, want %d", p.Bottom(), tt.wantBottom)
			}
			if p.Left() != tt.wantLeft {
				t.Errorf("Left() = %d, want %d", p.Left(), tt.wantLeft)
			}
		})
	}
}

func TestUniformPadding(t *testing.T) {
	tests := []struct {
		name string
		all  int
		want int
	}{
		{name: "zero padding", all: 0, want: 0},
		{name: "positive padding", all: 5, want: 5},
		{name: "negative padding clamped", all: -3, want: 0},
		{name: "large padding", all: 100, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := UniformPadding(tt.all)

			if p.Top() != tt.want {
				t.Errorf("Top() = %d, want %d", p.Top(), tt.want)
			}
			if p.Right() != tt.want {
				t.Errorf("Right() = %d, want %d", p.Right(), tt.want)
			}
			if p.Bottom() != tt.want {
				t.Errorf("Bottom() = %d, want %d", p.Bottom(), tt.want)
			}
			if p.Left() != tt.want {
				t.Errorf("Left() = %d, want %d", p.Left(), tt.want)
			}
		})
	}
}

func TestVerticalHorizontal(t *testing.T) {
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
			p := VerticalHorizontal(tt.vertical, tt.horizontal)

			if p.Top() != tt.wantV {
				t.Errorf("Top() = %d, want %d", p.Top(), tt.wantV)
			}
			if p.Bottom() != tt.wantV {
				t.Errorf("Bottom() = %d, want %d", p.Bottom(), tt.wantV)
			}
			if p.Left() != tt.wantH {
				t.Errorf("Left() = %d, want %d", p.Left(), tt.wantH)
			}
			if p.Right() != tt.wantH {
				t.Errorf("Right() = %d, want %d", p.Right(), tt.wantH)
			}
		})
	}
}

// --- Getter Tests ---

func TestPaddingGetters(t *testing.T) {
	p := NewPadding(1, 2, 3, 4)

	if got := p.Top(); got != 1 {
		t.Errorf("Top() = %d, want 1", got)
	}
	if got := p.Right(); got != 2 {
		t.Errorf("Right() = %d, want 2", got)
	}
	if got := p.Bottom(); got != 3 {
		t.Errorf("Bottom() = %d, want 3", got)
	}
	if got := p.Left(); got != 4 {
		t.Errorf("Left() = %d, want 4", got)
	}
}

// --- Calculation Tests ---

func TestPaddingHorizontal(t *testing.T) {
	tests := []struct {
		name string
		p    Padding
		want int
	}{
		{name: "zero padding", p: UniformPadding(0), want: 0},
		{name: "uniform padding", p: UniformPadding(5), want: 10},
		{name: "asymmetric padding", p: NewPadding(1, 2, 3, 4), want: 6}, // left(4) + right(2)
		{name: "only horizontal", p: VerticalHorizontal(0, 5), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Horizontal(); got != tt.want {
				t.Errorf("Horizontal() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPaddingVertical(t *testing.T) {
	tests := []struct {
		name string
		p    Padding
		want int
	}{
		{name: "zero padding", p: UniformPadding(0), want: 0},
		{name: "uniform padding", p: UniformPadding(5), want: 10},
		{name: "asymmetric padding", p: NewPadding(1, 2, 3, 4), want: 4}, // top(1) + bottom(3)
		{name: "only vertical", p: VerticalHorizontal(5, 0), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Vertical(); got != tt.want {
				t.Errorf("Vertical() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPaddingTotal(t *testing.T) {
	tests := []struct {
		name     string
		p        Padding
		wantVert int
		wantHorz int
	}{
		{
			name:     "zero padding",
			p:        UniformPadding(0),
			wantVert: 0, wantHorz: 0,
		},
		{
			name:     "uniform padding",
			p:        UniformPadding(5),
			wantVert: 10, wantHorz: 10,
		},
		{
			name:     "asymmetric padding",
			p:        NewPadding(1, 2, 3, 4),
			wantVert: 4, wantHorz: 6, // top(1)+bottom(3)=4, left(4)+right(2)=6
		},
		{
			name:     "vertical horizontal constructor",
			p:        VerticalHorizontal(3, 7),
			wantVert: 6, wantHorz: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVert, gotHorz := tt.p.Total()

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

func TestPaddingEqual(t *testing.T) {
	tests := []struct {
		name string
		p1   Padding
		p2   Padding
		want bool
	}{
		{
			name: "equal uniform",
			p1:   UniformPadding(5),
			p2:   UniformPadding(5),
			want: true,
		},
		{
			name: "equal asymmetric",
			p1:   NewPadding(1, 2, 3, 4),
			p2:   NewPadding(1, 2, 3, 4),
			want: true,
		},
		{
			name: "not equal - different top",
			p1:   NewPadding(1, 2, 3, 4),
			p2:   NewPadding(5, 2, 3, 4),
			want: false,
		},
		{
			name: "not equal - different right",
			p1:   NewPadding(1, 2, 3, 4),
			p2:   NewPadding(1, 5, 3, 4),
			want: false,
		},
		{
			name: "not equal - different bottom",
			p1:   NewPadding(1, 2, 3, 4),
			p2:   NewPadding(1, 2, 5, 4),
			want: false,
		},
		{
			name: "not equal - different left",
			p1:   NewPadding(1, 2, 3, 4),
			p2:   NewPadding(1, 2, 3, 5),
			want: false,
		},
		{
			name: "zero vs non-zero",
			p1:   UniformPadding(0),
			p2:   UniformPadding(5),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.Equal(tt.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- String Tests ---

func TestPaddingString(t *testing.T) {
	tests := []struct {
		name string
		p    Padding
		want string
	}{
		{
			name: "zero padding",
			p:    UniformPadding(0),
			want: "Padding(top=0, right=0, bottom=0, left=0)",
		},
		{
			name: "uniform padding",
			p:    UniformPadding(5),
			want: "Padding(top=5, right=5, bottom=5, left=5)",
		},
		{
			name: "asymmetric padding",
			p:    NewPadding(1, 2, 3, 4),
			want: "Padding(top=1, right=2, bottom=3, left=4)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Edge Case Tests ---

func TestPaddingEdgeCases(t *testing.T) {
	t.Run("large values don't overflow", func(t *testing.T) {
		p := NewPadding(1000000, 2000000, 3000000, 4000000)
		vert, horz := p.Total()

		if vert != 4000000 {
			t.Errorf("Total() vertical = %d, want 4000000", vert)
		}
		if horz != 6000000 {
			t.Errorf("Total() horizontal = %d, want 6000000", horz)
		}
	})

	t.Run("all negative becomes all zero", func(t *testing.T) {
		p := NewPadding(-10, -20, -30, -40)
		vert, horz := p.Total()

		if vert != 0 {
			t.Errorf("Total() vertical = %d, want 0", vert)
		}
		if horz != 0 {
			t.Errorf("Total() horizontal = %d, want 0", horz)
		}
	})
}
