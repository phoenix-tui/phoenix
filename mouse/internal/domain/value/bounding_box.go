package value

// BoundingBox represents a rectangular area in terminal coordinates.
// This is a value object used for hover detection.
type BoundingBox struct {
	x      int
	y      int
	width  int
	height int
}

// NewBoundingBox creates a new BoundingBox.
func NewBoundingBox(x, y, width, height int) BoundingBox {
	// Ensure non-negative dimensions
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return BoundingBox{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

// X returns the x coordinate of the top-left corner.
func (b BoundingBox) X() int {
	return b.x
}

// Y returns the y coordinate of the top-left corner.
func (b BoundingBox) Y() int {
	return b.y
}

// Width returns the width of the bounding box.
func (b BoundingBox) Width() int {
	return b.width
}

// Height returns the height of the bounding box.
func (b BoundingBox) Height() int {
	return b.height
}

// Contains returns true if the given position is inside this bounding box.
func (b BoundingBox) Contains(pos Position) bool {
	x := pos.X()
	y := pos.Y()
	return x >= b.x && x < b.x+b.width &&
		y >= b.y && y < b.y+b.height
}

// Overlaps returns true if this bounding box overlaps with another.
func (b BoundingBox) Overlaps(other BoundingBox) bool {
	return b.x < other.x+other.width &&
		b.x+b.width > other.x &&
		b.y < other.y+other.height &&
		b.y+b.height > other.y
}

// Equals returns true if this bounding box is equal to another.
func (b BoundingBox) Equals(other BoundingBox) bool {
	return b.x == other.x &&
		b.y == other.y &&
		b.width == other.width &&
		b.height == other.height
}

// IsEmpty returns true if the bounding box has zero area.
func (b BoundingBox) IsEmpty() bool {
	return b.width == 0 || b.height == 0
}

// Area returns the area of the bounding box.
func (b BoundingBox) Area() int {
	return b.width * b.height
}
