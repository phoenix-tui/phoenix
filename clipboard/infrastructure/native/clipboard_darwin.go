//go:build darwin

package native

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
)

// Provider implements clipboard operations using macOS pbcopy/pbpaste.
type Provider struct{}

// NewProvider creates a new macOS native clipboard provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Read reads content from the macOS clipboard using pbpaste.
func (p *Provider) Read() (*model.ClipboardContent, error) {
	// Use pbpaste to read from clipboard
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to read from clipboard: %w", err)
	}

	text := out.String()
	if text == "" {
		return nil, fmt.Errorf("clipboard is empty")
	}

	return model.NewTextContent(text)
}

// Write writes content to the macOS clipboard using pbcopy.
func (p *Provider) Write(content *model.ClipboardContent) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	// Get text from content
	text, err := content.Text()
	if err != nil {
		return fmt.Errorf("only text content is supported: %w", err)
	}

	// Use pbcopy to write to clipboard
	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewBufferString(text)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return nil
}

// IsAvailable returns true if pbcopy/pbpaste are available.
func (p *Provider) IsAvailable() bool {
	// Check if pbcopy and pbpaste are available
	_, err := exec.LookPath("pbcopy")
	if err != nil {
		return false
	}

	_, err = exec.LookPath("pbpaste")
	if err != nil {
		return false
	}

	return true
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "macOS Native (pbcopy/pbpaste)"
}
