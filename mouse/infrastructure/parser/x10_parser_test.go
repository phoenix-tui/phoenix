package parser

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

func TestX10Parser_Parse(t *testing.T) {
	parser := NewX10Parser()

	tests := []struct {
		name          string
		data          []byte
		wantButton    value.Button
		wantX         int
		wantY         int
		wantModifiers value.Modifiers
	}{
		{
			name:          "left button at 10,5",
			data:          []byte{32, 42, 37}, // button=0, x=10, y=5 (each +32, +1 for 1-based)
			wantButton:    value.ButtonLeft,
			wantX:         9, // 0-based
			wantY:         4, // 0-based
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "middle button",
			data:          []byte{33, 52, 42}, // button=1, x=20, y=10
			wantButton:    value.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "right button",
			data:          []byte{34, 62, 47}, // button=2, x=30, y=15
			wantButton:    value.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "wheel up",
			data:          []byte{96, 42, 37}, // button=64, x=10, y=5
			wantButton:    value.ButtonWheelUp,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "wheel down",
			data:          []byte{97, 42, 37}, // button=65, x=10, y=5
			wantButton:    value.ButtonWheelDown,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "left + shift",
			data:          []byte{36, 42, 37}, // button=4 (0+4), x=10, y=5
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift,
		},
		{
			name:          "left + alt",
			data:          []byte{40, 42, 37}, // button=8 (0+8), x=10, y=5
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierAlt,
		},
		{
			name:          "left + ctrl",
			data:          []byte{48, 42, 37}, // button=16 (0+16), x=10, y=5
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierCtrl,
		},
		{
			name:          "left + shift + ctrl",
			data:          []byte{52, 42, 37}, // button=20 (0+4+16), x=10, y=5
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift | value.ModifierCtrl,
		},
		{
			name:          "left + all modifiers",
			data:          []byte{60, 42, 37}, // button=28 (0+4+8+16), x=10, y=5
			wantButton:    value.ButtonLeft,
			wantX:         9,
			wantY:         4,
			wantModifiers: value.ModifierShift | value.ModifierCtrl | value.ModifierAlt,
		},
		{
			name:          "middle + ctrl",
			data:          []byte{49, 52, 42}, // button=17 (1+16), x=20, y=10
			wantButton:    value.ButtonMiddle,
			wantX:         19,
			wantY:         9,
			wantModifiers: value.ModifierCtrl,
		},
		{
			name:          "right + shift",
			data:          []byte{38, 62, 47}, // button=6 (2+4), x=30, y=15
			wantButton:    value.ButtonRight,
			wantX:         29,
			wantY:         14,
			wantModifiers: value.ModifierShift,
		},
		{
			name:          "at origin (0,0)",
			data:          []byte{32, 33, 33}, // button=0, x=1, y=1 (becomes 0,0 in 0-based)
			wantButton:    value.ButtonLeft,
			wantX:         0,
			wantY:         0,
			wantModifiers: value.ModifierNone,
		},
		{
			name:          "maximum coordinate (223,223)",
			data:          []byte{32, 255, 255}, // button=0, x=223, y=223 (32+223+1=256, but byte wraps)
			wantButton:    value.ButtonLeft,
			wantX:         222, // 255-32-1 = 222
			wantY:         222,
			wantModifiers: value.ModifierNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := parser.Parse(tt.data)
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

func TestX10Parser_ParseErrors(t *testing.T) {
	parser := NewX10Parser()

	tests := []struct {
		name string
		data []byte
	}{
		{"empty data", []byte{}},
		{"too short (1 byte)", []byte{32}},
		{"too short (2 bytes)", []byte{32, 42}},
		{"too long (4 bytes)", []byte{32, 42, 37, 50}},
		{"too long (5 bytes)", []byte{32, 42, 37, 50, 60}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.data)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestX10Parser_DecodeButton(t *testing.T) {
	parser := NewX10Parser()

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
		{"left + shift", 4, value.ButtonLeft, value.ModifierShift},
		{"left + alt", 8, value.ButtonLeft, value.ModifierAlt},
		{"left + ctrl", 16, value.ButtonLeft, value.ModifierCtrl},
		{"middle + shift", 5, value.ButtonMiddle, value.ModifierShift},
		{"right + ctrl", 18, value.ButtonRight, value.ModifierCtrl},
		{"wheel up + shift", 68, value.ButtonWheelUp, value.ModifierShift},
		{"all modifiers", 28, value.ButtonLeft, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt},
		{"unknown button code", 99, value.ButtonNone, value.ModifierNone}, // 99 & 0x63 = 99 (unknown)
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

func TestX10Parser_EncodeButton(t *testing.T) {
	parser := NewX10Parser()

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
		{"all modifiers", value.ButtonLeft, value.ModifierShift | value.ModifierCtrl | value.ModifierAlt, 28},
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

func TestX10Parser_FormatSequence(t *testing.T) {
	parser := NewX10Parser()

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
		{"large coords (within limit)", value.ButtonLeft, 100, 50, value.ModifierNone},
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

			// Verify format: \x1b[M + 3 bytes
			if len(sequence) != 6 {
				t.Errorf("Sequence length = %d, want 6", len(sequence))
			}

			if sequence[:3] != "\x1b[M" {
				t.Errorf("Sequence prefix = %q, want %q", sequence[:3], "\x1b[M")
			}

			// Parse back
			data := []byte(sequence[3:])
			parsed, err := parser.Parse(data)
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

func TestX10Parser_RoundTrip(t *testing.T) {
	parser := NewX10Parser()

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
		// Note: X10 limited to 223x223, so no large coordinate test
		{"near limit", value.ButtonLeft, 200, 200, value.ModifierNone},
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

			// Parse back (remove \x1b[M prefix)
			data := []byte(sequence[3:])
			parsed, err := parser.Parse(data)
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

func TestX10Parser_CoordinateBoundary(t *testing.T) {
	parser := NewX10Parser()

	// X10 protocol limitation: coordinates limited to 223 (255-32)
	// because we use byte(coord + 1 + 32) and byte max is 255

	tests := []struct {
		name string
		x    int
		y    int
	}{
		{"max valid", 222, 222}, // 222+1+32 = 255 (max byte value)
		{"at boundary", 200, 200},
		{"origin", 0, 0}, // 0+1+32 = 33
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
			data := []byte(sequence[3:])
			parsed, err := parser.Parse(data)

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if !parsed.Position().Equals(event.Position()) {
				t.Errorf("Position = %v, want %v", parsed.Position(), event.Position())
			}
		})
	}
}
