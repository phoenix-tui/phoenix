package service

import (
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// Test basic drag flow: press → motion → release
func TestDragTracker_BasicDragFlow(t *testing.T) {
	tracker := NewDragTracker(2) // 2-cell threshold
	startPos := value.NewPosition(10, 10)
	endPos := value.NewPosition(15, 10)

	// Press event - start drag
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierNone,
	)
	tracker.ProcessPress(pressEvent)

	// Should be active now
	if !tracker.IsActive() {
		t.Error("Expected drag to be active after press")
	}

	// Should not be a drag yet (no movement)
	if tracker.IsDrag() {
		t.Error("Expected not to be a drag before threshold reached")
	}

	// Motion event - beyond threshold
	motionEvent := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		endPos,
		value.ModifierNone,
	)
	dragEvent := tracker.ProcessMotion(motionEvent)

	// Should be a drag now
	if !tracker.IsDrag() {
		t.Error("Expected to be a drag after threshold reached")
	}

	// Should return drag event
	if dragEvent == nil {
		t.Fatal("Expected drag event, got nil")
	}

	if dragEvent.Type() != value.EventDrag {
		t.Errorf("Expected EventDrag, got %v", dragEvent.Type())
	}

	// Release event - end drag
	releaseEvent := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		endPos,
		value.ModifierNone,
	)
	wasDrag, start, end := tracker.ProcessRelease(releaseEvent)

	if !wasDrag {
		t.Error("Expected wasDrag to be true")
	}

	if start != startPos {
		t.Errorf("Expected start position %v, got %v", startPos, start)
	}

	if end != endPos {
		t.Errorf("Expected end position %v, got %v", endPos, end)
	}

	// Should not be active anymore
	if tracker.IsActive() {
		t.Error("Expected drag to be inactive after release")
	}
}

// Test threshold detection (exactly 2 cells)
func TestDragTracker_ThresholdDetection(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)

	// Start drag
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierNone,
	)
	tracker.ProcessPress(pressEvent)

	tests := []struct {
		name       string
		pos        value.Position
		expectDrag bool
	}{
		{
			name:       "No movement",
			pos:        value.NewPosition(10, 10),
			expectDrag: false,
		},
		{
			name:       "1 cell distance (below threshold)",
			pos:        value.NewPosition(11, 10),
			expectDrag: false,
		},
		{
			name:       "2 cells distance (at threshold)",
			pos:        value.NewPosition(12, 10),
			expectDrag: true,
		},
		{
			name:       "5 cells distance (above threshold)",
			pos:        value.NewPosition(15, 10),
			expectDrag: true,
		},
		{
			name:       "Vertical movement (2 cells)",
			pos:        value.NewPosition(10, 12),
			expectDrag: true,
		},
		{
			name:       "Diagonal movement (Manhattan=2, at threshold)",
			pos:        value.NewPosition(11, 11),
			expectDrag: true, // Manhattan: |11-10| + |11-10| = 2
		},
		{
			name:       "Diagonal movement (Manhattan=4, above threshold)",
			pos:        value.NewPosition(12, 12),
			expectDrag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and start fresh
			tracker.Reset()
			tracker.ProcessPress(pressEvent)

			// Move to test position
			motionEvent := model.NewMouseEvent(
				value.EventMotion,
				value.ButtonLeft,
				tt.pos,
				value.ModifierNone,
			)
			dragEvent := tracker.ProcessMotion(motionEvent)

			if tt.expectDrag { //nolint:nestif // Test validation requires nested checks
				if !tracker.IsDrag() {
					t.Errorf("Expected IsDrag() = true for position %v", tt.pos)
				}
				if dragEvent == nil {
					t.Errorf("Expected drag event for position %v", tt.pos)
				}
			} else {
				if tracker.IsDrag() {
					t.Errorf("Expected IsDrag() = false for position %v", tt.pos)
				}
				if dragEvent != nil {
					t.Errorf("Expected no drag event for position %v", tt.pos)
				}
			}
		})
	}
}

// Test motion without active drag
func TestDragTracker_MotionWithoutPress(t *testing.T) {
	tracker := NewDragTracker(2)

	// Motion event without press
	motionEvent := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(10, 10),
		value.ModifierNone,
	)
	dragEvent := tracker.ProcessMotion(motionEvent)

	if dragEvent != nil {
		t.Error("Expected nil drag event when no press occurred")
	}

	if tracker.IsActive() {
		t.Error("Expected drag to be inactive")
	}

	if tracker.IsDrag() {
		t.Error("Expected IsDrag() = false when no press occurred")
	}
}

// Test release without active drag
func TestDragTracker_ReleaseWithoutPress(t *testing.T) {
	tracker := NewDragTracker(2)

	// Release event without press
	releaseEvent := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		value.NewPosition(10, 10),
		value.ModifierNone,
	)
	wasDrag, start, end := tracker.ProcessRelease(releaseEvent)

	if wasDrag {
		t.Error("Expected wasDrag = false when no press occurred")
	}

	// Should return zero positions
	if start != value.NewPosition(0, 0) {
		t.Errorf("Expected zero start position, got %v", start)
	}

	if end != value.NewPosition(0, 0) {
		t.Errorf("Expected zero end position, got %v", end)
	}
}

// Test button and modifiers preservation
func TestDragTracker_PreservesButtonAndModifiers(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)
	button := value.ButtonRight
	modifiers := value.ModifierShift | value.ModifierCtrl

	// Press with specific button and modifiers
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		button,
		startPos,
		modifiers,
	)
	tracker.ProcessPress(pressEvent)

	// Check state
	state := tracker.State()
	if state.Button() != button {
		t.Errorf("Expected button %v, got %v", button, state.Button())
	}

	if state.Modifiers() != modifiers {
		t.Errorf("Expected modifiers %v, got %v", modifiers, state.Modifiers())
	}

	// Motion event (trigger drag)
	motionEvent := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft, // Different button (shouldn't matter)
		value.NewPosition(15, 10),
		value.ModifierNone, // Different modifiers (shouldn't matter)
	)
	dragEvent := tracker.ProcessMotion(motionEvent)

	// Drag event should preserve original button and modifiers
	if dragEvent == nil {
		t.Fatal("Expected drag event")
	}

	if dragEvent.Button() != button {
		t.Errorf("Drag event: expected button %v, got %v", button, dragEvent.Button())
	}

	if dragEvent.Modifiers() != modifiers {
		t.Errorf("Drag event: expected modifiers %v, got %v", modifiers, dragEvent.Modifiers())
	}
}

// Test reset functionality
func TestDragTracker_Reset(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)

	// Start drag
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierShift,
	)
	tracker.ProcessPress(pressEvent)

	// Move beyond threshold
	motionEvent := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(15, 10),
		value.ModifierShift,
	)
	tracker.ProcessMotion(motionEvent)

	if !tracker.IsActive() {
		t.Error("Expected drag to be active before reset")
	}

	if !tracker.IsDrag() {
		t.Error("Expected IsDrag() = true before reset")
	}

	// Reset
	tracker.Reset()

	// Check all state is cleared
	if tracker.IsActive() {
		t.Error("Expected drag to be inactive after reset")
	}

	if tracker.IsDrag() {
		t.Error("Expected IsDrag() = false after reset")
	}

	state := tracker.State()
	if state.StartPosition() != value.NewPosition(0, 0) {
		t.Errorf("Expected zero start position after reset, got %v", state.StartPosition())
	}

	if state.Current() != value.NewPosition(0, 0) {
		t.Errorf("Expected zero current position after reset, got %v", state.Current())
	}

	if state.Button() != value.ButtonNone {
		t.Errorf("Expected ButtonNone after reset, got %v", state.Button())
	}

	if state.Modifiers() != value.ModifierNone {
		t.Errorf("Expected ModifierNone after reset, got %v", state.Modifiers())
	}

	if state.Distance() != 0 {
		t.Errorf("Expected distance 0 after reset, got %d", state.Distance())
	}
}

// Test constructor with invalid threshold
func TestDragTracker_Constructor(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		expected  int // Expected threshold after defaulting
	}{
		{
			name:      "Valid threshold",
			threshold: 5,
			expected:  5,
		},
		{
			name:      "Zero threshold (should use 2)",
			threshold: 0,
			expected:  2,
		},
		{
			name:      "Negative threshold (should default to 2)",
			threshold: -10,
			expected:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewDragTracker(tt.threshold)

			// Start drag to check threshold
			pressEvent := model.NewMouseEvent(
				value.EventPress,
				value.ButtonLeft,
				value.NewPosition(10, 10),
				value.ModifierNone,
			)
			tracker.ProcessPress(pressEvent)

			// Move exactly tt.expected cells
			motionEvent := model.NewMouseEvent(
				value.EventMotion,
				value.ButtonLeft,
				value.NewPosition(10+tt.expected, 10),
				value.ModifierNone,
			)
			tracker.ProcessMotion(motionEvent)

			// Should be exactly at threshold
			if !tracker.IsDrag() {
				t.Errorf("Expected IsDrag() = true at threshold %d", tt.expected)
			}

			// Move 1 less than threshold
			tracker.Reset()
			tracker.ProcessPress(pressEvent)
			motionEvent2 := model.NewMouseEvent(
				value.EventMotion,
				value.ButtonLeft,
				value.NewPosition(10+tt.expected-1, 10),
				value.ModifierNone,
			)
			tracker.ProcessMotion(motionEvent2)

			if tracker.IsDrag() {
				t.Errorf("Expected IsDrag() = false below threshold %d", tt.expected)
			}
		})
	}
}

// Test distance calculation
func TestDragTracker_DistanceCalculation(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)

	// Start drag
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierNone,
	)
	tracker.ProcessPress(pressEvent)

	tests := []struct {
		name             string
		pos              value.Position
		expectedDistance int
	}{
		{
			name:             "No movement",
			pos:              value.NewPosition(10, 10),
			expectedDistance: 0,
		},
		{
			name:             "3 cells horizontal",
			pos:              value.NewPosition(13, 10),
			expectedDistance: 3,
		},
		{
			name:             "4 cells vertical",
			pos:              value.NewPosition(10, 14),
			expectedDistance: 4,
		},
		{
			name:             "Manhattan distance 3+4=7",
			pos:              value.NewPosition(13, 14),
			expectedDistance: 7, // Manhattan: |13-10| + |14-10| = 7
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and start
			tracker.Reset()
			tracker.ProcessPress(pressEvent)

			// Move to test position
			motionEvent := model.NewMouseEvent(
				value.EventMotion,
				value.ButtonLeft,
				tt.pos,
				value.ModifierNone,
			)
			tracker.ProcessMotion(motionEvent)

			state := tracker.State()
			if state.Distance() != tt.expectedDistance {
				t.Errorf("Expected distance %d, got %d", tt.expectedDistance, state.Distance())
			}
		})
	}
}

// Test continuous motion updates
func TestDragTracker_ContinuousMotion(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)

	// Start drag
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierNone,
	)
	tracker.ProcessPress(pressEvent)

	// Simulate continuous motion
	positions := []value.Position{
		value.NewPosition(11, 10), // 1 cell (no drag)
		value.NewPosition(12, 10), // 2 cells (drag starts)
		value.NewPosition(13, 10), // 3 cells (drag continues)
		value.NewPosition(14, 10), // 4 cells (drag continues)
		value.NewPosition(15, 10), // 5 cells (drag continues)
	}

	expectedDrag := []bool{false, true, true, true, true}

	for i, pos := range positions {
		motionEvent := model.NewMouseEvent(
			value.EventMotion,
			value.ButtonLeft,
			pos,
			value.ModifierNone,
		)
		dragEvent := tracker.ProcessMotion(motionEvent)

		if expectedDrag[i] { //nolint:nestif // Test validation requires nested checks
			if dragEvent == nil {
				t.Errorf("Position %d: expected drag event", i)
			}
			if !tracker.IsDrag() {
				t.Errorf("Position %d: expected IsDrag() = true", i)
			}
		} else {
			if dragEvent != nil {
				t.Errorf("Position %d: expected no drag event", i)
			}
			if tracker.IsDrag() {
				t.Errorf("Position %d: expected IsDrag() = false", i)
			}
		}

		// Check current position updated
		state := tracker.State()
		if state.Current() != pos {
			t.Errorf("Position %d: expected current %v, got %v", i, pos, state.Current())
		}
	}
}

// Test drag timestamp preservation
func TestDragTracker_TimestampPreservation(t *testing.T) {
	tracker := NewDragTracker(2)
	startPos := value.NewPosition(10, 10)
	endPos := value.NewPosition(15, 10)
	timestamp := time.Now().Add(100 * time.Millisecond)

	// Press event
	pressEvent := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		startPos,
		value.ModifierNone,
	)
	tracker.ProcessPress(pressEvent)

	// Motion event with specific timestamp
	motionEvent := model.NewMouseEventWithTimestamp(
		value.EventMotion,
		value.ButtonLeft,
		endPos,
		value.ModifierNone,
		timestamp,
	)
	dragEvent := tracker.ProcessMotion(motionEvent)

	if dragEvent == nil {
		t.Fatal("Expected drag event")
	}

	// Drag event should preserve motion event's timestamp
	if !dragEvent.Timestamp().Equal(timestamp) {
		t.Errorf("Expected timestamp %v, got %v", timestamp, dragEvent.Timestamp())
	}
}

// Test non-press events don't start drag
func TestDragTracker_OnlyPressStartsDrag(t *testing.T) {
	tracker := NewDragTracker(2)
	pos := value.NewPosition(10, 10)

	nonPressEvents := []value.EventType{
		value.EventRelease,
		value.EventMotion,
		value.EventClick,
		value.EventDoubleClick,
		value.EventScroll,
		value.EventDrag,
	}

	for _, eventType := range nonPressEvents {
		t.Run(eventType.String(), func(t *testing.T) {
			tracker.Reset()

			event := model.NewMouseEvent(
				eventType,
				value.ButtonLeft,
				pos,
				value.ModifierNone,
			)
			tracker.ProcessPress(event)

			if tracker.IsActive() {
				t.Errorf("Event type %v should not start drag", eventType)
			}
		})
	}
}

// Test update without start does nothing
func TestDragTracker_UpdateWithoutStart(t *testing.T) {
	tracker := NewDragTracker(2)

	// Try to update without starting
	state := tracker.State()
	initialCurrent := state.Current()

	// Simulate direct state update (shouldn't happen in practice, but test defensive code)
	state.Update(value.NewPosition(100, 100))

	// Current should not have changed because drag wasn't active
	if state.Current() != initialCurrent {
		t.Error("Current position should not change when drag is inactive")
	}
}
