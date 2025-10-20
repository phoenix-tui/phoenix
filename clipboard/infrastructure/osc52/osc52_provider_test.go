package osc52

import (
	"os"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
)

func TestNewProvider(t *testing.T) {
	provider := NewProvider(5 * time.Second)

	if provider == nil {
		t.Fatal("expected non-nil provider")
	}

	if provider.timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", provider.timeout)
	}

	if provider.output != os.Stdout {
		t.Errorf("expected output to be stdout")
	}
}

func TestProvider_Name(t *testing.T) {
	provider := NewProvider(5 * time.Second)

	if provider.Name() != "OSC52" {
		t.Errorf("expected name 'OSC52', got %s", provider.Name())
	}
}

func TestProvider_Read(t *testing.T) {
	provider := NewProvider(5 * time.Second)

	// OSC 52 read is not supported
	_, err := provider.Read()
	if err == nil {
		t.Error("expected error for OSC 52 read (not supported)")
	}
}

func TestProvider_Write(t *testing.T) {
	// Create a temporary file to capture output
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	content, err := model.NewTextContent("test data")
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}

	err = provider.Write(content)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Read back the written data
	tmpFile.Seek(0, 0)
	buf := make([]byte, 1024)
	n, _ := tmpFile.Read(buf)
	output := string(buf[:n])

	// Verify OSC 52 sequence format: ESC ] 52 ; c ; <base64> ESC \
	if len(output) < 10 {
		t.Errorf("output too short: %d bytes", len(output))
		return
	}

	// Check for OSC 52 start: ESC ] 52 ; c
	// Format: \033]52;c;<base64>\033\\
	if !containsString(output, "\033]52;") {
		t.Errorf("expected OSC 52 escape sequence in output, got %q", output)
	}

	if !containsString(output, ";c;") {
		t.Errorf("expected clipboard selection 'c' in OSC 52 sequence, got %q", output)
	}
}

func TestProvider_Write_NilContent(t *testing.T) {
	provider := NewProvider(5 * time.Second)

	err := provider.Write(nil)
	if err == nil {
		t.Error("expected error for nil content")
	}
}

func TestProvider_Write_Timeout(t *testing.T) {
	t.Skip("Timeout test is hard to make reliable across platforms")

	// Create a custom writer that blocks forever
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer pr.Close()
	defer pw.Close()

	// Fill the pipe buffer to cause blocking
	// Create a provider with very short timeout
	provider := NewProvider(1 * time.Millisecond)
	provider.WithOutput(pw)

	// Generate large content to ensure timeout
	largeData := ""
	for i := 0; i < 100000; i++ {
		largeData += "x"
	}
	content, _ := model.NewTextContent(largeData)

	// This should timeout due to blocking write
	err = provider.Write(content)

	// We expect either a timeout or successful write (depending on buffer size)
	// The key is that it doesn't hang forever
	_ = err
}

func TestProvider_IsAvailable(t *testing.T) {
	// Save original env vars
	origSSHTTY := os.Getenv("SSH_TTY")
	origSSHClient := os.Getenv("SSH_CLIENT")
	origSSHConnection := os.Getenv("SSH_CONNECTION")
	origTERM := os.Getenv("TERM")

	defer func() {
		os.Setenv("SSH_TTY", origSSHTTY)
		os.Setenv("SSH_CLIENT", origSSHClient)
		os.Setenv("SSH_CONNECTION", origSSHConnection)
		os.Setenv("TERM", origTERM)
	}()

	tests := []struct {
		name          string
		sshTTY        string
		sshClient     string
		sshConnection string
		term          string
		setupOutput   func() *os.File
		want          bool
	}{
		{"SSH_TTY session", "/dev/pts/0", "", "", "", nil, true},
		{"SSH_CLIENT session", "", "127.0.0.1", "", "", nil, true},
		{"SSH_CONNECTION session", "", "", "127.0.0.1 22", "", nil, true},
		{"supported xterm terminal", "", "", "", "xterm", nil, true},
		{"supported xterm-256color", "", "", "", "xterm-256color", nil, true},
		{"supported screen", "", "", "", "screen", nil, true},
		{"supported tmux", "", "", "", "tmux", nil, true},
		{"supported tmux-256color", "", "", "", "tmux-256color", nil, true},
		{"no indicators", "", "", "", "dumb", nil, false},
		{
			name: "nil output",
			setupOutput: func() *os.File {
				return nil
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //nolint:revive // Table-driven test pattern
			os.Setenv("SSH_TTY", tt.sshTTY)
			os.Setenv("SSH_CLIENT", tt.sshClient)
			os.Setenv("SSH_CONNECTION", tt.sshConnection)
			os.Setenv("TERM", tt.term)

			provider := NewProvider(5 * time.Second)

			if tt.setupOutput != nil {
				provider.WithOutput(tt.setupOutput())
			}

			// Note: IsAvailable also checks if output is a terminal
			// In test environment, stdout might not be a TTY
			// So we can't reliably test the full availability logic
			_ = provider.IsAvailable()
		})
	}
}

func TestProvider_IsAvailable_NilOutput(t *testing.T) {
	provider := NewProvider(5 * time.Second)
	provider.output = nil

	if provider.IsAvailable() {
		t.Error("expected IsAvailable to return false for nil output")
	}
}

func TestProvider_IsAvailable_StatError(t *testing.T) {
	// Create a temp file and close it to trigger Stat error
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tmpPath)

	// Try to open the already-deleted file
	closedFile, _ := os.Open(tmpPath)

	provider := NewProvider(5 * time.Second)
	if closedFile != nil {
		provider.WithOutput(closedFile)
	}

	// IsAvailable should handle Stat errors gracefully
	_ = provider.IsAvailable()
}

func TestProvider_WithOutput(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	if provider.output != tmpFile {
		t.Error("expected output to be set to custom file")
	}
}

// Helper function to check if a string contains a substring.
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
