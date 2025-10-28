package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// URxvtParser parses URxvt (1015) mouse protocol sequences.
// Format: \x1b[button;x;yM
// Similar to SGR but always ends with 'M' (no press/release distinction).
type URxvtParser struct{}

// NewURxvtParser creates a new URxvt parser.
func NewURxvtParser() *URxvtParser {
	return &URxvtParser{}
}

// Parse parses a URxvt mouse sequence.
// Input format: "button;x;y" (no angle brackets, always ends with M).
func (p *URxvtParser) Parse(sequence string) (*model.MouseEvent, error) {
	// Split by semicolon
	parts := strings.Split(sequence, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid URxvt sequence: expected 3 parts, got %d", len(parts))
	}

	// Parse button
	buttonCode, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid button code: %w", err)
	}

	// Parse X (1-based in protocol, convert to 0-based)
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid X coordinate: %w", err)
	}
	x-- // Convert to 0-based

	// Parse Y (1-based in protocol, convert to 0-based)
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid Y coordinate: %w", err)
	}
	y-- // Convert to 0-based

	// Extract button and modifiers
	button, modifiers := p.decodeButton(buttonCode)
	position := value2.NewPosition(x, y)

	// URxvt doesn't distinguish press/release
	var eventType value2.EventType
	if button.IsWheel() {
		eventType = value2.EventScroll
	} else {
		eventType = value2.EventPress
	}

	event := model.NewMouseEvent(eventType, button, position, modifiers)
	return &event, nil
}

// decodeButton extracts button and modifiers from URxvt button code.
func (p *URxvtParser) decodeButton(code int) (value2.Button, value2.Modifiers) {
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

// FormatSequence formats a mouse event as a URxvt sequence (for testing).
func (p *URxvtParser) FormatSequence(event model.MouseEvent) string {
	// Encode button
	buttonCode := p.encodeButton(event.Button(), event.Modifiers())

	// Get position (convert to 1-based)
	x := event.Position().X() + 1
	y := event.Position().Y() + 1

	return fmt.Sprintf("\x1b[%d;%d;%dM", buttonCode, x, y)
}

// encodeButton encodes button and modifiers into URxvt button code.
func (p *URxvtParser) encodeButton(button value2.Button, modifiers value2.Modifiers) int {
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
