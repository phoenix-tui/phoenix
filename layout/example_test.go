package layout_test

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/layout"
)

// Example demonstrates creating a simple box with padding.
// This is the most common use case for the layout package.
func Example() {
	// Create a box with content
	box := layout.NewBox("Hello")

	// Add padding (1 cell on all sides)
	box = box.PaddingAll(1)

	// Render the box
	result := box.Render()
	fmt.Printf("Box has content: %v\n", len(result) > 0)

	// Output:
	// Box has content: true
}

// Example_row demonstrates creating a horizontal row layout.
// This shows how to arrange multiple items side by side using flexbox.
func Example_row() {
	// Create a row with three items
	row := layout.Row().
		AddRaw("Left").
		AddRaw("Center").
		AddRaw("Right")

	// Render with width and height
	result := row.Render(30, 3)
	fmt.Printf("Row rendered: %v\n", len(result) > 0)

	// Output:
	// Row rendered: true
}

// Example_column demonstrates creating a vertical column layout.
// This shows how to stack items vertically.
func Example_column() {
	// Create a column with three items
	col := layout.Column().
		AddRaw("Top").
		AddRaw("Middle").
		AddRaw("Bottom")

	// Render with width and height
	result := col.Render(20, 10)
	fmt.Printf("Column rendered: %v\n", len(result) > 0)

	// Output:
	// Column rendered: true
}
