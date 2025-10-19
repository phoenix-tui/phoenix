# phoenix/tea

Elm Architecture implementation for Phoenix TUI Framework.

**Status**: 🟡 Week 6 Complete - v0.1.0-alpha.0 (Development)
**Coverage**: 95.7% overall, 100% domain layer
**License**: MIT (planned)
**Maturity**: Alpha - API may change before v0.1.0 release (Week 20)

---

## What is phoenix/tea?

Phoenix/tea is a modern, type-safe implementation of the Elm Architecture (Model-View-Update) for building terminal user interfaces in Go.

**⚠️ Alpha Quality**: Week 6 complete, but API may change before v0.1.0 release (Week 20). Use for experimentation and feedback.

**Features**:
- ✅ **Type-safe** - Generic constraints for compile-time safety
- ✅ **Clean API** - Natural return types, no interface{} casts
- ✅ **DDD Architecture** - Domain-driven design with clear layers
- ✅ **High Performance** - Efficient event loop, minimal allocations
- ✅ **Well-tested** - 95.7% coverage, 100% domain layer
- ✅ **Go 1.25+** - Modern Go patterns and generics

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

// Define your model
type Model struct {
    message string
}

// Initialize the model
func (m Model) Init() api.Cmd {
    return nil
}

// Handle messages and update state
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        if msg.String() == "q" {
            return m, api.Quit()
        }
    }
    return m, nil
}

// Render the view
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

## Examples

See the [examples/](examples/) directory for complete applications:

### 1. [Counter](examples/counter/) - Beginner
Simple increment/decrement counter demonstrating basic concepts.

```bash
cd examples/counter
go run main.go
```

### 2. [Todo List](examples/todo/) - Intermediate
Todo list with add/delete/toggle, demonstrating complex state management.

```bash
cd examples/todo
go run main.go
```

### 3. [Countdown Timer](examples/timer/) - Advanced
Timer with async commands, demonstrating time-based updates and state machines.

```bash
cd examples/timer
go run main.go
```

---

## API Reference

### Types

#### Messages
```go
type Msg interface{}              // Base message type

type KeyMsg struct {              // Keyboard input
    Type  KeyType
    Rune  rune
    Alt   bool
    Ctrl  bool
    Shift bool
}

type MouseMsg struct {            // Mouse events
    X, Y   int
    Button MouseButton
    Action MouseAction
    Alt, Ctrl, Shift bool
}

type WindowSizeMsg struct {       // Terminal resize
    Width  int
    Height int
}

type QuitMsg struct{}             // Application quit

type BatchMsg struct {            // Batch command results
    Messages []Msg
}

type SequenceMsg struct {         // Sequence command results
    Messages []Msg
}
```

#### Key Types
```go
const (
    KeyRune       KeyType = iota  // Regular character
    KeyEnter                      // Enter/Return
    KeyTab                        // Tab
    KeyBackspace                  // Backspace
    KeyEscape                     // Escape
    KeyUp                         // Arrow up
    KeyDown                       // Arrow down
    KeyLeft                       // Arrow left
    KeyRight                      // Arrow right
    KeyHome                       // Home
    KeyEnd                        // End
    KeyPageUp                     // Page up
    KeyPageDown                   // Page down
    KeyDelete                     // Delete
    KeyInsert                     // Insert
    KeyF1 - KeyF12                // Function keys
)
```

### Functions

#### Commands
```go
func Quit() Cmd
// Returns a command that quits the application.

func Batch(cmds ...Cmd) Cmd
// Executes multiple commands in parallel.

func Sequence(cmds ...Cmd) Cmd
// Executes multiple commands sequentially.
```

#### Built-in Commands (domain/service)
```go
import "github.com/phoenix-tui/phoenix/tea/domain/service"

func service.Tick(d time.Duration) Cmd
// Waits for duration d, then sends TickMsg.

func service.Println(msg string) Cmd
// Sends PrintlnMsg with the given message.
```

### Program

#### Creation
```go
func New[T modelConstraint[T]](m T, opts ...ProgramOption[T]) *Program[T]
// Creates a new Program with the given model and options.
```

#### Methods
```go
func (p *Program[T]) Run() error
// Runs the program synchronously until quit.

func (p *Program[T]) Start() error
// Starts the program asynchronously in a goroutine.

func (p *Program[T]) Stop()
// Stops a running program gracefully.

func (p *Program[T]) Quit()
// Signals the program to quit.

func (p *Program[T]) Send(msg Msg) error
// Sends a message to the program's event loop.

func (p *Program[T]) IsRunning() bool
// Returns true if the program is currently running.
```

#### Options
```go
func WithInput[T any](r io.Reader) ProgramOption[T]
// Sets a custom input source (default: os.Stdin).

func WithOutput[T any](w io.Writer) ProgramOption[T]
// Sets a custom output destination (default: os.Stdout).

func WithAltScreen[T any]() ProgramOption[T]
// Enables alternate screen buffer (takes over terminal).

func WithMouseAllMotion[T any]() ProgramOption[T]
// Enables mouse tracking for all motion events.
```

---

## Advanced Usage

### Asynchronous Commands

Commands can perform asynchronous operations:

```go
import "github.com/phoenix-tui/phoenix/tea/domain/service"

func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        if msg.String() == "t" {
            // Start a timer
            return m, service.Tick(1 * time.Second)
        }

    case service.TickMsg:
        // Timer fired!
        m.lastTick = msg.Time

        if m.keepTicking {
            // Continue ticking
            return m, service.Tick(1 * time.Second)
        }
    }
    return m, nil
}
```

### Batch Commands (Parallel)

Execute multiple commands in parallel:

```go
func (m Model) Init() api.Cmd {
    return api.Batch(
        loadUserData(),
        loadSettings(),
        loadHistory(),
    )
}
```

All commands run concurrently, and their results arrive as separate messages.

### Sequence Commands (Sequential)

Execute commands in order, waiting for each to complete:

```go
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    if msg == startWorkflow {
        return m, api.Sequence(
            stepOne(),
            stepTwo(),
            stepThree(),
        )
    }
    return m, nil
}
```

Results arrive as a single `SequenceMsg` with all messages in order.

### Custom Commands

Create your own commands for side effects:

```go
func fetchWeather(city string) api.Cmd {
    return func() api.Msg {
        // Perform HTTP request
        data, err := http.Get("https://api.weather.com/" + city)
        if err != nil {
            return WeatherErrorMsg{err}
        }

        // Parse response
        weather := parseWeather(data)
        return WeatherDataMsg{weather}
    }
}

// Usage
func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    case api.KeyMsg:
        if msg.String() == "f" {
            return m, fetchWeather(m.city)
        }

    case WeatherDataMsg:
        m.weather = msg.data
        return m, nil
}
```

### Input Modes

Handle different input modes with state:

```go
type Mode int

const (
    ModeNormal Mode = iota
    ModeInsert
    ModeCommand
)

type Model struct {
    mode Mode
    // ...
}

func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        switch m.mode {
        case ModeNormal:
            return m.handleNormalMode(msg)
        case ModeInsert:
            return m.handleInsertMode(msg)
        case ModeCommand:
            return m.handleCommandMode(msg)
        }
    }
    return m, nil
}
```

---

## Development Status (Week 6 Complete!)

- [x] **Day 1**: Message System (domain/model) - 100% coverage ✅
- [x] **Day 2**: Model Interface & Commands (domain/model) - 100% coverage ✅
- [x] **Day 3**: Program Core (application/program) - 100% coverage ✅
- [x] **Day 4**: Event Loop (application/program) - 100% coverage ✅
- [x] **Day 5**: Input & ANSI Parser (infrastructure) - 100% coverage ✅
- [x] **Day 6**: Public API (api/) - 95.7% coverage ✅
- [x] **Day 7**: Examples & Documentation - Complete! ✅

**Overall**: 🎉 Week 6 Complete - v0.1.0-alpha.0 (API may change)

---

## Project Structure

```
tea/
├── domain/              # Domain layer (100% coverage)
│   ├── model/          # Message types, Model interface, Cmd
│   └── service/        # Built-in commands (Quit, Tick, Println)
│
├── application/         # Application layer (100% coverage)
│   └── program/        # Program type with lifecycle management
│
├── infrastructure/      # Infrastructure layer (100% coverage)
│   ├── input/          # Input reader
│   ├── parser/         # ANSI escape sequence parser
│   └── terminal/       # Terminal operations
│
├── api/                # Public API (95.7% coverage)
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

Phoenix/tea has extensive test coverage:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# View coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Coverage by Layer**:
- Domain layer: 100%
- Application layer: 100%
- Infrastructure layer: 100%
- API layer: 95.7%
- **Overall**: 95.7%

---

## Comparison with Bubbletea

| Feature | Phoenix/tea | Bubbletea |
|---------|-------------|-----------|
| Type Safety | ✅ Generic constraints | ⚠️ interface{} casts |
| Return Types | ✅ Natural (MyModel) | ⚠️ tea.Model interface |
| DDD Architecture | ✅ Clear layers | ❌ Monolithic |
| Test Coverage | ✅ 95.7% | ⚠️ Variable |
| API Stability | ⚠️ Alpha (may change) | ⚠️ Breaking changes |
| Performance | ✅ Optimized event loop | ✅ Good |
| Mouse Support | ✅ Planned (Week 16) | ✅ Full |
| Components | 🚧 Planned (Week 11-12) | ✅ Bubbles library |
| Maturity | ⚠️ Alpha (Week 6/20) | ✅ Battle-tested |

**Migration Path**: Phoenix/tea will provide compatibility layer for Bubbletea migration (Week 17-18).
**Current Status**: Alpha development - API subject to change before v0.1.0 release (Week 20).

---

## Roadmap

### Completed
- ✅ Week 6: Elm Architecture implementation (phoenix/tea)
- ✅ Week 5: Styling system (phoenix/style)
- ✅ Week 4: Unicode foundation (phoenix/core)
- ✅ Week 3: Core primitives (phoenix/core)

### Coming Soon
- 🚧 Week 7-8: Alternative event loop implementations
- 🚧 Week 9-10: Layout system (phoenix/layout)
- 🚧 Week 11-12: Component library (phoenix/components)
- 🚧 Week 13-14: High-performance renderer (phoenix/render)
- 🚧 Week 15-16: Mouse & clipboard support
- 🚧 Week 17-18: Migration tools from Bubbletea
- 🚧 Week 19-20: Polish & v0.1.0 release

### Future (Post v0.1.0)
- Collect community feedback
- API refinements based on real-world usage
- Breaking changes as needed
- v1.0.0 (6-12 months after v0.1.0) - API stability guarantee

---

## Contributing

Phoenix/tea is part of the Phoenix TUI Framework project. Contributions welcome!

See [../../docs/dev/](../../docs/dev/) for:
- Architecture documentation
- Development roadmap
- API design principles
- Master plan

---

## License

MIT (planned)

---

## Resources

- **Examples**: [examples/](examples/) - Complete example applications
- **Architecture**: [../../docs/dev/ARCHITECTURE.md](../../docs/dev/ARCHITECTURE.md)
- **API Design**: [../../docs/dev/API_DESIGN.md](../../docs/dev/API_DESIGN.md)
- **Roadmap**: [../../docs/dev/ROADMAP.md](../../docs/dev/ROADMAP.md)
- **Phoenix Core**: [../core/](../core/) - Terminal primitives
- **Phoenix Style**: [../style/](../style/) - Styling system

---

*Built with ❤️ using Domain-Driven Design and Modern Go*
