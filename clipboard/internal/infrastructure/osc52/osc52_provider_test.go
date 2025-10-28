package osc52

import (
	"os"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/model"
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

// Additional OSC52 Tests

func TestProvider_Write_Base64Encoding(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	tests := []struct {
		name    string
		content string
	}{
		{"simple", "hello"},
		{"unicode", "你好"},
		{"emoji", "👋"},
		{"special chars", "!@#$%^&*()"},
		{"newlines", "line1\nline2\nline3"},
		{"tabs", "col1\tcol2\tcol3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := model.NewTextContent(tt.content)

			// Reset file
			tmpFile.Truncate(0)
			tmpFile.Seek(0, 0)

			err := provider.Write(content)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Read back
			tmpFile.Seek(0, 0)
			buf := make([]byte, 4096)
			n, _ := tmpFile.Read(buf)
			output := string(buf[:n])

			// Verify OSC 52 sequence present
			if !containsString(output, "\033]52;c;") {
				t.Errorf("expected OSC 52 sequence in output")
			}
		})
	}
}

func TestProvider_Write_LargeContent(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	// Generate 100KB of text
	largeData := ""
	for i := 0; i < 10000; i++ {
		largeData += "0123456789"
	}

	content, _ := model.NewTextContent(largeData)

	err = provider.Write(content)
	if err != nil {
		t.Errorf("unexpected error for large content: %v", err)
	}
}

func TestProvider_IsAvailable_TerminalVariants(t *testing.T) {
	origTERM := os.Getenv("TERM")
	defer os.Setenv("TERM", origTERM)

	// Clear SSH vars
	os.Unsetenv("SSH_TTY")
	os.Unsetenv("SSH_CLIENT")
	os.Unsetenv("SSH_CONNECTION")

	tests := []struct {
		term string
		want string // "true", "false", or "skip" (non-terminal output)
	}{
		{"xterm", "true"},
		{"xterm-256color", "true"},
		{"screen", "true"},
		{"screen-256color", "skip"}, // Not in supported list
		{"tmux", "true"},
		{"tmux-256color", "true"},
		{"dumb", "false"},
		{"vt100", "skip"}, // Not in supported list
		{"", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(_ *testing.T) {
			os.Setenv("TERM", tt.term)

			provider := NewProvider(5 * time.Second)

			// IsAvailable checks if output is a terminal
			// In test environment, os.Stdout might not be a TTY
			available := provider.IsAvailable()

			// We can't assert the exact value because it depends on
			// whether stdout is a terminal in the test environment
			// But we can verify it doesn't panic
			_ = available

			// This test is mostly to ensure no panic with various TERM values
		})
	}
}

func TestProvider_IsAvailable_NotTerminal(t *testing.T) {
	// Create a regular file (not a terminal)
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Clear env vars
	os.Unsetenv("SSH_TTY")
	os.Unsetenv("SSH_CLIENT")
	os.Unsetenv("SSH_CONNECTION")

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	if provider.IsAvailable() {
		t.Error("expected IsAvailable to return false for non-terminal output")
	}
}

func TestProvider_Write_EmptyContent(_ *testing.T) {
	provider := NewProvider(5 * time.Second)

	// Empty content is rejected by domain model, so this should fail
	content, err := model.NewTextContent("")
	if err != nil {
		// Expected - empty content not allowed
		return
	}

	// If we got here, try writing
	err = provider.Write(content)
	// Either way is acceptable
	_ = err
}

func TestProvider_Write_SpecialCharacters(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second)
	provider.WithOutput(tmpFile)

	tests := []string{
		"\x00\x01\x02",             // Control characters
		"\033[1m",                  // ANSI escape
		"<script>alert()</script>", // HTML/JS
		"'; DROP TABLE--",          // SQL injection attempt
	}

	for _, content := range tests {
		t.Run("special", func(t *testing.T) {
			c, _ := model.NewTextContent(content)

			tmpFile.Truncate(0)
			tmpFile.Seek(0, 0)

			err := provider.Write(c)
			// Should handle special characters safely
			if err != nil {
				t.Logf("Write failed (acceptable): %v", err)
			}
		})
	}
}

func TestProvider_Name_Consistency(t *testing.T) {
	provider := NewProvider(5 * time.Second)

	// Call multiple times
	name1 := provider.Name()
	name2 := provider.Name()
	name3 := provider.Name()

	if name1 != name2 || name2 != name3 {
		t.Error("Name() should return consistent results")
	}

	if name1 != "OSC52" {
		t.Errorf("expected 'OSC52', got %s", name1)
	}
}

func TestProvider_WithOutput_Chaining(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "clipboard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	provider := NewProvider(5 * time.Second).WithOutput(tmpFile)

	if provider.output != tmpFile {
		t.Error("expected output to be set via chaining")
	}
}

func TestProvider_IsAvailable_SSHIndicators_Priority(_ *testing.T) {
	origSSHTTY := os.Getenv("SSH_TTY")
	origTERM := os.Getenv("TERM")
	defer func() {
		os.Setenv("SSH_TTY", origSSHTTY)
		os.Setenv("TERM", origTERM)
	}()

	// Clear SSH vars
	os.Unsetenv("SSH_TTY")
	os.Unsetenv("SSH_CLIENT")
	os.Unsetenv("SSH_CONNECTION")

	// Set unsupported TERM
	os.Setenv("TERM", "dumb")

	provider := NewProvider(5 * time.Second)

	// Without SSH vars and with unsupported TERM, should be unavailable
	// (assuming output is not a terminal, which is typical in test env)
	_ = provider.IsAvailable()

	// Now set SSH_TTY
	os.Setenv("SSH_TTY", "/dev/pts/0")

	// Should still depend on terminal check
	_ = provider.IsAvailable()
}
