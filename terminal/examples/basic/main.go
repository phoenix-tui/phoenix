// Package main demonstrates Phoenix Terminal API usage.
//
// This example shows:.
//   - Auto-detection of terminal platform.
//   - Cursor positioning and visibility.
//   - ClearLines() operation (critical for multiline input).
//   - Capability discovery.
//   - ANSI escape code generation.
package main

import (
	"fmt"
	"time"

	"github.com/phoenix-tui/phoenix/terminal"
)

func main() {
	// Create platform-optimized terminal with auto-detection.
	term := terminal.New()

	fmt.Println("=== Phoenix Terminal API Demo ===")

	// Show platform detection.
	fmt.Printf("Platform: %s\n", term.Platform())
	fmt.Printf("Supports direct positioning: %v\n", term.SupportsDirectPositioning())
	fmt.Printf("Supports readback: %v\n", term.SupportsReadback())
	fmt.Printf("Supports true color: %v\n", term.SupportsTrueColor())
	fmt.Printf("Color depth: %d colors\n", term.ColorDepth())

	w, h, _ := term.Size()
	fmt.Printf("Terminal size: %dx%d\n\n", w, h)

	fmt.Println("Press Enter to start demo...")
	fmt.Scanln()

	// Demo 1: Cursor positioning.
	fmt.Println("\n--- Demo 1: Cursor Positioning ---")
	term.Clear()
	term.SetCursorPosition(0, 0)
	term.Write("Top-left (0,0)")

	term.SetCursorPosition(30, 5)
	term.Write("Middle (30,5)")

	term.SetCursorPosition(0, 10)
	term.Write("Press Enter to continue...")
	fmt.Scanln()

	// Demo 2: Cursor visibility.
	fmt.Println("\n--- Demo 2: Cursor Visibility ---")
	term.Clear()
	term.SetCursorPosition(0, 0)
	term.Write("Hiding cursor...")
	term.HideCursor()
	time.Sleep(2 * time.Second)

	term.SetCursorPosition(0, 1)
	term.Write("Showing cursor...")
	term.ShowCursor()
	time.Sleep(2 * time.Second)

	term.SetCursorPosition(0, 3)
	term.Write("Press Enter to continue...")
	fmt.Scanln()

	// Demo 3: ClearLines() - Critical for multiline input.
	fmt.Println("\n--- Demo 3: ClearLines (Multiline Clearing) ---")
	term.Clear()
	term.SetCursorPosition(0, 0)

	// Write some multiline content.
	lines := []string{
		"Line 1: This is the first line",
		"Line 2: This is the second line",
		"Line 3: This is the third line",
		"Line 4: This is the fourth line",
		"Line 5: This is the fifth line",
	}

	for _, line := range lines {
		term.Write(line + "\n")
	}

	term.SetCursorPosition(0, 7)
	term.Write("Press Enter to clear 3 lines...")
	fmt.Scanln()

	// Position cursor at line 3 (index 2).
	term.SetCursorPosition(0, 2)

	// Clear 3 lines (lines 2, 3, 4).
	term.ClearLines(3)
	term.Write("Cleared 3 lines! Notice lines 2-4 are gone.")

	term.SetCursorPosition(0, 10)
	term.Write("Press Enter to continue...")
	fmt.Scanln()

	// Demo 4: Screen operations.
	fmt.Println("\n--- Demo 4: Screen Operations ---")
	term.Clear()
	term.SetCursorPosition(0, 0)
	term.Write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	term.SetCursorPosition(0, 7)
	term.Write("Press Enter to clear from cursor...")
	fmt.Scanln()

	term.SetCursorPosition(0, 3)
	term.ClearFromCursor()
	term.Write("Cleared from line 3 to end!")

	term.SetCursorPosition(0, 5)
	term.Write("Press Enter to clear entire screen...")
	fmt.Scanln()

	term.Clear()
	term.SetCursorPosition(0, 0)
	term.Write("Screen cleared!")

	term.SetCursorPosition(0, 2)
	term.Write("Press Enter to continue...")
	fmt.Scanln()

	// Demo 5: WriteAt (optimized write at position).
	fmt.Println("\n--- Demo 5: WriteAt ---")
	term.Clear()

	// Draw a box using WriteAt.
	for i := 0; i < 20; i++ {
		term.WriteAt(10+i, 5, "-")
		term.WriteAt(10+i, 15, "-")
	}
	for i := 0; i < 10; i++ {
		term.WriteAt(10, 5+i, "|")
		term.WriteAt(29, 5+i, "|")
	}

	term.WriteAt(15, 10, "Hello!")
	term.SetCursorPosition(0, 17)
	term.Write("Press Enter to finish...")
	fmt.Scanln()

	// Cleanup.
	term.Clear()
	term.SetCursorPosition(0, 0)
	term.ShowCursor()

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nPhoenix Terminal API successfully demonstrated!")
	fmt.Printf("Platform: %s\n", term.Platform())
	fmt.Println("\nWeek 16 will add Windows Console API for 10x performance!")
}
