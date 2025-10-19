package application

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/infrastructure/parser"
	"github.com/phoenix-tui/phoenix/mouse/infrastructure/platform"
)

// MouseHandler is the main application service for mouse handling.
// It coordinates parsing, processing, and terminal mode management.
type MouseHandler struct {
	terminalMode   *platform.TerminalMode
	sgrParser      *parser.SGRParser
	x10Parser      *parser.X10Parser
	urxvtParser    *parser.URxvtParser
	eventProcessor *EventProcessor
}

// NewMouseHandler creates a new MouseHandler.
func NewMouseHandler() *MouseHandler {
	return &MouseHandler{
		terminalMode:   platform.NewTerminalMode(),
		sgrParser:      parser.NewSGRParser(),
		x10Parser:      parser.NewX10Parser(),
		urxvtParser:    parser.NewURxvtParser(),
		eventProcessor: NewEventProcessor(),
	}
}

// Enable enables mouse tracking.
func (h *MouseHandler) Enable() error {
	return h.terminalMode.EnableAll()
}

// Disable disables mouse tracking.
func (h *MouseHandler) Disable() error {
	return h.terminalMode.DisableAll()
}

// IsEnabled returns true if mouse tracking is enabled.
func (h *MouseHandler) IsEnabled() bool {
	return h.terminalMode.IsEnabled()
}

// ParseSequence parses a mouse input sequence and returns enriched events.
// Automatically detects protocol (SGR, X10, URxvt).
func (h *MouseHandler) ParseSequence(sequence string) ([]model.MouseEvent, error) {
	// Detect protocol and parse
	rawEvent, err := h.parseRawEvent(sequence)
	if err != nil {
		return nil, err
	}

	// Process event (add click detection, drag tracking, etc.)
	events := h.eventProcessor.ProcessEvent(*rawEvent)
	return events, nil
}

// parseRawEvent parses a raw mouse sequence into a basic mouse event.
func (h *MouseHandler) parseRawEvent(sequence string) (*model.MouseEvent, error) {
	// Remove ESC prefix if present
	sequence = strings.TrimPrefix(sequence, "\x1b")
	sequence = strings.TrimPrefix(sequence, "[")

	// Detect protocol
	if strings.HasPrefix(sequence, "<") {
		// SGR protocol
		return h.parseSGR(sequence)
	} else if strings.HasPrefix(sequence, "M") {
		// X10 protocol
		return h.parseX10(sequence)
	} else if strings.Contains(sequence, ";") && strings.HasSuffix(sequence, "M") {
		// URxvt protocol
		return h.parseURxvt(sequence)
	}

	return nil, fmt.Errorf("unknown mouse protocol: %s", sequence)
}

// parseSGR parses an SGR mouse sequence.
func (h *MouseHandler) parseSGR(sequence string) (*model.MouseEvent, error) {
	// Extract press/release indicator
	isPress := strings.HasSuffix(sequence, "M")

	// Remove suffix
	sequence = strings.TrimSuffix(sequence, "M")
	sequence = strings.TrimSuffix(sequence, "m")

	return h.sgrParser.Parse(sequence, isPress)
}

// parseX10 parses an X10 mouse sequence.
func (h *MouseHandler) parseX10(sequence string) (*model.MouseEvent, error) {
	// Remove 'M' prefix
	sequence = strings.TrimPrefix(sequence, "M")

	if len(sequence) != 3 {
		return nil, fmt.Errorf("invalid X10 sequence length: %d", len(sequence))
	}

	return h.x10Parser.Parse([]byte(sequence))
}

// parseURxvt parses a URxvt mouse sequence.
func (h *MouseHandler) parseURxvt(sequence string) (*model.MouseEvent, error) {
	// Remove 'M' suffix
	sequence = strings.TrimSuffix(sequence, "M")

	return h.urxvtParser.Parse(sequence)
}

// Processor returns the event processor (for advanced use cases).
func (h *MouseHandler) Processor() *EventProcessor {
	return h.eventProcessor
}

// Reset resets the handler state (useful for testing).
func (h *MouseHandler) Reset() {
	h.eventProcessor.Reset()
}
