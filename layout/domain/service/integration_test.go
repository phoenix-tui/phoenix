package service

import (
	"strings"
	"testing"

	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/layout/domain/model"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_MeasureAndLayout tests measure + layout pipeline.
func TestIntegration_MeasureAndLayout(t *testing.T) {
	// Setup services
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)

	t.Run("simple box centered", func(t *testing.T) {
		box := model.NewBox("Hello").
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Measure
		size := measureService.Measure(box)
		assert.Equal(t, 5, size.Width())
		assert.Equal(t, 1, size.Height())

		// Layout
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 37, position.X()) // (80 - 5) / 2
		assert.Equal(t, 11, position.Y()) // (24 - 1) / 2
	})

	t.Run("box with padding and border", func(t *testing.T) {
		box := model.NewBox("Test").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Measure: 4 (content) + 2 (explicit padding) + 2 (implicit padding) + 2 (border) = 10 wide, 7 tall
		size := measureService.Measure(box)
		assert.Equal(t, 10, size.Width())
		assert.Equal(t, 7, size.Height())

		// Layout (centered)
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 35, position.X()) // (80 - 10) / 2
		assert.Equal(t, 8, position.Y())  // (24 - 7) / 2
	})

	t.Run("Unicode content", func(t *testing.T) {
		box := model.NewBox("你好世界"). // 4 CJK chars = 8 cells
						WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Measure
		size := measureService.Measure(box)
		assert.Equal(t, 8, size.Width()) // 4 × 2

		// Layout
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 36, position.X()) // (80 - 8) / 2
	})
}

// TestIntegration_LayoutAndRender tests layout + render pipeline.
func TestIntegration_LayoutAndRender(t *testing.T) {
	// Setup services
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)
	renderService := NewRenderService()

	t.Run("simple box", func(t *testing.T) {
		box := model.NewBox("Hello").WithBorder(true)

		// Layout (not critical for render, but verifies integration)
		parentSize := value.NewSizeExact(80, 24)
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())

		// Render
		output := renderService.Render(box)
		expected := strings.Join([]string{
			"┌───────┐",
			"│ Hello │",
			"└───────┘",
		}, "\n")
		assert.Equal(t, expected, output)
	})

	t.Run("complex box", func(t *testing.T) {
		box := model.NewBox("Test").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithMargin(value.NewSpacingAll(1))

		// Layout
		parentSize := value.NewSizeExact(80, 24)
		layoutService.Layout(box, parentSize)

		// Render
		output := renderService.Render(box)

		// Verify structure
		lines := strings.Split(output, "\n")
		assert.Equal(t, 7, len(lines)) // 1 margin top + 1 border + 1 pad + 1 content + 1 pad + 1 border + 1 margin bottom
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")
		assert.Contains(t, output, "Test")
	})
}

// TestIntegration_FullPipeline tests measure → layout → render.
func TestIntegration_FullPipeline(t *testing.T) {
	// Setup services
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)
	renderService := NewRenderService()

	t.Run("complete workflow", func(t *testing.T) {
		// Create box
		box := model.NewBox("Hello\nWorld").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Step 1: Measure: 5 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) = 11 wide, 8 tall
		size := measureService.Measure(box)
		assert.Equal(t, 11, size.Width())
		assert.Equal(t, 8, size.Height())

		// Step 2: Layout
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 34, position.X()) // (80 - 11) / 2
		assert.Equal(t, 8, position.Y())  // (24 - 8) / 2

		// Step 3: Render
		output := renderService.Render(box)
		require.NotEmpty(t, output)

		// Verify output structure
		lines := strings.Split(output, "\n")
		assert.Equal(t, 6, len(lines))
		assert.Contains(t, output, "Hello")
		assert.Contains(t, output, "World")
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")

		// Log for manual verification
		t.Logf("Rendered output:\n%s", output)
		t.Logf("Position: (%d, %d)", position.X(), position.Y())
		t.Logf("Size: %dx%d", size.Width(), size.Height())
	})

	t.Run("Unicode content workflow", func(t *testing.T) {
		box := model.NewBox("你好 👋").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(40, 10)

		// Measure (Unicode-aware): 7 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) = 13 wide
		size := measureService.Measure(box)
		assert.Equal(t, 13, size.Width())

		// Layout
		position := layoutService.Layout(box, parentSize)
		assert.Equal(t, 13, position.X()) // (40 - 13) / 2

		// Render
		output := renderService.Render(box)
		require.NotEmpty(t, output)
		assert.Contains(t, output, "你好")
		assert.Contains(t, output, "👋")

		t.Logf("Unicode output:\n%s", output)
	})

	t.Run("multi-line dialog workflow", func(t *testing.T) {
		content := "Are you sure?\n\nThis action cannot be undone.\n\n[OK] [Cancel]"
		box := model.NewBox(content).
			WithPadding(value.NewSpacingVH(1, 2)).
			WithBorder(true).
			WithMargin(value.NewSpacingAll(1)).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Measure
		size := measureService.Measure(box)
		require.True(t, size.Width() > 0)
		require.True(t, size.Height() > 0)

		// Layout
		position := layoutService.Layout(box, parentSize)
		require.True(t, position.X() >= 0)
		require.True(t, position.Y() >= 0)

		// Render
		output := renderService.Render(box)
		require.NotEmpty(t, output)

		// Verify dialog structure
		assert.Contains(t, output, "Are you sure?")
		assert.Contains(t, output, "This action cannot be undone.")
		assert.Contains(t, output, "[OK]")
		assert.Contains(t, output, "[Cancel]")
		assert.Contains(t, output, "┌")
		assert.Contains(t, output, "└")

		t.Logf("Dialog output:\n%s", output)
	})
}

// TestIntegration_NodeTree tests node tree with all services.
func TestIntegration_NodeTree(t *testing.T) {
	// Setup services
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)
	renderService := NewRenderService()

	t.Run("single node tree", func(t *testing.T) {
		box := model.NewBox("Root").
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		node := model.NewNode(box)
		parentSize := value.NewSizeExact(80, 24)

		// Layout node
		positionedNode := layoutService.LayoutNode(node, parentSize)
		require.NotNil(t, positionedNode)

		// Verify position: 4 (content) + 2 (implicit pad) + 2 (border) = 8 wide, 5 tall
		position := positionedNode.Position()
		assert.Equal(t, 36, position.X()) // (80 - 8) / 2
		assert.Equal(t, 9, position.Y())  // (24 - 5) / 2

		// Render
		output := renderService.RenderNode(positionedNode)
		expected := strings.Join([]string{
			"┌──────┐",
			"│ Root │",
			"└──────┘",
		}, "\n")
		assert.Equal(t, expected, output)
	})

	t.Run("node with children", func(t *testing.T) {
		parentBox := model.NewBox("Parent").WithBorder(true)
		child1Box := model.NewBox("Child1")
		child2Box := model.NewBox("Child2")

		child1 := model.NewNode(child1Box)
		child2 := model.NewNode(child2Box)
		root := model.NewNode(parentBox).AddChild(child1).AddChild(child2)

		parentSize := value.NewSizeExact(80, 24)

		// Layout tree
		positionedRoot := layoutService.LayoutNode(root, parentSize)
		require.NotNil(t, positionedRoot)
		require.Len(t, positionedRoot.Children(), 2)

		// Verify parent position
		assert.Equal(t, 0, positionedRoot.Position().X())
		assert.Equal(t, 0, positionedRoot.Position().Y())

		// Verify children are positioned
		child1Positioned := positionedRoot.Children()[0]
		child2Positioned := positionedRoot.Children()[1]

		assert.True(t, child1Positioned.Position().Y() >= 0)
		assert.True(t, child2Positioned.Position().Y() >= child1Positioned.Position().Y())

		// Render root (simplified for Day 3)
		output := renderService.RenderNode(positionedRoot)
		require.NotEmpty(t, output)
		assert.Contains(t, output, "Parent")
	})
}

// TestIntegration_RealWorldScenarios tests complete real-world use cases.
func TestIntegration_RealWorldScenarios(t *testing.T) {
	// Setup services
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)
	renderService := NewRenderService()

	t.Run("centered modal dialog", func(t *testing.T) {
		box := model.NewBox("Confirm Action\n\nAre you sure you want to continue?\n\n[Yes] [No]").
			WithPadding(value.NewSpacingVH(1, 3)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)

		// Full pipeline
		size := measureService.Measure(box)
		position := layoutService.Layout(box, parentSize)
		output := renderService.Render(box)

		// Verify
		require.NotEmpty(t, output)
		assert.True(t, size.Width() > 0)
		assert.True(t, position.X() > 0) // Centered
		assert.True(t, position.Y() > 0) // Centered

		t.Logf("\n=== Centered Modal Dialog ===")
		t.Logf("Size: %dx%d", size.Width(), size.Height())
		t.Logf("Position: (%d, %d)", position.X(), position.Y())
		t.Logf("Output:\n%s", output)
	})

	t.Run("status bar (bottom-left)", func(t *testing.T) {
		box := model.NewBox("Ready | Line 42 | UTF-8").
			WithPadding(value.NewSpacing(0, 1, 0, 1)).
			WithAlignment(value.NewAlignment(value.AlignLeft, value.AlignBottom))

		parentSize := value.NewSizeExact(80, 24)

		// Full pipeline
		size := measureService.Measure(box)
		position := layoutService.Layout(box, parentSize)
		output := renderService.Render(box)

		// Verify
		assert.Equal(t, 0, position.X())  // Left-aligned
		assert.Equal(t, 23, position.Y()) // Bottom (24 - 1)
		assert.Equal(t, " Ready | Line 42 | UTF-8 ", output)

		t.Logf("\n=== Status Bar (Bottom-Left) ===")
		t.Logf("Size: %dx%d", size.Width(), size.Height())
		t.Logf("Position: (%d, %d)", position.X(), position.Y())
		t.Logf("Output: %q", output)
	})

	t.Run("title bar (top-center)", func(t *testing.T) {
		box := model.NewBox("Phoenix TUI Editor v1.0").
			WithPadding(value.NewSpacing(0, 2, 0, 2)).
			WithAlignment(value.NewAlignment(value.AlignCenter, value.AlignTop))

		parentSize := value.NewSizeExact(80, 24)

		// Full pipeline
		size := measureService.Measure(box)
		position := layoutService.Layout(box, parentSize)
		output := renderService.Render(box)

		// Verify
		assert.True(t, position.X() > 20 && position.X() < 40) // Roughly centered
		assert.Equal(t, 0, position.Y())                       // Top

		t.Logf("\n=== Title Bar (Top-Center) ===")
		t.Logf("Size: %dx%d", size.Width(), size.Height())
		t.Logf("Position: (%d, %d)", position.X(), position.Y())
		t.Logf("Output: %q", output)
	})

	t.Run("sidebar menu (right-aligned)", func(t *testing.T) {
		box := model.NewBox("File\nEdit\nView\nHelp").
			WithPadding(value.NewSpacingVH(0, 1)).
			WithBorder(true).
			WithAlignment(value.NewAlignment(value.AlignRight, value.AlignTop))

		parentSize := value.NewSizeExact(80, 24)

		// Full pipeline
		size := measureService.Measure(box)
		position := layoutService.Layout(box, parentSize)
		output := renderService.Render(box)

		// Verify
		assert.True(t, position.X() > 60) // Right side
		assert.Equal(t, 0, position.Y())  // Top

		t.Logf("\n=== Sidebar Menu (Right-Aligned) ===")
		t.Logf("Size: %dx%d", size.Width(), size.Height())
		t.Logf("Position: (%d, %d)", position.X(), position.Y())
		t.Logf("Output:\n%s", output)
	})
}

// TestIntegration_Performance tests basic performance characteristics.
func TestIntegration_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance tests in short mode")
	}

	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	layoutService := NewLayoutService(measureService)
	renderService := NewRenderService()

	t.Run("measure 1000 boxes", func(t *testing.T) {
		box := model.NewBox("Test content").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true)

		for i := 0; i < 1000; i++ {
			_ = measureService.Measure(box)
		}
	})

	t.Run("layout 1000 boxes", func(t *testing.T) {
		box := model.NewBox("Test").
			WithAlignment(value.NewAlignmentCenter())
		parentSize := value.NewSizeExact(80, 24)

		for i := 0; i < 1000; i++ {
			_ = layoutService.Layout(box, parentSize)
		}
	})

	t.Run("render 1000 boxes", func(t *testing.T) {
		box := model.NewBox("Test").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true)

		for i := 0; i < 1000; i++ {
			_ = renderService.Render(box)
		}
	})

	t.Run("full pipeline 1000 times", func(t *testing.T) {
		box := model.NewBox("Performance test").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())
		parentSize := value.NewSizeExact(80, 24)

		for i := 0; i < 1000; i++ {
			_ = measureService.Measure(box)
			_ = layoutService.Layout(box, parentSize)
			_ = renderService.Render(box)
		}
	})
}
