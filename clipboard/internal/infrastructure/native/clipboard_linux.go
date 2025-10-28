//go:build linux

package native

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/model"
)

// Provider implements clipboard operations using Linux clipboard tools.
// Supports both X11 (xclip/xsel) and Wayland (wl-copy/wl-paste)
type Provider struct {
	readCmd  string
	writeCmd string
}

// NewProvider creates a new Linux native clipboard provider.
// Automatically detects available clipboard tools.
func NewProvider() *Provider {
	p := &Provider{}

	// Try Wayland first (wl-clipboard)
	if _, err := exec.LookPath("wl-copy"); err == nil {
		if _, err := exec.LookPath("wl-paste"); err == nil {
			p.writeCmd = "wl-copy"
			p.readCmd = "wl-paste"
			return p
		}
	}

	// Try X11 tools
	// xclip is more common
	if _, err := exec.LookPath("xclip"); err == nil {
		p.writeCmd = "xclip"
		p.readCmd = "xclip"
		return p
	}

	// xsel as fallback
	if _, err := exec.LookPath("xsel"); err == nil {
		p.writeCmd = "xsel"
		p.readCmd = "xsel"
		return p
	}

	return p
}

// Read reads content from the Linux clipboard.
func (p *Provider) Read() (*model.ClipboardContent, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("no clipboard tool available (install xclip, xsel, or wl-clipboard)")
	}

	var cmd *exec.Cmd

	switch p.readCmd {
	case "xclip":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "xsel":
		cmd = exec.Command("xsel", "--clipboard", "--output")
	case "wl-paste":
		cmd = exec.Command("wl-paste", "--no-newline")
	default:
		return nil, fmt.Errorf("unknown read command: %s", p.readCmd)
	}

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

// Write writes content to the Linux clipboard.
func (p *Provider) Write(content *model.ClipboardContent) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	if !p.IsAvailable() {
		return fmt.Errorf("no clipboard tool available (install xclip, xsel, or wl-clipboard)")
	}

	// Get text from content
	text, err := content.Text()
	if err != nil {
		return fmt.Errorf("only text content is supported: %w", err)
	}

	var cmd *exec.Cmd

	switch p.writeCmd {
	case "xclip":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-i")
	case "xsel":
		cmd = exec.Command("xsel", "--clipboard", "--input")
	case "wl-copy":
		cmd = exec.Command("wl-copy")
	default:
		return fmt.Errorf("unknown write command: %s", p.writeCmd)
	}

	cmd.Stdin = bytes.NewBufferString(text)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return nil
}

// IsAvailable returns true if a clipboard tool is available.
func (p *Provider) IsAvailable() bool {
	return p.readCmd != "" && p.writeCmd != ""
}

// Name returns the provider name.
func (p *Provider) Name() string {
	if !p.IsAvailable() {
		return "Linux Native (no tool available)"
	}

	if p.readCmd == "wl-paste" {
		return "Linux Native (Wayland wl-clipboard)"
	}

	return fmt.Sprintf("Linux Native (X11 %s)", p.readCmd)
}
