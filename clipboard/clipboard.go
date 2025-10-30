// Package clipboard provides cross-platform clipboard operations with native and OSC 52 support.
package clipboard

import (
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/application"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/native"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/osc52"
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

// ReadImage reads image data from the clipboard.
// Returns the image bytes, MIME type, and any error.
// Note: Image clipboard support requires native provider implementation.
func (c *Clipboard) ReadImage() ([]byte, string, error) {
	return c.manager.ReadImage()
}

// WriteImage writes image data to the clipboard.
// The data parameter should contain the image bytes in PNG, JPEG, or GIF format.
// The mimeType parameter should be the MIME type (e.g., "image/png", "image/jpeg").
// Note: Image clipboard support requires native provider implementation.
func (c *Clipboard) WriteImage(data []byte, mimeType string) error {
	return c.manager.WriteImage(data, mimeType)
}

// ReadImagePNG reads PNG image data from the clipboard (convenience method).
func (c *Clipboard) ReadImagePNG() ([]byte, error) {
	data, mimeType, err := c.ReadImage()
	if err != nil {
		return nil, err
	}
	if mimeType != "image/png" {
		// Convert to PNG if needed
		codec := service.NewImageCodec()
		img, _, err := codec.Decode(data)
		if err != nil {
			return nil, err
		}
		return codec.EncodePNG(img)
	}
	return data, nil
}

// WriteImagePNG writes PNG image data to the clipboard (convenience method).
func (c *Clipboard) WriteImagePNG(data []byte) error {
	return c.WriteImage(data, "image/png")
}

// ReadImageJPEG reads JPEG image data from the clipboard (convenience method).
func (c *Clipboard) ReadImageJPEG() ([]byte, error) {
	data, mimeType, err := c.ReadImage()
	if err != nil {
		return nil, err
	}
	if mimeType != "image/jpeg" {
		// Convert to JPEG if needed
		codec := service.NewImageCodec()
		img, _, err := codec.Decode(data)
		if err != nil {
			return nil, err
		}
		return codec.EncodeJPEG(img, 90)
	}
	return data, nil
}

// WriteImageJPEG writes JPEG image data to the clipboard (convenience method).
func (c *Clipboard) WriteImageJPEG(data []byte) error {
	return c.WriteImage(data, "image/jpeg")
}

// ReadHTML reads HTML content from the clipboard.
// Returns the HTML string.
func (c *Clipboard) ReadHTML() (string, error) {
	return c.manager.ReadHTML()
}

// WriteHTML writes HTML content to the clipboard.
// The html parameter should contain valid HTML markup.
func (c *Clipboard) WriteHTML(html string) error {
	return c.manager.WriteHTML(html)
}

// ReadHTMLAsPlainText reads HTML from clipboard and strips all tags, returning plain text.
// This is a convenience method for getting text content from HTML.
func (c *Clipboard) ReadHTMLAsPlainText() (string, error) {
	html, err := c.ReadHTML()
	if err != nil {
		return "", err
	}

	codec := service.NewRichTextCodec()
	return codec.StripHTMLTags(html)
}

// ReadRTF reads RTF content from the clipboard.
// Returns the RTF string.
func (c *Clipboard) ReadRTF() (string, error) {
	return c.manager.ReadRTF()
}

// WriteRTF writes RTF content to the clipboard.
// The rtf parameter should contain valid RTF markup.
func (c *Clipboard) WriteRTF(rtf string) error {
	return c.manager.WriteRTF(rtf)
}

// ReadRTFAsPlainText reads RTF from clipboard and strips all formatting, returning plain text.
// This is a convenience method for getting text content from RTF.
func (c *Clipboard) ReadRTFAsPlainText() (string, error) {
	rtf, err := c.ReadRTF()
	if err != nil {
		return "", err
	}

	codec := service.NewRichTextCodec()
	return codec.StripRTFFormatting(rtf)
}

// ConvertHTMLToRTF converts HTML content to RTF format.
// This is a convenience method that uses the rich text codec.
func (c *Clipboard) ConvertHTMLToRTF(html string) (string, error) {
	return c.manager.ConvertHTMLToRTF(html)
}

// ConvertRTFToHTML converts RTF content to HTML format.
// This is a convenience method that uses the rich text codec.
func (c *Clipboard) ConvertRTFToHTML(rtf string) (string, error) {
	return c.manager.ConvertRTFToHTML(rtf)
}

// HistoryEntry represents a single clipboard history item in the public API.
type HistoryEntry struct {
	ID        string
	Content   []byte
	MIMEType  string
	Timestamp time.Time
	Size      int
}

// EnableHistory enables clipboard history tracking with the given limits.
// maxSize: maximum number of entries (0 = unlimited)
// maxAge: maximum age of entries (0 = no expiration)
//
// Example:
//
//	// Track last 100 entries for 24 hours
//	c.EnableHistory(100, 24*time.Hour)
//
//	// Unlimited entries, no expiration
//	c.EnableHistory(0, 0)
func (c *Clipboard) EnableHistory(maxSize int, maxAge time.Duration) {
	c.manager.EnableHistory(maxSize, maxAge)
}

// DisableHistory disables clipboard history tracking and clears existing history.
func (c *Clipboard) DisableHistory() {
	c.manager.DisableHistory()
}

// IsHistoryEnabled returns true if clipboard history tracking is active.
func (c *Clipboard) IsHistoryEnabled() bool {
	return c.manager.IsHistoryEnabled()
}

// GetHistory returns all clipboard history entries sorted by timestamp (newest first).
// Returns empty slice if history is not enabled.
func (c *Clipboard) GetHistory() []HistoryEntry {
	entries := c.manager.GetHistory()
	if entries == nil {
		return []HistoryEntry{}
	}

	// Convert internal entries to public API
	result := make([]HistoryEntry, len(entries))
	for i, entry := range entries {
		result[i] = HistoryEntry{
			ID:        entry.ID(),
			Content:   entry.Content(),
			MIMEType:  entry.MIMEType().String(),
			Timestamp: entry.Timestamp(),
			Size:      entry.Size(),
		}
	}
	return result
}

// GetHistoryEntry returns a specific history entry by ID.
// Returns error if history is not enabled or entry not found.
func (c *Clipboard) GetHistoryEntry(id string) (HistoryEntry, error) {
	entry, err := c.manager.GetHistoryEntry(id)
	if err != nil {
		return HistoryEntry{}, err
	}

	return HistoryEntry{
		ID:        entry.ID(),
		Content:   entry.Content(),
		MIMEType:  entry.MIMEType().String(),
		Timestamp: entry.Timestamp(),
		Size:      entry.Size(),
	}, nil
}

// GetRecentHistory returns the N most recent history entries.
// If count is 0 or negative, returns all entries.
// Returns empty slice if history is not enabled.
func (c *Clipboard) GetRecentHistory(count int) []HistoryEntry {
	entries := c.manager.GetRecentHistory(count)
	if entries == nil {
		return []HistoryEntry{}
	}

	// Convert internal entries to public API
	result := make([]HistoryEntry, len(entries))
	for i, entry := range entries {
		result[i] = HistoryEntry{
			ID:        entry.ID(),
			Content:   entry.Content(),
			MIMEType:  entry.MIMEType().String(),
			Timestamp: entry.Timestamp(),
			Size:      entry.Size(),
		}
	}
	return result
}

// ClearHistory removes all history entries.
// Does nothing if history is not enabled.
func (c *Clipboard) ClearHistory() {
	c.manager.ClearHistory()
}

// GetHistorySize returns the number of entries in history.
// Returns 0 if history is not enabled.
func (c *Clipboard) GetHistorySize() int {
	return c.manager.GetHistorySize()
}

// GetHistoryTotalSize returns the total memory usage of history in bytes.
// Returns 0 if history is not enabled.
func (c *Clipboard) GetHistoryTotalSize() int {
	return c.manager.GetHistoryTotalSize()
}

// RemoveExpiredHistory removes expired entries from history.
// Returns the number of entries removed.
// Returns 0 if history is not enabled.
//
// This method should be called periodically to clean up old entries.
// Consider calling it in a background goroutine:
//
//	// ticker := time.NewTicker(1 * time.Hour)
//	// go func() {
//	//     for range ticker.C {
//	//         removed := c.RemoveExpiredHistory()
//	//         log.Printf("Removed %d expired entries", removed)
//	//     }
//	// }()
func (c *Clipboard) RemoveExpiredHistory() int {
	return c.manager.RemoveExpiredHistory()
}

// RestoreFromHistory restores a history entry to the clipboard by writing it back.
// This is a convenience method that combines GetHistoryEntry and Write/WriteImage/WriteHTML/WriteRTF.
func (c *Clipboard) RestoreFromHistory(id string) error {
	entry, err := c.manager.GetHistoryEntry(id)
	if err != nil {
		return err
	}

	// Determine how to write based on MIME type
	mimeType := entry.MIMEType()
	if mimeType.IsText() {
		text, err := entry.Text()
		if err != nil {
			return err
		}
		// Choose write method based on specific MIME type
		switch mimeType.String() {
		case "text/html":
			return c.WriteHTML(text)
		case "text/rtf":
			return c.WriteRTF(text)
		default:
			return c.Write(text)
		}
	} else if mimeType.IsImage() {
		return c.WriteImage(entry.Content(), mimeType.String())
	}

	// Fallback to plain text
	return c.Write(string(entry.Content()))
}
