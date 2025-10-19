package application

import (
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/native"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/osc52"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/platform"
)

// ClipboardManager is the application service that manages clipboard operations
// It coordinates between domain services and infrastructure providers
type ClipboardManager struct {
	service  *service.ClipboardService
	detector *platform.Detector
}

// NewClipboardManager creates a new clipboard manager with auto-detected providers
func NewClipboardManager() (*ClipboardManager, error) {
	detector := platform.NewDetector()

	// Build provider chain based on platform and environment
	providers := buildProviderChain(detector)

	// Create domain service with providers
	clipboardService, err := service.NewClipboardService(providers)
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		service:  clipboardService,
		detector: detector,
	}, nil
}

// NewClipboardManagerWithProviders creates a clipboard manager with custom providers
func NewClipboardManagerWithProviders(providers []service.Provider) (*ClipboardManager, error) {
	clipboardService, err := service.NewClipboardService(providers)
	if err != nil {
		return nil, err
	}

	return &ClipboardManager{
		service:  clipboardService,
		detector: platform.NewDetector(),
	}, nil
}

// Read reads text from the clipboard
func (m *ClipboardManager) Read() (string, error) {
	return m.service.ReadText()
}

// Write writes text to the clipboard
func (m *ClipboardManager) Write(text string) error {
	return m.service.WriteText(text)
}

// IsAvailable returns true if clipboard is available
func (m *ClipboardManager) IsAvailable() bool {
	return m.service.IsAvailable()
}

// GetProviderName returns the name of the active provider
func (m *ClipboardManager) GetProviderName() string {
	return m.service.GetAvailableProviderName()
}

// IsSSH returns true if running in an SSH session
func (m *ClipboardManager) IsSSH() bool {
	return m.detector.IsSSH()
}

// buildProviderChain builds a prioritized list of clipboard providers
func buildProviderChain(detector *platform.Detector) []service.Provider {
	var providers []service.Provider

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
