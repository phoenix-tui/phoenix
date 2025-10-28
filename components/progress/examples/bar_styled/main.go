// Package main demonstrates styled progress bar usage.
// This example shows progress bars with custom colors and formatting.
package main

import (
	"fmt"
	"time"

	"github.com/phoenix-tui/phoenix/components/progress"
)

// Styled progress bar example.
// Demonstrates customization with labels, percentage, and custom characters.
func main() {
	fmt.Println("Styled Progress Bar Example")
	fmt.Println("============================")

	// Create a styled progress bar.
	bar := progress.NewBar(50).
		FillChar('█').
		EmptyChar('░').
		ShowPercent(true).
		Label("Downloading...")

	// Simulate download progress.
	for i := 0; i <= 100; i += 5 {
		bar.SetProgress(i)
		fmt.Printf("\r%s", bar.View())
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()

	// Example 2: Custom styling with different characters.
	bar2 := progress.NewBar(50).
		FillChar('▓').
		EmptyChar('▒').
		ShowPercent(true).
		Label("Processing...")

	for i := 0; i <= 100; i += 10 {
		bar2.SetProgress(i)
		fmt.Printf("\r%s", bar2.View())
		time.Sleep(150 * time.Millisecond)
	}

	fmt.Println()

	// Example 3: Minimal bar (no label, no percentage)
	bar3 := progress.NewBar(30)

	for i := 0; i <= 100; i += 20 {
		bar3.SetProgress(i)
		fmt.Printf("\r%s", bar3.View())
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n\nDone!")
}
