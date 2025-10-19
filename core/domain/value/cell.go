package value

// Cell represents a single terminal cell with content and visual width.
// This is an immutable value object.
//
// Content is a grapheme cluster (not just a rune) - may be:
//   - Single ASCII character: "a"
//   - Single emoji: "👋"
//   - Emoji with modifier: "👋🏻"
//   - Combined character: "é" (e + combining acute)
//
// Width is visual width in terminal columns:
//   - ASCII: 1 column
//   - Emoji: 2 columns
//   - East Asian Wide: 2 columns
//   - Zero-width joiners: 0 columns
//
// Invariants:
//   - Content is valid UTF-8
//   - Width >= 0
//   - Cell is immutable after creation
type Cell struct {
	content string // Grapheme cluster (not rune!)
	width   int    // Visual width in terminal columns
}

// NewCell creates a cell from a grapheme cluster with specified width.
// This is the universal constructor - use it for:
//   - Manual width control (advanced use cases)
//   - Automatic width (when width is pre-calculated by UnicodeService)
//
// Width will be clamped to 0 if negative.
//
// Example:
//
//	cell := value.NewCell("A", 1)          // Manual width
//	width := unicodeService.StringWidth("👋")
//	cell := value.NewCell("👋", width)     // Pre-calculated width
func NewCell(content string, width int) Cell {
	if width < 0 {
		width = 0
	}
	return Cell{
		content: content,
		width:   width,
	}
}

// Content returns the grapheme cluster content.
func (c Cell) Content() string {
	return c.content
}

// Width returns visual width in terminal columns.
func (c Cell) Width() int {
	return c.width
}

// IsEmpty returns true if cell has no visible content.
func (c Cell) IsEmpty() bool {
	return c.content == "" || c.content == " " || c.width == 0
}

// Equal returns true if cells are equal (same content and width).
func (c Cell) Equal(other Cell) bool {
	return c.content == other.content && c.width == other.width
}
