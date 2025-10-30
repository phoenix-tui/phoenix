package model

import (
	"time"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// HoverState tracks the current hover state for mouse interactions.
// This is an entity that maintains which component (if any) the mouse is currently hovering over.
type HoverState struct {
	componentID string         // ID of the currently hovered component (empty if none)
	position    value.Position // Current mouse position
	lastUpdate  time.Time      // When the hover state was last updated
	isActive    bool           // Whether hover tracking is active
}

// NewHoverState creates a new HoverState with no active hover.
func NewHoverState() *HoverState {
	return &HoverState{
		componentID: "",
		position:    value.NewPosition(0, 0),
		lastUpdate:  time.Now(),
		isActive:    false,
	}
}

// ComponentID returns the ID of the currently hovered component (empty if none).
func (h *HoverState) ComponentID() string {
	return h.componentID
}

// Position returns the current mouse position.
func (h *HoverState) Position() value.Position {
	return h.position
}

// LastUpdate returns the timestamp of the last hover state change.
func (h *HoverState) LastUpdate() time.Time {
	return h.lastUpdate
}

// IsActive returns true if hover tracking is currently active.
func (h *HoverState) IsActive() bool {
	return h.isActive
}

// IsHovering returns true if a component is currently being hovered.
func (h *HoverState) IsHovering() bool {
	return h.isActive && h.componentID != ""
}

// Enter updates the hover state when the mouse enters a component.
func (h *HoverState) Enter(componentID string, position value.Position) {
	h.componentID = componentID
	h.position = position
	h.lastUpdate = time.Now()
	h.isActive = true
}

// Move updates the hover state when the mouse moves within the same component.
func (h *HoverState) Move(position value.Position) {
	if !h.isActive {
		return
	}
	h.position = position
	h.lastUpdate = time.Now()
}

// Leave updates the hover state when the mouse leaves the current component.
func (h *HoverState) Leave(position value.Position) {
	if !h.isActive {
		return
	}
	h.componentID = ""
	h.position = position
	h.lastUpdate = time.Now()
	h.isActive = false
}

// Reset clears all hover state.
func (h *HoverState) Reset() {
	h.componentID = ""
	h.position = value.NewPosition(0, 0)
	h.lastUpdate = time.Now()
	h.isActive = false
}

// Equals returns true if this hover state is equal to another.
func (h *HoverState) Equals(other *HoverState) bool {
	if other == nil {
		return false
	}
	return h.componentID == other.componentID &&
		h.position.Equals(other.position) &&
		h.isActive == other.isActive
}
