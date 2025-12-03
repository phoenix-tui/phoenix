package model

import (
	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// DefaultTheme returns the default Phoenix theme.
// This is a balanced, modern theme with blue accents suitable for most applications.
//nolint:dupl // Theme presets have similar structure but different values - duplication is acceptable.
func DefaultTheme() *Theme {
	colors := ColorPalette{
		Primary:    value.RGB(59, 130, 246),   // Blue-500
		Secondary:  value.RGB(139, 92, 246),   // Purple-500
		Background: value.RGB(0, 0, 0),        // Black
		Surface:    value.RGB(30, 30, 30),     // Dark gray
		Text:       value.RGB(255, 255, 255),  // White
		TextMuted:  value.RGB(156, 163, 175),  // Gray-400
		Error:      value.RGB(239, 68, 68),    // Red-500
		Warning:    value.RGB(245, 158, 11),   // Amber-500
		Success:    value.RGB(34, 197, 94),    // Green-500
		Info:       value.RGB(59, 130, 246),   // Blue-500
		Border:     value.RGB(75, 85, 99),     // Gray-600
		Focus:      value.RGB(96, 165, 250),   // Blue-400
		Disabled:   value.RGB(107, 114, 128),  // Gray-500
	}

	borders := BorderStyles{
		Default: value.RoundedBorder,
		Input:   value.RoundedBorder,
		Modal:   value.ThickBorder,
		Table:   value.NormalBorder,
		Panel:   value.RoundedBorder,
	}

	spacing := SpacingScale{
		XS: 2,
		SM: 4,
		MD: 8,
		LG: 12,
		XL: 16,
	}

	typography := Typography{
		PlaceholderColor: value.RGB(107, 114, 128), // Gray-500
		CodeColor:        value.RGB(236, 72, 153),  // Pink-500
		LinkColor:        value.RGB(96, 165, 250),  // Blue-400
		HeadingColor:     value.RGB(255, 255, 255), // White
	}

	return NewTheme("Default", colors, borders, spacing, typography)
}

// DarkTheme returns a dark theme optimized for low-light environments.
// Uses darker backgrounds and muted colors to reduce eye strain.
//nolint:dupl // Theme presets have similar structure but different values - duplication is acceptable.
func DarkTheme() *Theme {
	colors := ColorPalette{
		Primary:    value.RGB(96, 165, 250),   // Blue-400
		Secondary:  value.RGB(167, 139, 250),  // Purple-400
		Background: value.RGB(0, 0, 0),        // True black
		Surface:    value.RGB(17, 24, 39),     // Gray-900
		Text:       value.RGB(229, 231, 235),  // Gray-200
		TextMuted:  value.RGB(107, 114, 128),  // Gray-500
		Error:      value.RGB(248, 113, 113),  // Red-400
		Warning:    value.RGB(251, 191, 36),   // Amber-400
		Success:    value.RGB(74, 222, 128),   // Green-400
		Info:       value.RGB(96, 165, 250),   // Blue-400
		Border:     value.RGB(55, 65, 81),     // Gray-700
		Focus:      value.RGB(147, 197, 253),  // Blue-300
		Disabled:   value.RGB(75, 85, 99),     // Gray-600
	}

	borders := BorderStyles{
		Default: value.RoundedBorder,
		Input:   value.RoundedBorder,
		Modal:   value.DoubleBorder,
		Table:   value.NormalBorder,
		Panel:   value.RoundedBorder,
	}

	spacing := SpacingScale{
		XS: 2,
		SM: 4,
		MD: 8,
		LG: 12,
		XL: 16,
	}

	typography := Typography{
		PlaceholderColor: value.RGB(75, 85, 99),   // Gray-600
		CodeColor:        value.RGB(249, 168, 212), // Pink-300
		LinkColor:        value.RGB(147, 197, 253), // Blue-300
		HeadingColor:     value.RGB(243, 244, 246), // Gray-100
	}

	return NewTheme("Dark", colors, borders, spacing, typography)
}

// LightTheme returns a light theme optimized for bright environments.
// Uses light backgrounds and darker text for readability in daylight.
//nolint:dupl // Theme presets have similar structure but different values - duplication is acceptable.
func LightTheme() *Theme {
	colors := ColorPalette{
		Primary:    value.RGB(37, 99, 235),    // Blue-600
		Secondary:  value.RGB(124, 58, 237),   // Purple-600
		Background: value.RGB(255, 255, 255),  // White
		Surface:    value.RGB(249, 250, 251),  // Gray-50
		Text:       value.RGB(17, 24, 39),     // Gray-900
		TextMuted:  value.RGB(107, 114, 128),  // Gray-500
		Error:      value.RGB(220, 38, 38),    // Red-600
		Warning:    value.RGB(217, 119, 6),    // Amber-600
		Success:    value.RGB(22, 163, 74),    // Green-600
		Info:       value.RGB(37, 99, 235),    // Blue-600
		Border:     value.RGB(209, 213, 219),  // Gray-300
		Focus:      value.RGB(59, 130, 246),   // Blue-500
		Disabled:   value.RGB(156, 163, 175),  // Gray-400
	}

	borders := BorderStyles{
		Default: value.NormalBorder,
		Input:   value.NormalBorder,
		Modal:   value.ThickBorder,
		Table:   value.NormalBorder,
		Panel:   value.NormalBorder,
	}

	spacing := SpacingScale{
		XS: 2,
		SM: 4,
		MD: 8,
		LG: 12,
		XL: 16,
	}

	typography := Typography{
		PlaceholderColor: value.RGB(156, 163, 175), // Gray-400
		CodeColor:        value.RGB(219, 39, 119),  // Pink-600
		LinkColor:        value.RGB(37, 99, 235),   // Blue-600
		HeadingColor:     value.RGB(17, 24, 39),    // Gray-900
	}

	return NewTheme("Light", colors, borders, spacing, typography)
}

// HighContrastTheme returns a high-contrast theme for accessibility.
// Uses maximum contrast ratios (black/white) and bold borders for visibility.
//nolint:dupl // Theme presets have similar structure but different values - duplication is acceptable.
func HighContrastTheme() *Theme {
	colors := ColorPalette{
		Primary:    value.RGB(0, 0, 255),      // Pure blue
		Secondary:  value.RGB(255, 0, 255),    // Pure magenta
		Background: value.RGB(0, 0, 0),        // Pure black
		Surface:    value.RGB(0, 0, 0),        // Pure black
		Text:       value.RGB(255, 255, 255),  // Pure white
		TextMuted:  value.RGB(192, 192, 192),  // Light gray
		Error:      value.RGB(255, 0, 0),      // Pure red
		Warning:    value.RGB(255, 255, 0),    // Pure yellow
		Success:    value.RGB(0, 255, 0),      // Pure green
		Info:       value.RGB(0, 255, 255),    // Pure cyan
		Border:     value.RGB(255, 255, 255),  // Pure white
		Focus:      value.RGB(255, 255, 0),    // Pure yellow
		Disabled:   value.RGB(128, 128, 128),  // Medium gray
	}

	borders := BorderStyles{
		Default: value.DoubleBorder,
		Input:   value.DoubleBorder,
		Modal:   value.ThickBorder,
		Table:   value.DoubleBorder,
		Panel:   value.DoubleBorder,
	}

	spacing := SpacingScale{
		XS: 2,
		SM: 6,  // Slightly larger for better visibility
		MD: 10,
		LG: 14,
		XL: 20,
	}

	typography := Typography{
		PlaceholderColor: value.RGB(192, 192, 192), // Light gray
		CodeColor:        value.RGB(255, 0, 255),   // Pure magenta
		LinkColor:        value.RGB(0, 255, 255),   // Pure cyan
		HeadingColor:     value.RGB(255, 255, 255), // Pure white
	}

	return NewTheme("HighContrast", colors, borders, spacing, typography)
}

// AllPresets returns all available preset themes.
// Useful for theme selection UIs.
func AllPresets() []*Theme {
	return []*Theme{
		DefaultTheme(),
		DarkTheme(),
		LightTheme(),
		HighContrastTheme(),
	}
}

// PresetByName returns a preset theme by name, or nil if not found.
// Case-insensitive matching.
func PresetByName(name string) *Theme {
	presets := map[string]func() *Theme{
		"default":      DefaultTheme,
		"dark":         DarkTheme,
		"light":        LightTheme,
		"highcontrast": HighContrastTheme,
	}

	// Normalize to lowercase
	nameLower := ""
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			nameLower += string(r + 32)
		} else {
			nameLower += string(r)
		}
	}

	if fn, ok := presets[nameLower]; ok {
		return fn()
	}

	return nil
}
