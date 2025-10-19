package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/layout/domain/model"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRenderService tests constructor.
func TestNewRenderService(t *testing.T) {
	rs := NewRenderService()
	assert.NotNil(t, rs)
}

// TestRender_SimpleContent tests rendering without borders or padding.
func TestRender_SimpleContent(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "single line",
			content:  "Hello",
			expected: "Hello",
		},
		{
			name:     "single char",
			content:  "X",
			expected: "X",
		},
		{
			name:     "with spaces",
			content:  "Hello World",
			expected: "Hello World",
		},
		{
			name:     "two lines",
			content:  "Hello\nWorld",
			expected: "Hello\nWorld",
		},
		{
			name:     "three lines",
			content:  "A\nB\nC",
			expected: "A\nB\nC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content)
			output := rs.Render(box)
			assert.Equal(t, tt.expected, output)
		})
	}
}

// TestRender_WithBorder tests rendering with borders.
func TestRender_WithBorder(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:    "single line with border",
			content: "Hi",
			expected: strings.Join([]string{
				"┌────┐",
				"│ Hi │",
				"└────┘",
			}, "\n"),
		},
		{
			name:    "single char with border",
			content: "X",
			expected: strings.Join([]string{
				"┌───┐",
				"│ X │",
				"└───┘",
			}, "\n"),
		},
		{
			name:    "two lines with border",
			content: "Hi\nBye",
			expected: strings.Join([]string{
				"┌─────┐",
				"│ Hi  │",
				"│ Bye │",
				"└─────┘",
			}, "\n"),
		},
		{
			name:    "different line lengths",
			content: "A\nBB\nCCC",
			expected: strings.Join([]string{
				"┌─────┐",
				"│ A   │",
				"│ BB  │",
				"│ CCC │",
				"└─────┘",
			}, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithBorder(true)
			output := rs.Render(box)
			assert.Equal(t, tt.expected, output,
				"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, tt.expected)
		})
	}
}

// TestRender_WithPadding tests rendering with padding.
func TestRender_WithPadding(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name     string
		content  string
		padding  value.Spacing
		expected string
	}{
		{
			name:    "uniform padding 1",
			content: "Hi",
			padding: value.NewSpacingAll(1),
			expected: strings.Join([]string{
				"    ",
				" Hi ",
				"    ",
			}, "\n"),
		},
		{
			name:    "vertical horizontal padding",
			content: "Hi",
			padding: value.NewSpacingVH(1, 2),
			expected: strings.Join([]string{
				"      ",
				"  Hi  ",
				"      ",
			}, "\n"),
		},
		{
			name:     "left padding only",
			content:  "Test",
			padding:  value.NewSpacing(0, 0, 0, 2),
			expected: "  Test",
		},
		{
			name:     "right padding only",
			content:  "Test",
			padding:  value.NewSpacing(0, 2, 0, 0),
			expected: "Test  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithPadding(tt.padding)
			output := rs.Render(box)
			assert.Equal(t, tt.expected, output,
				"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, tt.expected)
		})
	}
}

// TestRender_WithMargin tests rendering with margin.
func TestRender_WithMargin(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name     string
		content  string
		margin   value.Spacing
		expected string
	}{
		{
			name:    "uniform margin 1",
			content: "Hi",
			margin:  value.NewSpacingAll(1),
			expected: strings.Join([]string{
				"    ",
				" Hi ",
				"    ",
			}, "\n"),
		},
		{
			name:    "vertical horizontal margin",
			content: "X",
			margin:  value.NewSpacingVH(1, 2),
			expected: strings.Join([]string{
				"     ",
				"  X  ",
				"     ",
			}, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithMargin(tt.margin)
			output := rs.Render(box)
			assert.Equal(t, tt.expected, output,
				"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, tt.expected)
		})
	}
}

// TestRender_BorderAndPadding tests combining border with padding.
func TestRender_BorderAndPadding(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name     string
		content  string
		padding  value.Spacing
		expected string
	}{
		{
			name:    "border with uniform padding",
			content: "Hi",
			padding: value.NewSpacingAll(1),
			expected: strings.Join([]string{
				"┌──────┐",
				"│      │",
				"│  Hi  │",
				"│      │",
				"└──────┘",
			}, "\n"),
		},
		{
			name:    "border with vertical padding",
			content: "Test",
			padding: value.NewSpacingVH(1, 0),
			expected: strings.Join([]string{
				"┌──────┐",
				"│      │",
				"│ Test │",
				"│      │",
				"└──────┘",
			}, "\n"),
		},
		{
			name:    "border with horizontal padding",
			content: "Hi",
			padding: value.NewSpacingVH(0, 2),
			// Content: 2, explicit padding: 2+2, implicit padding (border): 1+1
			// Total: 1 + 2 + 2 + 2 + 1 = 8 inner width (symmetric)
			expected: strings.Join([]string{
				"┌────────┐",
				"│   Hi   │",
				"└────────┘",
			}, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithPadding(tt.padding).
				WithBorder(true)
			output := rs.Render(box)
			assert.Equal(t, tt.expected, output,
				"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, tt.expected)
		})
	}
}

// TestRender_AllLayers tests rendering with all box model layers.
func TestRender_AllLayers(t *testing.T) {
	rs := NewRenderService()

	t.Run("padding border margin all 1", func(t *testing.T) {
		box := model.NewBox("X").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithMargin(value.NewSpacingAll(1))

		expected := strings.Join([]string{
			"         ",
			" ┌─────┐ ",
			" │     │ ",
			" │  X  │ ",
			" │     │ ",
			" └─────┘ ",
			"         ",
		}, "\n")

		output := rs.Render(box)
		assert.Equal(t, expected, output,
			"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, expected)
	})

	t.Run("different values for each layer", func(t *testing.T) {
		box := model.NewBox("A").
			WithPadding(value.NewSpacingVH(0, 1)).
			WithBorder(true).
			WithMargin(value.NewSpacingVH(1, 2))

		// Content: 1, explicit padding: 1+1, implicit padding (border): 1+1, border: 1+1, margin: 2+2
		// Inner width: 1 + 1 + 1 + 1 + 1 = 5
		// Total width: 5 + 2 (border) + 4 (margin) = 11
		expected := strings.Join([]string{
			"           ",
			"  ┌─────┐  ",
			"  │  A  │  ",
			"  └─────┘  ",
			"           ",
		}, "\n")

		output := rs.Render(box)
		assert.Equal(t, expected, output,
			"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, expected)
	})
}

// TestRender_MultiLine tests multi-line content rendering.
func TestRender_MultiLine(t *testing.T) {
	rs := NewRenderService()

	t.Run("three lines with border and padding", func(t *testing.T) {
		box := model.NewBox("Line 1\nLine 2\nLine 3").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true)

		expected := strings.Join([]string{
			"┌──────────┐",
			"│          │",
			"│  Line 1  │",
			"│  Line 2  │",
			"│  Line 3  │",
			"│          │",
			"└──────────┘",
		}, "\n")

		output := rs.Render(box)
		assert.Equal(t, expected, output,
			"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, expected)
	})

	t.Run("different line lengths with border", func(t *testing.T) {
		box := model.NewBox("Short\nMedium line\nLong content here").
			WithBorder(true)

		// Max content width = 17 ("Long content here")
		// With implicit padding: 1 left + 17 + 1 right = 19 inner width
		expected := strings.Join([]string{
			"┌───────────────────┐",
			"│ Short             │",
			"│ Medium line       │",
			"│ Long content here │",
			"└───────────────────┘",
		}, "\n")

		output := rs.Render(box)
		assert.Equal(t, expected, output,
			"Output mismatch:\nGot:\n%s\n\nWant:\n%s", output, expected)
	})
}

// TestRender_EdgeCases tests boundary conditions.
func TestRender_EdgeCases(t *testing.T) {
	rs := NewRenderService()

	t.Run("single newline", func(t *testing.T) {
		box := model.NewBox("\n")
		output := rs.Render(box)
		expected := "\n"
		assert.Equal(t, expected, output)
	})

	t.Run("trailing newline", func(t *testing.T) {
		box := model.NewBox("Text\n")
		output := rs.Render(box)
		expected := "Text\n"
		assert.Equal(t, expected, output)
	})

	t.Run("only spaces", func(t *testing.T) {
		box := model.NewBox("   ")
		output := rs.Render(box)
		assert.Equal(t, "   ", output)
	})

	t.Run("zero padding", func(t *testing.T) {
		box := model.NewBox("Test").WithPadding(value.NewSpacingZero())
		output := rs.Render(box)
		assert.Equal(t, "Test", output)
	})

	t.Run("zero margin", func(t *testing.T) {
		box := model.NewBox("Test").WithMargin(value.NewSpacingZero())
		output := rs.Render(box)
		assert.Equal(t, "Test", output)
	})

	t.Run("border only (no padding/margin)", func(t *testing.T) {
		box := model.NewBox("X").WithBorder(true)
		expected := strings.Join([]string{
			"┌───┐",
			"│ X │",
			"└───┘",
		}, "\n")
		output := rs.Render(box)
		assert.Equal(t, expected, output)
	})
}

// TestRenderNode tests node rendering.
func TestRenderNode(t *testing.T) {
	rs := NewRenderService()

	t.Run("single node", func(t *testing.T) {
		box := model.NewBox("Hello").WithBorder(true)
		node := model.NewNode(box)

		expected := strings.Join([]string{
			"┌───────┐",
			"│ Hello │",
			"└───────┘",
		}, "\n")

		output := rs.RenderNode(node)
		assert.Equal(t, expected, output)
	})

	t.Run("node with children (simplified, renders root only)", func(t *testing.T) {
		parentBox := model.NewBox("Parent").WithBorder(true)
		childBox := model.NewBox("Child")

		child := model.NewNode(childBox)
		parent := model.NewNode(parentBox).AddChild(child)

		// For Day 3, RenderNode only renders root box (ignores children)
		// "Parent" (6 chars) + implicit padding (1 left, 1 right) = 8 inner width
		expected := strings.Join([]string{
			"┌────────┐",
			"│ Parent │",
			"└────────┘",
		}, "\n")

		output := rs.RenderNode(parent)
		assert.Equal(t, expected, output)
	})
}

// TestRender_RealWorldScenarios tests realistic usage patterns.
func TestRender_RealWorldScenarios(t *testing.T) {
	rs := NewRenderService()

	t.Run("dialog box", func(t *testing.T) {
		box := model.NewBox("Are you sure?").
			WithPadding(value.NewSpacingVH(1, 2)).
			WithBorder(true).
			WithMargin(value.NewSpacingAll(1))

		output := rs.Render(box)
		require.NotEmpty(t, output)

		// Check structure
		lines := strings.Split(output, "\n")
		assert.True(t, len(lines) >= 5, "should have at least 5 lines")
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")
		assert.Contains(t, output, "Are you sure?")
	})

	t.Run("status bar", func(t *testing.T) {
		box := model.NewBox("Ready").
			WithPadding(value.NewSpacing(0, 1, 0, 1))

		output := rs.Render(box)
		assert.Equal(t, " Ready ", output)
	})

	t.Run("menu item", func(t *testing.T) {
		box := model.NewBox("File").
			WithPadding(value.NewSpacing(0, 2, 0, 2))

		output := rs.Render(box)
		assert.Equal(t, "  File  ", output)
	})

	t.Run("card with title and content", func(t *testing.T) {
		content := "Title\n\nDescription here\n\n[OK] [Cancel]"
		box := model.NewBox(content).
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true)

		output := rs.Render(box)
		require.NotEmpty(t, output)

		lines := strings.Split(output, "\n")
		assert.True(t, len(lines) >= 9, "should have multiple lines")
		assert.Contains(t, output, "Title")
		assert.Contains(t, output, "Description")
		assert.Contains(t, output, "[OK]")
	})
}

// TestRender_GoldenFiles tests visual output with golden files.
func TestRender_GoldenFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping golden file tests in short mode")
	}

	rs := NewRenderService()
	testdataDir := "testdata"

	tests := []struct {
		name       string
		box        *model.Box
		goldenFile string
	}{
		{
			name:       "simple",
			box:        model.NewBox("Hello"),
			goldenFile: "render_simple.golden",
		},
		{
			name:       "border",
			box:        model.NewBox("Hello").WithBorder(true),
			goldenFile: "render_border.golden",
		},
		{
			name: "padding",
			box: model.NewBox("Test").
				WithPadding(value.NewSpacingAll(1)),
			goldenFile: "render_padding.golden",
		},
		{
			name: "complex",
			box: model.NewBox("Hello\nWorld").
				WithPadding(value.NewSpacingAll(1)).
				WithBorder(true).
				WithMargin(value.NewSpacingAll(1)),
			goldenFile: "render_complex.golden",
		},
		{
			name: "dialog",
			box: model.NewBox("Are you sure?").
				WithPadding(value.NewSpacingVH(1, 2)).
				WithBorder(true).
				WithMargin(value.NewSpacingVH(1, 0)),
			goldenFile: "render_dialog.golden",
		},
	}

	// Update mode: set UPDATE_GOLDEN=1 to regenerate golden files
	updateGolden := os.Getenv("UPDATE_GOLDEN") == "1"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := rs.Render(tt.box)
			goldenPath := filepath.Join(testdataDir, tt.goldenFile)

			if updateGolden {
				// Update golden file
				err := os.WriteFile(goldenPath, []byte(output), 0644)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", goldenPath)
				return
			}

			// Compare with golden file
			golden, err := os.ReadFile(goldenPath)
			if os.IsNotExist(err) {
				// Golden file doesn't exist, create it
				err := os.WriteFile(goldenPath, []byte(output), 0644)
				require.NoError(t, err)
				t.Logf("Created golden file: %s", goldenPath)
				return
			}
			require.NoError(t, err, "failed to read golden file")

			// Normalize line endings for cross-platform compatibility
			normalizedGolden := strings.ReplaceAll(string(golden), "\r\n", "\n")
			normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")

			assert.Equal(t, normalizedGolden, normalizedOutput,
				"Output differs from golden file %s:\n\nGot:\n%s\n\nWant:\n%s",
				goldenPath, normalizedOutput, normalizedGolden)
		})
	}
}

// TestRender_Unicode tests rendering with Unicode content.
func TestRender_Unicode(t *testing.T) {
	rs := NewRenderService()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "CJK characters",
			content: "你好",
		},
		{
			name:    "emoji",
			content: "👋",
		},
		{
			name:    "mixed",
			content: "Hello 世界 👋",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithBorder(true)
			output := rs.Render(box)

			// Verify structure (borders present)
			assert.Contains(t, output, "┌")
			assert.Contains(t, output, "└")
			assert.Contains(t, output, "│")
			assert.Contains(t, output, tt.content)

			// Verify multi-line structure
			lines := strings.Split(output, "\n")
			assert.Equal(t, 3, len(lines), "should have 3 lines (top, content, bottom)")
		})
	}
}

// TestRender_VisualVerification prints output for manual verification.
func TestRender_VisualVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping visual verification in short mode")
	}

	rs := NewRenderService()

	t.Run("print examples", func(t *testing.T) {
		examples := []struct {
			name string
			box  *model.Box
		}{
			{
				name: "Simple border",
				box:  model.NewBox("Hello, World!").WithBorder(true),
			},
			{
				name: "With padding",
				box: model.NewBox("Padded Box").
					WithPadding(value.NewSpacingAll(2)).
					WithBorder(true),
			},
			{
				name: "Full box model",
				box: model.NewBox("Complete\nExample").
					WithPadding(value.NewSpacingAll(1)).
					WithBorder(true).
					WithMargin(value.NewSpacingAll(1)),
			},
			{
				name: "Dialog",
				box: model.NewBox("Are you sure you want to quit?").
					WithPadding(value.NewSpacingVH(1, 2)).
					WithBorder(true).
					WithMargin(value.NewSpacingVH(2, 4)),
			},
		}

		for _, ex := range examples {
			t.Logf("\n=== %s ===\n%s\n", ex.name, rs.Render(ex.box))
		}
	})
}
