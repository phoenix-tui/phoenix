package service

import "github.com/phoenix-tui/phoenix/clipboard/domain/model"

// Provider is the interface that clipboard implementations must satisfy.
// This is a domain interface (hexagonal architecture port).
type Provider interface {
	// Read reads content from the clipboard
	Read() (*model.ClipboardContent, error)

	// Write writes content to the clipboard
	Write(content *model.ClipboardContent) error

	// IsAvailable returns true if the clipboard provider is available
	IsAvailable() bool

	// Name returns the name of the provider (e.g., "OSC52", "Windows Native")
	Name() string
}
