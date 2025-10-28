package parser

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
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
	position := value2.NewPosition(x, y)

	// X10 doesn't distinguish press/release for normal buttons
	// We'll treat it as a press event
	var eventType value2.EventType
	if button.IsWheel() {
		eventType = value2.EventScroll
	} else {
		eventType = value2.EventPress
	}

	event := model.NewMouseEvent(eventType, button, position, modifiers)
	return &event, nil
}

// decodeButton extracts button and modifiers from X10 button code.
func (p *X10Parser) decodeButton(code int) (value2.Button, value2.Modifiers) {
	// Extract modifiers
	shift := (code & 4) != 0
	alt := (code & 8) != 0
	ctrl := (code & 16) != 0
	modifiers := value2.NewModifiers(shift, ctrl, alt)

	// Extract base button (remove modifier bits)
	baseButton := code & 0x63 // bits 0-1 and 5-6

	var button value2.Button
	switch baseButton {
	case 0:
		button = value2.ButtonLeft
	case 1:
		button = value2.ButtonMiddle
	case 2:
		button = value2.ButtonRight
	case 64:
		button = value2.ButtonWheelUp
	case 65:
		button = value2.ButtonWheelDown
	case 32, 35: // Motion events
		button = value2.ButtonNone
	default:
		button = value2.ButtonNone
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
func (p *X10Parser) encodeButton(button value2.Button, modifiers value2.Modifiers) int {
	var code int

	// Base button
	switch button {
	case value2.ButtonLeft:
		code = 0
	case value2.ButtonMiddle:
		code = 1
	case value2.ButtonRight:
		code = 2
	case value2.ButtonWheelUp:
		code = 64
	case value2.ButtonWheelDown:
		code = 65
	case value2.ButtonNone:
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
