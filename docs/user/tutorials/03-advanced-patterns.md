# Tutorial 3: Advanced Patterns

> **Time to Complete**: 60-75 minutes
> **Difficulty Level**: Advanced
> **Prerequisites**: Tutorials 1-2 complete, solid understanding of Elm Architecture
> **What You'll Build**: Note-taking app with mouse, clipboard, layout, and custom components

---

## Table of Contents

1. [What You'll Learn](#what-youll-learn)
2. [Prerequisites](#prerequisites)
3. [Project Overview](#project-overview)
4. [Part 1: Custom Component Architecture](#part-1-custom-component-architecture-15-minutes)
5. [Part 2: Mouse Integration](#part-2-mouse-integration-15-minutes)
6. [Part 3: Clipboard Operations](#part-3-clipboard-operations-10-minutes)
7. [Part 4: Flexbox Layout System](#part-4-flexbox-layout-system-10-minutes)
8. [Part 5: Complex State Management](#part-5-complex-state-management-10-minutes)
9. [Part 6: Performance Optimization](#part-6-performance-optimization-10-minutes)
10. [Putting It All Together](#putting-it-all-together)
11. [Advanced Patterns Reference](#advanced-patterns-reference)
12. [Exercises](#exercises)
13. [Common Issues](#common-issues)
14. [Summary](#summary)

---

## What You'll Learn

By the end of this tutorial, you will:

- Build **custom reusable components** following Phoenix patterns
- Handle **mouse events** (click, drag, hover, scroll)
- Implement **clipboard operations** (copy, paste, cut)
- Use **Flexbox layout** for responsive UIs
- Manage **complex state** across multiple component levels
- **Optimize performance** for large datasets
- Use **async commands** for non-blocking operations
- Build **production-ready** TUI applications

---

## Prerequisites

### Knowledge Requirements

- **Tutorials 1-2 complete** - Solid MVU understanding
- **Go advanced** - Interfaces, generics, channels
- **Component patterns** - Composition, delegation
- **State management** - Immutable updates

### Software Requirements

```bash
# Verify setup
go version  # 1.25+

# Install all Phoenix packages
go get github.com/phoenix-tui/phoenix/tea
go get github.com/phoenix-tui/phoenix/style
go get github.com/phoenix-tui/phoenix/components/input
go get github.com/phoenix-tui/phoenix/components/list
go get github.com/phoenix-tui/phoenix/mouse
go get github.com/phoenix-tui/phoenix/clipboard
```

---

## Project Overview

We're building a **note-taking application** with advanced features:

**Features:**
- Multiple notes with tabs (switchable)
- Rich text editor with mouse support
- Click to position cursor
- Select text with mouse drag
- Copy/paste with Ctrl+C/Ctrl+V
- Responsive layout (adjusts to terminal size)
- Auto-save (async command)
- Syntax highlighting (custom component)

**Final UI Preview:**

```
╔═══════════════════════════════════════════════════════════════╗
║  Notes                                         [x] Clipboard   ║
╟───┬─────────┬─────────┬─────────┬──────────────────────────────╢
║ 📝│ Note 1  │ Note 2  │ Note 3  │                             ║
╟───┴─────────┴─────────┴─────────┴──────────────────────────────╢
║                                                                 ║
║  Phoenix TUI Framework is amazing!                              ║
║  [Selected text highlighted]                                    ║
║                                                                 ║
║  I can click to position cursor,                                ║
║  drag to select text,                                           ║
║  and copy/paste with keyboard.                                  ║
║                                                                 ║
╟─────────────────────────────────────────────────────────────────╢
║  Ctrl+C: Copy  Ctrl+V: Paste  Ctrl+S: Save  q: Quit           ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## Part 1: Custom Component Architecture (15 minutes)

Let's start by building a **custom TabBar component** following Phoenix DDD patterns.

### Step 1: Component Structure

Create `components/tabbar/tabbar.go`:

```go
package tabbar

import (
    "fmt"
    "strings"

    "github.com/phoenix-tui/phoenix/style"
    "github.com/phoenix-tui/phoenix/tea"
)

// Tab represents a single tab.
type Tab struct {
    ID    string
    Label string
    Icon  string  // Optional emoji/icon
}

// TabBar is a custom component for tab navigation.
type TabBar struct {
    // State
    tabs          []Tab
    activeIndex   int
    hoveredIndex  int  // -1 if no hover
    width         int

    // Styles
    activeStyle   style.Style
    inactiveStyle style.Style
    hoverStyle    style.Style
}

// New creates a new TabBar.
func New(tabs []Tab) *TabBar {
    return &TabBar{
        tabs:          tabs,
        activeIndex:   0,
        hoveredIndex:  -1,
        width:         80,
        activeStyle:   defaultActiveStyle(),
        inactiveStyle: defaultInactiveStyle(),
        hoverStyle:    defaultHoverStyle(),
    }
}

// defaultActiveStyle creates the active tab style.
func defaultActiveStyle() style.Style {
    return style.New().
        Foreground(style.RGB(100, 150, 255)).
        Bold(true).
        Underline(true)
}

// defaultInactiveStyle creates the inactive tab style.
func defaultInactiveStyle() style.Style {
    return style.New().
        Foreground(style.Color256(240))
}

// defaultHoverStyle creates the hover tab style.
func defaultHoverStyle() style.Style {
    return style.New().
        Foreground(style.RGB(150, 180, 255)).
        Underline(true)
}
```

### Step 2: Component API (Fluent Interface)

Add configuration methods:

```go
// Width sets the total width of the tab bar.
func (t *TabBar) Width(width int) *TabBar {
    t.width = width
    return t
}

// ActiveIndex sets the active tab index.
func (t *TabBar) ActiveIndex(index int) *TabBar {
    if index >= 0 && index < len(t.tabs) {
        t.activeIndex = index
    }
    return t
}

// ActiveStyle sets the style for active tabs.
func (t *TabBar) ActiveStyle(s style.Style) *TabBar {
    t.activeStyle = s
    return t
}

// InactiveStyle sets the style for inactive tabs.
func (t *TabBar) InactiveStyle(s style.Style) *TabBar {
    t.inactiveStyle = s
    return t
}

// HoverStyle sets the style for hovered tabs.
func (t *TabBar) HoverStyle(s style.Style) *TabBar {
    t.hoverStyle = s
    return t
}

// GetActiveTab returns the currently active tab.
func (t *TabBar) GetActiveTab() Tab {
    if t.activeIndex >= 0 && t.activeIndex < len(t.tabs) {
        return t.tabs[t.activeIndex]
    }
    return Tab{}
}
```

### Step 3: Implement Init/Update/View

```go
// Init implements tea.Model.
func (t *TabBar) Init() tea.Cmd {
    return nil
}

// Update implements tea.Model.
func (t *TabBar) Update(msg tea.Msg) (*TabBar, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "left", "h":
            // Previous tab
            if t.activeIndex > 0 {
                t.activeIndex--
            }
            return t, nil

        case "right", "l":
            // Next tab
            if t.activeIndex < len(t.tabs)-1 {
                t.activeIndex++
            }
            return t, nil

        case "1", "2", "3", "4", "5", "6", "7", "8", "9":
            // Jump to tab by number
            index := int(msg.Rune - '1')
            if index >= 0 && index < len(t.tabs) {
                t.activeIndex = index
            }
            return t, nil
        }

    case tea.MouseMsg:
        // Handle mouse clicks on tabs
        if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
            // Calculate which tab was clicked
            clickedTab := t.getTabAtPosition(msg.X)
            if clickedTab >= 0 {
                t.activeIndex = clickedTab
                return t, nil
            }
        }

        // Handle mouse hover
        if msg.Action == tea.MouseActionMotion {
            t.hoveredIndex = t.getTabAtPosition(msg.X)
            return t, nil
        }
    }

    return t, nil
}

// View implements tea.Model.
func (t *TabBar) View() string {
    var b strings.Builder

    // Calculate tab width (distribute evenly)
    tabWidth := t.width / len(t.tabs)

    for i, tab := range t.tabs {
        // Determine style
        var s style.Style
        if i == t.activeIndex {
            s = t.activeStyle
        } else if i == t.hoveredIndex {
            s = t.hoverStyle
        } else {
            s = t.inactiveStyle
        }

        // Build tab label
        label := fmt.Sprintf(" %s %s ", tab.Icon, tab.Label)

        // Pad to tab width
        if len(label) < tabWidth {
            padding := (tabWidth - len(label)) / 2
            label = strings.Repeat(" ", padding) + label + strings.Repeat(" ", tabWidth-len(label)-padding)
        } else if len(label) > tabWidth {
            label = label[:tabWidth-3] + "..."
        }

        // Render with style
        b.WriteString(style.Render(s, label))

        // Separator
        if i < len(t.tabs)-1 {
            b.WriteString("│")
        }
    }

    return b.String()
}

// getTabAtPosition calculates which tab is at the given X coordinate.
func (t *TabBar) getTabAtPosition(x int) int {
    tabWidth := t.width / len(t.tabs)
    index := x / tabWidth
    if index >= 0 && index < len(t.tabs) {
        return index
    }
    return -1
}
```

### Step 4: Using Custom Component

In your main app:

```go
package main

import (
    "example.com/notes/components/tabbar"
    "github.com/phoenix-tui/phoenix/tea"
)

type Model struct {
    tabBar *tabbar.TabBar
    // ... other fields
}

func initialModel() Model {
    tabs := []tabbar.Tab{
        {ID: "note1", Label: "Note 1", Icon: "📝"},
        {ID: "note2", Label: "Note 2", Icon: "📄"},
        {ID: "note3", Label: "Note 3", Icon: "✏️"},
    }

    return Model{
        tabBar: tabbar.New(tabs).Width(80),
    }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Delegate to tab bar
    updatedTabBar, cmd := m.tabBar.Update(msg)
    m.tabBar = updatedTabBar
    return m, cmd
}

func (m Model) View() string {
    return m.tabBar.View() + "\n\n" + m.renderContent()
}
```

**Key Custom Component Patterns:**

1. **Pointer receiver** for stateful components
2. **Fluent API** for configuration (Width(), Style(), etc.)
3. **GetX() methods** for accessing internal state
4. **Implement tea.Model** (Init, Update, View)
5. **Handle both keyboard and mouse** events

---

## Part 2: Mouse Integration (15 minutes)

Now let's add mouse support to our note editor.

### Step 1: Enable Mouse Events

```go
func main() {
    p := tea.New(
        initialModel(),
        tea.WithAltScreen[Model](),
        tea.WithMouseAllMotion[Model](),  // Enable mouse!
    )

    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Step 2: Handle Mouse Messages

```go
// TextEditor component with mouse support
type TextEditor struct {
    content      string
    cursorPos    int
    selectionStart int  // -1 if no selection
    selectionEnd   int
    bounds       Rect   // X, Y, Width, Height
    isDragging   bool
}

type Rect struct {
    X, Y, Width, Height int
}

func (e *TextEditor) Update(msg tea.Msg) (*TextEditor, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Check if click is within bounds
        if !e.isWithinBounds(msg.X, msg.Y) {
            return e, nil
        }

        switch msg.Action {
        case tea.MouseActionPress:
            if msg.Button == tea.MouseButtonLeft {
                // Click to position cursor
                e.cursorPos = e.getCursorPosFromClick(msg.X, msg.Y)
                e.selectionStart = -1  // Clear selection
                e.selectionEnd = -1
                e.isDragging = true
                return e, nil
            }

        case tea.MouseActionRelease:
            if msg.Button == tea.MouseButtonLeft {
                e.isDragging = false
                return e, nil
            }

        case tea.MouseActionMotion:
            if e.isDragging {
                // Drag to select text
                newPos := e.getCursorPosFromClick(msg.X, msg.Y)
                if e.selectionStart == -1 {
                    e.selectionStart = e.cursorPos
                }
                e.selectionEnd = newPos
                e.cursorPos = newPos
                return e, nil
            }

            // Hover effects (optional)
            // e.hoveredPos = e.getCursorPosFromClick(msg.X, msg.Y)
            return e, nil

        case tea.MouseActionPress:
            if msg.Button == tea.MouseButtonWheelUp {
                // Scroll up
                // e.scrollOffset--
                return e, nil
            } else if msg.Button == tea.MouseButtonWheelDown {
                // Scroll down
                // e.scrollOffset++
                return e, nil
            }
        }

    case tea.KeyMsg:
        // Keyboard handling...
        return e.handleKeyboard(msg)
    }

    return e, nil
}

// isWithinBounds checks if coordinates are within component bounds.
func (e *TextEditor) isWithinBounds(x, y int) bool {
    return x >= e.bounds.X && x < e.bounds.X+e.bounds.Width &&
           y >= e.bounds.Y && y < e.bounds.Y+e.bounds.Height
}

// getCursorPosFromClick calculates cursor position from mouse coordinates.
func (e *TextEditor) getCursorPosFromClick(x, y int) int {
    // Convert screen coordinates to text position
    // Simplified: assumes monospace font, no wrapping
    relativeX := x - e.bounds.X
    relativeY := y - e.bounds.Y

    lines := strings.Split(e.content, "\n")
    if relativeY >= len(lines) {
        // Click below content - go to end
        return len(e.content)
    }

    // Find position in clicked line
    line := lines[relativeY]
    if relativeX >= len(line) {
        // Click past end of line
        return e.getLineStartPos(relativeY) + len(line)
    }

    return e.getLineStartPos(relativeY) + relativeX
}

// getLineStartPos returns the byte offset of the start of line n.
func (e *TextEditor) getLineStartPos(lineNum int) int {
    lines := strings.Split(e.content, "\n")
    pos := 0
    for i := 0; i < lineNum && i < len(lines); i++ {
        pos += len(lines[i]) + 1  // +1 for newline
    }
    return pos
}
```

### Step 3: Visual Feedback for Selection

```go
func (e *TextEditor) View() string {
    if e.selectionStart == -1 {
        // No selection - render normally
        return e.renderContent()
    }

    // Render with selection highlighting
    start := min(e.selectionStart, e.selectionEnd)
    end := max(e.selectionStart, e.selectionEnd)

    before := e.content[:start]
    selected := e.content[start:end]
    after := e.content[end:]

    selectionStyle := style.New().
        Background(style.RGB(50, 100, 200)).
        Foreground(style.White)

    return before +
           style.Render(selectionStyle, selected) +
           after
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Mouse Event Types:**

```go
// Actions
tea.MouseActionPress    // Button pressed
tea.MouseActionRelease  // Button released
tea.MouseActionMotion   // Mouse moved

// Buttons
tea.MouseButtonLeft
tea.MouseButtonMiddle
tea.MouseButtonRight
tea.MouseButtonWheelUp
tea.MouseButtonWheelDown

// Modifiers
msg.Ctrl   // Ctrl key held
msg.Alt    // Alt key held
msg.Shift  // Shift key held
```

---

## Part 3: Clipboard Operations (10 minutes)

Add copy/paste functionality using Phoenix clipboard.

### Step 1: Install Clipboard Package

```bash
go get github.com/phoenix-tui/phoenix/clipboard
```

### Step 2: Implement Copy/Paste

```go
import (
    "github.com/phoenix-tui/phoenix/clipboard"
    "github.com/phoenix-tui/phoenix/tea"
)

// CopyCmd returns a command that copies text to clipboard.
func CopyCmd(text string) tea.Cmd {
    return func() tea.Msg {
        err := clipboard.WriteAll(text)
        if err != nil {
            return ClipboardErrorMsg{Err: err}
        }
        return ClipboardCopiedMsg{Text: text}
    }
}

// PasteCmd returns a command that reads from clipboard.
func PasteCmd() tea.Cmd {
    return func() tea.Msg {
        text, err := clipboard.ReadAll()
        if err != nil {
            return ClipboardErrorMsg{Err: err}
        }
        return ClipboardPastedMsg{Text: text}
    }
}

// Custom messages
type ClipboardCopiedMsg struct {
    Text string
}

type ClipboardPastedMsg struct {
    Text string
}

type ClipboardErrorMsg struct {
    Err error
}
```

### Step 3: Handle Copy/Paste in Update

```go
func (e *TextEditor) Update(msg tea.Msg) (*TextEditor, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Check for Ctrl+C (copy)
        if msg.String() == "ctrl+c" && e.hasSelection() {
            selectedText := e.getSelectedText()
            return e, CopyCmd(selectedText)
        }

        // Check for Ctrl+X (cut)
        if msg.String() == "ctrl+x" && e.hasSelection() {
            selectedText := e.getSelectedText()
            e.deleteSelection()
            return e, CopyCmd(selectedText)
        }

        // Check for Ctrl+V (paste)
        if msg.String() == "ctrl+v" {
            return e, PasteCmd()
        }

    case ClipboardCopiedMsg:
        // Show feedback (optional)
        e.statusMessage = "Copied to clipboard"
        return e, nil

    case ClipboardPastedMsg:
        // Insert pasted text at cursor
        e.insertText(msg.Text)
        return e, nil

    case ClipboardErrorMsg:
        // Handle error (optional)
        e.statusMessage = fmt.Sprintf("Clipboard error: %v", msg.Err)
        return e, nil
    }

    return e, nil
}

// hasSelection returns true if text is selected.
func (e *TextEditor) hasSelection() bool {
    return e.selectionStart != -1 && e.selectionEnd != -1 &&
           e.selectionStart != e.selectionEnd
}

// getSelectedText returns the currently selected text.
func (e *TextEditor) getSelectedText() string {
    if !e.hasSelection() {
        return ""
    }

    start := min(e.selectionStart, e.selectionEnd)
    end := max(e.selectionStart, e.selectionEnd)

    return e.content[start:end]
}

// deleteSelection removes the selected text.
func (e *TextEditor) deleteSelection() {
    if !e.hasSelection() {
        return
    }

    start := min(e.selectionStart, e.selectionEnd)
    end := max(e.selectionStart, e.selectionEnd)

    e.content = e.content[:start] + e.content[end:]
    e.cursorPos = start
    e.selectionStart = -1
    e.selectionEnd = -1
}

// insertText inserts text at cursor position.
func (e *TextEditor) insertText(text string) {
    // Delete selection first if any
    if e.hasSelection() {
        e.deleteSelection()
    }

    // Insert text
    e.content = e.content[:e.cursorPos] + text + e.content[e.cursorPos:]
    e.cursorPos += len(text)
}
```

**Clipboard Best Practices:**

1. **Always use commands** for clipboard operations (async)
2. **Handle errors gracefully** (clipboard might be unavailable)
3. **Provide visual feedback** (status message)
4. **Respect platform conventions** (Ctrl+C/V on Windows/Linux, Cmd+C/V on macOS)

---

## Part 4: Flexbox Layout System (10 minutes)

Phoenix will have a flexbox layout system (currently in development). Here's how it will work:

### Step 1: Basic Flexbox Layout

```go
import "github.com/phoenix-tui/phoenix/layout"

// Create a horizontal flex container
func (m Model) View() string {
    // Left panel (sidebar)
    sidebar := m.renderSidebar()

    // Right panel (content)
    content := m.renderContent()

    // Create flex layout: 20% sidebar, 80% content
    container := layout.NewFlex().
        Direction(layout.Row).  // Horizontal
        Children(
            layout.NewFlexItem(sidebar).Grow(1),   // 20%
            layout.NewFlexItem(content).Grow(4),   // 80%
        )

    return container.Render(m.windowWidth, m.windowHeight)
}
```

### Step 2: Complex Nested Layout

```go
// Complex layout: header, sidebar + content, footer
func (m Model) View() string {
    header := m.renderHeader()
    sidebar := m.renderSidebar()
    content := m.renderContent()
    footer := m.renderFooter()

    // Main container (vertical)
    mainLayout := layout.NewFlex().
        Direction(layout.Column).
        Children(
            // Header (fixed height)
            layout.NewFlexItem(header).
                Basis(3).  // 3 lines tall
                Grow(0).   // Don't grow
                Shrink(0), // Don't shrink

            // Middle section (horizontal, fills remaining space)
            layout.NewFlexItem(
                layout.NewFlex().
                    Direction(layout.Row).
                    Children(
                        layout.NewFlexItem(sidebar).
                            Basis(20).  // 20 cols wide
                            Grow(0),
                        layout.NewFlexItem(content).
                            Grow(1),  // Fill remaining width
                    ),
            ).Grow(1),  // Fill remaining height

            // Footer (fixed height)
            layout.NewFlexItem(footer).
                Basis(1).
                Grow(0).
                Shrink(0),
        )

    return mainLayout.Render(m.windowWidth, m.windowHeight)
}
```

### Step 3: Responsive Layout

```go
// Layout adapts to terminal size
func (m Model) View() string {
    // Small terminal: stack vertically
    if m.windowWidth < 60 {
        return m.renderVerticalLayout()
    }

    // Large terminal: side-by-side
    return m.renderHorizontalLayout()
}

func (m Model) renderVerticalLayout() string {
    return layout.NewFlex().
        Direction(layout.Column).
        Children(
            layout.NewFlexItem(m.renderPanel1()),
            layout.NewFlexItem(m.renderPanel2()),
            layout.NewFlexItem(m.renderPanel3()),
        ).
        Render(m.windowWidth, m.windowHeight)
}

func (m Model) renderHorizontalLayout() string {
    return layout.NewFlex().
        Direction(layout.Row).
        Children(
            layout.NewFlexItem(m.renderPanel1()).Grow(1),
            layout.NewFlexItem(m.renderPanel2()).Grow(1),
            layout.NewFlexItem(m.renderPanel3()).Grow(1),
        ).
        Render(m.windowWidth, m.windowHeight)
}
```

**Flexbox Properties:**

- **Direction**: Row (horizontal) or Column (vertical)
- **Grow**: How much to grow when extra space available
- **Shrink**: How much to shrink when space is tight
- **Basis**: Initial size before grow/shrink
- **Align**: Cross-axis alignment (start, center, end)
- **Justify**: Main-axis alignment (start, center, end, space-between)

---

## Part 5: Complex State Management (10 minutes)

Managing state across multiple component levels.

### Step 1: State Tree Structure

```go
// Root Model
type Model struct {
    // UI State
    activeView   string  // "editor", "settings", "help"
    windowWidth  int
    windowHeight int

    // Components
    tabBar       *TabBar
    editor       *TextEditor
    sidebar      *Sidebar
    statusBar    *StatusBar

    // Application State
    notes        []Note
    activeNote   int
    settings     Settings
    clipboard    ClipboardState

    // Transient State
    isDirty      bool
    lastSaved    time.Time
    errorMessage string
}

type Note struct {
    ID       string
    Title    string
    Content  string
    Tags     []string
    Modified time.Time
}

type Settings struct {
    Theme         string
    AutoSave      bool
    AutoSaveDelay time.Duration
    FontSize      int
}

type ClipboardState struct {
    LastCopied   string
    CopiedAt     time.Time
}
```

### Step 2: State Access Patterns

```go
// Getters for nested state
func (m Model) GetActiveNote() *Note {
    if m.activeNote >= 0 && m.activeNote < len(m.notes) {
        return &m.notes[m.activeNote]
    }
    return nil
}

// Setters for nested state (immutable)
func (m Model) SetActiveNote(note Note) Model {
    if m.activeNote >= 0 && m.activeNote < len(m.notes) {
        m.notes[m.activeNote] = note
        m.isDirty = true
    }
    return m
}

// Batch updates
func (m Model) ApplySettings(s Settings) Model {
    m.settings = s
    m.isDirty = true

    // Propagate to components
    m.editor = m.editor.UpdateSettings(s)
    m.statusBar = m.statusBar.UpdateSettings(s)

    return m
}
```

### Step 3: Message Routing

```go
// Custom message types
type (
    NoteModifiedMsg struct {
        NoteID  string
        Content string
    }

    AutoSaveMsg struct{}

    SaveCompleteMsg struct {
        NoteID string
        Err    error
    }

    SettingsChangedMsg struct {
        Settings Settings
    }
)

// Update with message routing
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Global message handling
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.windowWidth = msg.Width
        m.windowHeight = msg.Height
        // Update all components with new size
        m.tabBar = m.tabBar.Width(msg.Width)
        m.editor = m.editor.SetBounds(Rect{
            X: 0, Y: 3,
            Width: msg.Width, Height: msg.Height - 5,
        })
        return m, nil

    case NoteModifiedMsg:
        // Update note content
        note := m.GetActiveNote()
        if note != nil {
            note.Content = msg.Content
            note.Modified = time.Now()
            m = m.SetActiveNote(*note)

            // Trigger auto-save if enabled
            if m.settings.AutoSave {
                cmds = append(cmds, m.scheduleAutoSave())
            }
        }
        return m, tea.Batch(cmds...)

    case AutoSaveMsg:
        // Save current note
        note := m.GetActiveNote()
        if note != nil && m.isDirty {
            cmds = append(cmds, m.saveNote(*note))
        }
        return m, tea.Batch(cmds...)

    case SaveCompleteMsg:
        if msg.Err == nil {
            m.isDirty = false
            m.lastSaved = time.Now()
            m.statusBar = m.statusBar.ShowMessage("Saved successfully")
        } else {
            m.errorMessage = fmt.Sprintf("Save failed: %v", msg.Err)
        }
        return m, nil

    case SettingsChangedMsg:
        m = m.ApplySettings(msg.Settings)
        return m, nil
    }

    // Route to components based on active view
    switch m.activeView {
    case "editor":
        updatedEditor, cmd := m.editor.Update(msg)
        m.editor = updatedEditor
        cmds = append(cmds, cmd)

    case "settings":
        // updatedSettings, cmd := m.settingsView.Update(msg)
        // m.settingsView = updatedSettings
        // cmds = append(cmds, cmd)
    }

    // Always update tab bar and status bar
    updatedTabBar, cmd := m.tabBar.Update(msg)
    m.tabBar = updatedTabBar
    cmds = append(cmds, cmd)

    updatedStatusBar, cmd := m.statusBar.Update(msg)
    m.statusBar = updatedStatusBar
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}
```

### Step 4: Async Commands

```go
// scheduleAutoSave returns a command that saves after a delay.
func (m Model) scheduleAutoSave() tea.Cmd {
    return func() tea.Msg {
        time.Sleep(m.settings.AutoSaveDelay)
        return AutoSaveMsg{}
    }
}

// saveNote returns a command that saves a note asynchronously.
func (m Model) saveNote(note Note) tea.Cmd {
    return func() tea.Msg {
        // Simulate async save operation
        err := saveNoteToFile(note)
        return SaveCompleteMsg{
            NoteID: note.ID,
            Err:    err,
        }
    }
}

// saveNoteToFile writes a note to disk.
func saveNoteToFile(note Note) error {
    filename := fmt.Sprintf("notes/%s.md", note.ID)
    return os.WriteFile(filename, []byte(note.Content), 0644)
}

// loadNote returns a command that loads a note asynchronously.
func LoadNoteCmd(noteID string) tea.Cmd {
    return func() tea.Msg {
        filename := fmt.Sprintf("notes/%s.md", noteID)
        content, err := os.ReadFile(filename)
        if err != nil {
            return NoteLoadErrorMsg{NoteID: noteID, Err: err}
        }
        return NoteLoadedMsg{
            NoteID:  noteID,
            Content: string(content),
        }
    }
}
```

---

## Part 6: Performance Optimization (10 minutes)

Techniques for building fast, responsive TUIs.

### Optimization 1: Lazy Rendering

```go
// Only render visible content
type Viewport struct {
    content      []string  // All lines
    scrollOffset int       // Current scroll position
    height       int       // Visible height
}

func (v *Viewport) View() string {
    // Only render visible lines
    start := v.scrollOffset
    end := min(v.scrollOffset+v.height, len(v.content))

    visibleLines := v.content[start:end]
    return strings.Join(visibleLines, "\n")
}
```

### Optimization 2: Memoization

```go
// Cache expensive computations
type MemoizedView struct {
    lastContent string
    lastWidth   int
    cachedView  string
}

func (m *MemoizedView) View(content string, width int) string {
    // Return cached view if inputs haven't changed
    if content == m.lastContent && width == m.lastWidth {
        return m.cachedView
    }

    // Recompute view
    view := m.computeExpensiveView(content, width)

    // Update cache
    m.lastContent = content
    m.lastWidth = width
    m.cachedView = view

    return view
}

func (m *MemoizedView) computeExpensiveView(content string, width int) string {
    // Expensive operations: word wrap, syntax highlighting, etc.
    // ...
    return result
}
```

### Optimization 3: Debouncing

```go
// Debounce rapid updates (e.g., typing)
type DebouncedUpdate struct {
    timer      *time.Timer
    delay      time.Duration
    lastUpdate time.Time
}

func (d *DebouncedUpdate) Debounce(fn func()) tea.Cmd {
    return func() tea.Msg {
        // Cancel previous timer
        if d.timer != nil {
            d.timer.Stop()
        }

        // Start new timer
        d.timer = time.AfterFunc(d.delay, fn)

        return nil
    }
}

// Usage in Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Update content immediately (for visual feedback)
        m.editor.InsertRune(msg.Rune)

        // Debounce expensive operations (syntax highlighting, save)
        return m, m.debouncer.Debounce(func() {
            // This runs after delay (e.g., 500ms)
            m.performExpensiveUpdate()
        })
    }
    return m, nil
}
```

### Optimization 4: Virtual Scrolling

```go
// Render only visible items in large lists
type VirtualList struct {
    items        []Item   // All items (could be 100,000+)
    scrollOffset int      // First visible item index
    visibleCount int      // Number of visible items
}

func (v *VirtualList) View() string {
    var b strings.Builder

    // Calculate visible range
    start := v.scrollOffset
    end := min(v.scrollOffset+v.visibleCount, len(v.items))

    // Render only visible items
    for i := start; i < end; i++ {
        b.WriteString(v.items[i].Render())
        b.WriteRune('\n')
    }

    return b.String()
}

// Handle scroll efficiently
func (v *VirtualList) ScrollDown() {
    if v.scrollOffset+v.visibleCount < len(v.items) {
        v.scrollOffset++
    }
}

func (v *VirtualList) ScrollUp() {
    if v.scrollOffset > 0 {
        v.scrollOffset--
    }
}
```

### Optimization 5: Batch Updates

```go
// Batch multiple state changes into one render
func (m Model) ApplyMultipleUpdates(updates []Update) (Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Apply all updates
    for _, update := range updates {
        switch u := update.(type) {
        case ContentUpdate:
            m.editor.SetContent(u.Content)
        case CursorUpdate:
            m.editor.SetCursor(u.Position)
        case SelectionUpdate:
            m.editor.SetSelection(u.Start, u.End)
        }
    }

    // Single render after all updates
    return m, tea.Batch(cmds...)
}
```

**Performance Checklist:**

- Use lazy rendering for large datasets
- Cache expensive computations (memoization)
- Debounce rapid events (typing, scrolling)
- Virtual scrolling for long lists (>1000 items)
- Batch multiple updates together
- Profile with Go's pprof if needed

---

## Putting It All Together

Here's a complete example integrating all advanced patterns:

```go
package main

import (
    "fmt"
    "os"
    "time"

    "example.com/notes/components/tabbar"
    "github.com/phoenix-tui/phoenix/clipboard"
    "github.com/phoenix-tui/phoenix/components/input"
    "github.com/phoenix-tui/phoenix/style"
    "github.com/phoenix-tui/phoenix/tea"
)

// Model represents the complete application state
type Model struct {
    // Components
    tabBar    *tabbar.TabBar
    editor    *TextEditor
    statusBar *StatusBar

    // State
    notes         []Note
    activeNote    int
    windowWidth   int
    windowHeight  int
    isDirty       bool

    // Styles
    styles        Styles
}

type Note struct {
    ID       string
    Title    string
    Content  string
    Modified time.Time
}

type TextEditor struct {
    content        string
    cursorPos      int
    selectionStart int
    selectionEnd   int
    bounds         Rect
    isDragging     bool
}

type Rect struct {
    X, Y, Width, Height int
}

type StatusBar struct {
    message     string
    messageTime time.Time
}

type Styles struct {
    Title   style.Style
    Editor  style.Style
    Status  style.Style
}

func initialModel() Model {
    // Create sample notes
    notes := []Note{
        {
            ID:      "note1",
            Title:   "Welcome",
            Content: "Welcome to Phoenix Notes!\n\nClick to edit, drag to select.",
        },
        {
            ID:      "note2",
            Title:   "Features",
            Content: "Mouse support\nClipboard integration\nCustom components",
        },
    }

    // Create tab bar
    tabs := []tabbar.Tab{
        {ID: "note1", Label: "Welcome", Icon: "📝"},
        {ID: "note2", Label: "Features", Icon: "✨"},
    }

    tabBar := tabbar.New(tabs).Width(80)

    // Create editor
    editor := &TextEditor{
        content:        notes[0].Content,
        cursorPos:      0,
        selectionStart: -1,
        selectionEnd:   -1,
        bounds:         Rect{X: 2, Y: 5, Width: 76, Height: 15},
        isDragging:     false,
    }

    // Create status bar
    statusBar := &StatusBar{
        message:     "Ready",
        messageTime: time.Now(),
    }

    return Model{
        tabBar:       tabBar,
        editor:       editor,
        statusBar:    statusBar,
        notes:        notes,
        activeNote:   0,
        windowWidth:  80,
        windowHeight: 24,
        isDirty:      false,
        styles:       newStyles(),
    }
}

func newStyles() Styles {
    return Styles{
        Title: style.New().
            Foreground(style.RGB(100, 150, 255)).
            Bold(true),
        Editor: style.New().
            Border(style.RoundedBorder).
            BorderForeground(style.RGB(80, 120, 200)).
            Padding(style.NewPadding(1, 2, 1, 2)),
        Status: style.New().
            Foreground(style.Color256(240)),
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.windowWidth = msg.Width
        m.windowHeight = msg.Height
        m.tabBar = m.tabBar.Width(msg.Width)
        m.editor.bounds.Width = msg.Width - 4
        m.editor.bounds.Height = msg.Height - 10
        return m, nil

    case tea.KeyMsg:
        // Global keys
        switch msg.String() {
        case "q", "esc":
            if m.isDirty {
                // TODO: Show save prompt
            }
            return m, tea.Quit()

        case "ctrl+s":
            // Save current note
            note := m.notes[m.activeNote]
            note.Content = m.editor.content
            note.Modified = time.Now()
            m.notes[m.activeNote] = note
            m.isDirty = false
            m.statusBar.message = "Saved"
            m.statusBar.messageTime = time.Now()
            return m, nil

        case "ctrl+c":
            // Copy selection
            if m.editor.hasSelection() {
                text := m.editor.getSelectedText()
                return m, clipboard.WriteCmd(text)
            }

        case "ctrl+v":
            // Paste from clipboard
            return m, clipboard.ReadCmd()
        }

    case clipboard.CopiedMsg:
        m.statusBar.message = "Copied"
        m.statusBar.messageTime = time.Now()
        return m, nil

    case clipboard.PastedMsg:
        m.editor.insertText(msg.Text)
        m.isDirty = true
        return m, nil
    }

    // Delegate to components
    updatedTabBar, cmd := m.tabBar.Update(msg)
    m.tabBar = updatedTabBar
    cmds = append(cmds, cmd)

    // Check if active tab changed
    activeTab := m.tabBar.GetActiveTab()
    for i, note := range m.notes {
        if note.ID == activeTab.ID && i != m.activeNote {
            // Switch note
            m.activeNote = i
            m.editor.content = note.Content
            m.editor.cursorPos = 0
            m.editor.selectionStart = -1
            m.editor.selectionEnd = -1
        }
    }

    updatedEditor, cmd := m.editor.Update(msg)
    m.editor = updatedEditor
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    // Title
    title := style.Render(m.styles.Title, "Phoenix Notes")

    // Tab bar
    tabs := m.tabBar.View()

    // Editor
    editorContent := m.editor.View()
    editor := style.Render(m.styles.Editor, editorContent)

    // Status bar
    statusText := fmt.Sprintf("%s | %s",
        m.statusBar.message,
        m.notes[m.activeNote].Modified.Format("15:04:05"))
    if m.isDirty {
        statusText += " [modified]"
    }
    status := style.Render(m.styles.Status, statusText)

    // Compose
    return fmt.Sprintf(
        "%s\n\n%s\n\n%s\n\n%s\n",
        title,
        tabs,
        editor,
        status,
    )
}

func main() {
    p := tea.New(
        initialModel(),
        tea.WithAltScreen[Model](),
        tea.WithMouseAllMotion[Model](),
    )

    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

// TextEditor methods (Update, View, helpers)
// ... (implementation from previous sections)
```

---

## Advanced Patterns Reference

### Pattern 1: Command Composition

```go
// Run multiple commands in parallel
cmd := tea.Batch(
    LoadDataCmd(),
    StartTimerCmd(),
    CheckUpdatesCmd(),
)

// Run commands in sequence
cmd := SequenceCmd(
    LoginCmd(),
    LoadProfileCmd(),
    LoadDashboardCmd(),
)
```

### Pattern 2: Middleware Pattern

```go
// Wrap Update with middleware (logging, analytics, etc.)
type Middleware func(tea.Msg, Model) (Model, tea.Cmd)

func LoggingMiddleware(msg tea.Msg, m Model) (Model, tea.Cmd) {
    log.Printf("Message: %T", msg)
    return m.Update(msg)
}

func AnalyticsMiddleware(msg tea.Msg, m Model) (Model, tea.Cmd) {
    // Track user interactions
    analytics.Track(msg)
    return m.Update(msg)
}

func (m Model) UpdateWithMiddleware(msg tea.Msg) (Model, tea.Cmd) {
    m, cmd := LoggingMiddleware(msg, m)
    m, cmd2 := AnalyticsMiddleware(msg, m)
    return m, tea.Batch(cmd, cmd2)
}
```

### Pattern 3: Pub/Sub Pattern

```go
// Event bus for component communication
type EventBus struct {
    subscribers map[string][]chan tea.Msg
}

func (b *EventBus) Subscribe(event string, ch chan tea.Msg) {
    b.subscribers[event] = append(b.subscribers[event], ch)
}

func (b *EventBus) Publish(event string, msg tea.Msg) {
    for _, ch := range b.subscribers[event] {
        ch <- msg
    }
}

// Usage
bus := &EventBus{subscribers: make(map[string][]chan tea.Msg)}
bus.Subscribe("note.modified", noteModifiedCh)
bus.Publish("note.modified", NoteModifiedMsg{ID: "note1"})
```

### Pattern 4: State Machine

```go
// Finite state machine for complex flows
type State int

const (
    StateIdle State = iota
    StateEditing
    StateSaving
    StateError
)

type StateMachine struct {
    current State
}

func (sm *StateMachine) Transition(newState State) {
    sm.current = newState
}

func (sm *StateMachine) CanTransition(newState State) bool {
    // Define valid transitions
    validTransitions := map[State][]State{
        StateIdle:    {StateEditing},
        StateEditing: {StateIdle, StateSaving},
        StateSaving:  {StateIdle, StateError},
        StateError:   {StateIdle},
    }

    for _, valid := range validTransitions[sm.current] {
        if valid == newState {
            return true
        }
    }
    return false
}
```

---

## Exercises

### Exercise 1: Add Undo/Redo

Implement undo/redo functionality with Ctrl+Z/Ctrl+Y.

<details>
<summary>Hint</summary>

Add history stacks:
```go
type TextEditor struct {
    content     string
    undoStack   []string
    redoStack   []string
    maxHistory  int
}

func (e *TextEditor) pushUndo() {
    e.undoStack = append(e.undoStack, e.content)
    if len(e.undoStack) > e.maxHistory {
        e.undoStack = e.undoStack[1:]
    }
    e.redoStack = []string{}  // Clear redo on new change
}
```
</details>

### Exercise 2: Add Syntax Highlighting

Highlight code syntax based on file extension.

<details>
<summary>Hint</summary>

Create a syntax highlighter:
```go
func highlightSyntax(content string, lang string) string {
    keywords := getKeywordsForLang(lang)
    // Apply styles to keywords
    for _, keyword := range keywords {
        content = strings.ReplaceAll(
            content,
            keyword,
            style.Render(keywordStyle, keyword),
        )
    }
    return content
}
```
</details>

### Exercise 3: Add Search and Replace

Implement Ctrl+F for search, Ctrl+H for replace.

<details>
<summary>Hint</summary>

Add search mode:
```go
type Model struct {
    searchMode    bool
    searchQuery   string
    searchResults []int  // Positions of matches
}

func (m Model) performSearch() Model {
    m.searchResults = []int{}
    offset := 0
    for {
        idx := strings.Index(m.editor.content[offset:], m.searchQuery)
        if idx == -1 {
            break
        }
        m.searchResults = append(m.searchResults, offset+idx)
        offset += idx + len(m.searchQuery)
    }
    return m
}
```
</details>

---

## Common Issues

### Issue 1: "Mouse events not received"

**Cause:** Mouse support not enabled.

**Solution:**

```go
// Enable mouse
p := tea.New(
    model,
    tea.WithMouseAllMotion[Model](),  // Add this!
)
```

### Issue 2: "Clipboard copy fails silently"

**Cause:** Clipboard command error not handled.

**Solution:**

```go
case clipboard.ErrorMsg:
    m.statusBar = m.statusBar.ShowError(fmt.Sprintf("Clipboard error: %v", msg.Err))
    return m, nil
```

### Issue 3: "Layout doesn't adapt to terminal resize"

**Cause:** Not handling WindowSizeMsg.

**Solution:**

```go
case tea.WindowSizeMsg:
    m.windowWidth = msg.Width
    m.windowHeight = msg.Height
    // Update component bounds
    m.updateComponentBounds()
    return m, nil
```

### Issue 4: "Slow performance with large text"

**Cause:** Rendering entire content every frame.

**Solution:**

Use viewport with lazy rendering:
```go
// Only render visible lines
visibleLines := content[scrollOffset:scrollOffset+height]
return strings.Join(visibleLines, "\n")
```

### Issue 5: "Component state not persisting"

**Cause:** Not reassigning updated component.

**Solution:**

```go
// WRONG
m.component.Update(msg)

// CORRECT
updated, cmd := m.component.Update(msg)
m.component = updated
return m, cmd
```

---

## Summary

Congratulations! You've mastered advanced Phoenix TUI patterns.

### What You Learned

**Custom Components:**
- DDD architecture for components
- Fluent API pattern
- Pointer vs value semantics
- Component lifecycle (Init/Update/View)

**Mouse Integration:**
- Enable mouse events
- Handle click, drag, hover, scroll
- Calculate positions from coordinates
- Visual feedback for interactions

**Clipboard:**
- Async clipboard operations
- Copy/paste/cut implementation
- Error handling
- Cross-platform support

**Layout:**
- Flexbox container system
- Responsive design patterns
- Nested layouts
- Dynamic sizing

**State Management:**
- Complex state trees
- Message routing strategies
- Async command patterns
- Batch updates

**Performance:**
- Lazy rendering (viewports)
- Memoization (caching)
- Debouncing (rapid events)
- Virtual scrolling (large lists)
- Batch updates

### Key Patterns

1. **Custom Component**:
   ```go
   type Component struct { state }
   func (c *Component) Init() tea.Cmd
   func (c *Component) Update(msg) (*Component, tea.Cmd)
   func (c *Component) View() string
   ```

2. **Mouse Handling**:
   ```go
   case tea.MouseMsg:
       if msg.Action == tea.MouseActionPress {
           // Handle click
       }
   ```

3. **Clipboard Async**:
   ```go
   return m, clipboard.WriteCmd(text)
   case clipboard.CopiedMsg:
       // Handle success
   ```

4. **State Delegation**:
   ```go
   updated, cmd := m.component.Update(msg)
   m.component = updated
   return m, cmd
   ```

### Architecture

```
Model (root)
  ├─ TabBar (custom component)
  ├─ TextEditor (custom component)
  │   ├─ Mouse handling
  │   ├─ Clipboard integration
  │   └─ Performance optimizations
  ├─ StatusBar (custom component)
  └─ Application state
      ├─ Notes
      ├─ Settings
      └─ Transient state
```

### Best Practices

1. **Component Design**:
   - Single responsibility
   - Immutable updates
   - Clear API boundaries
   - Proper error handling

2. **Performance**:
   - Profile before optimizing
   - Use lazy rendering for large content
   - Memoize expensive computations
   - Debounce rapid events

3. **State Management**:
   - Single source of truth
   - Predictable updates
   - Clear message types
   - Proper delegation

4. **User Experience**:
   - Visual feedback for all actions
   - Error messages that help
   - Responsive to window resize
   - Keyboard shortcuts + mouse

### Production Checklist

Before shipping your Phoenix TUI app:

- Handle all error cases gracefully
- Test on Windows, macOS, Linux
- Test with different terminal sizes
- Test with high data volumes
- Add keyboard shortcuts documentation
- Implement proper logging
- Add configuration file support
- Handle signals (SIGTERM, SIGINT)
- Write integration tests
- Profile for performance bottlenecks

### Next Steps

You're now ready to build production-quality Phoenix TUI applications!

**Resources:**
- [Phoenix API Documentation](../../api/)
- [Component Library](../../components/)
- [Example Applications](../../../examples/)
- [Phoenix GitHub](https://github.com/phoenix-tui/phoenix)

**Community:**
- GitHub Discussions
- Discord Server
- Stack Overflow (tag: phoenix-tui)

**Contributing:**
- Report bugs
- Suggest features
- Submit PRs
- Write tutorials

---

## Additional Resources

### Complete Example Apps

- [Phoenix Notes](../../../examples/notes/) - Full note-taking app
- [Phoenix Code](../../../examples/code/) - Code editor with syntax highlighting
- [Phoenix Chat](../../../examples/chat/) - Real-time chat client
- [Phoenix Monitor](../../../examples/monitor/) - System monitor dashboard

### API Reference

- [tea Package](../../api/tea.md) - Event loop
- [style Package](../../api/style.md) - Styling
- [mouse Package](../../api/mouse.md) - Mouse handling
- [clipboard Package](../../api/clipboard.md) - Clipboard operations
- [layout Package](../../api/layout.md) - Flexbox layout

### Guides

- [Performance Guide](../guides/performance.md)
- [Testing Guide](../guides/testing.md)
- [Deployment Guide](../guides/deployment.md)
- [Migration from Charm](../guides/migration.md)

---

*Tutorial created for Phoenix TUI Framework v0.1.0*
*Last updated: 2025-01-04*
*Time to complete: 60-75 minutes*
