package service

import (
	model2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// DragTracker is a domain service that tracks drag operations.
type DragTracker struct {
	state *model2.DragState
}

// NewDragTracker creates a new DragTracker.
// threshold: minimum distance (in cells) to consider it a drag (typically 2)
func NewDragTracker(threshold int) *DragTracker {
	return &DragTracker{
		state: model2.NewDragState(threshold),
	}
}

// ProcessPress handles a press event, potentially starting a drag.
func (d *DragTracker) ProcessPress(pressEvent model2.MouseEvent) {
	if pressEvent.Type() == value2.EventPress {
		d.state.Start(
			pressEvent.Position(),
			pressEvent.Button(),
			pressEvent.Modifiers(),
		)
	}
}

// ProcessMotion handles a motion event during a drag.
// Returns a drag event if the motion is beyond the threshold, nil otherwise.
func (d *DragTracker) ProcessMotion(motionEvent model2.MouseEvent) *model2.MouseEvent {
	if !d.state.IsActive() {
		return nil
	}

	d.state.Update(motionEvent.Position())

	// Only emit drag events if beyond threshold
	if d.state.IsDrag() {
		dragEvent := model2.NewMouseEventWithTimestamp(
			value2.EventDrag,
			d.state.Button(),
			motionEvent.Position(),
			d.state.Modifiers(),
			motionEvent.Timestamp(),
		)
		return &dragEvent
	}

	return nil
}

// ProcessRelease handles a release event, ending the drag.
// Returns the final drag state information.
func (d *DragTracker) ProcessRelease(_ model2.MouseEvent) (wasDrag bool, start, end value2.Position) {
	if !d.state.IsActive() {
		return false, value2.NewPosition(0, 0), value2.NewPosition(0, 0)
	}

	wasDrag = d.state.IsDrag()
	start = d.state.StartPosition()
	end = d.state.Current()

	d.state.End()
	return wasDrag, start, end
}

// IsActive returns true if a drag is currently active.
func (d *DragTracker) IsActive() bool {
	return d.state.IsActive()
}

// IsDrag returns true if the current drag is beyond the threshold.
func (d *DragTracker) IsDrag() bool {
	return d.state.IsDrag()
}

// State returns the current drag state (read-only).
func (d *DragTracker) State() *model2.DragState {
	return d.state
}

// Reset resets the drag tracker.
func (d *DragTracker) Reset() {
	d.state.Reset()
}
