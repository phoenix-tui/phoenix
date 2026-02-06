# Phoenix TUI - TTY Control Guide

> **Guide to running external processes from TUI applications**

---

## Overview

Phoenix provides three levels of TTY control for running external commands:

| Level | API | When to use |
|-------|-----|-------------|
| **Level 1** | `ExecProcess()` | Simple commands (vim, less, git) |
| **Level 1+** | `Suspend()` / `Resume()` | Manual control (custom workflows) |
| **Level 2** | `ExecProcessWithTTY()` | Shells, job control (bash, zsh) |

---

## Level 1: ExecProcess (Simple)

The simplest way to run an external command. Suitable for 90% of use cases.

### Example: Launching vim

```go
package main

import (
    "os/exec"

    tea "github.com/phoenix-tui/phoenix/tea"
)

type Model struct {
    filename string
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "e":
            // Launch vim to edit file
            return m, m.editFile()
        case "q":
            return m, tea.Quit
        }

    case EditDoneMsg:
        // vim finished
        if msg.Err != nil {
            // Handle error
        }
    }

    return m, nil
}

func (m Model) View() string {
    return "Press 'e' to edit file, 'q' to quit\n"
}

// editFile returns a Cmd to launch vim
func (m Model) editFile() tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("vim", m.filename)
        err := tea.ExecProcess(cmd)
        return EditDoneMsg{Err: err}
    }
}

type EditDoneMsg struct {
    Err error
}

func main() {
    p := tea.NewProgram(Model{filename: "test.txt"})
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

### What happens inside ExecProcess:

1. **Suspend** - stops TUI (exit raw mode, alt screen, show cursor)
2. **Execute** - runs command with stdin/stdout/stderr
3. **Resume** - restores TUI (raw mode, alt screen, hide cursor)

---

## Level 1+: Suspend/Resume (Manual Control)

Use when you need full control over the suspension process.

### Example: Custom workflow

```go
func (m Model) runCustomWorkflow() tea.Cmd {
    return func() tea.Msg {
        // Step 1: Suspend TUI
        if err := m.program.Suspend(); err != nil {
            return WorkflowErrorMsg{err}
        }

        // Step 2: Execute multiple commands
        fmt.Println("=== Starting workflow ===")

        cmd1 := exec.Command("git", "status")
        cmd1.Stdout = os.Stdout
        cmd1.Stderr = os.Stderr
        _ = cmd1.Run()

        fmt.Println("\nPress Enter to continue...")
        bufio.NewReader(os.Stdin).ReadLine()

        cmd2 := exec.Command("git", "diff")
        cmd2.Stdout = os.Stdout
        cmd2.Stderr = os.Stderr
        _ = cmd2.Run()

        fmt.Println("\nPress Enter to return to TUI...")
        bufio.NewReader(os.Stdin).ReadLine()

        // Step 3: Resume TUI
        if err := m.program.Resume(); err != nil {
            return WorkflowErrorMsg{err}
        }

        return WorkflowDoneMsg{}
    }
}
```

### API

```go
// Suspend pauses the TUI
// - Stops input reading
// - Exits raw mode
// - Exits alt screen (if active)
// - Shows cursor
func (p *Program[T]) Suspend() error

// Resume restores the TUI
// - Enters raw mode
// - Enters alt screen (if was active)
// - Hides cursor
// - Restarts input reading
// - Redraws screen
func (p *Program[T]) Resume() error

// IsSuspended checks current state
func (p *Program[T]) IsSuspended() bool
```

### Important Notes

- `Suspend()` and `Resume()` are **idempotent** - safe to call multiple times
- Always call `Resume()` after `Suspend()`, even on errors
- Use `defer` to guarantee restoration:

```go
if err := program.Suspend(); err != nil {
    return err
}
defer program.Resume() // Always restore!

// ... your code ...
```

---

## Level 2: ExecProcessWithTTY (Advanced)

Use for launching interactive shells with full job control.

### When do you need Level 2?

| Scenario | Level 1 | Level 2 |
|----------|---------|---------|
| Launch vim/nano | ✅ | ✅ |
| Launch git/less | ✅ | ✅ |
| Launch bash/zsh | ⚠️ | ✅ |
| Ctrl+Z in child → suspend child | ❌ | ✅ |
| Nested shells | ⚠️ | ✅ |
| Job control (&, fg, bg) | ❌ | ✅ |

### Example: Launching an interactive shell

```go
func (m Model) launchShell() tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("bash")

        opts := tea.TTYOptions{
            TransferForeground: true,  // Transfer TTY control to child
            CreateProcessGroup: true,  // Create separate process group
        }

        err := m.program.ExecProcessWithTTY(cmd, opts)
        return ShellExitMsg{Err: err}
    }
}
```

### TTYOptions

```go
type TTYOptions struct {
    // TransferForeground transfers foreground process group to child
    //
    // Unix/Linux: Uses tcsetpgrp() to transfer TTY control
    // Windows: Ignored (no Unix-style process groups)
    //
    // When true:
    // - Ctrl+Z in child suspends child, NOT parent
    // - Child can use its own signal handlers
    // - Job control works correctly (fg, bg, jobs)
    TransferForeground bool

    // CreateProcessGroup creates a new process group for child
    //
    // Unix/Linux: Uses setpgid()
    // Windows: Uses CREATE_NEW_PROCESS_GROUP flag
    //
    // Recommended to enable with TransferForeground
    CreateProcessGroup bool
}
```

### Example: SSH session

```go
func (m Model) connectSSH(host string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("ssh", host)

        // SSH needs full TTY control
        opts := tea.TTYOptions{
            TransferForeground: true,
            CreateProcessGroup: true,
        }

        err := m.program.ExecProcessWithTTY(cmd, opts)
        return SSHDoneMsg{Host: host, Err: err}
    }
}
```

### Example: Python REPL

```go
func (m Model) launchPython() tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("python3")

        opts := tea.TTYOptions{
            TransferForeground: true,
            CreateProcessGroup: true,
        }

        err := m.program.ExecProcessWithTTY(cmd, opts)
        return PythonExitMsg{Err: err}
    }
}
```

---

## Level Comparison

### Level 1: ExecProcess

```
┌─────────────────────────────────────┐
│         Phoenix TUI (Parent)        │
│                                     │
│  1. Suspend()                       │
│     ├─ Stop input reader            │
│     ├─ Exit raw mode                │
│     ├─ Exit alt screen              │
│     └─ Show cursor                  │
│                                     │
│  2. cmd.Run() ─────────────────────►│ Child Process
│     (blocking)                      │ (vim, less, etc.)
│                                     │
│  3. Resume()                        │
│     ├─ Enter raw mode               │
│     ├─ Enter alt screen             │
│     ├─ Hide cursor                  │
│     └─ Restart input reader         │
└─────────────────────────────────────┘
```

### Level 2: ExecProcessWithTTY

```
┌─────────────────────────────────────┐
│         Phoenix TUI (Parent)        │
│                                     │
│  1. Suspend()                       │
│                                     │
│  2. tcsetpgrp(child_pid) ──────────►│ Child Process
│     (transfer TTY to child)         │ (bash, ssh, etc.)
│                                     │   │
│  3. Wait for child                  │   ├─ Owns TTY
│                                     │   ├─ Ctrl+Z works
│  4. tcsetpgrp(parent_pid)           │   └─ Job control
│     (reclaim TTY)                   │
│                                     │
│  5. Resume()                        │
└─────────────────────────────────────┘
```

---

## Error Handling

### Common Errors

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case ExecDoneMsg:
        if msg.Err != nil {
            switch {
            case errors.Is(msg.Err, exec.ErrNotFound):
                // Command not found
                m.status = "Command not found"

            case strings.Contains(msg.Err.Error(), "exit status"):
                // Command exited with error
                m.status = fmt.Sprintf("Command failed: %v", msg.Err)

            default:
                // Other error
                m.status = fmt.Sprintf("Error: %v", msg.Err)
            }
        } else {
            m.status = "Command completed successfully"
        }
    }
    return m, nil
}
```

### Terminal Recovery

If something goes wrong and terminal is "broken":

```bash
# In terminal (may not display correctly):
reset

# Or:
stty sane

# Or press Ctrl+J, type reset, press Ctrl+J
```

---

## Best Practices

### 1. Always use Cmd goroutine

```go
// ✅ Correct - in Cmd goroutine
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if msg.(tea.KeyMsg).String() == "e" {
        return m, func() tea.Msg {
            err := m.program.ExecProcess(cmd)
            return DoneMsg{err}
        }
    }
}

// ❌ Wrong - directly in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if msg.(tea.KeyMsg).String() == "e" {
        m.program.ExecProcess(cmd) // BLOCKS EVENT LOOP!
    }
}
```

### 2. Handle all messages

```go
type EditDoneMsg struct{ Err error }
type ShellExitMsg struct{ Err error }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case EditDoneMsg:
        // Handle editor completion
    case ShellExitMsg:
        // Handle shell exit
    }
}
```

### 3. Use Level 2 only when needed

```go
// For vim, Level 1 is sufficient
cmd := exec.Command("vim", "file.txt")
err := program.ExecProcess(cmd)

// For bash, use Level 2
cmd := exec.Command("bash")
opts := tea.TTYOptions{TransferForeground: true, CreateProcessGroup: true}
err := program.ExecProcessWithTTY(cmd, opts)
```

### 4. Test in different environments

- Linux terminal
- macOS Terminal.app / iTerm2
- Windows Terminal / PowerShell
- SSH sessions
- tmux/screen
- CI/CD (no TTY)

---

## Compatibility

| Platform | ExecProcess | Suspend/Resume | ExecProcessWithTTY |
|----------|-------------|----------------|-------------------|
| Linux | ✅ | ✅ | ✅ (tcsetpgrp) |
| macOS | ✅ | ✅ | ✅ (tcsetpgrp) |
| Windows | ✅ | ✅ | ✅ (SetConsoleMode) |
| WSL/WSL2 | ✅ | ✅ | ✅ |
| SSH | ✅ | ✅ | ✅ |
| tmux/screen | ✅ | ✅ | ✅ |
| CI (no TTY) | ✅* | ✅* | ✅* |

*Graceful fallback when TTY unavailable

---

## Real-World Example: Shell Application

Here's how a shell application might use Phoenix TTY Control:

```go
// shell/internal/ui/model.go

func (m *Model) executeCommand(input string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("bash", "-c", input)

        // Use Level 2 for interactive commands
        if m.isInteractive(input) {
            opts := tea.TTYOptions{
                TransferForeground: true,
                CreateProcessGroup: true,
            }
            err := m.program.ExecProcessWithTTY(cmd, opts)
            return CommandDoneMsg{Err: err}
        }

        // Level 1 is sufficient for simple commands
        err := m.program.ExecProcess(cmd)
        return CommandDoneMsg{Err: err}
    }
}

func (m *Model) isInteractive(input string) bool {
    interactive := []string{"vim", "nano", "less", "ssh", "python", "node", "bash", "zsh"}
    for _, cmd := range interactive {
        if strings.HasPrefix(input, cmd) {
            return true
        }
    }
    return false
}
```

---

## Troubleshooting

### Problem: Terminal doesn't restore after command

**Solution**: Ensure `Resume()` is always called:

```go
if err := program.Suspend(); err != nil {
    return err
}
defer program.Resume() // <-- Guaranteed restoration
```

### Problem: Ctrl+Z suspends TUI instead of child

**Solution**: Use Level 2 with `TransferForeground: true`:

```go
opts := tea.TTYOptions{
    TransferForeground: true,
    CreateProcessGroup: true,
}
```

### Problem: Command hangs in CI/CD

**Solution**: The API automatically detects missing TTY and uses fallback. For tests:

```go
if !term.IsTerminal(int(os.Stdin.Fd())) {
    // Run command without TTY control
    return cmd.Run()
}
```

### Problem: "not a tty" error

**Cause**: stdin not connected to TTY (pipe, redirect, CI)

**Solution**: API automatically handles this case via fallback

---

## See Also

- [Phoenix Tea Documentation](../api/tea.md)
- [ExecProcess Implementation](../../tea/internal/application/program/program.go)
- [TTY Control Research](../dev/research/PLATFORM_SPECIFIC_TTY_CONTROL.md)

---

*Version: v0.2.1*
*Updated: 2026-02-06*
