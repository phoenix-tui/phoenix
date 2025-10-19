package value

import "github.com/unilibs/uniwidth"

// UnicodeConfig represents configuration for Unicode text width calculation.
// This is a value object - immutable and compared by value.
//
// Used to configure locale-specific behavior, particularly East Asian Ambiguous
// character width (e.g., ±, ½, °, ×).
type UnicodeConfig struct {
	eastAsianAmbiguous uniwidth.EAWidth
}

// NewUnicodeConfig creates a new Unicode configuration with default settings.
// Default: East Asian Ambiguous characters are NARROW (width 1).
//
// Example:
//
//	config := NewUnicodeConfig()
//	// ± will be width 1 (neutral locale)
func NewUnicodeConfig() UnicodeConfig {
	return UnicodeConfig{
		eastAsianAmbiguous: uniwidth.EANarrow,
	}
}

// WithEastAsianWide returns a new config with East Asian Ambiguous characters as WIDE (width 2).
// Use this for East Asian locales (Japanese, Chinese, Korean).
//
// Affected characters: ±, ½, °, ×, §, etc.
//
// Example:
//
//	config := NewUnicodeConfig().WithEastAsianWide()
//	// ± will be width 2 (East Asian locale)
func (c UnicodeConfig) WithEastAsianWide() UnicodeConfig {
	return UnicodeConfig{
		eastAsianAmbiguous: uniwidth.EAWide,
	}
}

// WithEastAsianNarrow returns a new config with East Asian Ambiguous characters as NARROW (width 1).
// This is the default for neutral locales (English, etc.).
//
// Example:
//
//	config := NewUnicodeConfig().WithEastAsianNarrow()
//	// ± will be width 1 (neutral locale)
func (c UnicodeConfig) WithEastAsianNarrow() UnicodeConfig {
	return UnicodeConfig{
		eastAsianAmbiguous: uniwidth.EANarrow,
	}
}

// EastAsianAmbiguous returns the current East Asian Ambiguous width setting.
// For internal use by UnicodeService.
func (c UnicodeConfig) EastAsianAmbiguous() uniwidth.EAWidth {
	return c.eastAsianAmbiguous
}

// IsEastAsianWide returns true if East Asian Ambiguous characters are configured as wide.
//
// Example:
//
//	config := NewUnicodeConfig().WithEastAsianWide()
//	if config.IsEastAsianWide() {
//	    // Use East Asian locale rendering
//	}
func (c UnicodeConfig) IsEastAsianWide() bool {
	return c.eastAsianAmbiguous == uniwidth.EAWide
}

// Equal returns true if two configs are equal.
// Value objects must implement equality for comparison.
func (c UnicodeConfig) Equal(other UnicodeConfig) bool {
	return c.eastAsianAmbiguous == other.eastAsianAmbiguous
}
