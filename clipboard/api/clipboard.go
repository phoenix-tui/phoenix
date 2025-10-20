// Package api provides the public interface for cross-platform clipboard operations.
package api

import (
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/application"
	"github.com/phoenix-tui/phoenix/clipboard/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/native"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/osc52"
)

// Clipboard is the public API for clipboard operations.
type Clipboard struct {
	manager *application.ClipboardManager
}

// New creates a new clipboard instance with auto-detected providers.
// This is the recommended way to use the clipboard.
func New() (*Clipboard, error) {
	manager, err := application.NewClipboardManager()
	if err != nil {
		return nil, err
	}

	return &Clipboard{
		manager: manager,
	}, nil
}

// Read reads text from the clipboard.
func (c *Clipboard) Read() (string, error) {
	return c.manager.Read()
}

// Write writes text to the clipboard.
func (c *Clipboard) Write(text string) error {
	return c.manager.Write(text)
}

// IsAvailable returns true if clipboard is available.
func (c *Clipboard) IsAvailable() bool {
	return c.manager.IsAvailable()
}

// GetProviderName returns the name of the active provider.
func (c *Clipboard) GetProviderName() string {
	return c.manager.GetProviderName()
}

// IsSSH returns true if running in an SSH session.
func (c *Clipboard) IsSSH() bool {
	return c.manager.IsSSH()
}

// Builder provides a fluent interface for creating a clipboard instance.
type Builder struct {
	providers     []service.Provider
	osc52Enabled  bool
	osc52Timeout  time.Duration
	nativeEnabled bool
}

// NewBuilder creates a new clipboard builder.
func NewBuilder() *Builder {
	return &Builder{
		osc52Enabled:  true,
		osc52Timeout:  5 * time.Second,
		nativeEnabled: true,
	}
}

// WithOSC52 enables or disables OSC 52 provider.
func (b *Builder) WithOSC52(enabled bool) *Builder {
	b.osc52Enabled = enabled
	return b
}

// WithOSC52Timeout sets the timeout for OSC 52 operations.
func (b *Builder) WithOSC52Timeout(timeout time.Duration) *Builder {
	b.osc52Timeout = timeout
	return b
}

// WithNative enables or disables native platform clipboard.
func (b *Builder) WithNative(enabled bool) *Builder {
	b.nativeEnabled = enabled
	return b
}

// WithProvider adds a custom provider to the clipboard.
func (b *Builder) WithProvider(provider service.Provider) *Builder {
	b.providers = append(b.providers, provider)
	return b
}

// Build creates the clipboard instance.
func (b *Builder) Build() (*Clipboard, error) {
	var providers []service.Provider

	// Add custom providers first (highest priority)
	providers = append(providers, b.providers...)

	// Add OSC 52 if enabled
	if b.osc52Enabled {
		osc52Provider := osc52.NewProvider(b.osc52Timeout)
		providers = append(providers, osc52Provider)
	}

	// Add native provider if enabled
	if b.nativeEnabled {
		nativeProvider := native.NewProvider()
		providers = append(providers, nativeProvider)
	}

	manager, err := application.NewClipboardManagerWithProviders(providers)
	if err != nil {
		return nil, err
	}

	return &Clipboard{
		manager: manager,
	}, nil
}

// Global clipboard instance (convenience functions).
var globalClipboard *Clipboard

// init initializes the global clipboard instance.
func init() {
	clipboard, err := New()
	if err != nil {
		// Don't panic, just leave it nil
		// Users can still create their own instances
		return
	}
	globalClipboard = clipboard
}

// Read reads text from the global clipboard instance.
func Read() (string, error) {
	if globalClipboard == nil {
		clipboard, err := New()
		if err != nil {
			return "", err
		}
		globalClipboard = clipboard
	}
	return globalClipboard.Read()
}

// Write writes text to the global clipboard instance.
func Write(text string) error {
	if globalClipboard == nil {
		clipboard, err := New()
		if err != nil {
			return err
		}
		globalClipboard = clipboard
	}
	return globalClipboard.Write(text)
}

// IsAvailable returns true if the global clipboard is available.
func IsAvailable() bool {
	if globalClipboard == nil {
		clipboard, err := New()
		if err != nil {
			return false
		}
		globalClipboard = clipboard
	}
	return globalClipboard.IsAvailable()
}

// GetProviderName returns the name of the active provider for the global clipboard.
func GetProviderName() string {
	if globalClipboard == nil {
		clipboard, err := New()
		if err != nil {
			return "none"
		}
		globalClipboard = clipboard
	}
	return globalClipboard.GetProviderName()
}
