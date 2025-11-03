# Phoenix TUI Framework - Tutorial Series

Welcome to the Phoenix TUI Framework tutorial series! These tutorials will take you from complete beginner to advanced Phoenix developer.

---

## Overview

Phoenix TUI is a next-generation Terminal User Interface framework for Go, featuring:

- **Perfect Unicode support** (emoji, CJK, grapheme clusters)
- **10x performance** improvement over existing frameworks
- **DDD architecture** with rich domain models
- **Type-safe** with Go 1.25+ generics
- **Comprehensive component library**
- **Mouse and clipboard support**
- **Flexbox layout system**

---

## Tutorial Series

### Tutorial 1: Getting Started (15-20 minutes)

**File**: [01-getting-started.md](01-getting-started.md)

**What You'll Learn:**
- Elm Architecture (Model-View-Update pattern)
- Create your first Phoenix app
- Handle keyboard events
- Manage application state
- Render dynamic content

**What You'll Build:**
- Interactive counter app with keyboard controls

**Prerequisites:**
- Go 1.25+
- Basic Go knowledge (structs, methods)
- Terminal familiarity

---

### Tutorial 2: Building Components (30-40 minutes)

**File**: [02-building-components.md](02-building-components.md)

**What You'll Learn:**
- Use `phoenix/style` for colors, borders, spacing
- Integrate `phoenix/components/input` for text input
- Use `phoenix/components/list` for selectable lists
- Compose multiple components
- Handle component messages and delegation
- Focus management

**What You'll Build:**
- Full-featured TODO list application with:
  - Text input for adding TODOs
  - List component for displaying/navigating
  - Styled UI with borders and colors
  - Toggle completion, delete items

**Prerequisites:**
- Tutorial 1 complete
- Understanding of Elm Architecture

---

### Tutorial 3: Advanced Patterns (60-75 minutes)

**File**: [03-advanced-patterns.md](03-advanced-patterns.md)

**What You'll Learn:**
- Build custom reusable components
- Handle mouse events (click, drag, hover, scroll)
- Implement clipboard operations (copy/paste/cut)
- Use Flexbox layout for responsive UIs
- Manage complex state across components
- Optimize performance for large datasets
- Async commands and non-blocking operations

**What You'll Build:**
- Note-taking application with:
  - Multiple notes with tab navigation
  - Rich text editor with mouse support
  - Click to position cursor
  - Drag to select text
  - Copy/paste with keyboard shortcuts
  - Responsive layout
  - Auto-save functionality

**Prerequisites:**
- Tutorials 1-2 complete
- Advanced Go knowledge (interfaces, channels)
- Solid understanding of component patterns

---

## Learning Path

### For Beginners

1. **Start with Tutorial 1** - Get comfortable with the basics
   - Spend time understanding the Elm Architecture
   - Complete all exercises
   - Experiment with the counter app

2. **Move to Tutorial 2** - Learn component integration
   - Build the TODO app step-by-step
   - Don't skip the "Understanding" sections
   - Try the exercises before looking at solutions

3. **Finish with Tutorial 3** - Master advanced patterns
   - This is dense - take breaks!
   - Build the note-taking app in parts
   - Reference this tutorial when building your own apps

### For Experienced Developers

If you're already familiar with TUI frameworks or Elm Architecture:

1. **Skim Tutorial 1** - Focus on Phoenix-specific API
2. **Read Tutorial 2** - Learn component composition patterns
3. **Deep dive Tutorial 3** - Advanced patterns you'll use daily

### For Charm/Bubbletea Users

Migrating from Charm ecosystem?

1. **Quick read Tutorial 1** - See the differences in API
2. **Study Tutorial 2** - Component patterns differ significantly
3. **Master Tutorial 3** - Performance and architecture patterns

**Migration Guide**: See [../guides/migration-from-charm.md](../guides/migration-from-charm.md)

---

## Code Examples

All tutorials include **complete, working code examples** that you can:

- Copy and paste directly
- Run immediately (`go run main.go`)
- Modify and experiment with
- Use as templates for your apps

### Running Tutorial Code

```bash
# Tutorial 1: Counter app
mkdir phoenix-counter && cd phoenix-counter
go mod init example.com/counter
go get github.com/phoenix-tui/phoenix/tea
# Copy code from tutorial 1
go run main.go

# Tutorial 2: TODO app
mkdir phoenix-todo && cd phoenix-todo
go mod init example.com/todo
go get github.com/phoenix-tui/phoenix/tea
go get github.com/phoenix-tui/phoenix/style
go get github.com/phoenix-tui/phoenix/components/input
go get github.com/phoenix-tui/phoenix/components/list
# Copy code from tutorial 2
go run main.go

# Tutorial 3: Notes app
mkdir phoenix-notes && cd phoenix-notes
go mod init example.com/notes
# Install all Phoenix packages
go get github.com/phoenix-tui/phoenix/tea
go get github.com/phoenix-tui/phoenix/style
go get github.com/phoenix-tui/phoenix/components/input
go get github.com/phoenix-tui/phoenix/components/list
go get github.com/phoenix-tui/phoenix/mouse
go get github.com/phoenix-tui/phoenix/clipboard
# Copy code from tutorial 3
go run main.go
```

---

## Tutorial Structure

Each tutorial follows a consistent structure:

1. **Overview** - What you'll learn and build
2. **Prerequisites** - Required knowledge and software
3. **Step-by-Step Instructions** - Clear, numbered steps
4. **Understanding Sections** - Deep dives into concepts
5. **Exercises** - Practice challenges with solutions
6. **Common Issues** - Troubleshooting guide
7. **Summary** - Recap and next steps

### Time Estimates

All time estimates assume:
- You're typing the code (not copy-pasting)
- You read the "Understanding" sections
- You try at least one exercise

If you're just reading through or copy-pasting, you can complete much faster.

---

## Tips for Success

### 1. Type the Code Yourself

**Don't just copy-paste!** Typing helps you:
- Remember the APIs
- Catch typos (learn from errors)
- Build muscle memory
- Understand the patterns

### 2. Read the "Understanding" Sections

These explain **why**, not just **how**:
- How Phoenix event loop works
- Why immutability matters
- How components communicate
- When to use which patterns

### 3. Do the Exercises

Each tutorial has exercises that:
- Reinforce concepts
- Introduce variations
- Build confidence
- Prepare you for real projects

### 4. Experiment and Break Things

**Best way to learn:**
- Change values and see what happens
- Remove code and observe errors
- Add features not in the tutorial
- Break the app, then fix it

### 5. Use the Common Issues Sections

**Everyone hits these problems:**
- Component not updating
- Keys not responding
- Styles not applying
- Performance issues

Check the Common Issues section before asking for help!

---

## Additional Resources

### Documentation

- [API Reference](../../api/) - Complete API documentation
- [Component Guide](../../components/) - All built-in components
- [User Guides](../guides/) - Best practices, patterns, deployment

### Examples

- [Phoenix Examples](../../../examples/) - Full example applications
- [Component Examples](../../../components/) - Isolated component demos

### Community

- **GitHub**: [github.com/phoenix-tui/phoenix](https://github.com/phoenix-tui/phoenix)
- **Discussions**: GitHub Discussions for Q&A
- **Issues**: Bug reports and feature requests
- **Discord**: Real-time chat (link in README)

### Getting Help

1. **Check Common Issues** in tutorials
2. **Search GitHub Issues** - someone may have asked already
3. **Read API Documentation** - detailed parameter descriptions
4. **Ask in Discussions** - friendly community!
5. **Join Discord** - real-time help

---

## What's Next?

After completing all tutorials, you're ready to:

### Build Real Applications

Ideas to try:
- File manager with mouse support
- Log viewer with search/filter
- Database client with tables
- Chat application
- System monitor dashboard
- Git TUI client
- Music player interface

### Contribute to Phoenix

Ways to contribute:
- Report bugs you find
- Suggest new features
- Write more tutorials
- Improve documentation
- Submit PRs for fixes/features
- Help others in Discussions

### Deep Dive Topics

Advanced topics to explore:
- [Performance Optimization](../guides/performance.md)
- [Testing TUI Applications](../guides/testing.md)
- [Deployment Strategies](../guides/deployment.md)
- [Custom Component Architecture](../guides/custom-components.md)
- [Styling Best Practices](../guides/styling.md)

---

## Feedback

We want to make these tutorials better!

**Found an issue?**
- Typo or error: Open a PR
- Unclear explanation: Open an issue
- Missing topic: Suggest in Discussions

**Have a suggestion?**
- New tutorial topic
- Better examples
- Additional exercises
- More diagrams/visuals

**Share your creations!**
- Built something cool with Phoenix?
- Share in GitHub Discussions
- We might feature it in examples!

---

## Tutorial Versions

These tutorials are for **Phoenix TUI v0.1.0**.

**API Stability:**
- v0.1.x: Minor breaking changes possible
- v1.0.0: API stable (semantic versioning)

**Updating Tutorials:**
- We update tutorials with each release
- Check tutorial footer for "Last updated" date
- If API changed, tutorials will reflect it

---

## Comparison with Other Frameworks

### Phoenix vs Bubbletea/Charm

**Performance:**
- Phoenix: 29,000 FPS (render benchmark)
- Bubbletea: ~60 FPS typical
- **Result**: 489x faster

**Unicode Support:**
- Phoenix: Perfect (fixes Lipgloss #562)
- Lipgloss: Broken emoji/CJK width for months

**Architecture:**
- Phoenix: DDD + Rich Models + Hexagonal
- Charm: Traditional MVC
- **Result**: More testable, maintainable

**Coverage:**
- Phoenix: 90%+ (domain 95%+)
- Charm: Variable
- **Result**: Higher quality, fewer bugs

See [Comparison Guide](../guides/comparison.md) for details.

---

## Tutorial Authors

These tutorials were created by the Phoenix TUI team with input from:
- Early adopters and beta testers
- Community feedback
- Real-world project experience

**Contributing:**
If you find ways to improve these tutorials, please open a PR!

---

## License

Tutorials are part of Phoenix TUI Framework documentation.

- **Code examples**: MIT License (use freely)
- **Documentation**: CC BY 4.0 (attribute to Phoenix TUI)

---

*Phoenix TUI Framework - Tutorial Series*
*Version: 1.0.0 (for Phoenix v0.1.0)*
*Last updated: 2025-01-04*

Happy coding! 🚀
