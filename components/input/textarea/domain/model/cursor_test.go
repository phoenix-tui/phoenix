package model

import (
	"testing"
)

func TestNewCursor(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
	}{
		{
			name: "zero position",
			row:  0,
			col:  0,
		},
		{
			name: "positive position",
			row:  5,
			col:  10,
		},
		{
			name: "large position",
			row:  1000,
			col:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor := NewCursor(tt.row, tt.col)

			if cursor.Row() != tt.row {
				t.Errorf("Row() = %d, want %d", cursor.Row(), tt.row)
			}
			if cursor.Col() != tt.col {
				t.Errorf("Col() = %d, want %d", cursor.Col(), tt.col)
			}
		})
	}
}

func TestCursor_Position(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor := NewCursor(tt.row, tt.col)
			row, col := cursor.Position()

			if row != tt.wantRow {
				t.Errorf("Position() row = %d, want %d", row, tt.wantRow)
			}
			if col != tt.wantCol {
				t.Errorf("Position() col = %d, want %d", col, tt.wantCol)
			}
		})
	}
}

func TestCursor_MoveTo(t *testing.T) {
	tests := []struct {
		name    string
		initial *Cursor
		newRow  int
		newCol  int
		wantRow int
		wantCol int
		origRow int
		origCol int
	}{
		{
			name:    "move to valid position",
			initial: NewCursor(0, 0),
			newRow:  5,
			newCol:  10,
			wantRow: 5,
			wantCol: 10,
			origRow: 0,
			origCol: 0,
		},
		{
			name:    "move to same position",
			initial: NewCursor(3, 7),
			newRow:  3,
			newCol:  7,
			wantRow: 3,
			wantCol: 7,
			origRow: 3,
			origCol: 7,
		},
		{
			name:    "move to zero",
			initial: NewCursor(10, 20),
			newRow:  0,
			newCol:  0,
			wantRow: 0,
			wantCol: 0,
			origRow: 10,
			origCol: 20,
		},
		{
			name:    "move to large position",
			initial: NewCursor(5, 5),
			newRow:  1000,
			newCol:  2000,
			wantRow: 1000,
			wantCol: 2000,
			origRow: 5,
			origCol: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveTo(tt.newRow, tt.newCol)

			if result.Row() != tt.wantRow {
				t.Errorf("MoveTo() row = %d, want %d", result.Row(), tt.wantRow)
			}
			if result.Col() != tt.wantCol {
				t.Errorf("MoveTo() col = %d, want %d", result.Col(), tt.wantCol)
			}

			// Verify immutability
			if tt.initial.Row() != tt.origRow {
				t.Errorf("Original cursor row was modified: %d, want %d", tt.initial.Row(), tt.origRow)
			}
			if tt.initial.Col() != tt.origCol {
				t.Errorf("Original cursor col was modified: %d, want %d", tt.initial.Col(), tt.origCol)
			}
		})
	}
}

func TestCursor_MoveBy(t *testing.T) {
	tests := []struct {
		name     string
		initial  *Cursor
		deltaRow int
		deltaCol int
		wantRow  int
		wantCol  int
		origRow  int
		origCol  int
	}{
		{
			name:     "move by positive delta",
			initial:  NewCursor(5, 10),
			deltaRow: 2,
			deltaCol: 3,
			wantRow:  7,
			wantCol:  13,
			origRow:  5,
			origCol:  10,
		},
		{
			name:     "move by negative delta",
			initial:  NewCursor(10, 20),
			deltaRow: -3,
			deltaCol: -5,
			wantRow:  7,
			wantCol:  15,
			origRow:  10,
			origCol:  20,
		},
		{
			name:     "move by zero delta",
			initial:  NewCursor(5, 10),
			deltaRow: 0,
			deltaCol: 0,
			wantRow:  5,
			wantCol:  10,
			origRow:  5,
			origCol:  10,
		},
		{
			name:     "move by mixed delta",
			initial:  NewCursor(10, 10),
			deltaRow: 5,
			deltaCol: -5,
			wantRow:  15,
			wantCol:  5,
			origRow:  10,
			origCol:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveBy(tt.deltaRow, tt.deltaCol)

			if result.Row() != tt.wantRow {
				t.Errorf("MoveBy() row = %d, want %d", result.Row(), tt.wantRow)
			}
			if result.Col() != tt.wantCol {
				t.Errorf("MoveBy() col = %d, want %d", result.Col(), tt.wantCol)
			}

			// Verify immutability
			if tt.initial.Row() != tt.origRow {
				t.Errorf("Original cursor row was modified: %d, want %d", tt.initial.Row(), tt.origRow)
			}
			if tt.initial.Col() != tt.origCol {
				t.Errorf("Original cursor col was modified: %d, want %d", tt.initial.Col(), tt.origCol)
			}
		})
	}
}

func TestCursor_Copy(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
	}{
		{
			name: "copy zero position",
			row:  0,
			col:  0,
		},
		{
			name: "copy positive position",
			row:  5,
			col:  10,
		},
		{
			name: "copy large position",
			row:  1000,
			col:  2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := NewCursor(tt.row, tt.col)
			copy := original.Copy()

			if copy.Row() != original.Row() {
				t.Errorf("Copy() row = %d, want %d", copy.Row(), original.Row())
			}
			if copy.Col() != original.Col() {
				t.Errorf("Copy() col = %d, want %d", copy.Col(), original.Col())
			}

			// Verify it's a new instance (not same pointer)
			if copy == original {
				t.Error("Copy() returned same instance, want new instance")
			}

			// Verify modifying copy doesn't affect original
			newCopy := copy.MoveTo(999, 888)
			if original.Row() != tt.row {
				t.Errorf("Original cursor row was modified: %d, want %d", original.Row(), tt.row)
			}
			if original.Col() != tt.col {
				t.Errorf("Original cursor col was modified: %d, want %d", original.Col(), tt.col)
			}
			if newCopy.Row() != 999 || newCopy.Col() != 888 {
				t.Error("Modified copy has incorrect values")
			}
		})
	}
}

func TestCursor_Immutability(t *testing.T) {
	original := NewCursor(5, 10)

	// Test all mutation operations preserve original
	moved := original.MoveTo(20, 30)
	movedBy := original.MoveBy(3, 5)
	copied := original.Copy()

	// Original should remain unchanged
	if original.Row() != 5 {
		t.Errorf("Original row changed after operations: %d, want 5", original.Row())
	}
	if original.Col() != 10 {
		t.Errorf("Original col changed after operations: %d, want 10", original.Col())
	}

	// Results should have correct values
	if moved.Row() != 20 || moved.Col() != 30 {
		t.Error("MoveTo() produced incorrect values")
	}
	if movedBy.Row() != 8 || movedBy.Col() != 15 {
		t.Error("MoveBy() produced incorrect values")
	}
	if copied.Row() != 5 || copied.Col() != 10 {
		t.Error("Copy() produced incorrect values")
	}
}
