// Package main demonstrates basic usage of the Phoenix layout API.
package main

import (
	"fmt"

	layout "github.com/phoenix-tui/phoenix/layout/api"
)

func main() {
	fmt.Println("=== Basic Layout Examples ===")

	// Example 1: Simple text
	fmt.Println("1. Simple Text:")
	simple := layout.NewBox("Hello World").Render()
	fmt.Println(simple)
	fmt.Println()

	// Example 2: Text with padding
	fmt.Println("2. Text with Padding:")
	padded := layout.NewBox("Hello World").
		PaddingAll(2).
		Render()
	fmt.Println(padded)
	fmt.Println()

	// Example 3: Text with border
	fmt.Println("3. Text with Border:")
	bordered := layout.NewBox("Hello World").
		Border().
		Render()
	fmt.Println(bordered)
	fmt.Println()

	// Example 4: Text with border and padding
	fmt.Println("4. Text with Border and Padding:")
	full := layout.NewBox("Hello World").
		PaddingAll(1).
		Border().
		Render()
	fmt.Println(full)
	fmt.Println()

	// Example 5: Multi-line content
	fmt.Println("5. Multi-line Content:")
	multiline := layout.NewBox("Line 1\nLine 2\nLine 3").
		PaddingAll(1).
		Border().
		Render()
	fmt.Println(multiline)
	fmt.Println()

	// Example 6: With margin
	fmt.Println("6. With Margin:")
	withMargin := layout.NewBox("Centered Text").
		PaddingAll(1).
		Border().
		MarginAll(2).
		Render()
	fmt.Println(withMargin)
	fmt.Println()

	// Example 7: Unicode content
	fmt.Println("7. Unicode Content:")
	unicode := layout.NewBox("Hello 世界 👋").
		PaddingVH(1, 2).
		Border().
		Render()
	fmt.Println(unicode)
	fmt.Println()

	// Example 8: Size constraints
	fmt.Println("8. Size Constraints:")
	sized := layout.NewBox("Short").
		Width(30).
		PaddingAll(1).
		Border().
		Render()
	fmt.Println(sized)
}
