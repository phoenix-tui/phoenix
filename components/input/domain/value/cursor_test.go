package value

import "testing"

func TestNewCursor(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		want   int
	}{
		{"zero offset", 0, 0},
		{"positive offset", 5, 5},
		{"negative offset clamped", -5, 0},
		{"large offset", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.offset)
			if c.Offset() != tt.want {
				t.Errorf("NewCursor(%d).Offset() = %d, want %d", tt.offset, c.Offset(), tt.want)
			}
		})
	}
}

func TestCursor_MoveBy(t *testing.T) {
	tests := []struct {
		name      string
		initial   int
		delta     int
		maxOffset int
		want      int
	}{
		{"move right", 5, 3, 20, 8},
		{"move left", 5, -3, 20, 2},
		{"move beyond max clamped", 5, 20, 10, 10},
		{"move before zero clamped", 5, -10, 20, 0},
		{"zero delta", 5, 0, 20, 5},
		{"at zero move left", 0, -5, 20, 0},
		{"at max move right", 10, 5, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.initial)
			result := c.MoveBy(tt.delta, tt.maxOffset)

			// Check immutability.
			if c.Offset() != tt.initial {
				t.Errorf("original cursor modified: got %d, want %d", c.Offset(), tt.initial)
			}

			// Check result.
			if result.Offset() != tt.want {
				t.Errorf("MoveBy(%d, %d) = %d, want %d", tt.delta, tt.maxOffset, result.Offset(), tt.want)
			}
		})
	}
}

func TestCursor_MoveTo(t *testing.T) {
	tests := []struct {
		name      string
		initial   int
		target    int
		maxOffset int
		want      int
	}{
		{"move to valid position", 5, 10, 20, 10},
		{"move to zero", 5, 0, 20, 0},
		{"move to max", 5, 20, 20, 20},
		{"move beyond max clamped", 5, 30, 20, 20},
		{"move before zero clamped", 5, -10, 20, 0},
		{"move to same position", 5, 5, 20, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCursor(tt.initial)
			result := c.MoveTo(tt.target, tt.maxOffset)

			// Check immutability.
			if c.Offset() != tt.initial {
				t.Errorf("original cursor modified: got %d, want %d", c.Offset(), tt.initial)
			}

			// Check result.
			if result.Offset() != tt.want {
				t.Errorf("MoveTo(%d, %d) = %d, want %d", tt.target, tt.maxOffset, result.Offset(), tt.want)
			}
		})
	}
}

func TestCursor_Clone(t *testing.T) {
	original := NewCursor(5)
	clone := original.Clone()

	// Check values match.
	if clone.Offset() != original.Offset() {
		t.Errorf("Clone().Offset() = %d, want %d", clone.Offset(), original.Offset())
	}

	// Check they're different instances.
	if clone == original {
		t.Error("Clone() returned same instance, want different instance")
	}

	// Modify clone and verify original unchanged.
	clone.offset = 10
	if original.Offset() != 5 {
		t.Errorf("modifying clone affected original: got %d, want 5", original.Offset())
	}
}
