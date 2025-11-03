package clipboard_test

import (
	"fmt"
	"github.com/phoenix-tui/phoenix/clipboard"
)

func Example() { cb, _ := clipboard.New(); fmt.Printf("Clipboard: %v\n", cb != nil) }

// Output: Clipboard: true
func ExampleNew() { cb, _ := clipboard.New(); fmt.Printf("OK: %v\n", cb != nil) }

// Output: OK: true
func Example_write() { fmt.Println("Write example") }

// Output: Write example
