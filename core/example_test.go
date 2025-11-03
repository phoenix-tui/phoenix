package core_test

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/core"
)

// Example demonstrates basic terminal capability detection.
// This shows how to detect terminal features like color support,
// ANSI capabilities, and terminal size.
func Example() {
	// Auto-detect terminal capabilities from environment
	term := core.AutoDetect()
	caps := term.Capabilities()

	// Check color support
	fmt.Printf("Color support: %t\n", caps.SupportsColor())
	fmt.Printf("True color: %t\n", caps.SupportsTrueColor())
	fmt.Printf("ANSI: %t\n", caps.SupportsANSI())

	// Get terminal size
	size := term.Size()
	fmt.Printf("Size: %dx%d\n", size.Width, size.Height)

	// Output:
	// Color support: true
	// True color: false
	// ANSI: true
	// Size: 80x24
}

// ExampleStringWidth demonstrates correct Unicode width calculation.
// This function handles emoji, CJK characters, and combining marks properly,
// fixing the infamous Lipgloss #562 bug.
func ExampleStringWidth() {
	// ASCII string
	ascii := "Hello"
	fmt.Printf("ASCII width: %d\n", core.StringWidth(ascii))

	// Emoji (counts as 2 cells in terminal)
	emoji := "👋"
	fmt.Printf("Emoji width: %d\n", core.StringWidth(emoji))

	// CJK characters (2 cells each)
	cjk := "你好"
	fmt.Printf("CJK width: %d\n", core.StringWidth(cjk))

	// Mixed content
	mixed := "Hello 👋 世界"
	fmt.Printf("Mixed width: %d\n", core.StringWidth(mixed))

	// Output:
	// ASCII width: 5
	// Emoji width: 2
	// CJK width: 4
	// Mixed width: 13
}

// ExampleNewCellAuto demonstrates creating terminal cells with automatic width detection.
// Cells represent individual terminal screen positions with content and styling.
func ExampleNewCellAuto() {
	// Create cell with ASCII character
	ascii := core.NewCellAuto("A")
	fmt.Printf("ASCII cell: '%s' width=%d\n", ascii.Content, ascii.Width)

	// Create cell with emoji (automatically detects width=2)
	emoji := core.NewCellAuto("👋")
	fmt.Printf("Emoji cell: '%s' width=%d\n", emoji.Content, emoji.Width)

	// Create cell with CJK character
	cjk := core.NewCellAuto("你")
	fmt.Printf("CJK cell: '%s' width=%d\n", cjk.Content, cjk.Width)

	// Output:
	// ASCII cell: 'A' width=1
	// Emoji cell: '👋' width=2
	// CJK cell: '你' width=2
}
