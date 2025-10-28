package service

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/layout/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/layout/internal/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMeasureService tests constructor.
func TestNewMeasureService(t *testing.T) {
	ms := NewMeasureService()
	assert.NotNil(t, ms)
}

// TestMeasureContent tests raw content measurement (Unicode-aware).
func TestMeasureContent(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		expectedWidth  int
		expectedHeight int
	}{
		// Basic cases
		{
			name:           "empty string",
			content:        "",
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "single ASCII char",
			content:        "A",
			expectedWidth:  1,
			expectedHeight: 1,
		},
		{
			name:           "ASCII word",
			content:        "Hello",
			expectedWidth:  5,
			expectedHeight: 1,
		},
		{
			name:           "ASCII sentence",
			content:        "Hello, World!",
			expectedWidth:  13,
			expectedHeight: 1,
		},

		// Unicode cases (CJK)
		{
			name:           "single CJK char",
			content:        "中",
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "CJK word",
			content:        "你好",
			expectedWidth:  4,
			expectedHeight: 1,
		},
		{
			name:           "CJK sentence",
			content:        "こんにちは",
			expectedWidth:  10,
			expectedHeight: 1,
		},

		// Unicode cases (emoji)
		{
			name:           "single emoji",
			content:        "👋",
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "emoji with modifier",
			content:        "👋🏻",
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "multiple emoji",
			content:        "🎉🚀✨",
			expectedWidth:  6,
			expectedHeight: 1,
		},

		// Mixed content
		{
			name:           "ASCII and CJK",
			content:        "Hi世界",
			expectedWidth:  6,
			expectedHeight: 1,
		},
		{
			name:           "ASCII and emoji",
			content:        "Hello 👋",
			expectedWidth:  8,
			expectedHeight: 1,
		},
		{
			name:           "CJK and emoji",
			content:        "你好 🎉",
			expectedWidth:  7,
			expectedHeight: 1,
		},

		// Multi-line content
		{
			name:           "two lines ASCII",
			content:        "Hello\nWorld",
			expectedWidth:  5,
			expectedHeight: 2,
		},
		{
			name:           "three lines different widths",
			content:        "A\nBC\nDEF",
			expectedWidth:  3,
			expectedHeight: 3,
		},
		{
			name:           "multi-line with CJK",
			content:        "Hello\n世界\nTest",
			expectedWidth:  5,
			expectedHeight: 3,
		},
		{
			name:           "multi-line with emoji",
			content:        "Hi\n👋🏻\nBye",
			expectedWidth:  3,
			expectedHeight: 3,
		},

		// Edge cases
		{
			name:           "only newline",
			content:        "\n",
			expectedWidth:  0,
			expectedHeight: 2,
		},
		{
			name:           "multiple newlines",
			content:        "\n\n\n",
			expectedWidth:  0,
			expectedHeight: 4,
		},
		{
			name:           "trailing newline",
			content:        "Hello\n",
			expectedWidth:  5,
			expectedHeight: 2,
		},
		{
			name:           "very long line",
			content:        strings.Repeat("A", 100),
			expectedWidth:  100,
			expectedHeight: 1,
		},
		{
			name:           "very long line with CJK",
			content:        strings.Repeat("中", 50),
			expectedWidth:  100,
			expectedHeight: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := ms.MeasureContent(tt.content)
			assert.Equal(t, tt.expectedWidth, width,
				"width for %q: got %d, want %d", tt.content, width, tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, height,
				"height for %q: got %d, want %d", tt.content, height, tt.expectedHeight)
		})
	}
}

// TestMeasure_SimpleBox tests basic box measurement without padding/border/margin.
func TestMeasure_SimpleBox(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "ASCII text",
			content:        "Hello",
			expectedWidth:  5,
			expectedHeight: 1,
		},
		{
			name:           "CJK text",
			content:        "你好",
			expectedWidth:  4,
			expectedHeight: 1,
		},
		{
			name:           "emoji text",
			content:        "👋",
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "multi-line",
			content:        "A\nBC\nDEF",
			expectedWidth:  3,
			expectedHeight: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content)
			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_WithPadding tests box measurement with padding.
func TestMeasure_WithPadding(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		padding        value2.Spacing
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "uniform padding",
			content:        "Hi",
			padding:        value2.NewSpacingAll(1),
			expectedWidth:  4, // 2 + 1 left + 1 right
			expectedHeight: 3, // 1 + 1 top + 1 bottom
		},
		{
			name:           "vertical horizontal padding",
			content:        "Hi",
			padding:        value2.NewSpacingVH(1, 2),
			expectedWidth:  6, // 2 + 2 left + 2 right
			expectedHeight: 3, // 1 + 1 top + 1 bottom
		},
		{
			name:           "individual padding",
			content:        "Test",
			padding:        value2.NewSpacing(1, 2, 3, 4),
			expectedWidth:  10, // 4 + 4 left + 2 right
			expectedHeight: 5,  // 1 + 1 top + 3 bottom
		},
		{
			name:           "zero padding",
			content:        "Test",
			padding:        value2.NewSpacingZero(),
			expectedWidth:  4,
			expectedHeight: 1,
		},
		{
			name:           "padding with CJK",
			content:        "你好",
			padding:        value2.NewSpacingAll(1),
			expectedWidth:  6, // 4 + 1 + 1
			expectedHeight: 3, // 1 + 1 + 1
		},
		{
			name:           "padding with emoji",
			content:        "👋",
			padding:        value2.NewSpacingAll(2),
			expectedWidth:  6, // 2 + 2 + 2
			expectedHeight: 5, // 1 + 2 + 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithPadding(tt.padding)
			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_WithBorder tests box measurement with border.
func TestMeasure_WithBorder(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		hasBorder      bool
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "border enabled",
			content:        "Hi",
			hasBorder:      true,
			expectedWidth:  6, // 2 (content) + 2 (implicit padding) + 2 (border)
			expectedHeight: 5, // 1 (content) + 2 (implicit padding) + 2 (border)
		},
		{
			name:           "border disabled",
			content:        "Hi",
			hasBorder:      false,
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "border with CJK",
			content:        "你好",
			hasBorder:      true,
			expectedWidth:  8, // 4 (content) + 2 (implicit padding) + 2 (border)
			expectedHeight: 5, // 1 (content) + 2 (implicit padding) + 2 (border)
		},
		{
			name:           "border with emoji",
			content:        "👋🏻",
			hasBorder:      true,
			expectedWidth:  6, // 2 (content) + 2 (implicit padding) + 2 (border)
			expectedHeight: 5, // 1 (content) + 2 (implicit padding) + 2 (border)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithBorder(tt.hasBorder)
			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_WithMargin tests box measurement with margin.
func TestMeasure_WithMargin(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		margin         value2.Spacing
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "uniform margin",
			content:        "Hi",
			margin:         value2.NewSpacingAll(1),
			expectedWidth:  4, // 2 + 1 + 1
			expectedHeight: 3, // 1 + 1 + 1
		},
		{
			name:           "vertical horizontal margin",
			content:        "Test",
			margin:         value2.NewSpacingVH(2, 3),
			expectedWidth:  10, // 4 + 3 + 3
			expectedHeight: 5,  // 1 + 2 + 2
		},
		{
			name:           "zero margin",
			content:        "Test",
			margin:         value2.NewSpacingZero(),
			expectedWidth:  4,
			expectedHeight: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithMargin(tt.margin)
			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_ComplexBox tests box measurement with all layers.
func TestMeasure_ComplexBox(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name           string
		content        string
		padding        value2.Spacing
		hasBorder      bool
		margin         value2.Spacing
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "all layers uniform",
			content:        "Hi",
			padding:        value2.NewSpacingAll(1),
			hasBorder:      true,
			margin:         value2.NewSpacingAll(1),
			expectedWidth:  10, // 2 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
			expectedHeight: 9,  // 1 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
		},
		{
			name:           "all layers different",
			content:        "Test",
			padding:        value2.NewSpacing(1, 2, 1, 2),
			hasBorder:      true,
			margin:         value2.NewSpacing(2, 3, 2, 3),
			expectedWidth:  18, // 4 (content) + 4 (explicit pad) + 2 (implicit pad) + 2 (border) + 6 (margin)
			expectedHeight: 11, // 1 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 4 (margin)
		},
		{
			name:           "CJK with all layers",
			content:        "世界",
			padding:        value2.NewSpacingAll(1),
			hasBorder:      true,
			margin:         value2.NewSpacingAll(1),
			expectedWidth:  12, // 4 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
			expectedHeight: 9,  // 1 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
		},
		{
			name:           "emoji with all layers",
			content:        "👋🏻",
			padding:        value2.NewSpacingAll(1),
			hasBorder:      true,
			margin:         value2.NewSpacingVH(1, 2),
			expectedWidth:  12, // 2 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 4 (margin)
			expectedHeight: 9,  // 1 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
		},
		{
			name:           "multi-line with all layers",
			content:        "A\nBC\nDEF",
			padding:        value2.NewSpacingAll(1),
			hasBorder:      true,
			margin:         value2.NewSpacingAll(1),
			expectedWidth:  11, // 3 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
			expectedHeight: 11, // 3 (content) + 2 (explicit pad) + 2 (implicit pad) + 2 (border) + 2 (margin)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).
				WithPadding(tt.padding).
				WithBorder(tt.hasBorder).
				WithMargin(tt.margin)

			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_WithSizeConstraints tests size constraint application.
func TestMeasure_WithSizeConstraints(t *testing.T) {
	ms := NewMeasureService()

	tests := []struct {
		name            string
		content         string
		sizeConstraints value2.Size
		expectedWidth   int
		expectedHeight  int
	}{
		{
			name:            "exact size",
			content:         "Hello",
			sizeConstraints: value2.NewSizeExact(10, 5),
			expectedWidth:   10,
			expectedHeight:  5,
		},
		{
			name:            "min width enforced",
			content:         "Hi",
			sizeConstraints: value2.NewSize(-1, -1, 10, -1, -1, -1),
			expectedWidth:   10, // Natural 2, min 10
			expectedHeight:  1,
		},
		{
			name:            "max width enforced",
			content:         strings.Repeat("A", 100),
			sizeConstraints: value2.NewSize(-1, -1, -1, 50, -1, -1),
			expectedWidth:   50, // Natural 100, max 50
			expectedHeight:  1,
		},
		{
			name:            "min height enforced",
			content:         "Single line",
			sizeConstraints: value2.NewSize(-1, -1, -1, -1, 5, -1),
			expectedWidth:   11,
			expectedHeight:  5, // Natural 1, min 5
		},
		{
			name:            "max height enforced",
			content:         "A\nB\nC\nD\nE\nF\nG\nH\nI\nJ",
			sizeConstraints: value2.NewSize(-1, -1, -1, -1, -1, 5),
			expectedWidth:   1,
			expectedHeight:  5, // Natural 10, max 5
		},
		{
			name:            "both min and max",
			content:         "Test",
			sizeConstraints: value2.NewSize(-1, -1, 10, 20, 3, 8),
			expectedWidth:   10, // Natural 4, clamped to min 10
			expectedHeight:  3,  // Natural 1, clamped to min 3
		},
		{
			name:            "unconstrained",
			content:         "Hello World",
			sizeConstraints: value2.NewSizeUnconstrained(),
			expectedWidth:   11,
			expectedHeight:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := model.NewBox(tt.content).WithSize(tt.sizeConstraints)
			size := ms.Measure(box)

			assert.Equal(t, tt.expectedWidth, size.Width(),
				"width: got %d, want %d", size.Width(), tt.expectedWidth)
			assert.Equal(t, tt.expectedHeight, size.Height(),
				"height: got %d, want %d", size.Height(), tt.expectedHeight)
		})
	}
}

// TestMeasure_RealWorldScenarios tests realistic usage patterns.
func TestMeasure_RealWorldScenarios(t *testing.T) {
	ms := NewMeasureService()

	t.Run("dialog box", func(t *testing.T) {
		box := model.NewBox("Are you sure?").
			WithPadding(value2.NewSpacingVH(1, 2)).
			WithBorder(true).
			WithMargin(value2.NewSpacingAll(1))

		size := ms.Measure(box)
		// Content: 13
		// Explicit padding: +4 horizontal, +2 vertical
		// Implicit padding (border): +2 horizontal, +2 vertical
		// Border: +2 horizontal, +2 vertical
		// Margin: +2 horizontal, +2 vertical
		assert.Equal(t, 23, size.Width()) // 13 + 4 + 2 + 2 + 2
		assert.Equal(t, 9, size.Height()) // 1 + 2 + 2 + 2 + 2
	})

	t.Run("status bar", func(t *testing.T) {
		box := model.NewBox("Ready 🎉").
			WithPadding(value2.NewSpacing(0, 1, 0, 1))

		size := ms.Measure(box)
		// Content: "Ready " (6) + "🎉" (2) = 8
		// Padding: +2 horizontal, +0 vertical
		assert.Equal(t, 10, size.Width())
		assert.Equal(t, 1, size.Height())
	})

	t.Run("multi-line card", func(t *testing.T) {
		content := "Title\n\nDescription here\n\n[OK] [Cancel]"
		box := model.NewBox(content).
			WithPadding(value2.NewSpacingAll(1)).
			WithBorder(true)

		size := ms.Measure(box)
		// Max line width: 16 ("Description here")
		// Height: 5 lines
		// Explicit padding: +2 horizontal, +2 vertical
		// Implicit padding (border): +2 horizontal, +2 vertical
		// Border: +2 horizontal, +2 vertical
		assert.Equal(t, 22, size.Width())  // 16 + 2 + 2 + 2
		assert.Equal(t, 11, size.Height()) // 5 + 2 + 2 + 2
	})

	t.Run("Japanese menu item", func(t *testing.T) {
		box := model.NewBox("ファイル (File)"). // "File" in Japanese
							WithPadding(value2.NewSpacing(0, 2, 0, 2))

		size := ms.Measure(box)
		// Content width depends on Unicode width calculation
		// Padding: +4 horizontal (2 left + 2 right)
		require.NotNil(t, size)
		// Verify it has meaningful size (content + padding)
		assert.True(t, size.Width() > 4, "width should be greater than padding alone")
		assert.Equal(t, 1, size.Height())
	})
}

// TestMeasure_EdgeCases tests boundary conditions and edge cases.
func TestMeasure_EdgeCases(t *testing.T) {
	ms := NewMeasureService()

	t.Run("single newline character", func(t *testing.T) {
		box := model.NewBox("\n")
		size := ms.Measure(box)
		assert.Equal(t, 0, size.Width())
		assert.Equal(t, 2, size.Height())
	})

	t.Run("only spaces", func(t *testing.T) {
		box := model.NewBox("     ")
		size := ms.Measure(box)
		assert.Equal(t, 5, size.Width())
		assert.Equal(t, 1, size.Height())
	})

	t.Run("very large padding", func(t *testing.T) {
		box := model.NewBox("X").
			WithPadding(value2.NewSpacingAll(100))

		size := ms.Measure(box)
		assert.Equal(t, 201, size.Width())  // 1 + 200
		assert.Equal(t, 201, size.Height()) // 1 + 200
	})

	t.Run("zero-width characters", func(t *testing.T) {
		// Zero-width space
		box := model.NewBox("A\u200BB")
		size := ms.Measure(box)
		assert.Equal(t, 2, size.Width()) // A + B, zero-width ignored
		assert.Equal(t, 1, size.Height())
	})

	t.Run("combining characters", func(t *testing.T) {
		// e + combining acute accent = é
		box := model.NewBox("e\u0301")
		size := ms.Measure(box)
		assert.Equal(t, 1, size.Width()) // Single grapheme cluster
		assert.Equal(t, 1, size.Height())
	})
}
