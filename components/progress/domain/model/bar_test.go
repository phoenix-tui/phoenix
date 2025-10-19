package model

import "testing"

func TestNewBar(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		expectedWidth int
	}{
		{"Normal width", 40, 40},
		{"Large width", 100, 100},
		{"Minimum width", 1, 1},
		{"Zero width - defaults to 1", 0, 1},
		{"Negative width - defaults to 1", -10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBar(tt.width)
			if bar.Width() != tt.expectedWidth {
				t.Errorf("Width() = %d, expected %d", bar.Width(), tt.expectedWidth)
			}
			if bar.Percentage() != 0 {
				t.Errorf("Percentage() = %d, expected 0", bar.Percentage())
			}
			if bar.FillChar() != '█' {
				t.Errorf("FillChar() = %c, expected '█'", bar.FillChar())
			}
			if bar.EmptyChar() != '░' {
				t.Errorf("EmptyChar() = %c, expected '░'", bar.EmptyChar())
			}
		})
	}
}

func TestNewBarWithPercentage(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		percentage int
		expected   int
	}{
		{"Zero percent", 40, 0, 0},
		{"Half progress", 40, 50, 50},
		{"Full progress", 40, 100, 100},
		{"Over 100 - clamp", 40, 150, 100},
		{"Negative - clamp", 40, -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarWithPercentage(tt.width, tt.percentage)
			if bar.Percentage() != tt.expected {
				t.Errorf("Percentage() = %d, expected %d", bar.Percentage(), tt.expected)
			}
		})
	}
}

func TestBarWithPercentage(t *testing.T) {
	bar := NewBar(40)

	tests := []struct {
		name     string
		pct      int
		expected int
	}{
		{"Set to 50", 50, 50},
		{"Set to 100", 100, 100},
		{"Set to 0", 0, 0},
		{"Set over 100", 150, 100},
		{"Set negative", -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newBar := bar.WithPercentage(tt.pct)
			if newBar.Percentage() != tt.expected {
				t.Errorf("WithPercentage(%d) = %d, expected %d", tt.pct, newBar.Percentage(), tt.expected)
			}
			// Verify immutability
			if bar.Percentage() != 0 {
				t.Errorf("WithPercentage() mutated original")
			}
		})
	}
}

func TestBarWithChars(t *testing.T) {
	bar := NewBar(40)

	// Test fill char
	newBar := bar.WithFillChar('▓')
	if newBar.FillChar() != '▓' {
		t.Errorf("WithFillChar('▓') = %c, expected '▓'", newBar.FillChar())
	}
	if bar.FillChar() != '█' {
		t.Errorf("WithFillChar() mutated original")
	}

	// Test empty char
	newBar = bar.WithEmptyChar('▒')
	if newBar.EmptyChar() != '▒' {
		t.Errorf("WithEmptyChar('▒') = %c, expected '▒'", newBar.EmptyChar())
	}
	if bar.EmptyChar() != '░' {
		t.Errorf("WithEmptyChar() mutated original")
	}
}

func TestBarWithShowPercent(t *testing.T) {
	bar := NewBar(40)

	newBar := bar.WithShowPercent(true)
	if !newBar.ShowPercent() {
		t.Errorf("WithShowPercent(true) = false, expected true")
	}
	if bar.ShowPercent() {
		t.Errorf("WithShowPercent() mutated original")
	}

	newBar2 := newBar.WithShowPercent(false)
	if newBar2.ShowPercent() {
		t.Errorf("WithShowPercent(false) = true, expected false")
	}
}

func TestBarWithLabel(t *testing.T) {
	bar := NewBar(40)

	newBar := bar.WithLabel("Downloading...")
	if newBar.Label() != "Downloading..." {
		t.Errorf("WithLabel() = %s, expected 'Downloading...'", newBar.Label())
	}
	if bar.Label() != "" {
		t.Errorf("WithLabel() mutated original")
	}
}

func TestBarIncrement(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		delta    int
		expected int
	}{
		{"Increment from 0", 0, 10, 10},
		{"Increment from 50", 50, 20, 70},
		{"Increment to 100", 90, 10, 100},
		{"Increment past 100", 95, 20, 100},
		{"Increment by zero", 50, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarWithPercentage(40, tt.initial)
			newBar := bar.Increment(tt.delta)
			if newBar.Percentage() != tt.expected {
				t.Errorf("Increment(%d) from %d = %d, expected %d",
					tt.delta, tt.initial, newBar.Percentage(), tt.expected)
			}
			// Verify immutability
			if bar.Percentage() != tt.initial {
				t.Errorf("Increment() mutated original")
			}
		})
	}
}

func TestBarDecrement(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		delta    int
		expected int
	}{
		{"Decrement from 100", 100, 10, 90},
		{"Decrement from 50", 50, 20, 30},
		{"Decrement to 0", 10, 10, 0},
		{"Decrement past 0", 10, 30, 0},
		{"Decrement by zero", 50, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarWithPercentage(40, tt.initial)
			newBar := bar.Decrement(tt.delta)
			if newBar.Percentage() != tt.expected {
				t.Errorf("Decrement(%d) from %d = %d, expected %d",
					tt.delta, tt.initial, newBar.Percentage(), tt.expected)
			}
			// Verify immutability
			if bar.Percentage() != tt.initial {
				t.Errorf("Decrement() mutated original")
			}
		})
	}
}

func TestBarSetComplete(t *testing.T) {
	bar := NewBarWithPercentage(40, 50)
	newBar := bar.SetComplete()

	if newBar.Percentage() != 100 {
		t.Errorf("SetComplete() = %d, expected 100", newBar.Percentage())
	}
	if !newBar.IsComplete() {
		t.Errorf("SetComplete() bar should be complete")
	}
	if bar.Percentage() != 50 {
		t.Errorf("SetComplete() mutated original")
	}
}

func TestBarReset(t *testing.T) {
	bar := NewBarWithPercentage(40, 75)
	newBar := bar.Reset()

	if newBar.Percentage() != 0 {
		t.Errorf("Reset() = %d, expected 0", newBar.Percentage())
	}
	if newBar.IsComplete() {
		t.Errorf("Reset() bar should not be complete")
	}
	if bar.Percentage() != 75 {
		t.Errorf("Reset() mutated original")
	}
}

func TestBarIsComplete(t *testing.T) {
	tests := []struct {
		name       string
		percentage int
		expected   bool
	}{
		{"0%", 0, false},
		{"50%", 50, false},
		{"99%", 99, false},
		{"100%", 100, true},
		{"Over 100 (clamped)", 150, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarWithPercentage(40, tt.percentage)
			if bar.IsComplete() != tt.expected {
				t.Errorf("IsComplete() for %d%% = %v, expected %v",
					tt.percentage, bar.IsComplete(), tt.expected)
			}
		})
	}
}

func TestBarImmutability(t *testing.T) {
	bar := NewBar(40).
		WithPercentage(50).
		WithFillChar('▓').
		WithEmptyChar('▒').
		WithShowPercent(true).
		WithLabel("Test")

	// Apply all mutations
	_ = bar.WithPercentage(75)
	_ = bar.WithFillChar('█')
	_ = bar.WithEmptyChar('░')
	_ = bar.WithShowPercent(false)
	_ = bar.WithLabel("Changed")
	_ = bar.Increment(10)
	_ = bar.Decrement(10)
	_ = bar.SetComplete()
	_ = bar.Reset()

	// Original should be unchanged
	if bar.Percentage() != 50 {
		t.Errorf("Percentage mutated: %d != 50", bar.Percentage())
	}
	if bar.FillChar() != '▓' {
		t.Errorf("FillChar mutated: %c != '▓'", bar.FillChar())
	}
	if bar.EmptyChar() != '▒' {
		t.Errorf("EmptyChar mutated: %c != '▒'", bar.EmptyChar())
	}
	if !bar.ShowPercent() {
		t.Errorf("ShowPercent mutated")
	}
	if bar.Label() != "Test" {
		t.Errorf("Label mutated: %s != 'Test'", bar.Label())
	}
}

func TestBarFluentInterface(t *testing.T) {
	bar := NewBar(40).
		WithPercentage(50).
		WithFillChar('▓').
		WithEmptyChar('▒').
		WithShowPercent(true).
		WithLabel("Loading...")

	if bar.Percentage() != 50 {
		t.Errorf("Fluent Percentage() = %d, expected 50", bar.Percentage())
	}
	if bar.FillChar() != '▓' {
		t.Errorf("Fluent FillChar() = %c, expected '▓'", bar.FillChar())
	}
	if bar.EmptyChar() != '▒' {
		t.Errorf("Fluent EmptyChar() = %c, expected '▒'", bar.EmptyChar())
	}
	if !bar.ShowPercent() {
		t.Errorf("Fluent ShowPercent() = false, expected true")
	}
	if bar.Label() != "Loading..." {
		t.Errorf("Fluent Label() = %s, expected 'Loading...'", bar.Label())
	}
}
