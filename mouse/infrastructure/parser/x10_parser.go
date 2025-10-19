package parser

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// X10Parser parses X10 (1000) mouse protocol sequences.
// Format: \x1b[M<button><x><y>
// button, x, y are single bytes with value = byte - 32
// Limited to 223x223 terminal size.
type X10Parser struct{}

// NewX10Parser creates a new X10 parser.
func NewX10Parser() *X10Parser {
	return &X10Parser{}
}

// Parse parses an X10 mouse sequence.
// Input: 3 bytes (button, x, y) where each byte value = actual_value + 32.
func (p *X10Parser) Parse(data []byte) (*model.MouseEvent, error) {
	if len(data) != 3 {
		return nil, fmt.Errorf("invalid X10 sequence: expected 3 bytes, got %d", len(data))
	}

	// Decode button (subtract 32 offset)
	buttonCode := int(data[0]) - 32
	x := int(data[1]) - 32 - 1 // -1 to convert to 0-based
	y := int(data[2]) - 32 - 1 // -1 to convert to 0-based

	// Extract button and modifiers
	button, modifiers := p.decodeButton(buttonCode)
	position := value.NewPosition(x, y)

	// X10 doesn't distinguish press/release for normal buttons
	// We'll treat it as a press event
	var eventType value.EventType
	if button.IsWheel() {
		eventType = value.EventScroll
	} else {
		eventType = value.EventPress
	}

	event := model.NewMouseEvent(eventType, button, position, modifiers)
	return &event, nil
}

// decodeButton extracts button and modifiers from X10 button code.
func (p *X10Parser) decodeButton(code int) (value.Button, value.Modifiers) {
	// Extract modifiers
	shift := (code & 4) != 0
	alt := (code & 8) != 0
	ctrl := (code & 16) != 0
	modifiers := value.NewModifiers(shift, ctrl, alt)

	// Extract base button (remove modifier bits)
	baseButton := code & 0x63 // bits 0-1 and 5-6

	var button value.Button
	switch baseButton {
	case 0:
		button = value.ButtonLeft
	case 1:
		button = value.ButtonMiddle
	case 2:
		button = value.ButtonRight
	case 64:
		button = value.ButtonWheelUp
	case 65:
		button = value.ButtonWheelDown
	case 32, 35: // Motion events
		button = value.ButtonNone
	default:
		button = value.ButtonNone
	}

	return button, modifiers
}

// FormatSequence formats a mouse event as an X10 sequence (for testing).
func (p *X10Parser) FormatSequence(event model.MouseEvent) string {
	// Encode button
	buttonCode := p.encodeButton(event.Button(), event.Modifiers())

	// Get position (convert to 1-based, then add 32 offset)
	x := byte(event.Position().X() + 1 + 32)
	y := byte(event.Position().Y() + 1 + 32)

	// Build raw byte sequence (avoid UTF-8 encoding issues with %c)
	return string([]byte{0x1b, '[', 'M', byte(buttonCode + 32), x, y})
}

// encodeButton encodes button and modifiers into X10 button code.
func (p *X10Parser) encodeButton(button value.Button, modifiers value.Modifiers) int {
	var code int

	// Base button
	switch button {
	case value.ButtonLeft:
		code = 0
	case value.ButtonMiddle:
		code = 1
	case value.ButtonRight:
		code = 2
	case value.ButtonWheelUp:
		code = 64
	case value.ButtonWheelDown:
		code = 65
	case value.ButtonNone:
		code = 32
	}

	// Add modifiers
	if modifiers.HasShift() {
		code |= 4
	}
	if modifiers.HasAlt() {
		code |= 8
	}
	if modifiers.HasCtrl() {
		code |= 16
	}

	return code
}
