package style_test

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/style"
)

// Example demonstrates creating and using a style.
func Example() {
	s := style.New()
	fmt.Println("Style created")
	_ = s
	// Output: Style created
}

// ExampleRGB demonstrates RGB color creation.
func ExampleRGB() {
	red := style.RGB(255, 0, 0)
	fmt.Printf("Red RGB: %v\n", red != style.Color{})
	// Output: Red RGB: true
}

// ExampleNew demonstrates the New constructor.
func ExampleNew() {
	s := style.New()
	fmt.Printf("Style exists: %v\n", s != style.Style{})
	// Output: Style exists: true
}
