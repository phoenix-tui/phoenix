# Changelog

All notable changes to Phoenix TUI Framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned for v0.2.0
- Additional TUI components (Spinner, ProgressBar, Form)
- Theme system with presets
- Animation framework
- Advanced layout features (Grid)

---

## [0.1.0-beta.1] - 2025-10-19 (First Public Beta)

**Status**: 🎉 FIRST PUBLIC RELEASE

This is the first public beta release of Phoenix TUI Framework. The framework is ready for community testing and feedback. While labeled as beta, the code is production-ready with 93.5-96.4% test coverage.

### What's Included

All 8 core libraries are complete and tested:

- ✅ **phoenix/core** - Terminal primitives with perfect Unicode/Emoji support
- ✅ **phoenix/style** - CSS-like styling system
- ✅ **phoenix/tea** - Elm Architecture (Model-Update-View) with 95.7% coverage
- ✅ **phoenix/layout** - Flexbox layout system with 97.9% coverage
- ✅ **phoenix/render** - High-performance differential renderer (29,000 FPS!)
- ✅ **phoenix/components** - 6 production-ready components (TextInput, List, Viewport, Table, Modal, Progress)
- ✅ **phoenix/mouse** - Complete mouse event support
- ✅ **phoenix/clipboard** - Cross-platform clipboard (OSC 52 for SSH)

### Documentation

- ✅ Migration guide from Bubbletea/Lipgloss/Bubbles
- ✅ Testing guide with helpers (NullTerminal, MockTerminal)
- ✅ Troubleshooting guide for common issues
- ✅ Comprehensive API documentation

### Dependencies

- Migrated to published `github.com/unilibs/uniwidth@v0.1.0-beta` (3.9-46x faster than alternatives)

### Known Limitations

- API may change based on community feedback (semver allows breaking changes in beta)
- Some advanced components planned for v0.2.0
- CI/CD workflows to be added post-beta

### Community

- GitHub: https://github.com/phoenix-tui/phoenix
- Organization: phoenix-tui
- License: MIT

---

## [0.1.0] - 2025-10-19 (Internal Production Ready)

**Status**: ✅ PRODUCTION READY (93.5% test coverage, all libraries complete)

This is the first production-ready release of Phoenix TUI Framework. All 8 core libraries are complete, tested, and ready for real-world use. The framework solves all 7 critical pain points from the Charm ecosystem.

### Added

#### Core Libraries (Weeks 3-8)

**phoenix/core** (Week 3-4)
- Terminal primitives (ANSI escape sequences, cursor positioning)
- **Unicode/Emoji width calculation** (THE fix for Lipgloss #562)
- Grapheme cluster support (correct handling of 👋🏽 as 1 cluster, 2 cells)
- Terminal capabilities detection
- Position and Size value objects
- Test coverage: 98.4%

**phoenix/style** (Week 5-6)
- CSS-like styling system (bold, italic, underline, strikethrough)
- Color support (foreground, background, RGB, ANSI 256)
- Border rendering (single, double, rounded, thick styles)
- Padding and margin support
- Fluent builder API
- 8-stage rendering pipeline
- Test coverage: 90%+

**phoenix/tea** (Week 7-8)
- Elm Architecture implementation (Model-Update-View)
- Type-safe event loop with Go 1.25+ generics
- Message system (KeyMsg, MouseMsg, WindowSizeMsg, custom messages)
- Command system (Quit, Batch, Sequence, Tick)
- Program lifecycle management (Init, Update, View)
- Bubbletea-compatible API
- Test coverage: 95.7%

#### Layout & Components (Weeks 9-12)

**phoenix/layout** (Week 9-10)
- Box model (padding, margin, border, sizing)
- Flexbox layout system (row/column direction, gap, alignment)
- Flex grow/shrink support
- Responsive sizing
- Test coverage: 97.9% (highest!)

**phoenix/components** (Week 11-12)
- **TextInput** component (90.0% coverage)
  - Single-line text input
  - **Public cursor API** (syntax highlighting possible!)
  - Grapheme-aware cursor movement
  - Horizontal scrolling for long input
  - Selection support
  - Validation hooks
- **List** component (94.7% coverage)
  - Single/multi selection modes
  - Keyboard navigation (j/k Vim-style)
  - Custom item rendering
  - Filtering (built-in + custom)
- **Viewport** component (94.5% coverage)
  - Scrollable content (tested with 10K+ lines)
  - Keyboard scrolling (PgUp/PgDn, Ctrl-U/Ctrl-D)
  - Mouse wheel support
  - Auto-scrolling (follow mode for logs)
- **Table** component (92.0% coverage)
  - Column definitions (width, alignment, sortable)
  - Sorting (ascending/descending, toggle)
  - Custom cell renderers
  - Keyboard navigation
- **Modal** component (96.5% coverage)
  - Overlay rendering (centered)
  - Focus trap (modal captures all input)
  - Button support with keyboard shortcuts
  - Background dimming
- **Progress** component (98.5% coverage)
  - Progress bar with customizable width/character
  - 15 spinner styles (dots, line, arrow, circle, bounce, etc.)
  - Label support
  - Configurable FPS

Average component coverage: **94.5%**

#### High-Performance Rendering (Weeks 13-14)

**phoenix/render** (Week 13-14)
- Differential rendering engine (virtual buffer)
- **29,000 FPS achieved** (489x faster than 60 FPS target!)
- Only renders changed cells (70% I/O reduction)
- Zero allocations in hot paths
- Cell-based abstraction
- Test coverage: 91.7%

#### Advanced Input (Week 16)

**phoenix/mouse** (Week 16)
- **All mouse buttons**: Left, Right, Middle, WheelUp, WheelDown
- **Click detection**: Single, double, triple (automatic!)
- **Drag & drop**: Built-in state tracking with threshold
- **Scroll wheel**: Viewport scrolling support
- **Multi-protocol**: SGR (1006), X10 (1000), URxvt (1015)
- **Motion events**: Mouse movement tracking
- **Modifiers**: Shift, Ctrl, Alt support
- Comprehensive README (588 lines)
- Test coverage: **99.7%** (highest!)

**phoenix/clipboard** (Week 16)
- **Cross-platform**: Windows (user32.dll), macOS (pbcopy/pbpaste), Linux (xclip/xsel)
- **SSH support**: OSC 52 protocol (auto-detects $SSH_TTY)
- **Primary selection**: X11 middle-click paste
- Read and write operations
- DDD architecture with clear layers
- Test coverage: 82% average (domain 100%)

### Fixed

#### Test Coverage Sprint (Post-Week 16)

**Critical Bugs Found and Fixed**:

1. **Parser bitmask error** (CRITICAL)
   - **Affected**: sgr_parser.go, x10_parser.go, urxvt_parser.go
   - **Impact**: Motion events completely broken
   - **Issue**: Bitmask `0x43` missing bit 5 for motion events (codes 32, 35)
   - **Fix**: Changed to `0x63` in all three parsers
   - **Result**: Motion tracking now works correctly

2. **X10 FormatSequence UTF-8 encoding bug** (HIGH)
   - **Affected**: x10_parser.go
   - **Impact**: Large coordinates (>95) created invalid sequences
   - **Issue**: `fmt.Sprintf` with `%c` converts bytes >127 to multi-byte UTF-8
   - **Fix**: Use raw byte array construction instead of format string
   - **Result**: All coordinates work correctly (tested up to 9999)

3. **SGR IsMotion always returned false** (MEDIUM)
   - **Affected**: sgr_parser.go
   - **Impact**: Motion detection completely broken
   - **Issue**: Same incorrect bitmask (`0x43`) in IsMotion() method
   - **Fix**: Updated to correct `0x63` mask
   - **Result**: Motion detection now functional

**Coverage Improvements**:
- **phoenix/mouse**: 60% → 99.7% (+39.7%)
  - 6,000+ lines of test code
  - 1,000+ test cases
  - All protocols tested (SGR, X10, URxvt)
  - All event types tested (press, release, click, drag, motion, scroll)
- **phoenix/clipboard**: 60-97% → 82% average (domain 100%)
  - 21 new test functions
  - Cross-platform scenarios covered
- **phoenix/render**: 87.1% → 91.7% (+4.6%)
  - 17 new comprehensive tests
  - Application layer improved (64.5% → 79.4%)

**Overall**: 93.5% average test coverage (36,000 lines test code, 4,340+ test cases)

### Changed

- **Architecture**: Consistent DDD (Domain-Driven Design) across all libraries
  - Domain layer: Pure business logic (95%+ coverage target)
  - Application layer: Use cases (90%+ coverage target)
  - Infrastructure layer: Technical details (80%+ coverage target)
  - API layer: Public interface (85%+ coverage target)

- **Testing Standards**: Raised minimum coverage from 80% to 90% project-wide
  - Achieved: 93.5% average (exceeds target)
  - Domain layers: 95%+ coverage consistently
  - Comprehensive test patterns: table-driven, property-based, round-trip

- **Performance**: Optimized for zero allocations in hot paths
  - Render loop: <0.04ms per frame (29,000 FPS)
  - Unicode width calculation: Cached results
  - ANSI sequence generation: Pre-allocated buffers

### Documentation

- **Strategic Documents**
  - MASTER_PLAN.md - Strategic vision and success metrics
  - ARCHITECTURE.md - Complete technical architecture (22,000 words)
  - API_DESIGN.md - API principles and examples
  - ROADMAP.md (technical) - Detailed 20-week timeline
  - ROADMAP.md (public) - High-level public roadmap

- **Quality Reports**
  - FINAL_V0.1.0_READINESS_REPORT.md - Production readiness assessment
  - MOUSE_COVERAGE_REPORT.md - Test coverage sprint analysis
  - PHOENIX_GOSH_READINESS.md - Migration readiness for GoSh

- **Research**
  - CHARM_PAIN_POINTS.md - Problems with Charm ecosystem ($72K cost analysis)
  - TUI_ECOSYSTEM_RESEARCH_REPORT.md - TUI frameworks analysis
  - SHELL_COMPONENTS_DESIGN.md - Shell-specific component design

- **Development**
  - CONTRIBUTING.md - Development guide (setup, tasks, workflow)
  - INDEX.md (root) - Quick navigation
  - docs/dev/INDEX.md - Complete documentation index (Kanban structure)

- **Library-Specific**
  - mouse/README.md - Comprehensive mouse library guide (588 lines)
  - Each library: Package documentation with examples

### Performance Benchmarks

| Metric | Target | Achieved | Improvement |
|--------|--------|----------|-------------|
| **Render Performance** | 60 FPS (16ms) | 29,000 FPS (0.034ms) | **489x faster** |
| **Unicode Width Calc** | <1ms | <0.1ms (cached) | **10x faster** |
| **Test Execution** | <2 min | <30 sec | **4x faster** |
| **Memory Allocations** | Minimal | Zero (hot paths) | **100% reduction** |

### Comparison with Charm Ecosystem

| Feature | Charm | Phoenix | Status |
|---------|-------|---------|--------|
| **Unicode/Emoji** | ❌ Broken | ✅ Correct | **Fixed** |
| **Performance** | ~60 FPS | 29,000 FPS | **489x faster** |
| **Cursor API** | ❌ Private | ✅ Public | **Enabled** |
| **Click Detection** | ⚠️ Manual | ✅ Automatic | **Improved** |
| **Drag & Drop** | ⚠️ Manual | ✅ Built-in | **Added** |
| **Clipboard** | ❌ None | ✅ Cross-platform | **Added** |
| **Test Coverage** | Unknown | 93.5% | **Transparent** |
| **Architecture** | Flat | DDD + Layers | **Modernized** |

**All 7 CHARM pain points solved** ✅

---

## [0.0.1] - 2025-10-13 (Initial Project Setup)

### Added
- Project structure (8 Go workspace libraries)
- Go 1.25+ configuration
- Task automation (Taskfile.yml)
- CI/CD foundation
- Documentation framework
- Git repository initialization

---

## Version Strategy

Phoenix follows semantic versioning with a cautious approach:

- **v0.1.0** (current) - First production-ready release
  - Collect community feedback
  - API can change based on real-world usage
  - Breaking changes acceptable with migration guides

- **v0.2.0 - v0.9.0** - Iterative improvements
  - Theme system (based on feedback)
  - Additional components
  - API refinements from real usage
  - Community-requested features

- **v1.0.0** - API stability (6-12 months after v0.1.0)
  - API frozen (semantic versioning enforced)
  - Backwards compatibility guaranteed
  - Production certification
  - Full migration tooling

**Philosophy**: We follow gosh's cautious approach - still on v0.1.0-beta.7 after extensive use. We won't rush to v1.0 until API is proven stable in production.

---

## Links

- **GitHub**: https://github.com/phoenix-tui/phoenix (coming soon)
- **Documentation**: [docs/dev/INDEX.md](docs/dev/INDEX.md)
- **Issues**: https://github.com/phoenix-tui/phoenix/issues (coming soon)
- **Discussions**: https://github.com/phoenix-tui/phoenix/discussions (coming soon)

---

**Phoenix TUI Framework** 🔥 - Rising from the ashes of legacy TUI frameworks

*The future of Terminal UI development in Go* 🚀
