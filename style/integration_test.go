package style_test

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/style"
	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// Integration tests for complete phoenix/style pipeline.
// These tests verify that all components work together correctly.

func TestIntegration_BasicStyling(t *testing.T) {
	s := style.New().
		Foreground(style.RGB(255, 255, 255)).
		Background(style.RGB(0, 0, 255))

	output := style.Render(s, "Hello")

	// Should contain ANSI codes
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI codes in output")
	}

	// Should contain content
	if !strings.Contains(output, "Hello") {
		t.Error("expected content in output")
	}
}

func TestIntegration_BorderWithPaddingAndMargin(t *testing.T) {
	s := style.New().
		Border(style.RoundedBorder).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 1, 1, 1))

	output := style.Render(s, "Boxed")

	lines := strings.Split(output, "\n")

	// Should have margin + border + padding + content + padding + border + margin
	if len(lines) < 5 {
		t.Errorf("expected at least 5 lines, got %d", len(lines))
	}

	// Should contain border characters
	if !strings.Contains(output, "╭") {
		t.Error("expected rounded border character")
	}
}

func TestIntegration_CompleteStyle(t *testing.T) {
	s := style.New().
		Foreground(style.RGB(255, 255, 255)).
		Background(style.RGB(0, 0, 255)).
		Border(style.RoundedBorder).
		BorderColor(style.RGB(255, 0, 0)).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 1, 1, 1)).
		Bold(true).
		Italic(true)

	output := style.Render(s, "Styled Text")

	// Should contain all elements
	tests := []string{
		"\x1b[",       // ANSI codes
		"Styled Text", // Content
		"╭",           // Border
		"\x1b[1m",     // Bold
		"\x1b[3m",     // Italic
	}

	for _, test := range tests {
		if !strings.Contains(output, test) {
			t.Errorf("expected %q in output", test)
		}
	}
}

func TestIntegration_AlignmentWithSize(t *testing.T) {
	s := style.New().
		Width(20).
		Height(5).
		Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle))

	output := style.Render(s, "Centered")

	lines := strings.Split(output, "\n")

	// Should have exactly 5 lines
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}

	// Content should be in middle line (index 2)
	if !strings.Contains(lines[2], "Centered") {
		t.Error("expected content in middle line")
	}
}

func TestIntegration_UnicodeContent(t *testing.T) {
	content := "Hello 👋 World 🌍"

	s := style.New().
		Border(style.NormalBorder).
		Padding(style.NewPadding(1, 1, 1, 1))

	output := style.Render(s, content)

	// Should contain Unicode content
	if !strings.Contains(output, "👋") {
		t.Error("expected emoji in output")
	}
	if !strings.Contains(output, "🌍") {
		t.Error("expected emoji in output")
	}

	// Should have proper structure
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d", len(lines))
	}
}

func TestIntegration_MultipleStyles(t *testing.T) {
	// Create multiple different styles
	style1 := style.New().Foreground(style.Red).Bold(true)
	style2 := style.New().Background(style.Blue).Italic(true)
	style3 := style.New().Border(style.RoundedBorder)

	// Render same content with different styles
	output1 := style.Render(style1, "Red Bold")
	output2 := style.Render(style2, "Blue Italic")
	output3 := style.Render(style3, "Boxed")

	// Each should be different
	if output1 == output2 || output2 == output3 || output1 == output3 {
		t.Error("different styles should produce different outputs")
	}

	// Verify each has expected characteristics
	if !strings.Contains(output1, "\x1b[1m") {
		t.Error("style1 should contain bold")
	}
	if !strings.Contains(output2, "\x1b[3m") {
		t.Error("style2 should contain italic")
	}
	if !strings.Contains(output3, "╭") {
		t.Error("style3 should contain border")
	}
}

func TestIntegration_EmptyContent(t *testing.T) {
	s := style.New().
		Border(style.NormalBorder).
		Padding(style.NewPadding(1, 1, 1, 1))

	output := style.Render(s, "")

	// Empty content should return empty (or just spacing if enabled)
	if output == "" {
		return // OK
	}

	// If not empty, should at least have valid structure
	lines := strings.Split(output, "\n")
	if len(lines) > 0 && !strings.Contains(lines[0], "┌") {
		t.Error("if rendering empty content with border, should have border structure")
	}
}

func TestIntegration_ColorPresets(t *testing.T) {
	// Test all color presets
	presets := []struct {
		name  string
		color style.Color
	}{
		{"Red", style.Red},
		{"Green", style.Green},
		{"Blue", style.Blue},
		{"Yellow", style.Yellow},
		{"Cyan", style.Cyan},
		{"Magenta", style.Magenta},
		{"White", style.White},
		{"Black", style.Black},
		{"Gray", style.Gray},
	}

	for _, preset := range presets {
		t.Run(preset.name, func(t *testing.T) {
			s := style.New().Foreground(preset.color)
			output := style.Render(s, "Test")

			// Should contain ANSI codes
			if !strings.Contains(output, "\x1b[") {
				t.Errorf("%s color should produce ANSI codes", preset.name)
			}
		})
	}
}

func TestIntegration_BorderPresets(t *testing.T) {
	// Test all border presets
	presets := []struct {
		name   string
		border style.Border
		char   string
	}{
		{"Normal", style.NormalBorder, "┌"},
		{"Rounded", style.RoundedBorder, "╭"},
		{"Thick", style.ThickBorder, "┏"},
		{"Double", style.DoubleBorder, "╔"},
		{"ASCII", style.ASCIIBorder, "+"},
	}

	for _, preset := range presets {
		t.Run(preset.name, func(t *testing.T) {
			s := style.New().Border(preset.border)
			output := style.Render(s, "Test")

			// Should contain expected border character
			if !strings.Contains(output, preset.char) {
				t.Errorf("%s border should contain %q", preset.name, preset.char)
			}
		})
	}
}

func TestIntegration_StylePresets(t *testing.T) {
	// Test all style presets
	tests := []struct {
		name     string
		preset   style.Style
		expected string
	}{
		{"Default", style.DefaultStyle, "Test"},
		{"Bold", style.BoldStyle, "\x1b[1m"},
		{"Italic", style.ItalicStyle, "\x1b[3m"},
		{"Underline", style.UnderlineStyle, "\x1b[4m"},
		{"Strikethrough", style.StrikethroughStyle, "\x1b[9m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := style.Render(tt.preset, "Test")

			if !strings.Contains(output, tt.expected) {
				t.Errorf("%s preset should contain %q", tt.name, tt.expected)
			}
		})
	}
}

func TestIntegration_FluentAPI(t *testing.T) {
	// Test method chaining
	s := style.New().
		Foreground(style.Red).
		Background(style.Blue).
		Bold(true).
		Italic(true).
		Underline(true).
		Border(style.RoundedBorder).
		BorderColor(style.Yellow).
		Padding(style.NewPadding(1, 1, 1, 1)).
		Margin(style.NewMargin(1, 1, 1, 1)).
		Width(30).
		Height(7).
		Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle))

	output := style.Render(s, "Fluent API Test")

	// Should contain all elements
	if !strings.Contains(output, "Fluent API Test") {
		t.Error("should contain content")
	}

	if !strings.Contains(output, "\x1b[") {
		t.Error("should contain ANSI codes")
	}

	if !strings.Contains(output, "╭") {
		t.Error("should contain border")
	}
}

func TestIntegration_TerminalCapabilities(t *testing.T) {
	capabilities := []style.TerminalCapability{
		style.TrueColor,
		style.ANSI256,
		style.ANSI16,
	}

	for _, cap := range capabilities {
		t.Run(cap.String(), func(t *testing.T) {
			s := style.New().
				Foreground(style.RGB(255, 0, 0)).
				TerminalCapability(value2.TerminalCapability(cap))

			output := style.Render(s, "Test")

			// Should contain ANSI codes
			if !strings.Contains(output, "\x1b[") {
				t.Errorf("%s should produce ANSI codes", cap.String())
			}
		})
	}
}

func TestIntegration_MultilineContent(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"

	s := style.New().
		Border(style.NormalBorder).
		Padding(style.NewPadding(1, 1, 1, 1))

	output := style.Render(s, content)

	// All lines should be present
	for i := 1; i <= 3; i++ {
		expected := "Line " + string(rune('0'+i))
		if !strings.Contains(output, expected) {
			t.Errorf("should contain %q", expected)
		}
	}

	// Should have proper structure
	lines := strings.Split(output, "\n")
	if len(lines) < 7 {
		t.Errorf("expected at least 7 lines (border+padding+3 content+padding+border), got %d", len(lines))
	}
}

func TestIntegration_SelectiveBorders(t *testing.T) {
	// Only left and right borders
	s := style.New().
		Border(style.NormalBorder).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(true).
		BorderRight(true)

	output := style.Render(s, "Sides Only")

	lines := strings.Split(output, "\n")

	// Should have left/right borders but no top/bottom
	for _, line := range lines {
		if len(line) > 0 {
			if !strings.HasPrefix(line, "│") || !strings.HasSuffix(line, "│") {
				t.Errorf("line should have left/right borders: %q", line)
			}
		}
	}
}

func TestIntegration_RealWorldExample(t *testing.T) {
	// Simulate real-world usage: styled notification box
	notificationStyle := style.New().
		Foreground(style.White).
		Background(style.RGB(0, 102, 204)). // Blue background
		Border(style.RoundedBorder).
		BorderColor(style.White).
		Padding(style.NewPadding(1, 2, 1, 2)).
		Margin(style.NewMargin(1, 0, 1, 0)).
		Bold(true).
		Width(50)

	message := "🔔 Notification: Your task is complete!"

	output := style.Render(notificationStyle, message)

	// Verify output has all expected elements
	if !strings.Contains(output, message) {
		t.Error("should contain notification message")
	}

	if !strings.Contains(output, "╭") {
		t.Error("should contain rounded border")
	}

	if !strings.Contains(output, "\x1b[") {
		t.Error("should contain ANSI styling")
	}

	// Should have proper structure
	lines := strings.Split(output, "\n")
	if len(lines) < 5 {
		t.Errorf("notification box should have at least 5 lines, got %d", len(lines))
	}
}
