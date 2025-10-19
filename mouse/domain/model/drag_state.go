// Package model contains domain models (entities and aggregates) for mouse handling.
// Models have identity and mutable state, unlike value objects.
package model

import (
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// DragState represents the state of a drag operation.
// This is an entity that tracks drag progress.
type DragState struct {
	active    bool
	start     value.Position
	current   value.Position
	button    value.Button
	modifiers value.Modifiers
	threshold int // Minimum distance to consider it a drag
}

// NewDragState creates a new DragState with the specified threshold.
func NewDragState(threshold int) *DragState {
	if threshold <= 0 {
		threshold = 2 // Default threshold
	}
	return &DragState{
		active:    false,
		threshold: threshold,
	}
}

// Start begins a drag operation.
func (d *DragState) Start(pos value.Position, button value.Button, modifiers value.Modifiers) {
	d.active = true
	d.start = pos
	d.current = pos
	d.button = button
	d.modifiers = modifiers
}

// Update updates the current position during a drag.
func (d *DragState) Update(pos value.Position) {
	if d.active {
		d.current = pos
	}
}

// End ends the drag operation.
func (d *DragState) End() {
	d.active = false
}

// IsActive returns true if a drag is currently active.
func (d *DragState) IsActive() bool {
	return d.active
}

// IsDrag returns true if the movement is beyond the threshold (actual drag).
func (d *DragState) IsDrag() bool {
	if !d.active {
		return false
	}
	distance := d.start.DistanceTo(d.current)
	return distance >= d.threshold
}

// StartPosition returns the starting position of the drag.
func (d *DragState) StartPosition() value.Position {
	return d.start
}

// Current returns the current position of the drag.
func (d *DragState) Current() value.Position {
	return d.current
}

// Button returns the button used for the drag.
func (d *DragState) Button() value.Button {
	return d.button
}

// Modifiers returns the modifiers held during the drag.
func (d *DragState) Modifiers() value.Modifiers {
	return d.modifiers
}

// Distance returns the distance from start to current position.
func (d *DragState) Distance() int {
	if !d.active {
		return 0
	}
	return d.start.DistanceTo(d.current)
}

// Reset resets the drag state to inactive.
func (d *DragState) Reset() {
	d.active = false
	d.start = value.NewPosition(0, 0)
	d.current = value.NewPosition(0, 0)
	d.button = value.ButtonNone
	d.modifiers = value.ModifierNone
}
