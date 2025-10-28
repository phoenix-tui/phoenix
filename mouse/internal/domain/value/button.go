// Package value contains value objects for the mouse domain.
// Value objects are immutable and represent domain concepts like Position, Button, Modifiers.
package value

// Button represents a mouse button.
type Button int

const (
	// ButtonNone represents no button (motion only).
	ButtonNone Button = iota
	// ButtonLeft represents the left mouse button.
	ButtonLeft
	// ButtonMiddle represents the middle mouse button (scroll wheel press).
	ButtonMiddle
	// ButtonRight represents the right mouse button.
	ButtonRight
	// ButtonWheelUp represents scroll wheel up.
	ButtonWheelUp
	// ButtonWheelDown represents scroll wheel down.
	ButtonWheelDown
)

// String returns the string representation of the button.
func (b Button) String() string {
	switch b {
	case ButtonNone:
		return "None"
	case ButtonLeft:
		return "Left"
	case ButtonMiddle:
		return "Middle"
	case ButtonRight:
		return "Right"
	case ButtonWheelUp:
		return "WheelUp"
	case ButtonWheelDown:
		return "WheelDown"
	default:
		return "Unknown"
	}
}

// IsWheel returns true if the button is a scroll wheel action.
func (b Button) IsWheel() bool {
	return b == ButtonWheelUp || b == ButtonWheelDown
}

// IsButton returns true if the button is an actual button (not wheel).
func (b Button) IsButton() bool {
	return b >= ButtonLeft && b <= ButtonRight
}
