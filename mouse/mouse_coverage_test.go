package mouse

import (
	"testing"
)

// ============================================================================
// Hover Event Processing Tests (Week 15 Day 8 - Coverage Sprint)
// ============================================================================

func TestMouse_ProcessHover_HoverEnter(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	areas := []ComponentArea{
		{ID: "button1", Area: NewBoundingBox(5, 10, 20, 3)},
	}

	pos := NewPosition(10, 11)
	eventType := m.ProcessHover(pos, areas)

	if eventType != EventHoverEnter {
		t.Errorf("expected EventHoverEnter, got %v", eventType)
	}
	if !m.IsHovering() {
		t.Error("IsHovering() should return true after hover enter")
	}
	if m.CurrentHoverComponent() != "button1" {
		t.Errorf("expected component ID 'button1', got %q", m.CurrentHoverComponent())
	}
}

func TestMouse_ProcessHover_HoverLeave(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	areas := []ComponentArea{
		{ID: "button1", Area: NewBoundingBox(5, 10, 20, 3)},
	}

	pos := NewPosition(10, 11)
	m.ProcessHover(pos, areas)

	posOutside := NewPosition(1, 1)
	eventType := m.ProcessHover(posOutside, areas)

	if eventType != EventHoverLeave {
		t.Errorf("expected EventHoverLeave, got %v", eventType)
	}
	if m.IsHovering() {
		t.Error("IsHovering() should return false after hover leave")
	}
	if m.CurrentHoverComponent() != "" {
		t.Errorf("expected empty component ID, got %q", m.CurrentHoverComponent())
	}
}

func TestMouse_ProcessHover_HoverMove(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	areas := []ComponentArea{
		{ID: "button1", Area: NewBoundingBox(5, 10, 20, 3)},
	}

	pos1 := NewPosition(10, 11)
	m.ProcessHover(pos1, areas)

	pos2 := NewPosition(15, 11)
	eventType := m.ProcessHover(pos2, areas)

	if eventType != EventHoverMove {
		t.Errorf("expected EventHoverMove, got %v", eventType)
	}
	if !m.IsHovering() {
		t.Error("IsHovering() should still return true after hover move")
	}
	if m.CurrentHoverComponent() != "button1" {
		t.Errorf("expected component ID 'button1', got %q", m.CurrentHoverComponent())
	}
}

func TestMouse_ProcessHover_EmptyAreas(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	pos := NewPosition(10, 11)
	eventType := m.ProcessHover(pos, []ComponentArea{})

	if eventType != EventMotion {
		t.Errorf("expected EventMotion for empty areas, got %v", eventType)
	}
	if m.IsHovering() {
		t.Error("IsHovering() should be false with empty areas")
	}
}

func TestMouse_ProcessHover_MultipleOverlappingAreas(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	areas := []ComponentArea{
		{ID: "button1", Area: NewBoundingBox(5, 10, 20, 5)},
		{ID: "button2", Area: NewBoundingBox(10, 11, 10, 3)},
	}

	pos := NewPosition(12, 12)
	eventType := m.ProcessHover(pos, areas)

	if eventType != EventHoverEnter {
		t.Errorf("expected EventHoverEnter, got %v", eventType)
	}
	if m.CurrentHoverComponent() != "button1" {
		t.Errorf("expected component ID 'button1', got %q", m.CurrentHoverComponent())
	}
}

func TestMouse_IsHovering_InitiallyFalse(t *testing.T) {
	m := New()
	if m.IsHovering() {
		t.Error("IsHovering() should be false initially")
	}
}

func TestMouse_CurrentHoverComponent_InitiallyEmpty(t *testing.T) {
	m := New()
	if m.CurrentHoverComponent() != "" {
		t.Errorf("CurrentHoverComponent() should be empty initially, got %q", m.CurrentHoverComponent())
	}
}

func TestMouse_CurrentHoverComponent_Persistence(t *testing.T) {
	m := New()
	m.Enable()
	defer m.Disable()

	areas := []ComponentArea{
		{ID: "persistent-button", Area: NewBoundingBox(5, 10, 20, 3)},
	}

	pos := NewPosition(10, 11)
	m.ProcessHover(pos, areas)

	for i := 0; i < 3; i++ {
		if m.CurrentHoverComponent() != "persistent-button" {
			t.Errorf("iteration %d: expected 'persistent-button', got %q", i, m.CurrentHoverComponent())
		}
	}
}

// ============================================================================
// Menu Positioning Tests (Week 15 Day 8 - Coverage Sprint)
// ============================================================================

func TestMouse_CalculateMenuPosition_NormalPositioning(t *testing.T) {
	m := New()

	cursorPos := NewPosition(20, 10)
	menuWidth, menuHeight := 15, 5
	screenWidth, screenHeight := 80, 24

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	if resultPos.X() != 20 {
		t.Errorf("expected X=20, got %d", resultPos.X())
	}
	if resultPos.Y() != 10 {
		t.Errorf("expected Y=10, got %d", resultPos.Y())
	}
}

func TestMouse_CalculateMenuPosition_RightEdgeAdjustment(t *testing.T) {
	m := New()

	cursorPos := NewPosition(70, 10)
	menuWidth, menuHeight := 25, 5
	screenWidth, screenHeight := 80, 24

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	if resultPos.X()+menuWidth > screenWidth {
		t.Errorf("menu overflows right edge: X=%d, width=%d, screen=%d", resultPos.X(), menuWidth, screenWidth)
	}
}

func TestMouse_CalculateMenuPosition_BottomEdgeAdjustment(t *testing.T) {
	m := New()

	cursorPos := NewPosition(20, 20)
	menuWidth, menuHeight := 15, 8
	screenWidth, screenHeight := 80, 24

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	if resultPos.Y()+menuHeight > screenHeight {
		t.Errorf("menu overflows bottom edge: Y=%d, height=%d, screen=%d", resultPos.Y(), menuHeight, screenHeight)
	}
}

func TestMouse_CalculateMenuPosition_CornerAdjustment(t *testing.T) {
	m := New()

	cursorPos := NewPosition(75, 22)
	menuWidth, menuHeight := 20, 6
	screenWidth, screenHeight := 80, 24

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	if resultPos.X()+menuWidth > screenWidth {
		t.Errorf("menu overflows right edge: X=%d, width=%d, screen=%d", resultPos.X(), menuWidth, screenWidth)
	}
	if resultPos.Y()+menuHeight > screenHeight {
		t.Errorf("menu overflows bottom edge: Y=%d, height=%d, screen=%d", resultPos.Y(), menuHeight, screenHeight)
	}
}

func TestMouse_CalculateMenuPosition_OversizedMenu(t *testing.T) {
	m := New()

	cursorPos := NewPosition(10, 10)
	menuWidth, menuHeight := 100, 30
	screenWidth, screenHeight := 80, 24

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	if resultPos.X() < 0 || resultPos.Y() < 0 {
		t.Errorf("position should not be negative: X=%d, Y=%d", resultPos.X(), resultPos.Y())
	}
}

func TestMouse_CalculateMenuPosition_ZeroScreenSize(t *testing.T) {
	m := New()

	cursorPos := NewPosition(10, 10)
	menuWidth, menuHeight := 15, 5
	screenWidth, screenHeight := 0, 0

	resultPos := m.CalculateMenuPosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)
	_ = resultPos
}

// ============================================================================
// NewBoundingBox Tests (Week 15 Day 8 - Coverage Sprint)
// ============================================================================

func TestNewBoundingBox_NormalCreation(t *testing.T) {
	bb := NewBoundingBox(10, 20, 30, 40)

	if bb.X() != 10 {
		t.Errorf("expected X=10, got %d", bb.X())
	}
	if bb.Y() != 20 {
		t.Errorf("expected Y=20, got %d", bb.Y())
	}
	if bb.Width() != 30 {
		t.Errorf("expected Width=30, got %d", bb.Width())
	}
	if bb.Height() != 40 {
		t.Errorf("expected Height=40, got %d", bb.Height())
	}
}

func TestNewBoundingBox_ZeroDimensions(t *testing.T) {
	bb := NewBoundingBox(0, 0, 0, 0)

	if bb.X() != 0 || bb.Y() != 0 || bb.Width() != 0 || bb.Height() != 0 {
		t.Error("expected all zero dimensions")
	}
}

func TestNewBoundingBox_NegativeDimensions(t *testing.T) {
	bb := NewBoundingBox(-5, -10, -20, -30)

	// Negative dimensions are normalized to zero
	if bb.X() != -5 || bb.Y() != -10 {
		t.Error("expected negative position coordinates to be preserved")
	}
	if bb.Width() != 0 || bb.Height() != 0 {
		t.Error("expected negative dimensions to be normalized to zero")
	}
}
func TestBoundingBox_Contains(t *testing.T) {
	bb := NewBoundingBox(10, 10, 20, 10)

	testCases := []struct {
		name     string
		pos      Position
		expected bool
	}{
		{"inside", NewPosition(15, 15), true},
		{"top_left_corner", NewPosition(10, 10), true},
		{"bottom_right_corner", NewPosition(29, 19), true},
		{"outside_left", NewPosition(5, 15), false},
		{"outside_right", NewPosition(35, 15), false},
		{"outside_above", NewPosition(15, 5), false},
		{"outside_below", NewPosition(15, 25), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := bb.Contains(tc.pos)
			if result != tc.expected {
				t.Errorf("expected %v, got %v for position %v", tc.expected, result, tc.pos)
			}
		})
	}
}
