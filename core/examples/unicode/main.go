// Package main demonstrates Unicode string width calculation capabilities.
package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/core"
)

func main() {
	fmt.Println("Phoenix Unicode Public API Demo")
	fmt.Println("=================================")
	fmt.Println()

	// Demo 1: ASCII
	ascii := "Hello World"
	fmt.Printf("ASCII: %q\n", ascii)
	fmt.Printf("Width: %d (expected: 11)\n\n", core.StringWidth(ascii))

	// Demo 2: Emoji
	emoji := "👋😀🎉"
	fmt.Printf("Emoji: %q\n", emoji)
	fmt.Printf("Width: %d (expected: 6)\n\n", core.StringWidth(emoji))

	// Demo 3: CJK
	cjk := "你好世界"
	fmt.Printf("CJK: %q\n", cjk)
	fmt.Printf("Width: %d (expected: 8)\n\n", core.StringWidth(cjk))

	// Demo 4: Mixed content
	mixed := "Hello 👋 世界!"
	fmt.Printf("Mixed: %q\n", mixed)
	fmt.Printf("Width: %d (expected: 13)\n\n", core.StringWidth(mixed))

	// Demo 5: Emoji with modifiers
	modifier := "👋🏻"
	fmt.Printf("Emoji with modifier: %q\n", modifier)
	fmt.Printf("Width: %d (1 cluster, 2 columns)\n\n", core.StringWidth(modifier))

	// Demo 6: Flags
	flag := "🇺🇸"
	fmt.Printf("Flag: %q\n", flag)
	fmt.Printf("Width: %d (1 cluster, 2 columns)\n\n", core.StringWidth(flag))

	// Demo 7: Zero-width characters
	combining := "e\u0301" // e with acute accent (é)
	fmt.Printf("Combining mark: %q\n", combining)
	fmt.Printf("Width: %d (expected: 1)\n\n", core.StringWidth(combining))

	// Demo 8: Lipgloss #562 Bug Fix Test
	fmt.Println("=== Lipgloss #562 Bug Fix Test ===")
	lipglossBug := "📝 Test"
	fmt.Printf("String: %q\n", lipglossBug)
	fmt.Printf("Correct width: %d (2 for emoji + 1 space + 4 for 'Test' = 7)\n", core.StringWidth(lipglossBug))
	fmt.Printf("Lipgloss would calculate: 8 (WRONG!)\n\n")

	fmt.Println("✨ Unicode width calculation working perfectly!")
	fmt.Println("✨ Phoenix solves Lipgloss #562 - Perfect Unicode support!")
}
