# Tutorial 2: Building Components

> **Time to Complete**: 30-40 minutes
> **Difficulty Level**: Intermediate
> **Prerequisites**: Tutorial 1 complete, understanding of Elm Architecture
> **What You'll Build**: Full-featured TODO list with styled components

---

## Table of Contents

1. [What You'll Learn](#what-youll-learn)
2. [Prerequisites](#prerequisites)
3. [Project Overview](#project-overview)
4. [Step 1: Project Setup](#step-1-project-setup-2-minutes)
5. [Step 2: Understanding phoenix/style](#step-2-understanding-phoenixstyle-5-minutes)
6. [Step 3: Model with Components](#step-3-model-with-components-8-minutes)
7. [Step 4: Styling Your App](#step-4-styling-your-app-5-minutes)
8. [Step 5: Component Integration](#step-5-component-integration-10-minutes)
9. [Step 6: View Composition](#step-6-view-composition-5-minutes)
10. [Step 7: Testing Your App](#step-7-testing-your-app-3-minutes)
11. [Understanding Component Communication](#understanding-component-communication)
12. [Exercises](#exercises)
13. [Common Issues](#common-issues)
14. [Summary](#summary)

---

## What You'll Learn

By the end of this tutorial, you will:

- Use **phoenix/style** for colors, borders, and spacing
- Integrate **phoenix/components/input** for text input
- Use **phoenix/components/list** for displaying items
- **Compose components** in your Model
- Handle **component messages** and delegation
- Build a **complete TODO application** with:
  - Add new TODOs
  - Display TODO list
  - Select/toggle TODOs
  - Delete TODOs
  - Styled UI with borders and colors

---

## Prerequisites

### Knowledge Requirements

- **Tutorial 1 complete** - Understanding of Model, Init, Update, View
- **Go intermediate** - Structs, methods, slices, error handling
- **Elm Architecture** - Message flow, commands

### Software Requirements

```bash
# Verify Go version
go version  # 1.25+

# Have your editor ready
code .  # VS Code
# or your preferred editor
```

---

## Project Overview

We're building a TODO list application with:

**Features:**
- Add TODO items via text input
- Display TODOs in a styled list
- Navigate with keyboard (Tab, arrows)
- Toggle TODO completion (Space)
- Delete selected TODO (Delete key)
- Beautiful borders and colors

**Components Used:**
- `phoenix/tea` - Event loop
- `phoenix/style` - CSS-like styling
- `phoenix/components/input` - Text input field
- `phoenix/components/list` - Selectable list

**Final UI Preview:**

```
╔════════════════════════════════════════╗
║            TODO List Manager           ║
╠════════════════════════════════════════╣
║                                        ║
║  Add TODO:                             ║
║  ┌────────────────────────────────┐   ║
║  │ Buy groceries_                 │   ║
║  └────────────────────────────────┘   ║
║                                        ║
║  Your TODOs:                           ║
║  > ☐ Write Phoenix tutorial            ║
║    ☑ Learn Elm Architecture            ║
║    ☐ Build awesome TUI app             ║
║                                        ║
╠════════════════════════════════════════╣
║  Tab: Focus  Space: Toggle  Del: Delete║
║  Enter: Add  q: Quit                   ║
╚════════════════════════════════════════╝
```

---

## Step 1: Project Setup (2 minutes)

Create a new project for our TODO app:

```bash
# Create project
mkdir phoenix-todo
cd phoenix-todo

# Initialize module
go mod init example.com/todo

# Install Phoenix packages
go get github.com/phoenix-tui/phoenix/tea
go get github.com/phoenix-tui/phoenix/style
go get github.com/phoenix-tui/phoenix/components/input
go get github.com/phoenix-tui/phoenix/components/list
```

**Create `main.go`:**

```bash
touch main.go  # macOS/Linux
# or
type nul > main.go  # Windows
```

---

## Step 2: Understanding phoenix/style (5 minutes)

Phoenix provides CSS-like styling for terminal UIs. Let's explore the basics:

### Colors

```go
import "github.com/phoenix-tui/phoenix/style"

// RGB colors (24-bit TrueColor)
red := style.RGB(255, 0, 0)
blue := style.RGB(0, 0, 255)
gray := style.RGB(128, 128, 128)

// Hex colors
purple, _ := style.Hex("#9945FF")
green, _ := style.Hex("#00FF00")

// ANSI 256-color palette
orange := style.Color256(208)

// ANSI 16-color palette
brightRed := style.Color16(9)

// Pre-defined colors
white := style.White
black := style.Black
```

### Styles

```go
// Create a style
s := style.New().
    Foreground(style.RGB(255, 255, 255)).  // White text
    Background(style.RGB(0, 0, 255)).      // Blue background
    Bold(true).                             // Bold text
    Padding(style.NewPadding(1, 2, 1, 2))  // Padding (top, right, bottom, left)

// Render styled content
output := style.Render(s, "Hello, Phoenix!")
fmt.Println(output)
```

### Borders

```go
// Pre-defined borders
normalBox := style.New().Border(style.NormalBorder)    // ┌─┐ │ └─┘
roundedBox := style.New().Border(style.RoundedBorder)  // ╭─╮ │ ╰─╯
thickBox := style.New().Border(style.ThickBorder)      // ┏━┓ ┃ ┗━┛
doubleBox := style.New().Border(style.DoubleBorder)    // ╔═╗ ║ ╚═╝
asciiBox := style.New().Border(style.ASCIIBorder)      // +-+ | +-+

// Custom border
customBox := style.New().Border(
    style.NewBorder("*", "*", "*", "*", "*", "*", "*", "*"),
)

// Border with color
coloredBox := style.New().
    Border(style.RoundedBorder).
    BorderForeground(style.RGB(100, 200, 255))
```

### Spacing and Alignment

```go
// Padding (inside border)
s := style.New().
    Padding(style.NewPadding(1, 2, 1, 2))  // top, right, bottom, left

// Margin (outside border)
s := style.New().
    Margin(style.NewMargin(0, 1, 0, 1))

// Width and height
s := style.New().
    Width(40).
    Height(10)

// Alignment
s := style.New().
    Width(40).
    Align(style.NewAlignment(style.AlignCenter, style.AlignMiddle))
```

### Immutable Updates

**IMPORTANT**: Styles are immutable! You must reassign:

```go
// WRONG - style not updated!
s := style.New()
s.Foreground(style.Red)  // Returns new style, but not assigned!

// CORRECT
s := style.New()
s = s.Foreground(style.Red)  // Reassign the new style

// Or chain methods
s := style.New().
    Foreground(style.Red).
    Bold(true).
    Padding(style.NewPadding(1, 1, 1, 1))
```

---

## Step 3: Model with Components (8 minutes)

Now let's build our Model with input and list components.

Add to `main.go`:

```go
package main

import (
    "fmt"
    "os"

    "github.com/phoenix-tui/phoenix/components/input"
    "github.com/phoenix-tui/phoenix/components/list"
    "github.com/phoenix-tui/phoenix/style"
    "github.com/phoenix-tui/phoenix/tea"
)

// TodoItem represents a single TODO item.
type TodoItem struct {
    ID        int
    Text      string
    Completed bool
}

// Model represents the application state.
type Model struct {
    // Components (note: input uses value semantics, list uses pointer)
    input input.Input  // Text input for adding TODOs
    list  *list.List   // List of TODOs

    // State
    todos       []TodoItem  // All TODO items
    nextID      int         // Next TODO ID
    focusedPane string      // "input" or "list"
    windowWidth int         // Terminal width (from WindowSizeMsg)
    windowHeight int        // Terminal height
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
    return nil
}
```

**Key Points:**

- **input.Input** - Uses **value semantics** (store as value, not pointer)
- **list.List** - Uses **pointer semantics** (store as pointer)
- **focusedPane** - Tracks which component has focus (for Tab key)
- **windowWidth/Height** - For responsive layout (we'll handle WindowSizeMsg)

### Initializing Components

Add a helper function to create the initial model:

```go
// initialModel creates the starting state.
func initialModel() Model {
    // Create input component
    inputField := input.New(40).            // Width: 40 characters
        Placeholder("Enter a TODO...").     // Placeholder text
        Focused(true)                       // Start focused

    // Create list component (empty at start)
    todoList := list.NewSingleSelect(
        []interface{}{},  // No items yet
        []string{},       // No labels yet
    ).Height(10)          // Visible height: 10 items

    return Model{
        input:       *inputField,  // Dereference pointer
        list:        todoList,
        todos:       []TodoItem{},
        nextID:      1,
        focusedPane: "input",
        windowWidth: 80,   // Default
        windowHeight: 24,
    }
}
```

**Understanding Component Creation:**

```go
// Input - fluent API
inputField := input.New(40).
    Placeholder("text").
    Focused(true).
    Validator(input.NotEmpty())

// List - constructor + fluent API
values := []interface{}{item1, item2}
labels := []string{"Item 1", "Item 2"}
todoList := list.NewSingleSelect(values, labels).
    Height(10).
    ShowFilter(false)
```

---

## Step 4: Styling Your App (5 minutes)

Create styles for different UI elements:

```go
// Styles for the application.
type Styles struct {
    Title      style.Style
    Box        style.Style
    InputBox   style.Style
    HelpBar    style.Style
    TodoItem   style.Style
    TodoDone   style.Style
    Cursor     style.Style
}

// newStyles creates the application styles.
func newStyles() Styles {
    // Color scheme
    primary := style.RGB(100, 150, 255)     // Blue
    success := style.RGB(50, 200, 100)      // Green
    muted := style.Color256(240)            // Gray

    return Styles{
        Title: style.New().
            Foreground(primary).
            Bold(true).
            Align(style.NewAlignment(style.AlignCenter, style.AlignTop)),

        Box: style.New().
            Border(style.RoundedBorder).
            BorderForeground(primary).
            Padding(style.NewPadding(1, 2, 1, 2)),

        InputBox: style.New().
            Border(style.NormalBorder).
            BorderForeground(muted).
            Padding(style.NewPadding(0, 1, 0, 1)),

        HelpBar: style.New().
            Foreground(muted).
            Align(style.NewAlignment(style.AlignCenter, style.AlignTop)),

        TodoItem: style.New().
            Foreground(style.White),

        TodoDone: style.New().
            Foreground(success).
            Strikethrough(true),

        Cursor: style.New().
            Foreground(primary).
            Bold(true),
    }
}
```

**Style Usage:**

```go
// Add styles to Model
type Model struct {
    // ... existing fields ...
    styles Styles  // NEW
}

// In initialModel()
return Model{
    // ... existing fields ...
    styles: newStyles(),  // NEW
}
```

---

## Step 5: Component Integration (10 minutes)

Now implement the Update logic to handle both components:

```go
// Update handles all messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        // Terminal was resized
        m.windowWidth = msg.Width
        m.windowHeight = msg.Height
        return m, nil

    case tea.KeyMsg:
        // Handle global keys first
        switch msg.String() {
        case "q", "ctrl+c":
            // Quit
            return m, tea.Quit()

        case "tab":
            // Switch focus between input and list
            if m.focusedPane == "input" {
                m.focusedPane = "list"
                m.input = m.input.Focused(false)
            } else {
                m.focusedPane = "input"
                m.input = m.input.Focused(true)
            }
            return m, nil

        case "enter":
            // Add TODO (only if input is focused and not empty)
            if m.focusedPane == "input" && m.input.Value() != "" {
                return m.addTodo(), nil
            }
            return m, nil

        case "delete":
            // Delete selected TODO (only if list is focused)
            if m.focusedPane == "list" {
                return m.deleteTodo(), nil
            }
            return m, nil

        case " ":
            // Toggle TODO completion (only if list is focused)
            if m.focusedPane == "list" {
                return m.toggleTodo(), nil
            }
            return m, nil
        }

        // Delegate to focused component
        if m.focusedPane == "input" {
            // Update input component
            updatedInput, cmd := m.input.Update(msg)
            m.input = updatedInput
            return m, cmd
        } else {
            // Update list component
            updatedList, cmd := m.list.Update(msg)
            m.list = updatedList
            return m, cmd
        }
    }

    return m, nil
}
```

**Component Message Delegation:**

The key pattern is **delegation**:

```
User presses key
    ↓
Update receives tea.KeyMsg
    ↓
Handle global keys (quit, tab, etc.)
    ↓
If not handled, delegate to focused component
    ↓
Component handles key and returns (component, cmd)
    ↓
Update model with new component state
```

### Helper Methods

Add helper methods for TODO operations:

```go
// addTodo creates a new TODO from input value.
func (m Model) addTodo() Model {
    // Create new TODO
    newTodo := TodoItem{
        ID:        m.nextID,
        Text:      m.input.Value(),
        Completed: false,
    }

    // Add to todos
    m.todos = append(m.todos, newTodo)
    m.nextID++

    // Clear input
    m.input = m.input.SetContent("", 0)

    // Update list component
    m.list = m.rebuildList()

    return m
}

// deleteTodo removes the selected TODO.
func (m Model) deleteTodo() Model {
    selectedIdx := m.list.FocusedIndex()
    if selectedIdx < 0 || selectedIdx >= len(m.todos) {
        return m
    }

    // Remove from todos
    m.todos = append(m.todos[:selectedIdx], m.todos[selectedIdx+1:]...)

    // Update list component
    m.list = m.rebuildList()

    return m
}

// toggleTodo toggles completion status of selected TODO.
func (m Model) toggleTodo() Model {
    selectedIdx := m.list.FocusedIndex()
    if selectedIdx < 0 || selectedIdx >= len(m.todos) {
        return m
    }

    // Toggle completed
    m.todos[selectedIdx].Completed = !m.todos[selectedIdx].Completed

    // Update list component
    m.list = m.rebuildList()

    return m
}

// rebuildList recreates the list component with current TODOs.
func (m Model) rebuildList() *list.List {
    values := make([]interface{}, len(m.todos))
    labels := make([]string, len(m.todos))

    for i, todo := range m.todos {
        values[i] = todo

        // Build label with checkbox
        checkbox := "☐"
        if todo.Completed {
            checkbox = "☑"
        }
        labels[i] = fmt.Sprintf("%s %s", checkbox, todo.Text)
    }

    return list.NewSingleSelect(values, labels).Height(10)
}
```

**Why Rebuild List?**

The list component is **immutable** - we can't modify its items directly. Instead, we rebuild it whenever TODOs change. This is the functional programming approach Phoenix uses.

---

## Step 6: View Composition (5 minutes)

Now implement the View to render everything:

```go
// View renders the application.
func (m Model) View() string {
    // Title
    title := style.Render(m.styles.Title, "TODO List Manager")

    // Input section
    inputLabel := "Add TODO:"
    inputView := m.input.View()
    inputBoxed := style.Render(m.styles.InputBox, inputView)

    // List section
    listLabel := "Your TODOs:"
    listView := m.list.View()
    if listView == "" {
        listView = "(no TODOs yet)"
    }

    // Help bar
    helpText := "Tab: Focus  Space: Toggle  Del: Delete  Enter: Add  q: Quit"
    helpBar := style.Render(m.styles.HelpBar, helpText)

    // Compose sections
    content := fmt.Sprintf(
        "\n%s\n\n  %s\n  %s\n\n  %s\n%s\n",
        title,
        inputLabel,
        inputBoxed,
        listLabel,
        listView,
    )

    // Wrap in main box
    boxed := style.Render(m.styles.Box, content)

    // Add help bar at bottom
    return fmt.Sprintf("%s\n%s\n", boxed, helpBar)
}
```

**View Composition Pattern:**

```
Build each section separately
    ↓
Apply styles to each section
    ↓
Compose sections with formatting
    ↓
Wrap in container (optional)
    ↓
Return final string
```

**Advanced View Composition:**

```go
// Use strings.Builder for complex views
func (m Model) View() string {
    var b strings.Builder

    // Header
    b.WriteString(m.renderHeader())
    b.WriteString("\n\n")

    // Content
    b.WriteString(m.renderContent())
    b.WriteString("\n\n")

    // Footer
    b.WriteString(m.renderFooter())

    return b.String()
}

func (m Model) renderHeader() string {
    title := "My App"
    return style.Render(m.styles.Title, title)
}

func (m Model) renderContent() string {
    // ... build content ...
    return content
}

func (m Model) renderFooter() string {
    help := "Press q to quit"
    return style.Render(m.styles.HelpBar, help)
}
```

---

## Step 7: Testing Your App (3 minutes)

Add the main function and run:

```go
func main() {
    // Create program
    p := tea.New(
        initialModel(),
        tea.WithAltScreen[Model](),
    )

    // Run
    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Run your app:**

```bash
go run main.go
```

**Test the features:**

1. Type text in input field
2. Press Enter to add TODO
3. Press Tab to switch to list
4. Use arrow keys to navigate
5. Press Space to toggle completion
6. Press Delete to remove TODO
7. Press Tab to go back to input
8. Press q to quit

---

## Understanding Component Communication

### Message Flow

```
┌─────────────────────────────────────┐
│  User presses key (e.g., "j")      │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Phoenix creates tea.KeyMsg         │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Model.Update(msg) called           │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Check global keys (q, tab, etc.)   │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Delegate to focused component      │
│  component.Update(msg)              │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Component returns (newComponent,   │
│  cmd)                               │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Update model with newComponent     │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Return (model, cmd)                │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Phoenix calls View()               │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│  Render to terminal                 │
└─────────────────────────────────────┘
```

### Parent-Child Pattern

```go
// Parent Model
type Model struct {
    child ChildComponent
}

// Parent Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Handle parent-level keys
    if msg.String() == "parent-key" {
        // Handle at parent level
        return m, nil
    }

    // Delegate to child
    updatedChild, cmd := m.child.Update(msg)
    m.child = updatedChild
    return m, cmd
}

// Parent View
func (m Model) View() string {
    // Compose child view into parent view
    childView := m.child.View()
    return fmt.Sprintf("Parent:\n%s", childView)
}
```

---

## Exercises

### Exercise 1: Add Priority Levels

Add high/medium/low priority to TODOs with colors.

<details>
<summary>Hint</summary>

Add `Priority int` to `TodoItem` (1=high, 2=medium, 3=low).

Add key handlers `1`, `2`, `3` to set priority.

Use different colors in `rebuildList()`:
```go
color := style.Red    // High
color := style.Yellow // Medium
color := style.Green  // Low
```
</details>

<details>
<summary>Solution</summary>

```go
type TodoItem struct {
    ID        int
    Text      string
    Completed bool
    Priority  int  // 1=high, 2=medium, 3=low
}

// In Update()
case "1", "2", "3":
    if m.focusedPane == "list" {
        selectedIdx := m.list.FocusedIndex()
        if selectedIdx >= 0 && selectedIdx < len(m.todos) {
            m.todos[selectedIdx].Priority = int(msg.Rune - '0')
            m.list = m.rebuildList()
        }
    }
    return m, nil

// In rebuildList()
priorityColor := style.White
switch todo.Priority {
case 1:
    priorityColor = style.Red
case 2:
    priorityColor = style.Yellow
case 3:
    priorityColor = style.Green
}
labels[i] = style.Render(
    style.New().Foreground(priorityColor),
    fmt.Sprintf("%s %s", checkbox, todo.Text),
)
```
</details>

### Exercise 2: Add Filter

Add a filter input to show only incomplete TODOs.

<details>
<summary>Hint</summary>

Add `showCompleted bool` to Model.

Add key handler `f` to toggle filter.

Filter TODOs in `rebuildList()`:
```go
filteredTodos := []TodoItem{}
for _, todo := range m.todos {
    if m.showCompleted || !todo.Completed {
        filteredTodos = append(filteredTodos, todo)
    }
}
```
</details>

### Exercise 3: Add Search

Add a search input to filter TODOs by text.

<details>
<summary>Hint</summary>

Add `searchInput input.Input` to Model.

Add mode `searching bool` to Model.

Add key `/` to enter search mode.

Filter TODOs by search query in `rebuildList()`:
```go
query := m.searchInput.Value()
if query != "" {
    // Filter by query
    if strings.Contains(strings.ToLower(todo.Text), strings.ToLower(query)) {
        // Include this TODO
    }
}
```
</details>

---

## Common Issues

### Issue 1: "Component not updating"

**Cause:** Forgot to reassign component after Update.

**Solution:**

```go
// WRONG
m.input.Update(msg)  // Returns new component, but not assigned!

// CORRECT
updatedInput, cmd := m.input.Update(msg)
m.input = updatedInput
return m, cmd
```

### Issue 2: "List shows old data"

**Cause:** Forgot to rebuild list after changing todos.

**Solution:**

```go
// WRONG
m.todos = append(m.todos, newTodo)
return m, nil  // List still has old data!

// CORRECT
m.todos = append(m.todos, newTodo)
m.list = m.rebuildList()  // Rebuild list!
return m, nil
```

### Issue 3: "Keys do nothing in list"

**Cause:** Not delegating messages to list component.

**Solution:**

```go
// WRONG
if m.focusedPane == "list" {
    // Handled at parent level only
    return m, nil
}

// CORRECT
if m.focusedPane == "list" {
    // Delegate to list
    updatedList, cmd := m.list.Update(msg)
    m.list = updatedList
    return m, cmd
}
```

### Issue 4: "Input cursor not visible"

**Cause:** Input component not focused.

**Solution:**

```go
// Make sure to focus input
m.input = m.input.Focused(true)

// And unfocus when switching away
m.input = m.input.Focused(false)
```

### Issue 5: "Styles not applying"

**Cause:** Not calling `style.Render()`.

**Solution:**

```go
// WRONG
text := "Hello"  // No styling applied!

// CORRECT
text := style.Render(m.styles.Title, "Hello")
```

---

## Summary

Congratulations! You've built a full-featured TODO app with Phoenix components.

### What You Learned

- **phoenix/style** - CSS-like styling (colors, borders, padding)
- **phoenix/components/input** - Text input with placeholder and validation
- **phoenix/components/list** - Selectable list with navigation
- **Component composition** - Combining multiple components in one Model
- **Message delegation** - Routing messages to focused components
- **Focus management** - Tab key to switch between components
- **View composition** - Building complex UIs from simple parts
- **Immutability** - Functional updates for components and styles

### Key Patterns

1. **Component Delegation**:
   ```go
   updatedComponent, cmd := m.component.Update(msg)
   m.component = updatedComponent
   return m, cmd
   ```

2. **Focus Management**:
   ```go
   m.input = m.input.Focused(true)   // Focus
   m.input = m.input.Focused(false)  // Unfocus
   ```

3. **View Composition**:
   ```go
   part1 := style.Render(s1, content1)
   part2 := style.Render(s2, content2)
   return fmt.Sprintf("%s\n%s", part1, part2)
   ```

4. **List Rebuilding**:
   ```go
   values := []interface{}{...}
   labels := []string{...}
   m.list = list.NewSingleSelect(values, labels)
   ```

### Architecture

```
Model (parent)
  ├─ input.Input (component)
  ├─ list.List (component)
  ├─ todos []TodoItem (state)
  └─ focusedPane string (state)

Message flow:
  User input → Update → Delegate to component → Update model → View → Render
```

### Complete Code Structure

```go
// Model
type Model struct {
    input       input.Input
    list        *list.List
    todos       []TodoItem
    focusedPane string
    styles      Styles
}

// Init
func (m Model) Init() tea.Cmd {
    return nil
}

// Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // 1. Handle global keys
    // 2. Delegate to focused component
    // 3. Return updated model + cmd
}

// View
func (m Model) View() string {
    // 1. Render each component
    // 2. Apply styles
    // 3. Compose layout
    // 4. Return final string
}
```

### Next Steps

Ready for advanced topics? In **Tutorial 3: Advanced Patterns**, you'll learn:

- **Mouse events** - Click, drag, hover interactions
- **Clipboard operations** - Copy/paste text
- **Custom components** - Build your own reusable components
- **Flexbox layout** - Advanced layout system
- **Performance optimization** - Efficient rendering for large lists
- **Complex state management** - Multi-level component trees

**Continue to**: [Tutorial 3: Advanced Patterns](03-advanced-patterns.md)

---

## Additional Resources

- [Phoenix Style API](../../api/style.md)
- [Input Component Guide](../../components/input.md)
- [List Component Guide](../../components/list.md)
- [Component Examples](../../../components/)

---

*Tutorial created for Phoenix TUI Framework v0.1.0*
*Last updated: 2025-01-04*
