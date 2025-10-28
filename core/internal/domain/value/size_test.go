package value_test

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/core/internal/domain/value"
)

func TestNewSize(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected value2.Size
	}{
		{
			name:     "valid size",
			width:    80,
			height:   24,
			expected: value2.Size{Width: 80, Height: 24},
		},
		{
			name:     "minimum size",
			width:    1,
			height:   1,
			expected: value2.Size{Width: 1, Height: 1},
		},
		{
			name:     "width < 1 clamped to 1",
			width:    0,
			height:   24,
			expected: value2.Size{Width: 1, Height: 24},
		},
		{
			name:     "height < 1 clamped to 1",
			width:    80,
			height:   0,
			expected: value2.Size{Width: 80, Height: 1},
		},
		{
			name:     "negative width clamped to 1",
			width:    -10,
			height:   24,
			expected: value2.Size{Width: 1, Height: 24},
		},
		{
			name:     "negative height clamped to 1",
			width:    80,
			height:   -10,
			expected: value2.Size{Width: 80, Height: 1},
		},
		{
			name:     "both negative clamped to 1",
			width:    -80,
			height:   -24,
			expected: value2.Size{Width: 1, Height: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := value2.NewSize(tt.width, tt.height)

			if size.Width != tt.expected.Width {
				t.Errorf("expected width %d, got %d", tt.expected.Width, size.Width)
			}
			if size.Height != tt.expected.Height {
				t.Errorf("expected height %d, got %d", tt.expected.Height, size.Height)
			}
		})
	}
}

func TestSize_Area(t *testing.T) {
	tests := []struct {
		name     string
		size     value2.Size
		expected int
	}{
		{
			name:     "standard terminal",
			size:     value2.NewSize(80, 24),
			expected: 1920,
		},
		{
			name:     "minimum size",
			size:     value2.NewSize(1, 1),
			expected: 1,
		},
		{
			name:     "large terminal",
			size:     value2.NewSize(200, 60),
			expected: 12000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area := tt.size.Area()

			if area != tt.expected {
				t.Errorf("expected area %d, got %d", tt.expected, area)
			}
		})
	}
}

func TestSize_Contains(t *testing.T) {
	size := value2.NewSize(80, 24)

	tests := []struct {
		name     string
		pos      value2.Position
		expected bool
	}{
		{
			name:     "top-left corner",
			pos:      value2.NewPosition(0, 0),
			expected: true,
		},
		{
			name:     "bottom-right corner (inside)",
			pos:      value2.NewPosition(23, 79),
			expected: true,
		},
		{
			name:     "middle",
			pos:      value2.NewPosition(12, 40),
			expected: true,
		},
		{
			name:     "row out of bounds",
			pos:      value2.NewPosition(24, 40),
			expected: false,
		},
		{
			name:     "col out of bounds",
			pos:      value2.NewPosition(12, 80),
			expected: false,
		},
		{
			name:     "both out of bounds",
			pos:      value2.NewPosition(30, 100),
			expected: false,
		},
		// Note: NewPosition clamps negative values to 0, so these become (0, 40) and (12, 0)
		// which ARE inside bounds. We can't test truly negative positions through NewPosition.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := size.Contains(tt.pos)

			if result != tt.expected {
				t.Errorf("expected %v, got %v for position (%d, %d)",
					tt.expected, result, tt.pos.Row, tt.pos.Col)
			}
		})
	}
}

func TestSize_Equal(t *testing.T) {
	tests := []struct {
		name     string
		size1    value2.Size
		size2    value2.Size
		expected bool
	}{
		{
			name:     "equal sizes",
			size1:    value2.NewSize(80, 24),
			size2:    value2.NewSize(80, 24),
			expected: true,
		},
		{
			name:     "different width",
			size1:    value2.NewSize(80, 24),
			size2:    value2.NewSize(100, 24),
			expected: false,
		},
		{
			name:     "different height",
			size1:    value2.NewSize(80, 24),
			size2:    value2.NewSize(80, 30),
			expected: false,
		},
		{
			name:     "both different",
			size1:    value2.NewSize(80, 24),
			size2:    value2.NewSize(100, 30),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.size1.Equal(tt.size2)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSize_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		size     value2.Size
		expected bool
	}{
		{
			name:     "normal size",
			size:     value2.NewSize(80, 24),
			expected: false,
		},
		{
			name:     "minimum size",
			size:     value2.NewSize(1, 1),
			expected: false,
		},
		{
			name:     "zero width (clamped)",
			size:     value2.NewSize(0, 24),
			expected: false, // Clamped to 1
		},
		{
			name:     "zero height (clamped)",
			size:     value2.NewSize(80, 0),
			expected: false, // Clamped to 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.size.IsEmpty()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
