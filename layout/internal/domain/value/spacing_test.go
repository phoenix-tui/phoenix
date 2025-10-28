package value

import "testing"

func TestNewSpacing(t *testing.T) {
	tests := []struct {
		name       string
		top        int
		right      int
		bottom     int
		left       int
		wantTop    int
		wantRight  int
		wantBottom int
		wantLeft   int
	}{
		{
			name: "all positive",
			top:  1, right: 2, bottom: 3, left: 4,
			wantTop: 1, wantRight: 2, wantBottom: 3, wantLeft: 4,
		},
		{
			name: "all zero",
			top:  0, right: 0, bottom: 0, left: 0,
			wantTop: 0, wantRight: 0, wantBottom: 0, wantLeft: 0,
		},
		{
			name: "negative values clamped to 0",
			top:  -1, right: -2, bottom: -3, left: -4,
			wantTop: 0, wantRight: 0, wantBottom: 0, wantLeft: 0,
		},
		{
			name: "mixed positive and negative",
			top:  1, right: -2, bottom: 3, left: -4,
			wantTop: 1, wantRight: 0, wantBottom: 3, wantLeft: 0,
		},
		{
			name: "large values",
			top:  100, right: 200, bottom: 300, left: 400,
			wantTop: 100, wantRight: 200, wantBottom: 300, wantLeft: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpacing(tt.top, tt.right, tt.bottom, tt.left)

			if s.Top() != tt.wantTop {
				t.Errorf("Top() = %d, want %d", s.Top(), tt.wantTop)
			}
			if s.Right() != tt.wantRight {
				t.Errorf("Right() = %d, want %d", s.Right(), tt.wantRight)
			}
			if s.Bottom() != tt.wantBottom {
				t.Errorf("Bottom() = %d, want %d", s.Bottom(), tt.wantBottom)
			}
			if s.Left() != tt.wantLeft {
				t.Errorf("Left() = %d, want %d", s.Left(), tt.wantLeft)
			}
		})
	}
}

func TestNewSpacingAll(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "positive value",
			value: 5,
			want:  5,
		},
		{
			name:  "zero",
			value: 0,
			want:  0,
		},
		{
			name:  "negative clamped to 0",
			value: -10,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpacingAll(tt.value)

			if s.Top() != tt.want || s.Right() != tt.want || s.Bottom() != tt.want || s.Left() != tt.want {
				t.Errorf("NewSpacingAll(%d) = {%d %d %d %d}, want all %d",
					tt.value, s.Top(), s.Right(), s.Bottom(), s.Left(), tt.want)
			}
			if !s.IsUniform() && tt.want != 0 {
				t.Error("NewSpacingAll() should create uniform spacing")
			}
		})
	}
}

func TestNewSpacingVH(t *testing.T) {
	tests := []struct {
		name           string
		vertical       int
		horizontal     int
		wantVertical   int
		wantHorizontal int
	}{
		{
			name:     "positive values",
			vertical: 2, horizontal: 4,
			wantVertical: 2, wantHorizontal: 4,
		},
		{
			name:     "zero",
			vertical: 0, horizontal: 0,
			wantVertical: 0, wantHorizontal: 0,
		},
		{
			name:     "negative clamped to 0",
			vertical: -1, horizontal: -2,
			wantVertical: 0, wantHorizontal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpacingVH(tt.vertical, tt.horizontal)

			if s.Top() != tt.wantVertical {
				t.Errorf("Top() = %d, want %d", s.Top(), tt.wantVertical)
			}
			if s.Bottom() != tt.wantVertical {
				t.Errorf("Bottom() = %d, want %d", s.Bottom(), tt.wantVertical)
			}
			if s.Left() != tt.wantHorizontal {
				t.Errorf("Left() = %d, want %d", s.Left(), tt.wantHorizontal)
			}
			if s.Right() != tt.wantHorizontal {
				t.Errorf("Right() = %d, want %d", s.Right(), tt.wantHorizontal)
			}
		})
	}
}

func TestNewSpacingZero(t *testing.T) {
	s := NewSpacingZero()

	if !s.IsZero() {
		t.Error("NewSpacingZero() should create zero spacing")
	}
	if s.Top() != 0 || s.Right() != 0 || s.Bottom() != 0 || s.Left() != 0 {
		t.Errorf("NewSpacingZero() = {%d %d %d %d}, want all 0",
			s.Top(), s.Right(), s.Bottom(), s.Left())
	}
}

func TestSpacing_Horizontal(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		want    int
	}{
		{
			name:    "left=2, right=3",
			spacing: NewSpacing(0, 3, 0, 2),
			want:    5,
		},
		{
			name:    "zero",
			spacing: NewSpacingZero(),
			want:    0,
		},
		{
			name:    "symmetric",
			spacing: NewSpacingVH(1, 5),
			want:    10, // left=5, right=5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.Horizontal()
			if got != tt.want {
				t.Errorf("Horizontal() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpacing_Vertical(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		want    int
	}{
		{
			name:    "top=2, bottom=3",
			spacing: NewSpacing(2, 0, 3, 0),
			want:    5,
		},
		{
			name:    "zero",
			spacing: NewSpacingZero(),
			want:    0,
		},
		{
			name:    "symmetric",
			spacing: NewSpacingVH(5, 1),
			want:    10, // top=5, bottom=5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.Vertical()
			if got != tt.want {
				t.Errorf("Vertical() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpacing_IsZero(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		want    bool
	}{
		{
			name:    "all zero",
			spacing: NewSpacingZero(),
			want:    true,
		},
		{
			name:    "explicit zero",
			spacing: NewSpacing(0, 0, 0, 0),
			want:    true,
		},
		{
			name:    "one non-zero",
			spacing: NewSpacing(1, 0, 0, 0),
			want:    false,
		},
		{
			name:    "all non-zero",
			spacing: NewSpacingAll(1),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpacing_IsUniform(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		want    bool
	}{
		{
			name:    "all same",
			spacing: NewSpacingAll(5),
			want:    true,
		},
		{
			name:    "all zero (uniform)",
			spacing: NewSpacingZero(),
			want:    true,
		},
		{
			name:    "top different",
			spacing: NewSpacing(1, 2, 2, 2),
			want:    false,
		},
		{
			name:    "vertical/horizontal pattern",
			spacing: NewSpacingVH(1, 2),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.IsUniform()
			if got != tt.want {
				t.Errorf("IsUniform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpacing_With(t *testing.T) {
	base := NewSpacing(1, 2, 3, 4)

	t.Run("WithTop", func(t *testing.T) {
		s := base.WithTop(10)
		if s.Top() != 10 || s.Right() != 2 || s.Bottom() != 3 || s.Left() != 4 {
			t.Error("WithTop() did not update correctly")
		}
		if base.Top() != 1 {
			t.Error("Original spacing was modified (not immutable)")
		}
	})

	t.Run("WithRight", func(t *testing.T) {
		s := base.WithRight(20)
		if s.Top() != 1 || s.Right() != 20 || s.Bottom() != 3 || s.Left() != 4 {
			t.Error("WithRight() did not update correctly")
		}
	})

	t.Run("WithBottom", func(t *testing.T) {
		s := base.WithBottom(30)
		if s.Top() != 1 || s.Right() != 2 || s.Bottom() != 30 || s.Left() != 4 {
			t.Error("WithBottom() did not update correctly")
		}
	})

	t.Run("WithLeft", func(t *testing.T) {
		s := base.WithLeft(40)
		if s.Top() != 1 || s.Right() != 2 || s.Bottom() != 3 || s.Left() != 40 {
			t.Error("WithLeft() did not update correctly")
		}
	})

	t.Run("chaining", func(t *testing.T) {
		s := base.WithTop(10).WithRight(20)
		if s.Top() != 10 || s.Right() != 20 {
			t.Error("Chaining did not work correctly")
		}
	})
}

func TestSpacing_Add(t *testing.T) {
	tests := []struct {
		name string
		s1   Spacing
		s2   Spacing
		want Spacing
	}{
		{
			name: "add positive",
			s1:   NewSpacing(1, 2, 3, 4),
			s2:   NewSpacing(5, 6, 7, 8),
			want: NewSpacing(6, 8, 10, 12),
		},
		{
			name: "add zero",
			s1:   NewSpacing(1, 2, 3, 4),
			s2:   NewSpacingZero(),
			want: NewSpacing(1, 2, 3, 4),
		},
		{
			name: "add uniform",
			s1:   NewSpacingAll(2),
			s2:   NewSpacingAll(3),
			want: NewSpacingAll(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s1.Add(tt.s2)
			if !got.Equals(tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpacing_Scale(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		factor  int
		want    Spacing
	}{
		{
			name:    "scale by 2",
			spacing: NewSpacing(1, 2, 3, 4),
			factor:  2,
			want:    NewSpacing(2, 4, 6, 8),
		},
		{
			name:    "scale by 0",
			spacing: NewSpacing(1, 2, 3, 4),
			factor:  0,
			want:    NewSpacingZero(),
		},
		{
			name:    "scale by negative (clamped to 0)",
			spacing: NewSpacing(1, 2, 3, 4),
			factor:  -5,
			want:    NewSpacingZero(),
		},
		{
			name:    "scale uniform",
			spacing: NewSpacingAll(5),
			factor:  3,
			want:    NewSpacingAll(15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.Scale(tt.factor)
			if !got.Equals(tt.want) {
				t.Errorf("Scale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpacing_Equals(t *testing.T) {
	tests := []struct {
		name string
		s1   Spacing
		s2   Spacing
		want bool
	}{
		{
			name: "equal",
			s1:   NewSpacing(1, 2, 3, 4),
			s2:   NewSpacing(1, 2, 3, 4),
			want: true,
		},
		{
			name: "different top",
			s1:   NewSpacing(1, 2, 3, 4),
			s2:   NewSpacing(10, 2, 3, 4),
			want: false,
		},
		{
			name: "all different",
			s1:   NewSpacing(1, 2, 3, 4),
			s2:   NewSpacing(5, 6, 7, 8),
			want: false,
		},
		{
			name: "both zero",
			s1:   NewSpacingZero(),
			s2:   NewSpacingZero(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s1.Equals(tt.s2)
			if got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpacing_String(t *testing.T) {
	tests := []struct {
		name    string
		spacing Spacing
		want    string
	}{
		{
			name:    "zero",
			spacing: NewSpacingZero(),
			want:    "Spacing{0}",
		},
		{
			name:    "uniform",
			spacing: NewSpacingAll(5),
			want:    "Spacing{5}",
		},
		{
			name:    "vertical/horizontal",
			spacing: NewSpacingVH(1, 2),
			want:    "Spacing{1 2}",
		},
		{
			name:    "all different",
			spacing: NewSpacing(1, 2, 3, 4),
			want:    "Spacing{1 2 3 4}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spacing.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
