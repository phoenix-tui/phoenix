package parser

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

func TestURxvtParser_Parse(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name          string
		sequence      string
		wantButton    value.Button
		wantX         int
		wantY         int
		wantModifiers value.Modifiers
	}{
		{
			name:          "left button at 10,5",
			sequence:      "0;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9, // 0-based
			wantY:         4, // 0-based
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "middle button",
			sequence:      "1;20;10",
			wantButton:    value.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "right button",
			sequence:      "2;30;15",
			wantButton:    value.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "wheel up",
			sequence:      "64;10;5",
			wantButton:    value.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "wheel down",
			sequence:      "65;10;5",
			wantButton:    value.ButtonWheelDown,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "left + shift",
			sequence:      "4;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift,
		},
		{
			name:          "left + alt",
			sequence:      "8;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierAlt,
		},
		{
			name:          "left + ctrl",
			sequence:      "16;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierCtrl,
		},
		{
			name:          "left + shift + ctrl",
			sequence:      "20;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift | value.ModifierCtrl,
		},
		{
			name:          "left + shift + alt",
			sequence:      "12;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift | value.ModifierAlt,
		},
		{
			name:          "left + ctrl + alt",
			sequence:      "24;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierCtrl | value.ModifierAlt,
		},
		{
			name:          "left + all modifiers",
			sequence:      "28;10;5",
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift | value.ModifierCtrl | value.ModifierAlt,
		},
		{
			name:          "middle + ctrl",
			sequence:      "17;20;10",
			wantButton:    value.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value.ModifierCtrl,
		},
		{
			name:          "right + shift",
			sequence:      "6;30;15",
			wantButton:    value.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value.ModifierShift,
		},
		{
			name:          "wheel up + shift",
			sequence:      "68;10;5",
			wantButton:    value.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift,
		},
		{
			name:          "at origin",
			sequence:      "0;1;1",
			wantButton:    value.ButtonLeft,
			wantX:         0,
			wantY:         0,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "large coordinates",
			sequence:      "0;500;300",
			wantButton:    value.ButtonLeft,
			wantX:         499,
			wantY:         299,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "motion event",
			sequence:      "32;10;5",
			wantButton:    value.ButtonNone,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := parser.Parse(tt.sequence)
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

			// Check event type (URxvt doesn't distinguish press/release)
			expectedType := value.EventPress
			if event.Button().IsWheel() {
				expectedType = value.EventScroll
			}

			if event.Type() != expectedType {
				t.Errorf("Type = %v, want %v", event.Type(), expectedType)
			}
		})
	}
}

func TestURxvtParser_ParseErrors(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name     string
		sequence string
	}{
		{"too few parts", "0;10"},
		{"too many parts", "0;10;5;20"},
		{"invalid button", "abc;10;5"},
		{"invalid X", "0;abc;5"},
		{"invalid Y", "0;10;abc"},
		{"empty", ""},
		{"single value", "0"},
		{"missing semicolons", "0 10 5"},
		// Note: Negative values parse successfully (not ideal but acceptable)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.sequence)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestURxvtParser_DecodeButton(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name       string
		buttonCode int
		wantButton value.Button
		wantMods   value.Modifiers
	}{
		{"left", 0, value.ButtonLeft, value.ModifierNone},
		{"middle", 1, value.ButtonMiddle, value.ModifierNone},
		{"right", 2, value.ButtonRight, value.ModifierNone},
		{"wheel up", 64, value.ButtonWheelUp, value.ModifierNone},
		{"wheel down", 65, value.ButtonWheelDown, value.ModifierNone},
		{"motion", 32, value.ButtonNone, value.ModifierNone},
		{"motion alt", 35, value.ButtonNone, value.ModifierNone},
		{"left + shift", 4, value.ButtonLeft, value.ModifierShift},
		{"left + alt", 8, value.ButtonLeft, value.ModifierAlt},
		{"left + ctrl", 16, value.ButtonLeft, value.ModifierCtrl},
		{"middle + shift", 5, value.ButtonMiddle, value.ModifierShift},
		{"right + ctrl", 18, value.ButtonRight, value.ModifierCtrl},
		{"wheel up + shift", 68, value.ButtonWheelUp, value.ModifierShift},
		{"wheel down + alt", 73, value.ButtonWheelDown, value.ModifierAlt},
		{"all modifiers", 28, value.ButtonLeft, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt},
		{"motion + shift", 36, value.ButtonNone, value.ModifierShift},
		{"motion + ctrl", 48, value.ButtonNone, value.ModifierCtrl},
		{"motion + all mods", 60, value.ButtonNone, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt},
		// Unknown button codes default to ButtonNone
		{"unknown button (99)", 99, value.ButtonNone, value.ModifierNone}, // 99 & 0x63 = 99 (unknown base, no modifiers)
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

func TestURxvtParser_EncodeButton(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name      string
		button    value.Button
		modifiers value.Modifiers
		wantCode  int
	}{
		{"left", value.ButtonLeft, value.ModifierNone, 0},
		{"middle", value.ButtonMiddle, value.ModifierNone, 1},
		{"right", value.ButtonRight, value.ModifierNone, 2},
		{"wheel up", value.ButtonWheelUp, value.ModifierNone, 64},
		{"wheel down", value.ButtonWheelDown, value.ModifierNone, 65},
		{"motion", value.ButtonNone, value.ModifierNone, 32},
		{"left + shift", value.ButtonLeft, value.ModifierShift, 4},
		{"left + alt", value.ButtonLeft, value.ModifierAlt, 8},
		{"left + ctrl", value.ButtonLeft, value.ModifierCtrl, 16},
		{"middle + shift", value.ButtonMiddle, value.ModifierShift, 5},
		{"right + ctrl", value.ButtonRight, value.ModifierCtrl, 18},
		{"wheel up + shift", value.ButtonWheelUp, value.ModifierShift, 68},
		{"wheel down + alt", value.ButtonWheelDown, value.ModifierAlt, 73},
		{"all modifiers", value.ButtonLeft, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt, 28},
		{"motion + shift", value.ButtonNone, value.ModifierShift, 36},
		{"motion + ctrl", value.ButtonNone, value.ModifierCtrl, 48},
		{"motion + all mods", value.ButtonNone, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt, 60},
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

func TestURxvtParser_FormatSequence(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name      string
		button    value.Button
		x         int
		y         int
		modifiers value.Modifiers
	}{
		{"left at 10,5", value.ButtonLeft, 9, 4, value.ModifierNone},
		{"middle at 20,10", value.ButtonMiddle, 19, 9, value.ModifierNone},
		{"right at 30,15", value.ButtonRight, 29, 14, value.ModifierNone},
		{"wheel up", value.ButtonWheelUp, 9, 4, value.ModifierNone},
		{"wheel down", value.ButtonWheelDown, 9, 4, value.ModifierNone},
		{"left + shift", value.ButtonLeft, 9, 4, value.ModifierShift},
		{"middle + ctrl", value.ButtonMiddle, 19, 9, value.ModifierCtrl},
		{"right + alt", value.ButtonRight, 29, 14, value.ModifierAlt},
		{"all modifiers", value.ButtonLeft, 9, 4, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt},
		{"at origin", value.ButtonLeft, 0, 0, value.ModifierNone},
		{"large coords", value.ButtonLeft, 500, 300, value.ModifierNone},
		{"motion", value.ButtonNone, 10, 5, value.ModifierNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create event
			eventType := value.EventPress
			if tt.button.IsWheel() {
				eventType = value.EventScroll
			}
			event := model.NewMouseEvent(eventType, tt.button, value.NewPosition(tt.x, tt.y), tt.modifiers)

			// Format sequence
			sequence := parser.FormatSequence(event)

			// Verify format starts with \x1b[
			if sequence[:2] != "\x1b[" {
				t.Errorf("Sequence prefix = %q, want %q", sequence[:2], "\x1b[")
			}

			// Verify format ends with M
			if sequence[len(sequence)-1] != 'M' {
				t.Errorf("Sequence suffix = %q, want 'M'", sequence[len(sequence)-1])
			}

			// Parse back (remove \x1b[ prefix and M suffix)
			content := sequence[2 : len(sequence)-1]
			parsed, err := parser.Parse(content)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// Verify round-trip
			if parsed.Button() != event.Button() {
				t.Errorf("Button = %v, want %v", parsed.Button(), event.Button())
			}
			if !parsed.Position().Equals(event.Position()) {
				t.Errorf("Position = %v, want %v", parsed.Position(), event.Position())
			}
			if parsed.Modifiers() != event.Modifiers() {
				t.Errorf("Modifiers = %v, want %v", parsed.Modifiers(), event.Modifiers())
			}
		})
	}
}

func TestURxvtParser_RoundTrip(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name      string
		button    value.Button
		x         int
		y         int
		modifiers value.Modifiers
	}{
		{"left press", value.ButtonLeft, 0, 0, value.ModifierNone},
		{"middle press", value.ButtonMiddle, 50, 30, value.ModifierNone},
		{"right press", value.ButtonRight, 100, 40, value.ModifierNone},
		{"wheel up", value.ButtonWheelUp, 20, 10, value.ModifierNone},
		{"wheel down", value.ButtonWheelDown, 30, 15, value.ModifierNone},
		{"left + shift", value.ButtonLeft, 5, 5, value.ModifierShift},
		{"middle + ctrl", value.ButtonMiddle, 15, 25, value.ModifierCtrl},
		{"right + alt", value.ButtonRight, 25, 35, value.ModifierAlt},
		{"left + all mods", value.ButtonLeft, 40, 60, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt},
		{"large coordinates", value.ButtonLeft, 500, 300, value.ModifierNone},
		{"motion", value.ButtonNone, 20, 15, value.ModifierNone},
		{"motion + mods", value.ButtonNone, 30, 25, value.ModifierShift | value.ModifierCtrl},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original event
			eventType := value.EventPress
			if tt.button.IsWheel() {
				eventType = value.EventScroll
			}

			original := model.NewMouseEvent(eventType, tt.button, value.NewPosition(tt.x, tt.y), tt.modifiers)

			// Format to sequence
			sequence := parser.FormatSequence(original)

			// Parse back (remove \x1b[ prefix and M suffix)
			content := sequence[2 : len(sequence)-1]
			parsed, err := parser.Parse(content)
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

func TestURxvtParser_NoPressReleaseDistinction(t *testing.T) {
	parser := NewURxvtParser()

	// URxvt protocol doesn't distinguish press/release
	// Always creates EventPress (except for wheel -> EventScroll)

	event := model.NewMouseEvent(
		value.EventPress,
		value.ButtonLeft,
		value.NewPosition(10, 5),
		value.ModifierNone,
	)

	sequence := parser.FormatSequence(event)

	// Verify sequence always ends with 'M' (never 'm')
	if sequence[len(sequence)-1] != 'M' {
		t.Errorf("Sequence suffix = %q, want 'M'", sequence[len(sequence)-1])
	}

	// Parse back
	content := sequence[2 : len(sequence)-1]
	parsed, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Always creates press event (not release)
	if parsed.Type() != value.EventPress {
		t.Errorf("Type = %v, want %v", parsed.Type(), value.EventPress)
	}
}

func TestURxvtParser_LargeCoordinateSupport(t *testing.T) {
	parser := NewURxvtParser()

	// URxvt supports large coordinates (unlike X10)
	tests := []struct {
		name string
		x    int
		y    int
	}{
		{"within X10 limit", 200, 200},
		{"beyond X10 limit", 500, 400},
		{"very large", 1000, 800},
		{"maximum practical", 9999, 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := model.NewMouseEvent(
				value.EventPress,
				value.ButtonLeft,
				value.NewPosition(tt.x, tt.y),
				value.ModifierNone,
			)

			sequence := parser.FormatSequence(event)
			content := sequence[2 : len(sequence)-1]
			parsed, err := parser.Parse(content)

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if !parsed.Position().Equals(event.Position()) {
				t.Errorf("Position = %v, want %v", parsed.Position(), event.Position())
			}
		})
	}
}
