// Package main provides a custom rendering example with structured data.
package main

import (
	"fmt"
	"log"

	"github.com/phoenix-tui/phoenix/components/list"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// Person represents a person with name and age.
type Person struct {
	Name string
	Age  int
	City string
}

func main() {
	// Create a list of people with custom rendering.
	people := []interface{}{
		Person{Name: "Alice", Age: 30, City: "New York"},
		Person{Name: "Bob", Age: 25, City: "San Francisco"},
		Person{Name: "Charlie", Age: 35, City: "Seattle"},
		Person{Name: "Diana", Age: 28, City: "Austin"},
		Person{Name: "Eve", Age: 32, City: "Boston"},
	}
	labels := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}

	l := list.New(people, labels, list.SelectionModeSingle).
		Height(10).
		ItemRenderer(func(item interface{}, _ int, selected, focused bool) string {
			p := item.(Person)

			// Custom rendering with colors and formatting.
			prefix := "  "
			if selected {
				prefix = "✓ "
			}
			if focused {
				prefix = "→ "
			}

			return fmt.Sprintf("%s%-12s | Age: %2d | %s",
				prefix, p.Name, p.Age, p.City)
		})

	fmt.Println("People Directory (Custom Rendering)")
	fmt.Println("Navigate with arrows or j/k, Space to select, Enter to confirm, q to quit")
	fmt.Println()

	p := tea.New(l)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Show selected person.
	selected := l.SelectedItems()
	if len(selected) > 0 {
		person := selected[0].(Person)
		fmt.Printf("\nYou selected: %s (age %d) from %s\n",
			person.Name, person.Age, person.City)
	}
}
