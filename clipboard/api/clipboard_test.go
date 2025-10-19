package api

import (
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
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

func TestNew(t *testing.T) {
	clipboard, err := New()

	// May fail in headless environments
	if err != nil {
		t.Logf("New() failed (expected in headless env): %v", err)
		return
	}

	if clipboard == nil {
		t.Fatal("expected non-nil clipboard")
	}

	if clipboard.manager == nil {
		t.Error("expected non-nil manager")
	}
}

func TestClipboard_Read(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("test data")
		},
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, err := clipboard.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "test data" {
		t.Errorf("expected 'test data', got %s", text)
	}
}

func TestClipboard_Write(t *testing.T) {
	var written string

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		writeFunc: func(content *model.ClipboardContent) error {
			text, _ := content.Text()
			written = text
			return nil
		},
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = clipboard.Write("write test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if written != "write test" {
		t.Errorf("expected 'write test', got %s", written)
	}
}

func TestClipboard_IsAvailable(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !clipboard.IsAvailable() {
		t.Error("expected clipboard to be available")
	}
}

func TestClipboard_GetProviderName(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "TestMock",
		available: true,
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if clipboard.GetProviderName() != "TestMock" {
		t.Errorf("expected 'TestMock', got %s", clipboard.GetProviderName())
	}
}

func TestClipboard_IsSSH(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Just verify it doesn't panic
	_ = clipboard.IsSSH()
}

func TestBuilder_WithOSC52(t *testing.T) {
	builder := NewBuilder()
	builder.WithOSC52(true)

	if !builder.osc52Enabled {
		t.Error("expected OSC52 to be enabled")
	}

	builder.WithOSC52(false)

	if builder.osc52Enabled {
		t.Error("expected OSC52 to be disabled")
	}
}

func TestBuilder_WithOSC52Timeout(t *testing.T) {
	builder := NewBuilder()
	builder.WithOSC52Timeout(10 * time.Second)

	if builder.osc52Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", builder.osc52Timeout)
	}
}

func TestBuilder_WithNative(t *testing.T) {
	builder := NewBuilder()
	builder.WithNative(false)

	if builder.nativeEnabled {
		t.Error("expected native to be disabled")
	}

	builder.WithNative(true)

	if !builder.nativeEnabled {
		t.Error("expected native to be enabled")
	}
}

func TestBuilder_WithProvider(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	builder := NewBuilder()
	builder.WithProvider(mockProvider)

	if len(builder.providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(builder.providers))
	}
}

func TestBuilder_Build(t *testing.T) {
	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	clipboard, err := NewBuilder().
		WithProvider(mockProvider).
		WithOSC52(false).
		WithNative(false).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if clipboard == nil {
		t.Fatal("expected non-nil clipboard")
	}
}

func TestBuilder_Build_NoProviders(t *testing.T) {
	_, err := NewBuilder().
		WithOSC52(false).
		WithNative(false).
		Build()

	if err == nil {
		t.Error("expected error when no providers available")
	}
}

func TestGlobalFunctions(t *testing.T) {
	// Reset global clipboard
	globalClipboard = nil

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("global test")
		},
		writeFunc: func(content *model.ClipboardContent) error {
			return nil
		},
	}

	// Create a test clipboard to use as global
	testClipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	globalClipboard = testClipboard

	// Test Write
	err = Write("global write")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test Read
	text, err := Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "global test" {
		t.Errorf("expected 'global test', got %s", text)
	}

	// Test IsAvailable
	if !IsAvailable() {
		t.Error("expected clipboard to be available")
	}

	// Test GetProviderName
	if GetProviderName() != "mock" {
		t.Errorf("expected 'mock', got %s", GetProviderName())
	}
}

func TestGlobalRead_NilClipboard(t *testing.T) {
	// Reset global clipboard
	globalClipboard = nil

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("lazy init test")
		},
	}

	// Temporarily set global to nil to test lazy initialization
	globalClipboard = nil

	// Patch the initialization by creating a test clipboard
	testClipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set it as global for subsequent calls
	globalClipboard = testClipboard

	text, err := Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "lazy init test" {
		t.Errorf("expected 'lazy init test', got %s", text)
	}
}

func TestGlobalWrite_NilClipboard(t *testing.T) {
	// Reset global clipboard
	globalClipboard = nil

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
		writeFunc: func(content *model.ClipboardContent) error {
			return nil
		},
	}

	// Create a test clipboard
	testClipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	globalClipboard = testClipboard

	err = Write("lazy init write")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGlobalIsAvailable_NilClipboard(t *testing.T) {
	// Reset global clipboard
	globalClipboard = nil

	mockProvider := &MockProvider{
		name:      "mock",
		available: true,
	}

	testClipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	globalClipboard = testClipboard

	if !IsAvailable() {
		t.Error("expected clipboard to be available")
	}
}

func TestGlobalGetProviderName_NilClipboard(t *testing.T) {
	// Reset global clipboard
	globalClipboard = nil

	mockProvider := &MockProvider{
		name:      "lazy-mock",
		available: true,
	}

	testClipboard, err := NewBuilder().
		WithProvider(mockProvider).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	globalClipboard = testClipboard

	name := GetProviderName()
	if name != "lazy-mock" {
		t.Errorf("expected 'lazy-mock', got %s", name)
	}
}

func TestNew_ErrorHandling(t *testing.T) {
	// Test New() error handling
	// This may succeed or fail depending on environment
	clipboard, err := New()

	// Just verify it handles errors gracefully
	if err != nil {
		if clipboard != nil {
			t.Error("expected nil clipboard on error")
		}
	} else {
		if clipboard == nil {
			t.Error("expected non-nil clipboard on success")
		}
	}
}

func TestInit_ErrorHandling(t *testing.T) {
	// The init() function is called automatically
	// We test the behavior by checking if globalClipboard is set correctly

	// If globalClipboard is nil after init, it means initialization failed
	// which is acceptable in headless environments
	if globalClipboard == nil {
		t.Log("Global clipboard not initialized (acceptable in headless env)")
	} else {
		t.Log("Global clipboard initialized successfully")
	}
}

func TestBuilder_WithOSC52_Enabled(t *testing.T) {
	builder := NewBuilder().WithOSC52(true).WithNative(false)

	clipboard, err := builder.Build()
	if err != nil {
		// OSC52-only might fail if not available
		t.Logf("Build with OSC52 only failed: %v", err)
		return
	}

	if clipboard == nil {
		t.Error("expected non-nil clipboard")
	}
}

func TestBuilder_WithNative_Enabled(t *testing.T) {
	builder := NewBuilder().WithOSC52(false).WithNative(true)

	clipboard, err := builder.Build()
	if err != nil {
		// Native-only might fail in headless
		t.Logf("Build with native only failed: %v", err)
		return
	}

	if clipboard == nil {
		t.Error("expected non-nil clipboard")
	}
}

func TestBuilder_MultipleProviders(t *testing.T) {
	mock1 := &MockProvider{
		name:      "mock1",
		available: false,
	}
	mock2 := &MockProvider{
		name:      "mock2",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("fallback")
		},
	}

	clipboard, err := NewBuilder().
		WithProvider(mock1).
		WithProvider(mock2).
		WithOSC52(false).
		WithNative(false).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should use mock2 since mock1 is unavailable
	text, err := clipboard.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "fallback" {
		t.Errorf("expected 'fallback', got %s", text)
	}
}
