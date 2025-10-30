// Package application implements clipboard use cases and coordinates between domain services and providers.
package application

import (
	"fmt"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/clipboard/internal/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/native"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/osc52"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/platform"
)

// ClipboardManager is the application service that manages clipboard operations.
// It coordinates between domain services and infrastructure providers.
type ClipboardManager struct {
	service        *service2.ClipboardService
	detector       *platform.Detector
	richTextCodec  *service2.RichTextCodec
	history        *service2.ClipboardHistory
	historyEnabled bool
}

// NewClipboardManager creates a new clipboard manager with auto-detected providers.
func NewClipboardManager() (*ClipboardManager, error) {
	detector := platform.NewDetector()

	// Build provider chain based on platform and environment
	providers := buildProviderChain(detector)

	// Create domain service with providers
	clipboardService, err := service2.NewClipboardService(providers)
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		service:        clipboardService,
		detector:       detector,
		richTextCodec:  service2.NewRichTextCodec(),
		history:        nil, // History disabled by default
		historyEnabled: false,
	}, nil
}

// NewClipboardManagerWithProviders creates a clipboard manager with custom providers.
func NewClipboardManagerWithProviders(providers []service2.Provider) (*ClipboardManager, error) {
	clipboardService, err := service2.NewClipboardService(providers)
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		service:        clipboardService,
		detector:       platform.NewDetector(),
		richTextCodec:  service2.NewRichTextCodec(),
		history:        nil, // History disabled by default
		historyEnabled: false,
	}, nil
}

// Read reads text from the clipboard.
func (m *ClipboardManager) Read() (string, error) {
	return m.service.ReadText()
}

// Write writes text to the clipboard.
// If history is enabled, the content is added to history.
func (m *ClipboardManager) Write(text string) error {
	// Write to clipboard
	if err := m.service.WriteText(text); err != nil {
		return err
	}

	// Add to history if enabled
	if m.historyEnabled && m.history != nil {
		_ = m.history.Add([]byte(text), value.MIMETypePlainText)
	}

	return nil
}

// IsAvailable returns true if clipboard is available.
func (m *ClipboardManager) IsAvailable() bool {
	return m.service.IsAvailable()
}

// GetProviderName returns the name of the active provider.
func (m *ClipboardManager) GetProviderName() string {
	return m.service.GetAvailableProviderName()
}

// IsSSH returns true if running in an SSH session.
func (m *ClipboardManager) IsSSH() bool {
	return m.detector.IsSSH()
}

// buildProviderChain builds a prioritized list of clipboard providers.
func buildProviderChain(detector *platform.Detector) []service2.Provider {
	var providers []service2.Provider

	// Priority 1: OSC 52 if in SSH session
	if detector.ShouldUseOSC52() {
		osc52Provider := osc52.NewProvider(5 * time.Second)
		providers = append(providers, osc52Provider)
	}

	// Priority 2: Native platform clipboard
	nativeProvider := native.NewProvider()
	providers = append(providers, nativeProvider)

	// Priority 3: OSC 52 as fallback (even if not in SSH)
	// Some local terminals support OSC 52
	if !detector.ShouldUseOSC52() {
		osc52Provider := osc52.NewProvider(2 * time.Second)
		providers = append(providers, osc52Provider)
	}

	return providers
}

// ReadImage reads image data from the clipboard.
// Returns the image bytes and the detected MIME type.
// Note: Image clipboard support is currently limited to native providers.
func (m *ClipboardManager) ReadImage() ([]byte, string, error) {
	// For now, this is a placeholder that returns an error
	// Full implementation requires provider-specific image support
	return nil, "", fmt.Errorf("image clipboard operations not yet implemented in providers")
}

// WriteImage writes image data to the clipboard.
// The data should be in PNG, JPEG, or GIF format.
// Note: Image clipboard support is currently limited to native providers.
func (m *ClipboardManager) WriteImage(data []byte, mimeType string) error {
	// For now, this is a placeholder that returns an error
	// Full implementation requires provider-specific image support
	return fmt.Errorf("image clipboard operations not yet implemented in providers")
}

// ReadHTML reads HTML content from the clipboard.
// Returns the HTML string.
// Note: This is a simplified implementation that reads plain text and treats it as HTML.
// Full HTML clipboard support requires provider-specific implementation.
func (m *ClipboardManager) ReadHTML() (string, error) {
	// For now, read as plain text
	// Future: providers should support reading HTML directly
	text, err := m.service.ReadText()
	if err != nil {
		return "", fmt.Errorf("failed to read HTML from clipboard: %w", err)
	}
	return text, nil
}

// WriteHTML writes HTML content to the clipboard.
// The html parameter should contain valid HTML markup.
// Note: This is a simplified implementation that writes HTML as plain text.
// Full HTML clipboard support requires provider-specific implementation.
func (m *ClipboardManager) WriteHTML(html string) error {
	// For now, write as plain text
	// Future: providers should support writing HTML directly
	if err := m.service.WriteText(html); err != nil {
		return fmt.Errorf("failed to write HTML to clipboard: %w", err)
	}

	// Add to history if enabled
	if m.historyEnabled && m.history != nil {
		_ = m.history.Add([]byte(html), value.MIMETypeHTML)
	}

	return nil
}

// ReadRTF reads RTF content from the clipboard.
// Returns the RTF string.
// Note: This is a simplified implementation that reads plain text and treats it as RTF.
// Full RTF clipboard support requires provider-specific implementation.
func (m *ClipboardManager) ReadRTF() (string, error) {
	// For now, read as plain text
	// Future: providers should support reading RTF directly
	text, err := m.service.ReadText()
	if err != nil {
		return "", fmt.Errorf("failed to read RTF from clipboard: %w", err)
	}
	return text, nil
}

// WriteRTF writes RTF content to the clipboard.
// The rtf parameter should contain valid RTF markup.
// Note: This is a simplified implementation that writes RTF as plain text.
// Full RTF clipboard support requires provider-specific implementation.
func (m *ClipboardManager) WriteRTF(rtf string) error {
	// For now, write as plain text
	// Future: providers should support writing RTF directly
	if err := m.service.WriteText(rtf); err != nil {
		return fmt.Errorf("failed to write RTF to clipboard: %w", err)
	}

	// Add to history if enabled
	if m.historyEnabled && m.history != nil {
		_ = m.history.Add([]byte(rtf), value.MIMETypeRTF)
	}

	return nil
}

// WriteStyledText writes text with styling to the clipboard in HTML format.
// This is a convenience method that encodes text with styles and writes it.
func (m *ClipboardManager) WriteStyledText(text string, styles value.TextStyles) error {
	html, err := m.richTextCodec.EncodeHTML(text, styles)
	if err != nil {
		return fmt.Errorf("failed to encode styled text: %w", err)
	}
	return m.WriteHTML(html)
}

// ReadStyledText reads HTML from the clipboard and returns plain text with detected styles.
// This is a convenience method that reads HTML and decodes it.
func (m *ClipboardManager) ReadStyledText() (string, value.TextStyles, error) {
	html, err := m.ReadHTML()
	if err != nil {
		return "", value.NewTextStyles(), fmt.Errorf("failed to read HTML: %w", err)
	}

	text, styles, err := m.richTextCodec.DecodeHTML(html)
	if err != nil {
		return "", value.NewTextStyles(), fmt.Errorf("failed to decode HTML: %w", err)
	}

	return text, styles, nil
}

// ConvertHTMLToRTF converts HTML content to RTF format using the codec.
func (m *ClipboardManager) ConvertHTMLToRTF(html string) (string, error) {
	return m.richTextCodec.HTMLToRTF(html)
}

// ConvertRTFToHTML converts RTF content to HTML format using the codec.
func (m *ClipboardManager) ConvertRTFToHTML(rtf string) (string, error) {
	return m.richTextCodec.RTFToHTML(rtf)
}

// EnableHistory enables clipboard history tracking with the given limits.
// maxSize: maximum number of entries (0 = unlimited)
// maxAge: maximum age of entries (0 = no expiration)
func (m *ClipboardManager) EnableHistory(maxSize int, maxAge time.Duration) {
	m.history = service2.NewClipboardHistory(maxSize, maxAge)
	m.historyEnabled = true
}

// DisableHistory disables clipboard history tracking.
func (m *ClipboardManager) DisableHistory() {
	m.historyEnabled = false
	m.history = nil
}

// IsHistoryEnabled returns true if clipboard history tracking is active.
func (m *ClipboardManager) IsHistoryEnabled() bool {
	return m.historyEnabled
}

// GetHistory returns all clipboard history entries.
// Returns nil if history is not enabled.
func (m *ClipboardManager) GetHistory() []*model.HistoryEntry {
	if !m.historyEnabled || m.history == nil {
		return nil
	}
	return m.history.GetAll()
}

// GetHistoryEntry returns a specific history entry by ID.
// Returns an error if history is not enabled or entry not found.
func (m *ClipboardManager) GetHistoryEntry(id string) (*model.HistoryEntry, error) {
	if !m.historyEnabled || m.history == nil {
		return nil, fmt.Errorf("clipboard history is not enabled")
	}
	return m.history.Get(id)
}

// GetRecentHistory returns the N most recent history entries.
// Returns nil if history is not enabled.
func (m *ClipboardManager) GetRecentHistory(count int) []*model.HistoryEntry {
	if !m.historyEnabled || m.history == nil {
		return nil
	}
	return m.history.GetRecent(count)
}

// ClearHistory removes all history entries.
// Does nothing if history is not enabled.
func (m *ClipboardManager) ClearHistory() {
	if m.historyEnabled && m.history != nil {
		m.history.Clear()
	}
}

// GetHistorySize returns the number of entries in history.
// Returns 0 if history is not enabled.
func (m *ClipboardManager) GetHistorySize() int {
	if !m.historyEnabled || m.history == nil {
		return 0
	}
	return m.history.Size()
}

// GetHistoryTotalSize returns the total memory usage of history in bytes.
// Returns 0 if history is not enabled.
func (m *ClipboardManager) GetHistoryTotalSize() int {
	if !m.historyEnabled || m.history == nil {
		return 0
	}
	return m.history.TotalSize()
}

// RemoveExpiredHistory removes expired entries from history.
// Returns the number of entries removed.
// Returns 0 if history is not enabled.
func (m *ClipboardManager) RemoveExpiredHistory() int {
	if !m.historyEnabled || m.history == nil {
		return 0
	}
	return m.history.RemoveExpired()
}
