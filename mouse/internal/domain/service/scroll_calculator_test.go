package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// Test basic scroll delta calculation
func TestScrollCalculator_CalculateDelta(t *testing.T) {
	calc := NewScrollCalculator(3) // 3 lines per scroll

	tests := []struct {
		name          string
		button        value2.Button
		expectedDelta int
	}{
		{
			name:          "Scroll up",
			button:        value2.ButtonWheelUp,
			expectedDelta: -3,
		},
		{
			name:          "Scroll down",
			button:        value2.ButtonWheelDown,
			expectedDelta: 3,
		},
		{
			name:          "Left button (not scroll)",
			button:        value2.ButtonLeft,
			expectedDelta: 0,
		},
		{
			name:          "Right button (not scroll)",
			button:        value2.ButtonRight,
			expectedDelta: 0,
		},
		{
			name:          "Middle button (not scroll)",
			button:        value2.ButtonMiddle,
			expectedDelta: 0,
		},
		{
			name:          "None button",
			button:        value2.ButtonNone,
			expectedDelta: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := model.NewMouseEvent(
				value2.EventScroll,
				tt.button,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)

			delta := calc.CalculateDelta(event)
			if delta != tt.expectedDelta {
				t.Errorf("Expected delta %d, got %d", tt.expectedDelta, delta)
			}
		})
	}
}

// Test non-scroll events return 0 delta
func TestScrollCalculator_NonScrollEvents(t *testing.T) {
	calc := NewScrollCalculator(3)

	nonScrollEvents := []value2.EventType{
		value2.EventPress,
		value2.EventRelease,
		value2.EventMotion,
		value2.EventClick,
		value2.EventDoubleClick,
		value2.EventTripleClick,
		value2.EventDrag,
	}

	for _, eventType := range nonScrollEvents {
		t.Run(eventType.String(), func(t *testing.T) {
			event := model.NewMouseEvent(
				eventType,
				value2.ButtonWheelUp, // Even with wheel button
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)

			delta := calc.CalculateDelta(event)
			if delta != 0 {
				t.Errorf("Expected delta 0 for %v, got %d", eventType, delta)
			}
		})
	}
}

// Test IsScrollUp
func TestScrollCalculator_IsScrollUp(t *testing.T) {
	calc := NewScrollCalculator(3)

	tests := []struct {
		name     string
		event    model.MouseEvent
		expected bool
	}{
		{
			name: "Scroll up event",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: true,
		},
		{
			name: "Scroll down event",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
		{
			name: "Non-scroll event with wheel up",
			event: model.NewMouseEvent(
				value2.EventPress,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
		{
			name: "Scroll event with left button",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonLeft,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.IsScrollUp(tt.event)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test IsScrollDown
func TestScrollCalculator_IsScrollDown(t *testing.T) {
	calc := NewScrollCalculator(3)

	tests := []struct {
		name     string
		event    model.MouseEvent
		expected bool
	}{
		{
			name: "Scroll down event",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: true,
		},
		{
			name: "Scroll up event",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
		{
			name: "Non-scroll event with wheel down",
			event: model.NewMouseEvent(
				value2.EventPress,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
		{
			name: "Scroll event with right button",
			event: model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonRight,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.IsScrollDown(tt.event)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test constructor with various lines per scroll values
func TestScrollCalculator_Constructor(t *testing.T) {
	tests := []struct {
		name              string
		linesPerScroll    int
		expectedLines     int
		expectedUpDelta   int
		expectedDownDelta int
	}{
		{
			name:              "Valid positive value",
			linesPerScroll:    5,
			expectedLines:     5,
			expectedUpDelta:   -5,
			expectedDownDelta: 5,
		},
		{
			name:              "Minimum value (1)",
			linesPerScroll:    1,
			expectedLines:     1,
			expectedUpDelta:   -1,
			expectedDownDelta: 1,
		},
		{
			name:              "Large value",
			linesPerScroll:    100,
			expectedLines:     100,
			expectedUpDelta:   -100,
			expectedDownDelta: 100,
		},
		{
			name:              "Zero (should default to 3)",
			linesPerScroll:    0,
			expectedLines:     3,
			expectedUpDelta:   -3,
			expectedDownDelta: 3,
		},
		{
			name:              "Negative (should default to 3)",
			linesPerScroll:    -10,
			expectedLines:     3,
			expectedUpDelta:   -3,
			expectedDownDelta: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewScrollCalculator(tt.linesPerScroll)

			// Check LinesPerScroll getter
			if calc.LinesPerScroll() != tt.expectedLines {
				t.Errorf("Expected lines per scroll %d, got %d", tt.expectedLines, calc.LinesPerScroll())
			}

			// Check scroll up delta
			upEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)
			upDelta := calc.CalculateDelta(upEvent)
			if upDelta != tt.expectedUpDelta {
				t.Errorf("Expected up delta %d, got %d", tt.expectedUpDelta, upDelta)
			}

			// Check scroll down delta
			downEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)
			downDelta := calc.CalculateDelta(downEvent)
			if downDelta != tt.expectedDownDelta {
				t.Errorf("Expected down delta %d, got %d", tt.expectedDownDelta, downDelta)
			}
		})
	}
}

// Test scroll direction consistency
func TestScrollCalculator_DirectionConsistency(t *testing.T) {
	calc := NewScrollCalculator(3)

	// Scroll up: negative delta
	upEvent := model.NewMouseEvent(
		value2.EventScroll,
		value2.ButtonWheelUp,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)

	if !calc.IsScrollUp(upEvent) {
		t.Error("Expected IsScrollUp to be true for scroll up event")
	}

	if calc.IsScrollDown(upEvent) {
		t.Error("Expected IsScrollDown to be false for scroll up event")
	}

	upDelta := calc.CalculateDelta(upEvent)
	if upDelta >= 0 {
		t.Errorf("Expected negative delta for scroll up, got %d", upDelta)
	}

	// Scroll down: positive delta
	downEvent := model.NewMouseEvent(
		value2.EventScroll,
		value2.ButtonWheelDown,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)

	if !calc.IsScrollDown(downEvent) {
		t.Error("Expected IsScrollDown to be true for scroll down event")
	}

	if calc.IsScrollUp(downEvent) {
		t.Error("Expected IsScrollUp to be false for scroll down event")
	}

	downDelta := calc.CalculateDelta(downEvent)
	if downDelta <= 0 {
		t.Errorf("Expected positive delta for scroll down, got %d", downDelta)
	}
}

// Test scroll with modifiers (shouldn't affect delta)
func TestScrollCalculator_WithModifiers(t *testing.T) {
	calc := NewScrollCalculator(3)

	modifiers := []value2.Modifiers{
		value2.ModifierNone,
		value2.ModifierShift,
		value2.ModifierCtrl,
		value2.ModifierAlt,
		value2.ModifierShift | value2.ModifierCtrl,
		value2.ModifierShift | value2.ModifierAlt,
		value2.ModifierCtrl | value2.ModifierAlt,
		value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt,
	}

	for _, mod := range modifiers {
		t.Run(mod.String(), func(t *testing.T) {
			// Scroll up with modifiers
			upEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				mod,
			)
			upDelta := calc.CalculateDelta(upEvent)
			if upDelta != -3 {
				t.Errorf("Expected delta -3 for scroll up with %v, got %d", mod, upDelta)
			}

			// Scroll down with modifiers
			downEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				mod,
			)
			downDelta := calc.CalculateDelta(downEvent)
			if downDelta != 3 {
				t.Errorf("Expected delta 3 for scroll down with %v, got %d", mod, downDelta)
			}
		})
	}
}

// Test scroll at different positions (shouldn't affect delta)
func TestScrollCalculator_DifferentPositions(t *testing.T) {
	calc := NewScrollCalculator(3)

	positions := []value2.Position{
		value2.NewPosition(0, 0),
		value2.NewPosition(10, 10),
		value2.NewPosition(100, 100),
		value2.NewPosition(-5, -5), // Negative positions (edge case)
	}

	for _, pos := range positions {
		t.Run(pos.String(), func(t *testing.T) {
			// Position shouldn't affect scroll delta
			upEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				pos,
				value2.ModifierNone,
			)
			upDelta := calc.CalculateDelta(upEvent)
			if upDelta != -3 {
				t.Errorf("Expected delta -3 at position %v, got %d", pos, upDelta)
			}

			downEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				pos,
				value2.ModifierNone,
			)
			downDelta := calc.CalculateDelta(downEvent)
			if downDelta != 3 {
				t.Errorf("Expected delta 3 at position %v, got %d", pos, downDelta)
			}
		})
	}
}

// Test multiple scroll calculations (stateless)
func TestScrollCalculator_MultipleScrollsStateless(t *testing.T) {
	calc := NewScrollCalculator(3)

	// Scroll multiple times - each should return same delta
	for i := 0; i < 10; i++ {
		upEvent := model.NewMouseEvent(
			value2.EventScroll,
			value2.ButtonWheelUp,
			value2.NewPosition(10, 5),
			value2.ModifierNone,
		)
		upDelta := calc.CalculateDelta(upEvent)
		if upDelta != -3 {
			t.Errorf("Scroll %d: expected delta -3, got %d", i, upDelta)
		}

		downEvent := model.NewMouseEvent(
			value2.EventScroll,
			value2.ButtonWheelDown,
			value2.NewPosition(10, 5),
			value2.ModifierNone,
		)
		downDelta := calc.CalculateDelta(downEvent)
		if downDelta != 3 {
			t.Errorf("Scroll %d: expected delta 3, got %d", i, downDelta)
		}
	}
}

// Test delta symmetry (up and down should be opposite)
func TestScrollCalculator_DeltaSymmetry(t *testing.T) {
	linesPerScrollValues := []int{1, 3, 5, 10, 100}

	for _, lines := range linesPerScrollValues {
		t.Run("lines="+string(rune(lines)), func(t *testing.T) {
			calc := NewScrollCalculator(lines)

			upEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)
			upDelta := calc.CalculateDelta(upEvent)

			downEvent := model.NewMouseEvent(
				value2.EventScroll,
				value2.ButtonWheelDown,
				value2.NewPosition(10, 5),
				value2.ModifierNone,
			)
			downDelta := calc.CalculateDelta(downEvent)

			// Up and down deltas should be opposite
			if upDelta != -downDelta {
				t.Errorf("Expected symmetric deltas, got up=%d, down=%d", upDelta, downDelta)
			}

			// Magnitude should match lines per scroll
			if upDelta != -lines {
				t.Errorf("Expected up delta -%d, got %d", lines, upDelta)
			}

			if downDelta != lines {
				t.Errorf("Expected down delta %d, got %d", lines, downDelta)
			}
		})
	}
}
