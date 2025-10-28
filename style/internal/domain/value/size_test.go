package value

import "testing"

// --- Constructor Tests ---

func TestNewSize(t *testing.T) {
	s := NewSize()

	if _, isSet := s.Width(); isSet {
		t.Errorf("Width() should not be set for NewSize()")
	}
	if _, isSet := s.Height(); isSet {
		t.Errorf("Height() should not be set for NewSize()")
	}
	if _, isSet := s.MinWidth(); isSet {
		t.Errorf("MinWidth() should not be set for NewSize()")
	}
	if _, isSet := s.MaxWidth(); isSet {
		t.Errorf("MaxWidth() should not be set for NewSize()")
	}
	if _, isSet := s.MinHeight(); isSet {
		t.Errorf("MinHeight() should not be set for NewSize()")
	}
	if _, isSet := s.MaxHeight(); isSet {
		t.Errorf("MaxHeight() should not be set for NewSize()")
	}
}

func TestWithWidth(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  int
	}{
		{name: "positive width", width: 100, want: 100},
		{name: "zero width", width: 0, want: 0},
		{name: "negative width clamped", width: -10, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := WithWidth(tt.width)

			if got, isSet := s.Width(); !isSet {
				t.Errorf("Width() should be set")
			} else if got != tt.want {
				t.Errorf("Width() = %d, want %d", got, tt.want)
			}

			// Other constraints should not be set
			if _, isSet := s.Height(); isSet {
				t.Errorf("Height() should not be set")
			}
		})
	}
}

func TestWithHeight(t *testing.T) {
	tests := []struct {
		name   string
		height int
		want   int
	}{
		{name: "positive height", height: 50, want: 50},
		{name: "zero height", height: 0, want: 0},
		{name: "negative height clamped", height: -5, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := WithHeight(tt.height)

			if got, isSet := s.Height(); !isSet {
				t.Errorf("Height() should be set")
			} else if got != tt.want {
				t.Errorf("Height() = %d, want %d", got, tt.want)
			}

			// Other constraints should not be set
			if _, isSet := s.Width(); isSet {
				t.Errorf("Width() should not be set")
			}
		})
	}
}

// --- Getter Tests ---

func TestSizeGetters(t *testing.T) {
	s := NewSize().
		SetWidth(100).
		SetHeight(50).
		SetMinWidth(10).
		SetMaxWidth(200).
		SetMinHeight(5).
		SetMaxHeight(100)

	// Test all getters
	if got, isSet := s.Width(); !isSet || got != 100 {
		t.Errorf("Width() = (%d, %v), want (100, true)", got, isSet)
	}
	if got, isSet := s.Height(); !isSet || got != 50 {
		t.Errorf("Height() = (%d, %v), want (50, true)", got, isSet)
	}
	if got, isSet := s.MinWidth(); !isSet || got != 10 {
		t.Errorf("MinWidth() = (%d, %v), want (10, true)", got, isSet)
	}
	if got, isSet := s.MaxWidth(); !isSet || got != 200 {
		t.Errorf("MaxWidth() = (%d, %v), want (200, true)", got, isSet)
	}
	if got, isSet := s.MinHeight(); !isSet || got != 5 {
		t.Errorf("MinHeight() = (%d, %v), want (5, true)", got, isSet)
	}
	if got, isSet := s.MaxHeight(); !isSet || got != 100 {
		t.Errorf("MaxHeight() = (%d, %v), want (100, true)", got, isSet)
	}
}

func TestSizeGettersUnset(t *testing.T) {
	s := NewSize()

	// All should return (0, false)
	if got, isSet := s.Width(); isSet {
		t.Errorf("Width() = (%d, %v), want (0, false)", got, isSet)
	}
	if got, isSet := s.Height(); isSet {
		t.Errorf("Height() = (%d, %v), want (0, false)", got, isSet)
	}
	if got, isSet := s.MinWidth(); isSet {
		t.Errorf("MinWidth() = (%d, %v), want (0, false)", got, isSet)
	}
	if got, isSet := s.MaxWidth(); isSet {
		t.Errorf("MaxWidth() = (%d, %v), want (0, false)", got, isSet)
	}
	if got, isSet := s.MinHeight(); isSet {
		t.Errorf("MinHeight() = (%d, %v), want (0, false)", got, isSet)
	}
	if got, isSet := s.MaxHeight(); isSet {
		t.Errorf("MaxHeight() = (%d, %v), want (0, false)", got, isSet)
	}
}

// --- Setter Tests (Immutability) ---

func TestSetWidthImmutability(t *testing.T) {
	s1 := NewSize()
	s2 := s1.SetWidth(100)

	// s1 should not be modified
	if _, isSet := s1.Width(); isSet {
		t.Errorf("s1.Width() should not be set (immutability violated)")
	}

	// s2 should have width set
	if got, isSet := s2.Width(); !isSet || got != 100 {
		t.Errorf("s2.Width() = (%d, %v), want (100, true)", got, isSet)
	}
}

func TestSetHeightImmutability(t *testing.T) {
	s1 := NewSize()
	s2 := s1.SetHeight(50)

	// s1 should not be modified
	if _, isSet := s1.Height(); isSet {
		t.Errorf("s1.Height() should not be set (immutability violated)")
	}

	// s2 should have height set
	if got, isSet := s2.Height(); !isSet || got != 50 {
		t.Errorf("s2.Height() = (%d, %v), want (50, true)", got, isSet)
	}
}

func TestSettersChaining(t *testing.T) {
	s := NewSize().
		SetWidth(100).
		SetHeight(50).
		SetMinWidth(10).
		SetMaxWidth(200)

	if got, _ := s.Width(); got != 100 {
		t.Errorf("Width() = %d, want 100", got)
	}
	if got, _ := s.Height(); got != 50 {
		t.Errorf("Height() = %d, want 50", got)
	}
	if got, _ := s.MinWidth(); got != 10 {
		t.Errorf("MinWidth() = %d, want 10", got)
	}
	if got, _ := s.MaxWidth(); got != 200 {
		t.Errorf("MaxWidth() = %d, want 200", got)
	}
}

func TestSettersClampNegative(t *testing.T) {
	s := NewSize().
		SetWidth(-10).
		SetHeight(-5).
		SetMinWidth(-3).
		SetMaxWidth(-7)

	// All should be clamped to 0
	if got, _ := s.Width(); got != 0 {
		t.Errorf("Width() = %d, want 0 (clamped)", got)
	}
	if got, _ := s.Height(); got != 0 {
		t.Errorf("Height() = %d, want 0 (clamped)", got)
	}
	if got, _ := s.MinWidth(); got != 0 {
		t.Errorf("MinWidth() = %d, want 0 (clamped)", got)
	}
	if got, _ := s.MaxWidth(); got != 0 {
		t.Errorf("MaxWidth() = %d, want 0 (clamped)", got)
	}
}

// --- Validation Tests ---

func TestValidateWidth(t *testing.T) {
	tests := []struct {
		name  string
		size  Size
		input int
		want  int
	}{
		{
			name:  "no constraints - pass through",
			size:  NewSize(),
			input: 100,
			want:  100,
		},
		{
			name:  "min constraint - clamp up",
			size:  NewSize().SetMinWidth(50),
			input: 30,
			want:  50,
		},
		{
			name:  "min constraint - pass through",
			size:  NewSize().SetMinWidth(50),
			input: 100,
			want:  100,
		},
		{
			name:  "max constraint - clamp down",
			size:  NewSize().SetMaxWidth(100),
			input: 150,
			want:  100,
		},
		{
			name:  "max constraint - pass through",
			size:  NewSize().SetMaxWidth(100),
			input: 50,
			want:  50,
		},
		{
			name:  "both constraints - within range",
			size:  NewSize().SetMinWidth(10).SetMaxWidth(100),
			input: 50,
			want:  50,
		},
		{
			name:  "both constraints - below min",
			size:  NewSize().SetMinWidth(10).SetMaxWidth(100),
			input: 5,
			want:  10,
		},
		{
			name:  "both constraints - above max",
			size:  NewSize().SetMinWidth(10).SetMaxWidth(100),
			input: 150,
			want:  100,
		},
		{
			name:  "both constraints - at min",
			size:  NewSize().SetMinWidth(10).SetMaxWidth(100),
			input: 10,
			want:  10,
		},
		{
			name:  "both constraints - at max",
			size:  NewSize().SetMinWidth(10).SetMaxWidth(100),
			input: 100,
			want:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.ValidateWidth(tt.input); got != tt.want {
				t.Errorf("ValidateWidth(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateHeight(t *testing.T) {
	tests := []struct {
		name  string
		size  Size
		input int
		want  int
	}{
		{
			name:  "no constraints - pass through",
			size:  NewSize(),
			input: 50,
			want:  50,
		},
		{
			name:  "min constraint - clamp up",
			size:  NewSize().SetMinHeight(20),
			input: 10,
			want:  20,
		},
		{
			name:  "min constraint - pass through",
			size:  NewSize().SetMinHeight(20),
			input: 50,
			want:  50,
		},
		{
			name:  "max constraint - clamp down",
			size:  NewSize().SetMaxHeight(50),
			input: 100,
			want:  50,
		},
		{
			name:  "max constraint - pass through",
			size:  NewSize().SetMaxHeight(50),
			input: 30,
			want:  30,
		},
		{
			name:  "both constraints - within range",
			size:  NewSize().SetMinHeight(5).SetMaxHeight(50),
			input: 25,
			want:  25,
		},
		{
			name:  "both constraints - below min",
			size:  NewSize().SetMinHeight(5).SetMaxHeight(50),
			input: 2,
			want:  5,
		},
		{
			name:  "both constraints - above max",
			size:  NewSize().SetMinHeight(5).SetMaxHeight(50),
			input: 100,
			want:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.ValidateHeight(tt.input); got != tt.want {
				t.Errorf("ValidateHeight(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// --- Edge Case Tests ---

func TestSizeEdgeCases(t *testing.T) {
	t.Run("min greater than max width", func(t *testing.T) {
		// If user sets min > max, ValidateWidth applies min first, then max
		// This is an edge case where the user has conflicting constraints
		s := NewSize().SetMinWidth(100).SetMaxWidth(50)

		// Input below min: should clamp to min (100), then max (50) = 50
		if got := s.ValidateWidth(30); got != 50 {
			t.Errorf("ValidateWidth(30) = %d, want 50 (min clamps to 100, then max clamps to 50)", got)
		}

		// Input above max: should clamp to max (50)
		if got := s.ValidateWidth(150); got != 50 {
			t.Errorf("ValidateWidth(150) = %d, want 50 (clamped to max)", got)
		}

		// Input in the conflicting range: should clamp to max (50)
		if got := s.ValidateWidth(75); got != 50 {
			t.Errorf("ValidateWidth(75) = %d, want 50 (clamped to max)", got)
		}
	})

	t.Run("zero width constraint", func(t *testing.T) {
		s := NewSize().SetWidth(0)

		if got, isSet := s.Width(); !isSet || got != 0 {
			t.Errorf("Width() = (%d, %v), want (0, true)", got, isSet)
		}
	})

	t.Run("large values", func(t *testing.T) {
		s := NewSize().SetWidth(1000000).SetHeight(500000)

		if got, _ := s.Width(); got != 1000000 {
			t.Errorf("Width() = %d, want 1000000", got)
		}
		if got, _ := s.Height(); got != 500000 {
			t.Errorf("Height() = %d, want 500000", got)
		}
	})
}

// --- String Tests ---

func TestSizeString(t *testing.T) {
	tests := []struct {
		name string
		size Size
		want string
	}{
		{
			name: "all unset",
			size: NewSize(),
			want: "Size(width=unset, height=unset, minW=unset, maxW=unset, minH=unset, maxH=unset)",
		},
		{
			name: "width and height set",
			size: NewSize().SetWidth(100).SetHeight(50),
			want: "Size(width=100, height=50, minW=unset, maxW=unset, minH=unset, maxH=unset)",
		},
		{
			name: "all constraints set",
			size: NewSize().SetWidth(100).SetHeight(50).SetMinWidth(10).SetMaxWidth(200).SetMinHeight(5).SetMaxHeight(100),
			want: "Size(width=100, height=50, minW=10, maxW=200, minH=5, maxH=100)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
