package style_test

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/style"
	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// API tests for phoenix/style public interface.
// These tests verify the user-facing API is intuitive and works correctly.

func TestAPI_New(t *testing.T) {
	s := style.New()

	// Should create valid style
	output := style.Render(s, "Test")
	if output != "Test" {
		t.Errorf("default style should render content unchanged, got %q", output)
	}
}

func TestAPI_ColorConstructors(t *testing.T) {
	tests := []struct {
		name string
		test func() (style.Color, error)
	}{
		{"RGB", func() (style.Color, error) {
			return style.RGB(255, 0, 0), nil
		}},
		{"Hex", func() (style.Color, error) {
			return style.Hex("#FF0000")
		}},
		{"Color256", func() (style.Color, error) {
			return style.Color256(196), nil
		}},
		{"Color16", func() (style.Color, error) {
			return style.Color16(1), nil
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := tt.test()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should be able to use color
			s := style.New().Foreground(color)
			output := style.Render(s, "Test")

			if !strings.Contains(output, "\x1b[") {
				t.Error("expected ANSI codes in styled output")
			}
		})
	}
}

func TestAPI_ColorPresets(t *testing.T) {
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

			if !strings.Contains(output, "\x1b[") {
				t.Errorf("%s preset should produce styled output", preset.name)
			}
		})
	}
}

func TestAPI_BorderConstructors(t *testing.T) {
	// Test pre-defined borders
	borders := []struct {
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

	for _, b := range borders {
		t.Run(b.name, func(t *testing.T) {
			s := style.New().Border(b.border)
			output := style.Render(s, "Test")

			if !strings.Contains(output, b.char) {
				t.Errorf("%s border should contain %q", b.name, b.char)
			}
		})
	}

	// Test custom border
	t.Run("Custom", func(t *testing.T) {
		customBorder := style.NewBorder("*", "*", "*", "*", "*", "*", "*", "*")
		s := style.New().Border(customBorder)
		output := style.Render(s, "Test")

		if !strings.Contains(output, "*") {
			t.Error("custom border should contain *")
		}
	})
}

func TestAPI_SpacingConstructors(t *testing.T) {
	t.Run("NewPadding", func(t *testing.T) {
		padding := style.NewPadding(1, 2, 3, 4)
		s := style.New().Padding(padding)
		output := style.Render(s, "Test")

		lines := strings.Split(output, "\n")
		if len(lines) < 3 {
			t.Error("padding should add lines")
		}
	})

	t.Run("NewMargin", func(t *testing.T) {
		margin := style.NewMargin(1, 2, 3, 4)
		s := style.New().Margin(margin)
		output := style.Render(s, "Test")

		lines := strings.Split(output, "\n")
		if len(lines) < 3 {
			t.Error("margin should add lines")
		}
	})
}

func TestAPI_SizeConstructor(_ *testing.T) {
	size := style.NewSize()

	// Size should be usable (even if empty)
	_ = size
}

func TestAPI_AlignmentConstructor(t *testing.T) {
	align := style.NewAlignment(style.AlignCenter, style.AlignMiddle)
	s := style.New().Width(20).Height(5).Align(align)
	output := style.Render(s, "Test")

	// Should produce aligned output
	lines := strings.Split(output, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines with height constraint, got %d", len(lines))
	}
}

func TestAPI_StylePresets(t *testing.T) {
	presets := []struct {
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

	for _, preset := range presets {
		t.Run(preset.name, func(t *testing.T) {
			output := style.Render(preset.preset, "Test")

			if !strings.Contains(output, preset.expected) {
				t.Errorf("%s preset should contain %q", preset.name, preset.expected)
			}
		})
	}
}

func TestAPI_FluentChaining(t *testing.T) {
	// Test that all methods are chainable
	s := style.New().
		Foreground(style.Red).
		Background(style.Blue).
		Border(style.RoundedBorder).
		BorderColor(style.Yellow).
		Padding(style.NewPadding(1, 1, 1, 1)).
		Margin(style.NewMargin(1, 1, 1, 1)).
		Width(30).
		Height(7).
		Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle)).
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		TerminalCapability(value2.TerminalCapability(style.TrueColor))

	output := style.Render(s, "Test")

	// Should render without errors
	if !strings.Contains(output, "Test") {
		t.Error("fluent chaining should produce valid output")
	}
}

func TestAPI_Render(t *testing.T) {
	t.Run("PlainText", func(t *testing.T) {
		output := style.Render(style.New(), "Hello")
		if output != "Hello" {
			t.Errorf("plain text should be unchanged, got %q", output)
		}
	})

	t.Run("WithColor", func(t *testing.T) {
		s := style.New().Foreground(style.Red)
		output := style.Render(s, "Red")

		if !strings.Contains(output, "\x1b[") {
			t.Error("colored text should contain ANSI codes")
		}
	})

	t.Run("WithBorder", func(t *testing.T) {
		s := style.New().Border(style.RoundedBorder)
		output := style.Render(s, "Boxed")

		if !strings.Contains(output, "╭") {
			t.Error("bordered text should contain border characters")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		output := style.Render(style.New(), "")
		if output != "" {
			t.Errorf("empty content should render empty, got %q", output)
		}
	})
}

func TestAPI_TerminalCapabilityConstants(t *testing.T) {
	// Test that all terminal capability constants are accessible
	capabilities := []style.TerminalCapability{
		style.TrueColor,
		style.ANSI256,
		style.ANSI16,
	}

	for _, cap := range capabilities {
		t.Run(cap.String(), func(t *testing.T) {
			s := style.New().
				Foreground(style.Red).
				TerminalCapability(value2.TerminalCapability(cap))

			output := style.Render(s, "Test")

			if !strings.Contains(output, "\x1b[") {
				t.Errorf("%s should produce ANSI codes", cap.String())
			}
		})
	}
}

func TestAPI_AlignmentConstants(t *testing.T) {
	// Test horizontal alignments
	horizontals := []style.HorizontalAlignment{
		style.AlignLeft,
		style.AlignCenter,
		style.AlignRight,
	}

	for _, h := range horizontals {
		alignment := style.NewAlignment(h, style.AlignTop)
		s := style.New().Width(20).Align(alignment)
		output := style.Render(s, "Test")

		if len(output) == 0 {
			t.Errorf("alignment %v should produce output", h)
		}
	}

	// Test vertical alignments
	verticals := []style.VerticalAlignment{
		style.AlignTop,
		style.AlignMiddle,
		style.AlignBottom,
	}

	for _, v := range verticals {
		alignment := style.NewAlignment(style.AlignLeft, v)
		s := style.New().Height(5).Align(alignment)
		output := style.Render(s, "Test")

		if len(output) == 0 {
			t.Errorf("alignment %v should produce output", v)
		}
	}
}

func TestAPI_ImmutabilityThroughRender(t *testing.T) {
	original := style.New().Foreground(style.Red)

	// Render multiple times
	output1 := style.Render(original, "Test 1")
	output2 := style.Render(original, "Test 2")

	// Original style should still work
	output3 := style.Render(original, "Test 3")

	// All should contain red color codes
	for i, output := range []string{output1, output2, output3} {
		if !strings.Contains(output, "\x1b[") {
			t.Errorf("output %d should contain ANSI codes", i+1)
		}
	}
}

func TestAPI_ErrorHandling(t *testing.T) {
	t.Run("InvalidHex", func(t *testing.T) {
		_, err := style.Hex("invalid")
		if err == nil {
			t.Error("invalid hex should return error")
		}
	})

	t.Run("InvalidStyle", func(_ *testing.T) {
		// Style with border sides but no border
		s := style.New().BorderTop(true)
		output := style.Render(s, "Test")

		// Should handle gracefully (return content or empty)
		_ = output
	})
}

func TestAPI_RealWorldUsage(t *testing.T) {
	// Simulate typical usage patterns

	t.Run("HeaderStyle", func(t *testing.T) {
		headerStyle := style.New().
			Foreground(style.White).
			Background(style.Blue).
			Bold(true).
			Padding(style.NewPadding(0, 2, 0, 2))

		output := style.Render(headerStyle, "Application Header")

		if !strings.Contains(output, "Application Header") {
			t.Error("header should contain text")
		}
	})

	t.Run("ErrorStyle", func(t *testing.T) {
		errorStyle := style.New().
			Foreground(style.Red).
			Bold(true).
			Border(style.ThickBorder)

		output := style.Render(errorStyle, "Error: Something went wrong")

		if !strings.Contains(output, "Error") {
			t.Error("error style should contain error message")
		}
	})

	t.Run("InfoBox", func(t *testing.T) {
		infoBox := style.New().
			Border(style.RoundedBorder).
			BorderColor(style.Cyan).
			Padding(style.NewPadding(1, 2, 1, 2)).
			Margin(style.NewMargin(1, 0, 1, 0))

		output := style.Render(infoBox, "ℹ️  Information: Task completed successfully")

		if !strings.Contains(output, "Information") {
			t.Error("info box should contain message")
		}

		if !strings.Contains(output, "╭") {
			t.Error("info box should have rounded border")
		}
	})
}

func TestAPI_Documentation_Examples(t *testing.T) {
	// Test examples from documentation

	t.Run("QuickStart", func(t *testing.T) {
		// Example from package docs
		s := style.New().
			Foreground(style.RGB(255, 255, 255)).
			Background(style.RGB(0, 0, 255)).
			Bold(true)

		output := style.Render(s, "Hello, World!")

		if !strings.Contains(output, "Hello, World!") {
			t.Error("quick start example should work")
		}
	})

	t.Run("CompleteExample", func(t *testing.T) {
		s := style.New().
			Foreground(style.RGB(255, 255, 255)).
			Background(style.RGB(0, 0, 255)).
			Padding(style.NewPadding(1, 2, 1, 2)).
			Border(style.RoundedBorder)

		output := style.Render(s, "Hello, World!")

		if !strings.Contains(output, "Hello, World!") {
			t.Error("complete example should work")
		}

		if !strings.Contains(output, "╭") {
			t.Error("complete example should have border")
		}
	})
}
