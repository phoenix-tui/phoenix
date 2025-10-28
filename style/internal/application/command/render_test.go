package command

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/core"
	"github.com/phoenix-tui/phoenix/style/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/style/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/phoenix-tui/phoenix/style/internal/infrastructure/ansi"
)

// Helper to create RenderCommand with all services
func newRenderCommand() *RenderCommand {
	colorAdapter := service2.NewColorAdapter()
	spacingCalculator := service2.NewSpacingCalculator()
	textAligner := service2.NewTextAligner()
	ansiGenerator := ansi.NewANSICodeGenerator()

	return NewRenderCommand(
		colorAdapter,
		spacingCalculator,
		textAligner,
		ansiGenerator,
	)
}

// TestRenderCommand_Execute_EmptyContent tests rendering empty content
func TestRenderCommand_Execute_EmptyContent(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle()

	output, err := cmd.Execute(style, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

// TestRenderCommand_Execute_PlainText tests rendering plain text without styling
func TestRenderCommand_Execute_PlainText(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle()
	content := "Hello, World!"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Plain text should be unchanged
	if output != content {
		t.Errorf("expected %q, got %q", content, output)
	}
}

// TestRenderCommand_Execute_ForegroundColor tests foreground color application
func TestRenderCommand_Execute_ForegroundColor(t *testing.T) {
	cmd := newRenderCommand()
	color := value2.RGB(255, 0, 0) // Red
	style := model.NewStyle().Foreground(color)
	content := "Red text"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain ANSI escape codes
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI escape codes in output")
	}

	// Should contain reset code
	if !strings.Contains(output, "\x1b[0m") {
		t.Error("expected ANSI reset code in output")
	}
}

// TestRenderCommand_Execute_BackgroundColor tests background color application
func TestRenderCommand_Execute_BackgroundColor(t *testing.T) {
	cmd := newRenderCommand()
	color := value2.RGB(0, 0, 255) // Blue
	style := model.NewStyle().Background(color)
	content := "Blue background"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain ANSI escape codes
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI escape codes in output")
	}
}

// TestRenderCommand_Execute_BothColors tests both foreground and background colors
func TestRenderCommand_Execute_BothColors(t *testing.T) {
	cmd := newRenderCommand()
	fg := value2.RGB(255, 255, 255) // White
	bg := value2.RGB(0, 0, 0)       // Black
	style := model.NewStyle().
		Foreground(fg).
		Background(bg)
	content := "White on black"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain multiple ANSI codes
	ansiCount := strings.Count(output, "\x1b[")
	if ansiCount < 2 {
		t.Errorf("expected at least 2 ANSI codes, got %d", ansiCount)
	}
}

// TestRenderCommand_Execute_Bold tests bold text decoration
func TestRenderCommand_Execute_Bold(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().Bold(true)
	content := "Bold text"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain bold ANSI code (1m)
	if !strings.Contains(output, "\x1b[1m") {
		t.Error("expected bold ANSI code in output")
	}
}

// TestRenderCommand_Execute_Italic tests italic text decoration
func TestRenderCommand_Execute_Italic(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().Italic(true)
	content := "Italic text"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain italic ANSI code (3m)
	if !strings.Contains(output, "\x1b[3m") {
		t.Error("expected italic ANSI code in output")
	}
}

// TestRenderCommand_Execute_Underline tests underline text decoration
func TestRenderCommand_Execute_Underline(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().Underline(true)
	content := "Underlined text"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain underline ANSI code (4m)
	if !strings.Contains(output, "\x1b[4m") {
		t.Error("expected underline ANSI code in output")
	}
}

// TestRenderCommand_Execute_Strikethrough tests strikethrough text decoration
func TestRenderCommand_Execute_Strikethrough(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().Strikethrough(true)
	content := "Strikethrough text"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain strikethrough ANSI code (9m)
	if !strings.Contains(output, "\x1b[9m") {
		t.Error("expected strikethrough ANSI code in output")
	}
}

// TestRenderCommand_Execute_AllDecorations tests all text decorations together
func TestRenderCommand_Execute_AllDecorations(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true)
	content := "All decorations"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain all decoration codes
	decorations := []string{"\x1b[1m", "\x1b[3m", "\x1b[4m", "\x1b[9m"}
	for _, dec := range decorations {
		if !strings.Contains(output, dec) {
			t.Errorf("expected decoration code %q in output", dec)
		}
	}
}

// TestRenderCommand_Execute_PaddingOnly tests padding application
func TestRenderCommand_Execute_PaddingOnly(t *testing.T) {
	cmd := newRenderCommand()
	padding := value2.NewPadding(1, 2, 1, 2)
	style := model.NewStyle().Padding(padding)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 3 lines (top padding + content + bottom padding)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// First line should be empty (top padding)
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("expected empty top padding line, got %q", lines[0])
	}

	// Middle line should have left and right padding
	if !strings.HasPrefix(lines[1], "  ") {
		t.Error("expected left padding")
	}
	if !strings.HasSuffix(lines[1], "  ") {
		t.Error("expected right padding")
	}

	// Last line should be empty (bottom padding)
	if strings.TrimSpace(lines[2]) != "" {
		t.Errorf("expected empty bottom padding line, got %q", lines[2])
	}
}

// TestRenderCommand_Execute_MarginOnly tests margin application
func TestRenderCommand_Execute_MarginOnly(t *testing.T) {
	cmd := newRenderCommand()
	margin := value2.NewMargin(1, 2, 1, 2)
	style := model.NewStyle().Margin(margin)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 3 lines (top margin + content + bottom margin)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// First line should be empty (top margin)
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("expected empty top margin line, got %q", lines[0])
	}

	// Middle line should have left and right margin
	if !strings.HasPrefix(lines[1], "  ") {
		t.Error("expected left margin")
	}
	if !strings.HasSuffix(lines[1], "  ") {
		t.Error("expected right margin")
	}

	// Last line should be empty (bottom margin)
	if strings.TrimSpace(lines[2]) != "" {
		t.Errorf("expected empty bottom margin line, got %q", lines[2])
	}
}

// TestRenderCommand_Execute_BorderOnly tests border application
func TestRenderCommand_Execute_BorderOnly(t *testing.T) {
	cmd := newRenderCommand()
	style := model.NewStyle().Border(value2.NormalBorder)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 3 lines (top border + content + bottom border)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// First line should be top border (NormalBorder uses Unicode: ┌─┐)
	if !strings.Contains(lines[0], "┌") || !strings.Contains(lines[0], "─") {
		t.Errorf("expected top border with ┌─┐, got %q", lines[0])
	}

	// Middle line should have left and right border (│)
	if !strings.HasPrefix(lines[1], "│") {
		t.Error("expected left border │")
	}
	if !strings.HasSuffix(lines[1], "│") {
		t.Error("expected right border │")
	}

	// Last line should be bottom border (└─┘)
	if !strings.Contains(lines[2], "└") || !strings.Contains(lines[2], "─") {
		t.Errorf("expected bottom border with └─┘, got %q", lines[2])
	}
}

// TestRenderCommand_Execute_SelectiveBorders tests selective border sides
func TestRenderCommand_Execute_SelectiveBorders(t *testing.T) {
	tests := []struct {
		name      string
		style     model.Style
		topBorder bool
		botBorder bool
		leftSide  bool
		rightSide bool
	}{
		{
			name:      "top and bottom only",
			style:     model.NewStyle().Border(value2.NormalBorder).BorderTop(true).BorderBottom(true),
			topBorder: true,
			botBorder: true,
			leftSide:  false,
			rightSide: false,
		},
		{
			name:      "left and right only",
			style:     model.NewStyle().Border(value2.NormalBorder).BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(true),
			topBorder: false,
			botBorder: false,
			leftSide:  true,
			rightSide: true,
		},
		{
			name:      "top only",
			style:     model.NewStyle().Border(value2.NormalBorder).BorderTop(true),
			topBorder: true,
			botBorder: false,
			leftSide:  false,
			rightSide: false,
		},
	}

	cmd := newRenderCommand()
	content := "Hello"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := cmd.Execute(tt.style, content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lines := strings.Split(output, "\n")

			if tt.topBorder {
				if !strings.Contains(lines[0], "─") {
					t.Error("expected top border with ─")
				}
			}

			if tt.leftSide || tt.rightSide {
				for _, line := range lines {
					if tt.leftSide && len(line) > 0 && !strings.HasPrefix(line, "│") {
						t.Errorf("expected left border │ in line %q", line)
						break
					}
					if tt.rightSide && len(line) > 0 && !strings.HasSuffix(line, "│") {
						t.Errorf("expected right border │ in line %q", line)
						break
					}
				}
			}

			if tt.botBorder {
				if !strings.Contains(lines[len(lines)-1], "─") {
					t.Error("expected bottom border with ─")
				}
			}
		})
	}
}

// TestRenderCommand_Execute_BorderPlusPadding tests border + padding combination
func TestRenderCommand_Execute_BorderPlusPadding(t *testing.T) {
	cmd := newRenderCommand()
	padding := value2.NewPadding(1, 1, 1, 1)
	style := model.NewStyle().
		Border(value2.NormalBorder).
		Padding(padding)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 5 lines:
	// 1. Top border
	// 2. Top padding (inside border)
	// 3. Content with left/right padding
	// 4. Bottom padding (inside border)
	// 5. Bottom border
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d: %v", len(lines), lines)
	}

	// First and last lines should be borders (NormalBorder uses Unicode)
	if !strings.Contains(lines[0], "┌") {
		t.Error("expected top border ┌")
	}
	if !strings.Contains(lines[4], "└") {
		t.Error("expected bottom border └")
	}

	// Middle lines (padding) should have border sides (│)
	for i := 1; i < 4; i++ {
		if !strings.HasPrefix(lines[i], "│") || !strings.HasSuffix(lines[i], "│") {
			t.Errorf("line %d should have border sides │: %q", i, lines[i])
		}
	}
}

// TestRenderCommand_Execute_BorderPlusMargin tests border + margin combination
func TestRenderCommand_Execute_BorderPlusMargin(t *testing.T) {
	cmd := newRenderCommand()
	margin := value2.NewMargin(1, 2, 1, 2)
	style := model.NewStyle().
		Border(value2.NormalBorder).
		Margin(margin)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 5 lines:
	// 1. Top margin
	// 2. Top border (with left/right margin)
	// 3. Content with border and margin
	// 4. Bottom border (with left/right margin)
	// 5. Bottom margin
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}

	// First and last lines should be empty (top/bottom margin)
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("expected empty top margin, got %q", lines[0])
	}
	if strings.TrimSpace(lines[4]) != "" {
		t.Errorf("expected empty bottom margin, got %q", lines[4])
	}

	// Middle lines should have left margin
	for i := 1; i < 4; i++ {
		if !strings.HasPrefix(lines[i], "  ") {
			t.Errorf("line %d should have left margin: %q", i, lines[i])
		}
	}
}

// TestRenderCommand_Execute_AllSpacing tests border + padding + margin together
func TestRenderCommand_Execute_AllSpacing(t *testing.T) {
	cmd := newRenderCommand()
	padding := value2.NewPadding(1, 1, 1, 1)
	margin := value2.NewMargin(1, 1, 1, 1)
	style := model.NewStyle().
		Border(value2.NormalBorder).
		Padding(padding).
		Margin(margin)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have 7 lines:
	// 1. Top margin
	// 2. Top border (with left/right margin)
	// 3. Top padding (with border and margin)
	// 4. Content (with padding, border, margin)
	// 5. Bottom padding (with border and margin)
	// 6. Bottom border (with left/right margin)
	// 7. Bottom margin
	if len(lines) != 7 {
		t.Errorf("expected 7 lines, got %d", len(lines))
	}
}

// TestRenderCommand_Execute_UnicodeContent tests Unicode content (emoji, CJK)
func TestRenderCommand_Execute_UnicodeContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"emoji", "Hello 👋 World 🌍"},
		{"CJK", "你好世界"},
		{"mixed", "Hello 世界 👋"},
		{"grapheme clusters", "e\u0301"}, // é as combining characters
	}

	cmd := newRenderCommand()
	padding := value2.NewPadding(1, 1, 1, 1)
	style := model.NewStyle().
		Border(value2.NormalBorder).
		Padding(padding)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := cmd.Execute(style, tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should contain the original content
			if !strings.Contains(output, tt.content) {
				t.Errorf("output should contain original content %q", tt.content)
			}

			// Should have proper structure (border + padding)
			lines := strings.Split(output, "\n")
			if len(lines) < 3 {
				t.Errorf("expected at least 3 lines, got %d", len(lines))
			}
		})
	}
}

// TestRenderCommand_Execute_WidthConstraint tests width constraint with alignment
func TestRenderCommand_Execute_WidthConstraint(t *testing.T) {
	cmd := newRenderCommand()
	align := value2.NewAlignment(value2.AlignLeft, value2.AlignTop)
	style := model.NewStyle().Width(10).Align(align)
	content := "Short"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Content should be aligned within width 10 (padded with spaces)
	actualWidth := core.StringWidth(strings.TrimRight(output, "\n"))
	if actualWidth != 10 {
		t.Errorf("expected width 10 with alignment, got %d", actualWidth)
	}
}

// TestRenderCommand_Execute_HeightConstraint tests height constraint with alignment
func TestRenderCommand_Execute_HeightConstraint(t *testing.T) {
	cmd := newRenderCommand()
	align := value2.NewAlignment(value2.AlignLeft, value2.AlignTop)
	style := model.NewStyle().Height(5).Align(align)
	content := "Line1\nLine2"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Should have exactly 5 lines with alignment (empty lines added)
	if len(lines) != 5 {
		t.Errorf("expected 5 lines with alignment, got %d", len(lines))
	}
}

// TestRenderCommand_Execute_Alignment tests text alignment
func TestRenderCommand_Execute_Alignment(t *testing.T) {
	tests := []struct {
		name      string
		alignment value2.Alignment
		width     int
		content   string
	}{
		{
			name:      "left aligned",
			alignment: value2.NewAlignment(value2.AlignLeft, value2.AlignTop),
			width:     20,
			content:   "Hello",
		},
		{
			name:      "center aligned",
			alignment: value2.NewAlignment(value2.AlignCenter, value2.AlignMiddle),
			width:     20,
			content:   "Hello",
		},
		{
			name:      "right aligned",
			alignment: value2.NewAlignment(value2.AlignRight, value2.AlignBottom),
			width:     20,
			content:   "Hello",
		},
	}

	cmd := newRenderCommand()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := model.NewStyle().
				Width(tt.width).
				Align(tt.alignment)

			output, err := cmd.Execute(style, tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Output should contain the content
			if !strings.Contains(output, tt.content) {
				t.Errorf("output should contain %q", tt.content)
			}

			// Verify alignment based on position of content
			line := strings.TrimRight(output, "\n")
			idx := strings.Index(line, tt.content)

			switch tt.alignment.Horizontal() {
			case value2.AlignLeft:
				if idx != 0 {
					t.Errorf("expected left alignment, content at position %d", idx)
				}
			case value2.AlignCenter:
				expected := (tt.width - len(tt.content)) / 2
				if idx < expected-1 || idx > expected+1 {
					t.Errorf("expected center alignment at ~%d, got %d", expected, idx)
				}
			case value2.AlignRight:
				expected := tt.width - len(tt.content)
				if idx != expected {
					t.Errorf("expected right alignment at %d, got %d", expected, idx)
				}
			}
		})
	}
}

// TestRenderCommand_Execute_ColorAdaptation tests color adaptation (TrueColor → ANSI256 → ANSI16)
func TestRenderCommand_Execute_ColorAdaptation(t *testing.T) {
	tests := []struct {
		name       string
		capability value2.TerminalCapability
	}{
		{"TrueColor", value2.TrueColor},
		{"ANSI256", value2.ANSI256},
		{"ANSI16", value2.ANSI16},
	}

	color := value2.RGB(255, 0, 0) // Red

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRenderCommand()
			style := model.NewStyle().
				Foreground(color).
				TerminalCapability(tt.capability)
			content := "Colored text"

			output, err := cmd.Execute(style, content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should contain ANSI codes (adapted to capability)
			if !strings.Contains(output, "\x1b[") {
				t.Error("expected ANSI escape codes in output")
			}

			// Different capabilities produce different codes
			switch tt.capability {
			case value2.TrueColor:
				// Should contain 38;2 (TrueColor foreground)
				if !strings.Contains(output, "38;2") {
					t.Error("expected TrueColor ANSI code")
				}
			case value2.ANSI256:
				// Should contain 38;5 (256-color foreground)
				if !strings.Contains(output, "38;5") {
					t.Error("expected ANSI256 color code")
				}
			case value2.ANSI16:
				// Should contain basic color code (30-37 or 90-97)
				hasBasicColor := false
				// Check normal colors (30-37)
				for i := 30; i <= 37; i++ {
					code := fmt.Sprintf("\x1b[%dm", i)
					if strings.Contains(output, code) {
						hasBasicColor = true
						break
					}
				}
				// Check bright colors (90-97)
				if !hasBasicColor {
					for i := 90; i <= 97; i++ {
						code := fmt.Sprintf("\x1b[%dm", i)
						if strings.Contains(output, code) {
							hasBasicColor = true
							break
						}
					}
				}
				if !hasBasicColor {
					t.Errorf("expected ANSI16 basic color code, got output: %q", output)
				}
			}
		})
	}
}

// TestRenderCommand_Execute_BorderColor tests border color application
func TestRenderCommand_Execute_BorderColor(t *testing.T) {
	cmd := newRenderCommand()
	borderColor := value2.RGB(255, 0, 0) // Red
	style := model.NewStyle().
		Border(value2.NormalBorder).
		BorderColor(borderColor)
	content := "Hello"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain ANSI codes for border color
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI codes for border color")
	}

	// Should contain the border characters (NormalBorder uses Unicode)
	if !strings.Contains(output, "┌") || !strings.Contains(output, "│") {
		t.Error("expected border characters ┌ and │")
	}
}

// TestRenderCommand_Execute_CompleteExample tests everything together
func TestRenderCommand_Execute_CompleteExample(t *testing.T) {
	cmd := newRenderCommand()
	fg := value2.RGB(255, 255, 255)      // White
	bg := value2.RGB(0, 0, 255)          // Blue
	borderColor := value2.RGB(255, 0, 0) // Red
	padding := value2.NewPadding(1, 2, 1, 2)
	margin := value2.NewMargin(1, 1, 1, 1)
	alignment := value2.NewAlignment(value2.AlignCenter, value2.AlignMiddle)

	style := model.NewStyle().
		Foreground(fg).
		Background(bg).
		Border(value2.RoundedBorder).
		BorderColor(borderColor).
		Padding(padding).
		Margin(margin).
		Width(30).
		Height(7).
		Align(alignment).
		Bold(true).
		Italic(true)

	content := "Complete Example 🎨"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain content
	if !strings.Contains(output, content) {
		t.Error("output should contain original content")
	}

	// Should contain ANSI codes
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI codes")
	}

	// Should contain border characters
	if !strings.Contains(output, "╭") || !strings.Contains(output, "│") {
		t.Error("expected rounded border characters")
	}

	// Should have correct number of lines (height constraint)
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	// Height 7 + top margin (1) + bottom margin (1) = 9 total
	if len(lines) < 7 {
		t.Errorf("expected at least 7 lines, got %d", len(lines))
	}
}

// TestRenderCommand_Execute_InvalidStyle tests validation error handling
func TestRenderCommand_Execute_InvalidStyle(t *testing.T) {
	cmd := newRenderCommand()
	// Create invalid style: border sides enabled but no border
	style := model.NewStyle().BorderTop(true)
	content := "Invalid"

	_, err := cmd.Execute(style, content)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "border") {
		t.Errorf("expected error about border, got %v", err)
	}
}

// TestRenderCommand_Execute_MultilineContent tests multiline content handling
func TestRenderCommand_Execute_MultilineContent(t *testing.T) {
	cmd := newRenderCommand()
	padding := value2.NewPadding(1, 1, 1, 1)
	style := model.NewStyle().
		Border(value2.NormalBorder).
		Padding(padding)
	content := "Line 1\nLine 2\nLine 3"

	output, err := cmd.Execute(style, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain all lines
	if !strings.Contains(output, "Line 1") {
		t.Error("missing Line 1")
	}
	if !strings.Contains(output, "Line 2") {
		t.Error("missing Line 2")
	}
	if !strings.Contains(output, "Line 3") {
		t.Error("missing Line 3")
	}

	lines := strings.Split(output, "\n")
	// Should have: top border + top padding + 3 content lines + bottom padding + bottom border = 7 lines
	if len(lines) != 7 {
		t.Errorf("expected 7 lines, got %d", len(lines))
	}
}

// TestRenderCommand_Execute_ZeroDimensions tests zero width/height edge case
func TestRenderCommand_Execute_ZeroDimensions(t *testing.T) {
	cmd := newRenderCommand()

	tests := []struct {
		name  string
		style model.Style
	}{
		{
			name:  "zero width",
			style: model.NewStyle().Width(0),
		},
		{
			name:  "zero height",
			style: model.NewStyle().Height(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := cmd.Execute(tt.style, "content")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should handle gracefully (likely return empty or minimal output)
			if len(output) > 100 {
				t.Errorf("expected minimal output for zero dimensions, got %d chars", len(output))
			}
		})
	}
}
