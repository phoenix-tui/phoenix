# phoenix/tea Examples

This directory contains example applications demonstrating phoenix/tea's Elm Architecture implementation.

## Examples Overview

### 1. Counter (`counter/`)
**Difficulty**: Beginner
**Concepts**: Basic Model, Update, View, keyboard input

A simple increment/decrement counter demonstrating the fundamental Elm Architecture pattern.

**Features**:
- Basic state management (single integer)
- Keyboard input handling
- Simple view rendering

**Run**:
```bash
cd counter
go run main.go
```

**Controls**:
- `+` or `=` : Increment counter
- `-` or `_` : Decrement counter
- `q` or Ctrl+C : Quit

---

### 2. Todo List (`todo/`)
**Difficulty**: Intermediate
**Concepts**: Complex state, list navigation, input modes

A todo list application with add/delete/toggle functionality.

**Features**:
- Complex state (list of items)
- Multiple input modes (viewing vs editing)
- List navigation with cursor
- Item completion tracking

**Run**:
```bash
cd todo
go run main.go
```

**Controls**:
- `a` : Add new todo
- `d` : Delete selected todo
- `j` or `↓` : Move down
- `k` or `↑` : Move up
- `Space` : Toggle completion
- `q` : Quit

---

### 3. Countdown Timer (`timer/`)
**Difficulty**: Intermediate/Advanced
**Concepts**: Async commands, time-based updates, state machine

A countdown timer with start/pause/reset functionality.

**Features**:
- Asynchronous commands (Tick)
- State machine (stopped/running/paused)
- Time-based updates
- Progress bar visualization

**Run**:
```bash
cd timer
go run main.go
```

**Controls**:
- `Space` : Start/Pause timer
- `r` : Reset timer
- `+` : Add 10 seconds
- `-` : Subtract 10 seconds
- `q` : Quit

---

## Building and Running

### Run Directly
```bash
cd [example-name]
go run main.go
```

### Build Binary
```bash
cd [example-name]
go build -o example.exe
./example.exe
```

### Run All Examples (Testing)
```bash
# From tea/examples directory
for dir in */; do
  (cd "$dir" && go build && echo "✓ $dir builds successfully")
done
```

---

## Learning Path

**Recommended order**:

1. **Counter** - Learn basics:
   - Model implementation
   - Init/Update/View pattern
   - Keyboard message handling

2. **Todo** - Complex state:
   - List management
   - Input modes
   - Navigation patterns
   - State updates

3. **Timer** - Async operations:
   - Command execution
   - Time-based updates
   - State machines
   - Command chaining

---

## Key Concepts Demonstrated

### Elm Architecture (MVU)
All examples follow the Model-View-Update pattern:

```go
type Model struct {
    // Your state here
}

func (m Model) Init() tea.Cmd {
    // Initialize, optionally return a command
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Handle messages, update state, optionally return command
    return m, nil
}

func (m Model) View() string {
    // Render current state as string
    return "Your view here"
}
```

### Message Handling
Messages are the way events flow through your app:

```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
    case service.TickMsg:
        // Handle timer ticks
    }
    return m, nil
}
```

### Commands
Commands are side effects that return messages:

```go
// Built-in commands
tea.Quit()                        // Quit the application
service.Tick(1 * time.Second)     // Wait and send TickMsg

// Batch commands (parallel execution)
tea.Batch(cmd1, cmd2, cmd3)

// Sequence commands (sequential execution)
tea.Sequence(cmd1, cmd2, cmd3)
```

### Program Options
Configure your program with functional options:

```go
p := tea.New(initialModel,
    tea.WithAltScreen[Model](),      // Take over full terminal
    tea.WithMouseAllMotion[Model](), // Enable mouse tracking
    tea.WithInput[Model](reader),    // Custom input source
    tea.WithOutput[Model](writer),   // Custom output destination
)
```

---

## Common Patterns

### State Machine
See `timer/` for example of state transitions:

```go
type State int

const (
    StateStopped State = iota
    StateRunning
    StatePaused
)

type Model struct {
    state State
    // ... other fields
}
```

### Input Modes
See `todo/` for example of mode switching:

```go
type Model struct {
    addMode bool
    // ... other fields
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    if m.addMode {
        return m.handleAddMode(msg)
    }
    return m.handleNormalMode(msg)
}
```

### List Navigation
See `todo/` for example of cursor-based navigation:

```go
type Model struct {
    items  []Item
    cursor int
}

// Move cursor
case "j", "down":
    if m.cursor < len(m.items)-1 {
        m.cursor++
    }
```

---

## Next Steps

After trying these examples:

1. **Modify them** - Change behavior, add features
2. **Combine patterns** - Mix timer + todo, etc.
3. **Build your own** - Use these as templates
4. **Explore components** - Check `phoenix/components` (coming soon)

---

## Troubleshooting

### Application doesn't quit
Make sure you return `tea.Quit()` command:

```go
case "q", "ctrl+c":
    return m, tea.Quit()
```

### Terminal messed up after crash
Run `reset` command in terminal to restore normal mode.

### Async command not working
Ensure you're returning the Tick command:

```go
case service.TickMsg:
    if m.running {
        // Continue ticking
        return m, service.Tick(1 * time.Second)
    }
```

---

## Additional Resources

- [phoenix/tea README](../README.md) - Main documentation
- [phoenix/tea Architecture](../../../docs/dev/ARCHITECTURE.md) - Technical design
- [API Documentation](../api/) - Public API reference

---

*Examples created for Week 6 Day 7 - phoenix/tea v1.0*
