package terminal_test

import (
	"fmt"
	"github.com/phoenix-tui/phoenix/terminal"
)

func Example() { t := terminal.New(); fmt.Printf("Terminal: %v\n", t != nil) }

// Output: Terminal: true
func ExampleNew() { t := terminal.New(); fmt.Printf("OK: %v\n", t != nil) }

// Output: OK: true
func Example_cursor() { fmt.Println("Cursor example") }

// Output: Cursor example
