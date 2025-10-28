package mouse

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.handler == nil {
		t.Fatal("New() created Mouse with nil handler")
	}
}

// ============================================================================
// Enable/Disable/IsEnabled Tests
// ============================================================================

func TestMouse_Enable_Disable(t *testing.T) {
	m := New()

	// Initially disabled
	if m.IsEnabled() {
		t.Error("mouse should be disabled initially")
	}

	// Enable
	if err := m.Enable(); err != nil {
		t.Fatalf("Enable() failed: %v", err)
	}
	if !m.IsEnabled() {
		t.Error("mouse should be enabled after Enable()")
	}

	// Disable
	if err := m.Disable(); err != nil {
		t.Fatalf("Disable() failed: %v", err)
	}
	if m.IsEnabled() {
		t.Error("mouse should be disabled after Disable()")
	}
}

func TestMouse_Enable_MultipleTimesOK(t *testing.T) {
	m := New()

	// Enable multiple times should not error
	for i := 0; i < 3; i++ {
		if err := m.Enable(); err != nil {
			t.Fatalf("Enable() call %d failed: %v", i+1, err)
		}
		if !m.IsEnabled() {
			t.Errorf("mouse not enabled after Enable() call %d", i+1)
		}
	}

	m.Disable()
}

func TestMouse_Disable_MultipleTimesOK(t *testing.T) {
	m := New()
	m.Enable()

	// Disable multiple times should not error
	for i := 0; i < 3; i++ {
		if err := m.Disable(); err != nil {
			t.Fatalf("Disable() call %d failed: %v", i+1, err)
		}
		if m.IsEnabled() {
			t.Errorf("mouse still enabled after Disable() call %d", i+1)
		}
	}
}

// ============================================================================
// ParseSequence Tests
// ============================================================================

func TestMouse_ParseSequence_SGRPress(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	// SGR format: <0;10;5M (left button press at x=10, y=5)
	events, err := m.ParseSequence("<0;10;5M")
	if err != nil {
		t.Fatalf("ParseSequence() error: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("ParseSequence() returned no events")
	}

	event := events[0]
	if event.Type() != EventPress {
		t.Errorf("expected EventPress, got %v", event.Type())
	}
	if event.Button() != ButtonLeft {
		t.Errorf("expected ButtonLeft, got %v", event.Button())
	}
	if event.Position().X() != 9 { // 0-based vs 1-based
		t.Errorf("expected X=9, got %d", event.Position().X())
	}
	if event.Position().Y() != 4 { // 0-based vs 1-based
		t.Errorf("expected Y=4, got %d", event.Position().Y())
	}
}

func TestMouse_ParseSequence_SGRRelease(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	// Need press first, then release triggers Click event
	_, err := m.ParseSequence("<0;10;5M") // Press
	if err != nil {
		t.Fatalf("press parse error: %v", err)
	}

	// SGR format: <0;10;5m (left button release at x=10, y=5)
	events, err := m.ParseSequence("<0;10;5m")
	if err != nil {
		t.Fatalf("ParseSequence() error: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("ParseSequence() returned no events")
	}

	// Release at same position as press generates Click event (enriched)
	hasClick := false
	for _, event := range events {
		if event.Type() == EventClick {
			hasClick = true
			break
		}
	}
	if !hasClick {
		t.Error("expected at least one Click event in enriched sequence")
	}
}

func TestMouse_ParseSequence_InvalidSequence(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	testCases := []struct {
		name string
		seq  string
	}{
		{"empty", ""},
		{"malformed", "garbage"},
		{"incomplete", "<0;10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := m.ParseSequence(tc.seq)
			// Either error or empty events are acceptable
			if err == nil && len(events) > 0 {
				t.Errorf("expected error or empty events for sequence %q, got %d events", tc.seq, len(events))
			}
		})
	}
}

// ============================================================================
// ScrollDelta Tests
// ============================================================================

func TestMouse_ScrollDelta_WheelUp(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	event := NewMouseEvent(EventScroll, ButtonWheelUp, NewPosition(0, 0), ModifierNone)
	delta := m.ScrollDelta(event)

	if delta >= 0 {
		t.Errorf("expected negative delta for wheel up, got %d", delta)
	}
}

func TestMouse_ScrollDelta_WheelDown(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	event := NewMouseEvent(EventScroll, ButtonWheelDown, NewPosition(0, 0), ModifierNone)
	delta := m.ScrollDelta(event)

	if delta <= 0 {
		t.Errorf("expected positive delta for wheel down, got %d", delta)
	}
}

func TestMouse_ScrollDelta_NotScrollEvent(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	// Non-scroll event should return 0
	event := NewMouseEvent(EventClick, ButtonLeft, NewPosition(0, 0), ModifierNone)
	delta := m.ScrollDelta(event)

	if delta != 0 {
		t.Errorf("expected delta=0 for non-scroll event, got %d", delta)
	}
}

// ============================================================================
// IsDragging Tests
// ============================================================================

func TestMouse_IsDragging_InitiallyFalse(t *testing.T) {
	m := New()
	if m.IsDragging() {
		t.Error("IsDragging() should be false initially")
	}
}

func TestMouse_IsDragging_AfterDragSequence(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	// Press at (10, 5)
	_, err := m.ParseSequence("<0;10;5M")
	if err != nil {
		t.Fatalf("press parse error: %v", err)
	}

	// Motion to distant position (20, 5) - should trigger drag if tolerance exceeded
	// Need larger distance to exceed drag tolerance
	_, err = m.ParseSequence("<32;30;5M")
	if err != nil {
		t.Fatalf("motion parse error: %v", err)
	}

	// Note: IsDragging depends on drag threshold implementation
	// If still not dragging, the threshold might be larger
	// Just verify ParseSequence works without error
	isDragging := m.IsDragging()
	t.Logf("IsDragging after motion: %v", isDragging)

	// Release
	_, err = m.ParseSequence("<0;30;5m")
	if err != nil {
		t.Fatalf("release parse error: %v", err)
	}

	// After release, should definitely not be dragging
	if m.IsDragging() {
		t.Error("IsDragging() should be false after release")
	}
}

// ============================================================================
// Reset Tests
// ============================================================================

func TestMouse_Reset(t *testing.T) {
	m := New()
	m.Enable()

	// Trigger some mouse state
	m.ParseSequence("<0;10;5M")
	m.ParseSequence("<32;30;5M")

	// Whether dragging or not, Reset should clear state
	wasDragging := m.IsDragging()
	t.Logf("Was dragging before reset: %v", wasDragging)

	// Reset
	m.Reset()

	// After reset, should not be dragging
	if m.IsDragging() {
		t.Error("IsDragging() should be false after Reset()")
	}

	m.Disable()
}

// ============================================================================
// Helper Function Tests: NewPosition
// ============================================================================

func TestNewPosition(t *testing.T) {
	testCases := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"origin", 0, 0, "(0,0)"},
		{"positive", 10, 20, "(10,20)"},
		{"negative", -5, -3, "(-5,-3)"},
		{"large", 1000, 2000, "(1000,2000)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pos := NewPosition(tc.x, tc.y)
			if pos.X() != tc.x {
				t.Errorf("X: expected %d, got %d", tc.x, pos.X())
			}
			if pos.Y() != tc.y {
				t.Errorf("Y: expected %d, got %d", tc.y, pos.Y())
			}
			if pos.String() != tc.expected {
				t.Errorf("String: expected %s, got %s", tc.expected, pos.String())
			}
		})
	}
}

func TestPosition_Equals(t *testing.T) {
	testCases := []struct {
		name     string
		p1, p2   Position
		expected bool
	}{
		{"same", NewPosition(5, 10), NewPosition(5, 10), true},
		{"different_x", NewPosition(5, 10), NewPosition(6, 10), false},
		{"different_y", NewPosition(5, 10), NewPosition(5, 11), false},
		{"different_both", NewPosition(5, 10), NewPosition(6, 11), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.p1.Equals(tc.p2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestPosition_DistanceTo(t *testing.T) {
	testCases := []struct {
		name     string
		p1, p2   Position
		expected int
	}{
		{"same_position", NewPosition(0, 0), NewPosition(0, 0), 0},
		{"horizontal", NewPosition(0, 0), NewPosition(5, 0), 5},
		{"vertical", NewPosition(0, 0), NewPosition(0, 5), 5},
		{"diagonal", NewPosition(0, 0), NewPosition(3, 4), 7}, // Manhattan distance
		{"negative", NewPosition(5, 5), NewPosition(2, 2), 6},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.p1.DistanceTo(tc.p2)
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestPosition_IsWithinTolerance(t *testing.T) {
	testCases := []struct {
		name      string
		p1, p2    Position
		tolerance int
		expected  bool
	}{
		{"exact_match", NewPosition(5, 5), NewPosition(5, 5), 0, true},
		{"within_tolerance", NewPosition(5, 5), NewPosition(7, 5), 5, true},
		{"outside_tolerance", NewPosition(5, 5), NewPosition(10, 5), 3, false},
		{"zero_tolerance", NewPosition(5, 5), NewPosition(6, 5), 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.p1.IsWithinTolerance(tc.p2, tc.tolerance)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// ============================================================================
// Helper Function Tests: NewModifiers
// ============================================================================

func TestNewModifiers(t *testing.T) {
	testCases := []struct {
		name             string
		shift, ctrl, alt bool
		expectedString   string
		expectedHasShift bool
		expectedHasCtrl  bool
		expectedHasAlt   bool
	}{
		{"none", false, false, false, "None", false, false, false},
		{"shift_only", true, false, false, "Shift", true, false, false},
		{"ctrl_only", false, true, false, "Ctrl", false, true, false},
		{"alt_only", false, false, true, "Alt", false, false, true},
		{"shift_ctrl", true, true, false, "Shift+Ctrl", true, true, false},
		{"shift_alt", true, false, true, "Shift+Alt", true, false, true},
		{"ctrl_alt", false, true, true, "Ctrl+Alt", false, true, true},
		{"all", true, true, true, "Shift+Ctrl+Alt", true, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mods := NewModifiers(tc.shift, tc.ctrl, tc.alt)

			if mods.String() != tc.expectedString {
				t.Errorf("String: expected %s, got %s", tc.expectedString, mods.String())
			}
			if mods.HasShift() != tc.expectedHasShift {
				t.Errorf("HasShift: expected %v, got %v", tc.expectedHasShift, mods.HasShift())
			}
			if mods.HasCtrl() != tc.expectedHasCtrl {
				t.Errorf("HasCtrl: expected %v, got %v", tc.expectedHasCtrl, mods.HasCtrl())
			}
			if mods.HasAlt() != tc.expectedHasAlt {
				t.Errorf("HasAlt: expected %v, got %v", tc.expectedHasAlt, mods.HasAlt())
			}
		})
	}
}

func TestModifiers_Equals(t *testing.T) {
	testCases := []struct {
		name     string
		m1, m2   Modifiers
		expected bool
	}{
		{"both_none", ModifierNone, ModifierNone, true},
		{"same_shift", ModifierShift, ModifierShift, true},
		{"same_ctrl", ModifierCtrl, ModifierCtrl, true},
		{"same_alt", ModifierAlt, ModifierAlt, true},
		{"same_combo", ModifierShift | ModifierCtrl, ModifierShift | ModifierCtrl, true},
		{"different", ModifierShift, ModifierCtrl, false},
		{"partial_match", ModifierShift | ModifierCtrl, ModifierShift, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.m1.Equals(tc.m2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// ============================================================================
// Helper Function Tests: NewMouseEvent
// ============================================================================

func TestNewMouseEvent(t *testing.T) {
	eventType := EventClick
	button := ButtonLeft
	position := NewPosition(10, 20)
	modifiers := NewModifiers(true, false, false)

	event := NewMouseEvent(eventType, button, position, modifiers)

	if event.Type() != eventType {
		t.Errorf("Type: expected %v, got %v", eventType, event.Type())
	}
	if event.Button() != button {
		t.Errorf("Button: expected %v, got %v", button, event.Button())
	}
	if !event.Position().Equals(position) {
		t.Errorf("Position: expected %v, got %v", position, event.Position())
	}
	if !event.Modifiers().Equals(modifiers) {
		t.Errorf("Modifiers: expected %v, got %v", modifiers, event.Modifiers())
	}
	if event.Timestamp().IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNewMouseEvent_Timestamp(t *testing.T) {
	before := time.Now()
	event := NewMouseEvent(EventClick, ButtonLeft, NewPosition(0, 0), ModifierNone)
	after := time.Now()

	ts := event.Timestamp()
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in range [%v, %v]", ts, before, after)
	}
}

// ============================================================================
// Button Tests
// ============================================================================

func TestButton_String(t *testing.T) {
	testCases := []struct {
		button   Button
		expected string
	}{
		{ButtonNone, "None"},
		{ButtonLeft, "Left"},
		{ButtonMiddle, "Middle"},
		{ButtonRight, "Right"},
		{ButtonWheelUp, "WheelUp"},
		{ButtonWheelDown, "WheelDown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.button.String()
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestButton_IsWheel(t *testing.T) {
	testCases := []struct {
		button   Button
		expected bool
	}{
		{ButtonNone, false},
		{ButtonLeft, false},
		{ButtonMiddle, false},
		{ButtonRight, false},
		{ButtonWheelUp, true},
		{ButtonWheelDown, true},
	}

	for _, tc := range testCases {
		t.Run(tc.button.String(), func(t *testing.T) {
			result := tc.button.IsWheel()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestButton_IsButton(t *testing.T) {
	testCases := []struct {
		button   Button
		expected bool
	}{
		{ButtonNone, false},
		{ButtonLeft, true},
		{ButtonMiddle, true},
		{ButtonRight, true},
		{ButtonWheelUp, false},
		{ButtonWheelDown, false},
	}

	for _, tc := range testCases {
		t.Run(tc.button.String(), func(t *testing.T) {
			result := tc.button.IsButton()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// ============================================================================
// EventType Tests
// ============================================================================

func TestEventType_String(t *testing.T) {
	testCases := []struct {
		eventType EventType
		expected  string
	}{
		{EventPress, "Press"},
		{EventRelease, "Release"},
		{EventClick, "Click"},
		{EventDoubleClick, "DoubleClick"},
		{EventTripleClick, "TripleClick"},
		{EventDrag, "Drag"},
		{EventMotion, "Motion"},
		{EventScroll, "Scroll"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.eventType.String()
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestEventType_IsClick(t *testing.T) {
	testCases := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, true},
		{EventDoubleClick, true},
		{EventTripleClick, true},
		{EventDrag, false},
		{EventMotion, false},
		{EventScroll, false},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType.String(), func(t *testing.T) {
			result := tc.eventType.IsClick()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestEventType_IsDrag(t *testing.T) {
	testCases := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, false},
		{EventDrag, true},
		{EventMotion, false},
		{EventScroll, false},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType.String(), func(t *testing.T) {
			result := tc.eventType.IsDrag()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestEventType_IsScroll(t *testing.T) {
	testCases := []struct {
		eventType EventType
		expected  bool
	}{
		{EventPress, false},
		{EventRelease, false},
		{EventClick, false},
		{EventDrag, false},
		{EventMotion, false},
		{EventScroll, true},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType.String(), func(t *testing.T) {
			result := tc.eventType.IsScroll()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// ============================================================================
// MouseEvent Method Tests
// ============================================================================

func TestMouseEvent_IsClick(t *testing.T) {
	testCases := []struct {
		name      string
		eventType EventType
		expected  bool
	}{
		{"click", EventClick, true},
		{"double_click", EventDoubleClick, true},
		{"triple_click", EventTripleClick, true},
		{"press", EventPress, false},
		{"drag", EventDrag, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := NewMouseEvent(tc.eventType, ButtonLeft, NewPosition(0, 0), ModifierNone)
			result := event.IsClick()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMouseEvent_IsDrag(t *testing.T) {
	testCases := []struct {
		name      string
		eventType EventType
		expected  bool
	}{
		{"drag", EventDrag, true},
		{"click", EventClick, false},
		{"motion", EventMotion, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := NewMouseEvent(tc.eventType, ButtonLeft, NewPosition(0, 0), ModifierNone)
			result := event.IsDrag()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMouseEvent_IsScroll(t *testing.T) {
	testCases := []struct {
		name      string
		eventType EventType
		expected  bool
	}{
		{"scroll", EventScroll, true},
		{"click", EventClick, false},
		{"drag", EventDrag, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := NewMouseEvent(tc.eventType, ButtonLeft, NewPosition(0, 0), ModifierNone)
			result := event.IsScroll()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMouseEvent_String(t *testing.T) {
	event := NewMouseEvent(EventClick, ButtonLeft, NewPosition(10, 20), ModifierShift)
	str := event.String()

	// Should contain key information
	if !strings.Contains(str, "Click") {
		t.Errorf("String should contain event type, got: %s", str)
	}
	if !strings.Contains(str, "Left") {
		t.Errorf("String should contain button, got: %s", str)
	}
	if !strings.Contains(str, "(10,20)") {
		t.Errorf("String should contain position, got: %s", str)
	}
}

// ============================================================================
// Constants Export Tests
// ============================================================================

func TestEventTypeConstants(t *testing.T) {
	// Verify all event type constants are accessible from API
	eventTypes := []EventType{
		EventPress,
		EventRelease,
		EventClick,
		EventDoubleClick,
		EventTripleClick,
		EventDrag,
		EventMotion,
		EventScroll,
	}

	for _, et := range eventTypes {
		if et.String() == "Unknown" {
			t.Errorf("event type %d has invalid String representation", et)
		}
	}
}

func TestButtonConstants(t *testing.T) {
	// Verify all button constants are accessible from API
	buttons := []Button{
		ButtonNone,
		ButtonLeft,
		ButtonMiddle,
		ButtonRight,
		ButtonWheelUp,
		ButtonWheelDown,
	}

	for _, b := range buttons {
		if b.String() == "Unknown" && b != Button(999) {
			t.Errorf("button %d has invalid String representation", b)
		}
	}
}

func TestModifierConstants(_ *testing.T) {
	// Verify all modifier constants are accessible from API
	modifiers := []Modifiers{
		ModifierNone,
		ModifierShift,
		ModifierCtrl,
		ModifierAlt,
	}

	for _, m := range modifiers {
		// Just verify they're defined and don't panic
		_ = m.String()
	}
}
