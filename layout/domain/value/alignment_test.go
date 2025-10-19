package value

import "testing"

func TestHorizontalAlignment_String(t *testing.T) {
	tests := []struct {
		name  string
		align HorizontalAlignment
		want  string
	}{
		{
			name:  "left",
			align: AlignLeft,
			want:  "Left",
		},
		{
			name:  "center",
			align: AlignCenter,
			want:  "Center",
		},
		{
			name:  "right",
			align: AlignRight,
			want:  "Right",
		},
		{
			name:  "unknown",
			align: HorizontalAlignment(999),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVerticalAlignment_String(t *testing.T) {
	tests := []struct {
		name  string
		align VerticalAlignment
		want  string
	}{
		{
			name:  "top",
			align: AlignTop,
			want:  "Top",
		},
		{
			name:  "middle",
			align: AlignMiddle,
			want:  "Middle",
		},
		{
			name:  "bottom",
			align: AlignBottom,
			want:  "Bottom",
		},
		{
			name:  "unknown",
			align: VerticalAlignment(999),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewAlignment(t *testing.T) {
	tests := []struct {
		name       string
		horizontal HorizontalAlignment
		vertical   VerticalAlignment
	}{
		{
			name:       "top-left",
			horizontal: AlignLeft,
			vertical:   AlignTop,
		},
		{
			name:       "center",
			horizontal: AlignCenter,
			vertical:   AlignMiddle,
		},
		{
			name:       "bottom-right",
			horizontal: AlignRight,
			vertical:   AlignBottom,
		},
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

func TestNewAlignmentDefault(t *testing.T) {
	a := NewAlignmentDefault()
	if a.Horizontal() != AlignLeft {
		t.Errorf("Default horizontal = %v, want %v", a.Horizontal(), AlignLeft)
	}
	if a.Vertical() != AlignTop {
		t.Errorf("Default vertical = %v, want %v", a.Vertical(), AlignTop)
	}
	if !a.IsDefault() {
		t.Error("NewAlignmentDefault() should return true for IsDefault()")
	}
}

func TestNewAlignmentCenter(t *testing.T) {
	a := NewAlignmentCenter()
	if a.Horizontal() != AlignCenter {
		t.Errorf("Center horizontal = %v, want %v", a.Horizontal(), AlignCenter)
	}
	if a.Vertical() != AlignMiddle {
		t.Errorf("Center vertical = %v, want %v", a.Vertical(), AlignMiddle)
	}
	if !a.IsCenter() {
		t.Error("NewAlignmentCenter() should return true for IsCenter()")
	}
}

func TestAlignment_Is(t *testing.T) {
	tests := []struct {
		name            string
		alignment       Alignment
		checkHorizontal HorizontalAlignment
		wantHorizontal  bool
		checkVertical   VerticalAlignment
		wantVertical    bool
	}{
		{
			name:            "top-left check left",
			alignment:       NewAlignment(AlignLeft, AlignTop),
			checkHorizontal: AlignLeft,
			wantHorizontal:  true,
			checkVertical:   AlignTop,
			wantVertical:    true,
		},
		{
			name:            "center check center",
			alignment:       NewAlignmentCenter(),
			checkHorizontal: AlignCenter,
			wantHorizontal:  true,
			checkVertical:   AlignMiddle,
			wantVertical:    true,
		},
		{
			name:            "mismatch",
			alignment:       NewAlignment(AlignLeft, AlignTop),
			checkHorizontal: AlignRight,
			wantHorizontal:  false,
			checkVertical:   AlignBottom,
			wantVertical:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotH := tt.alignment.IsHorizontal(tt.checkHorizontal)
			if gotH != tt.wantHorizontal {
				t.Errorf("IsHorizontal(%v) = %v, want %v", tt.checkHorizontal, gotH, tt.wantHorizontal)
			}
			gotV := tt.alignment.IsVertical(tt.checkVertical)
			if gotV != tt.wantVertical {
				t.Errorf("IsVertical(%v) = %v, want %v", tt.checkVertical, gotV, tt.wantVertical)
			}
		})
	}
}

func TestAlignment_IsCenter(t *testing.T) {
	tests := []struct {
		name  string
		align Alignment
		want  bool
	}{
		{
			name:  "center",
			align: NewAlignmentCenter(),
			want:  true,
		},
		{
			name:  "horizontal center only",
			align: NewAlignment(AlignCenter, AlignTop),
			want:  false,
		},
		{
			name:  "vertical center only",
			align: NewAlignment(AlignLeft, AlignMiddle),
			want:  false,
		},
		{
			name:  "default (top-left)",
			align: NewAlignmentDefault(),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.IsCenter()
			if got != tt.want {
				t.Errorf("IsCenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_IsDefault(t *testing.T) {
	tests := []struct {
		name  string
		align Alignment
		want  bool
	}{
		{
			name:  "default",
			align: NewAlignmentDefault(),
			want:  true,
		},
		{
			name:  "explicit top-left",
			align: NewAlignment(AlignLeft, AlignTop),
			want:  true,
		},
		{
			name:  "center",
			align: NewAlignmentCenter(),
			want:  false,
		},
		{
			name:  "bottom-right",
			align: NewAlignment(AlignRight, AlignBottom),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.IsDefault()
			if got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_With(t *testing.T) {
	base := NewAlignment(AlignLeft, AlignTop)

	t.Run("WithHorizontal", func(t *testing.T) {
		a := base.WithHorizontal(AlignCenter)
		if a.Horizontal() != AlignCenter {
			t.Errorf("Horizontal() = %v, want %v", a.Horizontal(), AlignCenter)
		}
		if a.Vertical() != AlignTop {
			t.Errorf("Vertical() = %v, want %v (should be unchanged)", a.Vertical(), AlignTop)
		}
		if base.Horizontal() != AlignLeft {
			t.Error("Original alignment was modified (not immutable)")
		}
	})

	t.Run("WithVertical", func(t *testing.T) {
		a := base.WithVertical(AlignMiddle)
		if a.Horizontal() != AlignLeft {
			t.Errorf("Horizontal() = %v, want %v (should be unchanged)", a.Horizontal(), AlignLeft)
		}
		if a.Vertical() != AlignMiddle {
			t.Errorf("Vertical() = %v, want %v", a.Vertical(), AlignMiddle)
		}
	})

	t.Run("chaining", func(t *testing.T) {
		a := base.WithHorizontal(AlignCenter).WithVertical(AlignMiddle)
		if !a.IsCenter() {
			t.Error("Chaining did not create centered alignment")
		}
	})
}

func TestAlignment_Equals(t *testing.T) {
	tests := []struct {
		name string
		a1   Alignment
		a2   Alignment
		want bool
	}{
		{
			name: "equal",
			a1:   NewAlignment(AlignCenter, AlignMiddle),
			a2:   NewAlignment(AlignCenter, AlignMiddle),
			want: true,
		},
		{
			name: "different horizontal",
			a1:   NewAlignment(AlignLeft, AlignMiddle),
			a2:   NewAlignment(AlignCenter, AlignMiddle),
			want: false,
		},
		{
			name: "different vertical",
			a1:   NewAlignment(AlignCenter, AlignTop),
			a2:   NewAlignment(AlignCenter, AlignMiddle),
			want: false,
		},
		{
			name: "both different",
			a1:   NewAlignment(AlignLeft, AlignTop),
			a2:   NewAlignment(AlignRight, AlignBottom),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a1.Equals(tt.a2)
			if got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignment_String(t *testing.T) {
	tests := []struct {
		name  string
		align Alignment
		want  string
	}{
		{
			name:  "default",
			align: NewAlignmentDefault(),
			want:  "Alignment{Left, Top}",
		},
		{
			name:  "center",
			align: NewAlignmentCenter(),
			want:  "Alignment{Center, Middle}",
		},
		{
			name:  "bottom-right",
			align: NewAlignment(AlignRight, AlignBottom),
			want:  "Alignment{Right, Bottom}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.align.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalculateHorizontalOffset(t *testing.T) {
	tests := []struct {
		name           string
		align          HorizontalAlignment
		contentWidth   int
		containerWidth int
		want           int
	}{
		{
			name:         "left alignment",
			align:        AlignLeft,
			contentWidth: 20, containerWidth: 80,
			want: 0,
		},
		{
			name:         "center alignment",
			align:        AlignCenter,
			contentWidth: 20, containerWidth: 80,
			want: 30, // (80 - 20) / 2 = 30
		},
		{
			name:         "right alignment",
			align:        AlignRight,
			contentWidth: 20, containerWidth: 80,
			want: 60, // 80 - 20 = 60
		},
		{
			name:         "content fills container",
			align:        AlignCenter,
			contentWidth: 80, containerWidth: 80,
			want: 0,
		},
		{
			name:         "content exceeds container",
			align:        AlignCenter,
			contentWidth: 100, containerWidth: 80,
			want: 0,
		},
		{
			name:         "center with odd difference",
			align:        AlignCenter,
			contentWidth: 21, containerWidth: 80,
			want: 29, // (80 - 21) / 2 = 29.5 -> 29 (integer division)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateHorizontalOffset(tt.align, tt.contentWidth, tt.containerWidth)
			if got != tt.want {
				t.Errorf("CalculateHorizontalOffset() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalculateVerticalOffset(t *testing.T) {
	tests := []struct {
		name            string
		align           VerticalAlignment
		contentHeight   int
		containerHeight int
		want            int
	}{
		{
			name:          "top alignment",
			align:         AlignTop,
			contentHeight: 10, containerHeight: 24,
			want: 0,
		},
		{
			name:          "middle alignment",
			align:         AlignMiddle,
			contentHeight: 10, containerHeight: 24,
			want: 7, // (24 - 10) / 2 = 7
		},
		{
			name:          "bottom alignment",
			align:         AlignBottom,
			contentHeight: 10, containerHeight: 24,
			want: 14, // 24 - 10 = 14
		},
		{
			name:          "content fills container",
			align:         AlignMiddle,
			contentHeight: 24, containerHeight: 24,
			want: 0,
		},
		{
			name:          "content exceeds container",
			align:         AlignMiddle,
			contentHeight: 30, containerHeight: 24,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateVerticalOffset(tt.align, tt.contentHeight, tt.containerHeight)
			if got != tt.want {
				t.Errorf("CalculateVerticalOffset() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAlignment_CalculateOffsets(t *testing.T) {
	tests := []struct {
		name            string
		align           Alignment
		contentWidth    int
		contentHeight   int
		containerWidth  int
		containerHeight int
		wantX           int
		wantY           int
	}{
		{
			name:         "center both axes",
			align:        NewAlignmentCenter(),
			contentWidth: 20, contentHeight: 10,
			containerWidth: 80, containerHeight: 24,
			wantX: 30, wantY: 7,
		},
		{
			name:         "top-left (default)",
			align:        NewAlignmentDefault(),
			contentWidth: 20, contentHeight: 10,
			containerWidth: 80, containerHeight: 24,
			wantX: 0, wantY: 0,
		},
		{
			name:         "bottom-right",
			align:        NewAlignment(AlignRight, AlignBottom),
			contentWidth: 20, contentHeight: 10,
			containerWidth: 80, containerHeight: 24,
			wantX: 60, wantY: 14,
		},
		{
			name:         "center horizontal, top vertical",
			align:        NewAlignment(AlignCenter, AlignTop),
			contentWidth: 20, contentHeight: 10,
			containerWidth: 80, containerHeight: 24,
			wantX: 30, wantY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := tt.align.CalculateOffsets(
				tt.contentWidth, tt.contentHeight,
				tt.containerWidth, tt.containerHeight,
			)
			if gotX != tt.wantX {
				t.Errorf("CalculateOffsets() X = %d, want %d", gotX, tt.wantX)
			}
			if gotY != tt.wantY {
				t.Errorf("CalculateOffsets() Y = %d, want %d", gotY, tt.wantY)
			}
		})
	}
}
