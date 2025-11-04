# Cobra + Phoenix Integration Example

This example demonstrates the **hybrid CLI+TUI pattern** - the recommended approach for production CLI tools.

## 🎯 Pattern: CLI + TUI Hybrid

### Why This Pattern?

Modern CLI tools should support **two modes**:

1. **CLI Mode** (flags) → For scripts, CI/CD, automation
2. **TUI Mode** (interactive) → For humans, onboarding, exploration

**Same tool, different UX depending on usage!**

### Real-World Examples

This pattern is used by popular tools:
- `kubectl` - flags for scripts, interactive prompts for setup
- `gh` (GitHub CLI) - flags for automation, interactive for PR creation
- `terraform` - flags for CI, interactive plan review

## 🚀 Usage

### CLI Mode (Scriptable)

```bash
# Use flags for automation
./cobra-cli --name "John Doe" --email "john@example.com" --message "Hello Phoenix!"

# Perfect for scripts
for user in $(cat users.txt); do
  ./cobra-cli --name "$user" --email "$user@example.com" --message "Welcome!"
done

# CI/CD friendly
./cobra-cli -n "Bot" -e "bot@ci.com" -m "Deploy successful" && deploy.sh
```

### TUI Mode (Interactive)

```bash
# No flags → launches beautiful TUI
./cobra-cli

# Phoenix TUI appears with:
# - Tab navigation between fields
# - Real-time validation
# - Beautiful styling with emojis 🎨
# - Guided UX (users can't make mistakes!)
```

## 📦 Installation

```bash
# Clone Phoenix repository
git clone https://github.com/phoenix-tui/phoenix
cd phoenix/examples/cobra-cli

# Get dependencies
go mod download

# Run
go run main.go
```

## 🏗️ Architecture

### Component Structure

```
main.go
├── rootCmd (Cobra)           # CLI framework
│   ├── Flags                 # --name, --email, --message
│   └── Run()
│       ├── CLI Mode          # if flags provided
│       └── TUI Mode          # if no flags
│
├── formModel (Phoenix)       # TUI implementation
│   ├── TextInput components  # Phoenix components
│   ├── Update() logic        # Event handling
│   └── View() rendering      # Beautiful UI
│
└── Result display            # After submission
```

### Key Components

1. **Cobra** (`github.com/spf13/cobra`)
   - Command structure
   - Flag parsing
   - Help generation
   - Subcommands support

2. **Phoenix** (`github.com/phoenix-tui/phoenix`)
   - `phoenix/components` - TextInput, List, etc.
   - `phoenix/style` - CSS-like styling (with correct Unicode!)
   - `phoenix/tea` - Elm Architecture (event loop)
   - `phoenix/layout` - Flexbox layouts

### Decision Logic

```go
func Run(cmd *cobra.Command, args []string) {
    // Simple check: are flags provided?
    if cmd.Flags().NFlag() > 0 {
        runCLIMode()  // Process flags
    } else {
        runTUIMode()  // Launch interactive Phoenix
    }
}
```

## 🎨 Phoenix Advantages Over Charm

This example could be built with Charm (Huh + Lipgloss), but Phoenix offers:

| Feature | Charm (Huh/Lipgloss) | Phoenix |
|---------|---------------------|---------|
| **Unicode/Emoji** | ❌ Broken ([#562](https://github.com/charmbracelet/lipgloss/issues/562)) | ✅ Perfect |
| **Performance** | ~60 FPS | 29,000 FPS |
| **Modularity** | Monolithic | DDD layers |
| **Customization** | Limited | Full control |
| **Test Coverage** | Unknown | 91.8% |

### Unicode Bug Example

```go
// Charm/Lipgloss - BROKEN:
prompt := "Name 👋: "  // Width counted wrong → layout breaks

// Phoenix - WORKS:
prompt := "Name 👋: "  // Perfect width calculation ✓
```

## 📚 Learn More

### Phoenix Documentation
- [Getting Started](https://github.com/phoenix-tui/phoenix/blob/main/docs/user/tutorials/01-getting-started.md)
- [Building Components](https://github.com/phoenix-tui/phoenix/blob/main/docs/user/tutorials/02-building-components.md)
- [API Reference](https://github.com/phoenix-tui/phoenix/tree/main/docs/api)

### Migration from Charm
- [Migration Guide](https://github.com/phoenix-tui/phoenix/blob/main/docs/user/MIGRATION_GUIDE.md)
- [From Bubbletea](https://github.com/phoenix-tui/phoenix/blob/main/docs/user/MIGRATION_FROM_BUBBLETEA.md)

### Cobra Documentation
- [Cobra GitHub](https://github.com/spf13/cobra)
- [User Guide](https://github.com/spf13/cobra/blob/main/user_guide.md)

## 🔧 Extending This Example

### Add More Components

```go
// Select dropdown (coming in v0.2.0)
selectInput := components.NewSelect(
    components.WithOptions([]string{"Option 1", "Option 2"}),
)

// Confirm dialog (coming in v0.2.0)
confirm := components.NewConfirm(
    components.WithMessage("Are you sure?"),
)

// Progress indicator (already available!)
progress := components.NewProgress(
    components.WithTotal(100),
)
```

### Add Subcommands

```go
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new resource",
    Run: func(cmd *cobra.Command, args []string) {
        if cmd.Flags().NFlag() > 0 {
            // CLI: create --name "foo"
        } else {
            // TUI: interactive creation wizard
        }
    },
}

rootCmd.AddCommand(createCmd)
```

### Add Validation

```go
func (m formModel) validateEmail(email string) bool {
    // Simple email validation
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (m formModel) View() string {
    // Show validation errors
    if !m.validateEmail(m.emailInput.Value()) {
        return "⚠ Invalid email format"
    }
    // ...
}
```

## 🎓 Best Practices

### 1. Always Support Both Modes

```bash
# ✅ GOOD - supports both
./tool --name "John"  # CLI mode
./tool                # TUI mode

# ❌ BAD - forces one mode
./tool               # Always TUI (bad for scripts!)
./tool --interactive # Requires flag for TUI (annoying!)
```

### 2. Validate in Both Modes

```go
// CLI mode validation
if name == "" {
    fmt.Fprintln(os.Stderr, "Error: --name is required")
    os.Exit(1)
}

// TUI mode validation
func (m formModel) isFormValid() bool {
    return m.nameInput.Value() != ""
}
```

### 3. Consistent Output

```go
// Both modes should produce same output format
// Good for piping: ./tool | jq .name
```

### 4. Respect Terminal Capabilities

```go
// Detect if terminal supports TUI
if !term.IsTerminal(int(os.Stdin.Fd())) {
    // Not a TTY → force CLI mode
    runCLIMode()
}
```

## 🚀 Production Deployment

### Build for Multiple Platforms

```bash
# Use Goreleaser (like in Habr article)
goreleaser release --snapshot --clean

# Or manual cross-compilation
GOOS=linux GOARCH=amd64 go build -o mytool-linux
GOOS=darwin GOARCH=arm64 go build -o mytool-macos
GOOS=windows GOARCH=amd64 go build -o mytool.exe
```

### Distribution

```bash
# Homebrew
brew tap yourorg/tap
brew install mytool

# Go install
go install github.com/yourorg/mytool@latest

# Docker
docker run yourorg/mytool --name "John"
```

## 💡 Tips & Tricks

### Auto-detect Mode

```go
// If piped → CLI mode automatically
if !term.IsTerminal(int(os.Stdout.Fd())) {
    runCLIMode()
}

// If environment variable set
if os.Getenv("CI") != "" {
    runCLIMode()  // CI/CD environment
}
```

### Help Integration

```bash
# Cobra generates help automatically
./cobra-cli --help
./cobra-cli create --help

# TUI can show help too (F1 key)
case "f1":
    return m, showHelpModal()
```

### Configuration Files

```go
// Support config files too!
// ~/.mytool.yaml
name: "John Doe"
email: "john@example.com"

// Precedence: CLI flags > ENV vars > Config file > TUI input
```

## 🐛 Troubleshooting

### TUI Not Displaying

```bash
# Check if terminal supports colors
echo $TERM  # Should be xterm-256color or similar

# Test terminal capabilities
./cobra-cli  # Should show TUI
```

### Unicode Issues

Phoenix handles Unicode correctly, but ensure:
- Terminal font supports emojis (e.g., Nerd Fonts)
- Terminal encoding is UTF-8

```bash
# Check encoding
locale  # Should show UTF-8
```

## 📄 License

This example is part of Phoenix TUI Framework.
See [LICENSE](../../LICENSE) for details.

---

**Questions?** Open an issue or discussion on [GitHub](https://github.com/phoenix-tui/phoenix)!

**Show us what you build!** We'd love to see your Cobra + Phoenix applications! 🚀
