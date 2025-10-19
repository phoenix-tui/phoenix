package application

import (
	"os"
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
	"github.com/phoenix-tui/phoenix/clipboard/domain/service"
	"github.com/phoenix-tui/phoenix/clipboard/infrastructure/platform"
)

// MockProvider for testing
type MockProvider struct {
	name      string
	available bool
	readFunc  func() (*model.ClipboardContent, error)
	writeFunc func(content *model.ClipboardContent) error
}

func (m *MockProvider) Read() (*model.ClipboardContent, error) {
	if m.readFunc != nil {
		return m.readFunc()
	}
	return model.NewTextContent("mock data")
}

func (m *MockProvider) Write(content *model.ClipboardContent) error {
	if m.writeFunc != nil {
		return m.writeFunc(content)
	}
	return nil
}

func (m *MockProvider) IsAvailable() bool {
	return m.available
}

func (m *MockProvider) Name() string {
	return m.name
}

func TestNewClipboardManager(t *testing.T) {
	// This will use auto-detected providers
	manager, err := NewClipboardManager()

	// We expect this to succeed on most platforms
	// but it might fail in headless environments
	if err != nil {
		t.Logf("NewClipboardManager failed (expected in headless env): %v", err)
	}

	if manager != nil {
		if manager.service == nil {
			t.Error("expected non-nil service")
		}
		if manager.detector == nil {
			t.Error("expected non-nil detector")
		}
	}
}

func TestNewClipboardManagerWithProviders(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if manager == nil {
		t.Fatal("expected non-nil manager")
	}

	if manager.service == nil {
		t.Error("expected non-nil service")
	}
}

func TestClipboardManager_Write(t *testing.T) {
	var written string

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		writeFunc: func(content *model.ClipboardContent) error {
			text, err := content.Text()
			if err != nil {
				return err
			}
			written = text
			return nil
		},
	}

	manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = manager.Write("test text")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if written != "test text" {
		t.Errorf("expected 'test text', got %s", written)
	}
}

func TestClipboardManager_Read(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("read text")
		},
	}

	manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, err := manager.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "read text" {
		t.Errorf("expected 'read text', got %s", text)
	}
}

func TestClipboardManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		available bool
		want      bool
	}{
		{"available", true, true},
		{"unavailable", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &MockProvider{
				name:      "mock",
				available: tt.available,
			}

			manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := manager.IsAvailable(); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClipboardManager_GetProviderName(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "TestProvider",
		available: true,
	}

	manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := manager.GetProviderName(); got != "TestProvider" {
		t.Errorf("GetProviderName() = %s, want TestProvider", got)
	}
}

func TestClipboardManager_IsSSH(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	manager, err := NewClipboardManagerWithProviders([]service.Provider{mockProvider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Just verify it doesn't panic
	_ = manager.IsSSH()
}

func TestNewClipboardManagerWithProviders_ErrorHandling(t *testing.T) {
	// Test with empty provider list (should fail in service creation)
	_, err := NewClipboardManagerWithProviders([]service.Provider{})
	if err == nil {
		t.Error("expected error with empty provider list")
	}
}

func TestBuildProviderChain_SSHSession(t *testing.T) {
	// Save original env
	origSSHTTY := os.Getenv("SSH_TTY")
	defer os.Setenv("SSH_TTY", origSSHTTY)

	// Simulate SSH session
	os.Setenv("SSH_TTY", "/dev/pts/0")

	detector := platform.NewDetector()
	providers := buildProviderChain(detector)

	// Should have providers (OSC52 first in SSH, then native, then OSC52 fallback)
	if len(providers) == 0 {
		t.Error("expected at least one provider")
	}

	// Verify we have multiple providers for redundancy
	if len(providers) < 2 {
		t.Logf("Warning: expected multiple providers for redundancy, got %d", len(providers))
	}
}

func TestBuildProviderChain_NonSSH(t *testing.T) {
	// Save original env
	origSSHTTY := os.Getenv("SSH_TTY")
	origSSHClient := os.Getenv("SSH_CLIENT")
	origSSHConnection := os.Getenv("SSH_CONNECTION")
	defer func() {
		os.Setenv("SSH_TTY", origSSHTTY)
		os.Setenv("SSH_CLIENT", origSSHClient)
		os.Setenv("SSH_CONNECTION", origSSHConnection)
	}()

	// Clear SSH environment variables
	os.Setenv("SSH_TTY", "")
	os.Setenv("SSH_CLIENT", "")
	os.Setenv("SSH_CONNECTION", "")

	detector := platform.NewDetector()
	providers := buildProviderChain(detector)

	// Should have providers (native first, then OSC52 fallback)
	if len(providers) == 0 {
		t.Error("expected at least one provider")
	}

	// Verify we have multiple providers for fallback
	if len(providers) < 2 {
		t.Logf("Warning: expected multiple providers for fallback, got %d", len(providers))
	}
}

func TestNewClipboardManager_ErrorPropagation(t *testing.T) {
	// NewClipboardManager should succeed on most platforms
	// but may fail in truly headless environments
	manager, err := NewClipboardManager()

	if err != nil {
		// Error is acceptable in headless environments
		t.Logf("NewClipboardManager failed (acceptable in headless): %v", err)
		return
	}

	if manager == nil {
		t.Error("expected non-nil manager when no error")
	}
}
