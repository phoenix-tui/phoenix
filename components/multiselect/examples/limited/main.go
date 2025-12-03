package main

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/components/multiselect"
	"github.com/phoenix-tui/phoenix/tea"
)

// Feature represents a product feature
type Feature struct {
	ID   int
	Name string
	Desc string
}

// Model wraps the multiselect component
type Model struct {
	multi *multiselect.MultiSelect[Feature]
	done  bool
}

func main() {
	features := []Feature{
		{1, "Authentication", "User login and registration"},
		{2, "API Gateway", "RESTful API management"},
		{3, "Database", "PostgreSQL with migrations"},
		{4, "Cache", "Redis caching layer"},
		{5, "Queue", "Background job processing"},
		{6, "Search", "Full-text search with Elasticsearch"},
		{7, "Storage", "S3-compatible object storage"},
		{8, "Monitoring", "Metrics and alerting"},
		{9, "Logging", "Centralized log aggregation"},
		{10, "CI/CD", "Automated deployment pipeline"},
	}

	// Create options with descriptions
	opts := make([]*multiselect.Opt[Feature], len(features))
	for i, f := range features {
		opts[i] = multiselect.Opt[Feature](f.Name, f, f.Desc)
	}

	m := Model{
		multi: multiselect.New[Feature]("Select 2-5 features for your project:").
			Options(opts...).
			WithFilterable(true).
			WithHeight(8).
			Min(2). // At least 2 features required
			Max(5), // At most 5 features allowed
		done: false,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return m.multi.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit()
		}

	case multiselect.ConfirmSelectionMsg[Feature]:
		m.done = true
		fmt.Println("\nYou selected these features:")
		for _, f := range msg.Values {
			fmt.Printf("  [%d] %s - %s\n", f.ID, f.Name, f.Desc)
		}
		fmt.Printf("\nTotal: %d features\n", len(msg.Values))
		return m, tea.Quit()
	}

	// Update the multiselect component
	newMulti, cmd := m.multi.Update(msg)
	m.multi = newMulti
	return m, cmd
}

func (m Model) View() string {
	if m.done {
		return "" // Don't show anything after selection
	}
	return m.multi.View()
}
