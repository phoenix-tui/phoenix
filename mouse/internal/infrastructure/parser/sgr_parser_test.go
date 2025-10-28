package parser

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

func TestSGRParser_ParsePress(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name          string
		sequence      string
		isPress       bool
		wantButton    value2.Button
		wantX         int
		wantY         int
		wantModifiers value2.Modifiers
	}{
		{
			name:          "left button press",
			sequence:      "<0;10;5",
			isPress:       true,
			wantButton:    value2.ButtonLeft,
			wantX:         9, // 0-based
			wantY:         4, // 0-based
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "left button release",
			sequence:      "<0;10;5",
			isPress:       false,
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "middle button",
			sequence:      "<1;20;10",
			isPress:       true,
			wantButton:    value2.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "right button",
			sequence:      "<2;30;15",
			isPress:       true,
			wantButton:    value2.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "wheel up",
			sequence:      "<64;10;5",
			isPress:       true,
			wantButton:    value2.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "wheel down",
			sequence:      "<65;10;5",
			isPress:       true,
			wantButton:    value2.ButtonWheelDown,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "with shift",
			sequence:      "<4;10;5",
			isPress:       true,
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift,
		},
		{
			name:          "with ctrl",
			sequence:      "<16;10;5",
			isPress:       true,
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierCtrl,
		},
		{
			name:          "with alt",
			sequence:      "<8;10;5",
			isPress:       true,
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierAlt,
		},
		{
			name:          "shift+ctrl",
			sequence:      "<20;10;5",
			isPress:       true,
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift | value2.ModifierCtrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := parser.Parse(tt.sequence, tt.isPress)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if event.Button() != tt.wantButton {
				t.Errorf("Button = %v, want %v", event.Button(), tt.wantButton)
			}

			if event.Position().X() != tt.wantX {
				t.Errorf("X = %d, want %d", event.Position().X(), tt.wantX)
			}

			if event.Position().Y() != tt.wantY {
				t.Errorf("Y = %d, want %d", event.Position().Y(), tt.wantY)
			}

			if event.Modifiers() != tt.wantModifiers {
				t.Errorf("Modifiers = %v, want %v", event.Modifiers(), tt.wantModifiers)
			}

			// Check event type
			expectedType := value2.EventPress
			if !tt.isPress {
				expectedType = value2.EventRelease
			}
			if event.Button().IsWheel() {
				expectedType = value2.EventScroll
			}

			if event.Type() != expectedType {
				t.Errorf("Type = %v, want %v", event.Type(), expectedType)
			}
		})
	}
}

func TestSGRParser_ParseErrors(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name     string
		sequence string
		isPress  bool
	}{
		{"too few parts", "<0;10", true},
		{"too many parts", "<0;10;5;20", true},
		{"invalid button", "<abc;10;5", true},
		{"invalid X", "<0;abc;5", true},
		{"invalid Y", "<0;10;abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.sequence, tt.isPress)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestSGRParser_FormatSequence(t *testing.T) {
	parser := NewSGRParser()

	// Create an event
	event := model.NewMouseEvent(
		value2.EventPress,
		value2.ButtonLeft,
		value2.NewPosition(9, 4),
		value2.ModifierNone,
	)

	// Format it
	sequence := parser.FormatSequence(event, true)

	// Parse it back
	parsed, err := parser.Parse(sequence[3:len(sequence)-1], true) // Remove \x1b[< and M
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsed.Button() != event.Button() {
		t.Errorf("Button mismatch after round-trip: got %v, want %v", parsed.Button(), event.Button())
	}

	if !parsed.Position().Equals(event.Position()) {
		t.Errorf("Position mismatch after round-trip: got %v, want %v", parsed.Position(), event.Position())
	}
}

func TestSGRParser_IsMotion(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name       string
		buttonCode int
		wantMotion bool
	}{
		{"motion base", 32, true},
		{"motion with button", 35, true},
		{"motion with shift", 36, true},
		{"motion with ctrl", 48, true},
		{"motion with alt", 40, true},
		{"left button", 0, false},
		{"middle button", 1, false},
		{"right button", 2, false},
		{"wheel up", 64, false},
		{"wheel down", 65, false},
		{"left with shift", 4, false},
		{"left with ctrl", 16, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.IsMotion(tt.buttonCode)
			if got != tt.wantMotion {
				t.Errorf("IsMotion(%d) = %v, want %v", tt.buttonCode, got, tt.wantMotion)
			}
		})
	}
}

func TestSGRParser_EncodeButton(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name      string
		button    value2.Button
		modifiers value2.Modifiers
		wantCode  int
	}{
		{"left button", value2.ButtonLeft, value2.ModifierNone, 0},
		{"middle button", value2.ButtonMiddle, value2.ModifierNone, 1},
		{"right button", value2.ButtonRight, value2.ModifierNone, 2},
		{"wheel up", value2.ButtonWheelUp, value2.ModifierNone, 64},
		{"wheel down", value2.ButtonWheelDown, value2.ModifierNone, 65},
		{"motion (none)", value2.ButtonNone, value2.ModifierNone, 32},
		{"left + shift", value2.ButtonLeft, value2.ModifierShift, 4},
		{"left + alt", value2.ButtonLeft, value2.ModifierAlt, 8},
		{"left + ctrl", value2.ButtonLeft, value2.ModifierCtrl, 16},
		{"left + shift + ctrl", value2.ButtonLeft, value2.ModifierShift | value2.ModifierCtrl, 20},
		{"left + shift + alt", value2.ButtonLeft, value2.ModifierShift | value2.ModifierAlt, 12},
		{"left + ctrl + alt", value2.ButtonLeft, value2.ModifierCtrl | value2.ModifierAlt, 24},
		{"left + all mods", value2.ButtonLeft, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt, 28},
		{"middle + shift", value2.ButtonMiddle, value2.ModifierShift, 5},
		{"right + ctrl", value2.ButtonRight, value2.ModifierCtrl, 18},
		{"wheel up + shift", value2.ButtonWheelUp, value2.ModifierShift, 68},
		{"motion + ctrl", value2.ButtonNone, value2.ModifierCtrl, 48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.encodeButton(tt.button, tt.modifiers)
			if got != tt.wantCode {
				t.Errorf("encodeButton(%v, %v) = %d, want %d", tt.button, tt.modifiers, got, tt.wantCode)
			}
		})
	}
}

func TestSGRParser_RoundTrip(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name      string
		button    value2.Button
		x         int
		y         int
		modifiers value2.Modifiers
		isPress   bool
	}{
		{"left press", value2.ButtonLeft, 0, 0, value2.ModifierNone, true},
		{"left release", value2.ButtonLeft, 10, 20, value2.ModifierNone, false},
		{"middle press", value2.ButtonMiddle, 50, 30, value2.ModifierNone, true},
		{"right press", value2.ButtonRight, 100, 40, value2.ModifierNone, true},
		{"wheel up", value2.ButtonWheelUp, 20, 10, value2.ModifierNone, true},
		{"wheel down", value2.ButtonWheelDown, 30, 15, value2.ModifierNone, true},
		{"left + shift", value2.ButtonLeft, 5, 5, value2.ModifierShift, true},
		{"middle + ctrl", value2.ButtonMiddle, 15, 25, value2.ModifierCtrl, true},
		{"right + alt", value2.ButtonRight, 25, 35, value2.ModifierAlt, true},
		{"left + all mods", value2.ButtonLeft, 40, 60, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt, true},
		{"large coordinates", value2.ButtonLeft, 500, 300, value2.ModifierNone, true}, // SGR supports large coords
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original event
			eventType := value2.EventPress
			if !tt.isPress {
				eventType = value2.EventRelease
			}
			if tt.button.IsWheel() {
				eventType = value2.EventScroll
			}

			original := model.NewMouseEvent(eventType, tt.button, value2.NewPosition(tt.x, tt.y), tt.modifiers)

			// Format to sequence
			sequence := parser.FormatSequence(original, tt.isPress)

			// Parse back (remove \x1b[< prefix and M/m suffix)
			content := sequence[3 : len(sequence)-1]
			parsed, err := parser.Parse(content, tt.isPress)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// Verify all fields
			if parsed.Button() != original.Button() {
				t.Errorf("Button = %v, want %v", parsed.Button(), original.Button())
			}
			if !parsed.Position().Equals(original.Position()) {
				t.Errorf("Position = %v, want %v", parsed.Position(), original.Position())
			}
			if parsed.Modifiers() != original.Modifiers() {
				t.Errorf("Modifiers = %v, want %v", parsed.Modifiers(), original.Modifiers())
			}
			if parsed.Type() != original.Type() {
				t.Errorf("Type = %v, want %v", parsed.Type(), original.Type())
			}
		})
	}
}

func TestSGRParser_DecodeButton_MotionEvents(t *testing.T) {
	parser := NewSGRParser()

	tests := []struct {
		name       string
		buttonCode int
		wantButton value2.Button
		wantMods   value2.Modifiers
	}{
		{"motion base", 32, value2.ButtonNone, value2.ModifierNone},
		{"motion with button", 35, value2.ButtonNone, value2.ModifierNone},
		{"motion + shift", 36, value2.ButtonNone, value2.ModifierShift},
		{"motion + ctrl", 48, value2.ButtonNone, value2.ModifierCtrl},
		{"motion + alt", 40, value2.ButtonNone, value2.ModifierAlt},
		{"motion + all mods", 60, value2.ButtonNone, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			button, mods := parser.decodeButton(tt.buttonCode)
			if button != tt.wantButton {
				t.Errorf("Button = %v, want %v", button, tt.wantButton)
			}
			if mods != tt.wantMods {
				t.Errorf("Modifiers = %v, want %v", mods, tt.wantMods)
			}
		})
	}
}
