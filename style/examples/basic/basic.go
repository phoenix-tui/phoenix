package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/style/api"
)

func main() {
	fmt.Println("=== Phoenix Style Library - Basic Examples ===\n")

	// Example 1: Simple colored text
	fmt.Println("1. Colored Text:")
	redStyle := style.New().Foreground(style.Red)
	fmt.Println(style.Render(redStyle, "This is red text"))

	greenStyle := style.New().Foreground(style.Green)
	fmt.Println(style.Render(greenStyle, "This is green text"))

	blueStyle := style.New().Foreground(style.Blue)
	fmt.Println(style.Render(blueStyle, "This is blue text"))
	fmt.Println()

	// Example 2: Text decorations
	fmt.Println("2. Text Decorations:")
	fmt.Println(style.Render(style.BoldStyle, "Bold text"))
	fmt.Println(style.Render(style.ItalicStyle, "Italic text"))
	fmt.Println(style.Render(style.UnderlineStyle, "Underlined text"))
	fmt.Println(style.Render(style.StrikethroughStyle, "Strikethrough text"))
	fmt.Println()

	// Example 3: Combined decorations
	fmt.Println("3. Combined:")
	combined := style.New().
		Foreground(style.Magenta).
		Bold(true).
		Underline(true)
	fmt.Println(style.Render(combined, "Bold + Underline + Magenta"))
	fmt.Println()

	// Example 4: Background colors
	fmt.Println("4. Background Colors:")
	whiteOnBlue := style.New().
		Foreground(style.White).
		Background(style.Blue)
	fmt.Println(style.Render(whiteOnBlue, " White text on blue background "))
	fmt.Println()

	// Example 5: Simple border
	fmt.Println("5. Simple Border:")
	bordered := style.New().Border(style.RoundedBorder)
	fmt.Println(style.Render(bordered, "Boxed text"))
	fmt.Println()

	// Example 6: Border with color
	fmt.Println("6. Colored Border:")
	coloredBorder := style.New().
		Border(style.RoundedBorder).
		BorderColor(style.Cyan).
		Foreground(style.Yellow)
	fmt.Println(style.Render(coloredBorder, "Cyan border + Yellow text"))
	fmt.Println()

	// Example 7: Padding
	fmt.Println("7. With Padding:")
	padded := style.New().
		Border(style.NormalBorder).
		Padding(style.NewPadding(1, 2, 1, 2))
	fmt.Println(style.Render(padded, "Padded content"))
	fmt.Println()

	// Example 8: Margin
	fmt.Println("8. With Margin:")
	margined := style.New().
		Border(style.RoundedBorder).
		Margin(style.NewMargin(1, 0, 1, 0))
	fmt.Println(style.Render(margined, "Margined box"))
	fmt.Println()

	// Example 9: Alignment
	fmt.Println("9. Centered Text:")
	centered := style.New().
		Width(40).
		Align(style.NewAlignment(style.AlignCenter, style.AlignTop)).
		Foreground(style.Yellow)
	fmt.Println(style.Render(centered, "This is centered"))
	fmt.Println()

	// Example 10: Unicode content
	fmt.Println("10. Unicode Support:")
	unicode := style.New().
		Border(style.RoundedBorder).
		Padding(style.NewPadding(0, 1, 0, 1)).
		Foreground(style.Cyan)
	fmt.Println(style.Render(unicode, "Hello 👋 World 🌍"))
	fmt.Println()

	fmt.Println("=== End of Basic Examples ===")
}
