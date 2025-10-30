package service

import (
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// HoverTracker tracks hover state across component areas.
// This is a domain service that detects enter/leave/move events based on mouse position.
type HoverTracker struct {
	state *model.HoverState
}

// NewHoverTracker creates a new HoverTracker.
func NewHoverTracker() *HoverTracker {
	return &HoverTracker{
		state: model.NewHoverState(),
	}
}

// ComponentArea represents a component's hover-detection area.
type ComponentArea struct {
	ID   string
	Area value.BoundingBox
}

// Update processes a mouse position update and returns hover events.
// It checks all registered component areas to determine if the mouse has:
// - Entered a new component (HoverEnter)
// - Left the current component (HoverLeave)
// - Moved within the current component (HoverMove)
func (h *HoverTracker) Update(position value.Position, areas []ComponentArea) value.EventType {
	// Find which component (if any) contains the current position
	var hoveredComponentID string
	for _, area := range areas {
		if area.Area.Contains(position) {
			hoveredComponentID = area.ID
			break
		}
	}

	// Determine what happened
	wasHovering := h.state.IsHovering()
	previousComponentID := h.state.ComponentID()

	if hoveredComponentID == "" {
		// Mouse is not over any component
		if wasHovering {
			// Was hovering, now left
			h.state.Leave(position)
			return value.EventHoverLeave
		}
		// Was not hovering, still not hovering - no event
		return value.EventMotion
	}

	// Mouse is over a component
	if !wasHovering {
		// Was not hovering, now entered
		h.state.Enter(hoveredComponentID, position)
		return value.EventHoverEnter
	}

	if hoveredComponentID != previousComponentID {
		// Switched components - this is a leave + enter
		// For simplicity, we return HoverEnter (caller can detect component change)
		h.state.Enter(hoveredComponentID, position)
		return value.EventHoverEnter
	}

	// Still hovering over the same component - move event
	h.state.Move(position)
	return value.EventHoverMove
}

// State returns the current hover state (for testing/inspection).
func (h *HoverTracker) State() *model.HoverState {
	return h.state
}

// Reset clears the hover tracking state.
func (h *HoverTracker) Reset() {
	h.state.Reset()
}

// CurrentComponentID returns the ID of the currently hovered component (empty if none).
func (h *HoverTracker) CurrentComponentID() string {
	return h.state.ComponentID()
}

// IsHovering returns true if a component is currently being hovered.
func (h *HoverTracker) IsHovering() bool {
	return h.state.IsHovering()
}
