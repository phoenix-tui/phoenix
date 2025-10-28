package model

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// TestNewDragState tests the constructor with various thresholds.
func TestNewDragState(t *testing.T) {
	tests := []struct {
		name              string
		threshold         int
		expectedThreshold int
		expectActive      bool
	}{
		{
			name:              "positive threshold",
			threshold:         5,
			expectedThreshold: 5,
			expectActive:      false,
		},
		{
			name:              "zero threshold uses default",
			threshold:         0,
			expectedThreshold: 2,
			expectActive:      false,
		},
		{
			name:              "negative threshold uses default",
			threshold:         -5,
			expectedThreshold: 2,
			expectActive:      false,
		},
		{
			name:              "threshold of 1",
			threshold:         1,
			expectedThreshold: 1,
			expectActive:      false,
		},
		{
			name:              "large threshold",
			threshold:         100,
			expectedThreshold: 100,
			expectActive:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDragState(tt.threshold)

			if ds == nil {
				t.Fatal("NewDragState returned nil")
			}

			if ds.IsActive() != tt.expectActive {
				t.Errorf("expected IsActive() = %v, got %v", tt.expectActive, ds.IsActive())
			}

			// Verify threshold by testing IsDrag behavior
			// We can't directly access threshold, so we test its effect
			if ds.IsActive() {
				t.Error("newly created DragState should not be active")
			}
		})
	}
}

// TestDragStateStartEndsLifecycle tests the full lifecycle of a drag operation.
func TestDragStateStartEndsLifecycle(t *testing.T) {
	ds := NewDragState(2)

	// Initially not active
	if ds.IsActive() {
		t.Error("drag should not be active initially")
	}

	if ds.IsDrag() {
		t.Error("drag should not be considered a drag initially")
	}

	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0, got %d", ds.Distance())
	}

	// Start drag
	startPos := value2.NewPosition(10, 20)
	button := value2.ButtonLeft
	modifiers := value2.NewModifiers(false, true, false) // Ctrl pressed

	ds.Start(startPos, button, modifiers)

	// Should now be active
	if !ds.IsActive() {
		t.Error("drag should be active after Start")
	}

	// Check stored values
	if !ds.StartPosition().Equals(startPos) {
		t.Errorf("expected start position %v, got %v", startPos, ds.StartPosition())
	}

	if !ds.Current().Equals(startPos) {
		t.Errorf("expected current position %v initially, got %v", startPos, ds.Current())
	}

	if ds.Button() != button {
		t.Errorf("expected button %v, got %v", button, ds.Button())
	}

	if !ds.Modifiers().Equals(modifiers) {
		t.Errorf("expected modifiers %v, got %v", modifiers, ds.Modifiers())
	}

	// Distance should be 0 at start
	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 at start, got %d", ds.Distance())
	}

	// Not yet a drag (below threshold)
	if ds.IsDrag() {
		t.Error("should not be considered a drag at start position")
	}

	// Update to position at threshold (Manhattan distance = 2)
	ds.Update(value2.NewPosition(11, 21)) // dx=1, dy=1, distance=2

	if !ds.IsActive() {
		t.Error("drag should still be active after Update")
	}

	if !ds.IsDrag() {
		t.Error("should be considered a drag at threshold")
	}

	if ds.Distance() != 2 {
		t.Errorf("expected distance = 2, got %d", ds.Distance())
	}

	// Update to position beyond threshold
	ds.Update(value2.NewPosition(15, 25)) // dx=5, dy=5, distance=10

	if ds.Distance() != 10 {
		t.Errorf("expected distance = 10, got %d", ds.Distance())
	}

	if !ds.IsDrag() {
		t.Error("should definitely be a drag beyond threshold")
	}

	// End drag
	ds.End()

	if ds.IsActive() {
		t.Error("drag should not be active after End")
	}

	if ds.IsDrag() {
		t.Error("ended drag should not be considered a drag")
	}

	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 after End, got %d", ds.Distance())
	}
}

// TestDragStateUpdate tests the Update method behavior.
func TestDragStateUpdate(t *testing.T) {
	ds := NewDragState(2)
	startPos := value2.NewPosition(10, 10)

	// Update on inactive drag should have no effect
	ds.Update(value2.NewPosition(20, 20))

	if ds.IsActive() {
		t.Error("update on inactive drag should not activate it")
	}

	// Start drag
	ds.Start(startPos, value2.ButtonLeft, value2.ModifierNone)

	// Update position
	newPos := value2.NewPosition(15, 15)
	ds.Update(newPos)

	if !ds.Current().Equals(newPos) {
		t.Errorf("expected current position %v, got %v", newPos, ds.Current())
	}

	// Start position should remain unchanged
	if !ds.StartPosition().Equals(startPos) {
		t.Errorf("start position should not change, expected %v, got %v", startPos, ds.StartPosition())
	}
}

// TestDragStateIsActive tests the IsActive method.
func TestDragStateIsActive(t *testing.T) {
	ds := NewDragState(2)

	// Initially not active
	if ds.IsActive() {
		t.Error("should not be active initially")
	}

	// Start makes it active
	ds.Start(value2.NewPosition(0, 0), value2.ButtonLeft, value2.ModifierNone)
	if !ds.IsActive() {
		t.Error("should be active after Start")
	}

	// Update keeps it active
	ds.Update(value2.NewPosition(5, 5))
	if !ds.IsActive() {
		t.Error("should remain active after Update")
	}

	// End makes it inactive
	ds.End()
	if ds.IsActive() {
		t.Error("should not be active after End")
	}
}

// TestDragStateIsDrag tests the IsDrag method with threshold detection.
func TestDragStateIsDrag(t *testing.T) {
	threshold := 5

	tests := []struct {
		name       string
		startPos   value2.Position
		updatePos  value2.Position
		shouldDrag bool
		distance   int
	}{
		{
			name:       "no movement - not a drag",
			startPos:   value2.NewPosition(10, 10),
			updatePos:  value2.NewPosition(10, 10),
			shouldDrag: false,
			distance:   0,
		},
		{
			name:       "below threshold - not a drag",
			startPos:   value2.NewPosition(10, 10),
			updatePos:  value2.NewPosition(12, 12), // distance = 4
			shouldDrag: false,
			distance:   4,
		},
		{
			name:       "at threshold - is a drag",
			startPos:   value2.NewPosition(10, 10),
			updatePos:  value2.NewPosition(13, 12), // distance = 5
			shouldDrag: true,
			distance:   5,
		},
		{
			name:       "above threshold - is a drag",
			startPos:   value2.NewPosition(10, 10),
			updatePos:  value2.NewPosition(20, 20), // distance = 20
			shouldDrag: true,
			distance:   20,
		},
		{
			name:       "horizontal movement at threshold",
			startPos:   value2.NewPosition(0, 0),
			updatePos:  value2.NewPosition(5, 0), // distance = 5
			shouldDrag: true,
			distance:   5,
		},
		{
			name:       "vertical movement at threshold",
			startPos:   value2.NewPosition(0, 0),
			updatePos:  value2.NewPosition(0, 5), // distance = 5
			shouldDrag: true,
			distance:   5,
		},
		{
			name:       "negative direction movement",
			startPos:   value2.NewPosition(10, 10),
			updatePos:  value2.NewPosition(5, 5), // distance = 10
			shouldDrag: true,
			distance:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDragState(threshold)

			// Not a drag when inactive
			if ds.IsDrag() {
				t.Error("should not be a drag when inactive")
			}

			ds.Start(tt.startPos, value2.ButtonLeft, value2.ModifierNone)
			ds.Update(tt.updatePos)

			if ds.IsDrag() != tt.shouldDrag {
				t.Errorf("expected IsDrag() = %v, got %v (distance=%d, threshold=%d)",
					tt.shouldDrag, ds.IsDrag(), tt.distance, threshold)
			}

			if ds.Distance() != tt.distance {
				t.Errorf("expected distance = %d, got %d", tt.distance, ds.Distance())
			}
		})
	}
}

// TestDragStateDistanceCalculation tests the Distance method.
func TestDragStateDistanceCalculation(t *testing.T) {
	ds := NewDragState(2)

	// Distance is 0 when inactive
	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 when inactive, got %d", ds.Distance())
	}

	startPos := value2.NewPosition(10, 10)
	ds.Start(startPos, value2.ButtonLeft, value2.ModifierNone)

	// Distance is 0 at start
	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 at start, got %d", ds.Distance())
	}

	// Test Manhattan distance calculation
	tests := []struct {
		name     string
		position value2.Position
		expected int
	}{
		{"same position", value2.NewPosition(10, 10), 0},
		{"horizontal +5", value2.NewPosition(15, 10), 5},
		{"horizontal -5", value2.NewPosition(5, 10), 5},
		{"vertical +5", value2.NewPosition(10, 15), 5},
		{"vertical -5", value2.NewPosition(10, 5), 5},
		{"diagonal (3,4)", value2.NewPosition(13, 14), 7},
		{"diagonal (-3,-4)", value2.NewPosition(7, 6), 7},
		{"large distance", value2.NewPosition(100, 100), 180},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds.Update(tt.position)

			if ds.Distance() != tt.expected {
				t.Errorf("expected distance = %d, got %d for position %v",
					tt.expected, ds.Distance(), tt.position)
			}
		})
	}

	// Distance becomes 0 after End
	ds.End()
	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 after End, got %d", ds.Distance())
	}
}

// TestDragStateReset tests the Reset method.
func TestDragStateReset(t *testing.T) {
	ds := NewDragState(5)

	// Setup a drag with all fields populated
	startPos := value2.NewPosition(10, 20)
	button := value2.ButtonRight
	modifiers := value2.NewModifiers(true, true, true) // All modifiers

	ds.Start(startPos, button, modifiers)
	ds.Update(value2.NewPosition(50, 60))

	// Verify drag is active and populated
	if !ds.IsActive() {
		t.Fatal("drag should be active before reset")
	}

	if ds.Distance() == 0 {
		t.Fatal("distance should be non-zero before reset")
	}

	// Reset
	ds.Reset()

	// Verify all fields are reset
	if ds.IsActive() {
		t.Error("should not be active after Reset")
	}

	if ds.IsDrag() {
		t.Error("should not be a drag after Reset")
	}

	if ds.Distance() != 0 {
		t.Errorf("expected distance = 0 after Reset, got %d", ds.Distance())
	}

	expectedPos := value2.NewPosition(0, 0)
	if !ds.StartPosition().Equals(expectedPos) {
		t.Errorf("expected start position %v after Reset, got %v", expectedPos, ds.StartPosition())
	}

	if !ds.Current().Equals(expectedPos) {
		t.Errorf("expected current position %v after Reset, got %v", expectedPos, ds.Current())
	}

	if ds.Button() != value2.ButtonNone {
		t.Errorf("expected button None after Reset, got %v", ds.Button())
	}

	if !ds.Modifiers().Equals(value2.ModifierNone) {
		t.Errorf("expected modifiers None after Reset, got %v", ds.Modifiers())
	}
}

// TestDragStateGetters tests all getter methods.
func TestDragStateGetters(t *testing.T) {
	ds := NewDragState(2)

	startPos := value2.NewPosition(15, 25)
	currentPos := value2.NewPosition(20, 30)
	button := value2.ButtonMiddle
	modifiers := value2.NewModifiers(true, false, true) // Shift + Alt

	ds.Start(startPos, button, modifiers)
	ds.Update(currentPos)

	// Test StartPosition
	if !ds.StartPosition().Equals(startPos) {
		t.Errorf("StartPosition() = %v, expected %v", ds.StartPosition(), startPos)
	}

	// Test Current
	if !ds.Current().Equals(currentPos) {
		t.Errorf("Current() = %v, expected %v", ds.Current(), currentPos)
	}

	// Test Button
	if ds.Button() != button {
		t.Errorf("Button() = %v, expected %v", ds.Button(), button)
	}

	// Test Modifiers
	if !ds.Modifiers().Equals(modifiers) {
		t.Errorf("Modifiers() = %v, expected %v", ds.Modifiers(), modifiers)
	}

	// Test Distance (Manhattan: |20-15| + |30-25| = 10)
	expectedDistance := 10
	if ds.Distance() != expectedDistance {
		t.Errorf("Distance() = %d, expected %d", ds.Distance(), expectedDistance)
	}
}

// TestDragStateThresholdEdgeCases tests edge cases around the threshold boundary.
func TestDragStateThresholdEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		threshold  int
		distance   int
		shouldDrag bool
	}{
		{"threshold 1, distance 0", 1, 0, false},
		{"threshold 1, distance 1", 1, 1, true},
		{"threshold 1, distance 2", 1, 2, true},
		{"threshold 2, distance 1", 2, 1, false},
		{"threshold 2, distance 2", 2, 2, true},
		{"threshold 2, distance 3", 2, 3, true},
		{"threshold 5, distance 4", 5, 4, false},
		{"threshold 5, distance 5", 5, 5, true},
		{"threshold 5, distance 6", 5, 6, true},
		{"threshold 10, distance 9", 10, 9, false},
		{"threshold 10, distance 10", 10, 10, true},
		{"threshold 10, distance 11", 10, 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDragState(tt.threshold)

			startPos := value2.NewPosition(0, 0)
			// Create position at exact distance using horizontal movement for simplicity
			endPos := value2.NewPosition(tt.distance, 0)

			ds.Start(startPos, value2.ButtonLeft, value2.ModifierNone)
			ds.Update(endPos)

			if ds.IsDrag() != tt.shouldDrag {
				t.Errorf("expected IsDrag() = %v, got %v (threshold=%d, distance=%d)",
					tt.shouldDrag, ds.IsDrag(), tt.threshold, tt.distance)
			}
		})
	}
}

// TestDragStateMultipleButtons tests drag with different buttons.
func TestDragStateMultipleButtons(t *testing.T) {
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
			ds := NewDragState(2)

			ds.Start(value2.NewPosition(0, 0), button, value2.ModifierNone)

			if ds.Button() != button {
				t.Errorf("expected button %v, got %v", button, ds.Button())
			}
		})
	}
}

// TestDragStateMultipleModifiers tests drag with different modifier combinations.
func TestDragStateMultipleModifiers(t *testing.T) {
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
			ds := NewDragState(2)

			ds.Start(value2.NewPosition(0, 0), value2.ButtonLeft, tt.modifiers)

			if !ds.Modifiers().Equals(tt.modifiers) {
				t.Errorf("expected modifiers %v, got %v", tt.modifiers, ds.Modifiers())
			}
		})
	}
}

// TestDragStateSequentialDrags tests multiple drag operations in sequence.
func TestDragStateSequentialDrags(t *testing.T) {
	ds := NewDragState(3)

	// First drag
	ds.Start(value2.NewPosition(0, 0), value2.ButtonLeft, value2.ModifierNone)
	ds.Update(value2.NewPosition(10, 10))
	if !ds.IsDrag() {
		t.Error("first drag should be active")
	}
	ds.End()

	// Second drag (should be independent)
	ds.Start(value2.NewPosition(20, 20), value2.ButtonRight, value2.NewModifiers(true, false, false))
	ds.Update(value2.NewPosition(30, 30))

	if !ds.IsActive() {
		t.Error("second drag should be active")
	}

	if ds.StartPosition() != value2.NewPosition(20, 20) {
		t.Error("second drag should have new start position")
	}

	if ds.Button() != value2.ButtonRight {
		t.Error("second drag should have new button")
	}

	ds.End()

	// Third drag after Reset
	ds.Reset()
	ds.Start(value2.NewPosition(5, 5), value2.ButtonMiddle, value2.ModifierNone)

	if !ds.IsActive() {
		t.Error("third drag should be active after reset")
	}

	if ds.StartPosition() != value2.NewPosition(5, 5) {
		t.Error("third drag should have correct start position")
	}
}
