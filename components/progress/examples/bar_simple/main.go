package main

import (
	"fmt"
	"time"

	progress "github.com/phoenix-tui/phoenix/components/progress/api"
)

// Simple progress bar example
// Demonstrates basic usage without tea.Program
func main() {
	fmt.Println("Simple Progress Bar Example")
	fmt.Println("============================")

	// Create a 40-character wide progress bar
	bar := progress.NewBar(40)

	// Simulate progress from 0% to 100%
	for i := 0; i <= 100; i += 10 {
		bar.SetProgress(i)
		fmt.Printf("\r%s", bar.View())
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n\nDone!")
}
