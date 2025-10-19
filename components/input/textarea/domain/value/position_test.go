package value

import (
	"testing"
)

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name    string
		row     int
		col     int
		wantRow int
		wantCol int
	}{
		{
			name:    "zero position",
			row:     0,
			col:     0,
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "positive position",
			row:     5,
			col:     10,
			wantRow: 5,
			wantCol: 10,
		},
		{
			name:    "large position",
			row:     1000,
			col:     2000,
			wantRow: 1000,
			wantCol: 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition(tt.row, tt.col)

			if pos.Row() != tt.wantRow {
				t.Errorf("Row() = %d, want %d", pos.Row(), tt.wantRow)
			}
			if pos.Col() != tt.wantCol {
				t.Errorf("Col() = %d, want %d", pos.Col(), tt.wantCol)
			}
		})
	}
}

func TestPosition_Row(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
		want int
	}{
		{name: "zero", row: 0, col: 5, want: 0},
		{name: "positive", row: 42, col: 10, want: 42},
		{name: "large", row: 9999, col: 0, want: 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition(tt.row, tt.col)
			if pos.Row() != tt.want {
				t.Errorf("Row() = %d, want %d", pos.Row(), tt.want)
			}
		})
	}
}

func TestPosition_Col(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
		want int
	}{
		{name: "zero", row: 5, col: 0, want: 0},
		{name: "positive", row: 10, col: 42, want: 42},
		{name: "large", row: 0, col: 9999, want: 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition(tt.row, tt.col)
			if pos.Col() != tt.want {
				t.Errorf("Col() = %d, want %d", pos.Col(), tt.want)
			}
		})
	}
}

func TestPosition_IsBefore(t *testing.T) {
	tests := []struct {
		name  string
		pos   Position
		other Position
		want  bool
	}{
		{
			name:  "earlier row",
			pos:   NewPosition(1, 5),
			other: NewPosition(2, 3),
			want:  true,
		},
		{
			name:  "same row, earlier col",
			pos:   NewPosition(1, 3),
			other: NewPosition(1, 5),
			want:  true,
		},
		{
			name:  "same position",
			pos:   NewPosition(1, 5),
			other: NewPosition(1, 5),
			want:  false,
		},
		{
			name:  "later row",
			pos:   NewPosition(5, 3),
			other: NewPosition(2, 10),
			want:  false,
		},
		{
			name:  "same row, later col",
			pos:   NewPosition(1, 10),
			other: NewPosition(1, 5),
			want:  false,
		},
		{
			name:  "row 0 col 0 before anything",
			pos:   NewPosition(0, 0),
			other: NewPosition(0, 1),
			want:  true,
		},
		{
			name:  "large positions",
			pos:   NewPosition(999, 888),
			other: NewPosition(1000, 0),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsBefore(tt.other)
			if result != tt.want {
				t.Errorf("IsBefore() = %v, want %v (pos=%v, other=%v)", result, tt.want, tt.pos, tt.other)
			}
		})
	}
}

func TestPosition_IsAfter(t *testing.T) {
	tests := []struct {
		name  string
		pos   Position
		other Position
		want  bool
	}{
		{
			name:  "later row",
			pos:   NewPosition(5, 3),
			other: NewPosition(2, 10),
			want:  true,
		},
		{
			name:  "same row, later col",
			pos:   NewPosition(1, 10),
			other: NewPosition(1, 5),
			want:  true,
		},
		{
			name:  "same position",
			pos:   NewPosition(1, 5),
			other: NewPosition(1, 5),
			want:  false,
		},
		{
			name:  "earlier row",
			pos:   NewPosition(1, 5),
			other: NewPosition(2, 3),
			want:  false,
		},
		{
			name:  "same row, earlier col",
			pos:   NewPosition(1, 3),
			other: NewPosition(1, 5),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsAfter(tt.other)
			if result != tt.want {
				t.Errorf("IsAfter() = %v, want %v (pos=%v, other=%v)", result, tt.want, tt.pos, tt.other)
			}
		})
	}
}

func TestPosition_Equals(t *testing.T) {
	tests := []struct {
		name  string
		pos   Position
		other Position
		want  bool
	}{
		{
			name:  "same position",
			pos:   NewPosition(1, 5),
			other: NewPosition(1, 5),
			want:  true,
		},
		{
			name:  "zero positions",
			pos:   NewPosition(0, 0),
			other: NewPosition(0, 0),
			want:  true,
		},
		{
			name:  "different row",
			pos:   NewPosition(1, 5),
			other: NewPosition(2, 5),
			want:  false,
		},
		{
			name:  "different col",
			pos:   NewPosition(1, 5),
			other: NewPosition(1, 6),
			want:  false,
		},
		{
			name:  "completely different",
			pos:   NewPosition(1, 5),
			other: NewPosition(10, 20),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.Equals(tt.other)
			if result != tt.want {
				t.Errorf("Equals() = %v, want %v (pos=%v, other=%v)", result, tt.want, tt.pos, tt.other)
			}
		})
	}
}

func TestPosition_Comparisons(t *testing.T) {
	// Test that IsBefore and IsAfter are inverses when positions are different
	tests := []struct {
		name string
		pos1 Position
		pos2 Position
	}{
		{
			name: "different rows",
			pos1: NewPosition(1, 5),
			pos2: NewPosition(2, 3),
		},
		{
			name: "same row different col",
			pos1: NewPosition(1, 3),
			pos2: NewPosition(1, 5),
		},
		{
			name: "large difference",
			pos1: NewPosition(0, 0),
			pos2: NewPosition(100, 200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If pos1 is before pos2, then pos2 should be after pos1
			if tt.pos1.IsBefore(tt.pos2) {
				if !tt.pos2.IsAfter(tt.pos1) {
					t.Error("IsBefore and IsAfter are not inverses")
				}
			}

			// If pos1 is after pos2, then pos2 should be before pos1
			if tt.pos1.IsAfter(tt.pos2) {
				if !tt.pos2.IsBefore(tt.pos1) {
					t.Error("IsAfter and IsBefore are not inverses")
				}
			}

			// If positions are equal, neither should be before or after
			if tt.pos1.Equals(tt.pos2) {
				if tt.pos1.IsBefore(tt.pos2) || tt.pos1.IsAfter(tt.pos2) {
					t.Error("Equal positions should not be before or after each other")
				}
			}
		})
	}
}

func TestPosition_ValueObject(t *testing.T) {
	// Test that Position behaves as a value object
	pos1 := NewPosition(5, 10)
	pos2 := NewPosition(5, 10)
	pos3 := NewPosition(5, 11)

	// Equal positions should be equal by value
	if !pos1.Equals(pos2) {
		t.Error("Equal positions should compare equal")
	}

	// Different positions should not be equal
	if pos1.Equals(pos3) {
		t.Error("Different positions should not compare equal")
	}

	// Position should be comparable by value (not by reference)
	if pos1 == pos2 {
		// This is fine - value types can be compared with ==
	}
}

func TestPosition_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "position at origin",
			test: func(t *testing.T) {
				pos := NewPosition(0, 0)
				if pos.Row() != 0 || pos.Col() != 0 {
					t.Error("Origin position should be (0, 0)")
				}
			},
		},
		{
			name: "position with zero row but positive col",
			test: func(t *testing.T) {
				pos := NewPosition(0, 10)
				other := NewPosition(1, 0)
				if !pos.IsBefore(other) {
					t.Error("Row 0 should be before row 1")
				}
			},
		},
		{
			name: "position with zero col but positive row",
			test: func(t *testing.T) {
				pos := NewPosition(10, 0)
				other := NewPosition(5, 100)
				if !pos.IsAfter(other) {
					t.Error("Row 10 should be after row 5")
				}
			},
		},
		{
			name: "consecutive columns",
			test: func(t *testing.T) {
				pos1 := NewPosition(5, 10)
				pos2 := NewPosition(5, 11)
				if !pos1.IsBefore(pos2) {
					t.Error("Column 10 should be before column 11 on same row")
				}
			},
		},
		{
			name: "consecutive rows",
			test: func(t *testing.T) {
				pos1 := NewPosition(5, 999)
				pos2 := NewPosition(6, 0)
				if !pos1.IsBefore(pos2) {
					t.Error("Row 5 should be before row 6 regardless of column")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestPosition_String(t *testing.T) {
	// Test String() method exists and returns something
	pos := NewPosition(3, 5)
	str := pos.String()

	if str == "" {
		t.Error("String() should return non-empty string")
	}

	// Just verify it doesn't panic and returns a string
	// The exact format is not critical for value objects
}
