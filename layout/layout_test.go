package layout

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/layout/internal/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewBox tests box creation.
func TestNewBox(t *testing.T) {
	t.Run("creates box with content", func(t *testing.T) {
		box := NewBox("Hello")
		assert.NotNil(t, box)
		assert.NotNil(t, box.domain)
		assert.Equal(t, "Hello", box.domain.Content())
	})

	t.Run("panics with empty content", func(t *testing.T) {
		assert.Panics(t, func() {
			NewBox("")
		}, "should panic with empty content")
	})

	t.Run("creates box with multi-line content", func(t *testing.T) {
		box := NewBox("Line1\nLine2\nLine3")
		assert.NotNil(t, box)
		assert.Equal(t, "Line1\nLine2\nLine3", box.domain.Content())
	})
}

// TestBox_SizeConstraints tests size constraint methods.
func TestBox_SizeConstraints(t *testing.T) {
	t.Run("Width sets exact width", func(t *testing.T) {
		box := NewBox("Hi").Width(20)
		assert.Equal(t, 20, box.domain.Size().Width())
	})

	t.Run("Height sets exact height", func(t *testing.T) {
		box := NewBox("Hi").Height(10)
		assert.Equal(t, 10, box.domain.Size().Height())
	})

	t.Run("MinWidth sets minimum width", func(t *testing.T) {
		box := NewBox("Hi").MinWidth(15)
		assert.Equal(t, 15, box.domain.Size().MinWidth())
	})

	t.Run("MinHeight sets minimum height", func(t *testing.T) {
		box := NewBox("Hi").MinHeight(5)
		assert.Equal(t, 5, box.domain.Size().MinHeight())
	})

	t.Run("MaxWidth sets maximum width", func(t *testing.T) {
		box := NewBox("Hi").MaxWidth(30)
		assert.Equal(t, 30, box.domain.Size().MaxWidth())
	})

	t.Run("MaxHeight sets maximum height", func(t *testing.T) {
		box := NewBox("Hi").MaxHeight(8)
		assert.Equal(t, 8, box.domain.Size().MaxHeight())
	})

	t.Run("chain multiple size constraints", func(t *testing.T) {
		box := NewBox("Hi").
			Width(20).
			Height(10).
			MinWidth(10).
			MaxWidth(40)

		assert.Equal(t, 20, box.domain.Size().Width())
		assert.Equal(t, 10, box.domain.Size().Height())
		assert.Equal(t, 10, box.domain.Size().MinWidth())
		assert.Equal(t, 40, box.domain.Size().MaxWidth())
	})
}

// TestBox_Padding tests padding methods.
func TestBox_Padding(t *testing.T) {
	t.Run("Padding sets all sides individually", func(t *testing.T) {
		box := NewBox("Hi").Padding(1, 2, 3, 4)
		padding := box.domain.Padding()
		assert.Equal(t, 1, padding.Top())
		assert.Equal(t, 2, padding.Right())
		assert.Equal(t, 3, padding.Bottom())
		assert.Equal(t, 4, padding.Left())
	})

	t.Run("PaddingAll sets same padding for all sides", func(t *testing.T) {
		box := NewBox("Hi").PaddingAll(2)
		padding := box.domain.Padding()
		assert.Equal(t, 2, padding.Top())
		assert.Equal(t, 2, padding.Right())
		assert.Equal(t, 2, padding.Bottom())
		assert.Equal(t, 2, padding.Left())
	})

	t.Run("PaddingVH sets vertical and horizontal", func(t *testing.T) {
		box := NewBox("Hi").PaddingVH(1, 3)
		padding := box.domain.Padding()
		assert.Equal(t, 1, padding.Top())
		assert.Equal(t, 3, padding.Right())
		assert.Equal(t, 1, padding.Bottom())
		assert.Equal(t, 3, padding.Left())
	})

	t.Run("chain padding methods", func(t *testing.T) {
		box := NewBox("Hi").
			PaddingAll(1).
			Padding(2, 2, 2, 2) // Override with different values

		padding := box.domain.Padding()
		assert.Equal(t, 2, padding.Top())
		assert.Equal(t, 2, padding.Right())
	})
}

// TestBox_Border tests border methods.
func TestBox_Border(t *testing.T) {
	t.Run("Border enables border", func(t *testing.T) {
		box := NewBox("Hi").Border()
		assert.True(t, box.domain.HasBorder())
	})

	t.Run("NoBorder disables border", func(t *testing.T) {
		box := NewBox("Hi").Border().NoBorder()
		assert.False(t, box.domain.HasBorder())
	})

	t.Run("default has no border", func(t *testing.T) {
		box := NewBox("Hi")
		assert.False(t, box.domain.HasBorder())
	})
}

// TestBox_Margin tests margin methods.
func TestBox_Margin(t *testing.T) {
	t.Run("Margin sets all sides individually", func(t *testing.T) {
		box := NewBox("Hi").Margin(1, 2, 3, 4)
		margin := box.domain.Margin()
		assert.Equal(t, 1, margin.Top())
		assert.Equal(t, 2, margin.Right())
		assert.Equal(t, 3, margin.Bottom())
		assert.Equal(t, 4, margin.Left())
	})

	t.Run("MarginAll sets same margin for all sides", func(t *testing.T) {
		box := NewBox("Hi").MarginAll(3)
		margin := box.domain.Margin()
		assert.Equal(t, 3, margin.Top())
		assert.Equal(t, 3, margin.Right())
		assert.Equal(t, 3, margin.Bottom())
		assert.Equal(t, 3, margin.Left())
	})

	t.Run("MarginVH sets vertical and horizontal", func(t *testing.T) {
		box := NewBox("Hi").MarginVH(2, 5)
		margin := box.domain.Margin()
		assert.Equal(t, 2, margin.Top())
		assert.Equal(t, 5, margin.Right())
		assert.Equal(t, 2, margin.Bottom())
		assert.Equal(t, 5, margin.Left())
	})
}

// TestBox_Alignment tests alignment methods.
func TestBox_Alignment(t *testing.T) {
	t.Run("AlignLeft sets left alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignLeft()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignLeft, align.Horizontal())
		assert.Equal(t, value.AlignTop, align.Vertical())
	})

	t.Run("AlignCenter sets center alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignCenter()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignCenter, align.Horizontal())
		assert.Equal(t, value.AlignMiddle, align.Vertical())
	})

	t.Run("AlignRight sets right alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignRight()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignRight, align.Horizontal())
		assert.Equal(t, value.AlignTop, align.Vertical())
	})

	t.Run("AlignTop sets top alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignTop()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignTop, align.Vertical())
	})

	t.Run("AlignMiddle sets middle alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignMiddle()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignMiddle, align.Vertical())
	})

	t.Run("AlignBottom sets bottom alignment", func(t *testing.T) {
		box := NewBox("Hi").AlignBottom()
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignBottom, align.Vertical())
	})

	t.Run("Align sets both horizontal and vertical", func(t *testing.T) {
		box := NewBox("Hi").Align(value.AlignRight, value.AlignMiddle)
		align := box.domain.Alignment()
		assert.Equal(t, value.AlignRight, align.Horizontal())
		assert.Equal(t, value.AlignMiddle, align.Vertical())
	})
}

// TestBox_FluentAPI tests fluent API chaining.
func TestBox_FluentAPI(t *testing.T) {
	t.Run("chain all methods", func(t *testing.T) {
		box := NewBox("Hello").
			Width(20).
			Height(10).
			PaddingAll(1).
			Border().
			MarginAll(2).
			AlignCenter()

		assert.NotNil(t, box)
		assert.Equal(t, "Hello", box.domain.Content())
		assert.Equal(t, 20, box.domain.Size().Width())
		assert.Equal(t, 1, box.domain.Padding().Top())
		assert.True(t, box.domain.HasBorder())
		assert.Equal(t, 2, box.domain.Margin().Top())
		assert.Equal(t, value.AlignCenter, box.domain.Alignment().Horizontal())
	})

	t.Run("methods return same box instance for chaining", func(t *testing.T) {
		box1 := NewBox("Test")
		box2 := box1.Width(10)
		box3 := box2.PaddingAll(1)

		// All should be the same Box instance (wrapping different domain boxes)
		assert.Same(t, box1, box2)
		assert.Same(t, box2, box3)
	})
}

// TestBox_Render tests rendering output.
func TestBox_Render(t *testing.T) {
	t.Run("render simple content", func(t *testing.T) {
		output := NewBox("Hello").Render()
		assert.Equal(t, "Hello", output)
	})

	t.Run("render with padding", func(t *testing.T) {
		output := NewBox("Hi").PaddingAll(1).Render()
		expected := strings.Join([]string{
			"    ",
			" Hi ",
			"    ",
		}, "\n")
		assert.Equal(t, expected, output)
	})

	t.Run("render with border", func(t *testing.T) {
		output := NewBox("Hi").Border().Render()
		expected := strings.Join([]string{
			"┌────┐",
			"│ Hi │",
			"└────┘",
		}, "\n")
		assert.Equal(t, expected, output)
	})

	t.Run("render with border and padding", func(t *testing.T) {
		output := NewBox("X").
			PaddingAll(1).
			Border().
			Render()

		// Content: 1, explicit padding: 1+1, implicit padding (border): 1+1, border: 1+1
		lines := strings.Split(output, "\n")
		assert.Equal(t, 5, len(lines)) // Border + implicit padding + explicit padding + content
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "X")
		assert.Contains(t, output, "└")
	})

	t.Run("render multi-line content", func(t *testing.T) {
		output := NewBox("A\nB\nC").Render()
		expected := "A\nB\nC"
		assert.Equal(t, expected, output)
	})

	t.Run("render with margin", func(t *testing.T) {
		output := NewBox("X").MarginAll(1).Render()

		lines := strings.Split(output, "\n")
		assert.Equal(t, 3, len(lines)) // 1 top margin + 1 content + 1 bottom margin

		// First and last lines should be margin spaces (width = content width + left + right margin)
		// Content "X" = 1 char + 1 left margin + 1 right margin = 3 spaces
		assert.Equal(t, "   ", lines[0]) // Top margin line
		assert.Equal(t, " X ", lines[1]) // Content with left/right margin
		assert.Equal(t, "   ", lines[2]) // Bottom margin line
	})
}

// TestBox_String tests String method.
func TestBox_String(t *testing.T) {
	t.Run("String returns same as Render", func(t *testing.T) {
		box := NewBox("Test").PaddingAll(1)
		assert.Equal(t, box.Render(), box.String())
	})
}

// TestBox_Layout tests layout positioning.
func TestBox_Layout(t *testing.T) {
	t.Run("layout centered box", func(t *testing.T) {
		box := NewBox("Hello").AlignCenter()
		pos := box.Layout(80, 24)

		// "Hello" = 5 chars, centered in 80: (80-5)/2 = 37
		assert.Equal(t, 37, pos.X())
		// 1 line, centered in 24: (24-1)/2 = 11
		assert.Equal(t, 11, pos.Y())
	})

	t.Run("layout left-aligned box", func(t *testing.T) {
		box := NewBox("Test").AlignLeft()
		pos := box.Layout(80, 24)

		assert.Equal(t, 0, pos.X())
		assert.Equal(t, 0, pos.Y())
	})

	t.Run("layout right-aligned box", func(t *testing.T) {
		box := NewBox("Hi").AlignRight()
		pos := box.Layout(80, 24)

		// "Hi" = 2 chars, right in 80: 80 - 2 = 78
		assert.Equal(t, 78, pos.X())
		assert.Equal(t, 0, pos.Y())
	})

	t.Run("layout box with border", func(t *testing.T) {
		box := NewBox("Hi").Border().AlignCenter()
		pos := box.Layout(20, 10)

		// "Hi" = 2 + 2 (implicit pad) + 2 (border) = 6 wide, 5 tall
		// Centered: (20-6)/2 = 7, (10-5)/2 = 2
		assert.Equal(t, 7, pos.X())
		assert.Equal(t, 2, pos.Y())
	})
}

// TestBox_Domain tests Domain accessor.
func TestBox_Domain(t *testing.T) {
	t.Run("Domain returns underlying model", func(t *testing.T) {
		box := NewBox("Test")
		domain := box.Domain()

		assert.NotNil(t, domain)
		assert.Equal(t, "Test", domain.Content())
	})

	t.Run("Domain can be modified directly", func(t *testing.T) {
		box := NewBox("Original")
		domain := box.Domain()

		// Advanced: modify domain directly
		newDomain := domain.WithContent("Modified")
		assert.Equal(t, "Modified", newDomain.Content())
		assert.Equal(t, "Original", box.domain.Content()) // Original unchanged (immutable)
	})
}

// TestBox_RealWorldExamples tests real-world usage patterns.
func TestBox_RealWorldExamples(t *testing.T) {
	t.Run("dialog box", func(t *testing.T) {
		output := NewBox("Are you sure?").
			PaddingVH(1, 2).
			Border().
			AlignCenter().
			Render()

		require.NotEmpty(t, output)
		assert.Contains(t, output, "Are you sure?")
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")
	})

	t.Run("status bar", func(t *testing.T) {
		output := NewBox("Ready | Line 42").
			PaddingVH(0, 1).
			Render()

		assert.Equal(t, " Ready | Line 42 ", output)
	})

	t.Run("title bar", func(t *testing.T) {
		box := NewBox("Phoenix TUI v1.0").
			PaddingVH(0, 2).
			Align(value.AlignCenter, value.AlignTop) // Center horizontally, top vertically

		pos := box.Layout(80, 24)
		assert.True(t, pos.X() > 20 && pos.X() < 40) // Roughly centered horizontally
		assert.Equal(t, 0, pos.Y())                  // Top
	})

	t.Run("menu item", func(t *testing.T) {
		output := NewBox("File\nEdit\nView\nHelp").
			PaddingVH(0, 1).
			Border().
			Render()

		assert.Contains(t, output, "File")
		assert.Contains(t, output, "Edit")
		assert.Contains(t, output, "View")
		assert.Contains(t, output, "Help")
		assert.Contains(t, output, "┌")
	})

	t.Run("card with all layers", func(t *testing.T) {
		output := NewBox("Welcome!").
			PaddingAll(2).
			Border().
			MarginAll(1).
			Render()

		lines := strings.Split(output, "\n")
		assert.True(t, len(lines) >= 5) // Margin + border + padding + content

		// First line should be blank (top margin)
		assert.True(t, strings.TrimSpace(lines[0]) == "")

		// Should contain border characters
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")

		// Should contain content
		assert.Contains(t, output, "Welcome!")
	})
}

// TestBox_Unicode tests Unicode content handling.
func TestBox_Unicode(t *testing.T) {
	t.Run("CJK characters", func(t *testing.T) {
		output := NewBox("你好").Render()
		assert.Equal(t, "你好", output)
	})

	t.Run("emoji", func(t *testing.T) {
		output := NewBox("👋").Render()
		assert.Equal(t, "👋", output)
	})

	t.Run("mixed content", func(t *testing.T) {
		output := NewBox("Hello 世界 👋").Render()
		assert.Equal(t, "Hello 世界 👋", output)
	})

	t.Run("CJK with border", func(t *testing.T) {
		output := NewBox("你好").Border().Render()
		assert.Contains(t, output, "你好")
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")
	})
}
