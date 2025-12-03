// Package main demonstrates the Phoenix Theme system.
//
// This example shows how to:
// - Use preset themes (Default, Dark, Light, HighContrast)
// - Create custom themes
// - Switch themes at runtime
// - Use ThemeManager for thread-safe theme management
package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/style"
)

func main() {
	fmt.Println("=== Phoenix Theme System Demo ===")
	fmt.Println()

	// 1. Using preset themes
	fmt.Println("1. Preset Themes:")
	printTheme(style.DefaultTheme())
	printTheme(style.DarkTheme())
	printTheme(style.LightTheme())
	printTheme(style.HighContrastTheme())

	// 2. List all available themes
	fmt.Println("\n2. All Available Themes:")
	for _, theme := range style.AllThemes() {
		fmt.Printf("   - %s\n", theme.Name())
	}

	// 3. Theme lookup by name
	fmt.Println("\n3. Theme Lookup:")
	if theme := style.ThemeByName("dark"); theme != nil {
		fmt.Printf("   Found theme: %s\n", theme.Name())
	}

	// 4. Create custom theme
	fmt.Println("\n4. Custom Theme:")
	customTheme := style.NewTheme(
		"CustomBlue",
		style.NewColorPalette(
			style.RGB(30, 144, 255),  // Primary: Dodger Blue
			style.RGB(138, 43, 226),  // Secondary: Blue Violet
			style.RGB(15, 15, 35),    // Background: Dark Blue
			style.RGB(25, 25, 45),    // Surface: Slightly lighter
			style.RGB(240, 248, 255), // Text: Alice Blue
			style.RGB(176, 196, 222), // TextMuted: Light Steel Blue
			style.RGB(220, 20, 60),   // Error: Crimson
			style.RGB(255, 215, 0),   // Warning: Gold
			style.RGB(50, 205, 50),   // Success: Lime Green
			style.RGB(135, 206, 235), // Info: Sky Blue
			style.RGB(70, 130, 180),  // Border: Steel Blue
			style.RGB(0, 191, 255),   // Focus: Deep Sky Blue
			style.RGB(105, 105, 105), // Disabled: Dim Gray
		),
		style.NewBorderStyles(
			style.RoundedBorder,
			style.RoundedBorder,
			style.DoubleBorder,
			style.NormalBorder,
			style.RoundedBorder,
		),
		style.NewSpacingScale(2, 4, 8, 12, 16),
		style.NewTypography(
			style.RGB(128, 128, 128),
			style.RGB(255, 105, 180),
			style.RGB(0, 191, 255),
			style.RGB(240, 248, 255),
		),
	)
	printTheme(customTheme)

	// 5. Theme Manager (runtime switching)
	fmt.Println("\n5. Theme Manager (Runtime Switching):")
	tm := style.NewThemeManager(nil) // Starts with Default
	fmt.Printf("   Initial: %s\n", tm.Current().Name())

	// Switch themes
	tm.SetTheme(style.DarkTheme())
	fmt.Printf("   After SetTheme(Dark): %s\n", tm.Current().Name())

	// Switch by name
	tm.SetPreset("light")
	fmt.Printf("   After SetPreset(\"light\"): %s\n", tm.Current().Name())

	// Reset to default
	tm.Reset()
	fmt.Printf("   After Reset(): %s\n", tm.Current().Name())

	// 6. Theme inheritance (merge)
	fmt.Println("\n6. Theme Inheritance:")
	base := style.DefaultTheme()
	override := style.NewTheme(
		"CustomPrimary",
		style.ColorPalette{Primary: style.RGB(255, 0, 0)}, // Only override primary color
		style.BorderStyles{},                              // Keep base borders
		style.SpacingScale{},                              // Keep base spacing
		style.Typography{},                                // Keep base typography
	)
	merged := base.Merge(override)
	fmt.Printf("   Base: %s, Override Primary: Red\n", base.Name())
	fmt.Printf("   Merged name: %s\n", merged.Name())
	r, g, b := merged.Colors().Primary.RGB()
	fmt.Printf("   Merged primary: RGB(%d, %d, %d)\n", r, g, b)

	// 7. Demonstrating theme properties
	fmt.Println("\n7. Theme Properties (Default Theme):")
	theme := style.DefaultTheme()
	colors := theme.Colors()
	spacing := theme.Spacing()
	borders := theme.Borders()

	fmt.Printf("   Primary Color: %s\n", colors.Primary.Hex())
	fmt.Printf("   Background Color: %s\n", colors.Background.Hex())
	fmt.Printf("   Text Color: %s\n", colors.Text.Hex())
	fmt.Printf("   Spacing MD: %d\n", spacing.MD)
	fmt.Printf("   Default Border: %s\n", borders.Default)

	fmt.Println("\n=== Demo Complete ===")
}

func printTheme(theme *style.Theme) {
	colors := theme.Colors()
	fmt.Printf("   %s:\n", theme.Name())
	fmt.Printf("      Primary: %s\n", colors.Primary.Hex())
	fmt.Printf("      Background: %s\n", colors.Background.Hex())
	fmt.Printf("      Text: %s\n", colors.Text.Hex())
}
