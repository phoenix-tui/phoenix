package value

import "testing"

func TestNewPosition(t *testing.T) {
	p := NewPosition(10, 20)
	if p.X() != 10 {
		t.Errorf("Position.X() = %d, want 10", p.X())
	}
	if p.Y() != 20 {
		t.Errorf("Position.Y() = %d, want 20", p.Y())
	}
}

func TestPosition_String(t *testing.T) {
	p := NewPosition(5, 15)
	expected := "(5,15)"
	if got := p.String(); got != expected {
		t.Errorf("Position.String() = %s, want %s", got, expected)
	}
}

func TestPosition_Equals(t *testing.T) {
	p1 := NewPosition(10, 20)
	p2 := NewPosition(10, 20)
	p3 := NewPosition(10, 21)

	if !p1.Equals(p2) {
		t.Errorf("Position.Equals() = false, want true for equal positions")
	}
	if p1.Equals(p3) {
		t.Errorf("Position.Equals() = true, want false for different positions")
	}
}

func TestPosition_DistanceTo(t *testing.T) {
	tests := []struct {
		name     string
		p1       Position
		p2       Position
		expected int
	}{
		{"same position", NewPosition(0, 0), NewPosition(0, 0), 0},
		{"horizontal", NewPosition(0, 0), NewPosition(5, 0), 5},
		{"vertical", NewPosition(0, 0), NewPosition(0, 5), 5},
		{"diagonal", NewPosition(0, 0), NewPosition(3, 4), 7}, // Manhattan distance
		{"negative", NewPosition(5, 5), NewPosition(2, 3), 5}, // |5-2| + |5-3| = 5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.DistanceTo(tt.p2); got != tt.expected {
				t.Errorf("Position.DistanceTo() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestPosition_IsWithinTolerance(t *testing.T) {
	p1 := NewPosition(10, 10)

	tests := []struct {
		name      string
		p2        Position
		tolerance int
		expected  bool
	}{
		{"same position", NewPosition(10, 10), 0, true},
		{"within tolerance", NewPosition(11, 10), 1, true},
		{"at tolerance", NewPosition(11, 11), 2, true},
		{"beyond tolerance", NewPosition(13, 10), 2, false},
		{"large tolerance", NewPosition(15, 15), 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p1.IsWithinTolerance(tt.p2, tt.tolerance); got != tt.expected {
				t.Errorf("Position.IsWithinTolerance() = %v, want %v (distance=%d, tolerance=%d)",
					got, tt.expected, p1.DistanceTo(tt.p2), tt.tolerance)
			}
		})
	}
}
