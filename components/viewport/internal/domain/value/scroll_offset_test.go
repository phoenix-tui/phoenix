package value

import "testing"

func TestNewScrollOffset(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		want   int
	}{
		{"zero offset", 0, 0},
		{"positive offset", 10, 10},
		{"negative offset clamped to zero", -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrollOffset(tt.offset)
			if got := s.Offset(); got != tt.want {
				t.Errorf("NewScrollOffset(%d).Offset() = %d, want %d", tt.offset, got, tt.want)
			}
		})
	}
}

func TestScrollOffset_Offset(t *testing.T) {
	s := NewScrollOffset(42)
	if got := s.Offset(); got != 42 {
		t.Errorf("Offset() = %d, want 42", got)
	}
}

func TestScrollOffset_Add(t *testing.T) {
	tests := []struct {
		name      string
		offset    int
		delta     int
		maxOffset int
		want      int
	}{
		{"add positive within bounds", 5, 3, 20, 8},
		{"add negative within bounds", 10, -3, 20, 7},
		{"add exceeds max", 15, 10, 20, 20},
		{"add goes below zero", 5, -10, 20, 0},
		{"add zero", 5, 0, 20, 5},
		{"maxOffset is zero", 0, 5, 0, 0},
		{"negative maxOffset treated as zero", 5, 10, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrollOffset(tt.offset)
			result := s.Add(tt.delta, tt.maxOffset)
			if got := result.Offset(); got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.delta, tt.maxOffset, got, tt.want)
			}
			// Verify immutability.
			if s.Offset() != tt.offset {
				t.Errorf("Add() modified original offset: got %d, want %d", s.Offset(), tt.offset)
			}
		})
	}
}

func TestScrollOffset_Set(t *testing.T) {
	tests := []struct {
		name      string
		initial   int
		newOffset int
		maxOffset int
		want      int
	}{
		{"set within bounds", 5, 10, 20, 10},
		{"set exceeds max", 5, 25, 20, 20},
		{"set below zero", 5, -5, 20, 0},
		{"set to zero", 5, 0, 20, 0},
		{"set to max", 5, 20, 20, 20},
		{"maxOffset is zero", 5, 10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrollOffset(tt.initial)
			result := s.Set(tt.newOffset, tt.maxOffset)
			if got := result.Offset(); got != tt.want {
				t.Errorf("Set(%d, %d) = %d, want %d", tt.newOffset, tt.maxOffset, got, tt.want)
			}
			// Verify immutability.
			if s.Offset() != tt.initial {
				t.Errorf("Set() modified original offset: got %d, want %d", s.Offset(), tt.initial)
			}
		})
	}
}

func TestScrollOffset_Clamp(t *testing.T) {
	tests := []struct {
		name      string
		offset    int
		maxOffset int
		want      int
	}{
		{"within bounds", 10, 20, 10},
		{"exceeds max", 25, 20, 20},
		{"below zero", -5, 20, 0},
		{"at zero", 0, 20, 0},
		{"at max", 20, 20, 20},
		{"maxOffset is zero", 10, 0, 0},
		{"negative maxOffset", 10, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrollOffset(tt.offset)
			originalOffset := s.Offset() // Save actual offset after construction
			result := s.Clamp(tt.maxOffset)
			if got := result.Offset(); got != tt.want {
				t.Errorf("Clamp(%d) = %d, want %d", tt.maxOffset, got, tt.want)
			}
			// Verify immutability.
			if s.Offset() != originalOffset {
				t.Errorf("Clamp() modified original offset: got %d, want %d", s.Offset(), originalOffset)
			}
		})
	}
}

func TestScrollOffset_Immutability(t *testing.T) {
	original := NewScrollOffset(10)

	// Perform various operations.
	_ = original.Add(5, 20)
	_ = original.Set(15, 20)
	_ = original.Clamp(20)

	// Original should remain unchanged.
	if got := original.Offset(); got != 10 {
		t.Errorf("ScrollOffset was mutated: got %d, want 10", got)
	}
}
