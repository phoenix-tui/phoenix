# Phoenix TUI Framework

[![Go Version](https://img.shields.io/github/go-mod/go-version/phoenix-tui/phoenix)](https://github.com/phoenix-tui/phoenix)
[![Release](https://img.shields.io/github/v/release/phoenix-tui/phoenix?include_prereleases)](https://github.com/phoenix-tui/phoenix/releases)
[![CI](https://github.com/phoenix-tui/phoenix/actions/workflows/test.yml/badge.svg)](https://github.com/phoenix-tui/phoenix/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/phoenix-tui/phoenix)](https://goreportcard.com/report/github.com/phoenix-tui/phoenix)
[![Coverage](https://img.shields.io/badge/coverage-94.5%25-brightgreen)](https://github.com/phoenix-tui/phoenix)
[![License](https://img.shields.io/github/license/phoenix-tui/phoenix)](https://github.com/phoenix-tui/phoenix/blob/main/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/phoenix-tui/phoenix.svg)](https://pkg.go.dev/github.com/phoenix-tui/phoenix)

> **Multi-module monorepo** - 10 independent libraries. Full metrics in [CI](https://github.com/phoenix-tui/phoenix/actions).

> Next-generation Terminal User Interface framework for Go

**Status**: ✅ v0.1.0-beta.4 READY - API Modernization + Quality Improvements
**Progress**: 85% complete (17/20 weeks) 🎉
**Organization**: [github.com/phoenix-tui](https://github.com/phoenix-tui)
**Go Version**: 1.25+
**Test Coverage**: **93.5-96.4%** (Excellent across all modules)
**Performance**: 29,000 FPS (489x faster than 60 FPS target)
**Latest**: Industry-standard API patterns (Relica-style) + Cross-platform reliability 🚀

## Why Phoenix?

Phoenix rises from the ashes of legacy TUI frameworks, solving critical problems:

- ✅ **Perfect Unicode/Emoji support** - No more layout bugs
- ✅ **10x Performance** - Differential rendering, caching, zero allocations
- ✅ **DDD Architecture** - Clean, testable, extendable
- ✅ **Rich Component Library** - Everything you need out of the box
- ✅ **Public Cursor API** - Full control for shell applications

## Libraries

Phoenix is a modular framework with 8 independent libraries:

- **phoenix/core** ✅ - Terminal primitives, Unicode/Emoji support (CORRECT width calculation!)
- **phoenix/style** ✅ - CSS-like styling (bold, colors, borders, padding)
- **phoenix/layout** ✅ - Flexbox & grid layout (Box model, responsive sizing)
- **phoenix/tea** ✅ - Elm Architecture (Model-View-Update, type-safe)
- **phoenix/render** ✅ - High-performance differential renderer (29,000 FPS!)
- **phoenix/components** ✅ - Universal UI components (94.5% avg coverage):
  - TextArea (multiline editing, 96.4%) **NEW: Advanced cursor control API!**
  - TextInput (single-line, 90%)
  - List (94.7%) | Viewport (94.5%) | Table (92%) | Menu | Tabs
  - Modal (96.5%) | Progress (98.5%)
- **phoenix/mouse** ✅ - Mouse events (click, scroll, drag-drop, right-click support)
- **phoenix/clipboard** ✅ - Cross-platform clipboard (OSC 52 for SSH)

## Installation

### Install All Libraries (Recommended for new projects)

```bash
go get github.com/phoenix-tui/phoenix@latest
```

This installs the umbrella module with convenient access to all Phoenix libraries through a single import:

```go
import "github.com/phoenix-tui/phoenix"

// Use convenience API
term := phoenix.AutoDetectTerminal()
style := phoenix.NewStyle().Foreground("#00FF00").Bold()
p := phoenix.NewProgram(myModel, phoenix.WithAltScreen[MyModel]())
```

### Install Individual Libraries (For existing projects or selective use)

```bash
go get github.com/phoenix-tui/phoenix/tea@latest        # Elm Architecture
go get github.com/phoenix-tui/phoenix/components@latest # UI Components
go get github.com/phoenix-tui/phoenix/style@latest      # Styling
go get github.com/phoenix-tui/phoenix/core@latest       # Terminal primitives
```

Individual imports give you more control and smaller dependencies:

```go
import (
    tea "github.com/phoenix-tui/phoenix/tea/api"
    "github.com/phoenix-tui/phoenix/components/input/api"
)
```

## Quick Start

### Using the Umbrella Module

```bash
go get github.com/phoenix-tui/phoenix@latest
```

```go
package main

import (
    "fmt"
    "os"
    "github.com/phoenix-tui/phoenix"
    tea "github.com/phoenix-tui/phoenix/tea/api"
)

type Model struct {
    count int
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, phoenix.Quit()
        }
        m.count++
    }
    return m, nil
}

func (m Model) View() string {
    // Use Phoenix convenience API for styling
    style := phoenix.NewStyle().Foreground("#00FF00").Bold()
    return style.Render(fmt.Sprintf("Count: %d\n", m.count))
}

func main() {
    p := phoenix.NewProgram(Model{}, phoenix.WithAltScreen[Model]())
    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Using Individual Libraries

```bash
go get github.com/phoenix-tui/phoenix/tea@latest
```

```go
package main

import (
    "fmt"
    "os"
    "github.com/phoenix-tui/phoenix/tea/api"
)

type Model struct {
    count int
}

func (m Model) Init() api.Cmd { return nil }

func (m Model) Update(msg api.Msg) (Model, api.Cmd) {
    switch msg := msg.(type) {
    case api.KeyMsg:
        if msg.String() == "q" {
            return m, api.Quit()
        }
        m.count++
    }
    return m, nil
}

func (m Model) View() string {
    return fmt.Sprintf("Count: %d\nPress any key to increment, 'q' to quit\n", m.count)
}

func main() {
    p := api.New(Model{}, api.WithAltScreen[Model]())
    if err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Documentation

### Public Documents
- **[ROADMAP.md](ROADMAP.md)** - Public roadmap (milestones, progress, dates)
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and changes
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development guide

### Strategic
- [Master Plan](docs/dev/MASTER_PLAN.md) - Strategic vision & success metrics
- [Architecture](docs/dev/ARCHITECTURE.md) - DDD technical design
- [API Design](docs/dev/API_DESIGN.md) - API principles & examples
- [Technical Roadmap](docs/dev/ROADMAP.md) - Detailed 20-week plan

### Quality & Readiness
- **[v0.1.0 Readiness Report](docs/dev/FINAL_V0.1.0_READINESS_REPORT.md)** - 🔴 PRODUCTION READY assessment
- [Mouse Coverage Report](docs/dev/MOUSE_COVERAGE_REPORT.md) - Comprehensive test coverage analysis
- [GoSh Readiness](docs/dev/PHOENIX_GOSH_READINESS.md) - Migration readiness status

### Research
- [Charm Pain Points](docs/dev/research/CHARM_PAIN_POINTS.md) - Why we built Phoenix

## Development Status

| Library | Status | Coverage | Week | Notes |
|---------|--------|----------|------|-------|
| **core** | ✅ Complete | 98.4% | 3-4 | Unicode/Emoji CORRECT! |
| **style** | ✅ Complete | 90%+ | 5-6 | CSS-like styling |
| **tea** | ✅ Complete | 95.7% | 7-8 | Elm Architecture |
| **layout** | ✅ Complete | 97.9% | 9-10 | Flexbox + Box model |
| **components** | ✅ Complete | 94.5% | 11-12 | 6 universal components |
| **render** | ✅ Complete | **91.7%** | 13-14 | 29,000 FPS (489x faster!) |
| **mouse** | ✅ Complete | **99.7%** | 16 | 3 critical bugs fixed! |
| **clipboard** | ✅ Complete | 82.0% | 16 | Cross-platform + SSH |

### Overall Progress

```
Phase 1: Foundation     ████░░░░░░░░░░░░░░░░  (10%)  ✅ Weeks 1-2
Phase 2: Core Libs      ░░░░████████░░░░░░░░  (30%)  ✅ Weeks 3-8
Phase 3: Components     ░░░░░░░░░░░░████░░░░  (20%)  ✅ Weeks 9-12
Phase 4: Advanced       ░░░░░░░░░░░░░░░░████  (15%)  ✅ Weeks 13-16
Phase 5: Launch         ░░░░░░░░░░░░░░░░░░█░  (5%)   ✅ Coverage Sprint
                        ══════════════════════
                        Progress: 80% (Ready for v0.1.0!)
```

**Test Coverage Sprint Complete**: 36,000 lines of test code, 4,340+ test cases, 3 critical bugs fixed

### Completed Features

**Week 3-4: phoenix/core**
- Terminal primitives (ANSI, raw mode, capabilities)
- Unicode/Emoji width calculation (CORRECT - fixes Charm bug!)
- Grapheme cluster support (👋🏽 = 1 cluster, 2 cells)
- 98.4% test coverage

**Week 5-6: phoenix/style**
- CSS-like styling (bold, italic, colors)
- Border/padding/margin support
- 8-stage rendering pipeline
- Fluent builder API
- 90%+ test coverage

**Week 7-8: phoenix/tea**
- Elm Architecture (Model-Update-View)
- Type-safe event loop
- Command system (Quit, Batch, Sequence, Tick)
- Generic constraints (no interface{} casts!)
- 95.7% test coverage

**Week 9-10: phoenix/layout**
- Box model (padding, margin, border, sizing)
- Flexbox layout (row/column, gap, flex grow/shrink)
- Responsive sizing
- 97.9% test coverage (highest!)

**Week 11-12: phoenix/components**
- **TextInput** (90.0%) - Public cursor API, grapheme-aware, selection, validation
- **List** (94.7%) - Single/multi selection, filtering, custom rendering
- **Viewport** (94.5%) - Scrolling, follow mode, large content (10K+ lines)
- **Table** (92.0%) - Sortable columns, custom cell renderers, navigation
- **Modal** (96.5%) - Focus trap, buttons, dimming, keyboard shortcuts
- **Progress** (98.5%) - Bar + 15 spinner styles
- Average coverage: **94.5%** (exceeds 90% target!)

**Week 13-14: phoenix/render**
- Differential rendering (virtual buffer)
- 29,000 FPS achieved (489x faster than 60 FPS target!)
- Zero allocations in hot paths
- **91.7% test coverage** (improved from 87.1%)

**Week 16: phoenix/mouse** 🔥
- All buttons (Left, Right, Middle, WheelUp, WheelDown)
- Click detection (single/double/triple - automatic!)
- Drag & drop state tracking
- Multi-protocol (SGR, X10, URxvt)
- Comprehensive README (588 lines)
- **99.7% test coverage** - 6,000+ lines test code
- **3 critical bugs found and fixed** during coverage sprint!

**Week 16: phoenix/clipboard**
- Cross-platform (Windows/macOS/Linux)
- OSC 52 for SSH sessions (auto-detect)
- Native APIs (user32.dll, pbcopy/pbpaste, xclip/xsel)
- DDD architecture
- 82% average test coverage (domain 100%)

## Key Features

### 1. Perfect Unicode/Emoji Support ✅
**Problem**: Charm's Lipgloss has broken emoji width calculation ([issue #562](https://github.com/charmbracelet/lipgloss/issues/562))
**Solution**: Phoenix uses grapheme cluster detection with correct East Asian Width (UAX #11)

```go
// Phoenix: CORRECT
text := "Hello 👋 World 🌍"
width := style.Width(text)  // Returns 17 (correct!)

// Charm Lipgloss: BROKEN
width := lipgloss.Width(text)  // Returns 19 (wrong!)
```

### 2. 10x Performance ✅
**Benchmark**: 29,000 FPS (489x faster than 60 FPS target)
**Techniques**: Differential rendering, caching, zero allocations

### 3. DDD Architecture ✅
```
library/
├── domain/        # Business logic (95%+ coverage)
├── application/   # Use cases
├── infrastructure/ # Technical details
└── api/           # Public interface
```

### 4. Public Cursor API ✅
**Problem**: Bubbles TextArea has private cursor - syntax highlighting impossible
**Solution**: Phoenix TextInput exposes `CursorPosition()` and `ContentParts()`

```go
// Phoenix: PUBLIC API (syntax highlighting works!)
before, at, after := input.ContentParts()
highlighted := syntax.Highlight(before) +
               cursor.Render(at) +
               syntax.Highlight(after)

// Bubbles: PRIVATE (syntax highlighting impossible!)
// cursor is internal field - no access
```

### 5. Mouse & Clipboard Support ✅
**Mouse**: All buttons (Left, Right, Middle, Wheel), drag-drop, click detection
**Clipboard**: Cross-platform (Windows/macOS/Linux), SSH support (OSC 52)

### 6. Progress Component ✅
**Available**: Progress Bar + 15 Spinner Styles (Week 11-12, 98.5% coverage)
**Location**: `github.com/phoenix-tui/phoenix/components/progress/api`

Phoenix includes a comprehensive Progress component with both bars and animated spinners:

```go
import progress "github.com/phoenix-tui/phoenix/components/progress/api"

// Progress Bar
bar := progress.NewBar(100).  // Max value 100
    SetWidth(40).
    SetLabel("Downloading").
    SetValue(65)  // Current progress 65%

// Animated Spinner (15 styles available!)
spinner := progress.NewSpinner(progress.SpinnerDots).
    SetLabel("Loading").
    SetFPS(10)

// Example styles: SpinnerDots, SpinnerLine, SpinnerArrow, SpinnerCircle,
// SpinnerBounce, SpinnerPulse, SpinnerGrowHorizontal, SpinnerGrowVertical, etc.
```

**Features**:
- Progress bars with customizable width and characters
- 15 pre-built spinner styles (dots, lines, arrows, circles, bouncing, etc.)
- Label support for both bars and spinners
- Configurable FPS for smooth animations
- 98.5% test coverage

**Examples**: See [examples/progress/](examples/progress/) for working demonstrations:
- `bar_simple.go` - Basic progress bar
- `bar_styled.go` - Styled progress bar with colors
- `spinner_simple.go` - Animated spinner
- `multi_progress.go` - Multiple progress indicators

**Documentation**: See [components/progress/README.md](components/progress/README.md) for full API reference

## Next Steps

### Week 17-18: GoSh Migration (Real-World Validation)
- Migrate GoSh shell to Phoenix components
- Validate all libraries in production application
- Performance benchmarks vs Charm-based GoSh
- Expected: 90x faster rendering + new features!

### Week 19-20: Polish & v0.1.0 Launch
- Final performance optimizations
- Documentation polish
- Public release preparation
- Community engagement

## Contributing

Phoenix is part of an active development effort. See [docs/dev/](docs/dev/) for:
- Architecture documentation
- Development roadmap
- Contribution guidelines

## License

MIT License - see [LICENSE](LICENSE) file for details

## Special Thanks

**Professor Ancha Baranova** - This project would not have been possible without her invaluable help and support. Her assistance was crucial in bringing Phoenix to life.

---

*Rising from the ashes of legacy TUI frameworks* 🔥
**PRODUCTION READY - Test Coverage Sprint Complete!** ⭐
*93.5% average coverage | 36,000 lines test code | 3 critical bugs fixed*
