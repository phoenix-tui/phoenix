package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name string
		x, y int
		want Position
	}{
		{"origin", 0, 0, Position{0, 0}},
		{"positive", 10, 20, Position{10, 20}},
		{"negative", -5, -10, Position{-5, -10}},
		{"mixed", -5, 10, Position{-5, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPosition(tt.x, tt.y)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.x, got.X())
			assert.Equal(t, tt.y, got.Y())
		})
	}
}

func TestPosition_Add(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		dx, dy   int
		expected Position
	}{
		{"add positive", NewPosition(10, 20), 5, 10, NewPosition(15, 30)},
		{"add negative", NewPosition(10, 20), -5, -10, NewPosition(5, 10)},
		{"add zero", NewPosition(10, 20), 0, 0, NewPosition(10, 20)},
		{"add from zero", NewPosition(0, 0), 5, 10, NewPosition(5, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.Add(tt.dx, tt.dy)
			assert.Equal(t, tt.expected, result)
			// Verify immutability
			assert.Equal(t, tt.pos.X(), tt.pos.X()) // Original unchanged
		})
	}
}

func TestPosition_Equals(t *testing.T) {
	tests := []struct {
		name     string
		pos1     Position
		pos2     Position
		expected bool
	}{
		{"same position", NewPosition(10, 20), NewPosition(10, 20), true},
		{"different x", NewPosition(10, 20), NewPosition(15, 20), false},
		{"different y", NewPosition(10, 20), NewPosition(10, 25), false},
		{"both different", NewPosition(10, 20), NewPosition(15, 25), false},
		{"both zero", NewPosition(0, 0), NewPosition(0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos1.Equals(tt.pos2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPosition_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		expected bool
	}{
		{"zero", NewPosition(0, 0), true},
		{"non-zero x", NewPosition(1, 0), false},
		{"non-zero y", NewPosition(0, 1), false},
		{"both non-zero", NewPosition(1, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsZero()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPosition_WithX(t *testing.T) {
	pos := NewPosition(10, 20)
	result := pos.WithX(15)

	assert.Equal(t, 15, result.X())
	assert.Equal(t, 20, result.Y())
	// Verify immutability
	assert.Equal(t, 10, pos.X())
}

func TestPosition_WithY(t *testing.T) {
	pos := NewPosition(10, 20)
	result := pos.WithY(25)

	assert.Equal(t, 10, result.X())
	assert.Equal(t, 25, result.Y())
	// Verify immutability
	assert.Equal(t, 20, pos.Y())
}
