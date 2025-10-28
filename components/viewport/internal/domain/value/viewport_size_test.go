package value

import "testing"

func TestNewViewportSize(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		wantWidth  int
		wantHeight int
	}{
		{"positive dimensions", 80, 24, 80, 24},
		{"zero dimensions", 0, 0, 0, 0},
		{"negative width clamped", -10, 24, 0, 24},
		{"negative height clamped", 80, -10, 80, 0},
		{"both negative clamped", -10, -20, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewViewportSize(tt.width, tt.height)
			if got := s.Width(); got != tt.wantWidth {
				t.Errorf("Width() = %d, want %d", got, tt.wantWidth)
			}
			if got := s.Height(); got != tt.wantHeight {
				t.Errorf("Height() = %d, want %d", got, tt.wantHeight)
			}
		})
	}
}

func TestViewportSize_Width(t *testing.T) {
	s := NewViewportSize(100, 50)
	if got := s.Width(); got != 100 {
		t.Errorf("Width() = %d, want 100", got)
	}
}

func TestViewportSize_Height(t *testing.T) {
	s := NewViewportSize(100, 50)
	if got := s.Height(); got != 50 {
		t.Errorf("Height() = %d, want 50", got)
	}
}

func TestViewportSize_WithWidth(t *testing.T) {
	tests := []struct {
		name     string
		initial  *ViewportSize
		newWidth int
		want     int
	}{
		{"set positive width", NewViewportSize(80, 24), 100, 100},
		{"set zero width", NewViewportSize(80, 24), 0, 0},
		{"set negative width clamped", NewViewportSize(80, 24), -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.WithWidth(tt.newWidth)
			if got := result.Width(); got != tt.want {
				t.Errorf("WithWidth(%d).Width() = %d, want %d", tt.newWidth, got, tt.want)
			}
			// Height should remain unchanged.
			if got := result.Height(); got != tt.initial.Height() {
				t.Errorf("WithWidth() changed height: got %d, want %d", got, tt.initial.Height())
			}
			// Verify immutability.
			if tt.initial.Width() != 80 {
				t.Errorf("WithWidth() modified original width: got %d, want 80", tt.initial.Width())
			}
		})
	}
}

func TestViewportSize_WithHeight(t *testing.T) {
	tests := []struct {
		name      string
		initial   *ViewportSize
		newHeight int
		want      int
	}{
		{"set positive height", NewViewportSize(80, 24), 30, 30},
		{"set zero height", NewViewportSize(80, 24), 0, 0},
		{"set negative height clamped", NewViewportSize(80, 24), -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.WithHeight(tt.newHeight)
			if got := result.Height(); got != tt.want {
				t.Errorf("WithHeight(%d).Height() = %d, want %d", tt.newHeight, got, tt.want)
			}
			// Width should remain unchanged.
			if got := result.Width(); got != tt.initial.Width() {
				t.Errorf("WithHeight() changed width: got %d, want %d", got, tt.initial.Width())
			}
			// Verify immutability.
			if tt.initial.Height() != 24 {
				t.Errorf("WithHeight() modified original height: got %d, want 24", tt.initial.Height())
			}
		})
	}
}

func TestViewportSize_Immutability(t *testing.T) {
	original := NewViewportSize(80, 24)

	// Perform various operations.
	_ = original.WithWidth(100)
	_ = original.WithHeight(30)

	// Original should remain unchanged.
	if got := original.Width(); got != 80 {
		t.Errorf("ViewportSize width was mutated: got %d, want 80", got)
	}
	if got := original.Height(); got != 24 {
		t.Errorf("ViewportSize height was mutated: got %d, want 24", got)
	}
}

func TestViewportSize_ChainedOperations(t *testing.T) {
	original := NewViewportSize(80, 24)
	result := original.WithWidth(100).WithHeight(30)

	if got := result.Width(); got != 100 {
		t.Errorf("Chained operations: Width() = %d, want 100", got)
	}
	if got := result.Height(); got != 30 {
		t.Errorf("Chained operations: Height() = %d, want 30", got)
	}
	// Original should remain unchanged.
	if original.Width() != 80 || original.Height() != 24 {
		t.Errorf("Chained operations modified original: got (%d, %d), want (80, 24)", original.Width(), original.Height())
	}
}
