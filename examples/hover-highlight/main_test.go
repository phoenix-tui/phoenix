package main

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse"
	"github.com/phoenix-tui/phoenix/tea"
)

// TestModelInit verifies initial model setup.
func TestModelInit(t *testing.T) {
	m := initialModel()

	if m.mouse == nil {
		t.Error("Expected mouse handler to be initialized")
	}

	if len(m.buttons) != 6 {
		t.Errorf("Expected 6 buttons, got %d", len(m.buttons))
	}

	if m.hoveredID != "" {
		t.Errorf("Expected empty hoveredID, got %q", m.hoveredID)
	}

	if m.lastClickedID != "" {
		t.Errorf("Expected empty lastClickedID, got %q", m.lastClickedID)
	}

	if m.ready {
		t.Error("Expected ready=false before WindowSizeMsg")
	}
}

// TestLayoutButtons verifies button layout calculation.
func TestLayoutButtons(t *testing.T) {
	m := initialModel()
	m.width = 80
	m.height = 24
	m = m.layoutButtons()

	// Verify all buttons have valid bounding boxes
	for i, btn := range m.buttons {
		if btn.area.Width() == 0 || btn.area.Height() == 0 {
			t.Errorf("Button %d has invalid dimensions", i+1)
		}

		if btn.id == "" {
			t.Errorf("Button %d has empty ID", i+1)
		}
	}

	// Verify buttons are within terminal bounds
	for i, btn := range m.buttons {
		if btn.area.X() < 0 || btn.area.Y() < 0 {
			t.Errorf("Button %d has negative position", i+1)
		}

		if btn.area.X()+btn.area.Width() > m.width {
			t.Errorf("Button %d exceeds terminal width", i+1)
		}

		if btn.area.Y()+btn.area.Height() > m.height {
			t.Errorf("Button %d exceeds terminal height", i+1)
		}
	}
}

// TestWindowSizeMsg verifies terminal resize handling.
func TestWindowSizeMsg(t *testing.T) {
	m := initialModel()

	// Send WindowSizeMsg
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updated, _ := m.Update(msg)

	if !updated.ready {
		t.Error("Expected ready=true after WindowSizeMsg")
	}

	if updated.width != 100 {
		t.Errorf("Expected width=100, got %d", updated.width)
	}

	if updated.height != 30 {
		t.Errorf("Expected height=30, got %d", updated.height)
	}
}

// TestMouseHoverDetection verifies hover event handling.
func TestMouseHoverDetection(t *testing.T) {
	m := initialModel()
	m.width = 80
	m.height = 24
	m.ready = true
	m = m.layoutButtons()

	// Simulate mouse hovering over button 1
	btn1 := m.buttons[0]
	centerX := btn1.area.X() + btn1.area.Width()/2
	centerY := btn1.area.Y() + btn1.area.Height()/2

	mouseMsg := tea.MouseMsg{
		X:      centerX,
		Y:      centerY,
		Button: tea.MouseButtonNone,
		Action: tea.MouseActionMotion,
	}

	updated := m.handleMouseEvent(mouseMsg)

	if updated.hoveredID != "button1" {
		t.Errorf("Expected hoveredID=button1, got %q", updated.hoveredID)
	}
}

// TestMouseClick verifies click event handling.
func TestMouseClick(t *testing.T) {
	m := initialModel()
	m.width = 80
	m.height = 24
	m.ready = true
	m = m.layoutButtons()

	// First hover over button 2
	btn2 := m.buttons[1]
	centerX := btn2.area.X() + btn2.area.Width()/2
	centerY := btn2.area.Y() + btn2.area.Height()/2

	hoverMsg := tea.MouseMsg{
		X:      centerX,
		Y:      centerY,
		Button: tea.MouseButtonNone,
		Action: tea.MouseActionMotion,
	}
	m = m.handleMouseEvent(hoverMsg)

	// Then click
	clickMsg := tea.MouseMsg{
		X:      centerX,
		Y:      centerY,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionRelease,
	}
	updated := m.handleMouseEvent(clickMsg)

	if updated.lastClickedID != "button2" {
		t.Errorf("Expected lastClickedID=button2, got %q", updated.lastClickedID)
	}
}

// TestKeyboardInput verifies keyboard event handling.
func TestKeyboardInput(t *testing.T) {
	m := initialModel()
	m.hoveredID = "button1"
	m.lastClickedID = "button2"

	// Test reset ('r' key)
	resetMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'r'}
	updated, _ := m.Update(resetMsg)

	if updated.hoveredID != "" {
		t.Errorf("Expected hoveredID to be cleared, got %q", updated.hoveredID)
	}

	if updated.lastClickedID != "" {
		t.Errorf("Expected lastClickedID to be cleared, got %q", updated.lastClickedID)
	}
}

// TestViewRendering verifies View() doesn't panic.
func TestViewRendering(t *testing.T) {
	m := initialModel()
	m.width = 80
	m.height = 24
	m.ready = true
	m = m.layoutButtons()

	// Should not panic
	view := m.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Basic checks
	if len(view) < 100 {
		t.Error("View seems too short")
	}
}

// TestBoundingBoxOverlap verifies buttons don't overlap.
func TestBoundingBoxOverlap(t *testing.T) {
	m := initialModel()
	m.width = 80
	m.height = 24
	m = m.layoutButtons()

	// Check all pairs of buttons for overlap
	for i := 0; i < len(m.buttons); i++ {
		for j := i + 1; j < len(m.buttons); j++ {
			box1 := m.buttons[i].area
			box2 := m.buttons[j].area

			// Check if boxes overlap
			if overlaps(box1, box2) {
				t.Errorf("Buttons %d and %d overlap", i+1, j+1)
			}
		}
	}
}

// overlaps checks if two bounding boxes overlap.
func overlaps(a, b mouse.BoundingBox) bool {
	return a.X() < b.X()+b.Width() &&
		a.X()+a.Width() > b.X() &&
		a.Y() < b.Y()+b.Height() &&
		a.Y()+a.Height() > b.Y()
}
