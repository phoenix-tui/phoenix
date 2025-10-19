package value

import "testing"

func TestNewPercentage(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 0},
		{"Valid middle", 50, 50},
		{"Max value", 100, 100},
		{"Above max - clamp to 100", 150, 100},
		{"Below min - clamp to 0", -50, 0},
		{"Large positive", 999, 100},
		{"Large negative", -999, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPercentage(tt.input)
			if p.Value() != tt.expected {
				t.Errorf("NewPercentage(%d) = %d, expected %d", tt.input, p.Value(), tt.expected)
			}
		})
	}
}

func TestPercentageAdd(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		delta    int
		expected int
	}{
		{"Add to zero", 0, 10, 10},
		{"Add to middle", 50, 20, 70},
		{"Add to max", 100, 10, 100}, // Clamp
		{"Add causing overflow", 90, 50, 100},
		{"Add zero", 50, 0, 50},
		{"Add negative (subtract)", 50, -10, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPercentage(tt.initial)
			result := p.Add(tt.delta)
			if result.Value() != tt.expected {
				t.Errorf("Add(%d, %d) = %d, expected %d", tt.initial, tt.delta, result.Value(), tt.expected)
			}
			// Verify immutability
			if p.Value() != tt.initial {
				t.Errorf("Add() mutated original: %d != %d", p.Value(), tt.initial)
			}
		})
	}
}

func TestPercentageSubtract(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		delta    int
		expected int
	}{
		{"Subtract from max", 100, 10, 90},
		{"Subtract from middle", 50, 20, 30},
		{"Subtract to zero", 10, 10, 0},
		{"Subtract causing underflow", 10, 50, 0},
		{"Subtract zero", 50, 0, 50},
		{"Subtract negative (add)", 50, -10, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPercentage(tt.initial)
			result := p.Subtract(tt.delta)
			if result.Value() != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d, expected %d", tt.initial, tt.delta, result.Value(), tt.expected)
			}
			// Verify immutability
			if p.Value() != tt.initial {
				t.Errorf("Subtract() mutated original: %d != %d", p.Value(), tt.initial)
			}
		})
	}
}

func TestPercentageIsComplete(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{"Zero", 0, false},
		{"Middle", 50, false},
		{"99", 99, false},
		{"100", 100, true},
		{"Over 100 (clamped)", 150, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPercentage(tt.value)
			if p.IsComplete() != tt.expected {
				t.Errorf("IsComplete() for %d = %v, expected %v", tt.value, p.IsComplete(), tt.expected)
			}
		})
	}
}

func TestPercentageImmutability(t *testing.T) {
	p1 := NewPercentage(50)
	p2 := p1.Add(10)
	p3 := p1.Subtract(10)

	if p1.Value() != 50 {
		t.Errorf("Original percentage mutated: %d != 50", p1.Value())
	}
	if p2.Value() != 60 {
		t.Errorf("Add() result incorrect: %d != 60", p2.Value())
	}
	if p3.Value() != 40 {
		t.Errorf("Subtract() result incorrect: %d != 40", p3.Value())
	}
}
