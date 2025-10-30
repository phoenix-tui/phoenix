package service

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// TestNewHoverTracker tests the constructor.
func TestNewHoverTracker(t *testing.T) {
	tracker := NewHoverTracker()

	if tracker == nil {
		t.Fatal("NewHoverTracker() returned nil")
	}

	if tracker.IsHovering() {
		t.Error("Expected IsHovering() to be false for new tracker")
	}

	if tracker.CurrentComponentID() != "" {
		t.Errorf("Expected empty componentID, got %s", tracker.CurrentComponentID())
	}

	if tracker.State() == nil {
		t.Error("Expected State() to return non-nil")
	}
}

// TestHoverTrackerUpdate_Enter tests entering a component.
func TestHoverTrackerUpdate_Enter(t *testing.T) {
	tests := []struct {
		name        string
		position    value2.Position
		areas       []ComponentArea
		expectEvent value2.EventType
		expectID    string
		expectHover bool
	}{
		{
			name:     "enter single component",
			position: value2.NewPosition(10, 5),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverEnter,
			expectID:    "button1",
			expectHover: true,
		},
		{
			name:     "enter first of multiple components",
			position: value2.NewPosition(10, 5),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
				{ID: "button2", Area: value2.NewBoundingBox(5, 10, 20, 5)},
			},
			expectEvent: value2.EventHoverEnter,
			expectID:    "button1",
			expectHover: true,
		},
		{
			name:     "enter at exact boundary",
			position: value2.NewPosition(5, 3),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverEnter,
			expectID:    "button1",
			expectHover: true,
		},
		{
			name:     "no component at position",
			position: value2.NewPosition(50, 50),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventMotion,
			expectID:    "",
			expectHover: false,
		},
		{
			name:     "overlapping components - first wins",
			position: value2.NewPosition(10, 5),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
				{ID: "button2", Area: value2.NewBoundingBox(8, 4, 15, 4)}, // Overlaps
			},
			expectEvent: value2.EventHoverEnter,
			expectID:    "button1",
			expectHover: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewHoverTracker()

			event := tracker.Update(tt.position, tt.areas)

			if event != tt.expectEvent {
				t.Errorf("Expected event %s, got %s", tt.expectEvent, event)
			}

			if tracker.CurrentComponentID() != tt.expectID {
				t.Errorf("Expected componentID %s, got %s", tt.expectID, tracker.CurrentComponentID())
			}

			if tracker.IsHovering() != tt.expectHover {
				t.Errorf("Expected IsHovering() to be %v, got %v", tt.expectHover, tracker.IsHovering())
			}
		})
	}
}

// TestHoverTrackerUpdate_Move tests moving within a component.
func TestHoverTrackerUpdate_Move(t *testing.T) {
	tests := []struct {
		name        string
		setupPos    value2.Position
		movePos     value2.Position
		areas       []ComponentArea
		expectEvent value2.EventType
		expectID    string
	}{
		{
			name:     "move within same component",
			setupPos: value2.NewPosition(10, 5),
			movePos:  value2.NewPosition(15, 6),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverMove,
			expectID:    "button1",
		},
		{
			name:     "move to same position",
			setupPos: value2.NewPosition(10, 5),
			movePos:  value2.NewPosition(10, 5),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverMove,
			expectID:    "button1",
		},
		{
			name:     "move across component boundary",
			setupPos: value2.NewPosition(10, 5),
			movePos:  value2.NewPosition(24, 5),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverMove,
			expectID:    "button1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewHoverTracker()

			// Setup: enter component
			enterEvent := tracker.Update(tt.setupPos, tt.areas)
			if enterEvent != value2.EventHoverEnter {
				t.Fatalf("Setup failed: expected HoverEnter, got %s", enterEvent)
			}

			// Test: move within component
			event := tracker.Update(tt.movePos, tt.areas)

			if event != tt.expectEvent {
				t.Errorf("Expected event %s, got %s", tt.expectEvent, event)
			}

			if tracker.CurrentComponentID() != tt.expectID {
				t.Errorf("Expected componentID %s, got %s", tt.expectID, tracker.CurrentComponentID())
			}

			if !tracker.IsHovering() {
				t.Error("Expected IsHovering() to be true")
			}
		})
	}
}

// TestHoverTrackerUpdate_Leave tests leaving a component.
func TestHoverTrackerUpdate_Leave(t *testing.T) {
	tests := []struct {
		name        string
		setupPos    value2.Position
		leavePos    value2.Position
		areas       []ComponentArea
		expectEvent value2.EventType
		expectID    string
		expectHover bool
	}{
		{
			name:     "leave component area",
			setupPos: value2.NewPosition(10, 5),
			leavePos: value2.NewPosition(50, 50),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverLeave,
			expectID:    "",
			expectHover: false,
		},
		{
			name:     "leave at exact boundary",
			setupPos: value2.NewPosition(10, 5),
			leavePos: value2.NewPosition(25, 5), // Just outside (x + width)
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
			},
			expectEvent: value2.EventHoverLeave,
			expectID:    "",
			expectHover: false,
		},
		{
			name:     "leave one component, enter another",
			setupPos: value2.NewPosition(10, 5),
			leavePos: value2.NewPosition(10, 12),
			areas: []ComponentArea{
				{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
				{ID: "button2", Area: value2.NewBoundingBox(5, 10, 20, 5)},
			},
			expectEvent: value2.EventHoverEnter, // Treated as enter to new component
			expectID:    "button2",
			expectHover: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewHoverTracker()

			// Setup: enter component
			enterEvent := tracker.Update(tt.setupPos, tt.areas)
			if enterEvent != value2.EventHoverEnter {
				t.Fatalf("Setup failed: expected HoverEnter, got %s", enterEvent)
			}

			// Test: leave component
			event := tracker.Update(tt.leavePos, tt.areas)

			if event != tt.expectEvent {
				t.Errorf("Expected event %s, got %s", tt.expectEvent, event)
			}

			if tracker.CurrentComponentID() != tt.expectID {
				t.Errorf("Expected componentID %s, got %s", tt.expectID, tracker.CurrentComponentID())
			}

			if tracker.IsHovering() != tt.expectHover {
				t.Errorf("Expected IsHovering() to be %v, got %v", tt.expectHover, tracker.IsHovering())
			}
		})
	}
}

// TestHoverTrackerUpdate_ComponentSwitch tests switching between components.
func TestHoverTrackerUpdate_ComponentSwitch(t *testing.T) {
	tracker := NewHoverTracker()

	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 3)},
		{ID: "button2", Area: value2.NewBoundingBox(5, 8, 20, 3)},
		{ID: "button3", Area: value2.NewBoundingBox(5, 13, 20, 3)},
	}

	// Enter button1
	event := tracker.Update(value2.NewPosition(10, 4), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter for button1, got %s", event)
	}
	if tracker.CurrentComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", tracker.CurrentComponentID())
	}

	// Move to button2 (should be treated as enter)
	event = tracker.Update(value2.NewPosition(10, 9), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter for button2, got %s", event)
	}
	if tracker.CurrentComponentID() != "button2" {
		t.Errorf("Expected button2, got %s", tracker.CurrentComponentID())
	}

	// Move to button3
	event = tracker.Update(value2.NewPosition(10, 14), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter for button3, got %s", event)
	}
	if tracker.CurrentComponentID() != "button3" {
		t.Errorf("Expected button3, got %s", tracker.CurrentComponentID())
	}

	// Leave all components
	event = tracker.Update(value2.NewPosition(50, 50), areas)
	if event != value2.EventHoverLeave {
		t.Errorf("Expected HoverLeave, got %s", event)
	}
	if tracker.CurrentComponentID() != "" {
		t.Errorf("Expected empty componentID, got %s", tracker.CurrentComponentID())
	}
}

// TestHoverTrackerUpdate_EmptyAreas tests behavior with no component areas.
func TestHoverTrackerUpdate_EmptyAreas(t *testing.T) {
	tracker := NewHoverTracker()

	// Update with no areas
	event := tracker.Update(value2.NewPosition(10, 5), []ComponentArea{})

	if event != value2.EventMotion {
		t.Errorf("Expected EventMotion with no areas, got %s", event)
	}

	if tracker.IsHovering() {
		t.Error("Expected IsHovering() to be false with no areas")
	}

	if tracker.CurrentComponentID() != "" {
		t.Errorf("Expected empty componentID, got %s", tracker.CurrentComponentID())
	}
}

// TestHoverTrackerUpdate_ZeroSizeArea tests zero-size bounding boxes.
func TestHoverTrackerUpdate_ZeroSizeArea(t *testing.T) {
	tracker := NewHoverTracker()

	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(5, 3, 0, 0)}, // Zero size
	}

	// Try to hover over zero-size area
	event := tracker.Update(value2.NewPosition(5, 3), areas)

	// Zero-size areas should not contain any point
	if event != value2.EventMotion {
		t.Errorf("Expected EventMotion for zero-size area, got %s", event)
	}

	if tracker.IsHovering() {
		t.Error("Expected IsHovering() to be false for zero-size area")
	}
}

// TestHoverTrackerReset tests the Reset method.
func TestHoverTrackerReset(t *testing.T) {
	tracker := NewHoverTracker()

	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
	}

	// Setup: enter component
	event := tracker.Update(value2.NewPosition(10, 5), areas)
	if event != value2.EventHoverEnter {
		t.Fatalf("Setup failed: expected HoverEnter, got %s", event)
	}

	// Verify hovering
	if !tracker.IsHovering() {
		t.Error("Expected hovering before reset")
	}

	// Reset
	tracker.Reset()

	// Verify reset
	if tracker.IsHovering() {
		t.Error("Expected IsHovering() to be false after Reset()")
	}

	if tracker.CurrentComponentID() != "" {
		t.Errorf("Expected empty componentID after Reset(), got %s", tracker.CurrentComponentID())
	}

	// State should be reset
	state := tracker.State()
	if state.IsActive() {
		t.Error("Expected state to be inactive after Reset()")
	}
}

// TestHoverTrackerState tests the State accessor.
func TestHoverTrackerState(t *testing.T) {
	tracker := NewHoverTracker()

	state := tracker.State()
	if state == nil {
		t.Fatal("Expected State() to return non-nil")
	}

	// State should reflect tracker's current state
	if state.IsActive() {
		t.Error("Expected state to be inactive initially")
	}

	// Enter component
	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 5)},
	}
	tracker.Update(value2.NewPosition(10, 5), areas)

	// State should now be active
	if !state.IsActive() {
		t.Error("Expected state to be active after entering component")
	}

	if state.ComponentID() != "button1" {
		t.Errorf("Expected state componentID button1, got %s", state.ComponentID())
	}
}

// TestHoverTrackerSequence tests a complete hover sequence.
func TestHoverTrackerSequence(t *testing.T) {
	tracker := NewHoverTracker()

	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(5, 3, 20, 3)},
		{ID: "button2", Area: value2.NewBoundingBox(5, 8, 20, 3)},
	}

	// Start outside any component
	event := tracker.Update(value2.NewPosition(0, 0), areas)
	if event != value2.EventMotion {
		t.Errorf("Expected Motion, got %s", event)
	}

	// Enter button1
	event = tracker.Update(value2.NewPosition(10, 4), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter, got %s", event)
	}
	if tracker.CurrentComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", tracker.CurrentComponentID())
	}

	// Move within button1
	event = tracker.Update(value2.NewPosition(15, 4), areas)
	if event != value2.EventHoverMove {
		t.Errorf("Expected HoverMove, got %s", event)
	}
	if tracker.CurrentComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", tracker.CurrentComponentID())
	}

	// Leave to empty space
	event = tracker.Update(value2.NewPosition(0, 7), areas)
	if event != value2.EventHoverLeave {
		t.Errorf("Expected HoverLeave, got %s", event)
	}
	if tracker.IsHovering() {
		t.Error("Should not be hovering after leaving")
	}

	// Enter button2
	event = tracker.Update(value2.NewPosition(10, 9), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter, got %s", event)
	}
	if tracker.CurrentComponentID() != "button2" {
		t.Errorf("Expected button2, got %s", tracker.CurrentComponentID())
	}

	// Switch directly to button1
	event = tracker.Update(value2.NewPosition(10, 4), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter, got %s", event)
	}
	if tracker.CurrentComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", tracker.CurrentComponentID())
	}

	// Reset
	tracker.Reset()
	if tracker.IsHovering() {
		t.Error("Should not be hovering after reset")
	}
}

// TestHoverTrackerUpdate_NegativeCoordinates tests negative coordinate handling.
func TestHoverTrackerUpdate_NegativeCoordinates(t *testing.T) {
	tracker := NewHoverTracker()

	// Component area with negative coordinates
	areas := []ComponentArea{
		{ID: "button1", Area: value2.NewBoundingBox(-10, -5, 20, 10)},
	}

	// Test position inside negative coordinate area
	event := tracker.Update(value2.NewPosition(-5, 0), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter for negative coordinate area, got %s", event)
	}

	if tracker.CurrentComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", tracker.CurrentComponentID())
	}

	// Leave negative area
	event = tracker.Update(value2.NewPosition(50, 50), areas)
	if event != value2.EventHoverLeave {
		t.Errorf("Expected HoverLeave, got %s", event)
	}
}

// TestHoverTrackerUpdate_LargeCoordinates tests large coordinate handling.
func TestHoverTrackerUpdate_LargeCoordinates(t *testing.T) {
	tracker := NewHoverTracker()

	// Component area with large coordinates
	areas := []ComponentArea{
		{ID: "panel", Area: value2.NewBoundingBox(1000, 500, 2000, 1000)},
	}

	// Test position inside large coordinate area
	event := tracker.Update(value2.NewPosition(2000, 1000), areas)
	if event != value2.EventHoverEnter {
		t.Errorf("Expected HoverEnter for large coordinate area, got %s", event)
	}

	if tracker.CurrentComponentID() != "panel" {
		t.Errorf("Expected panel, got %s", tracker.CurrentComponentID())
	}
}
