// Package components provides rich, reusable UI components for Phoenix TUI framework.
//
// # Overview
//
// Package components is a comprehensive library of terminal UI widgets:
//   - 6 production-ready components (Input, List, Viewport, Table, Modal, Progress)
//   - Built on phoenix/tea (Elm Architecture pattern)
//   - DDD architecture (domain logic separate from presentation)
//   - Type-safe with Go generics
//   - Unicode-aware (perfect CJK and emoji support)
//   - Zero external TUI dependencies (pure Phoenix)
//
// # Features
//
//   - Rich component library (input fields, lists, tables, modals, progress bars)
//   - Fluent builder API (chainable method calls for styling)
//   - Keyboard navigation (arrow keys, tab, enter, escape)
//   - Mouse support (click, drag, scroll, hover)
//   - Validation (real-time input validation)
//   - Theming (customizable colors, borders, styles)
//   - 94.5% average test coverage (production-ready)
//   - Full Unicode support (emoji, CJK, combining characters)
//
// # Components
//
// Input - Single-line text input:
//
//	import input "github.com/phoenix-tui/phoenix/components/input/api"
//
//	field := input.New(40).
//		Placeholder("Enter name...").
//		Validate(func(s string) error {
//			if len(s) < 3 {
//				return errors.New("too short")
//			}
//			return nil
//		})
//
// List - Selectable item list:
//
//	import list "github.com/phoenix-tui/phoenix/components/list/api"
//
//	items := []string{"Item 1", "Item 2", "Item 3"}
//	l := list.New(items, 10).
//		Title("Select an item").
//		ShowNumbers(true)
//
// Viewport - Scrollable content area:
//
//	import viewport "github.com/phoenix-tui/phoenix/components/viewport/api"
//
//	vp := viewport.New(80, 24).
//		SetContent("Long content...\n" + strings.Repeat("Line\n", 100))
//
// Table - Data table with columns:
//
//	import table "github.com/phoenix-tui/phoenix/components/table/api"
//
//	tbl := table.New(
//		[]string{"Name", "Age", "City"},
//		[][]string{
//			{"Alice", "30", "NYC"},
//			{"Bob", "25", "SF"},
//		},
//	)
//
// Modal - Overlay dialog:
//
//	import modal "github.com/phoenix-tui/phoenix/components/modal/api"
//
//	m := modal.New("Confirm Action", "Are you sure?").
//		Buttons([]string{"Yes", "No"})
//
// Progress - Progress indicator:
//
//	import progress "github.com/phoenix-tui/phoenix/components/progress/api"
//
//	p := progress.New(100).
//		SetValue(50).
//		ShowPercentage(true)
//
// # Architecture
//
// Each component follows DDD structure:
//
//	components/
//	├── input/              # Single-line text input
//	│   ├── internal/domain/   # Business logic (90%+ coverage)
//	│   ├── examples/          # Usage examples
//	│   └── input.go           # Public API
//	├── list/              # Selectable list
//	│   ├── internal/domain/
//	│   ├── examples/
//	│   └── list.go
//	├── viewport/          # Scrollable area
//	├── table/             # Data table
//	├── modal/             # Overlay dialog
//	└── progress/          # Progress bar
//
// All components implement the tea.Model interface:
//   - Init() tea.Cmd
//   - Update(tea.Msg) (Model, tea.Cmd)
//   - View() string
//
// # Integration with phoenix/tea
//
// Components are tea.Model implementations:
//
//	type AppModel struct {
//		input input.Model
//	}
//
//	func (m AppModel) Init() tea.Cmd {
//		return m.input.Init()
//	}
//
//	func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//		var cmd tea.Cmd
//		m.input, cmd = m.input.Update(msg)
//		return m, cmd
//	}
//
//	func (m AppModel) View() string {
//		return m.input.View()
//	}
//
// # Performance
//
// Components are optimized for responsiveness:
//   - Efficient rendering (only changed content updated)
//   - Minimal allocations (object pooling where applicable)
//   - Unicode-correct width calculations (no visual glitches)
//   - Keyboard navigation <10ms latency
//   - 94.5% average test coverage across all components
//
// Component-specific coverage:
//   - input: 90.0%
//   - list: 94.7%
//   - viewport: 94.5%
//   - table: 92.0%
//   - modal: 96.5%
//   - progress: 98.5%
//
// See each component's README for detailed API documentation and examples.
package components
