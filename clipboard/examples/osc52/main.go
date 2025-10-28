// Package main demonstrates OSC 52 clipboard access for SSH sessions.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard"
)

func main() {
	fmt.Println("Phoenix Clipboard - OSC 52 Example")
	fmt.Println("===================================")
	fmt.Println()

	// Display environment information
	fmt.Println("Environment Information:")
	fmt.Println("------------------------")
	fmt.Printf("SSH_TTY: %s\n", os.Getenv("SSH_TTY"))
	fmt.Printf("SSH_CLIENT: %s\n", os.Getenv("SSH_CLIENT"))
	fmt.Printf("SSH_CONNECTION: %s\n", os.Getenv("SSH_CONNECTION"))
	fmt.Printf("TERM: %s\n", os.Getenv("TERM"))
	fmt.Println()

	// Create clipboard with explicit OSC 52 configuration
	cb, err := clipboard.NewBuilder().
		WithOSC52(true).
		WithOSC52Timeout(5 * time.Second).
		Build()
	if err != nil {
		log.Fatalf("Failed to create clipboard: %v", err)
	}

	// Check availability
	if !cb.IsAvailable() {
		log.Fatal("Clipboard is not available")
	}

	fmt.Printf("Using provider: %s\n", cb.GetProviderName())
	fmt.Printf("SSH session: %v\n", cb.IsSSH())
	fmt.Println()

	// Write text to clipboard using OSC 52
	textToWrite := "Phoenix TUI clipboard via OSC 52! 🎉"
	fmt.Printf("Writing to clipboard: %s\n", textToWrite)

	err = cb.Write(textToWrite)
	if err != nil {
		log.Fatalf("Failed to write to clipboard: %v", err)
	}

	fmt.Println("✓ Successfully sent OSC 52 escape sequence")
	fmt.Println()

	// Note about OSC 52
	fmt.Println("Note:")
	fmt.Println("-----")
	fmt.Println("OSC 52 sends an escape sequence to the terminal to set clipboard content.")
	fmt.Println("This works over SSH connections and allows clipboard sync between:")
	fmt.Println("  - Remote server → Local machine clipboard")
	fmt.Println()
	fmt.Println("Supported terminals:")
	fmt.Println("  - xterm, iTerm2, Windows Terminal, tmux, screen, and many others")
	fmt.Println()
	fmt.Println("Check your local clipboard to verify the text was synchronized!")
	fmt.Println()

	// Try to read (will likely fail as OSC 52 read is not widely supported)
	fmt.Println("Attempting to read (OSC 52 read not widely supported)...")
	text, err := cb.Read()
	if err != nil {
		fmt.Printf("✗ Read failed (expected): %v\n", err)
		fmt.Println()
		fmt.Println("This is normal - most terminals only support OSC 52 write, not read.")
	} else {
		fmt.Printf("✓ Successfully read: %s\n", text)
	}

	// Using OSC 52 only (no fallback to native)
	fmt.Println()
	fmt.Println("OSC 52 Only Mode:")
	fmt.Println("-----------------")

	osc52Only, err := clipboard.NewBuilder().
		WithOSC52(true).
		WithNative(false).
		Build()
	if err != nil {
		log.Fatalf("Failed to create OSC 52-only clipboard: %v", err)
	}

	if !osc52Only.IsAvailable() {
		fmt.Println("✗ OSC 52 not available (terminal may not support it)")
		return
	}

	fmt.Printf("Provider: %s\n", osc52Only.GetProviderName())
	if err := osc52Only.Write("OSC 52 only mode test"); err != nil {
		fmt.Printf("✗ Write failed: %v\n", err)
	} else {
		fmt.Println("✓ Write successful")
	}
}
