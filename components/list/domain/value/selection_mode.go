package value

// SelectionMode defines how items can be selected in a list
type SelectionMode int

const (
	// SelectionModeSingle allows only one item to be selected at a time (radio button behavior)
	SelectionModeSingle SelectionMode = iota
	// SelectionModeMulti allows multiple items to be selected (checkbox behavior)
	SelectionModeMulti
)

// IsSingle returns true if the selection mode is single-selection
func (m SelectionMode) IsSingle() bool {
	return m == SelectionModeSingle
}

// IsMulti returns true if the selection mode is multi-selection
func (m SelectionMode) IsMulti() bool {
	return m == SelectionModeMulti
}

// String returns a string representation of the selection mode
func (m SelectionMode) String() string {
	switch m {
	case SelectionModeSingle:
		return "single"
	case SelectionModeMulti:
		return "multi"
	default:
		return "unknown"
	}
}
