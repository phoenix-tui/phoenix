# phoenix/tea

Elm Architecture implementation for Phoenix TUI Framework.

[![Go Reference](https://pkg.go.dev/badge/github.com/phoenix-tui/phoenix/tea.svg)](https://pkg.go.dev/github.com/phoenix-tui/phoenix/tea)
[![CI](https://github.com/phoenix-tui/phoenix/actions/workflows/test.yml/badge.svg)](https://github.com/phoenix-tui/phoenix/actions)

---

## What is phoenix/tea?

Phoenix/tea is a modern, type-safe implementation of the Elm Architecture (Model-View-Update) for building terminal user interfaces in Go.

**Features**:
- **Type-safe** - Generic constraints for compile-time safety (no `interface{}` casts)
- **Inline & Alt-Screen rendering** - Differential inline renderer with per-line diffing, or full alternate screen
- **TTY Control** - Run external processes (vim, shells) with full terminal handoff
- **Suspend/Resume** - Pause and restore TUI state for job control
- **DDD Architecture** - Domain-driven design with clear layers
- **High Performance** - Efficient event loop, minimal allocations
- **Well-tested** - Extensive coverage across all layers
- **Go 1.25+** - Modern Go patterns and generics

---

## Quick Start

### Installation

```bash
go get github.com/phoenix-tui/phoenix/tea
```

### Hello World

```go
package main

import (
    "fmt"
    "os"
    "github.com/phoenix-tui/phoenix/tea/api"
)

type Model struct {
    message string
}

func (m Model) Init() api.Cmd { return nil }

func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        if msg.String() == "q" {
            return m, api.Quit()
        }
    }
    return m, nil
}

func (m Model) View() string {
    return fmt.Sprintf("Hello, %s!\n\nPress 'q' to quit.", m.message)
}

func main() {
    p := api.New(Model{message: "World"}, api.WithAltScreen[Model]())
    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## Architecture

Phoenix/tea implements the Elm Architecture (Model-View-Update):

```
┌─────────────────────────────────────┐
│         Your Application            │
│                                     │
│  ┌─────────┐                        │
│  │  Init   │ ─────┐                 │
│  └─────────┘      │                 │
│                   ↓                 │
│              ┌─────────┐            │
│         ┌────│  Model  │────┐       │
│         │    └─────────┘    │       │
│         │                   │       │
│    ┌────▼────┐         ┌───▼────┐  │
│    │ Update  │         │  View  │  │
│    └────┬────┘         └────────┘  │
│         │                           │
│         ↓                           │
│    ┌─────────┐                     │
│    │   Cmd   │                     │
│    └─────────┘                     │
│         │                           │
│         └──────► Msg ──────┐       │
│                            │       │
└────────────────────────────┼───────┘
                             │
                        (User Input,
                         Timer, etc.)
```

### Core Concepts

**Model** - Your application state
```go
type Model struct {
    count int
    name  string
}
```

**Init** - Initialize the model, optionally return a command
```go
func (m Model) Init() api.Cmd {
    return nil  // or return a command
}
```

**Update** - Handle messages, update state, optionally return a command
```go
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        // Handle keyboard input
    }
    return m, nil
}
```

**View** - Render current state as a string
```go
func (m Model) View() string {
    return fmt.Sprintf("Count: %d", m.count)
}
```

**Cmd** - Side effects that return messages
```go
api.Quit()                     // Quit the application
service.Tick(1 * time.Second)  // Wait and send TickMsg
```

---

## Rendering Modes

### Alt-Screen Mode

Takes over the full terminal. Best for full-screen TUI applications.

```go
p := api.New(myModel, api.WithAltScreen[Model]())
```

### Inline Mode (default)

Renders directly in the terminal scrollback. The built-in **InlineRenderer** uses ANSI cursor-up sequences to overwrite the previous frame, with per-line diffing to minimize I/O. Ideal for CLI tools that show a small live-updating view.

```go
p := api.New(myModel) // inline mode by default
```

InlineRenderer features:
- **Per-line diffing** - Only redraws lines that changed
- **Width truncation** - Prevents line wrap from corrupting cursor positioning
- **Height clipping** - Keeps output within terminal bounds
- **ANSI-preserving** - Color codes pass through without affecting width calculations
- **Thread-safe** - All methods safe for concurrent use

---

## TTY Control

Phoenix/tea can hand off terminal control to external processes (editors, shells, pagers) and restore TUI state afterwards.

### ExecProcess

Run an external command with full TTY access:

```go
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    case api.KeyMsg:
        if msg.String() == "e" {
            // Open vim — Phoenix suspends, vim takes over, then Phoenix resumes
            return m, api.ExecProcess("vim", "file.txt")
        }
    case api.ExecProcessFinishedMsg:
        // vim exited, TUI is restored
        m.status = fmt.Sprintf("Editor exited: %v", msg.Err)
}
```

### Suspend / Resume

Manually suspend and resume the TUI for job control:

```go
// Suspend: exits raw mode, restores alt screen, shows cursor
err := p.Suspend()

// ... user interacts with normal terminal ...

// Resume: re-enters raw mode, restores alt screen, hides cursor
err = p.Resume()
```

Platform support: Linux, macOS, Windows.

---

## Examples

See the [examples/](examples/) directory for complete applications:

### 1. [Counter](examples/counter/) - Beginner
Simple increment/decrement counter demonstrating basic concepts.

```bash
cd examples/counter && go run main.go
```

### 2. [Todo List](examples/todo/) - Intermediate
Todo list with add/delete/toggle, demonstrating complex state management.

```bash
cd examples/todo && go run main.go
```

### 3. [Countdown Timer](examples/timer/) - Advanced
Timer with async commands, demonstrating time-based updates and state machines.

```bash
cd examples/timer && go run main.go
```

---

## API Reference

### Messages

```go
type Msg interface{}              // Base message type

type KeyMsg struct {              // Keyboard input
    Type  KeyType
    Rune  rune
    Alt, Ctrl, Shift bool
}

type MouseMsg struct {            // Mouse events
    X, Y   int
    Button MouseButton
    Action MouseAction
}

type WindowSizeMsg struct {       // Terminal resize
    Width, Height int
}

type QuitMsg struct{}             // Application quit
type ExecProcessFinishedMsg struct { Err error }  // External process done
```

### Key Types
```go
const (
    KeyRune       KeyType = iota  // Regular character
    KeyEnter                      // Enter/Return
    KeyTab                        // Tab
    KeyBackspace                  // Backspace
    KeyEscape                     // Escape
    KeyUp, KeyDown, KeyLeft, KeyRight
    KeyHome, KeyEnd, KeyPageUp, KeyPageDown
    KeyDelete, KeyInsert
    KeyF1 ... KeyF12              // Function keys
)
```

### Commands

```go
func Quit() Cmd                   // Quit the application
func Batch(cmds ...Cmd) Cmd       // Execute commands in parallel
func Sequence(cmds ...Cmd) Cmd    // Execute commands sequentially
func ExecProcess(name string, args ...string) Cmd  // Run external process
```

### Program

```go
// Creation
func New[T modelConstraint[T]](m T, opts ...ProgramOption[T]) *Program[T]

// Lifecycle
func (p *Program[T]) Run() error      // Run synchronously until quit
func (p *Program[T]) Start() error    // Start asynchronously
func (p *Program[T]) Stop()           // Stop gracefully
func (p *Program[T]) Quit()           // Signal quit

// Communication
func (p *Program[T]) Send(msg Msg) error  // Send message to event loop

// State
func (p *Program[T]) IsRunning() bool
func (p *Program[T]) IsSuspended() bool

// Job control
func (p *Program[T]) Suspend() error  // Suspend TUI, restore terminal
func (p *Program[T]) Resume() error   // Resume TUI from suspension
```

### Options

```go
func WithAltScreen[T any]() ProgramOption[T]       // Alternate screen buffer
func WithInput[T any](r io.Reader) ProgramOption[T] // Custom input source
func WithOutput[T any](w io.Writer) ProgramOption[T] // Custom output
func WithMouseAllMotion[T any]() ProgramOption[T]   // Mouse tracking
```

---

## Advanced Usage

### Asynchronous Commands

```go
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        if msg.String() == "t" {
            return m, service.Tick(1 * time.Second)
        }
    case service.TickMsg:
        m.lastTick = msg.Time
        if m.keepTicking {
            return m, service.Tick(1 * time.Second)
        }
    }
    return m, nil
}
```

### Batch Commands (Parallel)

```go
func (m Model) Init() api.Cmd {
    return api.Batch(
        loadUserData(),
        loadSettings(),
        loadHistory(),
    )
}
```

### Sequence Commands (Sequential)

```go
return m, api.Sequence(stepOne(), stepTwo(), stepThree())
```

### Custom Commands

```go
func fetchWeather(city string) api.Cmd {
    return func() api.Msg {
        data, err := http.Get("https://api.weather.com/" + city)
        if err != nil {
            return WeatherErrorMsg{err}
        }
        return WeatherDataMsg{parseWeather(data)}
    }
}
```

---

## Project Structure

```
tea/
├── domain/              # Domain layer
│   ├── model/          # Message types, Model interface, Cmd
│   └── service/        # Built-in commands (Quit, Tick, Println)
│
├── application/         # Application layer
│   └── program/        # Program type with lifecycle management
│
├── infrastructure/      # Infrastructure layer
│   ├── input/          # Input reader (pipe-based cancellation)
│   ├── ansi/           # ANSI escape sequence parser
│   └── renderer/       # InlineRenderer (per-line diff rendering)
│
├── api/                # Public API
│   ├── tea.go          # Re-exports and wrappers
│   └── tea_test.go     # API tests
│
└── examples/           # Example applications
    ├── counter/        # Simple counter (beginner)
    ├── todo/           # Todo list (intermediate)
    └── timer/          # Countdown timer (advanced)
```

---

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/infrastructure/renderer/
```

---

## Comparison with Bubbletea

| Feature | Phoenix/tea | Bubbletea |
|---------|-------------|-----------|
| Type Safety | Generic constraints | `interface{}` casts |
| Return Types | Natural (`MyModel`) | `tea.Model` interface |
| Inline Rendering | Per-line diffing | Basic |
| TTY Control | ExecProcess + Suspend/Resume | ExecProcess only |
| DDD Architecture | Clear layers | Monolithic |
| API Stability | Semantic versioning | Breaking changes |

**Migration**: See [MIGRATION_GUIDE.md](../docs/user/MIGRATION_GUIDE.md) for migrating from Bubbletea.

---

## Resources

- **Examples**: [examples/](examples/)
- **Migration Guide**: [MIGRATION_GUIDE.md](../docs/user/MIGRATION_GUIDE.md)
- **Architecture**: [docs/dev/ARCHITECTURE.md](../docs/dev/ARCHITECTURE.md)
- **API Design**: [docs/dev/API_DESIGN.md](../docs/dev/API_DESIGN.md)
- **Phoenix Core**: [../core/](../core/)
- **Phoenix Style**: [../style/](../style/)

---

## License

MIT

---

*Built with Domain-Driven Design and Modern Go*
