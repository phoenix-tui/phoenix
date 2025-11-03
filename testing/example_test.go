package testing_test

import (
	"fmt"
	test "github.com/phoenix-tui/phoenix/testing"
)

func ExampleNullTerminal() {
	var t test.NullTerminal
	fmt.Printf("Null: %v\n", t != test.NullTerminal{})
}

// Output: Null: false
func ExampleMockTerminal() { t := &test.MockTerminal{}; fmt.Printf("Mock: %v\n", t != nil) }

// Output: Mock: true
func Example() { fmt.Println("Testing helpers") }

// Output: Testing helpers
