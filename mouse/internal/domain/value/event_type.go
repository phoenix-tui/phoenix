package value

// EventType represents the type of mouse event.
type EventType int

const (
	// EventPress represents a button press (start of click or drag).
	EventPress EventType = iota
	// EventRelease represents a button release (end of click or drag).
	EventRelease
	// EventClick represents a single click (press + release at same position).
	EventClick
	// EventDoubleClick represents a double click (two clicks within timeout).
	EventDoubleClick
	// EventTripleClick represents a triple click (three clicks within timeout).
	EventTripleClick
	// EventDrag represents mouse drag (motion with button pressed).
	EventDrag
	// EventMotion represents mouse motion (no button pressed).
	EventMotion
	// EventScroll represents scroll wheel action.
	EventScroll
	// EventHoverEnter represents mouse entering a component area.
	EventHoverEnter
	// EventHoverLeave represents mouse leaving a component area.
	EventHoverLeave
	// EventHoverMove represents mouse moving within a component area.
	EventHoverMove
)

// String returns the string representation of the event type.
func (e EventType) String() string {
	switch e {
	case EventPress:
		return "Press"
	case EventRelease:
		return "Release"
	case EventClick:
		return "Click"
	case EventDoubleClick:
		return "DoubleClick"
	case EventTripleClick:
		return "TripleClick"
	case EventDrag:
		return "Drag"
	case EventMotion:
		return "Motion"
	case EventScroll:
		return "Scroll"
	case EventHoverEnter:
		return "HoverEnter"
	case EventHoverLeave:
		return "HoverLeave"
	case EventHoverMove:
		return "HoverMove"
	default:
		return "Unknown"
	}
}

// IsClick returns true if the event is a click-related event.
func (e EventType) IsClick() bool {
	return e == EventClick || e == EventDoubleClick || e == EventTripleClick
}

// IsDrag returns true if the event is a drag event.
func (e EventType) IsDrag() bool {
	return e == EventDrag
}

// IsScroll returns true if the event is a scroll event.
func (e EventType) IsScroll() bool {
	return e == EventScroll
}

// IsHover returns true if the event is a hover-related event.
func (e EventType) IsHover() bool {
	return e == EventHoverEnter || e == EventHoverLeave || e == EventHoverMove
}
