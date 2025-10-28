package model

import (
	"strings"
	"testing"
	"time"

	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// TestNewMouseEvent tests the basic constructor.
func TestNewMouseEvent(t *testing.T) {
	eventType := value2.EventClick
	button := value2.ButtonLeft
	position := value2.NewPosition(10, 20)
	modifiers := value2.NewModifiers(true, false, true) // Shift + Alt

	before := time.Now()
	event := NewMouseEvent(eventType, button, position, modifiers)
	after := time.Now()

	// Test all fields are set correctly
	if event.Type() != eventType {
		t.Errorf("Type() = %v, expected %v", event.Type(), eventType)
	}

	if event.Button() != button {
		t.Errorf("Button() = %v, expected %v", event.Button(), button)
	}

	if !event.Position().Equals(position) {
		t.Errorf("Position() = %v, expected %v", event.Position(), position)
	}

	if !event.Modifiers().Equals(modifiers) {
		t.Errorf("Modifiers() = %v, expected %v", event.Modifiers(), modifiers)
	}

	// Test timestamp is set to now (within reasonable bounds)
	timestamp := event.Timestamp()
	if timestamp.Before(before) || timestamp.After(after) {
		t.Errorf("Timestamp() = %v, expected between %v and %v", timestamp, before, after)
	}
}

// TestNewMouseEventWithTimestamp tests the constructor with explicit timestamp.
func TestNewMouseEventWithTimestamp(t *testing.T) {
	eventType := value2.EventDrag
	button := value2.ButtonRight
	position := value2.NewPosition(50, 60)
	modifiers := value2.NewModifiers(false, true, false) // Ctrl only
	timestamp := time.Date(2025, 10, 17, 12, 30, 45, 0, time.UTC)

	event := NewMouseEventWithTimestamp(eventType, button, position, modifiers, timestamp)

	// Test all fields including timestamp
	if event.Type() != eventType {
		t.Errorf("Type() = %v, expected %v", event.Type(), eventType)
	}

	if event.Button() != button {
		t.Errorf("Button() = %v, expected %v", event.Button(), button)
	}

	if !event.Position().Equals(position) {
		t.Errorf("Position() = %v, expected %v", event.Position(), position)
	}

	if !event.Modifiers().Equals(modifiers) {
		t.Errorf("Modifiers() = %v, expected %v", event.Modifiers(), modifiers)
	}

	if !event.Timestamp().Equal(timestamp) {
		t.Errorf("Timestamp() = %v, expected %v", event.Timestamp(), timestamp)
	}
}

// TestMouseEventGetters tests all getter methods.
func TestMouseEventGetters(t *testing.T) {
	tests := []struct {
		name      string
		eventType value2.EventType
		button    value2.Button
		position  value2.Position
		modifiers value2.Modifiers
	}{
		{
			name:      "click event",
			eventType: value2.EventClick,
			button:    value2.ButtonLeft,
			position:  value2.NewPosition(10, 20),
			modifiers: value2.ModifierNone,
		},
		{
			name:      "double-click with modifiers",
			eventType: value2.EventDoubleClick,
			button:    value2.ButtonLeft,
			position:  value2.NewPosition(5, 15),
			modifiers: value2.NewModifiers(true, true, false),
		},
		{
			name:      "right-click with alt",
			eventType: value2.EventClick,
			button:    value2.ButtonRight,
			position:  value2.NewPosition(100, 50),
			modifiers: value2.NewModifiers(false, false, true),
		},
		{
			name:      "drag event",
			eventType: value2.EventDrag,
			button:    value2.ButtonLeft,
			position:  value2.NewPosition(25, 35),
			modifiers: value2.NewModifiers(true, false, false),
		},
		{
			name:      "motion event",
			eventType: value2.EventMotion,
			button:    value2.ButtonNone,
			position:  value2.NewPosition(0, 0),
			modifiers: value2.ModifierNone,
		},
		{
			name:      "scroll event",
			eventType: value2.EventScroll,
			button:    value2.ButtonWheelUp,
			position:  value2.NewPosition(50, 50),
			modifiers: value2.ModifierNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, tt.button, tt.position, tt.modifiers)

			if event.Type() != tt.eventType {
				t.Errorf("Type() = %v, expected %v", event.Type(), tt.eventType)
			}

			if event.Button() != tt.button {
				t.Errorf("Button() = %v, expected %v", event.Button(), tt.button)
			}

			if !event.Position().Equals(tt.position) {
				t.Errorf("Position() = %v, expected %v", event.Position(), tt.position)
			}

			if !event.Modifiers().Equals(tt.modifiers) {
				t.Errorf("Modifiers() = %v, expected %v", event.Modifiers(), tt.modifiers)
			}
		})
	}
}

// TestMouseEventWithType tests the WithType mutation method.
func TestMouseEventWithType(t *testing.T) {
	// Create original event
	originalType := value2.EventPress
	button := value2.ButtonLeft
	position := value2.NewPosition(10, 20)
	modifiers := value2.NewModifiers(true, false, false)
	timestamp := time.Date(2025, 10, 17, 10, 0, 0, 0, time.UTC)

	original := NewMouseEventWithTimestamp(originalType, button, position, modifiers, timestamp)

	// Create mutated event
	newType := value2.EventRelease
	mutated := original.WithType(newType)

	// Test that type changed
	if mutated.Type() != newType {
		t.Errorf("mutated Type() = %v, expected %v", mutated.Type(), newType)
	}

	// Test that other fields remain unchanged
	if mutated.Button() != button {
		t.Errorf("mutated Button() = %v, expected %v", mutated.Button(), button)
	}

	if !mutated.Position().Equals(position) {
		t.Errorf("mutated Position() = %v, expected %v", mutated.Position(), position)
	}

	if !mutated.Modifiers().Equals(modifiers) {
		t.Errorf("mutated Modifiers() = %v, expected %v", mutated.Modifiers(), modifiers)
	}

	if !mutated.Timestamp().Equal(timestamp) {
		t.Errorf("mutated Timestamp() = %v, expected %v", mutated.Timestamp(), timestamp)
	}

	// Test that original is unchanged (immutability)
	if original.Type() != originalType {
		t.Errorf("original Type() = %v, should remain %v", original.Type(), originalType)
	}
}

// TestMouseEventString tests the String() method.
func TestMouseEventString(t *testing.T) {
	timestamp := time.Date(2025, 10, 17, 14, 30, 45, 123000000, time.UTC)

	tests := []struct {
		name          string
		event         MouseEvent
		expectedParts []string // Parts that should be in the string
	}{
		{
			name: "click event",
			event: NewMouseEventWithTimestamp(
				value2.EventClick,
				value2.ButtonLeft,
				value2.NewPosition(10, 20),
				value2.ModifierNone,
				timestamp,
			),
			expectedParts: []string{
				"MouseEvent",
				"type=Click",
				"button=Left",
				"pos=(10,20)",
				"mods=None",
				"14:30:45.123",
			},
		},
		{
			name: "drag with modifiers",
			event: NewMouseEventWithTimestamp(
				value2.EventDrag,
				value2.ButtonRight,
				value2.NewPosition(50, 60),
				value2.NewModifiers(true, true, false),
				timestamp,
			),
			expectedParts: []string{
				"MouseEvent",
				"type=Drag",
				"button=Right",
				"pos=(50,60)",
				"mods=Shift+Ctrl",
				"14:30:45.123",
			},
		},
		{
			name: "scroll event",
			event: NewMouseEventWithTimestamp(
				value2.EventScroll,
				value2.ButtonWheelUp,
				value2.NewPosition(0, 0),
				value2.ModifierNone,
				timestamp,
			),
			expectedParts: []string{
				"MouseEvent",
				"type=Scroll",
				"button=WheelUp",
				"pos=(0,0)",
				"mods=None",
				"14:30:45.123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.event.String()

			for _, part := range tt.expectedParts {
				if !strings.Contains(str, part) {
					t.Errorf("String() = %q, expected to contain %q", str, part)
				}
			}
		})
	}
}

// TestMouseEventIsAt tests position checking with tolerance.
func TestMouseEventIsAt(t *testing.T) {
	eventPos := value2.NewPosition(10, 20)
	event := NewMouseEvent(value2.EventClick, value2.ButtonLeft, eventPos, value2.ModifierNone)

	tests := []struct {
		name      string
		checkPos  value2.Position
		tolerance int
		expected  bool
	}{
		{
			name:      "exact position, zero tolerance",
			checkPos:  value2.NewPosition(10, 20),
			tolerance: 0,
			expected:  true,
		},
		{
			name:      "exact position, with tolerance",
			checkPos:  value2.NewPosition(10, 20),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "within tolerance (distance 2, tolerance 5)",
			checkPos:  value2.NewPosition(11, 21),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "at tolerance boundary (distance 5, tolerance 5)",
			checkPos:  value2.NewPosition(13, 22),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "outside tolerance (distance 6, tolerance 5)",
			checkPos:  value2.NewPosition(14, 22),
			tolerance: 5,
			expected:  false,
		},
		{
			name:      "far outside tolerance",
			checkPos:  value2.NewPosition(100, 100),
			tolerance: 5,
			expected:  false,
		},
		{
			name:      "different position, zero tolerance",
			checkPos:  value2.NewPosition(11, 20),
			tolerance: 0,
			expected:  false,
		},
		{
			name:      "horizontal offset within tolerance",
			checkPos:  value2.NewPosition(12, 20),
			tolerance: 3,
			expected:  true,
		},
		{
			name:      "vertical offset within tolerance",
			checkPos:  value2.NewPosition(10, 23),
			tolerance: 3,
			expected:  true,
		},
		{
			name:      "negative offset within tolerance",
			checkPos:  value2.NewPosition(8, 19),
			tolerance: 3,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := event.IsAt(tt.checkPos, tt.tolerance)

			if result != tt.expected {
				distance := eventPos.DistanceTo(tt.checkPos)
				t.Errorf("IsAt(%v, %d) = %v, expected %v (distance=%d)",
					tt.checkPos, tt.tolerance, result, tt.expected, distance)
			}
		})
	}
}

// TestMouseEventIsClick tests the IsClick method.
func TestMouseEventIsClick(t *testing.T) {
	tests := []struct {
		name      string
		eventType value2.EventType
		expected  bool
	}{
		{"EventClick", value2.EventClick, true},
		{"EventDoubleClick", value2.EventDoubleClick, true},
		{"EventTripleClick", value2.EventTripleClick, true},
		{"EventPress", value2.EventPress, false},
		{"EventRelease", value2.EventRelease, false},
		{"EventDrag", value2.EventDrag, false},
		{"EventMotion", value2.EventMotion, false},
		{"EventScroll", value2.EventScroll, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)

			if event.IsClick() != tt.expected {
				t.Errorf("IsClick() = %v, expected %v for event type %v",
					event.IsClick(), tt.expected, tt.eventType)
			}
		})
	}
}

// TestMouseEventIsDrag tests the IsDrag method.
func TestMouseEventIsDrag(t *testing.T) {
	tests := []struct {
		name      string
		eventType value2.EventType
		expected  bool
	}{
		{"EventDrag", value2.EventDrag, true},
		{"EventClick", value2.EventClick, false},
		{"EventDoubleClick", value2.EventDoubleClick, false},
		{"EventTripleClick", value2.EventTripleClick, false},
		{"EventPress", value2.EventPress, false},
		{"EventRelease", value2.EventRelease, false},
		{"EventMotion", value2.EventMotion, false},
		{"EventScroll", value2.EventScroll, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)

			if event.IsDrag() != tt.expected {
				t.Errorf("IsDrag() = %v, expected %v for event type %v",
					event.IsDrag(), tt.expected, tt.eventType)
			}
		})
	}
}

// TestMouseEventIsScroll tests the IsScroll method.
func TestMouseEventIsScroll(t *testing.T) {
	tests := []struct {
		name      string
		eventType value2.EventType
		expected  bool
	}{
		{"EventScroll", value2.EventScroll, true},
		{"EventClick", value2.EventClick, false},
		{"EventDoubleClick", value2.EventDoubleClick, false},
		{"EventTripleClick", value2.EventTripleClick, false},
		{"EventPress", value2.EventPress, false},
		{"EventRelease", value2.EventRelease, false},
		{"EventDrag", value2.EventDrag, false},
		{"EventMotion", value2.EventMotion, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)

			if event.IsScroll() != tt.expected {
				t.Errorf("IsScroll() = %v, expected %v for event type %v",
					event.IsScroll(), tt.expected, tt.eventType)
			}
		})
	}
}

// TestMouseEventAllEventTypes tests events with all event types.
func TestMouseEventAllEventTypes(t *testing.T) {
	eventTypes := []value2.EventType{
		value2.EventPress,
		value2.EventRelease,
		value2.EventClick,
		value2.EventDoubleClick,
		value2.EventTripleClick,
		value2.EventDrag,
		value2.EventMotion,
		value2.EventScroll,
	}

	for _, eventType := range eventTypes {
		t.Run(eventType.String(), func(t *testing.T) {
			event := NewMouseEvent(
				eventType,
				value2.ButtonLeft,
				value2.NewPosition(10, 20),
				value2.ModifierNone,
			)

			if event.Type() != eventType {
				t.Errorf("Type() = %v, expected %v", event.Type(), eventType)
			}
		})
	}
}

// TestMouseEventAllButtons tests events with all button types.
func TestMouseEventAllButtons(t *testing.T) {
	buttons := []value2.Button{
		value2.ButtonNone,
		value2.ButtonLeft,
		value2.ButtonMiddle,
		value2.ButtonRight,
		value2.ButtonWheelUp,
		value2.ButtonWheelDown,
	}

	for _, button := range buttons {
		t.Run(button.String(), func(t *testing.T) {
			event := NewMouseEvent(
				value2.EventClick,
				button,
				value2.NewPosition(10, 20),
				value2.ModifierNone,
			)

			if event.Button() != button {
				t.Errorf("Button() = %v, expected %v", event.Button(), button)
			}
		})
	}
}

// TestMouseEventAllModifiers tests events with all modifier combinations.
func TestMouseEventAllModifiers(t *testing.T) {
	tests := []struct {
		name      string
		modifiers value2.Modifiers
	}{
		{"no modifiers", value2.ModifierNone},
		{"shift only", value2.NewModifiers(true, false, false)},
		{"ctrl only", value2.NewModifiers(false, true, false)},
		{"alt only", value2.NewModifiers(false, false, true)},
		{"shift+ctrl", value2.NewModifiers(true, true, false)},
		{"shift+alt", value2.NewModifiers(true, false, true)},
		{"ctrl+alt", value2.NewModifiers(false, true, true)},
		{"all modifiers", value2.NewModifiers(true, true, true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(
				value2.EventClick,
				value2.ButtonLeft,
				value2.NewPosition(10, 20),
				tt.modifiers,
			)

			if !event.Modifiers().Equals(tt.modifiers) {
				t.Errorf("Modifiers() = %v, expected %v", event.Modifiers(), tt.modifiers)
			}
		})
	}
}

// TestMouseEventPositionVariations tests events at various positions.
func TestMouseEventPositionVariations(t *testing.T) {
	positions := []value2.Position{
		value2.NewPosition(0, 0),      // top-left
		value2.NewPosition(100, 0),    // top-right area
		value2.NewPosition(0, 50),     // left side
		value2.NewPosition(100, 50),   // center area
		value2.NewPosition(1000, 500), // large coordinates
		value2.NewPosition(-10, -20),  // negative (edge case, might be invalid in practice)
	}

	for _, pos := range positions {
		t.Run(pos.String(), func(t *testing.T) {
			event := NewMouseEvent(
				value2.EventClick,
				value2.ButtonLeft,
				pos,
				value2.ModifierNone,
			)

			if !event.Position().Equals(pos) {
				t.Errorf("Position() = %v, expected %v", event.Position(), pos)
			}
		})
	}
}

// TestMouseEventImmutability tests that WithType doesn't mutate the original.
func TestMouseEventImmutability(t *testing.T) {
	original := NewMouseEvent(
		value2.EventPress,
		value2.ButtonLeft,
		value2.NewPosition(10, 20),
		value2.NewModifiers(true, false, false),
	)

	// Store original values
	originalType := original.Type()
	originalButton := original.Button()
	originalPos := original.Position()
	originalMods := original.Modifiers()
	originalTime := original.Timestamp()

	// Create multiple mutations
	_ = original.WithType(value2.EventRelease)
	_ = original.WithType(value2.EventClick)
	_ = original.WithType(value2.EventDrag)

	// Verify original is unchanged
	if original.Type() != originalType {
		t.Errorf("original Type mutated: got %v, expected %v", original.Type(), originalType)
	}

	if original.Button() != originalButton {
		t.Errorf("original Button mutated: got %v, expected %v", original.Button(), originalButton)
	}

	if !original.Position().Equals(originalPos) {
		t.Errorf("original Position mutated: got %v, expected %v", original.Position(), originalPos)
	}

	if !original.Modifiers().Equals(originalMods) {
		t.Errorf("original Modifiers mutated: got %v, expected %v", original.Modifiers(), originalMods)
	}

	if !original.Timestamp().Equal(originalTime) {
		t.Errorf("original Timestamp mutated: got %v, expected %v", original.Timestamp(), originalTime)
	}
}

// TestMouseEventTimestampOrdering tests that timestamps work correctly for ordering.
func TestMouseEventTimestampOrdering(t *testing.T) {
	// Create events in sequence
	event1 := NewMouseEvent(value2.EventPress, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)
	time.Sleep(2 * time.Millisecond) // Small delay to ensure different timestamps
	event2 := NewMouseEvent(value2.EventRelease, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)
	time.Sleep(2 * time.Millisecond)
	event3 := NewMouseEvent(value2.EventClick, value2.ButtonLeft, value2.NewPosition(0, 0), value2.ModifierNone)

	// Verify ordering
	if !event1.Timestamp().Before(event2.Timestamp()) {
		t.Errorf("event1 timestamp (%v) should be before event2 timestamp (%v)",
			event1.Timestamp(), event2.Timestamp())
	}

	if !event2.Timestamp().Before(event3.Timestamp()) {
		t.Errorf("event2 timestamp (%v) should be before event3 timestamp (%v)",
			event2.Timestamp(), event3.Timestamp())
	}

	if !event1.Timestamp().Before(event3.Timestamp()) {
		t.Errorf("event1 timestamp (%v) should be before event3 timestamp (%v)",
			event1.Timestamp(), event3.Timestamp())
	}
}

// TestMouseEventIsAtZeroTolerance tests IsAt with zero tolerance (exact match only).
func TestMouseEventIsAtZeroTolerance(t *testing.T) {
	eventPos := value2.NewPosition(50, 50)
	event := NewMouseEvent(value2.EventClick, value2.ButtonLeft, eventPos, value2.ModifierNone)

	// Exact match
	if !event.IsAt(value2.NewPosition(50, 50), 0) {
		t.Error("IsAt should return true for exact position with zero tolerance")
	}

	// Off by one should fail with zero tolerance
	if event.IsAt(value2.NewPosition(51, 50), 0) {
		t.Error("IsAt should return false for different position with zero tolerance")
	}

	if event.IsAt(value2.NewPosition(50, 51), 0) {
		t.Error("IsAt should return false for different position with zero tolerance")
	}
}

// TestMouseEventWithTypeChaining tests chaining multiple WithType calls.
func TestMouseEventWithTypeChaining(t *testing.T) {
	event := NewMouseEvent(
		value2.EventPress,
		value2.ButtonLeft,
		value2.NewPosition(10, 20),
		value2.ModifierNone,
	)

	// Chain multiple transformations
	transformed := event.
		WithType(value2.EventRelease).
		WithType(value2.EventClick).
		WithType(value2.EventDoubleClick)

	if transformed.Type() != value2.EventDoubleClick {
		t.Errorf("final Type() = %v, expected EventDoubleClick", transformed.Type())
	}

	// Original should be unchanged
	if event.Type() != value2.EventPress {
		t.Errorf("original Type() = %v, should remain EventPress", event.Type())
	}
}
