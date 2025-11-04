# Migrating from Bubbletea to Phoenix

**Target Audience**: Developers familiar with Charm ecosystem (Bubbletea, Lipgloss, Bubbles)
**Based On**: Real-world GoSh migration experience (130+ tests, 4 UI modes)
**Reading Time**: 20 minutes
**Difficulty**: Intermediate

---

## 🎯 Why Migrate to Phoenix?

Phoenix solves critical issues in the Charm ecosystem while maintaining the excellent Elm Architecture pattern:

### Problems Phoenix Solves

| Issue | Charm (Bubbletea/Lipgloss) | Phoenix |
|-------|----------------------------|---------|
| **Unicode/Emoji Width** | Broken for months ([lipgloss#562](https://github.com/charmbracelet/lipgloss/issues/562)) | ✅ Perfect with uniwidth |
| **Performance** | 450ms with 1000+ lines | ✅ 20-40ms (10x faster) |
| **PR Review Time** | 60-90 days average | ✅ Active development |
| **Breaking Changes** | Frequent, no migration path | ✅ Semver + migration guides |
| **Architecture** | Monolithic | ✅ Modular DDD |

### What You Keep

✅ **Elm Architecture (MVU)** - Same mental model
✅ **Familiar patterns** - Init/Update/View
✅ **Component composition** - Similar to Bubbles
✅ **Terminal abstractions** - Enhanced platform support

---

## 📋 Quick API Mapping

### Core Framework

| Bubbletea | Phoenix | Notes |
|-----------|---------|-------|
| `tea.Model` | `tea.Model[T]` | Generic type parameter |
| `tea.Cmd` | `tea.Cmd` | Same concept |
| `tea.Msg` | `tea.Msg` | Same interface |
| `tea.Program` | `tea.Program[T]` | Generic type |
| `tea.NewProgram()` | `tea.NewProgram()` | Same API |
| `tea.Quit()` | `tea.Quit()` | Same |
| `tea.Batch()` | `tea.Batch()` | Same |

### Styling

| Lipgloss | Phoenix | Notes |
|----------|---------|-------|
| `lipgloss.NewStyle()` | `style.New()` | Constructor |
| `.Foreground(color)` | `.Foreground(color)` | Same API |
| `.Background(color)` | `.Background(color)` | Same API |
| `.Bold(true)` | `.Bold()` | Simplified |
| `.Italic(true)` | `.Italic()` | Simplified |
| `.Render(text)` | `.Render(text)` | Same |
| `.Width(n)` | `.Width(n)` | Same |
| `.Padding(...)` | `.Padding(...)` | Same |

### Components

| Bubbles | Phoenix | Notes |
|---------|---------|-------|
| `textinput.Model` | `input.Input` | Renamed |
| `viewport.Model` | `viewport.Viewport` | Same name |
| `list.Model` | `list.List` | Same name |
| `table.Model` | `table.Table` | Same name |
| `spinner.Model` | `spinner.Spinner` | Same name |
| `progress.Model` | `progress.Progress` | Same name |

---

## 🚀 Step-by-Step Migration Guide

### Step 1: Update Imports

**Before (Bubbletea):**
```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
)
```

**After (Phoenix):**
```go
import (
    "github.com/phoenix-tui/phoenix/tea"
    "github.com/phoenix-tui/phoenix/style"
    "github.com/phoenix-tui/phoenix/components/input"
    "github.com/phoenix-tui/phoenix/components/viewport"
)
```

---

### Step 2: Update Model Definition

**Before (Bubbletea):**
```go
type Model struct {
    input    textinput.Model
    viewport viewport.Model
    ready    bool
}

func (m Model) Init() tea.Cmd {
    return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ...
    return m, nil
}

func (m Model) View() string {
    return "..."
}
```

**After (Phoenix):**
```go
type Model struct {
    input    *input.Input      // Note: pointer (for now, see API notes)
    viewport *viewport.Viewport
    ready    bool
}

func (m Model) Init() tea.Cmd {
    return input.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    // Note: tea.Model[Model] is generic
    // ...
    return m, nil
}

func (m Model) View() string {
    return "..."
}
```

**Key Changes**:
1. Generic type parameter: `tea.Model[Model]`
2. Components as pointers (value semantics available in v0.1.0)

---

### Step 3: Update Component Initialization

**Before (Bubbletea):**
```go
func initialModel() Model {
    ti := textinput.New()
    ti.Placeholder = "Type something..."
    ti.Focus()
    ti.CharLimit = 156
    ti.Width = 20

    vp := viewport.New(80, 24)
    vp.SetContent("Hello, World!")

    return Model{
        input:    ti,
        viewport: vp,
    }
}
```

**After (Phoenix - Current):**
```go
func initialModel() Model {
    // Phoenix uses pointer receivers (uses pointers in v0.1.0)
    ti := input.New()
    ti.SetPlaceholder("Type something...")
    ti.Focus()
    ti.SetCharLimit(156)
    ti.SetWidth(20)

    vp := viewport.New(80, 24)
    vp.SetContent("Hello, World!")

    return Model{
        input:    ti,
        viewport: vp,
    }
}
```

**After (Phoenix - v0.1.0 STABLE):**
```go
func initialModel() Model {
    // Future: Functional Options Pattern + Value Semantics
    ti := input.New(
        input.WithPlaceholder("Type something..."),
        input.WithCharLimit(156),
        input.WithWidth(20),
        input.WithFocus(true),
    )

    vp := viewport.New(
        viewport.WithSize(80, 24),
        viewport.WithContent("Hello, World!"),
    )

    return Model{
        input:    ti,  // Value, not pointer
        viewport: vp,
    }
}
```

---

### Step 4: Update Message Handling

**Before (Bubbletea):**
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Msg) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter":
            value := m.input.Value()
            m.input.SetValue("")  // Clear input
            // Process value...
        }

    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 3
    }

    // Update components
    m.input, cmd = m.input.Update(msg)
    return m, cmd
}
```

**After (Phoenix - Current):**
```go
func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter":
            value := m.input.Value()
            m.input.SetValue("")  // IMPORTANT: See note below!
            // Process value...
        }

    case tea.WindowSizeMsg:
        // Phoenix: Must reassign viewport on resize
        m.viewport = viewport.New(msg.Width, msg.Height-3)
        m.viewport.SetContent(m.content)  // Restore content
    }

    // Update components
    m.input, cmd = m.input.Update(msg)
    return m, cmd
}
```

**CRITICAL Note**: Current Phoenix API has a gotcha with `SetValue("")` - see Troubleshooting section!

---

### Step 5: Update Styling

**Before (Lipgloss):**
```go
var (
    titleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Bold(true).
        Padding(0, 1)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF0000"))
)

func (m Model) View() string {
    title := titleStyle.Render("My App")
    error := errorStyle.Render("Error!")
    return lipgloss.JoinVertical(lipgloss.Left, title, error)
}
```

**After (Phoenix):**
```go
var (
    titleStyle = style.New().
        Foreground(style.Color("#FAFAFA")).
        Background(style.Color("#7D56F4")).
        Bold().
        Padding(0, 1)

    errorStyle = style.New().
        Foreground(style.Color("#FF0000"))
)

func (m Model) View() string {
    title := titleStyle.Render("My App")
    error := errorStyle.Render("Error!")
    return style.JoinVertical(style.Left, title, error)
}
```

**Key Differences**:
- `.Bold()` instead of `.Bold(true)` (more idiomatic Go)
- `style.Color()` instead of `lipgloss.Color()`
- Same rendering API

---

### Step 6: Update Terminal Operations

**Before (Bubbletea - manual ANSI):**
```go
// In Bubbletea, you often wrote raw ANSI:
fmt.Print("\033[2J")       // Clear screen
fmt.Print("\033[H")        // Move cursor home
fmt.Print("\033[?25l")     // Hide cursor
```

**After (Phoenix - abstracted):**
```go
import "github.com/phoenix-tui/phoenix/terminal"

term := terminal.New()  // Auto-detects best implementation
term.Clear()
term.SetCursorPosition(0, 0)
term.HideCursor()

// Platform-optimized:
// - Windows Console: Win32 API calls (10x faster)
// - Unix/Git Bash: ANSI escape codes
```

---

## 🧪 Migrating Tests

### Problem: Terminal Nil Pointer in Tests

**Before (Bubbletea - no issue):**
```go
func TestModelUpdate(t *testing.T) {
    m := initialModel()
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
    // Works fine - no terminal operations
}
```

**After (Phoenix - WRONG):**
```go
func TestModelUpdate(t *testing.T) {
    m := initialModel()
    m.terminal = nil  // Tests don't need terminal
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
    // PANIC! nil pointer dereference if Update calls terminal.ClearLine()
}
```

**After (Phoenix - CORRECT):**
```go
import (
    phoenixtesting "github.com/phoenix-tui/phoenix/testing"
)

func TestModelUpdate(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()  // No-op implementation
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
    // Works! All terminal calls succeed silently
}
```

### Using MockTerminal for Verification

**Phoenix gives you powerful test helpers:**
```go
func TestRenderCallsClearLine(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{terminal: mock}

    m.Render()

    // Verify terminal operations
    assert.Contains(t, mock.Calls, "ClearLine")
    assert.Equal(t, 1, mock.CallCount("ClearLine"))
}
```

---

## ⚠️ Common Pitfalls (From GoSh Migration)

### 1. SetValue("") Doesn't Clear Input

**Problem**: After executing command, input still shows old text.

**Wrong:**
```go
case tea.KeyMsg:
    if msg.String() == "enter" {
        m.input.SetValue("")  // Doesn't work reliably!
    }
```

**Right (temporary workaround):**
```go
case tea.KeyMsg:
    if msg.String() == "enter" {
        // Create new input instead of clearing
        m.input = input.New(
            input.WithPlaceholder(m.input.Placeholder()),
            input.WithWidth(m.input.Width()),
        )
    }
```

**v0.1.0 STABLE:**
```go
case tea.KeyMsg:
    if msg.String() == "enter" {
        // Value semantics - reassignment required
        m.input = m.input.SetValue("")
        //        ^^^^^^^^^ Returns NEW Input
    }
}
```

**Status**: Fixed in v0.1.0 with value semantics API.

---

### 2. Viewport Resize Requires Recreate

**Problem**: Setting viewport size doesn't work as expected.

**Wrong:**
```go
case tea.WindowSizeMsg:
    m.viewport.Width = msg.Width
    m.viewport.Height = msg.Height
    // Doesn't resize properly!
```

**Right:**
```go
case tea.WindowSizeMsg:
    oldContent := m.viewport.Content()
    m.viewport = viewport.New(msg.Width, msg.Height)
    m.viewport.SetContent(oldContent)
```

**Why**: This is a Tea MVU pattern quirk (same in Bubbletea).

---

### 3. Tests Fail with Nil Terminal

**Problem**: Tests panic with `nil pointer dereference` on terminal operations.

**Wrong:**
```go
func TestExecuteCommand(t *testing.T) {
    m := &Model{terminal: nil}  // WILL PANIC!
    m.executeCommand("ls")
}
```

**Right:**
```go
import phoenixtesting "github.com/phoenix-tui/phoenix/testing"

func TestExecuteCommand(t *testing.T) {
    m := &Model{
        terminal: phoenixtesting.NewNullTerminal(),  // No-op terminal
    }
    m.executeCommand("ls")
}
```

---

### 4. Component Pointer Semantics (Temporary)

**Current Phoenix (v0.1.0 STABLE):**
```go
type Model struct {
    input *input.Input  // Pointer required
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    m.input.SetValue("text")  // Modifies in place
    return m, nil
}
```

**Future Phoenix (v0.2.0+):**
```go
type Model struct {
    input input.Input  // VALUE, not pointer
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    m.input = m.input.SetValue("text")  // Returns NEW Input
    //        ^^^^^^^^^ Reassignment required!
    return m, nil
}
```

**Migration**: Phoenix will provide migration guide and examples when API changes.

---

## 🎓 Real-World Example: GoSh Migration

Here's a simplified excerpt from GoSh's actual migration:

### Before (Bubbletea):
```go
type Model struct {
    shellInput   textinput.Model
    viewport     viewport.Model
    multilineMode bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            cmd := m.shellInput.Value()
            m.shellInput.SetValue("")
            return m, executeCommand(cmd)
        }
    }

    var cmd tea.Cmd
    m.shellInput, cmd = m.shellInput.Update(msg)
    return m, cmd
}
```

### After (Phoenix):
```go
import (
    "github.com/phoenix-tui/phoenix/tea"
    "github.com/phoenix-tui/phoenix/components/input"
    "github.com/phoenix-tui/phoenix/components/viewport"
    "github.com/phoenix-tui/phoenix/terminal"
)

type Model struct {
    shellInput    *input.Input
    viewport      *viewport.Viewport
    terminal      terminal.Terminal
    multilineMode bool
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            cmd := m.shellInput.Value()

            // Clear terminal lines (Phoenix-specific)
            if m.terminal != nil {
                if m.multilineMode {
                    _ = m.terminal.ClearLines(3)
                } else {
                    _ = m.terminal.ClearLine()
                }
            }

            m.shellInput.SetValue("")
            return m, executeCommand(cmd)
        }
    }

    var cmd tea.Cmd
    m.shellInput, cmd = m.shellInput.Update(msg)
    return m, cmd
}
```

**Results**:
- ✅ 10x performance improvement (450ms → 40ms)
- ✅ Perfect Unicode/Emoji support
- ✅ Platform-optimized terminal operations
- ⚠️ 10+ REPL tests failed (API differences) - addressed in v0.1.0

---

## 📊 Migration Checklist

Use this checklist to track your migration:

### Preparation
- [ ] Read this migration guide
- [ ] Read `TROUBLESHOOTING.md` for common issues
- [ ] Backup your project (git commit/tag)

### Code Migration
- [ ] Update imports (Bubbletea → Phoenix)
- [ ] Update Model definition (add generic type parameter)
- [ ] Update component types (textinput.Model → input.Input, etc.)
- [ ] Update component initialization
- [ ] Update message handling (check type assertions)
- [ ] Update styling (Lipgloss → Phoenix style)
- [ ] Add terminal abstraction (if needed)

### Test Migration
- [ ] Update test imports
- [ ] Add NullTerminal for tests with terminal operations
- [ ] Update component assertions (API differences)
- [ ] Verify all tests pass

### Verification
- [ ] Application runs without panics
- [ ] Unicode/Emoji render correctly
- [ ] Performance improved (measure with profiler)
- [ ] All features work as expected

### Optional Optimizations
- [ ] Use platform-specific terminal operations
- [ ] Optimize rendering with Phoenix render engine
- [ ] Apply Phoenix layout system for complex UIs

---

## 🚀 Performance Optimization

Phoenix gives you 10x performance out of the box, but here are tips for maximum speed:

### 1. Use Differential Rendering

**Phoenix render engine** automatically optimizes rendering:
```go
import "github.com/phoenix-tui/phoenix/render"

renderer := render.New()
renderer.Render(m.View())  // Only writes changed cells
```

### 2. Platform-Specific Terminal Operations

**Phoenix auto-detects** best implementation:
```go
term := terminal.New()
if term.SupportsDirectPositioning() {
    // Use fast absolute positioning (Windows Console API)
    term.WriteAt(x, y, content)
} else {
    // Use ANSI relative movements
    term.Write(content)
}
```

### 3. Optimize Large Content

**Viewport with thousands of lines**:
```go
// Phoenix viewport efficiently handles large content
vp := viewport.New(80, 24)
vp.SetContent(strings.Join(lines, "\n"))  // No lag with 10k+ lines

// Bubbletea viewport lagged at 1000+ lines
```

---

## 📚 Next Steps

After migration:

1. **Read `TESTING.md`** - Learn Phoenix testing patterns
2. **Read `TROUBLESHOOTING.md`** - Quick fixes for common issues
3. **Explore Phoenix components** - See what's new beyond Bubbles
4. **Join community** - GitHub Discussions, Discord (coming soon)

---

## 🆘 Getting Help

- **GitHub Issues**: [github.com/phoenix-tui/phoenix/issues](https://github.com/phoenix-tui/phoenix/issues)
- **Discussions**: [github.com/phoenix-tui/phoenix/discussions](https://github.com/phoenix-tui/phoenix/discussions)
- **Documentation**: [docs/](../docs/)
- **Examples**: [examples/](../../examples/)

---

## 📖 Appendix: API Compatibility Table

### Messages

| Bubbletea | Phoenix | Compatibility |
|-----------|---------|---------------|
| `tea.KeyMsg` | `tea.KeyMsg` | ✅ 100% |
| `tea.MouseMsg` | `tea.MouseMsg` | ✅ 100% |
| `tea.WindowSizeMsg` | `tea.WindowSizeMsg` | ✅ 100% |
| `tea.FocusMsg` | `tea.FocusMsg` | ✅ 100% |
| `tea.BlurMsg` | `tea.BlurMsg` | ✅ 100% |

### Commands

| Bubbletea | Phoenix | Compatibility |
|-----------|---------|---------------|
| `tea.Quit` | `tea.Quit` | ✅ 100% |
| `tea.Batch` | `tea.Batch` | ✅ 100% |
| `tea.Sequence` | `tea.Sequence` | ✅ 100% |
| `tea.Tick` | `tea.Tick` | ✅ 100% |
| `tea.Every` | `tea.Every` | ✅ 100% |

### Program Options

| Bubbletea | Phoenix | Compatibility |
|-----------|---------|---------------|
| `tea.WithAltScreen()` | `tea.WithAltScreen()` | ✅ 100% |
| `tea.WithMouseCellMotion()` | `tea.WithMouseCellMotion()` | ✅ 100% |
| `tea.WithMouseAllMotion()` | `tea.WithMouseAllMotion()` | ✅ 100% |
| `tea.WithoutRenderer()` | `tea.WithoutRenderer()` | ✅ 100% |
| `tea.WithFilter()` | `tea.WithFilter()` | ✅ 100% |

---

*Migration Guide Version: 1.0*
*Last Updated: 2025-11-04*
*Based on: GoSh production migration experience*
*Target Phoenix Version: v0.1.0 (STABLE)*
