package model

import (
	"fmt"
	"time"

	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// MouseEvent is the aggregate root representing a mouse event.
// It encapsulates all information about a single mouse interaction.
type MouseEvent struct {
	eventType value.EventType
	button    value.Button
	position  value.Position
	modifiers value.Modifiers
	timestamp time.Time
}

// NewMouseEvent creates a new MouseEvent.
func NewMouseEvent(
	eventType value.EventType,
	button value.Button,
	position value.Position,
	modifiers value.Modifiers,
) MouseEvent {
	return MouseEvent{
		eventType: eventType,
		button:    button,
		position:  position,
		modifiers: modifiers,
		timestamp: time.Now(),
	}
}

// NewMouseEventWithTimestamp creates a new MouseEvent with an explicit timestamp.
func NewMouseEventWithTimestamp(
	eventType value.EventType,
	button value.Button,
	position value.Position,
	modifiers value.Modifiers,
	timestamp time.Time,
) MouseEvent {
	return MouseEvent{
		eventType: eventType,
		button:    button,
		position:  position,
		modifiers: modifiers,
		timestamp: timestamp,
	}
}

// Type returns the event type.
func (e MouseEvent) Type() value.EventType {
	return e.eventType
}

// Button returns the button.
func (e MouseEvent) Button() value.Button {
	return e.button
}

// Position returns the position.
func (e MouseEvent) Position() value.Position {
	return e.position
}

// Modifiers returns the modifiers.
func (e MouseEvent) Modifiers() value.Modifiers {
	return e.modifiers
}

// Timestamp returns the timestamp.
func (e MouseEvent) Timestamp() time.Time {
	return e.timestamp
}

// String returns the string representation of the mouse event.
func (e MouseEvent) String() string {
	return fmt.Sprintf(
		"MouseEvent{type=%s, button=%s, pos=%s, mods=%s, time=%s}",
		e.eventType,
		e.button,
		e.position,
		e.modifiers,
		e.timestamp.Format("15:04:05.000"),
	)
}

// WithType returns a new MouseEvent with the specified type.
func (e MouseEvent) WithType(eventType value.EventType) MouseEvent {
	return MouseEvent{
		eventType: eventType,
		button:    e.button,
		position:  e.position,
		modifiers: e.modifiers,
		timestamp: e.timestamp,
	}
}

// IsAt checks if the event is at the specified position (with tolerance).
func (e MouseEvent) IsAt(pos value.Position, tolerance int) bool {
	return e.position.IsWithinTolerance(pos, tolerance)
}

// IsClick returns true if this is a click-related event.
func (e MouseEvent) IsClick() bool {
	return e.eventType.IsClick()
}

// IsDrag returns true if this is a drag event.
func (e MouseEvent) IsDrag() bool {
	return e.eventType.IsDrag()
}

// IsScroll returns true if this is a scroll event.
func (e MouseEvent) IsScroll() bool {
	return e.eventType.IsScroll()
}
