package tea_test

import (
	"fmt"
	"github.com/phoenix-tui/phoenix/tea"
)

type simpleModel struct{ value string }

func (m simpleModel) Init() tea.Cmd                             { return nil }
func (m simpleModel) Update(msg tea.Msg) (simpleModel, tea.Cmd) { return m, nil }
func (m simpleModel) View() string                              { return m.value }

// Example demonstrates creating a basic MVU model.
func Example() {
	m := simpleModel{value: "Hello"}
	fmt.Println(m.View())
	// Output: Hello
}

// ExampleQuit demonstrates the Quit command.
func ExampleQuit() {
	cmd := tea.Quit()
	fmt.Printf("Quit command: %v\n", cmd != nil)
	// Output: Quit command: true
}

// ExampleBatch demonstrates batching commands.
func ExampleBatch() {
	cmd1 := tea.Quit()
	cmd2 := tea.Quit()
	batched := tea.Batch(cmd1, cmd2)
	fmt.Printf("Batch exists: %v\n", batched != nil)
	// Output: Batch exists: true
}
