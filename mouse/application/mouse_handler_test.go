package application

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

func TestNewMouseHandler(t *testing.T) {
	handler := NewMouseHandler()

	if handler == nil {
		t.Fatal("Expected handler, got nil")
	}

	if handler.terminalMode == nil {
		t.Error("Expected terminal mode to be initialized")
	}

	if handler.sgrParser == nil {
		t.Error("Expected SGR parser to be initialized")
	}

	if handler.x10Parser == nil {
		t.Error("Expected X10 parser to be initialized")
	}

	if handler.urxvtParser == nil {
		t.Error("Expected URxvt parser to be initialized")
	}

	if handler.eventProcessor == nil {
		t.Error("Expected event processor to be initialized")
	}

	// Initially disabled
	if handler.IsEnabled() {
		t.Error("Expected mouse tracking disabled initially")
	}
}

func TestEnable(t *testing.T) {
	handler := NewMouseHandler()

	// Initially disabled
	if handler.IsEnabled() {
		t.Error("Expected disabled initially")
	}

	// Enable
	err := handler.Enable()
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	// Should be enabled
	if !handler.IsEnabled() {
		t.Error("Expected enabled after Enable()")
	}
}

func TestDisable(t *testing.T) {
	handler := NewMouseHandler()

	// Enable first
	err := handler.Enable()
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	if !handler.IsEnabled() {
		t.Error("Expected enabled")
	}

	// Disable
	err = handler.Disable()
	if err != nil {
		t.Fatalf("Disable failed: %v", err)
	}

	// Should be disabled
	if handler.IsEnabled() {
		t.Error("Expected disabled after Disable()")
	}
}

func TestIsEnabled(t *testing.T) {
	handler := NewMouseHandler()

	// Test initial state
	if handler.IsEnabled() {
		t.Error("Expected disabled initially")
	}

	// Enable
	handler.Enable()
	if !handler.IsEnabled() {
		t.Error("Expected enabled after Enable()")
	}

	// Disable
	handler.Disable()
	if handler.IsEnabled() {
		t.Error("Expected disabled after Disable()")
	}
}

func TestParseSequence_SGR_Press(t *testing.T) {
	handler := NewMouseHandler()

	// SGR format: ESC[<0;10;5M (left button press at 10,5)
	sequence := "\x1b[<0;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	// First event should be press
	event := events[0]
	if event.Type() != value.EventPress {
		t.Errorf("Expected EventPress, got %v", event.Type())
	}

	if event.Button() != value.ButtonLeft {
		t.Errorf("Expected ButtonLeft, got %v", event.Button())
	}

	// Terminal coordinates are 1-based, Position is 0-based
	if event.Position().X() != 9 || event.Position().Y() != 4 {
		t.Errorf("Expected position (9,4), got (%d,%d)",
			event.Position().X(), event.Position().Y())
	}
}

func TestParseSequence_SGR_Release(t *testing.T) {
	handler := NewMouseHandler()

	// SGR format: ESC[<0;10;5m (left button release at 10,5)
	sequence := "\x1b[<0;10;5m"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	// Should contain release event (and possibly click)
	hasRelease := false
	for _, event := range events {
		if event.Type() == value.EventRelease {
			hasRelease = true
			break
		}
	}

	if !hasRelease {
		t.Error("Expected EventRelease in results")
	}
}

func TestParseSequence_SGR_WithModifiers(t *testing.T) {
	handler := NewMouseHandler()

	// SGR with Shift modifier: ESC[<4;10;5M (0 + Shift(4))
	sequence := "\x1b[<4;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	event := events[0]
	if !event.Modifiers().HasShift() {
		t.Error("Expected Shift modifier")
	}
}

func TestParseSequence_SGR_ButtonNone(t *testing.T) {
	handler := NewMouseHandler()

	// SGR with button code 32 (ButtonNone/motion indicator)
	// Parser sets button to ButtonNone
	sequence := "\x1b[<32;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	event := events[0]
	if event.Button() != value.ButtonNone {
		t.Errorf("Expected ButtonNone, got %v", event.Button())
	}
}

func TestParseSequence_SGR_Scroll(t *testing.T) {
	handler := NewMouseHandler()

	tests := []struct {
		name     string
		sequence string
		button   value.Button
	}{
		{
			name:     "Scroll up",
			sequence: "\x1b[<64;10;5M",
			button:   value.ButtonWheelUp,
		},
		{
			name:     "Scroll down",
			sequence: "\x1b[<65;10;5M",
			button:   value.ButtonWheelDown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := handler.ParseSequence(tt.sequence)
			if err != nil {
				t.Fatalf("ParseSequence failed: %v", err)
			}

			if len(events) == 0 {
				t.Fatal("Expected at least 1 event")
			}

			event := events[0]
			if event.Type() != value.EventScroll {
				t.Errorf("Expected EventScroll, got %v", event.Type())
			}

			if event.Button() != tt.button {
				t.Errorf("Expected %v, got %v", tt.button, event.Button())
			}
		})
	}
}

func TestParseSequence_X10(t *testing.T) {
	handler := NewMouseHandler()

	// X10 format: ESC[Mbxy (left button press)
	// b=32+button (32=press), x=33+X, y=33+Y
	// For position (10,5): x=43, y=38
	sequence := "\x1b[M\x20\x2b\x26" // 32, 43, 38

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	event := events[0]
	if event.Button() != value.ButtonLeft {
		t.Errorf("Expected ButtonLeft, got %v", event.Button())
	}
}

func TestParseSequence_X10_InvalidLength(t *testing.T) {
	handler := NewMouseHandler()

	// X10 requires exactly 3 bytes after 'M'
	sequence := "\x1b[M\x20\x2b" // Only 2 bytes

	_, err := handler.ParseSequence(sequence)
	if err == nil {
		t.Error("Expected error for invalid X10 sequence length")
	}

	if !strings.Contains(err.Error(), "invalid X10 sequence length") {
		t.Errorf("Expected 'invalid X10 sequence length' error, got: %v", err)
	}
}

func TestParseSequence_URxvt(t *testing.T) {
	handler := NewMouseHandler()

	// URxvt format: ESC[32;10;5M (button;x;y;M)
	sequence := "\x1b[32;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}

	event := events[0]
	// Terminal coordinates are 1-based, Position is 0-based
	if event.Position().X() != 9 || event.Position().Y() != 4 {
		t.Errorf("Expected position (9,4), got (%d,%d)",
			event.Position().X(), event.Position().Y())
	}
}

func TestParseSequence_UnknownProtocol(t *testing.T) {
	handler := NewMouseHandler()

	// Invalid sequence
	sequence := "\x1b[X123"

	_, err := handler.ParseSequence(sequence)
	if err == nil {
		t.Error("Expected error for unknown protocol")
	}

	if !strings.Contains(err.Error(), "unknown mouse protocol") {
		t.Errorf("Expected 'unknown mouse protocol' error, got: %v", err)
	}
}

func TestParseSequence_WithoutESCPrefix(t *testing.T) {
	handler := NewMouseHandler()

	// Sequence without ESC prefix (should still work)
	sequence := "[<0;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least 1 event")
	}
}

func TestParseSequence_ProtocolDetection(t *testing.T) {
	handler := NewMouseHandler()

	tests := []struct {
		name     string
		sequence string
		protocol string
	}{
		{
			name:     "SGR (starts with <)",
			sequence: "\x1b[<0;10;5M",
			protocol: "SGR",
		},
		{
			name:     "X10 (starts with M)",
			sequence: "\x1b[M\x20\x2b\x26",
			protocol: "X10",
		},
		{
			name:     "URxvt (contains ; and ends with M)",
			sequence: "\x1b[32;10;5M",
			protocol: "URxvt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := handler.ParseSequence(tt.sequence)
			if err != nil {
				t.Fatalf("%s: ParseSequence failed: %v", tt.protocol, err)
			}

			if len(events) == 0 {
				t.Fatalf("%s: Expected at least 1 event", tt.protocol)
			}

			// Just verify parsing succeeded (protocol was detected)
		})
	}
}

func TestParseSequence_EnrichmentWithClickDetection(t *testing.T) {
	handler := NewMouseHandler()

	// First release (should generate click)
	sequence1 := "\x1b[<0;10;5m"
	events1, err := handler.ParseSequence(sequence1)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	// Should contain click event
	hasClick := false
	for _, e := range events1 {
		if e.Type() == value.EventClick {
			hasClick = true
			break
		}
	}

	if !hasClick {
		t.Error("Expected click event after release")
	}

	// Second release at same position (should generate double click)
	sequence2 := "\x1b[<0;10;5m"
	events2, err := handler.ParseSequence(sequence2)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	// Should contain double click event
	hasDoubleClick := false
	for _, e := range events2 {
		if e.Type() == value.EventDoubleClick {
			hasDoubleClick = true
			break
		}
	}

	if !hasDoubleClick {
		t.Error("Expected double click event after second release")
	}
}

func TestParseSequence_EnrichmentWithPressRelease(t *testing.T) {
	handler := NewMouseHandler()

	// Press at (10,10) terminal coords = (9,9) position
	sequence1 := "\x1b[<0;10;10M"
	events1, err := handler.ParseSequence(sequence1)
	if err != nil {
		t.Fatalf("ParseSequence press failed: %v", err)
	}

	if len(events1) != 1 || events1[0].Type() != value.EventPress {
		t.Error("Expected press event")
	}

	// Release at same position
	sequence2 := "\x1b[<0;10;10m"
	events2, err := handler.ParseSequence(sequence2)
	if err != nil {
		t.Fatalf("ParseSequence release failed: %v", err)
	}

	// Should contain release event
	hasRelease := false
	for _, e := range events2 {
		if e.Type() == value.EventRelease {
			hasRelease = true
			break
		}
	}

	if !hasRelease {
		t.Error("Expected release event")
	}
}

func TestProcessor(t *testing.T) {
	handler := NewMouseHandler()

	processor := handler.Processor()
	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	// Verify it's the same processor
	if processor != handler.eventProcessor {
		t.Error("Expected same processor instance")
	}

	// Verify processor functionality
	if processor.ClickCount() != 0 {
		t.Error("Expected initial click count 0")
	}
}

func TestMouseHandler_Reset(t *testing.T) {
	handler := NewMouseHandler()

	// Create some state via click
	sequence1 := "\x1b[<0;10;10m" // Release (generates click)
	handler.ParseSequence(sequence1)

	// Verify state exists
	if handler.Processor().ClickCount() == 0 {
		t.Error("Expected click count > 0 before reset")
	}

	// Reset
	handler.Reset()

	// Verify state cleared
	if handler.Processor().IsDragging() {
		t.Error("Expected not dragging after reset")
	}

	if handler.Processor().ClickCount() != 0 {
		t.Errorf("Expected click count 0 after reset, got %d",
			handler.Processor().ClickCount())
	}
}

// Integration test: Full click sequence
func TestIntegration_MouseHandler_ClickSequence(t *testing.T) {
	handler := NewMouseHandler()

	pos := "10;5"

	// Triple click sequence
	clicks := []string{
		"\x1b[<0;" + pos + "m", // Release 1
		"\x1b[<0;" + pos + "m", // Release 2
		"\x1b[<0;" + pos + "m", // Release 3
	}

	for i, seq := range clicks {
		events, err := handler.ParseSequence(seq)
		if err != nil {
			t.Fatalf("Click %d failed: %v", i+1, err)
		}

		// Each should generate events
		if len(events) == 0 {
			t.Fatalf("Click %d: expected events", i+1)
		}
	}

	// Verify triple click was detected
	if handler.Processor().ClickCount() != 3 {
		t.Errorf("Expected click count 3, got %d",
			handler.Processor().ClickCount())
	}
}

// Integration test: Press and release sequence
func TestIntegration_MouseHandler_PressReleaseSequence(t *testing.T) {
	handler := NewMouseHandler()

	// Press at terminal (10,10) = position (9,9)
	press := "\x1b[<0;10;10M"
	events, err := handler.ParseSequence(press)
	if err != nil {
		t.Fatalf("Press failed: %v", err)
	}

	if len(events) != 1 || events[0].Type() != value.EventPress {
		t.Error("Expected press event")
	}

	// Release at same position (generates click)
	release := "\x1b[<0;10;10m"
	events, err = handler.ParseSequence(release)
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	// Should contain both release and click
	hasRelease := false
	hasClick := false
	for _, e := range events {
		if e.Type() == value.EventRelease {
			hasRelease = true
		}
		if e.IsClick() {
			hasClick = true
		}
	}

	if !hasRelease {
		t.Error("Expected release event")
	}

	if !hasClick {
		t.Error("Expected click event after press+release")
	}

	// Click count should be 1
	if handler.Processor().ClickCount() != 1 {
		t.Errorf("Expected click count 1, got %d", handler.Processor().ClickCount())
	}
}

// Integration test: Scroll events
func TestIntegration_ScrollSequence(t *testing.T) {
	handler := NewMouseHandler()

	tests := []struct {
		name     string
		sequence string
		isUp     bool
		isDown   bool
		delta    int
	}{
		{
			name:     "Scroll up",
			sequence: "\x1b[<64;10;5M",
			isUp:     true,
			isDown:   false,
			delta:    -3,
		},
		{
			name:     "Scroll down",
			sequence: "\x1b[<65;10;5M",
			isUp:     false,
			isDown:   true,
			delta:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := handler.ParseSequence(tt.sequence)
			if err != nil {
				t.Fatalf("ParseSequence failed: %v", err)
			}

			if len(events) == 0 {
				t.Fatal("Expected at least 1 event")
			}

			event := events[0]

			// Verify scroll helpers
			processor := handler.Processor()
			if processor.IsScrollUp(event) != tt.isUp {
				t.Errorf("IsScrollUp: expected %v, got %v",
					tt.isUp, processor.IsScrollUp(event))
			}

			if processor.IsScrollDown(event) != tt.isDown {
				t.Errorf("IsScrollDown: expected %v, got %v",
					tt.isDown, processor.IsScrollDown(event))
			}

			delta := processor.ScrollDelta(event)
			if delta != tt.delta {
				t.Errorf("ScrollDelta: expected %d, got %d", tt.delta, delta)
			}
		})
	}
}

// Integration test: Mixed button types
func TestIntegration_MultipleButtons(t *testing.T) {
	handler := NewMouseHandler()

	pos := "10;5"

	buttons := []struct {
		name   string
		code   string
		button value.Button
	}{
		{"Left", "0", value.ButtonLeft},
		{"Middle", "1", value.ButtonMiddle},
		{"Right", "2", value.ButtonRight},
	}

	for _, btn := range buttons {
		t.Run(btn.name, func(t *testing.T) {
			// Press
			press := "\x1b[<" + btn.code + ";" + pos + "M"
			events, err := handler.ParseSequence(press)
			if err != nil {
				t.Fatalf("Press %s failed: %v", btn.name, err)
			}

			if len(events) == 0 {
				t.Fatal("Expected press event")
			}

			if events[0].Button() != btn.button {
				t.Errorf("Expected %v, got %v", btn.button, events[0].Button())
			}

			// Release
			release := "\x1b[<" + btn.code + ";" + pos + "m"
			events, err = handler.ParseSequence(release)
			if err != nil {
				t.Fatalf("Release %s failed: %v", btn.name, err)
			}

			// Should generate click
			hasClick := false
			for _, e := range events {
				if e.IsClick() {
					hasClick = true
					break
				}
			}

			if !hasClick {
				t.Errorf("Expected click for %s button", btn.name)
			}
		})
	}
}

// Edge case: Enable/disable multiple times
func TestEnableDisableMultipleTimes(t *testing.T) {
	handler := NewMouseHandler()

	for i := 0; i < 5; i++ {
		// Enable
		err := handler.Enable()
		if err != nil {
			t.Fatalf("Enable %d failed: %v", i, err)
		}

		if !handler.IsEnabled() {
			t.Errorf("Enable %d: expected enabled", i)
		}

		// Disable
		err = handler.Disable()
		if err != nil {
			t.Fatalf("Disable %d failed: %v", i, err)
		}

		if handler.IsEnabled() {
			t.Errorf("Disable %d: expected disabled", i)
		}
	}
}

// Edge case: Empty sequence
func TestParseSequence_Empty(t *testing.T) {
	handler := NewMouseHandler()

	sequence := ""

	_, err := handler.ParseSequence(sequence)
	if err == nil {
		t.Error("Expected error for empty sequence")
	}
}

// Edge case: Parse after reset
func TestParseSequence_AfterReset(t *testing.T) {
	handler := NewMouseHandler()

	// Create state
	sequence1 := "\x1b[<0;10;5M"
	handler.ParseSequence(sequence1)

	// Reset
	handler.Reset()

	// Parse again (should work)
	sequence2 := "\x1b[<0;10;5m"
	events, err := handler.ParseSequence(sequence2)
	if err != nil {
		t.Fatalf("ParseSequence after reset failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected events after reset")
	}

	// Should generate fresh click (count = 1)
	if handler.Processor().ClickCount() != 1 {
		t.Errorf("Expected click count 1 after reset, got %d",
			handler.Processor().ClickCount())
	}
}

// Edge case: Very long sequence
func TestParseSequence_LongSequence(t *testing.T) {
	handler := NewMouseHandler()

	// SGR with very large coordinates
	sequence := "\x1b[<0;9999;9999M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected events")
	}

	event := events[0]
	// Terminal coordinates are 1-based, Position is 0-based
	if event.Position().X() != 9998 || event.Position().Y() != 9998 {
		t.Errorf("Expected position (9998,9998), got (%d,%d)",
			event.Position().X(), event.Position().Y())
	}
}

// Edge case: Modifiers combination
func TestParseSequence_MultipleModifiers(t *testing.T) {
	handler := NewMouseHandler()

	// Button 0 + Shift(4) + Alt(8) + Ctrl(16) = 28
	sequence := "\x1b[<28;10;5M"

	events, err := handler.ParseSequence(sequence)
	if err != nil {
		t.Fatalf("ParseSequence failed: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected events")
	}

	event := events[0]
	mods := event.Modifiers()

	if !mods.HasShift() {
		t.Error("Expected Shift modifier")
	}

	if !mods.HasAlt() {
		t.Error("Expected Alt modifier")
	}

	if !mods.HasCtrl() {
		t.Error("Expected Ctrl modifier")
	}
}
