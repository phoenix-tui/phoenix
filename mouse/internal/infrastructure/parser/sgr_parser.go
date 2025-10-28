// Package parser contains parsers for different mouse protocol formats (SGR, X10, URxvt).
// Parsers convert raw ANSI sequences into domain MouseEvent objects.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// SGRParser parses SGR (1006) mouse protocol sequences.
// Format: \x1b[<button;x;y(M|m)
// M = press, m = release
// button: 0=left, 1=middle, 2=right, 64=wheel up, 65=wheel down, +4=shift, +8=alt, +16=ctrl, +32=motion
type SGRParser struct{}

// NewSGRParser creates a new SGR parser.
func NewSGRParser() *SGRParser {
	return &SGRParser{}
}

// Parse parses an SGR mouse sequence.
// Input format: "<button;x;y" followed by 'M' (press) or 'm' (release).
func (p *SGRParser) Parse(sequence string, isPress bool) (*model.MouseEvent, error) {
	// Remove leading "<" if present
	sequence = strings.TrimPrefix(sequence, "<")

	// Split by semicolon
	parts := strings.Split(sequence, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid SGR sequence: expected 3 parts, got %d", len(parts))
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

	// Determine event type
	var eventType value2.EventType
	if button.IsWheel() { //nolint:gocritic // ifElseChain: Event type determination requires conditional checks
		eventType = value2.EventScroll
	} else if isPress {
		eventType = value2.EventPress
	} else {
		eventType = value2.EventRelease
	}

	event := model.NewMouseEvent(eventType, button, position, modifiers)
	return &event, nil
}

// decodeButton extracts button and modifiers from SGR button code.
func (p *SGRParser) decodeButton(code int) (value2.Button, value2.Modifiers) {
	// Extract modifiers (bits 2-4)
	shift := (code & 4) != 0
	alt := (code & 8) != 0
	ctrl := (code & 16) != 0
	modifiers := value2.NewModifiers(shift, ctrl, alt)

	// Extract base button (remove modifier bits)
	// Mask: keep bits 0,1,5,6 (0x63 = 0b01100011)
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
	case 32, 35: // Motion events (bit 5 set)
		button = value2.ButtonNone
	default:
		button = value2.ButtonNone
	}

	return button, modifiers
}

// IsMotion checks if the button code represents a motion event.
func (p *SGRParser) IsMotion(buttonCode int) bool {
	baseButton := buttonCode & 0x63 // same mask as decodeButton
	return baseButton == 32 || baseButton == 35
}

// FormatSequence formats a mouse event as an SGR sequence (for testing).
func (p *SGRParser) FormatSequence(event model.MouseEvent, isPress bool) string {
	// Encode button
	buttonCode := p.encodeButton(event.Button(), event.Modifiers())

	// Get position (convert to 1-based)
	x := event.Position().X() + 1
	y := event.Position().Y() + 1

	// Format sequence
	suffix := "M"
	if !isPress {
		suffix = "m"
	}

	return fmt.Sprintf("\x1b[<%d;%d;%d%s", buttonCode, x, y, suffix)
}

// encodeButton encodes button and modifiers into SGR button code.
func (p *SGRParser) encodeButton(button value2.Button, modifiers value2.Modifiers) int {
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
		code = 32 // Motion
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
