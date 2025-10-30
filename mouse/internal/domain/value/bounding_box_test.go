package value

import "testing"

func TestNewBoundingBox(t *testing.T) {
	tests := []struct {
		name           string
		x              int
		y              int
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "positive dimensions",
			x:              5,
			y:              10,
			width:          20,
			height:         15,
			expectedWidth:  20,
			expectedHeight: 15,
		},
		{
			name:           "zero dimensions",
			x:              0,
			y:              0,
			width:          0,
			height:         0,
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "negative width normalized to zero",
			x:              5,
			y:              10,
			width:          -10,
			height:         15,
			expectedWidth:  0,
			expectedHeight: 15,
		},
		{
			name:           "negative height normalized to zero",
			x:              5,
			y:              10,
			width:          20,
			height:         -5,
			expectedWidth:  20,
			expectedHeight: 0,
		},
		{
			name:           "both negative dimensions normalized",
			x:              5,
			y:              10,
			width:          -10,
			height:         -5,
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "negative coordinates allowed",
			x:              -5,
			y:              -10,
			width:          20,
			height:         15,
			expectedWidth:  20,
			expectedHeight: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := NewBoundingBox(tt.x, tt.y, tt.width, tt.height)

			if bb.X() != tt.x {
				t.Errorf("X() = %d, want %d", bb.X(), tt.x)
			}
			if bb.Y() != tt.y {
				t.Errorf("Y() = %d, want %d", bb.Y(), tt.y)
			}
			if bb.Width() != tt.expectedWidth {
				t.Errorf("Width() = %d, want %d", bb.Width(), tt.expectedWidth)
			}
			if bb.Height() != tt.expectedHeight {
				t.Errorf("Height() = %d, want %d", bb.Height(), tt.expectedHeight)
			}
		})
	}
}

func TestBoundingBox_Contains(t *testing.T) {
	tests := []struct {
		name     string
		bb       BoundingBox
		pos      Position
		expected bool
	}{
		{
			name:     "point inside",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(10, 15),
			expected: true,
		},
		{
			name:     "point at top-left corner",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(5, 10),
			expected: true,
		},
		{
			name:     "point at bottom-right corner (exclusive)",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(25, 25),
			expected: false,
		},
		{
			name:     "point just inside right boundary",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(24, 15),
			expected: true,
		},
		{
			name:     "point just inside bottom boundary",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(10, 24),
			expected: true,
		},
		{
			name:     "point outside left",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(4, 15),
			expected: false,
		},
		{
			name:     "point outside top",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(10, 9),
			expected: false,
		},
		{
			name:     "point outside right",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(25, 15),
			expected: false,
		},
		{
			name:     "point outside bottom",
			bb:       NewBoundingBox(5, 10, 20, 15),
			pos:      NewPosition(10, 25),
			expected: false,
		},
		{
			name:     "zero size box contains nothing",
			bb:       NewBoundingBox(5, 10, 0, 0),
			pos:      NewPosition(5, 10),
			expected: false,
		},
		{
			name:     "negative coordinates",
			bb:       NewBoundingBox(-10, -5, 20, 10),
			pos:      NewPosition(-5, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb.Contains(tt.pos); got != tt.expected {
				t.Errorf("BoundingBox.Contains(%v) = %v, want %v", tt.pos, got, tt.expected)
			}
		})
	}
}

func TestBoundingBox_Overlaps(t *testing.T) {
	tests := []struct {
		name     string
		bb1      BoundingBox
		bb2      BoundingBox
		expected bool
	}{
		{
			name:     "overlapping boxes",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(15, 15, 20, 15),
			expected: true,
		},
		{
			name:     "touching boxes (edge to edge)",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(25, 10, 20, 15),
			expected: false,
		},
		{
			name:     "completely contained",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(10, 15, 5, 5),
			expected: true,
		},
		{
			name:     "completely separate",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(50, 50, 20, 15),
			expected: false,
		},
		{
			name:     "identical boxes",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(5, 10, 20, 15),
			expected: true,
		},
		{
			name:     "overlapping at corner",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(20, 20, 20, 15),
			expected: true,
		},
		{
			name:     "zero size boxes",
			bb1:      NewBoundingBox(5, 10, 0, 0),
			bb2:      NewBoundingBox(5, 10, 0, 0),
			expected: false,
		},
		{
			name:     "negative coordinates overlapping",
			bb1:      NewBoundingBox(-10, -5, 20, 10),
			bb2:      NewBoundingBox(-5, 0, 20, 10),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb1.Overlaps(tt.bb2); got != tt.expected {
				t.Errorf("BoundingBox.Overlaps() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoundingBox_Equals(t *testing.T) {
	tests := []struct {
		name     string
		bb1      BoundingBox
		bb2      BoundingBox
		expected bool
	}{
		{
			name:     "equal boxes",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(5, 10, 20, 15),
			expected: true,
		},
		{
			name:     "different x",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(6, 10, 20, 15),
			expected: false,
		},
		{
			name:     "different y",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(5, 11, 20, 15),
			expected: false,
		},
		{
			name:     "different width",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(5, 10, 21, 15),
			expected: false,
		},
		{
			name:     "different height",
			bb1:      NewBoundingBox(5, 10, 20, 15),
			bb2:      NewBoundingBox(5, 10, 20, 16),
			expected: false,
		},
		{
			name:     "zero size boxes",
			bb1:      NewBoundingBox(5, 10, 0, 0),
			bb2:      NewBoundingBox(5, 10, 0, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb1.Equals(tt.bb2); got != tt.expected {
				t.Errorf("BoundingBox.Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoundingBox_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		bb       BoundingBox
		expected bool
	}{
		{
			name:     "non-empty box",
			bb:       NewBoundingBox(5, 10, 20, 15),
			expected: false,
		},
		{
			name:     "zero width",
			bb:       NewBoundingBox(5, 10, 0, 15),
			expected: true,
		},
		{
			name:     "zero height",
			bb:       NewBoundingBox(5, 10, 20, 0),
			expected: true,
		},
		{
			name:     "both zero",
			bb:       NewBoundingBox(5, 10, 0, 0),
			expected: true,
		},
		{
			name:     "negative dimensions normalized to zero",
			bb:       NewBoundingBox(5, 10, -10, -5),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb.IsEmpty(); got != tt.expected {
				t.Errorf("BoundingBox.IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoundingBox_Area(t *testing.T) {
	tests := []struct {
		name     string
		bb       BoundingBox
		expected int
	}{
		{
			name:     "normal box",
			bb:       NewBoundingBox(5, 10, 20, 15),
			expected: 300,
		},
		{
			name:     "square box",
			bb:       NewBoundingBox(0, 0, 10, 10),
			expected: 100,
		},
		{
			name:     "zero width",
			bb:       NewBoundingBox(5, 10, 0, 15),
			expected: 0,
		},
		{
			name:     "zero height",
			bb:       NewBoundingBox(5, 10, 20, 0),
			expected: 0,
		},
		{
			name:     "both zero",
			bb:       NewBoundingBox(5, 10, 0, 0),
			expected: 0,
		},
		{
			name:     "large box",
			bb:       NewBoundingBox(0, 0, 1000, 500),
			expected: 500000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb.Area(); got != tt.expected {
				t.Errorf("BoundingBox.Area() = %d, want %d", got, tt.expected)
			}
		})
	}
}

// TestBoundingBox_BoundaryConditions tests edge cases with boundary conditions.
func TestBoundingBox_BoundaryConditions(t *testing.T) {
	bb := NewBoundingBox(10, 20, 30, 40)

	// Test corners
	corners := []struct {
		name     string
		pos      Position
		expected bool
	}{
		{"top-left inclusive", NewPosition(10, 20), true},
		{"top-right exclusive", NewPosition(40, 20), false},
		{"bottom-left exclusive", NewPosition(10, 60), false},
		{"bottom-right exclusive", NewPosition(40, 60), false},
		{"just inside top-right", NewPosition(39, 20), true},
		{"just inside bottom-left", NewPosition(10, 59), true},
		{"just inside bottom-right", NewPosition(39, 59), true},
	}

	for _, tc := range corners {
		t.Run(tc.name, func(t *testing.T) {
			if got := bb.Contains(tc.pos); got != tc.expected {
				t.Errorf("Contains(%v) = %v, want %v", tc.pos, got, tc.expected)
			}
		})
	}
}

// TestBoundingBox_NegativeDimensions tests that negative dimensions are handled correctly.
func TestBoundingBox_NegativeDimensions(t *testing.T) {
	// Negative width/height should be normalized to zero
	bb := NewBoundingBox(10, 20, -30, -40)

	if bb.Width() != 0 {
		t.Errorf("Negative width should be normalized to 0, got %d", bb.Width())
	}

	if bb.Height() != 0 {
		t.Errorf("Negative height should be normalized to 0, got %d", bb.Height())
	}

	if !bb.IsEmpty() {
		t.Error("Box with negative dimensions should be empty")
	}

	if bb.Area() != 0 {
		t.Errorf("Area of box with negative dimensions should be 0, got %d", bb.Area())
	}

	// Should not contain any point
	if bb.Contains(NewPosition(10, 20)) {
		t.Error("Box with negative dimensions should not contain any point")
	}
}
