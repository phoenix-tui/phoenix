package model

import (
	"strings"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// TestNewMouseEvent tests the basic constructor.
func TestNewMouseEvent(t *testing.T) {
	eventType := value.EventClick
	button := value.ButtonLeft
	position := value.NewPosition(10, 20)
	modifiers := value.NewModifiers(true, false, true) // Shift + Alt

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
	eventType := value.EventDrag
	button := value.ButtonRight
	position := value.NewPosition(50, 60)
	modifiers := value.NewModifiers(false, true, false) // Ctrl only
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
		eventType value.EventType
		button    value.Button
		position  value.Position
		modifiers value.Modifiers
	}{
		{
			name:      "click event",
			eventType: value.EventClick,
			button:    value.ButtonLeft,
			position:  value.NewPosition(10, 20),
			modifiers: value.ModifierNone,
		},
		{
			name:      "double-click with modifiers",
			eventType: value.EventDoubleClick,
			button:    value.ButtonLeft,
			position:  value.NewPosition(5, 15),
			modifiers: value.NewModifiers(true, true, false),
		},
		{
			name:      "right-click with alt",
			eventType: value.EventClick,
			button:    value.ButtonRight,
			position:  value.NewPosition(100, 50),
			modifiers: value.NewModifiers(false, false, true),
		},
		{
			name:      "drag event",
			eventType: value.EventDrag,
			button:    value.ButtonLeft,
			position:  value.NewPosition(25, 35),
			modifiers: value.NewModifiers(true, false, false),
		},
		{
			name:      "motion event",
			eventType: value.EventMotion,
			button:    value.ButtonNone,
			position:  value.NewPosition(0, 0),
			modifiers: value.ModifierNone,
		},
		{
			name:      "scroll event",
			eventType: value.EventScroll,
			button:    value.ButtonWheelUp,
			position:  value.NewPosition(50, 50),
			modifiers: value.ModifierNone,
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
	originalType := value.EventPress
	button := value.ButtonLeft
	position := value.NewPosition(10, 20)
	modifiers := value.NewModifiers(true, false, false)
	timestamp := time.Date(2025, 10, 17, 10, 0, 0, 0, time.UTC)

	original := NewMouseEventWithTimestamp(originalType, button, position, modifiers, timestamp)

	// Create mutated event
	newType := value.EventRelease
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
				value.EventClick,
				value.ButtonLeft,
				value.NewPosition(10, 20),
				value.ModifierNone,
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
				value.EventDrag,
				value.ButtonRight,
				value.NewPosition(50, 60),
				value.NewModifiers(true, true, false),
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
				value.EventScroll,
				value.ButtonWheelUp,
				value.NewPosition(0, 0),
				value.ModifierNone,
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
	eventPos := value.NewPosition(10, 20)
	event := NewMouseEvent(value.EventClick, value.ButtonLeft, eventPos, value.ModifierNone)

	tests := []struct {
		name      string
		checkPos  value.Position
		tolerance int
		expected  bool
	}{
		{
			name:      "exact position, zero tolerance",
			checkPos:  value.NewPosition(10, 20),
			tolerance: 0,
			expected:  true,
		},
		{
			name:      "exact position, with tolerance",
			checkPos:  value.NewPosition(10, 20),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "within tolerance (distance 2, tolerance 5)",
			checkPos:  value.NewPosition(11, 21),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "at tolerance boundary (distance 5, tolerance 5)",
			checkPos:  value.NewPosition(13, 22),
			tolerance: 5,
			expected:  true,
		},
		{
			name:      "outside tolerance (distance 6, tolerance 5)",
			checkPos:  value.NewPosition(14, 22),
			tolerance: 5,
			expected:  false,
		},
		{
			name:      "far outside tolerance",
			checkPos:  value.NewPosition(100, 100),
			tolerance: 5,
			expected:  false,
		},
		{
			name:      "different position, zero tolerance",
			checkPos:  value.NewPosition(11, 20),
			tolerance: 0,
			expected:  false,
		},
		{
			name:      "horizontal offset within tolerance",
			checkPos:  value.NewPosition(12, 20),
			tolerance: 3,
			expected:  true,
		},
		{
			name:      "vertical offset within tolerance",
			checkPos:  value.NewPosition(10, 23),
			tolerance: 3,
			expected:  true,
		},
		{
			name:      "negative offset within tolerance",
			checkPos:  value.NewPosition(8, 19),
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
		eventType value.EventType
		expected  bool
	}{
		{"EventClick", value.EventClick, true},
		{"EventDoubleClick", value.EventDoubleClick, true},
		{"EventTripleClick", value.EventTripleClick, true},
		{"EventPress", value.EventPress, false},
		{"EventRelease", value.EventRelease, false},
		{"EventDrag", value.EventDrag, false},
		{"EventMotion", value.EventMotion, false},
		{"EventScroll", value.EventScroll, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)

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
		eventType value.EventType
		expected  bool
	}{
		{"EventDrag", value.EventDrag, true},
		{"EventClick", value.EventClick, false},
		{"EventDoubleClick", value.EventDoubleClick, false},
		{"EventTripleClick", value.EventTripleClick, false},
		{"EventPress", value.EventPress, false},
		{"EventRelease", value.EventRelease, false},
		{"EventMotion", value.EventMotion, false},
		{"EventScroll", value.EventScroll, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)

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
		eventType value.EventType
		expected  bool
	}{
		{"EventScroll", value.EventScroll, true},
		{"EventClick", value.EventClick, false},
		{"EventDoubleClick", value.EventDoubleClick, false},
		{"EventTripleClick", value.EventTripleClick, false},
		{"EventPress", value.EventPress, false},
		{"EventRelease", value.EventRelease, false},
		{"EventDrag", value.EventDrag, false},
		{"EventMotion", value.EventMotion, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(tt.eventType, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)

			if event.IsScroll() != tt.expected {
				t.Errorf("IsScroll() = %v, expected %v for event type %v",
					event.IsScroll(), tt.expected, tt.eventType)
			}
		})
	}
}

// TestMouseEventAllEventTypes tests events with all event types.
func TestMouseEventAllEventTypes(t *testing.T) {
	eventTypes := []value.EventType{
		value.EventPress,
		value.EventRelease,
		value.EventClick,
		value.EventDoubleClick,
		value.EventTripleClick,
		value.EventDrag,
		value.EventMotion,
		value.EventScroll,
	}

	for _, eventType := range eventTypes {
		t.Run(eventType.String(), func(t *testing.T) {
			event := NewMouseEvent(
				eventType,
				value.ButtonLeft,
				value.NewPosition(10, 20),
				value.ModifierNone,
			)

			if event.Type() != eventType {
				t.Errorf("Type() = %v, expected %v", event.Type(), eventType)
			}
		})
	}
}

// TestMouseEventAllButtons tests events with all button types.
func TestMouseEventAllButtons(t *testing.T) {
	buttons := []value.Button{
		value.ButtonNone,
		value.ButtonLeft,
		value.ButtonMiddle,
		value.ButtonRight,
		value.ButtonWheelUp,
		value.ButtonWheelDown,
	}

	for _, button := range buttons {
		t.Run(button.String(), func(t *testing.T) {
			event := NewMouseEvent(
				value.EventClick,
				button,
				value.NewPosition(10, 20),
				value.ModifierNone,
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
		modifiers value.Modifiers
	}{
		{"no modifiers", value.ModifierNone},
		{"shift only", value.NewModifiers(true, false, false)},
		{"ctrl only", value.NewModifiers(false, true, false)},
		{"alt only", value.NewModifiers(false, false, true)},
		{"shift+ctrl", value.NewModifiers(true, true, false)},
		{"shift+alt", value.NewModifiers(true, false, true)},
		{"ctrl+alt", value.NewModifiers(false, true, true)},
		{"all modifiers", value.NewModifiers(true, true, true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMouseEvent(
				value.EventClick,
				value.ButtonLeft,
				value.NewPosition(10, 20),
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
	positions := []value.Position{
		value.NewPosition(0, 0),      // top-left
		value.NewPosition(100, 0),    // top-right area
		value.NewPosition(0, 50),     // left side
		value.NewPosition(100, 50),   // center area
		value.NewPosition(1000, 500), // large coordinates
		value.NewPosition(-10, -20),  // negative (edge case, might be invalid in practice)
	}

	for _, pos := range positions {
		t.Run(pos.String(), func(t *testing.T) {
			event := NewMouseEvent(
				value.EventClick,
				value.ButtonLeft,
				pos,
				value.ModifierNone,
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
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 20),
		value.NewModifiers(true, false, false),
	)

	// Store original values
	originalType := original.Type()
	originalButton := original.Button()
	originalPos := original.Position()
	originalMods := original.Modifiers()
	originalTime := original.Timestamp()

	// Create multiple mutations
	_ = original.WithType(value.EventRelease)
	_ = original.WithType(value.EventClick)
	_ = original.WithType(value.EventDrag)

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
	event1 := NewMouseEvent(value.EventPress, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)
	time.Sleep(2 * time.Millisecond) // Small delay to ensure different timestamps
	event2 := NewMouseEvent(value.EventRelease, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)
	time.Sleep(2 * time.Millisecond)
	event3 := NewMouseEvent(value.EventClick, value.ButtonLeft, value.NewPosition(0, 0), value.ModifierNone)

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
	eventPos := value.NewPosition(50, 50)
	event := NewMouseEvent(value.EventClick, value.ButtonLeft, eventPos, value.ModifierNone)

	// Exact match
	if !event.IsAt(value.NewPosition(50, 50), 0) {
		t.Error("IsAt should return true for exact position with zero tolerance")
	}

	// Off by one should fail with zero tolerance
	if event.IsAt(value.NewPosition(51, 50), 0) {
		t.Error("IsAt should return false for different position with zero tolerance")
	}

	if event.IsAt(value.NewPosition(50, 51), 0) {
		t.Error("IsAt should return false for different position with zero tolerance")
	}
}

// TestMouseEventWithTypeChaining tests chaining multiple WithType calls.
func TestMouseEventWithTypeChaining(t *testing.T) {
	event := NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 20),
		value.ModifierNone,
	)

	// Chain multiple transformations
	transformed := event.
		WithType(value.EventRelease).
		WithType(value.EventClick).
		WithType(value.EventDoubleClick)

	if transformed.Type() != value.EventDoubleClick {
		t.Errorf("final Type() = %v, expected EventDoubleClick", transformed.Type())
	}

	// Original should be unchanged
	if event.Type() != value.EventPress {
		t.Errorf("original Type() = %v, should remain EventPress", event.Type())
	}
}
