# Testing Phoenix Applications

**Target Audience**: Developers building Phoenix TUI applications
**Based On**: GoSh test suite (130+ tests), Phoenix internal testing
**Reading Time**: 25 minutes
**Prerequisites**: Basic Go testing knowledge (`go test`, `testing` package)

---

## 🎯 Testing Philosophy

Phoenix follows TDD (Test-Driven Development) with high coverage targets:

### Coverage Targets

| Layer | Coverage | Rationale |
|-------|----------|-----------|
| **Domain** | 95%+ | Pure business logic - easy to test, critical to reliability |
| **Application** | 90%+ | Use cases - testable with mocked infrastructure |
| **Infrastructure** | 80%+ | I/O operations - harder to test, focus on critical paths |
| **API** | 85%+ | Public interface - example-based tests |

### Why Test TUIs?

1. **State Management** - Complex state machines need verification
2. **User Input** - Edge cases in key handling, Unicode support
3. **Rendering** - Correct output generation
4. **Performance** - Ensure responsiveness (< 50ms frame time)
5. **Regressions** - Catch API changes, maintain quality

---

## 🏗️ Test Structure

### Recommended Test Organization

```
your-app/
├── internal/
│   ├── domain/           # Business logic
│   │   └── model_test.go
│   ├── application/      # Use cases
│   │   └── handler_test.go
│   └── interfaces/       # UI layer
│       ├── tui/
│       │   ├── model_test.go
│       │   ├── update_test.go
│       │   └── view_test.go
│       └── testdata/     # Test fixtures
├── tests/                # Integration tests
│   ├── e2e_test.go
│   └── fixtures/
└── go.mod
```

---

## 🧪 Unit Testing: Model/Update/View

### Testing Model Initialization

```go
package tui

import (
    "testing"

    "github.com/phoenix-tui/phoenix/tea"
    phoenixtesting "github.com/phoenix-tui/phoenix/testing"
    "github.com/stretchr/testify/assert"
)

func TestModel_Init(t *testing.T) {
    m := initialModel()

    cmd := m.Init()

    // Verify initial state
    assert.NotNil(t, m.input, "input should be initialized")
    assert.False(t, m.ready, "ready should be false before WindowSizeMsg")
    assert.NotNil(t, cmd, "Init should return a command")
}
```

### Testing Update Logic

```go
func TestModel_Update_KeyEnter(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()  // CRITICAL!
    m.input.SetValue("test command")

    // Send Enter key
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    updated, cmd := m.Update(msg)

    // Verify state changes
    assert.Empty(t, updated.input.Value(), "input should be cleared")
    assert.NotNil(t, cmd, "should return command to execute")
}

func TestModel_Update_WindowSize(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()

    // Send window resize
    msg := tea.WindowSizeMsg{Width: 100, Height: 50}
    updated, _ := m.Update(msg)

    // Verify viewport resized
    assert.Equal(t, 100, updated.viewport.Width)
    assert.Equal(t, 47, updated.viewport.Height) // 50 - 3 for input
    assert.True(t, updated.ready, "ready should be true after first WindowSizeMsg")
}

func TestModel_Update_Quit(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()

    tests := []struct {
        name string
        key  string
    }{
        {"ctrl+c", "ctrl+c"},
        {"q", "q"},
        {"esc", "esc"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
            _, cmd := m.Update(msg)

            assert.NotNil(t, cmd, "quit command should be returned")
            // Note: Can't easily test cmd == tea.Quit without executing it
        })
    }
}
```

### Testing View Rendering

```go
func TestModel_View_NotReady(t *testing.T) {
    m := initialModel()
    m.ready = false

    view := m.View()

    assert.Contains(t, view, "Initializing", "should show loading state")
}

func TestModel_View_Ready(t *testing.T) {
    m := initialModel()
    m.ready = true
    m.terminal = phoenixtesting.NewNullTerminal()

    view := m.View()

    // Verify expected content
    assert.Contains(t, view, m.input.View(), "should include input")
    assert.Contains(t, view, m.viewport.View(), "should include viewport")
}

func TestModel_View_Error(t *testing.T) {
    m := initialModel()
    m.ready = true
    m.error = errors.New("test error")

    view := m.View()

    assert.Contains(t, view, "test error", "should display error")
    assert.Contains(t, view, "Error:", "should include error label")
}
```

---

## 🔌 Testing with phoenix/testing Package

### Using NullTerminal

**Problem**: Production code calls `terminal.ClearLine()` but tests crash with nil pointer.

**Solution**: Use `NullTerminal` - all methods succeed silently.

```go
import phoenixtesting "github.com/phoenix-tui/phoenix/testing"

func TestExecuteCommand(t *testing.T) {
    m := &Model{
        terminal: phoenixtesting.NewNullTerminal(),  // ← No-op terminal
    }

    m.executeCommand("ls")

    // All terminal operations succeed without panics
    assert.NotNil(t, m)
}
```

**Best Practice**: Initialize NullTerminal in test helpers:

```go
func newTestModel() *Model {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()
    return m
}

func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name string
        cmd  string
    }{
        {"simple", "ls"},
        {"with args", "ls -la"},
        {"complex", "find . -name '*.go'"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := newTestModel()  // ← Reusable helper
            m.executeCommand(tt.cmd)
            // ...
        })
    }
}
```

### Using MockTerminal for Verification

**When to use**: When you need to verify terminal operations were called.

```go
func TestRender_ClearsScreen(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{terminal: mock}

    m.Render()

    // Verify terminal operations
    assert.Equal(t, 1, mock.CallCount("Clear"), "Clear should be called once")
    assert.Equal(t, 1, mock.CallCount("SetCursorPosition"), "SetCursorPosition should be called")

    // Verify specific arguments
    expectedCalls := []string{
        "Clear",
        "SetCursorPosition(0, 0)",
        "HideCursor",
        "Write(\"...\")",
        "ShowCursor",
    }
    assert.Equal(t, expectedCalls, mock.Calls)
}

func TestMultiline_ClearsMultipleLines(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{
        terminal:      mock,
        multilineMode: true,
    }

    m.clearInput()

    // Verify ClearLines called with correct count
    assert.Contains(t, mock.Calls, "ClearLines(3)")
}
```

### MockTerminal Advanced Usage

```go
func TestTerminalCallOrder(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{terminal: mock}

    m.RenderFrame()

    // Verify call order
    expectedOrder := []string{
        "HideCursor",
        "SetCursorPosition(0, 0)",
        "Write(\"Frame content\")",
        "ShowCursor",
    }

    assert.Equal(t, expectedOrder, mock.Calls, "calls should be in correct order")
}

func TestTerminal_Reset_BetweenTests(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()

    // First test phase
    mock.Clear()
    assert.Equal(t, 1, mock.CallCount("Clear"))

    // Reset for next phase
    mock.Reset()
    assert.Equal(t, 0, mock.CallCount("Clear"), "CallCount should be 0 after reset")
    assert.Empty(t, mock.Calls, "Calls should be empty after reset")

    // Second test phase
    mock.HideCursor()
    assert.Equal(t, 1, mock.CallCount("HideCursor"))
}
```

---

## 🎨 Testing Components

### Testing Input Component

```go
import "github.com/phoenix-tui/phoenix/components/input"

func TestInput_SetValue(t *testing.T) {
    inp := input.New()

    inp.SetValue("test")

    assert.Equal(t, "test", inp.Value())
}

func TestInput_CharLimit(t *testing.T) {
    inp := input.New()
    inp.SetCharLimit(5)

    inp.SetValue("toolong")

    assert.Equal(t, "toolo", inp.Value(), "should truncate to char limit")
}

func TestInput_Focus(t *testing.T) {
    inp := input.New()

    assert.False(t, inp.Focused(), "should not be focused initially")

    inp.Focus()
    assert.True(t, inp.Focused(), "should be focused after Focus()")

    inp.Blur()
    assert.False(t, inp.Focused(), "should not be focused after Blur()")
}

func TestInput_Update_KeyPress(t *testing.T) {
    inp := input.New()
    inp.Focus()

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
    updated, _ := inp.Update(msg)

    assert.Equal(t, "a", updated.Value())
}
```

### Testing Viewport Component

```go
import "github.com/phoenix-tui/phoenix/components/viewport"

func TestViewport_SetContent(t *testing.T) {
    vp := viewport.New(80, 24)

    vp.SetContent("Line 1\nLine 2\nLine 3")

    assert.Contains(t, vp.View(), "Line 1")
    assert.Contains(t, vp.View(), "Line 2")
}

func TestViewport_Scroll(t *testing.T) {
    vp := viewport.New(80, 10)
    content := strings.Repeat("Line\n", 100)
    vp.SetContent(content)

    // Scroll down
    vp.LineDown(5)
    assert.Equal(t, 5, vp.YOffset())

    // Scroll up
    vp.LineUp(2)
    assert.Equal(t, 3, vp.YOffset())

    // Scroll to bottom
    vp.GotoBottom()
    assert.True(t, vp.AtBottom())

    // Scroll to top
    vp.GotoTop()
    assert.Equal(t, 0, vp.YOffset())
    assert.True(t, vp.AtTop())
}

func TestViewport_Update_MouseScroll(t *testing.T) {
    vp := viewport.New(80, 10)
    vp.SetContent(strings.Repeat("Line\n", 100))

    // Mouse wheel down
    msg := tea.MouseMsg{Type: tea.MouseWheelDown}
    updated, _ := vp.Update(msg)

    assert.Greater(t, updated.YOffset(), 0, "should scroll down")
}
```

---

## 🔄 Testing State Machines

### Example: Multi-State Application

```go
type AppState int

const (
    StateMenu AppState = iota
    StateEditor
    StateHelp
)

type Model struct {
    state   AppState
    menu    *menu.Model
    editor  *editor.Model
    help    *help.Model
}

func TestModel_StateTransitions(t *testing.T) {
    tests := []struct {
        name       string
        initial    AppState
        key        string
        expected   AppState
    }{
        {"menu to editor", StateMenu, "e", StateEditor},
        {"menu to help", StateMenu, "?", StateHelp},
        {"help to menu", StateHelp, "q", StateMenu},
        {"editor to menu", StateEditor, "esc", StateMenu},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &Model{state: tt.initial}
            m.terminal = phoenixtesting.NewNullTerminal()

            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
            updated, _ := m.Update(msg)

            assert.Equal(t, tt.expected, updated.state)
        })
    }
}
```

---

## 🎭 Table-Driven Tests

**Best practice**: Use table-driven tests for comprehensive coverage.

```go
func TestModel_HandleCommand(t *testing.T) {
    tests := []struct {
        name        string
        command     string
        wantError   bool
        wantOutput  string
    }{
        {
            name:       "simple ls",
            command:    "ls",
            wantError:  false,
            wantOutput: "",  // Platform-specific
        },
        {
            name:       "invalid command",
            command:    "nonexistent-cmd-12345",
            wantError:  true,
            wantOutput: "",
        },
        {
            name:       "cd command",
            command:    "cd /tmp",
            wantError:  false,
            wantOutput: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := newTestModel()
            m.terminal = phoenixtesting.NewNullTerminal()

            err := m.handleCommand(tt.command)

            if tt.wantError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            if tt.wantOutput != "" {
                assert.Contains(t, m.output, tt.wantOutput)
            }
        })
    }
}
```

---

## 📊 Testing Performance

### Benchmark Rendering

```go
func BenchmarkModel_View(b *testing.B) {
    m := initialModel()
    m.ready = true
    m.terminal = phoenixtesting.NewNullTerminal()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

func BenchmarkViewport_LargeContent(b *testing.B) {
    content := strings.Repeat("Line of text\n", 10000)
    vp := viewport.New(80, 24)
    vp.SetContent(content)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = vp.View()
    }
}

// Expected: < 1ms per frame (60 FPS = 16ms budget)
```

### Measuring Frame Time

```go
func TestModel_RenderPerformance(t *testing.T) {
    m := newTestModel()
    m.viewport.SetContent(strings.Repeat("Line\n", 1000))

    start := time.Now()
    _ = m.View()
    elapsed := time.Since(start)

    assert.Less(t, elapsed, 50*time.Millisecond, "rendering should be fast")
}
```

---

## 🧩 Integration Testing

### Testing Full User Flows

```go
func TestUserFlow_ExecuteCommand(t *testing.T) {
    // Setup
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()

    // Step 1: Initialize
    cmd := m.Init()
    assert.NotNil(t, cmd)

    // Step 2: Window size
    m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
    assert.True(t, m.ready)

    // Step 3: Type command
    for _, r := range "ls -la" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        m, _ = m.Update(msg)
    }
    assert.Equal(t, "ls -la", m.input.Value())

    // Step 4: Execute command
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    m, cmd = m.Update(msg)

    assert.Empty(t, m.input.Value(), "input should be cleared")
    assert.NotNil(t, cmd, "should execute command")
}
```

### Testing Error Scenarios

```go
func TestErrorHandling_InvalidCommand(t *testing.T) {
    m := newTestModel()

    // Execute invalid command
    m.input.SetValue("invalid-command-xyz")
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    m, _ = m.Update(msg)

    // Verify error displayed
    assert.NotNil(t, m.error)
    assert.Contains(t, m.View(), "Error:")
}
```

---

## 🎯 Best Practices

### 1. Always Use Test Helpers

**Bad**:
```go
func TestA(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()
    m.ready = true
    // ...
}

func TestB(t *testing.T) {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()
    m.ready = true
    // ... same setup repeated!
}
```

**Good**:
```go
func newTestModel() *Model {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()
    m.ready = true
    return m
}

func TestA(t *testing.T) {
    m := newTestModel()
    // ...
}

func TestB(t *testing.T) {
    m := newTestModel()
    // ...
}
```

### 2. Test One Thing Per Test

**Bad**:
```go
func TestEverything(t *testing.T) {
    m := newTestModel()

    // Tests too many things
    m.input.SetValue("test")
    assert.Equal(t, "test", m.input.Value())

    m.input.SetCharLimit(5)
    assert.Equal(t, 5, m.input.CharLimit())

    m.viewport.SetContent("content")
    assert.Contains(t, m.viewport.View(), "content")
    // ... etc
}
```

**Good**:
```go
func TestInput_SetValue(t *testing.T) {
    m := newTestModel()
    m.input.SetValue("test")
    assert.Equal(t, "test", m.input.Value())
}

func TestInput_SetCharLimit(t *testing.T) {
    m := newTestModel()
    m.input.SetCharLimit(5)
    assert.Equal(t, 5, m.input.CharLimit())
}

func TestViewport_SetContent(t *testing.T) {
    m := newTestModel()
    m.viewport.SetContent("content")
    assert.Contains(t, m.viewport.View(), "content")
}
```

### 3. Use testdata/ for Fixtures

```
tests/
└── testdata/
    ├── input_commands.txt
    ├── expected_output.txt
    └── golden/
        ├── menu_view.golden
        └── editor_view.golden
```

```go
func TestModel_View_GoldenFile(t *testing.T) {
    m := newTestModel()
    view := m.View()

    goldenPath := "testdata/golden/menu_view.golden"

    if *update {
        // Update golden file
        os.WriteFile(goldenPath, []byte(view), 0644)
    }

    expected, _ := os.ReadFile(goldenPath)
    assert.Equal(t, string(expected), view)
}
```

### 4. Isolate Terminal Operations

**Pattern**: Wrap terminal in interface for easy mocking.

```go
// production.go
type TerminalOps interface {
    ClearLine() error
    Write(s string) error
}

type Model struct {
    term TerminalOps
}

// main.go
m := Model{term: terminal.New()}

// model_test.go
m := Model{term: phoenixtesting.NewNullTerminal()}
```

---

## 📝 Testing Checklist

Use this checklist for comprehensive test coverage:

### Model Tests
- [ ] Init() returns correct initial state
- [ ] Init() returns correct command
- [ ] Update handles all message types
- [ ] Update maintains state correctly
- [ ] View renders correctly in all states
- [ ] View handles errors gracefully

### Component Tests
- [ ] Components initialize with correct defaults
- [ ] Setters update state correctly
- [ ] Getters return correct values
- [ ] Update handles keyboard input
- [ ] Update handles mouse input
- [ ] View renders expected output

### Integration Tests
- [ ] Full user flows work end-to-end
- [ ] Error scenarios handled correctly
- [ ] State transitions work
- [ ] Commands execute correctly

### Performance Tests
- [ ] Rendering < 50ms per frame
- [ ] No memory leaks in long sessions
- [ ] Large content handles efficiently

---

## 🚀 Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestModel_Update ./internal/interfaces/tui

# Run benchmarks
go test -bench=. ./...

# Verbose output
go test -v ./...

# Race detection
go test -race ./...
```

---

## 🔍 Debugging Tests

### Print Debugging

```go
func TestModel_Debug(t *testing.T) {
    m := newTestModel()

    t.Logf("Before update: %+v", m)
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
    t.Logf("After update: %+v", m)
}
```

### Mock Terminal Debugging

```go
func TestDebug_TerminalCalls(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{terminal: mock}

    m.Render()

    // Print all terminal calls
    t.Logf("Terminal calls: %v", mock.Calls)
    t.Logf("ClearLine count: %d", mock.CallCount("ClearLine"))
}
```

---

## 📚 Resources

- **Phoenix testing package**: `github.com/phoenix-tui/phoenix/testing`
- **Testify assertions**: `github.com/stretchr/testify/assert`
- **Go testing docs**: https://pkg.go.dev/testing
- **Table-driven tests**: https://go.dev/wiki/TableDrivenTests

---

*Testing Guide Version: 1.0*
*Last Updated: 2025-10-19*
*Based on: GoSh test suite (130+ tests) and Phoenix internal testing*
