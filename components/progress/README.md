# Progress Component

Universal progress indicators for Phoenix TUI Framework - progress bars and animated spinners.

## Features

- **Progress Bars** - Visual percentage indicators with customizable styling
- **Spinners** - Animated loading indicators with 15+ pre-defined styles
- **tea.Model Integration** - Works seamlessly with Phoenix tea event loop
- **Fluent API** - Method chaining for clean, readable code
- **Universal Design** - Works for any application (file downloads, task processing, loading indicators, etc.)

## Installation

```go
import progress "github.com/phoenix-tui/phoenix/components/progress/api"
```

## Quick Start

### Progress Bar

```go
// Create a 40-character wide progress bar
bar := progress.NewBar(40)

// Update progress
bar.SetProgress(50)

// Render
fmt.Println(bar.View())
// Output: "████████████████████░░░░░░░░░░░░░░░░░░░░"
```

### Spinner

```go
type model struct {
    spinner *progress.Spinner
}

func (m model) Init() tea.Cmd {
    m.spinner = progress.NewSpinner("dots").Label("Loading...")
    return m.spinner.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.TickMsg:
        updated, cmd := m.spinner.Update(msg)
        m.spinner = updated.(*progress.Spinner)
        return m, cmd
    }
    return m, nil
}

func (m model) View() string {
    return m.spinner.View()
    // Output: "⠋ Loading..." (animated)
}
```

## API Reference

### Progress Bar

#### Constructor

```go
// Create bar with specified width
bar := progress.NewBar(width int) *Bar

// Create bar with initial percentage
bar := progress.NewBarWithProgress(width int, percentage int) *Bar
```

#### Configuration (Fluent API)

```go
bar.FillChar(char rune) *Bar          // Set filled character (default: '█')
bar.EmptyChar(char rune) *Bar         // Set empty character (default: '░')
bar.ShowPercent(show bool) *Bar       // Toggle percentage display
bar.Label(label string) *Bar          // Set label text
```

#### Progress Updates

```go
bar.SetProgress(pct int) *Bar         // Set progress (0-100)
bar.Increment(delta int) *Bar         // Increase progress
bar.Decrement(delta int) *Bar         // Decrease progress
```

#### Accessors

```go
bar.Progress() int                    // Get current percentage
bar.IsComplete() bool                 // Check if 100%
```

#### tea.Model Interface

```go
bar.Init() tea.Cmd                    // Initialize (returns nil)
bar.Update(msg tea.Msg) (tea.Model, tea.Cmd)  // Handle messages
bar.View() string                     // Render to string
```

### Spinner

#### Constructor

```go
// Create spinner with pre-defined style
spinner := progress.NewSpinner(style string) *Spinner

// Available styles:
// "dots", "line", "arrow", "circle", "bounce",
// "dot-pulse", "grow-vertical", "grow-horizontal",
// "box-bounce", "simple-dots", "clock", "earth",
// "moon", "toggle", "hamburger"
```

#### Configuration

```go
spinner.Label(label string) *Spinner  // Set label text
```

#### tea.Model Interface

```go
spinner.Init() tea.Cmd                // Initialize animation
spinner.Update(msg tea.Msg) (tea.Model, tea.Cmd)  // Handle tick messages
spinner.View() string                 // Render current frame
```

## Examples

### Example 1: Simple Progress Bar

```go
bar := progress.NewBar(40)

for i := 0; i <= 100; i += 10 {
    bar.SetProgress(i)
    fmt.Printf("\r%s", bar.View())
    time.Sleep(200 * time.Millisecond)
}
```

### Example 2: Styled Progress Bar

```go
bar := progress.NewBar(50).
    FillChar('█').
    EmptyChar('░').
    ShowPercent(true).
    Label("Downloading...")

bar.SetProgress(75)
fmt.Println(bar.View())
// Output: "Downloading... █████████████████████████████████████░░░░░░░░░░░░░ 075%"
```

### Example 3: Custom Characters

```go
bar := progress.NewBar(30).
    FillChar('▓').
    EmptyChar('▒').
    SetProgress(50)

fmt.Println(bar.View())
// Output: "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒"
```

### Example 4: Animated Spinner

```go
type AppModel struct {
    spinner *progress.Spinner
}

func (m AppModel) Init() tea.Cmd {
    return m.spinner.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.TickMsg:
        updated, cmd := m.spinner.Update(msg)
        m.spinner = updated.(*progress.Spinner)
        return m, cmd
    case tea.KeyMsg:
        return m, tea.Quit
    }
    return m, nil
}

func (m AppModel) View() string {
    return m.spinner.View()
}

func main() {
    model := AppModel{
        spinner: progress.NewSpinner("dots").Label("Loading..."),
    }
    p := tea.New(model)
    p.Run()
}
```

### Example 5: Multiple Progress Bars

```go
bars := []*progress.Bar{
    progress.NewBar(40).Label("Task 1").ShowPercent(true),
    progress.NewBar(40).Label("Task 2").ShowPercent(true),
    progress.NewBar(40).Label("Task 3").ShowPercent(true),
}

// Update progress
bars[0].SetProgress(75)
bars[1].SetProgress(50)
bars[2].SetProgress(25)

// Render all
for _, bar := range bars {
    fmt.Println(bar.View())
}
```

## Spinner Styles

Phoenix provides 15 pre-defined spinner styles:

| Style | Animation | Description |
|-------|-----------|-------------|
| `dots` | ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏ | Unicode Braille dots (most popular) |
| `line` | \| / - \\ | Classic ASCII line spinner |
| `arrow` | ← ↖ ↑ ↗ → ↘ ↓ ↙ | Rotating arrow |
| `circle` | ◐ ◓ ◑ ◒ | Rotating circle quarters |
| `bounce` | ⠁ ⠂ ⠄ ⡀ ⢀ ⠠ ⠐ ⠈ | Bouncing ball effect |
| `dot-pulse` | ⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷ | Pulsing dots |
| `grow-vertical` | ▁ ▃ ▄ ▅ ▆ ▇ █ | Vertical growth |
| `grow-horizontal` | ▏ ▎ ▍ ▌ ▋ ▊ ▉ █ | Horizontal growth |
| `box-bounce` | ▖ ▘ ▝ ▗ | Box bouncing |
| `simple-dots` | . .. ... | Simple ASCII dots |
| `clock` | 🕐 🕑 🕒 🕓 🕔 🕕 | Clock rotation |
| `earth` | 🌍 🌎 🌏 | Spinning earth |
| `moon` | 🌑 🌒 🌓 🌔 🌕 | Moon phases |
| `toggle` | ⊶ ⊷ | On/off toggle |
| `hamburger` | ☱ ☲ ☴ | Hamburger menu animation |

## Progress Bar Customization

### Character Sets

Common fill/empty character combinations:

```go
// Solid blocks
FillChar('█').EmptyChar('░')  // Default
FillChar('▓').EmptyChar('▒')  // Shaded
FillChar('■').EmptyChar('□')  // Squares
FillChar('●').EmptyChar('○')  // Circles

// ASCII-safe
FillChar('#').EmptyChar('-')  // Classic
FillChar('=').EmptyChar(' ')  // Minimal
FillChar('>').EmptyChar('.')  // Arrows
```

### Percentage Display

```go
// Without percentage
bar := progress.NewBar(40)
// Output: "████████████████████░░░░░░░░░░░░░░░░░░░░"

// With percentage
bar := progress.NewBar(40).ShowPercent(true)
// Output: "████████████████████░░░░░░░░░░░░░░░░░░░░ 050%"
```

### Labels

```go
// Without label
bar := progress.NewBar(40)
// Output: "████████████████████░░░░░░░░░░░░░░░░░░░░"

// With label
bar := progress.NewBar(40).Label("Downloading")
// Output: "Downloading ████████████████████░░░░░░░░░░░░░░░░░░░░"

// Label + percentage
bar := progress.NewBar(40).Label("Downloading").ShowPercent(true)
// Output: "Downloading ████████████████████░░░░░░░░░░░░░░░░░░░░ 050%"
```

## Integration with gosh

The Progress component was designed as a **universal** component and will be used by **gosh** (Phoenix's cross-platform shell) for:

- Long-running command indicators
- File transfer progress
- Batch operation tracking
- Loading states during command execution

Example (gosh Week 17-18):
```go
// gosh will use Progress for long-running commands
type CommandModel struct {
    progress *progress.Bar
}

func (m CommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case CommandProgressMsg:
        m.progress.SetProgress(msg.Percentage)
    }
    return m, nil
}
```

## Architecture

The Progress component follows Phoenix's DDD architecture:

```
progress/
├── domain/              # Pure business logic
│   ├── value/          # Percentage, SpinnerStyle
│   ├── model/          # Bar, Spinner (rich models)
│   └── service/        # RenderService
├── infrastructure/      # Pre-defined spinner styles
├── api/                # Public interface
│   ├── bar.go          # Bar API + tea.Model
│   └── spinner.go      # Spinner API + tea.Model
└── examples/           # Usage examples
```

**Design Principles:**
- Rich domain models (behavior + data)
- Immutability (all operations return new instances)
- Type safety
- 80%+ test coverage (domain: 95%+)

## Performance

- **Progress bars**: Static rendering (no animation overhead)
- **Spinners**: Configurable FPS (default 10 FPS for most styles)
- **Memory**: Minimal allocations in rendering path
- **Unicode**: Correct width calculation for all spinner characters

## Testing

Run tests:
```bash
go test ./components/progress/...

# With coverage
go test -cover ./components/progress/...

# Specific layer
go test ./components/progress/domain/...
go test ./components/progress/api/...
```

## Version

- **Week 12 Day 5-6** - Initial implementation
- **Status**: v0.1.0-alpha
- **Coverage**: 80%+ (domain: 95%+)

## License

Part of Phoenix TUI Framework
Organization: phoenix-tui
Repository: github.com/phoenix-tui/phoenix

---

**Next Component**: Week 13-14 - Render (High-performance renderer)
