// Package main provides a filtered file list example.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/phoenix-tui/phoenix/components/list"
	tea "github.com/phoenix-tui/phoenix/tea"
)

func main() {
	// Create a file list with custom filter.
	files := []interface{}{
		"main.go",
		"main_test.go",
		"README.md",
		"config.yaml",
		"list.go",
		"list_test.go",
		"api.go",
		"api_test.go",
	}
	labels := []string{
		"main.go",
		"main_test.go",
		"README.md",
		"config.yaml",
		"list.go",
		"list_test.go",
		"api.go",
		"api_test.go",
	}

	// Custom filter: show only Go files containing query.
	l := list.New(files, labels, list.SelectionModeSingle).
		Height(10).
		Filter(func(item interface{}, query string) bool {
			filename := item.(string)
			// Show files matching query.
			return strings.Contains(filename, query)
		})

	fmt.Println("File List (Filtered - Go files only)")
	fmt.Println("Navigate with arrows or j/k, Enter to select, q to quit")
	fmt.Println()

	p := tea.New(l)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Show selected file.
	selected := l.SelectedItems()
	if len(selected) > 0 {
		fmt.Printf("\nYou selected: %v\n", selected[0])
	}
}
