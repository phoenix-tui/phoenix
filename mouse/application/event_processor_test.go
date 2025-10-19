package application

import (
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

func TestNewEventProcessor(t *testing.T) {
	processor := NewEventProcessor()

	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	if processor.clickDetector == nil {
		t.Error("Expected click detector to be initialized")
	}

	if processor.dragTracker == nil {
		t.Error("Expected drag tracker to be initialized")
	}

	if processor.scrollCalculator == nil {
		t.Error("Expected scroll calculator to be initialized")
	}

	// Verify initial state
	if processor.IsDragging() {
		t.Error("Expected not dragging initially")
	}

	if processor.ClickCount() != 0 {
		t.Errorf("Expected click count 0, got %d", processor.ClickCount())
	}
}

func TestNewEventProcessorWithConfig(t *testing.T) {
	processor := NewEventProcessorWithConfig(
		300*time.Millisecond, // click timeout
		2,                    // click tolerance
		3,                    // drag threshold
		5,                    // lines per scroll
	)

	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	// Verify configuration is applied by testing behavior
	if processor.clickDetector == nil {
		t.Error("Expected click detector to be initialized")
	}
}

func TestProcessEvent_Press(t *testing.T) {
	processor := NewEventProcessor()

	event := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(event)

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type() != value.EventPress {
		t.Errorf("Expected EventPress, got %v", events[0].Type())
	}
}

func TestProcessEvent_Release_WithoutDrag(t *testing.T) {
	processor := NewEventProcessor()

	event := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(event)

	// Should return both click and release events
	if len(events) < 1 {
		t.Fatalf("Expected at least 1 event, got %d", len(events))
	}

	// Last event should be release
	lastEvent := events[len(events)-1]
	if lastEvent.Type() != value.EventRelease {
		t.Errorf("Expected last event to be EventRelease, got %v", lastEvent.Type())
	}
}

func TestProcessEvent_Release_AfterDrag(t *testing.T) {
	processor := NewEventProcessor()

	// Start drag with press
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	// Move enough to trigger drag
	motion := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(15, 5), // Move 5 cells (> threshold)
		value.ModifierNone,
	)
	processor.ProcessEvent(motion)

	// Release after drag
	release := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		value.NewPosition(15, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(release)

	// Should only return release, NOT click (drag suppresses click)
	if len(events) != 1 {
		t.Fatalf("Expected 1 event (release only), got %d", len(events))
	}

	if events[0].Type() != value.EventRelease {
		t.Errorf("Expected EventRelease, got %v", events[0].Type())
	}

	// Should not generate click after drag
	hasClick := false
	for _, e := range events {
		if e.IsClick() {
			hasClick = true
			break
		}
	}

	if hasClick {
		t.Error("Should not generate click after drag")
	}
}

func TestProcessEvent_Motion_NoDrag(t *testing.T) {
	processor := NewEventProcessor()

	event := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(event)

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type() != value.EventMotion {
		t.Errorf("Expected EventMotion, got %v", events[0].Type())
	}
}

func TestProcessEvent_Motion_WithDrag(t *testing.T) {
	processor := NewEventProcessor()

	// Start drag with press
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	// Move enough to trigger drag (threshold = 2)
	motion := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(13, 5), // Move 3 cells (> 2 threshold)
		value.ModifierNone,
	)

	events := processor.ProcessEvent(motion)

	// Should return drag event
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type() != value.EventDrag {
		t.Errorf("Expected EventDrag, got %v", events[0].Type())
	}

	// Verify dragging state
	if !processor.IsDragging() {
		t.Error("Expected to be dragging")
	}
}

func TestProcessEvent_Scroll(t *testing.T) {
	processor := NewEventProcessor()

	event := model.NewMouseEvent(
		value.EventScroll,
		value.ButtonWheelUp,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(event)

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type() != value.EventScroll {
		t.Errorf("Expected EventScroll, got %v", events[0].Type())
	}
}

func TestProcessEvent_UnknownType(t *testing.T) {
	processor := NewEventProcessor()

	// Create event with unknown type (255)
	event := model.NewMouseEvent(
		value.EventType(255),
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	events := processor.ProcessEvent(event)

	// Should pass through unknown events
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type() != value.EventType(255) {
		t.Errorf("Expected EventType(255), got %v", events[0].Type())
	}
}

func TestProcessEvent_ClickDetection(t *testing.T) {
	processor := NewEventProcessor()

	pos := value.NewPosition(10, 5)

	// First release (single click)
	event1 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	events1 := processor.ProcessEvent(event1)

	// Should contain click event
	hasClick := false
	for _, e := range events1 {
		if e.Type() == value.EventClick {
			hasClick = true
			break
		}
	}

	if !hasClick {
		t.Error("Expected click event in results")
	}

	if processor.ClickCount() != 1 {
		t.Errorf("Expected click count 1, got %d", processor.ClickCount())
	}

	// Second release (double click)
	event2 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	events2 := processor.ProcessEvent(event2)

	// Should contain double click event
	hasDoubleClick := false
	for _, e := range events2 {
		if e.Type() == value.EventDoubleClick {
			hasDoubleClick = true
			break
		}
	}

	if !hasDoubleClick {
		t.Error("Expected double click event in results")
	}

	if processor.ClickCount() != 2 {
		t.Errorf("Expected click count 2, got %d", processor.ClickCount())
	}
}

func TestProcessEvent_TripleClick(t *testing.T) {
	processor := NewEventProcessor()

	pos := value.NewPosition(10, 5)

	// Three rapid releases
	for i := 0; i < 3; i++ {
		event := model.NewMouseEvent(
			value.EventRelease,
			value.ButtonLeft,
			pos,
			value.ModifierNone,
		)
		processor.ProcessEvent(event)
	}

	if processor.ClickCount() != 3 {
		t.Errorf("Expected click count 3, got %d", processor.ClickCount())
	}
}

func TestReset(t *testing.T) {
	processor := NewEventProcessor()

	// Create some state
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	motion := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(15, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(motion)

	// Verify state exists
	if !processor.IsDragging() {
		t.Error("Expected to be dragging before reset")
	}

	// Reset
	processor.Reset()

	// Verify state cleared
	if processor.IsDragging() {
		t.Error("Expected not dragging after reset")
	}

	if processor.ClickCount() != 0 {
		t.Errorf("Expected click count 0 after reset, got %d", processor.ClickCount())
	}
}

func TestScrollDelta(t *testing.T) {
	processor := NewEventProcessor()

	tests := []struct {
		name     string
		button   value.Button
		expected int
	}{
		{
			name:     "Scroll up",
			button:   value.ButtonWheelUp,
			expected: -3, // Default 3 lines per scroll
		},
		{
			name:     "Scroll down",
			button:   value.ButtonWheelDown,
			expected: 3,
		},
		{
			name:     "Non-scroll button",
			button:   value.ButtonLeft,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := model.NewMouseEvent(
				value.EventScroll,
				tt.button,
				value.NewPosition(10, 5),
				value.ModifierNone,
			)

			delta := processor.ScrollDelta(event)
			if delta != tt.expected {
				t.Errorf("Expected delta %d, got %d", tt.expected, delta)
			}
		})
	}
}

func TestIsScrollUp(t *testing.T) {
	processor := NewEventProcessor()

	tests := []struct {
		name     string
		button   value.Button
		expected bool
	}{
		{
			name:     "Wheel up",
			button:   value.ButtonWheelUp,
			expected: true,
		},
		{
			name:     "Wheel down",
			button:   value.ButtonWheelDown,
			expected: false,
		},
		{
			name:     "Left button",
			button:   value.ButtonLeft,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := model.NewMouseEvent(
				value.EventScroll,
				tt.button,
				value.NewPosition(10, 5),
				value.ModifierNone,
			)

			result := processor.IsScrollUp(event)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsScrollDown(t *testing.T) {
	processor := NewEventProcessor()

	tests := []struct {
		name     string
		button   value.Button
		expected bool
	}{
		{
			name:     "Wheel down",
			button:   value.ButtonWheelDown,
			expected: true,
		},
		{
			name:     "Wheel up",
			button:   value.ButtonWheelUp,
			expected: false,
		},
		{
			name:     "Right button",
			button:   value.ButtonRight,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := model.NewMouseEvent(
				value.EventScroll,
				tt.button,
				value.NewPosition(10, 5),
				value.ModifierNone,
			)

			result := processor.IsScrollDown(event)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsDragging(t *testing.T) {
	processor := NewEventProcessor()

	// Initially not dragging
	if processor.IsDragging() {
		t.Error("Expected not dragging initially")
	}

	// Start drag
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	// Motion without threshold
	motion1 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(11, 5), // Only 1 cell (< threshold)
		value.ModifierNone,
	)
	processor.ProcessEvent(motion1)

	if processor.IsDragging() {
		t.Error("Expected not dragging yet (below threshold)")
	}

	// Motion exceeding threshold
	motion2 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(13, 5), // 3 cells (> 2 threshold)
		value.ModifierNone,
	)
	processor.ProcessEvent(motion2)

	if !processor.IsDragging() {
		t.Error("Expected to be dragging after threshold exceeded")
	}

	// Release
	release := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		value.NewPosition(13, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(release)

	if processor.IsDragging() {
		t.Error("Expected not dragging after release")
	}
}

func TestClickCount(t *testing.T) {
	processor := NewEventProcessor()

	if processor.ClickCount() != 0 {
		t.Errorf("Expected initial click count 0, got %d", processor.ClickCount())
	}

	pos := value.NewPosition(10, 5)

	// Single click
	event1 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	processor.ProcessEvent(event1)

	if processor.ClickCount() != 1 {
		t.Errorf("Expected click count 1, got %d", processor.ClickCount())
	}

	// Double click
	event2 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	processor.ProcessEvent(event2)

	if processor.ClickCount() != 2 {
		t.Errorf("Expected click count 2, got %d", processor.ClickCount())
	}

	// Triple click
	event3 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	processor.ProcessEvent(event3)

	if processor.ClickCount() != 3 {
		t.Errorf("Expected click count 3, got %d", processor.ClickCount())
	}
}

// Integration test: Complex drag scenario
func TestIntegration_DragSequence(t *testing.T) {
	processor := NewEventProcessor()

	// Press at (10, 10)
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 10),
		value.ModifierNone,
	)
	events := processor.ProcessEvent(press)

	if len(events) != 1 || events[0].Type() != value.EventPress {
		t.Error("Press event should pass through")
	}

	// Move within threshold (no drag yet)
	motion1 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(11, 10),
		value.ModifierNone,
	)
	events = processor.ProcessEvent(motion1)

	if len(events) != 1 || events[0].Type() != value.EventMotion {
		t.Error("Motion within threshold should pass through as motion")
	}

	// Move beyond threshold (drag starts)
	motion2 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(13, 10),
		value.ModifierNone,
	)
	events = processor.ProcessEvent(motion2)

	if len(events) != 1 || events[0].Type() != value.EventDrag {
		t.Error("Motion beyond threshold should become drag")
	}

	if !processor.IsDragging() {
		t.Error("Should be dragging")
	}

	// Continue dragging
	motion3 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(20, 10),
		value.ModifierNone,
	)
	events = processor.ProcessEvent(motion3)

	if len(events) != 1 || events[0].Type() != value.EventDrag {
		t.Error("Continued motion should remain drag")
	}

	// Release (end drag, no click)
	release := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		value.NewPosition(20, 10),
		value.ModifierNone,
	)
	events = processor.ProcessEvent(release)

	if len(events) != 1 || events[0].Type() != value.EventRelease {
		t.Error("Release after drag should not generate click")
	}

	if processor.IsDragging() {
		t.Error("Should not be dragging after release")
	}

	// Click count should remain 0 (drag suppressed click)
	if processor.ClickCount() != 0 {
		t.Errorf("Drag should suppress click, got count %d", processor.ClickCount())
	}
}

// Integration test: Multiple clicks with scrolling
func TestIntegration_ClicksAndScroll(t *testing.T) {
	processor := NewEventProcessor()

	pos := value.NewPosition(10, 5)

	// Double click
	event1 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	processor.ProcessEvent(event1)

	event2 := model.NewMouseEvent(
		value.EventRelease,
		value.ButtonLeft,
		pos,
		value.ModifierNone,
	)
	processor.ProcessEvent(event2)

	if processor.ClickCount() != 2 {
		t.Errorf("Expected click count 2, got %d", processor.ClickCount())
	}

	// Scroll event (should not affect clicks)
	scroll := model.NewMouseEvent(
		value.EventScroll,
		value.ButtonWheelDown,
		pos,
		value.ModifierNone,
	)
	events := processor.ProcessEvent(scroll)

	if len(events) != 1 || events[0].Type() != value.EventScroll {
		t.Error("Scroll should pass through")
	}

	// Verify scroll helpers
	if !processor.IsScrollDown(scroll) {
		t.Error("Expected scroll down")
	}

	if processor.IsScrollUp(scroll) {
		t.Error("Expected not scroll up")
	}

	delta := processor.ScrollDelta(scroll)
	if delta != 3 {
		t.Errorf("Expected delta 3, got %d", delta)
	}

	// Click count should not be affected by scroll
	if processor.ClickCount() != 2 {
		t.Errorf("Scroll should not affect click count, got %d", processor.ClickCount())
	}
}

// Edge case: Multiple resets
func TestMultipleResets(t *testing.T) {
	processor := NewEventProcessor()

	// Reset empty state
	processor.Reset()
	if processor.ClickCount() != 0 {
		t.Error("Empty reset should work")
	}

	// Create state
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	// Multiple resets
	for i := 0; i < 5; i++ {
		processor.Reset()
		if processor.IsDragging() {
			t.Errorf("Reset %d: should not be dragging", i)
		}
		if processor.ClickCount() != 0 {
			t.Errorf("Reset %d: click count should be 0", i)
		}
	}
}

// Edge case: Custom configuration
func TestCustomConfiguration(t *testing.T) {
	// Very short timeout, low tolerance, high threshold, many lines per scroll
	processor := NewEventProcessorWithConfig(
		50*time.Millisecond, // 50ms click timeout
		0,                   // 0 tolerance (exact position)
		10,                  // 10 cell drag threshold
		10,                  // 10 lines per scroll
	)

	// Verify scroll configuration
	scroll := model.NewMouseEvent(
		value.EventScroll,
		value.ButtonWheelDown,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	delta := processor.ScrollDelta(scroll)
	if delta != 10 {
		t.Errorf("Expected delta 10 (custom config), got %d", delta)
	}

	// Verify drag threshold (10 cells)
	press := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 10),
		value.ModifierNone,
	)
	processor.ProcessEvent(press)

	// Move 5 cells (below 10 threshold)
	motion1 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(15, 10),
		value.ModifierNone,
	)
	events := processor.ProcessEvent(motion1)

	if len(events) != 1 || events[0].Type() != value.EventMotion {
		t.Error("5 cell motion should not trigger drag with 10 cell threshold")
	}

	// Move 11 cells (above 10 threshold)
	motion2 := model.NewMouseEvent(
		value.EventMotion,
		value.ButtonLeft,
		value.NewPosition(21, 10),
		value.ModifierNone,
	)
	events = processor.ProcessEvent(motion2)

	if len(events) != 1 || events[0].Type() != value.EventDrag {
		t.Error("11 cell motion should trigger drag with 10 cell threshold")
	}
}
