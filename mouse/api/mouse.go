// Package api provides the public API for the phoenix/mouse library.
// It offers a clean, fluent interface for mouse event handling in terminal applications.
package api

import (
	"github.com/phoenix-tui/phoenix/mouse/application"
	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// Re-export domain types for convenience.
type (
	// MouseEvent represents a mouse event.
	MouseEvent = model.MouseEvent

	// EventType represents the type of mouse event.
	EventType = value.EventType

	// Button represents a mouse button.
	Button = value.Button

	// Position represents a mouse position.
	Position = value.Position

	// Modifiers represents keyboard modifiers.
	Modifiers = value.Modifiers
)

// Event types.
const (
	EventPress       = value.EventPress
	EventRelease     = value.EventRelease
	EventClick       = value.EventClick
	EventDoubleClick = value.EventDoubleClick
	EventTripleClick = value.EventTripleClick
	EventDrag        = value.EventDrag
	EventMotion      = value.EventMotion
	EventScroll      = value.EventScroll
)

// Buttons.
const (
	ButtonNone      = value.ButtonNone
	ButtonLeft      = value.ButtonLeft
	ButtonMiddle    = value.ButtonMiddle
	ButtonRight     = value.ButtonRight
	ButtonWheelUp   = value.ButtonWheelUp
	ButtonWheelDown = value.ButtonWheelDown
)

// Modifiers.
const (
	ModifierNone  = value.ModifierNone
	ModifierShift = value.ModifierShift
	ModifierCtrl  = value.ModifierCtrl
	ModifierAlt   = value.ModifierAlt
)

// Mouse is the main API for mouse handling.
// It provides a simple, fluent interface for enabling mouse support
// and processing mouse events.
type Mouse struct {
	handler *application.MouseHandler
}

// New creates a new Mouse instance.
func New() *Mouse {
	return &Mouse{
		handler: application.NewMouseHandler(),
	}
}

// Enable enables mouse tracking.
// This writes ANSI escape sequences to stdout to enable mouse reporting.
func (m *Mouse) Enable() error {
	return m.handler.Enable()
}

// Disable disables mouse tracking.
// This should be called on program exit or when mouse support is no longer needed.
func (m *Mouse) Disable() error {
	return m.handler.Disable()
}

// IsEnabled returns true if mouse tracking is currently enabled.
func (m *Mouse) IsEnabled() bool {
	return m.handler.IsEnabled()
}

// ParseSequence parses a mouse input sequence from the terminal.
// The sequence should be the raw ANSI escape sequence (without the ESC prefix).
// Returns enriched events with click detection, drag tracking, etc.
//
// Example sequences:
//   - SGR: "\x1b[<0;10;5M" (left button press at 10,5)
//   - X10: "\x1b[M !!    " (legacy format)
//
// The returned events may include:
//   - Press/Release events (raw button events)
//   - Click/DoubleClick/TripleClick events (detected clicks)
//   - Drag events (motion with button pressed)
//   - Scroll events (mouse wheel)
func (m *Mouse) ParseSequence(sequence string) ([]MouseEvent, error) {
	return m.handler.ParseSequence(sequence)
}

// ScrollDelta calculates the scroll delta (in lines) for a scroll event.
// Returns negative for scroll up, positive for scroll down.
func (m *Mouse) ScrollDelta(event MouseEvent) int {
	return m.handler.Processor().ScrollDelta(event)
}

// IsDragging returns true if a drag is currently in progress.
func (m *Mouse) IsDragging() bool {
	return m.handler.Processor().IsDragging()
}

// Reset resets the mouse handler state (useful for testing).
func (m *Mouse) Reset() {
	m.handler.Reset()
}

// Helper functions for creating values

// NewPosition creates a new Position.
func NewPosition(x, y int) Position {
	return value.NewPosition(x, y)
}

// NewModifiers creates a new Modifiers value.
func NewModifiers(shift, ctrl, alt bool) Modifiers {
	return value.NewModifiers(shift, ctrl, alt)
}

// NewMouseEvent creates a new MouseEvent.
func NewMouseEvent(eventType EventType, button Button, position Position, modifiers Modifiers) MouseEvent {
	return model.NewMouseEvent(eventType, button, position, modifiers)
}
