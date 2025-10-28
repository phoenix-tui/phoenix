// Package main demonstrates vertical split layout for shell applications.
// This example shows how to create a vertical layout with history and input areas.
package main

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/layout"
)

func main() {
	const (
		terminalWidth  = 80
		terminalHeight = 24
	)

	// Example 1: Classic shell (history + input)
	// - Top: Command history (scrollable area)
	// - Bottom: Input line (fixed height)

	historyLines := []string{
		"$ ls -la",
		"total 42",
		"drwxr-xr-x  5 user  staff  160 Oct 16 12:00 .",
		"drwxr-xr-x 20 user  staff  640 Oct 15 10:00 ..",
		"-rw-r--r--  1 user  staff 1234 Oct 16 11:30 README.md",
		"$ cd projects",
		"$ git status",
		"On branch main",
		"Your branch is up to date with 'origin/main'.",
		"nothing to commit, working tree clean",
	}

	history := layout.NewBox(strings.Join(historyLines, "\n"))

	input := layout.NewBox("$ echo 'Hello World'").
		Border()

	shell := layout.Column().
		Gap(0). // No gap between history and input
		JustifyStart().
		AlignStretch(). // Stretch to full width
		Add(history).
		Add(input)

	fmt.Println("=== Vertical Split (History + Input) ===")
	fmt.Println(shell.Render(terminalWidth, terminalHeight))
	fmt.Println()

	// Example 2: Three-panel layout (header + content + footer)
	header := layout.NewBox("=== GOSH Shell v1.0.0 ===").
		Border().
		AlignCenter()

	content := layout.NewBox("Command history appears here...\nExecuted commands will be logged.\nScroll up to view history.")

	footer := layout.NewBox("Press Ctrl+C to exit | Tab for completion | ↑/↓ for history").
		AlignCenter()

	app := layout.Column().
		Gap(1). // 1 line gap between sections
		JustifyStart().
		Add(header).
		Add(content).
		Add(footer)

	fmt.Println("=== Three-Panel Layout (Header + Content + Footer) ===")
	fmt.Println(app.Render(terminalWidth, 15))
	fmt.Println()

	// Example 3: Centered dialog
	dialog := layout.NewBox("Are you sure you want to exit?").
		Border().
		PaddingAll(2)

	centered := layout.Column().
		JustifyCenter(). // Center vertically
		AlignCenter().   // Center horizontally
		Add(dialog)

	fmt.Println("=== Centered Dialog ===")
	fmt.Println(centered.Render(terminalWidth, 10))
}
