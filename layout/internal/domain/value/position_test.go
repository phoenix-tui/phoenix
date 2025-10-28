package value

import "testing"

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name  string
		x     int
		y     int
		wantX int
		wantY int
	}{
		{
			name: "positive coordinates",
			x:    10, y: 5,
			wantX: 10, wantY: 5,
		},
		{
			name: "zero coordinates",
			x:    0, y: 0,
			wantX: 0, wantY: 0,
		},
		{
			name: "negative x clamped to 0",
			x:    -5, y: 10,
			wantX: 0, wantY: 10,
		},
		{
			name: "negative y clamped to 0",
			x:    10, y: -5,
			wantX: 10, wantY: 0,
		},
		{
			name: "both negative clamped to 0",
			x:    -10, y: -5,
			wantX: 0, wantY: 0,
		},
		{
			name: "large coordinates",
			x:    1000, y: 2000,
			wantX: 1000, wantY: 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition(tt.x, tt.y)
			if pos.X() != tt.wantX {
				t.Errorf("X() = %d, want %d", pos.X(), tt.wantX)
			}
			if pos.Y() != tt.wantY {
				t.Errorf("Y() = %d, want %d", pos.Y(), tt.wantY)
			}
		})
	}
}

func TestOrigin(t *testing.T) {
	pos := Origin()
	if pos.X() != 0 || pos.Y() != 0 {
		t.Errorf("Origin() = (%d, %d), want (0, 0)", pos.X(), pos.Y())
	}
	if !pos.IsOrigin() {
		t.Error("Origin() should return true for IsOrigin()")
	}
}

func TestPosition_Add(t *testing.T) {
	tests := []struct {
		name  string
		base  Position
		dx    int
		dy    int
		wantX int
		wantY int
	}{
		{
			name: "add positive",
			base: NewPosition(10, 5),
			dx:   2, dy: 3,
			wantX: 12, wantY: 8,
		},
		{
			name: "add zero",
			base: NewPosition(10, 5),
			dx:   0, dy: 0,
			wantX: 10, wantY: 5,
		},
		{
			name: "add negative (within bounds)",
			base: NewPosition(10, 5),
			dx:   -2, dy: -3,
			wantX: 8, wantY: 2,
		},
		{
			name: "add negative (clamp to origin)",
			base: NewPosition(10, 5),
			dx:   -20, dy: -10,
			wantX: 0, wantY: 0,
		},
		{
			name: "from origin",
			base: Origin(),
			dx:   5, dy: 10,
			wantX: 5, wantY: 10,
		},
		{
			name: "immutability check",
			base: NewPosition(10, 5),
			dx:   2, dy: 3,
			wantX: 12, wantY: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.base
			result := tt.base.Add(tt.dx, tt.dy)

			if result.X() != tt.wantX {
				t.Errorf("Add() X = %d, want %d", result.X(), tt.wantX)
			}
			if result.Y() != tt.wantY {
				t.Errorf("Add() Y = %d, want %d", result.Y(), tt.wantY)
			}

			// Check immutability
			if tt.name == "immutability check" {
				if original.X() != 10 || original.Y() != 5 {
					t.Error("Original position was modified (not immutable)")
				}
			}
		})
	}
}

func TestPosition_Sub(t *testing.T) {
	tests := []struct {
		name  string
		base  Position
		dx    int
		dy    int
		wantX int
		wantY int
	}{
		{
			name: "subtract positive (within bounds)",
			base: NewPosition(10, 5),
			dx:   2, dy: 3,
			wantX: 8, wantY: 2,
		},
		{
			name: "subtract zero",
			base: NewPosition(10, 5),
			dx:   0, dy: 0,
			wantX: 10, wantY: 5,
		},
		{
			name: "subtract negative (adds)",
			base: NewPosition(10, 5),
			dx:   -2, dy: -3,
			wantX: 12, wantY: 8,
		},
		{
			name: "subtract beyond origin (clamp)",
			base: NewPosition(5, 3),
			dx:   10, dy: 10,
			wantX: 0, wantY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Sub(tt.dx, tt.dy)
			if result.X() != tt.wantX {
				t.Errorf("Sub() X = %d, want %d", result.X(), tt.wantX)
			}
			if result.Y() != tt.wantY {
				t.Errorf("Sub() Y = %d, want %d", result.Y(), tt.wantY)
			}
		})
	}
}

func TestPosition_Offset(t *testing.T) {
	tests := []struct {
		name   string
		base   Position
		offset Position
		wantX  int
		wantY  int
	}{
		{
			name:   "offset by positive",
			base:   NewPosition(10, 5),
			offset: NewPosition(2, 3),
			wantX:  12, wantY: 8,
		},
		{
			name:   "offset by zero",
			base:   NewPosition(10, 5),
			offset: Origin(),
			wantX:  10, wantY: 5,
		},
		{
			name:   "offset from origin",
			base:   Origin(),
			offset: NewPosition(5, 10),
			wantX:  5, wantY: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Offset(tt.offset)
			if result.X() != tt.wantX {
				t.Errorf("Offset() X = %d, want %d", result.X(), tt.wantX)
			}
			if result.Y() != tt.wantY {
				t.Errorf("Offset() Y = %d, want %d", result.Y(), tt.wantY)
			}
		})
	}
}

func TestPosition_Distance(t *testing.T) {
	tests := []struct {
		name string
		p1   Position
		p2   Position
		want int
	}{
		{
			name: "same position",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(10, 5),
			want: 0,
		},
		{
			name: "horizontal distance",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(15, 5),
			want: 5,
		},
		{
			name: "vertical distance",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(10, 10),
			want: 5,
		},
		{
			name: "diagonal distance",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(15, 8),
			want: 8, // |15-10| + |8-5| = 5 + 3 = 8
		},
		{
			name: "distance from origin",
			p1:   Origin(),
			p2:   NewPosition(3, 4),
			want: 7, // 3 + 4 = 7
		},
		{
			name: "negative direction (same result)",
			p1:   NewPosition(15, 8),
			p2:   NewPosition(10, 5),
			want: 8, // Same as diagonal distance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.Distance(tt.p2)
			if got != tt.want {
				t.Errorf("Distance() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPosition_Equals(t *testing.T) {
	tests := []struct {
		name string
		p1   Position
		p2   Position
		want bool
	}{
		{
			name: "equal positions",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(10, 5),
			want: true,
		},
		{
			name: "different x",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(11, 5),
			want: false,
		},
		{
			name: "different y",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(10, 6),
			want: false,
		},
		{
			name: "both different",
			p1:   NewPosition(10, 5),
			p2:   NewPosition(11, 6),
			want: false,
		},
		{
			name: "both origin",
			p1:   Origin(),
			p2:   Origin(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.Equals(tt.p2)
			if got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsOrigin(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		{
			name: "origin",
			pos:  Origin(),
			want: true,
		},
		{
			name: "explicit (0, 0)",
			pos:  NewPosition(0, 0),
			want: true,
		},
		{
			name: "non-zero x",
			pos:  NewPosition(1, 0),
			want: false,
		},
		{
			name: "non-zero y",
			pos:  NewPosition(0, 1),
			want: false,
		},
		{
			name: "both non-zero",
			pos:  NewPosition(10, 5),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.IsOrigin()
			if got != tt.want {
				t.Errorf("IsOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want string
	}{
		{
			name: "origin",
			pos:  Origin(),
			want: "Position{x=0, y=0}",
		},
		{
			name: "positive coordinates",
			pos:  NewPosition(10, 5),
			want: "Position{x=10, y=5}",
		},
		{
			name: "large coordinates",
			pos:  NewPosition(123, 456),
			want: "Position{x=123, y=456}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
