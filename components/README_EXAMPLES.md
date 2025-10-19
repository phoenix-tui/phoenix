# Phoenix Component Examples

This directory contains working examples demonstrating Phoenix TUI components.

## Structure

Each example is in its own subdirectory with a `main.go` file:

```
components/
├── input/
│   └── examples/
│       ├── basic/main.go           ← Run: go run ./components/input/examples/basic
│       ├── cursor_api/main.go      ← Demonstrates cursor manipulation API
│       ├── styled/main.go          ← Shows styling patterns
│       └── validated/main.go       ← Input validation example
├── list/
│   └── examples/
│       ├── basic/main.go
│       ├── custom_render/main.go
│       ├── filtered/main.go
│       └── multi_select/main.go
├── progress/
│   └── examples/
│       ├── bar_simple/main.go
│       ├── bar_styled/main.go
│       ├── multi_progress/main.go
│       └── spinner_simple/main.go
└── ... (other components)
```

## Running Examples

Examples require the Phoenix workspace to be enabled (local development).

### From repository root:
```bash
# Run a specific example
go run ./components/input/examples/basic

# Or with full path
go run ./components/input/examples/cursor_api/main.go
```

### Building examples:
```bash
# Build single example
cd components/input/examples/basic
go build -o example.exe .

# Run it
./example.exe
```

## Why Subdirectories?

Each example is a separate `package main` program. Go requires that programs with
`package main` be in separate directories to avoid "main redeclared" errors.

This structure also:
- ✅ Keeps examples out of library code (excluded from `go vet`, `go build`, `go test` in CI)
- ✅ Makes each example independently runnable
- ✅ Follows Go ecosystem best practices (see stdlib, Kubernetes, Docker, etc.)
- ✅ Allows examples to have their own README/documentation

## CI/CD Notes

Examples are **excluded** from CI pipelines:
- `go vet $(go list ./... | grep -v "/examples")` - Linting library code only
- `go build $(go list ./... | grep -v "/examples")` - Building library code only
- `go test $(go list ./... | grep -v "/examples")` - Testing library code only

Examples are demonstration programs, not part of the Phoenix library itself.
They exist to show developers how to use Phoenix components.

## Adding New Examples

1. Create a subdirectory: `mkdir components/[component]/examples/[example-name]`
2. Create `main.go`: `package main` with `func main()`
3. Document what it demonstrates
4. Run from repository root: `go run ./components/[component]/examples/[example-name]`

Example template:
```go
package main

import (
    "fmt"
    "os"

    "github.com/phoenix-tui/phoenix/components/input/api"
    tea "github.com/phoenix-tui/phoenix/tea/api"
)

type model struct {
    input *input.Input
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "ctrl+c" {
            return m, tea.Quit
        }
    }

    var cmd tea.Cmd
    updated, cmd := m.input.Update(msg)
    m.input = updated.(*input.Input)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf("Example: %s\n\nPress Ctrl-C to quit", m.input.View())
}

func main() {
    p := tea.NewProgram(model{
        input: input.New(40).Placeholder("Type here...").Focused(true),
    })

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Documentation

For component API documentation, see:
- [Input Component](./input/README.md)
- [List Component](./list/README.md)
- [Progress Component](./progress/README.md)
- [Viewport Component](./viewport/README.md)

For Phoenix framework documentation:
- [Phoenix Documentation](../docs/)
- [API Design Guide](../docs/dev/API_DESIGN.md)
