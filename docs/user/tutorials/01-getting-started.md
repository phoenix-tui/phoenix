# Tutorial 1: Getting Started with Phoenix TUI

> **Time to Complete**: 15-20 minutes
> **Difficulty Level**: Beginner
> **Prerequisites**: Go 1.25+ installed, basic Go knowledge
> **What You'll Build**: Interactive counter app with keyboard handling

---

## Table of Contents

1. [What You'll Learn](#what-youll-learn)
2. [Prerequisites](#prerequisites)
3. [Understanding the Elm Architecture](#understanding-the-elm-architecture)
4. [Step 1: Project Setup](#step-1-project-setup-2-minutes)
5. [Step 2: Create Your Model](#step-2-create-your-model-3-minutes)
6. [Step 3: Implement Init](#step-3-implement-init-2-minutes)
7. [Step 4: Implement Update](#step-4-implement-update-5-minutes)
8. [Step 5: Implement View](#step-5-implement-view-3-minutes)
9. [Step 6: Run Your App](#step-6-run-your-app-2-minutes)
10. [Understanding What Happens](#understanding-what-happens)
11. [Exercises](#exercises)
12. [Common Issues](#common-issues)
13. [Summary](#summary)

---

## What You'll Learn

By the end of this tutorial, you will:

- Understand the **Model-View-Update (MVU) pattern** (Elm Architecture)
- Create a Phoenix TUI application from scratch
- Handle **keyboard events** (arrow keys, letters, Ctrl+C)
- Manage **application state** through immutable updates
- Render **dynamic content** to the terminal
- Use **tea.Cmd** for commands (quit)

---

## Prerequisites

### Required Knowledge

- **Go basics**: structs, methods, interfaces
- **Terminal/command line** familiarity
- **Text editor** (VS Code, GoLand, Vim, etc.)

### Software Requirements

- **Go 1.25 or later**
  ```bash
  go version  # Should show 1.25 or higher
  ```

- **Terminal emulator** (any modern terminal)
  - Windows: Windows Terminal, PowerShell, Git Bash
  - macOS: Terminal.app, iTerm2
  - Linux: gnome-terminal, konsole, xterm

---

## Understanding the Elm Architecture

Phoenix TUI uses the **Elm Architecture** (also known as MVU - Model-View-Update pattern). This is a simple, predictable way to build interactive applications.

### The Pattern

```
┌─────────────────────────────────────────────┐
│              ELM ARCHITECTURE               │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────┐                               │
│  │  Model   │  ← Your application state    │
│  └────┬─────┘                               │
│       │                                     │
│       ↓                                     │
│  ┌──────────┐                               │
│  │   View   │  ← Render state to string    │
│  └────┬─────┘                               │
│       │                                     │
│       ↓                                     │
│  ┌──────────┐                               │
│  │   Msg    │  ← Events (keyboard, timer)  │
│  └────┬─────┘                               │
│       │                                     │
│       ↓                                     │
│  ┌──────────┐                               │
│  │  Update  │  ← Handle msg, update state  │
│  └────┬─────┘                               │
│       │                                     │
│       ↓                                     │
│  ┌──────────┐                               │
│  │   Cmd    │  ← Side effects (optional)   │
│  └────┬─────┘                               │
│       │                                     │
│       └─────► (back to View)                │
│                                             │
└─────────────────────────────────────────────┘
```

### The Three Components

1. **Model**: Your application state (data)
   - Example: `{ count: 5 }`

2. **View**: Renders the Model to a string
   - Example: `"Count: 5\nPress + to increment"`

3. **Update**: Receives messages (events) and returns new Model
   - Example: Keyboard event `+` → new Model `{ count: 6 }`

### Why This Pattern?

- **Predictable**: Same input = same output (pure functions)
- **Testable**: Easy to test each component independently
- **Simple**: No complex state management, no lifecycle hooks
- **Functional**: Immutable updates, no side effects in Update

---

## Step 1: Project Setup (2 minutes)

Create a new Go project for your first Phoenix app:

```bash
# Create project directory
mkdir phoenix-hello
cd phoenix-hello

# Initialize Go module
go mod init example.com/hello

# Install Phoenix tea package
go get github.com/phoenix-tui/phoenix/tea
```

**Expected output:**

```
go: downloading github.com/phoenix-tui/phoenix/tea v0.1.0
go: added github.com/phoenix-tui/phoenix/tea v0.1.0
```

Create `main.go`:

```bash
# Windows
type nul > main.go

# macOS/Linux
touch main.go
```

---

## Step 2: Create Your Model (3 minutes)

Open `main.go` in your editor and define the Model:

```go
package main

import (
    "fmt"
    "os"

    "github.com/phoenix-tui/phoenix/tea"
)

// Model represents your application state.
// This holds ALL data your application needs.
type Model struct {
    count int     // The counter value
}
```

**Key Points:**

- The Model is just a plain Go struct
- It contains your application's state (data)
- In this example, we only need one field: `count`
- You can add any fields you need: strings, slices, maps, etc.

**Common Patterns:**

```go
// Simple state
type Model struct {
    count int
}

// Complex state
type Model struct {
    username   string
    todos      []string
    cursor     int
    loggedIn   bool
    lastError  error
}
```

---

## Step 3: Implement Init (2 minutes)

The `Init` method is called **once** when your program starts. It returns an optional `tea.Cmd` for initial side effects (like loading data).

Add this method to your `main.go`:

```go
// Init is called once when the program starts.
// It sets up the initial state and can return commands.
func (m Model) Init() tea.Cmd {
    // No initial commands needed for this simple app
    return nil
}
```

**Key Points:**

- `Init()` is called **exactly once** at program start
- Return `nil` if you don't need initial side effects
- Return a `tea.Cmd` if you need to run async operations (we'll cover this in Tutorial 2)

**Examples of Init Commands:**

```go
// No initial command (most simple apps)
func (m Model) Init() tea.Cmd {
    return nil
}

// Load data on startup (Tutorial 2)
func (m Model) Init() tea.Cmd {
    return LoadDataFromAPI()
}

// Start a timer (Tutorial 3)
func (m Model) Init() tea.Cmd {
    return tea.Tick(time.Second)
}
```

---

## Step 4: Implement Update (5 minutes)

The `Update` method handles **all events** in your application. It receives a message and returns:
- Updated Model
- Optional Cmd (for side effects)

Add this method to your `main.go`:

```go
// Update handles all messages (events) and updates the model.
// This is where ALL state changes happen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Type switch to handle different message types
    switch msg := msg.(type) {

    case tea.KeyMsg:
        // KeyMsg means a key was pressed
        switch msg.String() {

        case "q", "ctrl+c":
            // User pressed 'q' or Ctrl+C - quit the app
            return m, tea.Quit()

        case "up", "+", "=":
            // Increment counter
            m.count++
            return m, nil

        case "down", "-", "_":
            // Decrement counter
            if m.count > 0 {
                m.count--
            }
            return m, nil
        }
    }

    // No change - return model as-is
    return m, nil
}
```

**Key Points:**

- `Update()` is called **every time an event happens**
- Events include: keyboard, mouse, window resize, timer ticks, custom messages
- ALWAYS return a Model (either modified or unchanged)
- Return `nil` for Cmd if no side effects needed
- Return `tea.Quit()` to exit the application

**Understanding tea.KeyMsg:**

```go
// KeyMsg structure (internal)
type KeyMsg struct {
    Type  KeyType   // KeyRune, KeyEnter, KeyUp, etc.
    Rune  rune      // The actual character (for KeyRune)
    Alt   bool      // Alt key pressed?
    Ctrl  bool      // Ctrl key pressed?
    Shift bool      // Shift key pressed?
}

// String() returns human-readable key
msg.String() // Examples:
             // "a" - letter a
             // "A" - Shift+A
             // "ctrl+c" - Ctrl+C
             // "up" - Up arrow
             // "enter" - Enter key
```

**Common Key Patterns:**

```go
switch msg.String() {
case "q", "esc", "ctrl+c":
    // Quit
    return m, tea.Quit()

case "enter":
    // Confirm/Submit
    m.confirmed = true
    return m, nil

case "up", "k":
    // Vim-style navigation
    m.cursor--
    return m, nil

case "down", "j":
    m.cursor++
    return m, nil
}
```

---

## Step 5: Implement View (3 minutes)

The `View` method renders your Model to a string. Phoenix will display this string in the terminal.

Add this method to your `main.go`:

```go
// View renders the current state as a string.
// This is called after EVERY Update.
func (m Model) View() string {
    // Build the UI as a simple string
    return fmt.Sprintf(
        "╔════════════════════════════╗\n"+
        "║    Phoenix Counter Demo    ║\n"+
        "╠════════════════════════════╣\n"+
        "║                            ║\n"+
        "║    Count: %-15d ║\n"+
        "║                            ║\n"+
        "╠════════════════════════════╣\n"+
        "║  Controls:                 ║\n"+
        "║    ↑/+ : Increment         ║\n"+
        "║    ↓/- : Decrement         ║\n"+
        "║    q   : Quit              ║\n"+
        "╚════════════════════════════╝\n",
        m.count,
    )
}
```

**Key Points:**

- `View()` is called after **every Update()**
- Returns a plain string (no ANSI codes needed - Phoenix handles that)
- Can use any string building technique: `fmt.Sprintf`, `strings.Builder`, templates
- Box-drawing characters work perfectly (Unicode support is perfect!)

**View Techniques:**

```go
// Simple formatting
func (m Model) View() string {
    return fmt.Sprintf("Count: %d\nPress q to quit", m.count)
}

// Multi-line with strings.Builder
func (m Model) View() string {
    var b strings.Builder
    b.WriteString("My App\n")
    b.WriteString(fmt.Sprintf("Count: %d\n", m.count))
    b.WriteString("Press q to quit")
    return b.String()
}

// Conditional rendering
func (m Model) View() string {
    if m.count == 0 {
        return "Counter is zero! Press + to start."
    }
    return fmt.Sprintf("Count: %d", m.count)
}
```

---

## Step 6: Run Your App (2 minutes)

Now add the `main` function to tie everything together:

```go
func main() {
    // Create initial model (starting state)
    initialModel := Model{
        count: 0,
    }

    // Create program with alt screen (full-screen TUI mode)
    p := tea.New(initialModel, tea.WithAltScreen[Model]())

    // Run the program (blocks until quit)
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

**Expected Result:**

```
╔════════════════════════════╗
║    Phoenix Counter Demo    ║
╠════════════════════════════╣
║                            ║
║    Count: 0                ║
║                            ║
╠════════════════════════════╣
║  Controls:                 ║
║    ↑/+ : Increment         ║
║    ↓/- : Decrement         ║
║    q   : Quit              ║
╚════════════════════════════╝
```

**Try the controls:**

- Press `+` or `↑` - count increases
- Press `-` or `↓` - count decreases
- Press `q` - app exits cleanly

---

## Understanding What Happens

Let's trace what happens when you press the `+` key:

```
1. User presses '+'
   ↓
2. Phoenix captures keyboard event
   ↓
3. Phoenix creates tea.KeyMsg{Type: KeyRune, Rune: '+'}
   ↓
4. Phoenix calls Update(msg)
   ↓
5. Update receives KeyMsg, matches case "+"
   ↓
6. Update increments m.count (5 → 6)
   ↓
7. Update returns (Model{count: 6}, nil)
   ↓
8. Phoenix calls View()
   ↓
9. View renders "Count: 6"
   ↓
10. Phoenix displays new view in terminal
    ↓
11. Wait for next event... (back to step 1)
```

**The Event Loop:**

Phoenix runs an **infinite event loop**:

```
┌───────────────────────────────┐
│   Wait for Event (blocking)   │
└───────────┬───────────────────┘
            ↓
┌───────────────────────────────┐
│   Call Update(msg)            │
└───────────┬───────────────────┘
            ↓
┌───────────────────────────────┐
│   Call View()                 │
└───────────┬───────────────────┘
            ↓
┌───────────────────────────────┐
│   Render to Terminal          │
└───────────┬───────────────────┘
            ↓
┌───────────────────────────────┐
│   Execute Cmd (if any)        │
└───────────┬───────────────────┘
            │
            └───────► (loop back)
```

This loop continues until `tea.Quit()` is returned from `Update()`.

---

## Exercises

Ready to practice? Try these challenges:

### Exercise 1: Add a Reset Button

Add a key (e.g., `r`) that resets the counter to 0.

<details>
<summary>Hint</summary>

Add a case in Update():

```go
case "r":
    m.count = 0
    return m, nil
```

Update View() to show the new control.
</details>

<details>
<summary>Solution</summary>

```go
// In Update()
case "r":
    m.count = 0
    return m, nil

// In View() controls section
"║    r   : Reset             ║\n"+
```
</details>

### Exercise 2: Add Step Size

Allow the user to change increment/decrement size with `1`, `5`, `10` keys.

<details>
<summary>Hint</summary>

Add a `step` field to Model:

```go
type Model struct {
    count int
    step  int  // NEW
}
```

Use `m.step` in increment/decrement logic.
</details>

<details>
<summary>Solution</summary>

```go
// Model
type Model struct {
    count int
    step  int
}

// In main()
initialModel := Model{count: 0, step: 1}

// In Update()
case "1":
    m.step = 1
    return m, nil
case "5":
    m.step = 5
    return m, nil
case "0":
    m.step = 10
    return m, nil

case "up", "+", "=":
    m.count += m.step  // Use step
    return m, nil

case "down", "-", "_":
    m.count -= m.step  // Use step
    if m.count < 0 {
        m.count = 0
    }
    return m, nil

// In View()
"║    Count: %-15d ║\n"+
"║    Step:  %-15d ║\n"+
// ...
"║    1/5/0: Change step      ║\n"+
```
</details>

### Exercise 3: Add History

Track the last 5 values in a slice and display them.

<details>
<summary>Hint</summary>

Add `history []int` to Model. On each change, append old value to history (limit to 5).
</details>

<details>
<summary>Solution</summary>

```go
type Model struct {
    count   int
    history []int
}

// Helper to add to history
func (m Model) addToHistory(val int) Model {
    m.history = append(m.history, val)
    if len(m.history) > 5 {
        m.history = m.history[1:]  // Keep last 5
    }
    return m
}

// In Update() before changing count
case "up", "+", "=":
    m = m.addToHistory(m.count)
    m.count++
    return m, nil

// In View()
"║  History:                  ║\n"+
fmt.Sprintf("║    %v\n", m.history)+
```
</details>

---

## Common Issues

### Issue 1: "panic: runtime error: invalid memory address"

**Cause:** Model not initialized properly.

**Solution:**

```go
// WRONG
var p *tea.Program[Model]
p.Run()  // Panic!

// CORRECT
p := tea.New(Model{count: 0})
p.Run()
```

### Issue 2: "Program exits immediately"

**Cause:** Returning `tea.Quit()` in `Init()` or immediately in `Update()`.

**Solution:**

```go
// WRONG
func (m Model) Init() tea.Cmd {
    return tea.Quit()  // Exits immediately!
}

// CORRECT
func (m Model) Init() tea.Cmd {
    return nil
}
```

### Issue 3: "Key presses don't work"

**Cause:** Wrong key string comparison.

**Solution:**

Use `msg.String()` for human-readable keys:

```go
// WRONG
case msg.Rune == '+':  // Doesn't compile

// CORRECT
case msg.String() == "+":
```

### Issue 4: "Terminal doesn't restore after crash"

**Cause:** Program panicked before cleanup.

**Solution:**

Phoenix automatically cleans up on normal exit. For debugging crashes:

```bash
# Reset terminal manually if needed
reset

# Or on Windows
cls
```

### Issue 5: "Box-drawing characters look broken"

**Cause:** Terminal encoding not set to UTF-8.

**Solution:**

```bash
# Windows PowerShell
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Or use Windows Terminal (built-in UTF-8 support)
```

---

## Summary

Congratulations! You've built your first Phoenix TUI application. Let's recap:

### What You Learned

- **Elm Architecture (MVU)**: Model-View-Update pattern
- **Model**: Defines application state (data)
- **Init**: One-time initialization (returns Cmd)
- **Update**: Handles events, returns new Model + optional Cmd
- **View**: Renders Model to string
- **tea.KeyMsg**: Keyboard event handling
- **tea.Quit()**: Exit the application
- **Event Loop**: How Phoenix processes events continuously

### Key Concepts

1. **Immutability**: Model updates create new values (functional style)
2. **Pure Functions**: Update and View should be deterministic
3. **Single Source of Truth**: All state lives in Model
4. **Unidirectional Data Flow**: Event → Update → View → Render

### The Pattern

```
Event → Update(msg) → new Model → View() → Terminal
   ↑                                           │
   └───────────────────────────────────────────┘
```

### Full Code Example

Here's the complete `main.go` for reference:

```go
package main

import (
    "fmt"
    "os"

    "github.com/phoenix-tui/phoenix/tea"
)

type Model struct {
    count int
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit()
        case "up", "+", "=":
            m.count++
            return m, nil
        case "down", "-", "_":
            if m.count > 0 {
                m.count--
            }
            return m, nil
        }
    }
    return m, nil
}

func (m Model) View() string {
    return fmt.Sprintf(
        "╔════════════════════════════╗\n"+
        "║    Phoenix Counter Demo    ║\n"+
        "╠════════════════════════════╣\n"+
        "║                            ║\n"+
        "║    Count: %-15d ║\n"+
        "║                            ║\n"+
        "╠════════════════════════════╣\n"+
        "║  Controls:                 ║\n"+
        "║    ↑/+ : Increment         ║\n"+
        "║    ↓/- : Decrement         ║\n"+
        "║    q   : Quit              ║\n"+
        "╚════════════════════════════╝\n",
        m.count,
    )
}

func main() {
    p := tea.New(Model{count: 0}, tea.WithAltScreen[Model]())
    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Next Steps

Ready for more? In **Tutorial 2: Building Components**, you'll learn:

- Using `phoenix/style` for colors and styling
- Building text input with `phoenix/components/input`
- Creating lists with `phoenix/components/list`
- Composing components in your Model
- Building a complete TODO app

**Continue to**: [Tutorial 2: Building Components](02-building-components.md)

---

## Additional Resources

- [Phoenix API Documentation](../../api/tea.md)
- [Elm Architecture Guide](https://guide.elm-lang.org/architecture/)
- [Phoenix Examples](../../../tea/examples/)
- [Phoenix GitHub](https://github.com/phoenix-tui/phoenix)

---

*Tutorial created for Phoenix TUI Framework v0.1.0*
*Last updated: 2025-01-04*
