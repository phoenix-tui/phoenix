package value

import "testing"

func TestNewModifiers(t *testing.T) {
	tests := []struct {
		name  string
		shift bool
		ctrl  bool
		alt   bool
	}{
		{"none", false, false, false},
		{"shift only", true, false, false},
		{"ctrl only", false, true, false},
		{"alt only", false, false, true},
		{"shift+ctrl", true, true, false},
		{"all", true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModifiers(tt.shift, tt.ctrl, tt.alt)

			if m.HasShift() != tt.shift {
				t.Errorf("HasShift() = %v, want %v", m.HasShift(), tt.shift)
			}
			if m.HasCtrl() != tt.ctrl {
				t.Errorf("HasCtrl() = %v, want %v", m.HasCtrl(), tt.ctrl)
			}
			if m.HasAlt() != tt.alt {
				t.Errorf("HasAlt() = %v, want %v", m.HasAlt(), tt.alt)
			}
		})
	}
}

func TestModifiers_String(t *testing.T) {
	tests := []struct {
		modifiers Modifiers
		expected  string
	}{
		{ModifierNone, "None"},
		{ModifierShift, "Shift"},
		{ModifierCtrl, "Ctrl"},
		{ModifierAlt, "Alt"},
		{ModifierShift | ModifierCtrl, "Shift+Ctrl"},
		{ModifierShift | ModifierAlt, "Shift+Alt"},
		{ModifierCtrl | ModifierAlt, "Ctrl+Alt"},
		{ModifierShift | ModifierCtrl | ModifierAlt, "Shift+Ctrl+Alt"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.modifiers.String(); got != tt.expected {
				t.Errorf("Modifiers.String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestModifiers_Equals(t *testing.T) {
	m1 := NewModifiers(true, false, true)
	m2 := NewModifiers(true, false, true)
	m3 := NewModifiers(true, true, true)

	if !m1.Equals(m2) {
		t.Errorf("Equals() = false, want true for equal modifiers")
	}
	if m1.Equals(m3) {
		t.Errorf("Equals() = true, want false for different modifiers")
	}
}
