# TextArea Architecture

> **Domain-Driven Design** implementation for multiline text editing
>
> **Pattern**: DDD + Rich Models + Hexagonal Architecture
>
> **Version**: 1.0.0

---

## Table of Contents

1. [Overview](#overview)
2. [Layer Architecture](#layer-architecture)
3. [Domain Layer](#domain-layer)
4. [Application Layer](#application-layer)
5. [Infrastructure Layer](#infrastructure-layer)
6. [API Layer](#api-layer)
7. [Design Decisions](#design-decisions)
8. [Testing Strategy](#testing-strategy)

---

## Overview

TextArea follows **Domain-Driven Design (DDD)** principles with clear separation between:

- **Domain** - Pure business logic (text editing rules)
- **Application** - Use cases (future expansion)
- **Infrastructure** - Technical details (keybindings, rendering)
- **API** - Public interface (fluent builder pattern)

### Key Principles

1. **Rich Domain Models** - Models contain behavior, not just data
2. **Immutability** - All operations return new instances
3. **Layer Boundaries** - Domain never depends on infrastructure
4. **Testability** - Domain layer has 95%+ coverage target
5. **Type Safety** - Leverage Go 1.25+ features

---

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Public API Layer                       │
│  api/textarea.go - Fluent builder, Elm Architecture         │
└─────────────────────────────────────────────────────────────┘
                            ↓ uses
┌────────────────────┬─────────────────────┬──────────────────┐
│  Infrastructure    │    Application      │    Domain        │
│  (Technical)       │    (Use Cases)      │    (Business)    │
├────────────────────┼─────────────────────┼──────────────────┤
│ keybindings/       │ command/ (future)   │ model/           │
│ - emacs.go         │ - insert.go         │ - textarea.go    │
│ - vi.go (future)   │ - delete.go         │ - buffer.go      │
│                    │ - move.go           │ - cursor.go      │
│ renderer/          │                     │ - killring.go    │
│ - textarea.go      │ query/ (future)     │ - selection.go   │
│                    │ - content.go        │                  │
│                    │ - visible.go        │ value/           │
│                    │                     │ - position.go    │
│                    │                     │ - range.go       │
│                    │                     │                  │
│                    │                     │ service/         │
│                    │                     │ - navigation.go  │
│                    │                     │ - editing.go     │
└────────────────────┴─────────────────────┴──────────────────┘
```

### Dependency Rules

- **API** → **Infrastructure** → **Application** → **Domain**
- **Domain** never depends on other layers (pure business logic)
- **Application** may use domain services
- **Infrastructure** implements technical concerns

---

## Domain Layer

### Purpose

The domain layer contains **pure business logic** with zero external dependencies. It defines:

- What is a text buffer?
- How does cursor movement work?
- What are the rules for text editing?
- How does the kill ring behave?

### Components

#### 1. Value Objects (`domain/value/`)

**Position** - Immutable (row, col) tuple:

```go
type Position struct {
	row int
	col int
}

func (p Position) IsBefore(other Position) bool
func (p Position) IsAfter(other Position) bool
func (p Position) Equals(other Position) bool
```

**Range** - Immutable (start, end) tuple:

```go
type Range struct {
	start Position
	end   Position
}

func (r Range) Contains(pos Position) bool
func (r Range) IsEmpty() bool
func (r Range) IsSingleLine() bool
```

#### 2. Domain Models (`domain/model/`)

**Buffer** - Rich model for text storage:

```go
type Buffer struct {
	lines []string  // Text as array of lines
}

// Behavior (immutable operations)
func (b *Buffer) InsertChar(row, col int, ch rune) *Buffer
func (b *Buffer) DeleteChar(row, col int) *Buffer
func (b *Buffer) InsertNewline(row, col int) *Buffer
func (b *Buffer) DeleteLine(row int) (*Buffer, string)
func (b *Buffer) TextInRange(r value.Range) string
```

**Why rich?** Buffer encapsulates text editing rules (e.g., always keep at least one empty line, handle Unicode correctly).

**Cursor** - Rich model for position tracking:

```go
type Cursor struct {
	row int
	col int
}

func (c *Cursor) MoveTo(row, col int) *Cursor
func (c *Cursor) MoveBy(deltaRow, deltaCol int) *Cursor
```

**KillRing** - Emacs-style clipboard with history:

```go
type KillRing struct {
	items   []string
	maxSize int
	index   int
}

func (k *KillRing) Kill(text string) *KillRing
func (k *KillRing) Yank() string
func (k *KillRing) YankPop() *KillRing  // Rotate kill ring
```

**Selection** - Text selection state:

```go
type Selection struct {
	anchor value.Position  // Where selection started
	cursor value.Position  // Current cursor position
}

func (s *Selection) Range() value.Range  // Normalized range
```

**TextArea** - Aggregate root:

```go
type TextArea struct {
	// Core state
	buffer    *Buffer
	cursor    *Cursor
	selection *Selection
	killRing  *KillRing

	// Display config
	width, height int
	scrollRow, scrollCol int

	// Behavior config
	maxLines, maxChars int
	placeholder string
	wrap, readOnly bool
	showLineNumbers bool
}

// Public getters (CRITICAL for integration!)
func (t *TextArea) CursorPosition() (row, col int)
func (t *TextArea) Lines() []string
func (t *TextArea) Value() string
func (t *TextArea) ContentParts() (before, at, after string)

// Configuration methods (fluent builder)
func (t *TextArea) WithSize(width, height int) *TextArea
func (t *TextArea) WithMaxLines(max int) *TextArea
func (t *TextArea) WithReadOnly(readOnly bool) *TextArea
```

**Why TextArea is aggregate root?**
- Coordinates all components (buffer, cursor, killring, selection)
- Enforces invariants (e.g., cursor position always valid)
- Provides atomic operations (e.g., insert + move cursor)

#### 3. Domain Services (`domain/service/`)

**NavigationService** - Cursor movement logic:

```go
type NavigationService struct{}

func (s *NavigationService) MoveLeft(ta *model.TextArea) *model.TextArea
func (s *NavigationService) MoveRight(ta *model.TextArea) *model.TextArea
func (s *NavigationService) MoveUp(ta *model.TextArea) *model.TextArea
func (s *NavigationService) MoveDown(ta *model.TextArea) *model.TextArea
func (s *NavigationService) MoveToLineStart(ta *model.TextArea) *model.TextArea
func (s *NavigationService) MoveToLineEnd(ta *model.TextArea) *model.TextArea
func (s *NavigationService) ForwardWord(ta *model.TextArea) *model.TextArea
func (s *NavigationService) BackwardWord(ta *model.TextArea) *model.TextArea
```

**EditingService** - Text modification logic:

```go
type EditingService struct{}

func (s *EditingService) InsertChar(ta *model.TextArea, ch rune) *model.TextArea
func (s *EditingService) DeleteCharBackward(ta *model.TextArea) *model.TextArea
func (s *EditingService) DeleteCharForward(ta *model.TextArea) *model.TextArea
func (s *EditingService) InsertNewline(ta *model.TextArea) *model.TextArea
func (s *EditingService) KillLine(ta *model.TextArea) *model.TextArea
func (s *EditingService) Yank(ta *model.TextArea) *model.TextArea
```

**Why services?** Complex business logic that doesn't belong to a single entity.

---

## Application Layer

### Purpose

The application layer orchestrates domain services for specific use cases.

### Status

**Currently minimal** - Will expand when needed (future).

### Future Components

**Commands** (`application/command/`):

```go
// InsertCharCommand - Use case: insert character at cursor
type InsertCharCommand struct {
	char rune
}

func (c InsertCharCommand) Execute(ta *model.TextArea) (*model.TextArea, error)
```

**Queries** (`application/query/`):

```go
// GetVisibleLinesQuery - Use case: get visible lines for rendering
type GetVisibleLinesQuery struct {
	scrollRow int
	height    int
}

func (q GetVisibleLinesQuery) Execute(ta *model.TextArea) []string
```

---

## Infrastructure Layer

### Purpose

The infrastructure layer handles **technical concerns** (not business logic):

- Keybinding translation (KeyMsg → domain operations)
- Rendering (TextArea → string output)
- Platform-specific code (future: clipboard, mouse)

### Components

#### 1. Keybindings (`infrastructure/keybindings/`)

**EmacsKeybindings** - Translate KeyMsg to domain operations:

```go
type EmacsKeybindings struct {
	navigation *service.NavigationService
	editing    *service.EditingService
}

func (e *EmacsKeybindings) Handle(msg tea.KeyMsg, ta *model.TextArea) (*model.TextArea, tea.Cmd) {
	// Map Ctrl+A → MoveToLineStart
	// Map Ctrl+K → KillLine
	// Map Alt+F → ForwardWord
	// etc.
}
```

**Future: ViKeybindings** - Vi-style keybindings:

```go
type ViKeybindings struct {
	mode ViMode  // Normal, Insert, Visual
}
```

#### 2. Renderer (`infrastructure/renderer/`)

**TextAreaRenderer** - Render TextArea to string:

```go
type TextAreaRenderer struct{}

func (r *TextAreaRenderer) Render(ta *model.TextArea) string {
	// Get visible lines
	// Render line numbers (if enabled)
	// Render cursor
	// Return formatted string
}
```

---

## API Layer

### Purpose

The API layer provides the **public interface** for external users. It:

- Wraps domain models in user-friendly API
- Implements Elm Architecture (Model-View-Update)
- Provides fluent builder pattern
- Hides internal complexity

### Component

**TextArea** (`api/textarea.go`):

```go
// Public API
type TextArea struct {
	model       *model.TextArea           // Domain model
	keybindings KeybindingMode            // Emacs, Vi, etc.
	renderer    *renderer.TextAreaRenderer // Renderer
}

// Fluent builder
func New() TextArea
func (t TextArea) Size(width, height int) TextArea
func (t TextArea) Placeholder(text string) TextArea
func (t TextArea) ReadOnly(readOnly bool) TextArea

// State access (CRITICAL for integration!)
func (t TextArea) Value() string
func (t TextArea) CursorPosition() (row, col int)
func (t TextArea) ContentParts() (before, at, after string)

// Elm Architecture
func (t TextArea) Init() tea.Cmd
func (t TextArea) Update(msg tea.Msg) (TextArea, tea.Cmd)
func (t TextArea) View() string
```

### Why separate public API from domain model?

1. **Encapsulation** - Hide internal domain complexity
2. **Stability** - Public API can remain stable while domain evolves
3. **Usability** - Fluent API is easier to use than domain models directly
4. **Testability** - Can test domain layer without public API

---

## Design Decisions

### 1. Immutability

**Decision**: All operations return new instances (no mutation).

**Rationale**:
- Fits Elm Architecture (Model-View-Update)
- Easier to test (no hidden state changes)
- Prevents bugs from shared state
- Safe for concurrent access (future)

**Example**:

```go
// ✅ Immutable (correct)
newTA := ta.WithSize(80, 24)

// ❌ Mutable (wrong)
ta.SetSize(80, 24)  // Modifies in place
```

### 2. Rich Domain Models

**Decision**: Models contain behavior, not just data.

**Rationale**:
- Business logic belongs in domain layer
- Easier to test (logic is in models, not scattered)
- Prevents anemic models (data bags)

**Example**:

```go
// ✅ Rich model
type Buffer struct {
	lines []string
}

func (b *Buffer) InsertChar(row, col int, ch rune) *Buffer {
	// Business logic: how to insert character
	// Handles: bounds checking, Unicode, newlines
}

// ❌ Anemic model
type Buffer struct {
	Lines []string  // Exposed for external manipulation
}
// Logic scattered across services, hard to test
```

### 3. Public Cursor API

**Decision**: Expose `CursorPosition()`, `ContentParts()` publicly.

**Rationale**:
- **Critical for GoSh** - Enables syntax highlighting
- **Flexible** - External tools can work with cursor
- **Simple** - Just (row, col) integers

**Example**:

```go
// GoSh syntax highlighting
before, at, after := ta.ContentParts()
highlightedBefore := highlighter.Highlight(before)
highlightedAfter := highlighter.Highlight(after)
cursorStyled := style.New().Reverse(true).Render(at)
```

### 4. Kill Ring

**Decision**: Implement full Emacs-style kill ring (not just clipboard).

**Rationale**:
- **Emacs standard** - Ctrl+K/Ctrl+Y workflow
- **History** - Keep last N kills (default 10)
- **Power users** - Expected by terminal users

**Example**:

```go
// Ctrl+K three times
ta = editing.KillLine(ta)  // Kill "line1"
ta = editing.KillLine(ta)  // Kill "line2"
ta = editing.KillLine(ta)  // Kill "line3"

// Ctrl+Y
ta = editing.Yank(ta)      // Paste "line3" (most recent)

// Alt+Y (future)
ta = editing.YankPop(ta)   // Cycle to "line2"
ta = editing.YankPop(ta)   // Cycle to "line1"
```

### 5. Layer Boundaries

**Decision**: Domain layer has zero dependencies on infrastructure.

**Rationale**:
- **Testability** - Domain tests don't need mocks
- **Portability** - Can swap infrastructure (e.g., Vi keybindings)
- **Clarity** - Clear separation of concerns

**Example**:

```go
// ✅ Correct: Infrastructure depends on domain
package keybindings

import "domain/model"
import "domain/service"

// ❌ Wrong: Domain depends on infrastructure
package model

import "infrastructure/keybindings"  // NEVER!
```

---

## Testing Strategy

### Coverage Targets

- **Domain layer**: **95%+** (pure logic, easy to test)
- **Application layer**: **90%+** (use cases with mocked infrastructure)
- **Infrastructure layer**: **85%+** (integration tests)
- **Overall**: **90%+**

### Test Types

#### 1. Unit Tests (Domain Layer)

Test pure business logic:

```go
func TestBuffer_InsertChar(t *testing.T) {
	buf := model.NewBufferFromString("hello")
	result := buf.InsertChar(0, 2, 'X')

	assert.Equal(t, "heXllo", result.String())
	assert.Equal(t, "hello", buf.String())  // Original unchanged (immutability)
}
```

#### 2. Property-Based Tests

Test invariants:

```go
func TestBuffer_Invariants(t *testing.T) {
	quick.Check(func(ops []Operation) bool {
		buf := model.NewBuffer()

		// Apply random operations
		for _, op := range ops {
			buf = op.Apply(buf)
		}

		// Invariant: Always at least one line
		return buf.LineCount() >= 1
	}, nil)
}
```

#### 3. Integration Tests

Test layer interaction:

```go
func TestTextArea_EmacsWorkflow(t *testing.T) {
	ta := api.New()

	// Type "hello"
	for _, ch := range "hello" {
		ta, _ = ta.Update(tea.KeyMsg{Rune: ch})
	}

	// Ctrl+A (move to start)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyCtrlA})

	// Verify cursor at start
	row, col := ta.CursorPosition()
	assert.Equal(t, 0, row)
	assert.Equal(t, 0, col)
}
```

#### 4. Example-Based Tests

Test documentation examples work:

```go
func TestExample_BasicUsage(t *testing.T) {
	// Copy example from README
	ta := api.New().
		Size(80, 24).
		Placeholder("Type something...")

	assert.Equal(t, 80, ta.Width())
	assert.Equal(t, 24, ta.Height())
}
```

---

## Future Enhancements

### Phase 1 (Current)
- [x] Domain layer (models, services)
- [x] Infrastructure (Emacs keybindings, renderer)
- [x] Public API (fluent builder)
- [x] Basic tests
- [x] Documentation

### Phase 2 (Next)
- [ ] 95%+ test coverage for domain layer
- [ ] Integration tests for public API
- [ ] Performance benchmarking
- [ ] More examples (Emacs demo, GoSh integration)

### Phase 3 (Future)
- [ ] Selection support (visual mode)
- [ ] Vi keybindings
- [ ] Application layer (commands, queries)
- [ ] Undo/redo stack
- [ ] Search & replace
- [ ] Syntax highlighting hooks

---

## References

### DDD Resources
- **Eric Evans** - "Domain-Driven Design" (Blue Book)
- **Vaughn Vernon** - "Implementing Domain-Driven Design" (Red Book)

### Go Best Practices
- **Go Wiki** - [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Effective Go** - [golang.org/doc/effective_go](https://golang.org/doc/effective_go)

### Phoenix Framework
- **ARCHITECTURE.md** - Phoenix framework architecture
- **MASTER_PLAN.md** - Strategic vision
- **API_DESIGN.md** - API principles

---

*Built with Domain-Driven Design for production use.*
*Part of Phoenix TUI Framework v0.1.0-alpha*
