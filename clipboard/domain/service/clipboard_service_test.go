package service

import (
	"fmt"
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
)

// MockProvider is a mock implementation of Provider for testing
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

func TestNewClipboardService(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		wantError bool
	}{
		{
			"single provider",
			[]Provider{&MockProvider{name: "mock", available: true}},
			false,
		},
		{
			"multiple providers",
			[]Provider{
				&MockProvider{name: "mock1", available: true},
				&MockProvider{name: "mock2", available: true},
			},
			false,
		},
		{
			"no providers",
			[]Provider{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewClipboardService(tt.providers)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected non-nil service")
			}
		})
	}
}

func TestClipboardService_Read(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		wantError bool
		wantData  string
	}{
		{
			"available provider",
			[]Provider{&MockProvider{
				name:      "mock",
				available: true,
				readFunc: func() (*model.ClipboardContent, error) {
					return model.NewTextContent("test data")
				},
			}},
			false,
			"test data",
		},
		{
			"unavailable provider",
			[]Provider{&MockProvider{
				name:      "mock",
				available: false,
			}},
			true,
			"",
		},
		{
			"fallback to second provider",
			[]Provider{
				&MockProvider{name: "mock1", available: false},
				&MockProvider{
					name:      "mock2",
					available: true,
					readFunc: func() (*model.ClipboardContent, error) {
						return model.NewTextContent("fallback data")
					},
				},
			},
			false,
			"fallback data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService(tt.providers)
			if err != nil {
				t.Fatalf("unexpected error creating service: %v", err)
			}

			content, err := service.Read()

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			text, err := content.Text()
			if err != nil {
				t.Errorf("unexpected error getting text: %v", err)
				return
			}

			if text != tt.wantData {
				t.Errorf("expected data %s, got %s", tt.wantData, text)
			}
		})
	}
}

func TestClipboardService_Write(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		content   *model.ClipboardContent
		wantError bool
	}{
		{
			"available provider",
			[]Provider{&MockProvider{
				name:      "mock",
				available: true,
			}},
			mustNewTextContent("test"),
			false,
		},
		{
			"unavailable provider",
			[]Provider{&MockProvider{
				name:      "mock",
				available: false,
			}},
			mustNewTextContent("test"),
			true,
		},
		{
			"nil content",
			[]Provider{&MockProvider{
				name:      "mock",
				available: true,
			}},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService(tt.providers)
			if err != nil {
				t.Fatalf("unexpected error creating service: %v", err)
			}

			err = service.Write(tt.content)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestClipboardService_ReadText(t *testing.T) {
	provider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return model.NewTextContent("hello")
		},
	}

	service, err := NewClipboardService([]Provider{provider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, err := service.ReadText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "hello" {
		t.Errorf("expected 'hello', got %s", text)
	}
}

func TestClipboardService_WriteText(t *testing.T) {
	var written *model.ClipboardContent

	provider := &MockProvider{
		name:      "mock",
		available: true,
		writeFunc: func(content *model.ClipboardContent) error {
			written = content
			return nil
		},
	}

	service, err := NewClipboardService([]Provider{provider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = service.WriteText("test text")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if written == nil {
		t.Errorf("expected content to be written")
		return
	}

	text, err := written.Text()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if text != "test text" {
		t.Errorf("expected 'test text', got %s", text)
	}
}

func TestClipboardService_IsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		want      bool
	}{
		{
			"available provider",
			[]Provider{&MockProvider{name: "mock", available: true}},
			true,
		},
		{
			"unavailable provider",
			[]Provider{&MockProvider{name: "mock", available: false}},
			false,
		},
		{
			"mixed providers",
			[]Provider{
				&MockProvider{name: "mock1", available: false},
				&MockProvider{name: "mock2", available: true},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService(tt.providers)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := service.IsAvailable(); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClipboardService_GetAvailableProviderName(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		want      string
	}{
		{
			"single available provider",
			[]Provider{&MockProvider{name: "mock", available: true}},
			"mock",
		},
		{
			"no available provider",
			[]Provider{&MockProvider{name: "mock", available: false}},
			"none",
		},
		{
			"first available provider",
			[]Provider{
				&MockProvider{name: "mock1", available: false},
				&MockProvider{name: "mock2", available: true},
				&MockProvider{name: "mock3", available: true},
			},
			"mock2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService(tt.providers)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := service.GetAvailableProviderName(); got != tt.want {
				t.Errorf("GetAvailableProviderName() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestClipboardService_ProviderError(t *testing.T) {
	provider := &MockProvider{
		name:      "mock",
		available: true,
		readFunc: func() (*model.ClipboardContent, error) {
			return nil, fmt.Errorf("provider error")
		},
	}

	service, err := NewClipboardService([]Provider{provider})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.Read()
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestClipboardService_ReadText_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		provider *MockProvider
		wantErr  bool
	}{
		{
			name: "successful read",
			provider: &MockProvider{
				name:      "mock",
				available: true,
				readFunc: func() (*model.ClipboardContent, error) {
					return model.NewTextContent("success")
				},
			},
			wantErr: false,
		},
		{
			name: "read error",
			provider: &MockProvider{
				name:      "mock",
				available: true,
				readFunc: func() (*model.ClipboardContent, error) {
					return nil, fmt.Errorf("read failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService([]Provider{tt.provider})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, err = service.ReadText()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClipboardService_WriteText_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		provider *MockProvider
		wantErr  bool
	}{
		{
			name: "successful write",
			text: "test",
			provider: &MockProvider{
				name:      "mock",
				available: true,
				writeFunc: func(content *model.ClipboardContent) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "write error",
			text: "test",
			provider: &MockProvider{
				name:      "mock",
				available: true,
				writeFunc: func(content *model.ClipboardContent) error {
					return fmt.Errorf("write failed")
				},
			},
			wantErr: true,
		},
		{
			name: "empty text (invalid)",
			text: "",
			provider: &MockProvider{
				name:      "mock",
				available: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewClipboardService([]Provider{tt.provider})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			err = service.WriteText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mustNewTextContent(text string) *model.ClipboardContent {
	content, err := model.NewTextContent(text)
	if err != nil {
		panic(err)
	}
	return content
}
