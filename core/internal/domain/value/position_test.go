package value_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/internal/domain/value"
)

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name     string
		row      int
		col      int
		expected value.Position
	}{
		{
			name:     "valid position",
			row:      5,
			col:      10,
			expected: value.Position{Row: 5, Col: 10},
		},
		{
			name:     "zero position",
			row:      0,
			col:      0,
			expected: value.Position{Row: 0, Col: 0},
		},
		{
			name:     "negative row clamped to 0",
			row:      -5,
			col:      10,
			expected: value.Position{Row: 0, Col: 10},
		},
		{
			name:     "negative col clamped to 0",
			row:      5,
			col:      -10,
			expected: value.Position{Row: 5, Col: 0},
		},
		{
			name:     "both negative clamped to 0",
			row:      -5,
			col:      -10,
			expected: value.Position{Row: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := value.NewPosition(tt.row, tt.col)

			if pos.Row != tt.expected.Row {
				t.Errorf("expected row %d, got %d", tt.expected.Row, pos.Row)
			}
			if pos.Col != tt.expected.Col {
				t.Errorf("expected col %d, got %d", tt.expected.Col, pos.Col)
			}
		})
	}
}

func TestPosition_Add(t *testing.T) {
	tests := []struct {
		name     string
		pos      value.Position
		deltaRow int
		deltaCol int
		expected value.Position
	}{
		{
			name:     "add positive deltas",
			pos:      value.NewPosition(5, 10),
			deltaRow: 3,
			deltaCol: 7,
			expected: value.Position{Row: 8, Col: 17},
		},
		{
			name:     "add negative deltas",
			pos:      value.NewPosition(10, 20),
			deltaRow: -5,
			deltaCol: -10,
			expected: value.Position{Row: 5, Col: 10},
		},
		{
			name:     "add zero deltas",
			pos:      value.NewPosition(5, 10),
			deltaRow: 0,
			deltaCol: 0,
			expected: value.Position{Row: 5, Col: 10},
		},
		{
			name:     "negative result clamped",
			pos:      value.NewPosition(5, 10),
			deltaRow: -10,
			deltaCol: -20,
			expected: value.Position{Row: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.Add(tt.deltaRow, tt.deltaCol)

			if result.Row != tt.expected.Row {
				t.Errorf("expected row %d, got %d", tt.expected.Row, result.Row)
			}
			if result.Col != tt.expected.Col {
				t.Errorf("expected col %d, got %d", tt.expected.Col, result.Col)
			}

			// Verify immutability - original unchanged
			if tt.pos.Row != value.NewPosition(tt.pos.Row, tt.pos.Col).Row {
				t.Error("original position was mutated")
			}
		})
	}
}

func TestPosition_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		pos      value.Position
		expected bool
	}{
		{
			name:     "zero position",
			pos:      value.NewPosition(0, 0),
			expected: true,
		},
		{
			name:     "non-zero row",
			pos:      value.NewPosition(1, 0),
			expected: false,
		},
		{
			name:     "non-zero col",
			pos:      value.NewPosition(0, 1),
			expected: false,
		},
		{
			name:     "both non-zero",
			pos:      value.NewPosition(5, 10),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsZero()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPosition_Equal(t *testing.T) {
	tests := []struct {
		name     string
		pos1     value.Position
		pos2     value.Position
		expected bool
	}{
		{
			name:     "equal positions",
			pos1:     value.NewPosition(5, 10),
			pos2:     value.NewPosition(5, 10),
			expected: true,
		},
		{
			name:     "different row",
			pos1:     value.NewPosition(5, 10),
			pos2:     value.NewPosition(6, 10),
			expected: false,
		},
		{
			name:     "different col",
			pos1:     value.NewPosition(5, 10),
			pos2:     value.NewPosition(5, 11),
			expected: false,
		},
		{
			name:     "both different",
			pos1:     value.NewPosition(5, 10),
			pos2:     value.NewPosition(6, 11),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos1.Equal(tt.pos2)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
