# Phoenix TUI Troubleshooting Guide

**Target Audience**: Developers encountering issues with Phoenix
**Based On**: Real-world GoSh migration experience, user reports
**Last Updated**: 2025-11-04 (v0.1.0-beta.6)

---

## 📋 Quick Links

- [Input Not Clearing](#input-not-clearing-after-setvalue)
- [Tests Panic with Nil Pointer](#tests-panic-nil-pointer-dereference)
- [Viewport Doesn't Resize](#viewport-doesnt-resize-properly)
- [Poor Performance](#application-lags-or-freezes)
- [Build Errors](#build-errors)
- [Unicode/Emoji Issues](#unicodeemoji-rendering-broken)

---

## 🔴 Common Runtime Issues

### Input Not Clearing After SetValue("")

**Symptom**: Calling `input.SetValue("")` doesn't clear the input field visually.

**Example**:
```go
case tea.KeyMsg:
    if msg.String() == "enter" {
        m.input.SetValue("")  // Doesn't work!
        // Input still shows old text
    }
```

**Why This Happens**: Current Phoenix API uses pointer receivers with value semantics in Tea MVU pattern. The `SetValue()` call modifies a copy, not the original.

**Workaround (Current v0.1.0-alpha)**:
```go
// Option 1: Create new input
case tea.KeyMsg:
    if msg.String() == "enter" {
        oldPlaceholder := m.input.Placeholder()
        oldWidth := m.input.Width()

        m.input = input.New()
        m.input.SetPlaceholder(oldPlaceholder)
        m.input.SetWidth(oldWidth)
        m.input.Focus()
    }

// Option 2: Store as pointer and be careful
type Model struct {
    input *input.Input  // Note: pointer
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    // This works because pointer is copied but points to same object
    m.input.SetValue("")
    return m, nil
}
```

**Permanent Fix (Coming in v0.1.0-beta.1)**:
```go
// Future: Value semantics with explicit reassignment
type Model struct {
    input input.Input  // Value, not pointer
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    m.input = m.input.SetValue("")  // Returns new Input
    //        ^^^^^^^^^ Reassignment required!
    return m, nil
}
```

**Status**: P0 issue, will be fixed in beta.1 with value semantics API.

**See Also**: [MIGRATION_FROM_BUBBLETEA.md](MIGRATION_FROM_BUBBLETEA.md#1-setvalue-doesnt-clear-input)

---

### Tests Panic: Nil Pointer Dereference

**Symptom**: Tests crash with `panic: runtime error: invalid memory address or nil pointer dereference` when calling terminal methods.

**Example**:
```go
func TestExecuteCommand(t *testing.T) {
    m := &Model{
        terminal: nil,  // WRONG!
    }

    m.executeCommand("ls")
    // PANIC at m.terminal.ClearLine()
}
```

**Error Message**:
```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x...]

goroutine 6 [running]:
yourapp.(*Model).executeCommand(...)
    /path/to/your/model.go:123
yourapp.TestExecuteCommand(...)
    /path/to/your/model_test.go:45
```

**Why This Happens**: Production code calls `m.terminal.ClearLine()` but tests don't initialize a real terminal.

**Solution: Use phoenix/testing Package**

```go
import phoenixtesting "github.com/phoenix-tui/phoenix/testing"

func TestExecuteCommand(t *testing.T) {
    m := &Model{
        terminal: phoenixtesting.NewNullTerminal(),  // ← Fix!
    }

    m.executeCommand("ls")  // No panic!
}

// Better: Create test helper
func newTestModel() *Model {
    m := initialModel()
    m.terminal = phoenixtesting.NewNullTerminal()
    return m
}

func TestExecuteCommand(t *testing.T) {
    m := newTestModel()
    m.executeCommand("ls")
}
```

**See Also**: [TESTING.md](TESTING.md#using-nullterminal)

---

### Viewport Doesn't Resize Properly

**Symptom**: After window resize, viewport doesn't display correctly or content is cut off.

**Wrong Approach**:
```go
case tea.WindowSizeMsg:
    m.viewport.Width = msg.Width
    m.viewport.Height = msg.Height - 3
    // Viewport doesn't resize properly!
```

**Correct Approach**:
```go
case tea.WindowSizeMsg:
    // Must RECREATE viewport (Tea MVU pattern quirk)
    oldContent := m.viewport.Content()
    oldOffset := m.viewport.YOffset()

    m.viewport = viewport.New(msg.Width, msg.Height-3)
    m.viewport.SetContent(oldContent)
    m.viewport.SetYOffset(oldOffset)  // Restore scroll position
```

**Why**: This is a Tea MVU pattern requirement (same in Bubbletea). Components must be recreated on structural changes.

**See Also**: [MIGRATION_FROM_BUBBLETEA.md](MIGRATION_FROM_BUBBLETEA.md#2-viewport-resize-requires-recreate)

---

### Application Lags or Freezes

**Symptom**: TUI becomes sluggish, rendering takes > 100ms, input feels delayed.

**Diagnosis**:
```go
import "time"

func (m Model) View() string {
    start := time.Now()
    view := m.render()
    elapsed := time.Since(start)

    if elapsed > 50*time.Millisecond {
        log.Printf("WARNING: Slow render: %v", elapsed)
    }

    return view
}
```

**Common Causes & Solutions**:

#### 1. Large Viewport Content

**Problem**: Rendering thousands of lines on every frame.

**Solution**: Use viewport pagination:
```go
// Bad: Rendering all 10,000 lines
vp.SetContent(strings.Join(allLines, "\n"))

// Good: Only render visible portion
visibleLines := allLines[vp.YOffset():vp.YOffset()+vp.Height]
vp.SetContent(strings.Join(visibleLines, "\n"))
```

#### 2. Expensive Styling in View()

**Problem**: Creating new styles on every render.

**Wrong**:
```go
func (m Model) View() string {
    // Creating styles on EVERY render!
    titleStyle := style.New().Foreground(style.Color("#FF0000")).Bold()
    return titleStyle.Render("Title")
}
```

**Right**:
```go
// Create styles ONCE (package-level or in Init)
var titleStyle = style.New().Foreground(style.Color("#FF0000")).Bold()

func (m Model) View() string {
    return titleStyle.Render("Title")
}
```

#### 3. Unnecessary Re-renders

**Problem**: View() called even when state hasn't changed.

**Solution**: Implement caching:
```go
type Model struct {
    cachedView    string
    viewDirty     bool
}

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    // Mark view dirty when state changes
    m.viewDirty = true
    // ...
}

func (m Model) View() string {
    if !m.viewDirty && m.cachedView != "" {
        return m.cachedView
    }

    m.cachedView = m.render()
    m.viewDirty = false
    return m.cachedView
}
```

**Performance Target**: < 50ms per frame (20 FPS minimum, 60 FPS ideal).

---

## 🔧 Build & Compilation Issues

### Build Error: Cannot Find Package

**Error**:
```
could not import github.com/phoenix-tui/phoenix/tea (no required module provides package "github.com/phoenix-tui/phoenix/tea")
```

**Solution**:
```bash
# Initialize go module (if not already)
go mod init your-app

# Add Phoenix dependency
go get github.com/phoenix-tui/phoenix/tea
go get github.com/phoenix-tui/phoenix/style
go get github.com/phoenix-tui/phoenix/components/input

# Tidy dependencies
go mod tidy
```

---

### Workspace Module Error

**Error**:
```
current directory is contained in a module that is not one of the workspace modules listed in go.work
```

**Solution**:
```bash
# Add your module to workspace
go work use .

# Or create workspace if it doesn't exist
go work init
go work use .
```

---

### Replace Directive Issues

**Error**: Phoenix modules not found despite `replace` directive.

**Check go.mod**:
```go
module your-app

require (
    github.com/phoenix-tui/phoenix/tea v0.0.0
    github.com/phoenix-tui/phoenix/style v0.0.0
)

// For local development (before v0.1.0-beta.1 release)
replace (
    github.com/phoenix-tui/phoenix/tea => ../path/to/phoenix/tea
    github.com/phoenix-tui/phoenix/style => ../path/to/phoenix/style
)
```

**Solution**: Verify paths are correct relative to your module root.

---

## 🎨 Rendering Issues

### Unicode/Emoji Rendering Broken

**Symptom**: Emoji/CJK characters have wrong width, text misaligned.

**Phoenix Solution**: This is **fixed** in Phoenix (unlike Charm ecosystem).

**Verify Fix**:
```go
import "github.com/phoenix-tui/phoenix/core/unicode"

text := "Hello 👋 世界"
width := unicode.StringWidth(text)  // Correct width with Phoenix

// Should print: Width: 10 (not 8 like buggy implementations)
fmt.Printf("Width: %d\n", width)
```

**If Still Broken**: You might be using old Charm dependencies. Check:
```bash
go list -m all | grep charmbracelet
# Should show NO Charm dependencies!
```

---

### Colors Not Displaying

**Symptom**: Styled text appears without colors.

**Diagnosis**:
```go
import "github.com/phoenix-tui/phoenix/terminal"

term := terminal.New()
fmt.Printf("Color depth: %d\n", term.ColorDepth())
fmt.Printf("True color: %v\n", term.SupportsTrueColor())

// Expected: 256 or 16777216 (true color)
```

**Common Issues**:

1. **Terminal doesn't support colors**:
   ```bash
   # Check TERM environment variable
   echo $TERM
   # Should be: xterm-256color or similar
   ```

2. **Windows Console (cmd.exe)**: Limited color support
   - Solution: Use Windows Terminal instead
   - Or: Enable VT processing in code:
   ```go
   import "github.com/phoenix-tui/phoenix/terminal/infrastructure/windows"
   windows.EnableVirtualTerminalProcessing()
   ```

3. **SSH session**: Colors stripped
   - Solution: Set `TERM=xterm-256color` on client

---

### Text Overlapping or Misaligned

**Symptom**: Text renders on top of previous text, cursor in wrong position.

**Cause**: Insufficient clearing between frames.

**Solution 1: Use ClearLine**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    if m.terminal != nil {
        _ = m.terminal.ClearLine()  // Clear before rendering
    }
    return m, nil
}
```

**Solution 2: Use Full Clear**:
```go
func (m Model) View() string {
    if m.terminal != nil {
        _ = m.terminal.Clear()  // Clear entire screen
        _ = m.terminal.SetCursorPosition(0, 0)
    }
    return m.render()
}
```

---

## 🧪 Testing Issues

### Tests Fail After Migration from Bubbletea

**Symptom**: Tests that worked with Bubbletea now fail with Phoenix.

**Common Failures**:

#### 1. Component Behavior Differences

**Bubbletea**:
```go
func TestInput(t *testing.T) {
    m := textinput.New()
    m.SetValue("test")
    assert.Equal(t, "test", m.Value())  // ✅ Works
}
```

**Phoenix (Wrong)**:
```go
func TestInput(t *testing.T) {
    m := input.New()
    m.SetValue("test")
    assert.Equal(t, "test", m.Value())  // ❌ Fails (pointer issue)
}
```

**Phoenix (Right)**:
```go
func TestInput(t *testing.T) {
    m := input.New()
    m.SetValue("test")

    // Future (v0.1.0-beta.1+):
    // m = m.SetValue("test")

    assert.Equal(t, "test", m.Value())  // ✅ Works
}
```

#### 2. Terminal Operations in Tests

See: [Tests Panic: Nil Pointer Dereference](#tests-panic-nil-pointer-dereference)

#### 3. Rendering Output Differences

**Issue**: Phoenix render output format may differ slightly from Bubbletea.

**Solution**: Update golden files or use flexible assertions:
```go
// Instead of exact match:
assert.Equal(t, expectedView, actualView)

// Use substring checks:
assert.Contains(t, actualView, "Expected Text")
assert.Contains(t, actualView, "Another Part")
```

---

## 🐛 API Confusion

### SetValue vs Reset - Which to Use?

**Question**: What's the difference? When should I use each?

**Answer** (Current API):
```go
// SetValue - changes content, preserves other state
input.SetValue("new text")  // Sets content to "new text"

// Reset - clears everything (content + internal state)
input.Reset()  // Clears content, resets cursor, etc.
```

**Recommendation**: Use `SetValue("")` for clearing in most cases.

**Future API (beta.1+)**:
```go
// SetValue - returns new Input
m.input = m.input.SetValue("text")

// Clear method (more explicit)
m.input = m.input.Clear()
```

---

### Components as Pointers vs Values

**Question**: Should I store components as `*input.Input` or `input.Input`?

**Current (v0.1.0-alpha)**:
```go
type Model struct {
    input *input.Input  // Pointer required
}
```

**Future (v0.1.0-beta.1+)**:
```go
type Model struct {
    input input.Input  // Value (immutable-ish)
}
```

**Why the Change**: Value semantics fit Elm Architecture better (immutability, explicit state changes).

---

## 🔍 Debugging Tips

### Enable Debug Logging

```go
import "log"

func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    log.Printf("Update: %T %+v", msg, msg)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        log.Printf("Key: %s", msg.String())
    }

    return m, nil
}
```

### Print Model State

```go
func (m Model) Update(msg tea.Msg) (tea.Model[Model], tea.Cmd) {
    log.Printf("Before: input=%q ready=%v", m.input.Value(), m.ready)

    // ... update logic ...

    log.Printf("After: input=%q ready=%v", m.input.Value(), m.ready)
    return m, nil
}
```

### Inspect Terminal Calls (Tests)

```go
func TestDebug(t *testing.T) {
    mock := phoenixtesting.NewMockTerminal()
    m := &Model{terminal: mock}

    m.Render()

    // Print all calls
    for i, call := range mock.Calls {
        t.Logf("Call %d: %s", i, call)
    }
}
```

---

## 🚀 Performance Profiling

### CPU Profiling

```bash
# Run with profiling
go test -cpuprofile=cpu.prof -bench=.

# Analyze profile
go tool pprof cpu.prof
# In pprof: type 'top' to see hotspots
```

### Memory Profiling

```bash
# Run with memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze profile
go tool pprof mem.prof
```

### Simple Benchmarking

```go
func BenchmarkView(b *testing.B) {
    m := newTestModel()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

// Run: go test -bench=BenchmarkView
// Target: < 1ms per iteration
```

---

## 🆘 Getting Help

### Before Asking for Help

1. **Check this troubleshooting guide**
2. **Read error message carefully** - often tells you exactly what's wrong
3. **Search GitHub issues**: [github.com/phoenix-tui/phoenix/issues](https://github.com/phoenix-tui/phoenix/issues)
4. **Read relevant docs**:
   - [MIGRATION_FROM_BUBBLETEA.md](MIGRATION_FROM_BUBBLETEA.md)
   - [TESTING.md](TESTING.md)
   - [API documentation](https://pkg.go.dev/github.com/phoenix-tui/phoenix)

### How to Report a Bug

**Good Bug Report**:
```markdown
## Environment
- Phoenix version: v0.1.0-alpha
- Go version: 1.25
- OS: Windows 11 / Git Bash
- Terminal: Windows Terminal

## Expected Behavior
Input should clear after calling SetValue("")

## Actual Behavior
Input still shows old text after SetValue("")

## Minimal Reproduction
```go
m := Model{input: input.New()}
m.input.SetValue("test")
m.input.SetValue("")  // Doesn't clear
fmt.Println(m.input.Value())  // Prints "test"
```

## Additional Context
Migrating from Bubbletea where this worked fine.
```

### Where to Get Help

- **GitHub Issues**: [github.com/phoenix-tui/phoenix/issues](https://github.com/phoenix-tui/phoenix/issues)
- **Discussions**: [github.com/phoenix-tui/phoenix/discussions](https://github.com/phoenix-tui/phoenix/discussions)
- **Discord**: (coming soon)

---

## 📚 Appendix: Known Issues

### Current (v0.1.0-alpha)

| Issue | Severity | Workaround | ETA Fix |
|-------|----------|------------|---------|
| SetValue("") doesn't clear input | High | Create new input | beta.1 |
| Pointer vs value confusion | Medium | Use pointers | beta.1 |
| Terminal nil in tests | High | Use NullTerminal | ✅ Fixed (phoenix/testing) |
| Component API inconsistency | Medium | Check docs | beta.1 |

### Fixed in Latest Version

| Issue | Fixed In | Details |
|-------|----------|---------|
| Unicode width calculation | v0.1.0-alpha | Perfect emoji/CJK support |
| Terminal nil panics | v0.1.0-alpha | phoenix/testing package |

---

*Troubleshooting Guide Version: 1.0*
*Last Updated: 2025-10-19*
*Based on: Real-world GoSh migration issues and user reports*
