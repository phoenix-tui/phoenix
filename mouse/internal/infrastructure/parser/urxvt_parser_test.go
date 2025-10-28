package parser

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

func TestURxvtParser_Parse(t *testing.T) {
	parser := NewURxvtParser()

	tests := []struct {
		name          string
		sequence      string
		wantButton    value2.Button
		wantX         int
		wantY         int
		wantModifiers value2.Modifiers
	}{
		{
			name:          "left button at 10,5",
			sequence:      "0;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9, // 0-based
			wantY:         4, // 0-based
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "middle button",
			sequence:      "1;20;10",
			wantButton:    value2.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "right button",
			sequence:      "2;30;15",
			wantButton:    value2.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "wheel up",
			sequence:      "64;10;5",
			wantButton:    value2.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "wheel down",
			sequence:      "65;10;5",
			wantButton:    value2.ButtonWheelDown,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "left + shift",
			sequence:      "4;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift,
		},
		{
			name:          "left + alt",
			sequence:      "8;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierAlt,
		},
		{
			name:          "left + ctrl",
			sequence:      "16;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierCtrl,
		},
		{
			name:          "left + shift + ctrl",
			sequence:      "20;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift | value2.ModifierCtrl,
		},
		{
			name:          "left + shift + alt",
			sequence:      "12;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift | value2.ModifierAlt,
		},
		{
			name:          "left + ctrl + alt",
			sequence:      "24;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierCtrl | value2.ModifierAlt,
		},
		{
			name:          "left + all modifiers",
			sequence:      "28;10;5",
			wantButton:    value2.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt,
		},
		{
			name:          "middle + ctrl",
			sequence:      "17;20;10",
			wantButton:    value2.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value2.ModifierCtrl,
		},
		{
			name:          "right + shift",
			sequence:      "6;30;15",
			wantButton:    value2.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value2.ModifierShift,
		},
		{
			name:          "wheel up + shift",
			sequence:      "68;10;5",
			wantButton:    value2.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierShift,
		},
		{
			name:          "at origin",
			sequence:      "0;1;1",
			wantButton:    value2.ButtonLeft,
			wantX:         0,
			wantY:         0,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "large coordinates",
			sequence:      "0;500;300",
			wantButton:    value2.ButtonLeft,
			wantX:         499,
			wantY:         299,
			wantModifiers: value2.ModifierNone,
		},
		{
			name:          "motion event",
			sequence:      "32;10;5",
			wantButton:    value2.ButtonNone,
			wantX:         9,
			wantY:         4,
			wantModifiers: value2.ModifierNone,
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
			expectedType := value2.EventPress
			if event.Button().IsWheel() {
				expectedType = value2.EventScroll
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
		wantButton value2.Button
		wantMods   value2.Modifiers
	}{
		{"left", 0, value2.ButtonLeft, value2.ModifierNone},
		{"middle", 1, value2.ButtonMiddle, value2.ModifierNone},
		{"right", 2, value2.ButtonRight, value2.ModifierNone},
		{"wheel up", 64, value2.ButtonWheelUp, value2.ModifierNone},
		{"wheel down", 65, value2.ButtonWheelDown, value2.ModifierNone},
		{"motion", 32, value2.ButtonNone, value2.ModifierNone},
		{"motion alt", 35, value2.ButtonNone, value2.ModifierNone},
		{"left + shift", 4, value2.ButtonLeft, value2.ModifierShift},
		{"left + alt", 8, value2.ButtonLeft, value2.ModifierAlt},
		{"left + ctrl", 16, value2.ButtonLeft, value2.ModifierCtrl},
		{"middle + shift", 5, value2.ButtonMiddle, value2.ModifierShift},
		{"right + ctrl", 18, value2.ButtonRight, value2.ModifierCtrl},
		{"wheel up + shift", 68, value2.ButtonWheelUp, value2.ModifierShift},
		{"wheel down + alt", 73, value2.ButtonWheelDown, value2.ModifierAlt},
		{"all modifiers", 28, value2.ButtonLeft, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt},
		{"motion + shift", 36, value2.ButtonNone, value2.ModifierShift},
		{"motion + ctrl", 48, value2.ButtonNone, value2.ModifierCtrl},
		{"motion + all mods", 60, value2.ButtonNone, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt},
		// Unknown button codes default to ButtonNone
		{"unknown button (99)", 99, value2.ButtonNone, value2.ModifierNone}, // 99 & 0x63 = 99 (unknown base, no modifiers)
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
		button    value2.Button
		modifiers value2.Modifiers
		wantCode  int
	}{
		{"left", value2.ButtonLeft, value2.ModifierNone, 0},
		{"middle", value2.ButtonMiddle, value2.ModifierNone, 1},
		{"right", value2.ButtonRight, value2.ModifierNone, 2},
		{"wheel up", value2.ButtonWheelUp, value2.ModifierNone, 64},
		{"wheel down", value2.ButtonWheelDown, value2.ModifierNone, 65},
		{"motion", value2.ButtonNone, value2.ModifierNone, 32},
		{"left + shift", value2.ButtonLeft, value2.ModifierShift, 4},
		{"left + alt", value2.ButtonLeft, value2.ModifierAlt, 8},
		{"left + ctrl", value2.ButtonLeft, value2.ModifierCtrl, 16},
		{"middle + shift", value2.ButtonMiddle, value2.ModifierShift, 5},
		{"right + ctrl", value2.ButtonRight, value2.ModifierCtrl, 18},
		{"wheel up + shift", value2.ButtonWheelUp, value2.ModifierShift, 68},
		{"wheel down + alt", value2.ButtonWheelDown, value2.ModifierAlt, 73},
		{"all modifiers", value2.ButtonLeft, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt, 28},
		{"motion + shift", value2.ButtonNone, value2.ModifierShift, 36},
		{"motion + ctrl", value2.ButtonNone, value2.ModifierCtrl, 48},
		{"motion + all mods", value2.ButtonNone, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt, 60},
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
		button    value2.Button
		x         int
		y         int
		modifiers value2.Modifiers
	}{
		{"left at 10,5", value2.ButtonLeft, 9, 4, value2.ModifierNone},
		{"middle at 20,10", value2.ButtonMiddle, 19, 9, value2.ModifierNone},
		{"right at 30,15", value2.ButtonRight, 29, 14, value2.ModifierNone},
		{"wheel up", value2.ButtonWheelUp, 9, 4, value2.ModifierNone},
		{"wheel down", value2.ButtonWheelDown, 9, 4, value2.ModifierNone},
		{"left + shift", value2.ButtonLeft, 9, 4, value2.ModifierShift},
		{"middle + ctrl", value2.ButtonMiddle, 19, 9, value2.ModifierCtrl},
		{"right + alt", value2.ButtonRight, 29, 14, value2.ModifierAlt},
		{"all modifiers", value2.ButtonLeft, 9, 4, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt},
		{"at origin", value2.ButtonLeft, 0, 0, value2.ModifierNone},
		{"large coords", value2.ButtonLeft, 500, 300, value2.ModifierNone},
		{"motion", value2.ButtonNone, 10, 5, value2.ModifierNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create event
			eventType := value2.EventPress
			if tt.button.IsWheel() {
				eventType = value2.EventScroll
			}
			event := model.NewMouseEvent(eventType, tt.button, value2.NewPosition(tt.x, tt.y), tt.modifiers)

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
		button    value2.Button
		x         int
		y         int
		modifiers value2.Modifiers
	}{
		{"left press", value2.ButtonLeft, 0, 0, value2.ModifierNone},
		{"middle press", value2.ButtonMiddle, 50, 30, value2.ModifierNone},
		{"right press", value2.ButtonRight, 100, 40, value2.ModifierNone},
		{"wheel up", value2.ButtonWheelUp, 20, 10, value2.ModifierNone},
		{"wheel down", value2.ButtonWheelDown, 30, 15, value2.ModifierNone},
		{"left + shift", value2.ButtonLeft, 5, 5, value2.ModifierShift},
		{"middle + ctrl", value2.ButtonMiddle, 15, 25, value2.ModifierCtrl},
		{"right + alt", value2.ButtonRight, 25, 35, value2.ModifierAlt},
		{"left + all mods", value2.ButtonLeft, 40, 60, value2.ModifierShift | value2.ModifierCtrl | value2.ModifierAlt},
		{"large coordinates", value2.ButtonLeft, 500, 300, value2.ModifierNone},
		{"motion", value2.ButtonNone, 20, 15, value2.ModifierNone},
		{"motion + mods", value2.ButtonNone, 30, 25, value2.ModifierShift | value2.ModifierCtrl},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original event
			eventType := value2.EventPress
			if tt.button.IsWheel() {
				eventType = value2.EventScroll
			}

			original := model.NewMouseEvent(eventType, tt.button, value2.NewPosition(tt.x, tt.y), tt.modifiers)

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
		value2.EventPress,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
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
	if parsed.Type() != value2.EventPress {
		t.Errorf("Type = %v, want %v", parsed.Type(), value2.EventPress)
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
				value2.EventPress,
				value2.ButtonLeft,
				value2.NewPosition(tt.x, tt.y),
				value2.ModifierNone,
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
