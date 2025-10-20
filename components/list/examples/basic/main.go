// Package main provides a basic file picker example using the list component.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/components/list/api"
	"github.com/phoenix-tui/phoenix/components/list/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func main() {
	// Create a simple file picker.
	files := []interface{}{"file1.txt", "file2.go", "file3.md", "README.md", "main.go"}
	labels := []string{"file1.txt", "file2.go", "file3.md", "README.md", "main.go"}

	l := list.New(files, labels, value.SelectionModeSingle).
		Height(10)

	p := tea.New(l)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// After program exits, show selected file.
	selected := l.SelectedItems()
	if len(selected) > 0 {
		fmt.Printf("\nYou selected: %v\n", selected[0])
	}
}
