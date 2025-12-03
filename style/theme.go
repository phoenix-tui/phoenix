package style

import (
	"github.com/phoenix-tui/phoenix/style/internal/application"
	"github.com/phoenix-tui/phoenix/style/internal/domain/model"
)

// Theme represents a complete UI theme with colors, borders, spacing, and typography.
//
// Zero value: Not useful - use NewTheme() or preset functions.
//
//	var t style.Theme          // Zero value - NOT useful
//	t2 := style.DefaultTheme() // Correct - use preset
//	t3 := style.NewTheme(...)  // Correct - custom theme
//
// Thread safety: Theme is immutable and safe for concurrent use.
// All setter methods return new instances.
type Theme = model.Theme

// ColorPalette defines all colors used in a theme.
// Organized by semantic purpose rather than specific colors.
type ColorPalette = model.ColorPalette

// BorderStyles defines border styles for different UI elements.
type BorderStyles = model.BorderStyles

// SpacingScale defines the spacing values used throughout a theme.
// Follows a consistent scale for visual harmony.
type SpacingScale = model.SpacingScale

// Typography defines text-related theme settings.
type Typography = model.Typography

// ThemeManager manages the current theme and provides thread-safe theme switching.
//
// Zero value: Not valid - use NewThemeManager().
//
//	var tm style.ThemeManager          // Zero value - INVALID
//	tm2 := style.NewThemeManager(nil)  // Correct - uses DefaultTheme
//
// Thread safety: ThemeManager is safe for concurrent use.
type ThemeManager = application.ThemeManager

// NewTheme creates a new custom Theme.
// For common themes, use preset functions (DefaultTheme, DarkTheme, etc.).
//
// Example:
//
//	colors := style.ColorPalette{
//	    Primary:    style.RGB(59, 130, 246),
//	    Background: style.RGB(0, 0, 0),
//	    Text:       style.RGB(255, 255, 255),
//	    // ... other colors
//	}
//	theme := style.NewTheme("MyTheme", colors, borders, spacing, typography)
func NewTheme(
	name string,
	colors ColorPalette,
	borders BorderStyles,
	spacing SpacingScale,
	typography Typography,
) *Theme {
	return model.NewTheme(name, colors, borders, spacing, typography)
}

// DefaultTheme returns the default Phoenix theme.
// This is a balanced, modern theme with blue accents suitable for most applications.
//
// Example:
//
//	theme := style.DefaultTheme()
//	fmt.Println(theme.Name()) // "Default"
func DefaultTheme() *Theme {
	return model.DefaultTheme()
}

// DarkTheme returns a dark theme optimized for low-light environments.
// Uses darker backgrounds and muted colors to reduce eye strain.
//
// Example:
//
//	theme := style.DarkTheme()
//	fmt.Println(theme.Name()) // "Dark"
func DarkTheme() *Theme {
	return model.DarkTheme()
}

// LightTheme returns a light theme optimized for bright environments.
// Uses light backgrounds and darker text for readability in daylight.
//
// Example:
//
//	theme := style.LightTheme()
//	fmt.Println(theme.Name()) // "Light"
func LightTheme() *Theme {
	return model.LightTheme()
}

// HighContrastTheme returns a high-contrast theme for accessibility.
// Uses maximum contrast ratios and bold borders for visibility.
//
// Example:
//
//	theme := style.HighContrastTheme()
//	fmt.Println(theme.Name()) // "HighContrast"
func HighContrastTheme() *Theme {
	return model.HighContrastTheme()
}

// AllThemes returns all available preset themes.
// Useful for theme selection UIs.
//
// Example:
//
//	themes := style.AllThemes()
//	for _, theme := range themes {
//	    fmt.Println(theme.Name())
//	}
func AllThemes() []*Theme {
	return model.AllPresets()
}

// ThemeByName returns a preset theme by name, or nil if not found.
// Case-insensitive matching.
//
// Supported names: "Default", "Dark", "Light", "HighContrast"
//
// Example:
//
//	theme := style.ThemeByName("dark")
//	if theme != nil {
//	    fmt.Println("Found theme:", theme.Name())
//	}
func ThemeByName(name string) *Theme {
	return model.PresetByName(name)
}

// NewThemeManager creates a new ThemeManager with the given initial theme.
// If theme is nil, DefaultTheme is used.
//
// ThemeManager is safe for concurrent use from multiple goroutines.
//
// Example:
//
//	// Use default theme
//	tm := style.NewThemeManager(nil)
//
//	// Use dark theme
//	tm := style.NewThemeManager(style.DarkTheme())
//
//	// Later, switch themes
//	tm.SetTheme(style.LightTheme())
func NewThemeManager(theme *Theme) *ThemeManager {
	return application.NewThemeManager(theme)
}

// Helper functions for creating theme components

// NewColorPalette creates a ColorPalette with all colors specified.
// This is a convenience function for custom theme creation.
func NewColorPalette(
	primary, secondary, background, surface, text, textMuted Color,
	errorColor, warning, success, info, border, focus, disabled Color,
) ColorPalette {
	return ColorPalette{
		Primary:    primary,
		Secondary:  secondary,
		Background: background,
		Surface:    surface,
		Text:       text,
		TextMuted:  textMuted,
		Error:      errorColor,
		Warning:    warning,
		Success:    success,
		Info:       info,
		Border:     border,
		Focus:      focus,
		Disabled:   disabled,
	}
}

// NewBorderStyles creates BorderStyles with all borders specified.
// This is a convenience function for custom theme creation.
func NewBorderStyles(
	defaultBorder, input, modal, table, panel Border,
) BorderStyles {
	return BorderStyles{
		Default: defaultBorder,
		Input:   input,
		Modal:   modal,
		Table:   table,
		Panel:   panel,
	}
}

// NewSpacingScale creates a SpacingScale with all values specified.
// This is a convenience function for custom theme creation.
func NewSpacingScale(xs, sm, md, lg, xl int) SpacingScale {
	return SpacingScale{
		XS: xs,
		SM: sm,
		MD: md,
		LG: lg,
		XL: xl,
	}
}

// NewTypography creates Typography with all settings specified.
// This is a convenience function for custom theme creation.
func NewTypography(placeholder, code, link, heading Color) Typography {
	return Typography{
		PlaceholderColor: placeholder,
		CodeColor:        code,
		LinkColor:        link,
		HeadingColor:     heading,
	}
}
