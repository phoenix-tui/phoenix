// Package application implements clipboard use cases and coordinates between domain services and providers.
package application

import (
	"fmt"
	"time"

	service2 "github.com/phoenix-tui/phoenix/clipboard/internal/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/native"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/osc52"
	"github.com/phoenix-tui/phoenix/clipboard/internal/infrastructure/platform"
)

// ClipboardManager is the application service that manages clipboard operations.
// It coordinates between domain services and infrastructure providers.
type ClipboardManager struct {
	service       *service2.ClipboardService
	detector      *platform.Detector
	richTextCodec *service2.RichTextCodec
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
		service:       clipboardService,
		detector:      detector,
		richTextCodec: service2.NewRichTextCodec(),
	}, nil
}

// NewClipboardManagerWithProviders creates a clipboard manager with custom providers.
func NewClipboardManagerWithProviders(providers []service2.Provider) (*ClipboardManager, error) {
	clipboardService, err := service2.NewClipboardService(providers)
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		service:       clipboardService,
		detector:      platform.NewDetector(),
		richTextCodec: service2.NewRichTextCodec(),
	}, nil
}

// Read reads text from the clipboard.
func (m *ClipboardManager) Read() (string, error) {
	return m.service.ReadText()
}

// Write writes text to the clipboard.
func (m *ClipboardManager) Write(text string) error {
	return m.service.WriteText(text)
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
