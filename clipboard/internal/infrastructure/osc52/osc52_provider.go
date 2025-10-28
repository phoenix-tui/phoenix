// Package osc52 implements clipboard access via OSC 52 escape sequences for SSH sessions.
package osc52

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/model"
)

// Provider implements clipboard operations using OSC 52 escape sequences.
// This works over SSH connections by sending escape sequences to the terminal.
type Provider struct {
	timeout time.Duration
	output  *os.File // Terminal output (usually os.Stdout)
}

// NewProvider creates a new OSC 52 clipboard provider.
func NewProvider(timeout time.Duration) *Provider {
	return &Provider{
		timeout: timeout,
		output:  os.Stdout,
	}
}

// Read attempts to read from clipboard via OSC 52.
// Note: OSC 52 read is not widely supported, so this returns an error.
func (p *Provider) Read() (*model.ClipboardContent, error) {
	// OSC 52 clipboard reading is not widely supported by terminals
	// Most terminals only support writing via OSC 52
	return nil, fmt.Errorf("OSC 52 clipboard reading is not supported by most terminals")
}

// Write writes content to the clipboard using OSC 52 escape sequences.
func (p *Provider) Write(content *model.ClipboardContent) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	// Get the data to write
	data := content.Data()

	// OSC 52 format: ESC ] 52 ; c ; <base64 data> ESC \
	// Where 'c' is the clipboard selection (c = clipboard, p = primary)
	encoded := base64.StdEncoding.EncodeToString(data)

	// Build the OSC 52 sequence
	// Using 'c' for clipboard (as opposed to 'p' for primary selection)
	sequence := fmt.Sprintf("\033]52;c;%s\033\\", encoded)

	// Write to terminal with timeout
	done := make(chan error, 1)
	go func() {
		_, err := p.output.WriteString(sequence)
		if err != nil {
			done <- err
			return
		}
		// Flush is important to ensure the escape sequence is sent immediately
		done <- p.output.Sync()
	}()

	// Wait for write to complete or timeout
	select {
	case err := <-done:
		return err
	case <-time.After(p.timeout):
		return fmt.Errorf("OSC 52 write timeout after %v", p.timeout)
	}
}

// IsAvailable returns true if OSC 52 can be used.
func (p *Provider) IsAvailable() bool {
	// Check if we have a valid output file
	if p.output == nil {
		return false
	}

	// Check if the output is a terminal
	fileInfo, err := p.output.Stat()
	if err != nil {
		return false
	}

	// Must be a character device (terminal)
	isTerminal := (fileInfo.Mode() & os.ModeCharDevice) != 0
	if !isTerminal {
		return false
	}

	// Check for SSH session indicators
	sshIndicators := []string{"SSH_TTY", "SSH_CLIENT", "SSH_CONNECTION"}
	for _, indicator := range sshIndicators {
		if os.Getenv(indicator) != "" {
			return true
		}
	}

	// Also check if TERM supports OSC 52
	term := os.Getenv("TERM")
	supportedTerms := []string{"xterm", "xterm-256color", "screen", "tmux", "tmux-256color"}
	for _, supported := range supportedTerms {
		if term == supported {
			return true
		}
	}

	return false
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "OSC52"
}

// WithOutput sets a custom output file (useful for testing).
func (p *Provider) WithOutput(output *os.File) *Provider {
	p.output = output
	return p
}
