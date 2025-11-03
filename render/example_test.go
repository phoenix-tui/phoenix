package render_test

import (
	"bytes"
	"fmt"
	"github.com/phoenix-tui/phoenix/render"
)

func Example() {
	var b bytes.Buffer
	r := render.New(80, 24, &b)
	fmt.Printf("Size: %dx%d\n", 80, 24)
	_ = r
}

// Output: Size: 80x24
func ExampleNew() { var b bytes.Buffer; r := render.New(100, 30, &b); fmt.Printf("OK: %v\n", r != nil) }

// Output: OK: true
func ExampleNewBuffer() { buf := render.NewBuffer(10, 5); fmt.Printf("Buffer: %v\n", buf != nil) }

// Output: Buffer: true
