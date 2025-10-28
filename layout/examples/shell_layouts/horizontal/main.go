// Package main demonstrates horizontal split layout for shell applications.
// This example shows how to create a horizontal layout with prompt and input areas.
package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/layout"
)

func main() {
	// Simulate a shell with horizontal split:
	// - Left: Prompt area (e.g., "$ ")
	// - Right: Input area (user typing)

	const (
		terminalWidth  = 80
		terminalHeight = 1 // Single line for input
	)

	// Create prompt box (fixed width)
	prompt := layout.NewBox("$ ").
		Width(2)

	// Create input box (flexible width)
	input := layout.NewBox("echo 'Hello World'")

	// Create horizontal layout
	shell := layout.Row().
		Gap(0). // No gap between prompt and input
		JustifyStart().
		AlignStart().
		Add(prompt).
		Add(input)

	// Render the shell
	output := shell.Render(terminalWidth, terminalHeight)

	fmt.Println("=== Horizontal Split (Prompt + Input) ===")
	fmt.Println(output)
	fmt.Println()

	// Example 2: Status bar split (left info + right info)
	leftStatus := layout.NewBox("[gosh v1.0.0]")
	rightStatus := layout.NewBox("[12:34 PM]")

	statusBar := layout.Row().
		Gap(1).
		JustifySpaceBetween(). // Push items to edges
		AlignCenter().
		Add(leftStatus).
		Add(rightStatus)

	fmt.Println("=== Status Bar (Left + Right) ===")
	fmt.Println(statusBar.Render(terminalWidth, 1))
	fmt.Println()

	// Example 3: Three-column layout
	col1 := layout.NewBox("Column 1")
	col2 := layout.NewBox("Column 2")
	col3 := layout.NewBox("Column 3")

	threeColumn := layout.Row().
		Gap(5). // 5 cells between columns
		JustifyCenter().
		AlignStart().
		Add(col1).
		Add(col2).
		Add(col3)

	fmt.Println("=== Three Columns (with gaps) ===")
	fmt.Println(threeColumn.Render(terminalWidth, 3))
}
