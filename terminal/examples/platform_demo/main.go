package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/phoenix-tui/phoenix/terminal/api"
	"github.com/phoenix-tui/phoenix/terminal/infrastructure"
)

func main() {
	fmt.Println("Phoenix Terminal Platform Demo")
	fmt.Println("================================")
	fmt.Println()

	// Auto-detect platform
	platform := infrastructure.DetectPlatform()
	fmt.Printf("Detected Platform: %s\n", platform)
	fmt.Println()

	// Create terminal with auto-detection
	term := infrastructure.NewTerminal()

	// Display capabilities
	fmt.Println("Terminal Capabilities:")
	fmt.Printf("  Direct Positioning:  %v\n", term.SupportsDirectPositioning())
	fmt.Printf("  Cursor Readback:     %v\n", term.SupportsReadback())
	fmt.Printf("  TrueColor Support:   %v\n", term.SupportsTrueColor())
	fmt.Printf("  Color Depth:         %d colors\n", term.ColorDepth())
	fmt.Println()

	// Get terminal size
	width, height, err := term.Size()
	if err != nil {
		fmt.Printf("  Size:                Error: %v\n", err)
	} else {
		fmt.Printf("  Size:                %dx%d\n", width, height)
	}
	fmt.Println()

	// Wait for user to read capabilities
	fmt.Println("Press Enter to run performance demo...")
	fmt.Scanln()

	// Run performance demonstration
	runPerformanceDemo(term, platform)

	// Wait before cleanup
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}

func runPerformanceDemo(term api.Terminal, platform api.Platform) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Performance Demonstration")
	fmt.Println(strings.Repeat("=", 60))

	// Save cursor position before demo
	term.SaveCursorPosition()
	defer term.RestoreCursorPosition()

	// Demo 1: Cursor Positioning
	fmt.Println("\n1. Cursor Positioning Speed")
	fmt.Println("   Moving cursor to 100 different positions...")
	start := time.Now()
	for i := 0; i < 100; i++ {
		x := i % 40
		y := (i / 40) + 5
		term.SetCursorPosition(x, y)
	}
	duration := time.Since(start)
	fmt.Printf("   Completed in: %v (avg: %v per operation)\n", duration, duration/100)

	// Demo 2: Cursor Readback (if supported)
	if term.SupportsReadback() {
		fmt.Println("\n2. Cursor Position Readback (Windows Console API only)")
		term.SetCursorPosition(25, 10)
		start = time.Now()
		x, y, err := term.GetCursorPosition()
		duration = time.Since(start)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
		} else {
			fmt.Printf("   Position: (%d, %d) in %v\n", x, y, duration)
		}
	} else {
		fmt.Println("\n2. Cursor Readback: Not supported (ANSI limitation)")
	}

	// Demo 3: Multiline Clearing (CRITICAL for GoSh)
	fmt.Println("\n3. Multiline Clearing Speed (Critical for GoSh)")
	fmt.Println("   Clearing 10 lines, repeated 100 times...")

	// Set up test position
	term.SetCursorPosition(0, 15)

	start = time.Now()
	for i := 0; i < 100; i++ {
		term.ClearLines(10)
		term.SetCursorPosition(0, 15) // Reset for next iteration
	}
	duration = time.Since(start)
	fmt.Printf("   Completed in: %v (avg: %v per operation)\n", duration, duration/100)

	if platform == api.PlatformWindowsConsole {
		fmt.Println("   ✓ Using Windows Console API (10x faster than ANSI!)")
	} else {
		fmt.Println("   ℹ Using ANSI escape codes (universal compatibility)")
	}

	// Demo 4: Screen Buffer Readback (if supported)
	if term.SupportsReadback() {
		fmt.Println("\n4. Screen Buffer Readback (Windows Console API only)")
		start = time.Now()
		buffer, err := term.ReadScreenBuffer()
		duration = time.Since(start)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
		} else {
			fmt.Printf("   Read %dx%d buffer in %v\n", len(buffer[0]), len(buffer), duration)
			fmt.Println("   ✓ Differential rendering possible!")
		}
	} else {
		fmt.Println("\n4. Screen Buffer Readback: Not supported (ANSI limitation)")
		fmt.Println("   ℹ Differential rendering not available on this platform")
	}

	// Demo 5: Cursor Visibility
	fmt.Println("\n5. Cursor Visibility Control")
	term.SetCursorPosition(0, 20)
	term.Write("Hiding cursor...")
	term.HideCursor()
	time.Sleep(1 * time.Second)

	term.SetCursorPosition(0, 20)
	term.Write("Showing cursor...")
	term.ShowCursor()
	time.Sleep(1 * time.Second)
	fmt.Println(" Done!")

	// Demo 6: Text Output Performance
	fmt.Println("\n6. Text Output Speed")
	text := "Phoenix Terminal Framework - Next-generation TUI for Go!"
	fmt.Println("   Writing text 1000 times...")

	start = time.Now()
	term.SetCursorPosition(0, 22)
	for i := 0; i < 1000; i++ {
		term.Write(text)
		term.SetCursorPosition(0, 22) // Overwrite same line
	}
	duration = time.Since(start)
	fmt.Printf("   Completed in: %v (avg: %v per write)\n", duration, duration/1000)

	// Final cleanup
	term.SetCursorPosition(0, 23)
	term.ClearFromCursor()

	// Performance summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Platform Summary")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Platform Type: %s\n", platform)

	switch platform {
	case api.PlatformWindowsConsole:
		fmt.Println("\n✓ OPTIMAL PERFORMANCE")
		fmt.Println("  • Using native Windows Console API")
		fmt.Println("  • 10x faster than ANSI for cursor operations")
		fmt.Println("  • 10x faster than ANSI for multiline clearing")
		fmt.Println("  • Cursor and buffer readback supported")
		fmt.Println("  • Perfect for GoSh multiline shell rendering")

	case api.PlatformWindowsANSI:
		fmt.Println("\nℹ GOOD PERFORMANCE (ANSI Fallback)")
		fmt.Println("  • Using ANSI escape codes")
		fmt.Println("  • Detected Git Bash / MinTTY environment")
		fmt.Println("  • Full compatibility with all terminals")
		fmt.Println("  • Cursor readback not available")
		fmt.Println("  • Still perfectly usable for GoSh")

	case api.PlatformUnix:
		fmt.Println("\nℹ GOOD PERFORMANCE (ANSI Standard)")
		fmt.Println("  • Using ANSI escape codes")
		fmt.Println("  • Universal Unix/Linux/macOS compatibility")
		fmt.Println("  • Cursor readback not available")
		fmt.Println("  • Excellent for terminal applications")

	default:
		fmt.Println("\n⚠ UNKNOWN PLATFORM")
		fmt.Println("  • Detection failed - using ANSI fallback")
	}

	fmt.Println()
	fmt.Println("Performance comparison (Windows Console API vs ANSI):")
	fmt.Println("  SetCursorPosition:    10x faster")
	fmt.Println("  ClearLines(10):       10x faster")
	fmt.Println("  GetCursorPosition:    Only on Windows Console")
	fmt.Println("  ReadScreenBuffer:     Only on Windows Console")
}
