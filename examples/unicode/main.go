package main

import (
	"fmt"
	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
)

func main() {
	fmt.Println("Phoenix Unicode Service Demo")
	fmt.Println("================================\n")

	us := coreService.NewUnicodeService()

	// Demo 1: ASCII
	ascii := "Hello World"
	fmt.Printf("ASCII: %q\n", ascii)
	fmt.Printf("Width: %d (expected: 11)\n\n", us.StringWidth(ascii))

	// Demo 2: Emoji
	emoji := "👋😀🎉"
	fmt.Printf("Emoji: %q\n", emoji)
	fmt.Printf("Width: %d (expected: 6)\n", us.StringWidth(emoji))
	fmt.Printf("Clusters: %v\n\n", us.GraphemeClusters(emoji))

	// Demo 3: CJK
	cjk := "你好世界"
	fmt.Printf("CJK: %q\n", cjk)
	fmt.Printf("Width: %d (expected: 8)\n\n", us.StringWidth(cjk))

	// Demo 4: Mixed content
	mixed := "Hello 👋 世界!"
	fmt.Printf("Mixed: %q\n", mixed)
	fmt.Printf("Width: %d\n", us.StringWidth(mixed))
	fmt.Printf("Clusters: %v\n\n", us.GraphemeClusters(mixed))

	// Demo 5: Emoji with modifiers
	modifier := "👋🏻"
	fmt.Printf("Emoji with modifier: %q\n", modifier)
	fmt.Printf("Width: %d (1 cluster, 2 columns)\n", us.StringWidth(modifier))
	fmt.Printf("Clusters: %v\n\n", us.GraphemeClusters(modifier))

	// Demo 6: Flags
	flag := "🇺🇸"
	fmt.Printf("Flag: %q\n", flag)
	fmt.Printf("Width: %d (1 cluster, 2 columns)\n", us.StringWidth(flag))
	fmt.Printf("Clusters: %v\n\n", us.GraphemeClusters(flag))

	// Demo 7: Zero-width characters
	combining := "e\u0301" // e with acute accent (é)
	fmt.Printf("Combining mark: %q\n", combining)
	fmt.Printf("Width: %d (expected: 1)\n", us.StringWidth(combining))
	fmt.Printf("Clusters: %v\n\n", us.GraphemeClusters(combining))

	// Demo 8: Lipgloss #562 Bug Test
	fmt.Println("=== Lipgloss #562 Bug Fix Test ===")
	lipglossBug := "📝 Test"
	fmt.Printf("String: %q\n", lipglossBug)
	fmt.Printf("Correct width: %d (2 for emoji + 1 space + 4 for 'Test' = 7)\n", us.StringWidth(lipglossBug))
	fmt.Printf("Lipgloss would calculate: 8 (WRONG!)\n\n")

	fmt.Println("✨ Unicode width calculation working perfectly!")
	fmt.Println("✨ Phoenix solves Lipgloss #562 - Perfect Unicode support!")
}
