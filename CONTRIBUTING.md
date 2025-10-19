# Contributing to Phoenix TUI Framework

Thank you for your interest in contributing to Phoenix!

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.25+** - [Install Go](https://go.dev/doc/install)
- **Task** - [Install Task](https://taskfile.dev/installation/) (recommended)
- **golangci-lint** - [Install golangci-lint](https://golangci-lint.run/welcome/install/)

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/phoenix-tui/phoenix.git
cd phoenix

# Install Task (choose your platform)
# macOS/Linux:
brew install go-task

# Windows (Scoop):
scoop install task

# Windows (Chocolatey):
choco install go-task

# Or download from https://taskfile.dev/installation/

# Verify installation
task --version

# Show available tasks
task
```

---

## 🛠️ Development Workflow

### Running Tests

```bash
# Run all tests with coverage
task test

# Run tests for specific package
task test:core

# Generate HTML coverage report
task test:coverage
# Opens coverage.html in browser

# Watch mode (requires entr or fswatch)
task test:watch
```

### Code Quality

```bash
# Run linter
task lint

# Run linter and auto-fix issues
task lint:fix

# Format code
task fmt

# Run all quality checks
task check  # fmt + vet + lint + test

# Pre-commit checks (recommended before git commit)
task dev    # fmt + vet + lint:fix + test
```

### Building

```bash
# Build all packages
task build

# Build example applications
task build:examples
# Creates bin/basic.exe, bin/unicode.exe

# Run examples
task run:basic
task run:unicode
```

### Benchmarks

```bash
# Run all benchmarks
task bench

# Run core benchmarks only
task bench:core

# Run Unicode benchmarks (Week 4)
task bench:unicode

# Compare benchmarks
task bench:compare
```

### Dependencies

```bash
# Download and tidy dependencies
task deps

# Update all dependencies
task deps:update

# Verify dependencies
task deps:verify
```

### Cleaning

```bash
# Remove build artifacts and coverage files
task clean
```

---

## 📋 Before Submitting PR

Run the development check:

```bash
task dev
```

This will:
1. ✅ Format code (`gofmt`)
2. ✅ Run go vet
3. ✅ Auto-fix linter issues
4. ✅ Run all tests

If all checks pass, you're ready to commit!

---

## 🌿 Git-Flow Branching

Phoenix uses **Git-Flow** with `main` and `develop` branches:

```bash
# Main branches
main              # Production releases (v0.1.0-beta.1, v0.1.0, etc.)
develop           # Active development (default branch for PRs)

# Supporting branches
feature/*         # New features (branch from develop)
release/*         # Release preparation (branch from develop)
hotfix/*          # Critical fixes (branch from main)
```

### Creating a Feature Branch

```bash
# Start from develop
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/my-new-feature

# Work on your changes
# ... make commits ...

# Push to your fork
git push origin feature/my-new-feature

# Create Pull Request to develop branch
```

See **[WORKFLOW.md](WORKFLOW.md)** for complete git-flow documentation.

---

## 🔄 Pull Request Process

### 1. Before Creating PR

- ✅ Run `task dev` (all checks must pass)
- ✅ Update tests (coverage must not decrease)
- ✅ Update documentation if needed
- ✅ Follow commit message format (Conventional Commits)
- ✅ Rebase on latest `develop`

### 2. Creating PR

1. Push your feature branch to **your fork**
2. Open PR from your fork to `phoenix-tui/phoenix:develop`
3. Fill out PR template (if provided)
4. Add clear description of changes
5. Link related issues (if any)

### 3. PR Requirements

Your PR must:
- ✅ Pass all CI checks (tests, lint, format)
- ✅ Maintain or improve test coverage (90%+ minimum)
- ✅ Have at least 1 approval from maintainer
- ✅ No merge conflicts with `develop`
- ✅ Follow project code style

### 4. Code Review

- Maintainers will review within 1-3 business days
- Address feedback by pushing new commits
- Once approved, maintainer will merge (usually squash merge)

### 5. After Merge

- Your changes appear in next release
- Delete your feature branch
- Pull latest `develop`

---

## 🚦 CI/CD Requirements

All PRs must pass these automated checks:

### 1. Tests
```bash
go test -v -race -cover ./...
```
- All tests must pass
- No race conditions
- Coverage must not decrease below current (93.5%)

### 2. Linter
```bash
golangci-lint run --config .golangci.yml ./...
```
- Zero linter issues (enforced)
- See `.golangci.yml` for enabled linters

### 3. Format
```bash
gofmt -l .
```
- All code must be formatted with `gofmt`
- Zero unformatted files

### 4. Go Vet
```bash
go vet ./...
```
- Must pass with zero suspicious constructs

### 5. Build
```bash
go build ./...
```
- All packages must compile successfully
- Works on Linux, macOS, Windows

**CI Pipeline**: Automated via GitHub Actions on every push/PR

---

## 🔧 Alternative: Without Task

If you prefer not to use Task, here are the raw commands:

```bash
# Run tests
go test -v -race -cover ./...

# Run linter
golangci-lint run --config .golangci.yml ./...

# Format code
go fmt ./...

# Run vet
go vet ./...

# Build
go build ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## 📦 Project Structure

```
phoenix/
├── core/              # Foundation library (terminal, Unicode, capabilities)
├── style/             # Styling library (colors, borders, padding)
├── tea/               # Event loop (Elm Architecture)
├── layout/            # Layout system (flexbox, grid)
├── render/            # High-performance renderer
├── components/        # UI components (input, list, table, etc.)
├── mouse/             # Mouse events
├── clipboard/         # Clipboard operations
├── examples/          # Example applications
├── docs/              # Documentation
├── benchmarks/        # Performance benchmarks
├── Taskfile.yml       # Task automation
└── .golangci.yml      # Linter configuration
```

---

## 🧪 Testing Standards

Phoenix has **strict testing requirements**:

- **Domain Layer**: 95%+ coverage (pure business logic)
- **Application Layer**: 90%+ coverage (use cases)
- **Infrastructure Layer**: 80%+ coverage (integration tests)
- **API Layer**: 85%+ coverage (example-based tests)
- **Overall Project**: 90%+ minimum

**Current coverage: 93.5% average** ✅

Coverage by library:
- core: 98.4% | style: 90%+ | tea: 95.7% | layout: 97.9%
- components: 94.5% | render: 91.7% | mouse: 99.7% | clipboard: 82.0%

**Quality**: 36,000 lines of test code, 4,340+ test cases, 3 critical bugs found and fixed

---

## 📝 Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, not styling library)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks

**Examples:**
```bash
feat(core): add Unicode width calculation
fix(style): resolve emoji rendering bug
docs(api): update terminal API examples
test(core): add property-based tests for Cell
chore: add golangci-lint configuration
```

---

## 🎯 Code Style

- **Language**: All code comments and documentation in English
- **Formatting**: Use `gofmt` (enforced by CI)
- **Linting**: Pass `golangci-lint` checks (see `.golangci.yml`)
- **Naming**: Follow Go conventions (PascalCase for exported, camelCase for unexported)
- **Comments**:
  - Package comments required (revive)
  - Exported functions must have comments
  - Comments should end with period (godot)

---

## 🚫 What NOT to Commit

- `*.exe` - Build artifacts
- `coverage.out`, `coverage.html` - Coverage reports
- `.claude/settings.local.json` - Personal AI settings
- `bin/` - Build output directory
- `nul` - Windows temp files

See `.gitignore` for full list.

---

## 🤝 Getting Help

- **Documentation**: See [docs/dev/](docs/dev/) for architecture and design docs
- **Issues**: [GitHub Issues](https://github.com/phoenix-tui/phoenix/issues) (when repo is public)
- **Discussions**: [GitHub Discussions](https://github.com/phoenix-tui/phoenix/discussions) (when repo is public)

---

## 📖 Additional Resources

### Strategic Documents
- [MASTER_PLAN.md](docs/dev/MASTER_PLAN.md) - Strategic vision
- [ARCHITECTURE.md](docs/dev/ARCHITECTURE.md) - Technical architecture
- [API_DESIGN.md](docs/dev/API_DESIGN.md) - API principles and examples
- [ROADMAP.md](docs/dev/ROADMAP.md) - 20-week development timeline

### Quality & Status
- [FINAL_V0.1.0_READINESS_REPORT.md](docs/dev/FINAL_V0.1.0_READINESS_REPORT.md) - Production readiness (LATEST)
- [MOUSE_COVERAGE_REPORT.md](docs/dev/MOUSE_COVERAGE_REPORT.md) - Test coverage analysis

### Navigation
- [INDEX.md](docs/dev/INDEX.md) - Complete documentation index (Kanban structure)

---

*Last updated: 2025-10-19 | Status: PRODUCTION READY*
