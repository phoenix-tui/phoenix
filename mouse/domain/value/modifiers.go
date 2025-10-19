package value

// Modifiers represents keyboard modifiers held during a mouse event.
type Modifiers int

const (
	// ModifierNone represents no modifiers.
	ModifierNone Modifiers = 0
	// ModifierShift represents the Shift key.
	ModifierShift Modifiers = 1 << iota
	// ModifierCtrl represents the Ctrl key.
	ModifierCtrl
	// ModifierAlt represents the Alt key.
	ModifierAlt
)

// NewModifiers creates a new Modifiers value.
func NewModifiers(shift, ctrl, alt bool) Modifiers {
	m := ModifierNone
	if shift {
		m |= ModifierShift
	}
	if ctrl {
		m |= ModifierCtrl
	}
	if alt {
		m |= ModifierAlt
	}
	return m
}

// HasShift returns true if Shift is held.
func (m Modifiers) HasShift() bool {
	return m&ModifierShift != 0
}

// HasCtrl returns true if Ctrl is held.
func (m Modifiers) HasCtrl() bool {
	return m&ModifierCtrl != 0
}

// HasAlt returns true if Alt is held.
func (m Modifiers) HasAlt() bool {
	return m&ModifierAlt != 0
}

// String returns the string representation of the modifiers.
func (m Modifiers) String() string {
	if m == ModifierNone {
		return "None"
	}

	var mods []string
	if m.HasShift() {
		mods = append(mods, "Shift")
	}
	if m.HasCtrl() {
		mods = append(mods, "Ctrl")
	}
	if m.HasAlt() {
		mods = append(mods, "Alt")
	}

	result := ""
	for i, mod := range mods {
		if i > 0 {
			result += "+"
		}
		result += mod
	}
	return result
}

// Equals checks if two modifiers are equal.
func (m Modifiers) Equals(other Modifiers) bool {
	return m == other
}
