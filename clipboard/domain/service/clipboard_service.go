package service

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
)

// ClipboardService provides domain logic for clipboard operations
type ClipboardService struct {
	providers []Provider
}

// NewClipboardService creates a new clipboard service with a prioritized list of providers
func NewClipboardService(providers []Provider) (*ClipboardService, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("at least one provider must be specified")
	}

	return &ClipboardService{
		providers: providers,
	}, nil
}

// Read reads content from the first available provider
func (s *ClipboardService) Read() (*model.ClipboardContent, error) {
	provider := s.getAvailableProvider()
	if provider == nil {
		return nil, fmt.Errorf("no clipboard provider available")
	}

	return provider.Read()
}

// Write writes content using the first available provider
func (s *ClipboardService) Write(content *model.ClipboardContent) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	provider := s.getAvailableProvider()
	if provider == nil {
		return fmt.Errorf("no clipboard provider available")
	}

	return provider.Write(content)
}

// ReadText reads text content from the clipboard
func (s *ClipboardService) ReadText() (string, error) {
	content, err := s.Read()
	if err != nil {
		return "", err
	}

	return content.Text()
}

// WriteText writes text content to the clipboard
func (s *ClipboardService) WriteText(text string) error {
	content, err := model.NewTextContent(text)
	if err != nil {
		return err
	}

	return s.Write(content)
}

// IsAvailable returns true if any provider is available
func (s *ClipboardService) IsAvailable() bool {
	return s.getAvailableProvider() != nil
}

// GetAvailableProviderName returns the name of the first available provider
func (s *ClipboardService) GetAvailableProviderName() string {
	provider := s.getAvailableProvider()
	if provider == nil {
		return "none"
	}
	return provider.Name()
}

// getAvailableProvider returns the first available provider
func (s *ClipboardService) getAvailableProvider() Provider {
	for _, provider := range s.providers {
		if provider.IsAvailable() {
			return provider
		}
	}
	return nil
}
