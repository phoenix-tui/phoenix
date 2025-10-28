// Package main demonstrates basic clipboard read/write operations.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/clipboard"
)

func main() {
	fmt.Println("Phoenix Clipboard - Basic Example")
	fmt.Println("==================================")
	fmt.Println()

	// Create clipboard instance
	cb, err := clipboard.New()
	if err != nil {
		log.Fatalf("Failed to create clipboard: %v", err)
	}

	// Check if clipboard is available
	if !cb.IsAvailable() {
		log.Fatal("Clipboard is not available on this system")
	}

	fmt.Printf("Using provider: %s\n", cb.GetProviderName())
	fmt.Println()

	// Write to clipboard
	textToWrite := "Hello from Phoenix TUI Framework! 🚀"
	fmt.Printf("Writing to clipboard: %s\n", textToWrite)

	err = cb.Write(textToWrite)
	if err != nil {
		log.Fatalf("Failed to write to clipboard: %v", err)
	}

	fmt.Println("✓ Successfully written to clipboard")
	fmt.Println()

	// Read from clipboard
	fmt.Println("Reading from clipboard...")

	textRead, err := cb.Read()
	if err != nil {
		log.Fatalf("Failed to read from clipboard: %v", err)
	}

	fmt.Printf("✓ Read from clipboard: %s\n", textRead)
	fmt.Println()

	// Verify
	if textRead == textToWrite {
		fmt.Println("✓ Clipboard operation successful!")
	} else {
		fmt.Println("✗ Clipboard content doesn't match")
	}

	// Using global convenience functions
	fmt.Println()
	fmt.Println("Using global convenience functions:")
	fmt.Println("-----------------------------------")

	err = clipboard.Write("Global function test")
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}

	text, err := clipboard.Read()
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}

	fmt.Printf("✓ Read using global function: %s\n", text)
}
