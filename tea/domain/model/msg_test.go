package model

import (
	"testing"
)

// TestKeyMsg_String tests the String method for KeyMsg
func TestKeyMsg_String(t *testing.T) {
	tests := []struct {
		name string
		key  KeyMsg
		want string
	}{
		// Simple runes
		{
			name: "simple rune lowercase",
			key:  KeyMsg{Type: KeyRune, Rune: 'a'},
			want: "a",
		},
		{
			name: "simple rune uppercase",
			key:  KeyMsg{Type: KeyRune, Rune: 'A'},
			want: "A",
		},
		{
			name: "simple rune digit",
			key:  KeyMsg{Type: KeyRune, Rune: '5'},
			want: "5",
		},
		{
			name: "simple rune special char",
			key:  KeyMsg{Type: KeyRune, Rune: '@'},
			want: "@",
		},

		// Modifiers with runes
		{
			name: "ctrl+c rune",
			key:  KeyMsg{Type: KeyRune, Rune: 'c', Ctrl: true},
			want: "ctrl+c",
		},
		{
			name: "ctrl+a",
			key:  KeyMsg{Type: KeyRune, Rune: 'a', Ctrl: true},
			want: "ctrl+a",
		},
		{
			name: "alt+x",
			key:  KeyMsg{Type: KeyRune, Rune: 'x', Alt: true},
			want: "alt+x",
		},
		{
			name: "shift+a (implicit in uppercase)",
			key:  KeyMsg{Type: KeyRune, Rune: 'A', Shift: true},
			want: "A",
		},
		{
			name: "alt+ctrl+s",
			key:  KeyMsg{Type: KeyRune, Rune: 's', Alt: true, Ctrl: true},
			want: "alt+ctrl+s",
		},
		{
			name: "alt+shift+ctrl+z",
			key:  KeyMsg{Type: KeyRune, Rune: 'z', Alt: true, Ctrl: true, Shift: true},
			want: "alt+ctrl+z",
		},

		// Special keys
		{
			name: "enter",
			key:  KeyMsg{Type: KeyEnter},
			want: "enter",
		},
		{
			name: "backspace",
			key:  KeyMsg{Type: KeyBackspace},
			want: "backspace",
		},
		{
			name: "tab",
			key:  KeyMsg{Type: KeyTab},
			want: "tab",
		},
		{
			name: "escape",
			key:  KeyMsg{Type: KeyEsc},
			want: "esc",
		},
		{
			name: "space",
			key:  KeyMsg{Type: KeySpace},
			want: "space",
		},

		// Arrow keys
		{
			name: "arrow up",
			key:  KeyMsg{Type: KeyUp},
			want: "↑",
		},
		{
			name: "arrow down",
			key:  KeyMsg{Type: KeyDown},
			want: "↓",
		},
		{
			name: "arrow left",
			key:  KeyMsg{Type: KeyLeft},
			want: "←",
		},
		{
			name: "arrow right",
			key:  KeyMsg{Type: KeyRight},
			want: "→",
		},

		// Navigation keys
		{
			name: "home",
			key:  KeyMsg{Type: KeyHome},
			want: "home",
		},
		{
			name: "end",
			key:  KeyMsg{Type: KeyEnd},
			want: "end",
		},
		{
			name: "page up",
			key:  KeyMsg{Type: KeyPgUp},
			want: "pgup",
		},
		{
			name: "page down",
			key:  KeyMsg{Type: KeyPgDown},
			want: "pgdown",
		},
		{
			name: "delete",
			key:  KeyMsg{Type: KeyDelete},
			want: "delete",
		},
		{
			name: "insert",
			key:  KeyMsg{Type: KeyInsert},
			want: "insert",
		},

		// Function keys
		{
			name: "F1",
			key:  KeyMsg{Type: KeyF1},
			want: "F1",
		},
		{
			name: "F2",
			key:  KeyMsg{Type: KeyF2},
			want: "F2",
		},
		{
			name: "F3",
			key:  KeyMsg{Type: KeyF3},
			want: "F3",
		},
		{
			name: "F4",
			key:  KeyMsg{Type: KeyF4},
			want: "F4",
		},
		{
			name: "F5",
			key:  KeyMsg{Type: KeyF5},
			want: "F5",
		},
		{
			name: "F6",
			key:  KeyMsg{Type: KeyF6},
			want: "F6",
		},
		{
			name: "F7",
			key:  KeyMsg{Type: KeyF7},
			want: "F7",
		},
		{
			name: "F8",
			key:  KeyMsg{Type: KeyF8},
			want: "F8",
		},
		{
			name: "F9",
			key:  KeyMsg{Type: KeyF9},
			want: "F9",
		},
		{
			name: "F10",
			key:  KeyMsg{Type: KeyF10},
			want: "F10",
		},
		{
			name: "F11",
			key:  KeyMsg{Type: KeyF11},
			want: "F11",
		},
		{
			name: "F12",
			key:  KeyMsg{Type: KeyF12},
			want: "F12",
		},

		// Special keys with modifiers
		{
			name: "alt+enter",
			key:  KeyMsg{Type: KeyEnter, Alt: true},
			want: "alt+enter",
		},
		{
			name: "ctrl+backspace",
			key:  KeyMsg{Type: KeyBackspace, Ctrl: true},
			want: "ctrl+backspace",
		},
		{
			name: "shift+tab",
			key:  KeyMsg{Type: KeyTab, Shift: true},
			want: "shift+tab",
		},
		{
			name: "ctrl+arrow up",
			key:  KeyMsg{Type: KeyUp, Ctrl: true},
			want: "ctrl+↑",
		},
		{
			name: "shift+arrow down",
			key:  KeyMsg{Type: KeyDown, Shift: true},
			want: "shift+↓",
		},
		{
			name: "ctrl+home",
			key:  KeyMsg{Type: KeyHome, Ctrl: true},
			want: "ctrl+home",
		},
		{
			name: "ctrl+end",
			key:  KeyMsg{Type: KeyEnd, Ctrl: true},
			want: "ctrl+end",
		},
		{
			name: "ctrl+delete",
			key:  KeyMsg{Type: KeyDelete, Ctrl: true},
			want: "ctrl+delete",
		},
		{
			name: "shift+F1",
			key:  KeyMsg{Type: KeyF1, Shift: true},
			want: "shift+F1",
		},
		{
			name: "ctrl+F5",
			key:  KeyMsg{Type: KeyF5, Ctrl: true},
			want: "ctrl+F5",
		},
		{
			name: "alt+F12",
			key:  KeyMsg{Type: KeyF12, Alt: true},
			want: "alt+F12",
		},

		// Dedicated Ctrl+C
		{
			name: "ctrl+c dedicated",
			key:  KeyMsg{Type: KeyCtrlC},
			want: "ctrl+c",
		},
		{
			name: "ctrl+c dedicated with modifiers (ignored)",
			key:  KeyMsg{Type: KeyCtrlC, Alt: true},
			want: "ctrl+c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.key.String()
			if got != tt.want {
				t.Errorf("KeyMsg.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestKeyMsg_TypeAssertion tests that KeyMsg implements Msg interface
func TestKeyMsg_TypeAssertion(t *testing.T) {
	var msg Msg = KeyMsg{Type: KeyRune, Rune: 'a'}

	keyMsg, ok := msg.(KeyMsg)
	if !ok {
		t.Error("KeyMsg should implement Msg interface")
	}

	if keyMsg.Rune != 'a' {
		t.Errorf("expected rune 'a', got %v", keyMsg.Rune)
	}
}

// TestMouseMsg_String tests the String method for MouseMsg
func TestMouseMsg_String(t *testing.T) {
	tests := []struct {
		name  string
		mouse MouseMsg
		want  string
	}{
		// Left button
		{
			name:  "left press",
			mouse: MouseMsg{X: 10, Y: 5, Button: MouseButtonLeft, Action: MouseActionPress},
			want:  "left press at (10, 5)",
		},
		{
			name:  "left release",
			mouse: MouseMsg{X: 20, Y: 10, Button: MouseButtonLeft, Action: MouseActionRelease},
			want:  "left release at (20, 10)",
		},
		{
			name:  "left motion (drag)",
			mouse: MouseMsg{X: 15, Y: 8, Button: MouseButtonLeft, Action: MouseActionMotion},
			want:  "left motion at (15, 8)",
		},

		// Middle button
		{
			name:  "middle press",
			mouse: MouseMsg{X: 5, Y: 5, Button: MouseButtonMiddle, Action: MouseActionPress},
			want:  "middle press at (5, 5)",
		},
		{
			name:  "middle release",
			mouse: MouseMsg{X: 6, Y: 6, Button: MouseButtonMiddle, Action: MouseActionRelease},
			want:  "middle release at (6, 6)",
		},

		// Right button
		{
			name:  "right press",
			mouse: MouseMsg{X: 30, Y: 15, Button: MouseButtonRight, Action: MouseActionPress},
			want:  "right press at (30, 15)",
		},
		{
			name:  "right release",
			mouse: MouseMsg{X: 31, Y: 16, Button: MouseButtonRight, Action: MouseActionRelease},
			want:  "right release at (31, 16)",
		},

		// Mouse motion without button
		{
			name:  "mouse motion",
			mouse: MouseMsg{X: 15, Y: 8, Button: MouseButtonNone, Action: MouseActionMotion},
			want:  "mouse motion at (15, 8)",
		},
		{
			name:  "mouse motion at origin",
			mouse: MouseMsg{X: 0, Y: 0, Button: MouseButtonNone, Action: MouseActionMotion},
			want:  "mouse motion at (0, 0)",
		},

		// Wheel events
		{
			name:  "wheel up",
			mouse: MouseMsg{X: 0, Y: 0, Button: MouseButtonWheelUp, Action: MouseActionPress},
			want:  "wheel up at (0, 0)",
		},
		{
			name:  "wheel down",
			mouse: MouseMsg{X: 10, Y: 10, Button: MouseButtonWheelDown, Action: MouseActionPress},
			want:  "wheel down at (10, 10)",
		},
		{
			name:  "wheel up at position",
			mouse: MouseMsg{X: 25, Y: 12, Button: MouseButtonWheelUp, Action: MouseActionPress},
			want:  "wheel up at (25, 12)",
		},

		// Edge cases
		{
			name:  "origin click",
			mouse: MouseMsg{X: 0, Y: 0, Button: MouseButtonLeft, Action: MouseActionPress},
			want:  "left press at (0, 0)",
		},
		{
			name:  "large coordinates",
			mouse: MouseMsg{X: 999, Y: 999, Button: MouseButtonRight, Action: MouseActionPress},
			want:  "right press at (999, 999)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mouse.String()
			if got != tt.want {
				t.Errorf("MouseMsg.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestMouseMsg_TypeAssertion tests that MouseMsg implements Msg interface
func TestMouseMsg_TypeAssertion(t *testing.T) {
	var msg Msg = MouseMsg{X: 10, Y: 5, Button: MouseButtonLeft, Action: MouseActionPress}

	mouseMsg, ok := msg.(MouseMsg)
	if !ok {
		t.Error("MouseMsg should implement Msg interface")
	}

	if mouseMsg.X != 10 || mouseMsg.Y != 5 {
		t.Errorf("expected position (10, 5), got (%d, %d)", mouseMsg.X, mouseMsg.Y)
	}
}

// TestWindowSizeMsg_String tests the String method for WindowSizeMsg
func TestWindowSizeMsg_String(t *testing.T) {
	tests := []struct {
		name string
		msg  WindowSizeMsg
		want string
	}{
		{
			name: "standard size",
			msg:  WindowSizeMsg{Width: 80, Height: 24},
			want: "window resize: 80x24",
		},
		{
			name: "large terminal",
			msg:  WindowSizeMsg{Width: 200, Height: 60},
			want: "window resize: 200x60",
		},
		{
			name: "small terminal",
			msg:  WindowSizeMsg{Width: 40, Height: 10},
			want: "window resize: 40x10",
		},
		{
			name: "single digit",
			msg:  WindowSizeMsg{Width: 5, Height: 5},
			want: "window resize: 5x5",
		},
		{
			name: "wide terminal",
			msg:  WindowSizeMsg{Width: 300, Height: 24},
			want: "window resize: 300x24",
		},
		{
			name: "tall terminal",
			msg:  WindowSizeMsg{Width: 80, Height: 100},
			want: "window resize: 80x100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.String()
			if got != tt.want {
				t.Errorf("WindowSizeMsg.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestWindowSizeMsg_IsValid tests the IsValid method for WindowSizeMsg
func TestWindowSizeMsg_IsValid(t *testing.T) {
	tests := []struct {
		name string
		msg  WindowSizeMsg
		want bool
	}{
		{
			name: "valid standard size",
			msg:  WindowSizeMsg{Width: 80, Height: 24},
			want: true,
		},
		{
			name: "valid large size",
			msg:  WindowSizeMsg{Width: 200, Height: 60},
			want: true,
		},
		{
			name: "valid minimal size",
			msg:  WindowSizeMsg{Width: 1, Height: 1},
			want: true,
		},
		{
			name: "zero width",
			msg:  WindowSizeMsg{Width: 0, Height: 24},
			want: false,
		},
		{
			name: "zero height",
			msg:  WindowSizeMsg{Width: 80, Height: 0},
			want: false,
		},
		{
			name: "both zero",
			msg:  WindowSizeMsg{Width: 0, Height: 0},
			want: false,
		},
		{
			name: "negative width",
			msg:  WindowSizeMsg{Width: -1, Height: 24},
			want: false,
		},
		{
			name: "negative height",
			msg:  WindowSizeMsg{Width: 80, Height: -1},
			want: false,
		},
		{
			name: "both negative",
			msg:  WindowSizeMsg{Width: -10, Height: -20},
			want: false,
		},
		{
			name: "large negative width",
			msg:  WindowSizeMsg{Width: -999, Height: 24},
			want: false,
		},
		{
			name: "large negative height",
			msg:  WindowSizeMsg{Width: 80, Height: -999},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.IsValid()
			if got != tt.want {
				t.Errorf("WindowSizeMsg.IsValid() = %v, want %v (Width=%d, Height=%d)",
					got, tt.want, tt.msg.Width, tt.msg.Height)
			}
		})
	}
}

// TestWindowSizeMsg_TypeAssertion tests that WindowSizeMsg implements Msg interface
func TestWindowSizeMsg_TypeAssertion(t *testing.T) {
	var msg Msg = WindowSizeMsg{Width: 80, Height: 24}

	sizeMsg, ok := msg.(WindowSizeMsg)
	if !ok {
		t.Error("WindowSizeMsg should implement Msg interface")
	}

	if sizeMsg.Width != 80 || sizeMsg.Height != 24 {
		t.Errorf("expected size 80x24, got %dx%d", sizeMsg.Width, sizeMsg.Height)
	}
}

// TestQuitMsg_String tests the String method for QuitMsg
func TestQuitMsg_String(t *testing.T) {
	msg := QuitMsg{}
	want := "quit"

	got := msg.String()
	if got != want {
		t.Errorf("QuitMsg.String() = %q, want %q", got, want)
	}
}

// TestQuitMsg_TypeAssertion tests that QuitMsg implements Msg interface
func TestQuitMsg_TypeAssertion(t *testing.T) {
	var msg Msg = QuitMsg{}

	quitMsg, ok := msg.(QuitMsg)
	if !ok {
		t.Error("QuitMsg should implement Msg interface")
	}

	// QuitMsg is an empty struct, so just verify it was cast successfully
	_ = quitMsg
}

// TestMsg_AnyType tests that any type can be used as Msg
func TestMsg_AnyType(t *testing.T) {
	// Custom message type
	type CustomMsg struct {
		Value string
	}

	var msg Msg = CustomMsg{Value: "test"}

	customMsg, ok := msg.(CustomMsg)
	if !ok {
		t.Error("Custom type should implement Msg interface")
	}

	if customMsg.Value != "test" {
		t.Errorf("expected value 'test', got %q", customMsg.Value)
	}
}

// TestKeyType_Coverage tests all KeyType constants are defined
func TestKeyType_Coverage(t *testing.T) {
	// This test ensures all KeyType constants are defined and can be used
	keyTypes := []KeyType{
		KeyRune,
		KeyEnter,
		KeyBackspace,
		KeyTab,
		KeyEsc,
		KeySpace,
		KeyUp,
		KeyDown,
		KeyLeft,
		KeyRight,
		KeyHome,
		KeyEnd,
		KeyPgUp,
		KeyPgDown,
		KeyDelete,
		KeyInsert,
		KeyF1,
		KeyF2,
		KeyF3,
		KeyF4,
		KeyF5,
		KeyF6,
		KeyF7,
		KeyF8,
		KeyF9,
		KeyF10,
		KeyF11,
		KeyF12,
		KeyCtrlC,
	}

	if len(keyTypes) != 29 {
		t.Errorf("expected 29 key types, got %d", len(keyTypes))
	}

	// Verify each can be used in a KeyMsg
	for i, kt := range keyTypes {
		msg := KeyMsg{Type: kt}
		_ = msg.String() // Should not panic

		// Verify they have sequential values (iota)
		if kt != KeyType(i) {
			t.Errorf("KeyType[%d] = %d, expected sequential iota values", i, kt)
		}
	}
}

// TestMouseButton_Coverage tests all MouseButton constants are defined
func TestMouseButton_Coverage(t *testing.T) {
	buttons := []MouseButton{
		MouseButtonNone,
		MouseButtonLeft,
		MouseButtonMiddle,
		MouseButtonRight,
		MouseButtonWheelUp,
		MouseButtonWheelDown,
	}

	if len(buttons) != 6 {
		t.Errorf("expected 6 mouse buttons, got %d", len(buttons))
	}

	// Verify sequential iota values
	for i, btn := range buttons {
		if btn != MouseButton(i) {
			t.Errorf("MouseButton[%d] = %d, expected sequential iota values", i, btn)
		}
	}
}

// TestMouseAction_Coverage tests all MouseAction constants are defined
func TestMouseAction_Coverage(t *testing.T) {
	actions := []MouseAction{
		MouseActionPress,
		MouseActionRelease,
		MouseActionMotion,
	}

	if len(actions) != 3 {
		t.Errorf("expected 3 mouse actions, got %d", len(actions))
	}

	// Verify sequential iota values
	for i, action := range actions {
		if action != MouseAction(i) {
			t.Errorf("MouseAction[%d] = %d, expected sequential iota values", i, action)
		}
	}
}

// TestKeyTypeName_AllTypes tests the keyTypeName helper function
func TestKeyTypeName_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		keyType  KeyType
		rune     rune
		expected string
	}{
		{"KeyRune", KeyRune, 'a', "a"},
		{"KeyEnter", KeyEnter, 0, "enter"},
		{"KeyBackspace", KeyBackspace, 0, "backspace"},
		{"KeyTab", KeyTab, 0, "tab"},
		{"KeyEsc", KeyEsc, 0, "esc"},
		{"KeySpace", KeySpace, 0, "space"},
		{"KeyUp", KeyUp, 0, "↑"},
		{"KeyDown", KeyDown, 0, "↓"},
		{"KeyLeft", KeyLeft, 0, "←"},
		{"KeyRight", KeyRight, 0, "→"},
		{"KeyHome", KeyHome, 0, "home"},
		{"KeyEnd", KeyEnd, 0, "end"},
		{"KeyPgUp", KeyPgUp, 0, "pgup"},
		{"KeyPgDown", KeyPgDown, 0, "pgdown"},
		{"KeyDelete", KeyDelete, 0, "delete"},
		{"KeyInsert", KeyInsert, 0, "insert"},
		{"KeyF1", KeyF1, 0, "F1"},
		{"KeyF2", KeyF2, 0, "F2"},
		{"KeyF3", KeyF3, 0, "F3"},
		{"KeyF4", KeyF4, 0, "F4"},
		{"KeyF5", KeyF5, 0, "F5"},
		{"KeyF6", KeyF6, 0, "F6"},
		{"KeyF7", KeyF7, 0, "F7"},
		{"KeyF8", KeyF8, 0, "F8"},
		{"KeyF9", KeyF9, 0, "F9"},
		{"KeyF10", KeyF10, 0, "F10"},
		{"KeyF11", KeyF11, 0, "F11"},
		{"KeyF12", KeyF12, 0, "F12"},
		{"KeyCtrlC", KeyCtrlC, 0, "ctrl+c"},
		{"Unknown", KeyType(999), 0, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyTypeName(tt.keyType, tt.rune)
			if got != tt.expected {
				t.Errorf("keyTypeName(%v, %v) = %q, want %q", tt.keyType, tt.rune, got, tt.expected)
			}
		})
	}
}

// TestMouseButtonName_AllButtons tests the mouseButtonName helper function
func TestMouseButtonName_AllButtons(t *testing.T) {
	tests := []struct {
		name     string
		button   MouseButton
		expected string
	}{
		{"None", MouseButtonNone, "mouse"},
		{"Left", MouseButtonLeft, "left"},
		{"Middle", MouseButtonMiddle, "middle"},
		{"Right", MouseButtonRight, "right"},
		{"WheelUp", MouseButtonWheelUp, "wheel up"},
		{"WheelDown", MouseButtonWheelDown, "wheel down"},
		{"Unknown", MouseButton(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mouseButtonName(tt.button)
			if got != tt.expected {
				t.Errorf("mouseButtonName(%v) = %q, want %q", tt.button, got, tt.expected)
			}
		})
	}
}

// TestMouseActionName_AllActions tests the mouseActionName helper function
func TestMouseActionName_AllActions(t *testing.T) {
	tests := []struct {
		name     string
		action   MouseAction
		expected string
	}{
		{"Press", MouseActionPress, "press"},
		{"Release", MouseActionRelease, "release"},
		{"Motion", MouseActionMotion, "motion"},
		{"Unknown", MouseAction(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mouseActionName(tt.action)
			if got != tt.expected {
				t.Errorf("mouseActionName(%v) = %q, want %q", tt.action, got, tt.expected)
			}
		})
	}
}
