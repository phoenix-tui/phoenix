package value

import (
	"testing"

	"github.com/unilibs/uniwidth"
)

func TestNewUnicodeConfig(t *testing.T) {
	config := NewUnicodeConfig()

	if config.EastAsianAmbiguous() != uniwidth.EANarrow {
		t.Errorf("Default config should have EANarrow, got %v", config.EastAsianAmbiguous())
	}

	if config.IsEastAsianWide() {
		t.Error("Default config should not be East Asian wide")
	}
}

func TestWithEastAsianWide(t *testing.T) {
	config := NewUnicodeConfig().WithEastAsianWide()

	if config.EastAsianAmbiguous() != uniwidth.EAWide {
		t.Errorf("WithEastAsianWide should set EAWide, got %v", config.EastAsianAmbiguous())
	}

	if !config.IsEastAsianWide() {
		t.Error("IsEastAsianWide should return true after WithEastAsianWide")
	}
}

func TestWithEastAsianNarrow(t *testing.T) {
	// Start with wide, then set to narrow
	config := NewUnicodeConfig().WithEastAsianWide().WithEastAsianNarrow()

	if config.EastAsianAmbiguous() != uniwidth.EANarrow {
		t.Errorf("WithEastAsianNarrow should set EANarrow, got %v", config.EastAsianAmbiguous())
	}

	if config.IsEastAsianWide() {
		t.Error("IsEastAsianWide should return false after WithEastAsianNarrow")
	}
}

func TestUnicodeConfig_Immutability(t *testing.T) {
	config1 := NewUnicodeConfig()
	config2 := config1.WithEastAsianWide()

	// config1 should be unchanged
	if config1.IsEastAsianWide() {
		t.Error("config1 should not be modified (immutability violation)")
	}

	// config2 should be wide
	if !config2.IsEastAsianWide() {
		t.Error("config2 should be East Asian wide")
	}
}

func TestUnicodeConfig_Equal(t *testing.T) {
	config1 := NewUnicodeConfig()
	config2 := NewUnicodeConfig()
	config3 := NewUnicodeConfig().WithEastAsianWide()

	if !config1.Equal(config2) {
		t.Error("Two default configs should be equal")
	}

	if config1.Equal(config3) {
		t.Error("Narrow and wide configs should not be equal")
	}

	if config3.Equal(config1) {
		t.Error("Equality should be symmetric")
	}
}
