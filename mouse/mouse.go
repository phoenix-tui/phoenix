// Package api provides the public API for the phoenix/mouse library.
// It offers a clean, fluent interface for mouse event handling in terminal applications.
package mouse

import (
	"github.com/phoenix-tui/phoenix/mouse/internal/application"
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// Re-export domain types for convenience.
type (
	// MouseEvent represents a mouse event.
	MouseEvent = model.MouseEvent

	// EventType represents the type of mouse event.
	EventType = value2.EventType

	// Button represents a mouse button.
	Button = value2.Button

	// Position represents a mouse position.
	Position = value2.Position

	// Modifiers represents keyboard modifiers.
	Modifiers = value2.Modifiers
)

// Event types.
const (
	EventPress       = value2.EventPress
	EventRelease     = value2.EventRelease
	EventClick       = value2.EventClick
	EventDoubleClick = value2.EventDoubleClick
	EventTripleClick = value2.EventTripleClick
	EventDrag        = value2.EventDrag
	EventMotion      = value2.EventMotion
	EventScroll      = value2.EventScroll
)

// Buttons.
const (
	ButtonNone      = value2.ButtonNone
	ButtonLeft      = value2.ButtonLeft
	ButtonMiddle    = value2.ButtonMiddle
	ButtonRight     = value2.ButtonRight
	ButtonWheelUp   = value2.ButtonWheelUp
	ButtonWheelDown = value2.ButtonWheelDown
)

// Modifiers.
const (
	ModifierNone  = value2.ModifierNone
	ModifierShift = value2.ModifierShift
	ModifierCtrl  = value2.ModifierCtrl
	ModifierAlt   = value2.ModifierAlt
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
	return value2.NewPosition(x, y)
}

// NewModifiers creates a new Modifiers value.
func NewModifiers(shift, ctrl, alt bool) Modifiers {
	return value2.NewModifiers(shift, ctrl, alt)
}

// NewMouseEvent creates a new MouseEvent.
func NewMouseEvent(eventType EventType, button Button, position Position, modifiers Modifiers) MouseEvent {
	return model.NewMouseEvent(eventType, button, position, modifiers)
}
