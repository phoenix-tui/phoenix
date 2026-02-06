package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/style"
)

func main() {
	fmt.Println("=== Phoenix Style Library - Complete Examples ===\n")

	// Example 1: Styled Header.
	fmt.Println("1. Application Header:")
	headerStyle := style.New().
		Foreground(style.White).
		Background(style.RGB(0, 102, 204)). // Custom blue
		Bold(true).
		Padding(style.NewPadding(1, 3, 1, 3)).
		Width(60).
		Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle))

	fmt.Println(style.Render(headerStyle, "🚀 Phoenix TUI Framework"))
	fmt.Println()

	// Example 2: Info Box.
	fmt.Println("2. Info Box:")
	infoStyle := style.New().
		Border(style.RoundedBorder).
		BorderColor(style.Cyan).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 0, 1, 0)).
		Foreground(style.Cyan)

	infoMessage := `ℹ️  Information

Your task has been completed successfully.
All tests passed with 100% coverage.`

	fmt.Println(style.Render(infoStyle, infoMessage))
	fmt.Println()

	// Example 3: Error Box.
	fmt.Println("3. Error Alert:")
	errorStyle := style.New().
		Border(style.ThickBorder).
		BorderColor(style.Red).
		Foreground(style.Red).
		Background(style.RGB(50, 0, 0)). // Dark red background
		Bold(true).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(0, 0, 1, 0))

	errorMessage := `❌ Error: Operation Failed

Could not connect to database.
Please check your configuration.`

	fmt.Println(style.Render(errorStyle, errorMessage))
	fmt.Println()

	// Example 4: Success Box.
	fmt.Println("4. Success Message:")
	successStyle := style.New().
		Border(style.DoubleBorder).
		BorderColor(style.Green).
		Foreground(style.Green).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Width(50).
		Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle))

	fmt.Println(style.Render(successStyle, "✅ Build Successful"))
	fmt.Println()

	// Example 5: Menu Items.
	fmt.Println("5. Styled Menu:")
	menuItemStyle := style.New().
		Foreground(style.White).
		Background(style.RGB(40, 44, 52)). // Dark gray
		Padding(style.NewPadding(0, 2, 0, 2)).
		Margin(style.NewMargin(0, 0, 0, 0))

	selectedStyle := style.New().
		Foreground(style.Black).
		Background(style.Cyan).
		Bold(true).
		Padding(style.NewPadding(0, 2, 0, 2))

	fmt.Println(style.Render(selectedStyle, "▸ Home"))
	fmt.Println(style.Render(menuItemStyle, "  Settings"))
	fmt.Println(style.Render(menuItemStyle, "  About"))
	fmt.Println(style.Render(menuItemStyle, "  Exit"))
	fmt.Println()

	// Example 6: Table-like layout.
	fmt.Println("6. Data Display:")
	headerCellStyle := style.New().
		Background(style.RGB(60, 60, 60)).
		Foreground(style.White).
		Bold(true).
		Padding(style.NewPadding(0, 1, 0, 1)).
		Width(20).
		Align(style.NewAlignment(style.AlignLeft, style.AlignTop))

	cellStyle := style.New().
		Foreground(style.White).
		Padding(style.NewPadding(0, 1, 0, 1)).
		Width(20).
		Align(style.NewAlignment(style.AlignLeft, style.AlignTop))

	fmt.Print(style.Render(headerCellStyle, "Name"))
	fmt.Print(style.Render(headerCellStyle, "Status"))
	fmt.Println(style.Render(headerCellStyle, "Progress"))

	fmt.Print(style.Render(cellStyle, "Task 1"))
	fmt.Print(style.Render(cellStyle, "✅ Complete"))
	fmt.Println(style.Render(cellStyle, "100%"))

	fmt.Print(style.Render(cellStyle, "Task 2"))
	fmt.Print(style.Render(cellStyle, "⏳ Running"))
	fmt.Println(style.Render(cellStyle, "75%"))
	fmt.Println()

	// Example 7: Card Layout.
	fmt.Println("7. Card Component:")
	cardStyle := style.New().
		Border(style.RoundedBorder).
		BorderColor(style.RGB(128, 128, 128)).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 0, 1, 0)).
		Width(50)

	cardTitle := style.New().
		Foreground(style.Cyan).
		Bold(true)

	cardContent := `Phoenix TUI Framework

A modern, high-performance TUI library for Go
with DDD architecture and perfect Unicode support.

Features:
• 10x faster than competitors
• Rich domain models
• Fluent API design
• Perfect emoji rendering 👋🌍`

	fmt.Println(style.Render(cardTitle, "Project Info"))
	fmt.Println(style.Render(cardStyle, cardContent))
	fmt.Println()

	// Example 8: Progress Indicator.
	fmt.Println("8. Progress Bar:")
	progressBarStyle := style.New().
		Border(style.ASCIIBorder).
		Padding(style.NewPadding(0, 1, 0, 1)).
		Width(50)

	progressFilled := style.New().
		Background(style.Green).
		Foreground(style.Green)

	progressEmpty := style.New().
		Background(style.RGB(40, 40, 40)).
		Foreground(style.RGB(40, 40, 40))

	progress := style.Render(progressFilled, "████████████████████") +
		style.Render(progressEmpty, "            ")

	fmt.Println(style.Render(progressBarStyle, progress))
	fmt.Println(style.Render(style.New().Foreground(style.Cyan), "Progress: 62% complete"))
	fmt.Println()

	// Example 9: Notification.
	fmt.Println("9. System Notification:")
	notificationStyle := style.New().
		Border(style.RoundedBorder).
		BorderColor(style.Yellow).
		Foreground(style.Yellow).
		Background(style.RGB(40, 40, 0)). // Dark yellow
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 0, 1, 0)).
		Bold(true).
		Width(55)

	fmt.Println(style.Render(notificationStyle, "⚠️  Warning: Low disk space (15% remaining)"))
	fmt.Println()

	// Example 10: Complete Dashboard.
	fmt.Println("10. Dashboard Layout:")
	fmt.Println(style.Render(headerStyle, "System Dashboard"))

	statusOk := style.New().
		Foreground(style.Green).
		Bold(true)

	statusWarning := style.New().
		Foreground(style.Yellow).
		Bold(true)

	dashboardBox := style.New().
		Border(style.NormalBorder).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Width(60)

	dashboardContent := fmt.Sprintf(`%s  CPU Usage: %s
%s  Memory: %s
%s  Disk: %s
%s  Network: %s`,
		"📊", style.Render(statusOk, "45%"),
		"💾", style.Render(statusOk, "68%"),
		"💿", style.Render(statusWarning, "85%"),
		"🌐", style.Render(statusOk, "Connected"),
	)

	fmt.Println(style.Render(dashboardBox, dashboardContent))
	fmt.Println()

	fmt.Println("=== End of Complete Examples ===")
}
