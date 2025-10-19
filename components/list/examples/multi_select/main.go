// Package main provides a multi-selection todo list example.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/components/list/api"
	"github.com/phoenix-tui/phoenix/components/list/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func main() {
	// Create a todo list with multi-selection
	todos := []interface{}{
		"Buy milk",
		"Write Phoenix documentation",
		"Review pull requests",
		"Fix Unicode bugs",
		"Implement layout system",
	}
	labels := []string{
		"Buy milk",
		"Write Phoenix documentation",
		"Review pull requests",
		"Fix Unicode bugs",
		"Implement layout system",
	}

	l := list.New(todos, labels, value.SelectionModeMulti).
		Height(8)

	fmt.Println("Todo List (Multi-Select)")
	fmt.Println("Press Space to toggle, Ctrl+A to select all, Esc to clear, Enter to confirm, q to quit")
	fmt.Println()

	p := tea.New(l)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Show completed todos
	selected := l.SelectedItems()
	if len(selected) > 0 {
		fmt.Printf("\nCompleted tasks:\n")
		for _, task := range selected {
			fmt.Printf("  ✓ %v\n", task)
		}
	} else {
		fmt.Println("\nNo tasks completed")
	}
}
