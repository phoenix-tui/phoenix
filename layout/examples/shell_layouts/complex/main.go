// Package main demonstrates complex nested layouts for advanced shell UIs.
// This example shows how to combine row and column layouts.
package main

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/layout"
)

func main() {
	const (
		terminalWidth  = 80
		terminalHeight = 24
	)

	// Example 1: Split pane shell (like tmux)
	// - Left pane: Main terminal
	// - Right pane: Help sidebar

	fmt.Println("=== Split Pane Layout (Left: Terminal, Right: Help) ===")
	// For now, show a simpler version
	left := layout.NewBox("Main Terminal Area\n\nCommand history...\nMore history...\n\n$ current input")
	right := layout.NewBox("Help\n\nls - list\ncd - change\nexit - quit")

	split := layout.Row().
		Gap(2).
		JustifyStart().
		Add(left).
		Add(right)

	fmt.Println(split.Render(terminalWidth, 12))
	fmt.Println()

	// Example 2: Dashboard layout
	// - Top: Title bar
	// - Middle: Three columns (stats)
	// - Bottom: Status bar

	title := layout.NewBox("=== System Monitor ===").
		AlignCenter().
		Border()

	cpu := layout.NewBox("CPU\n45%").Border().PaddingAll(1)
	memory := layout.NewBox("Memory\n2.1GB").Border().PaddingAll(1)
	disk := layout.NewBox("Disk\n78%").Border().PaddingAll(1)

	stats := layout.Row().
		Gap(3).
		JustifyCenter().
		Add(cpu).
		Add(memory).
		Add(disk)

	// For now we need to manually render stats
	statsBox := layout.NewBox(stats.Render(terminalWidth, 5))

	status := layout.NewBox("Status: Running | Uptime: 3d 12h").
		AlignCenter()

	dashboard := layout.Column().
		Gap(1).
		JustifyStart().
		Add(title).
		Add(statsBox).
		Add(status)

	fmt.Println("=== Dashboard Layout ===")
	fmt.Println(dashboard.Render(terminalWidth, 15))
	fmt.Println()

	// Example 3: Overlay (completion menu over input)
	// This demonstrates absolute positioning concept

	baseInput := layout.NewBox("$ git che").
		Border()

	// Completion menu (would be positioned absolutely in real impl)
	completions := layout.NewBox("Completions:\n  checkout\n  cherry-pick\n  check").
		Border().
		PaddingAll(1)

	fmt.Println("=== Completion Menu Overlay Concept ===")
	fmt.Println("Base input:")
	fmt.Println(baseInput.Render())
	fmt.Println()
	fmt.Println("Completion menu (would overlay above):")
	fmt.Println(completions.Render())
	fmt.Println()
	fmt.Println("Note: True overlay requires absolute positioning (future feature)")
}
