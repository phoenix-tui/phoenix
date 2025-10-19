package service

import (
	"testing"

	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/layout/domain/model"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLayoutService tests constructor validation.
func TestNewLayoutService(t *testing.T) {
	t.Run("valid measure service", func(t *testing.T) {
		us := coreService.NewUnicodeService()
		ms := NewMeasureService(us)
		ls := NewLayoutService(ms)
		assert.NotNil(t, ls)
	})

	t.Run("nil measure service panics", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLayoutService(nil)
		}, "should panic with nil measure service")
	})
}

// TestLayout_LeftAlignment tests left alignment positioning.
func TestLayout_LeftAlignment(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name          string
		content       string
		verticalAlign value.VerticalAlignment
		parentWidth   int
		parentHeight  int
		expectedX     int
		expectedY     int
	}{
		{
			name:          "top-left",
			content:       "Hello",
			verticalAlign: value.AlignTop,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     0,
			expectedY:     0,
		},
		{
			name:          "middle-left",
			content:       "Hello",
			verticalAlign: value.AlignMiddle,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     0,
			expectedY:     11, // (24 - 1) / 2
		},
		{
			name:          "bottom-left",
			content:       "Hello",
			verticalAlign: value.AlignBottom,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     0,
			expectedY:     23, // 24 - 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithAlignment(value.NewAlignment(value.AlignLeft, tt.verticalAlign))

			parentSize := value.NewSizeExact(tt.parentWidth, tt.parentHeight)
			position := ls.Layout(box, parentSize)

			assert.Equal(t, tt.expectedX, position.X(),
				"X position: got %d, want %d", position.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, position.Y(),
				"Y position: got %d, want %d", position.Y(), tt.expectedY)
		})
	}
}

// TestLayout_CenterAlignment tests center alignment positioning.
func TestLayout_CenterAlignment(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name          string
		content       string
		verticalAlign value.VerticalAlignment
		parentWidth   int
		parentHeight  int
		expectedX     int
		expectedY     int
	}{
		{
			name:          "top-center",
			content:       "Hello",
			verticalAlign: value.AlignTop,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     37, // (80 - 5) / 2
			expectedY:     0,
		},
		{
			name:          "middle-center",
			content:       "Hello",
			verticalAlign: value.AlignMiddle,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     37, // (80 - 5) / 2
			expectedY:     11, // (24 - 1) / 2
		},
		{
			name:          "bottom-center",
			content:       "Hello",
			verticalAlign: value.AlignBottom,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     37, // (80 - 5) / 2
			expectedY:     23, // 24 - 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithAlignment(value.NewAlignment(value.AlignCenter, tt.verticalAlign))

			parentSize := value.NewSizeExact(tt.parentWidth, tt.parentHeight)
			position := ls.Layout(box, parentSize)

			assert.Equal(t, tt.expectedX, position.X(),
				"X position: got %d, want %d", position.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, position.Y(),
				"Y position: got %d, want %d", position.Y(), tt.expectedY)
		})
	}
}

// TestLayout_RightAlignment tests right alignment positioning.
func TestLayout_RightAlignment(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name          string
		content       string
		verticalAlign value.VerticalAlignment
		parentWidth   int
		parentHeight  int
		expectedX     int
		expectedY     int
	}{
		{
			name:          "top-right",
			content:       "Hello",
			verticalAlign: value.AlignTop,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     75, // 80 - 5
			expectedY:     0,
		},
		{
			name:          "middle-right",
			content:       "Hello",
			verticalAlign: value.AlignMiddle,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     75, // 80 - 5
			expectedY:     11, // (24 - 1) / 2
		},
		{
			name:          "bottom-right",
			content:       "Hello",
			verticalAlign: value.AlignBottom,
			parentWidth:   80,
			parentHeight:  24,
			expectedX:     75, // 80 - 5
			expectedY:     23, // 24 - 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithAlignment(value.NewAlignment(value.AlignRight, tt.verticalAlign))

			parentSize := value.NewSizeExact(tt.parentWidth, tt.parentHeight)
			position := ls.Layout(box, parentSize)

			assert.Equal(t, tt.expectedX, position.X(),
				"X position: got %d, want %d", position.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, position.Y(),
				"Y position: got %d, want %d", position.Y(), tt.expectedY)
		})
	}
}

// TestLayout_AllAlignmentCombinations tests all 9 alignment combinations.
func TestLayout_AllAlignmentCombinations(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	parentSize := value.NewSizeExact(80, 24)
	content := "Test" // 4 chars wide, 1 line tall

	tests := []struct {
		name       string
		horizontal value.HorizontalAlignment
		vertical   value.VerticalAlignment
		expectedX  int
		expectedY  int
	}{
		{"top-left", value.AlignLeft, value.AlignTop, 0, 0},
		{"top-center", value.AlignCenter, value.AlignTop, 38, 0},
		{"top-right", value.AlignRight, value.AlignTop, 76, 0},
		{"middle-left", value.AlignLeft, value.AlignMiddle, 0, 11},
		{"middle-center", value.AlignCenter, value.AlignMiddle, 38, 11},
		{"middle-right", value.AlignRight, value.AlignMiddle, 76, 11},
		{"bottom-left", value.AlignLeft, value.AlignBottom, 0, 23},
		{"bottom-center", value.AlignCenter, value.AlignBottom, 38, 23},
		{"bottom-right", value.AlignRight, value.AlignBottom, 76, 23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(content).
				WithAlignment(value.NewAlignment(tt.horizontal, tt.vertical))

			position := ls.Layout(box, parentSize)

			assert.Equal(t, tt.expectedX, position.X(),
				"X position: got %d, want %d", position.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, position.Y(),
				"Y position: got %d, want %d", position.Y(), tt.expectedY)
		})
	}
}

// TestLayout_WithPaddingBorderMargin tests layout with box model layers.
func TestLayout_WithPaddingBorderMargin(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("padding affects size for centering", func(t *testing.T) {
		box := model.NewBox("Hi").
			WithPadding(value.NewSpacingAll(1)).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(20, 10)
		position := ls.Layout(box, parentSize)

		// Box size: 2 (content) + 2 (padding) = 4 wide, 3 tall
		// Center: (20 - 4) / 2 = 8, (10 - 3) / 2 = 3
		assert.Equal(t, 8, position.X())
		assert.Equal(t, 3, position.Y())
	})

	t.Run("border affects size for centering", func(t *testing.T) {
		box := model.NewBox("Hi").
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(20, 10)
		position := ls.Layout(box, parentSize)

		// Box size: 2 (content) + 2 (implicit padding) + 2 (border) = 6 wide, 5 tall
		// Center: (20 - 6) / 2 = 7
		assert.Equal(t, 7, position.X())
		assert.Equal(t, 2, position.Y()) // (10 - 5) / 2
	})

	t.Run("margin affects size for centering", func(t *testing.T) {
		box := model.NewBox("Hi").
			WithMargin(value.NewSpacingAll(1)).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(20, 10)
		position := ls.Layout(box, parentSize)

		// Box size: 2 (content) + 2 (margin) = 4 wide, 3 tall
		// Center: (20 - 4) / 2 = 8, (10 - 3) / 2 = 3
		assert.Equal(t, 8, position.X())
		assert.Equal(t, 3, position.Y())
	})

	t.Run("all layers combined", func(t *testing.T) {
		box := model.NewBox("X").
			WithPadding(value.NewSpacingAll(1)).
			WithBorder(true).
			WithMargin(value.NewSpacingAll(1)).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(20, 10)
		position := ls.Layout(box, parentSize)

		// Box size: 1 (content) + 2 (explicit padding) + 2 (implicit padding) + 2 (border) + 2 (margin) = 9 wide, 9 tall
		// Center: (20 - 9) / 2 = 5 (rounded down)
		assert.Equal(t, 5, position.X())
		assert.Equal(t, 0, position.Y()) // (10 - 9) / 2
	})
}

// TestLayout_WithUnicode tests layout with Unicode content.
func TestLayout_WithUnicode(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name        string
		content     string
		parentWidth int
		expectedX   int
	}{
		{
			name:        "CJK centered",
			content:     "你好",
			parentWidth: 20,
			expectedX:   8, // (20 - 4) / 2
		},
		{
			name:        "emoji centered",
			content:     "👋",
			parentWidth: 20,
			expectedX:   9, // (20 - 2) / 2
		},
		{
			name:        "mixed ASCII and CJK",
			content:     "Hi世界",
			parentWidth: 20,
			expectedX:   7, // (20 - 6) / 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithAlignment(value.NewAlignmentCenter())

			parentSize := value.NewSizeExact(tt.parentWidth, 10)
			position := ls.Layout(box, parentSize)

			assert.Equal(t, tt.expectedX, position.X(),
				"X position: got %d, want %d", position.X(), tt.expectedX)
		})
	}
}

// TestLayout_Overflow tests overflow handling (clamping).
func TestLayout_Overflow(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("box wider than parent", func(t *testing.T) {
		box := model.NewBox("Very long content here").
			WithAlignment(value.NewAlignmentCenter())

		// Parent too narrow
		parentSize := value.NewSizeExact(10, 10)
		position := ls.Layout(box, parentSize)

		// Should clamp to left edge
		assert.Equal(t, 0, position.X(), "overflow should clamp to 0")
	})

	t.Run("box taller than parent", func(t *testing.T) {
		box := model.NewBox("A\nB\nC\nD\nE\nF\nG\nH\nI\nJ\nK\nL\nM\nN\nO").
			WithAlignment(value.NewAlignmentCenter())

		// Parent too short
		parentSize := value.NewSizeExact(10, 5)
		position := ls.Layout(box, parentSize)

		// Should clamp to top edge
		assert.Equal(t, 0, position.Y(), "overflow should clamp to 0")
	})

	t.Run("right alignment overflow", func(t *testing.T) {
		box := model.NewBox("Very long content").
			WithAlignment(value.NewAlignment(value.AlignRight, value.AlignTop))

		parentSize := value.NewSizeExact(10, 10)
		position := ls.Layout(box, parentSize)

		// Should clamp to 0 (can't be negative)
		assert.Equal(t, 0, position.X())
	})

	t.Run("bottom alignment overflow", func(t *testing.T) {
		box := model.NewBox("A\nB\nC\nD\nE\nF\nG\nH\nI\nJ").
			WithAlignment(value.NewAlignment(value.AlignLeft, value.AlignBottom))

		parentSize := value.NewSizeExact(10, 5)
		position := ls.Layout(box, parentSize)

		// Should clamp to 0 (can't be negative)
		assert.Equal(t, 0, position.Y())
	})
}

// TestLayout_ZeroSizedParent tests edge case of zero-sized parent.
func TestLayout_ZeroSizedParent(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("zero width parent", func(t *testing.T) {
		box := model.NewBox("Test")
		parentSize := value.NewSizeExact(0, 10)
		position := ls.Layout(box, parentSize)

		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())
	})

	t.Run("zero height parent", func(t *testing.T) {
		box := model.NewBox("Test")
		parentSize := value.NewSizeExact(10, 0)
		position := ls.Layout(box, parentSize)

		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())
	})

	t.Run("zero sized parent", func(t *testing.T) {
		box := model.NewBox("Test")
		parentSize := value.NewSizeExact(0, 0)
		position := ls.Layout(box, parentSize)

		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())
	})
}

// TestLayoutNode_SingleNode tests node tree positioning (no children).
func TestLayoutNode_SingleNode(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("single centered node", func(t *testing.T) {
		box := model.NewBox("Hello").
			WithAlignment(value.NewAlignmentCenter())

		node := model.NewNode(box)
		parentSize := value.NewSizeExact(80, 24)

		positionedNode := ls.LayoutNode(node, parentSize)

		require.NotNil(t, positionedNode)
		position := positionedNode.Position()
		assert.Equal(t, 37, position.X()) // (80 - 5) / 2
		assert.Equal(t, 11, position.Y()) // (24 - 1) / 2
	})

	t.Run("single top-left node", func(t *testing.T) {
		box := model.NewBox("Test")
		node := model.NewNode(box)
		parentSize := value.NewSizeExact(80, 24)

		positionedNode := ls.LayoutNode(node, parentSize)

		require.NotNil(t, positionedNode)
		position := positionedNode.Position()
		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())
	})
}

// TestLayoutNode_WithChildren tests node tree with children (vertical stacking).
func TestLayoutNode_WithChildren(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("parent with two children", func(t *testing.T) {
		// Parent box
		parentBox := model.NewBox("Parent")

		// Child boxes
		child1Box := model.NewBox("Child1")
		child2Box := model.NewBox("Child2")

		// Build tree
		child1 := model.NewNode(child1Box)
		child2 := model.NewNode(child2Box)
		root := model.NewNode(parentBox).AddChild(child1).AddChild(child2)

		// Layout
		parentSize := value.NewSizeExact(80, 24)
		positionedRoot := ls.LayoutNode(root, parentSize)

		require.NotNil(t, positionedRoot)
		require.Len(t, positionedRoot.Children(), 2)

		// Parent should be at 0,0
		assert.Equal(t, 0, positionedRoot.Position().X())
		assert.Equal(t, 0, positionedRoot.Position().Y())

		// Children should be stacked vertically
		posChild1 := positionedRoot.Children()[0]
		posChild2 := positionedRoot.Children()[1]

		// Child1 at Y=0
		assert.Equal(t, 0, posChild1.Position().Y())

		// Child2 at Y=1 (after Child1 which is 1 line tall)
		assert.Equal(t, 1, posChild2.Position().Y())
	})

	t.Run("nested children", func(t *testing.T) {
		// Three-level tree
		grandchildBox := model.NewBox("Grandchild")
		childBox := model.NewBox("Child")
		parentBox := model.NewBox("Parent")

		grandchild := model.NewNode(grandchildBox)
		child := model.NewNode(childBox).AddChild(grandchild)
		root := model.NewNode(parentBox).AddChild(child)

		// Layout
		parentSize := value.NewSizeExact(80, 24)
		positionedRoot := ls.LayoutNode(root, parentSize)

		require.NotNil(t, positionedRoot)
		require.Len(t, positionedRoot.Children(), 1)
		require.Len(t, positionedRoot.Children()[0].Children(), 1)

		// Verify positions are calculated
		assert.Equal(t, 0, positionedRoot.Position().X())
		assert.Equal(t, 0, positionedRoot.Position().Y())
	})
}

// TestCalculatePosition tests low-level position calculation.
func TestCalculatePosition(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name         string
		boxWidth     int
		boxHeight    int
		parentWidth  int
		parentHeight int
		alignment    value.Alignment
		expectedX    int
		expectedY    int
	}{
		{
			name:         "small box, centered",
			boxWidth:     10,
			boxHeight:    5,
			parentWidth:  80,
			parentHeight: 24,
			alignment:    value.NewAlignmentCenter(),
			expectedX:    35, // (80 - 10) / 2
			expectedY:    9,  // (24 - 5) / 2
		},
		{
			name:         "exact fit",
			boxWidth:     80,
			boxHeight:    24,
			parentWidth:  80,
			parentHeight: 24,
			alignment:    value.NewAlignmentCenter(),
			expectedX:    0,
			expectedY:    0,
		},
		{
			name:         "overflow",
			boxWidth:     100,
			boxHeight:    30,
			parentWidth:  80,
			parentHeight: 24,
			alignment:    value.NewAlignmentCenter(),
			expectedX:    0, // Clamped
			expectedY:    0, // Clamped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxSize := value.NewSizeExact(tt.boxWidth, tt.boxHeight)
			parentSize := value.NewSizeExact(tt.parentWidth, tt.parentHeight)

			position := ls.CalculatePosition(boxSize, parentSize, tt.alignment)

			assert.Equal(t, tt.expectedX, position.X(),
				"X: got %d, want %d", position.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, position.Y(),
				"Y: got %d, want %d", position.Y(), tt.expectedY)
		})
	}
}

// TestClampPosition tests position clamping.
func TestClampPosition(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	tests := []struct {
		name         string
		posX         int
		posY         int
		boxWidth     int
		boxHeight    int
		parentWidth  int
		parentHeight int
		expectedX    int
		expectedY    int
	}{
		{
			name:         "position within bounds",
			posX:         10,
			posY:         5,
			boxWidth:     20,
			boxHeight:    10,
			parentWidth:  80,
			parentHeight: 24,
			expectedX:    10,
			expectedY:    5,
		},
		{
			name:         "negative position clamped",
			posX:         -10,
			posY:         -5,
			boxWidth:     20,
			boxHeight:    10,
			parentWidth:  80,
			parentHeight: 24,
			expectedX:    0,
			expectedY:    0,
		},
		{
			name:         "position exceeds parent",
			posX:         70,
			posY:         20,
			boxWidth:     20,
			boxHeight:    10,
			parentWidth:  80,
			parentHeight: 24,
			expectedX:    60, // 80 - 20
			expectedY:    14, // 24 - 10
		},
		{
			name:         "box larger than parent",
			posX:         10,
			posY:         10,
			boxWidth:     100,
			boxHeight:    30,
			parentWidth:  80,
			parentHeight: 24,
			expectedX:    0,
			expectedY:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position := value.NewPosition(tt.posX, tt.posY)
			boxSize := value.NewSizeExact(tt.boxWidth, tt.boxHeight)
			parentSize := value.NewSizeExact(tt.parentWidth, tt.parentHeight)

			clamped := ls.ClampPosition(position, boxSize, parentSize)

			assert.Equal(t, tt.expectedX, clamped.X(),
				"X: got %d, want %d", clamped.X(), tt.expectedX)
			assert.Equal(t, tt.expectedY, clamped.Y(),
				"Y: got %d, want %d", clamped.Y(), tt.expectedY)
		})
	}
}

// TestLayout_RealWorldScenarios tests realistic usage patterns.
func TestLayout_RealWorldScenarios(t *testing.T) {
	us := coreService.NewUnicodeService()
	ms := NewMeasureService(us)
	ls := NewLayoutService(ms)

	t.Run("centered dialog", func(t *testing.T) {
		box := model.NewBox("Are you sure?").
			WithPadding(value.NewSpacingVH(1, 2)).
			WithBorder(true).
			WithAlignment(value.NewAlignmentCenter())

		parentSize := value.NewSizeExact(80, 24)
		position := ls.Layout(box, parentSize)

		// Box width: 13 (content) + 4 (explicit padding) + 2 (implicit padding) + 2 (border) = 21
		// Center: (80 - 21) / 2 = 29 (rounded down)
		assert.Equal(t, 29, position.X())
	})

	t.Run("bottom status bar", func(t *testing.T) {
		box := model.NewBox("Ready").
			WithAlignment(value.NewAlignment(value.AlignLeft, value.AlignBottom))

		parentSize := value.NewSizeExact(80, 24)
		position := ls.Layout(box, parentSize)

		assert.Equal(t, 0, position.X())
		assert.Equal(t, 23, position.Y()) // 24 - 1
	})

	t.Run("right-aligned menu", func(t *testing.T) {
		box := model.NewBox("Menu").
			WithAlignment(value.NewAlignment(value.AlignRight, value.AlignTop))

		parentSize := value.NewSizeExact(80, 24)
		position := ls.Layout(box, parentSize)

		assert.Equal(t, 76, position.X()) // 80 - 4
		assert.Equal(t, 0, position.Y())
	})
}
