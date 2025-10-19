// Package main demonstrates basic phoenix/core usage.
package main

import (
	"fmt"

	core "github.com/phoenix-tui/phoenix/core/api"
)

func main() {
	fmt.Println("Phoenix Core - Basic Example")
	fmt.Println("=============================\n")

	// Auto-detect terminal capabilities
	term := core.AutoDetect()
	caps := term.Capabilities()

	// Display terminal information
	fmt.Printf("Terminal Size: %dx%d\n", term.Size().Width, term.Size().Height)
	fmt.Printf("ANSI Support: %v\n", caps.SupportsANSI())
	fmt.Printf("Color Depth: %s\n", caps.ColorDepth().String())
	fmt.Printf("True Color: %v\n", caps.SupportsTrueColor())
	fmt.Printf("Mouse Support: %v\n", caps.SupportsMouse())
	fmt.Printf("Alt Screen: %v\n", caps.SupportsAltScreen())
	fmt.Printf("Cursor Control: %v\n", caps.SupportsCursorControl())

	fmt.Println("\n--- Position Operations ---")
	pos := core.NewPosition(5, 10)
	fmt.Printf("Initial position: Row=%d, Col=%d\n", pos.Row, pos.Col)

	moved := pos.Add(3, 5)
	fmt.Printf("After Add(3, 5): Row=%d, Col=%d\n", moved.Row, moved.Col)

	fmt.Println("\n--- Size Operations ---")
	size := core.NewSize(120, 40)
	fmt.Printf("Custom size: %dx%d\n", size.Width, size.Height)

	// Demonstrate immutability
	resized := term.WithSize(size)
	fmt.Printf("Original terminal size: %dx%d\n", term.Size().Width, term.Size().Height)
	fmt.Printf("Resized terminal size: %dx%d\n", resized.Size().Width, resized.Size().Height)

	fmt.Println("\n--- Cell Operations ---")

	// Manual width (old way - still supported for advanced use)
	ascii := core.NewCell("A", 1)
	fmt.Printf("ASCII cell (manual) '%s' has width %d\n", ascii.Content, ascii.Width)

	// Automatic width (new way - recommended!)
	fmt.Println("\n🔥 Auto-width calculation (NEW!):")
	asciiAuto := core.NewCellAuto("A")
	emojiAuto := core.NewCellAuto("😀")
	cjkAuto := core.NewCellAuto("中")
	mixedAuto := core.NewCellAuto("Hi 👋")
	fmt.Printf("  ASCII cell (auto) '%s' has width %d\n", asciiAuto.Content, asciiAuto.Width)
	fmt.Printf("  Emoji cell (auto) '%s' has width %d\n", emojiAuto.Content, emojiAuto.Width)
	fmt.Printf("  CJK cell (auto) '%s' has width %d\n", cjkAuto.Content, cjkAuto.Width)
	fmt.Printf("  Mixed cell (auto) '%s' has width %d\n", mixedAuto.Content, mixedAuto.Width)

	fmt.Println("\n--- Manual Capabilities ---")
	manual := core.NewCapabilities(true, core.ColorDepth256, true, true, true)
	fmt.Printf("Manual capabilities: %s with mouse support\n", manual.ColorDepth().String())

	fmt.Println("\n✨ Phoenix Core working perfectly!")
}
