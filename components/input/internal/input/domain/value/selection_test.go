package value

import "testing"

func TestNewSelection(t *testing.T) {
	tests := []struct {
		name      string
		start     int
		end       int
		wantStart int
		wantEnd   int
	}{
		{"forward selection", 5, 10, 5, 10},
		{"backward selection normalized", 10, 5, 5, 10},
		{"zero length", 5, 5, 5, 5},
		{"negative start clamped", -5, 10, 0, 10},
		{"negative end clamped", 5, -5, 0, 5},
		{"both negative clamped", -10, -5, 0, 0},
		{"zero to positive", 0, 10, 0, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelection(tt.start, tt.end)
			if s.Start() != tt.wantStart {
				t.Errorf("Start() = %d, want %d", s.Start(), tt.wantStart)
			}
			if s.End() != tt.wantEnd {
				t.Errorf("End() = %d, want %d", s.End(), tt.wantEnd)
			}
		})
	}
}

func TestSelection_Length(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
		want  int
	}{
		{"normal selection", 5, 10, 5},
		{"zero length", 5, 5, 0},
		{"single char", 5, 6, 1},
		{"large selection", 0, 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelection(tt.start, tt.end)
			if got := s.Length(); got != tt.want {
				t.Errorf("Length() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSelection_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
		want  bool
	}{
		{"empty at zero", 0, 0, true},
		{"empty at position", 5, 5, true},
		{"non-empty", 5, 10, false},
		{"single char", 5, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelection(tt.start, tt.end)
			if got := s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelection_Contains(t *testing.T) {
	s := NewSelection(5, 10)

	tests := []struct {
		name   string
		offset int
		want   bool
	}{
		{"before start", 4, false},
		{"at start", 5, true},
		{"in middle", 7, true},
		{"before end", 9, true},
		{"at end", 10, false}, // end is exclusive
		{"after end", 11, false},
		{"at zero", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Contains(tt.offset); got != tt.want {
				t.Errorf("Contains(%d) = %v, want %v", tt.offset, got, tt.want)
			}
		})
	}
}

func TestSelection_Clamp(t *testing.T) {
	tests := []struct {
		name      string
		start     int
		end       int
		maxOffset int
		wantStart int
		wantEnd   int
	}{
		{"within bounds", 5, 10, 20, 5, 10},
		{"end exceeds max", 5, 15, 10, 5, 10},
		{"both exceed max", 15, 20, 10, 10, 10},
		{"start exceeds max", 15, 18, 10, 10, 10},
		{"zero max", 5, 10, 0, 0, 0},
		{"exact max", 5, 10, 10, 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelection(tt.start, tt.end)
			result := s.Clamp(tt.maxOffset)

			// Check immutability.
			if s.Start() != tt.start || s.End() != tt.end {
				t.Errorf("original selection modified")
			}

			// Check result.
			if result.Start() != tt.wantStart {
				t.Errorf("Clamp(%d).Start() = %d, want %d", tt.maxOffset, result.Start(), tt.wantStart)
			}
			if result.End() != tt.wantEnd {
				t.Errorf("Clamp(%d).End() = %d, want %d", tt.maxOffset, result.End(), tt.wantEnd)
			}
		})
	}
}

func TestSelection_Clone(t *testing.T) {
	original := NewSelection(5, 10)
	clone := original.Clone()

	// Check values match.
	if clone.Start() != original.Start() {
		t.Errorf("Clone().Start() = %d, want %d", clone.Start(), original.Start())
	}
	if clone.End() != original.End() {
		t.Errorf("Clone().End() = %d, want %d", clone.End(), original.End())
	}

	// Check they're different instances.
	if clone == original {
		t.Error("Clone() returned same instance, want different instance")
	}

	// Modify clone and verify original unchanged.
	clone.start = 20
	clone.end = 30
	if original.Start() != 5 || original.End() != 10 {
		t.Errorf("modifying clone affected original: got Start=%d End=%d, want Start=5 End=10",
			original.Start(), original.End())
	}
}
