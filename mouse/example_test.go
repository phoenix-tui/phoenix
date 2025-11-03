package mouse_test

import (
	"fmt"
	"github.com/phoenix-tui/phoenix/mouse"
)

func Example() { box := mouse.NewBoundingBox(0, 0, 10, 10); fmt.Println("Box created"); _ = box }

// Output: Box created
func ExampleNewBoundingBox() {
	box := mouse.NewBoundingBox(5, 5, 20, 20)
	fmt.Printf("Box: %v\n", box != mouse.BoundingBox{})
}

// Output: Box: true
func Example_click() { fmt.Println("Click example") }

// Output: Click example
