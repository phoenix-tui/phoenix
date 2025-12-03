//go:build unix || darwin

// Package program provides TTY control for Unix/Linux/macOS platforms.
package program

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

// TTYOptions configures TTY control behavior for external commands.
type TTYOptions struct {
	// TransferForeground controls whether to transfer foreground process group.
	// Unix/Linux: Uses tcsetpgrp() to make child foreground.
	// This enables proper job control (Ctrl+Z in child won't affect parent).
	TransferForeground bool

	// CreateProcessGroup creates a new process group for the child.
	// Recommended when TransferForeground is true.
	CreateProcessGroup bool
}

// execWithTTYControl executes a command with full TTY control on Unix platforms.
//
// This provides Level 2 TTY control with proper foreground process group transfer.
// Uses tcsetpgrp() to give the child process full terminal control, enabling:
//   - Proper job control signals (Ctrl+Z suspends child, not parent)
//   - Child can use its own signal handlers
//   - Proper TTY ownership for nested shells
//
// Implementation follows the Unix pattern:
//  1. Suspend TUI (exit raw mode, exit alt screen, show cursor)
//  2. Get current foreground process group (to restore later)
//  3. Ignore SIGTTOU (prevent being stopped during tcsetpgrp)
//  4. Create new process group for child (if requested)
//  5. Start child process
//  6. Transfer foreground to child via tcsetpgrp (from parent!)
//  7. Wait for child to complete
//  8. Reclaim foreground via tcsetpgrp
//  9. Resume TUI (restore raw mode, alt screen, restart input)
//
// CRITICAL: tcsetpgrp() MUST be called from parent, not child!
// See: https://github.com/golang/go/issues/37217
//
// Returns error if command execution or TTY control fails.
func (p *Program[T]) execWithTTYControl(cmd *exec.Cmd, opts TTYOptions) error {
	// Validate command
	if cmd == nil {
		return fmt.Errorf("exec: cmd is nil")
	}

	// Get TTY file descriptor (stdin)
	ttyFD := int(os.Stdin.Fd())

	// Get current foreground process group (to restore later)
	parentPgid, err := unix.Tcgetpgrp(ttyFD)
	if err != nil {
		// Not a TTY - fall back to simple ExecProcess
		log.Printf("WARNING: tcgetpgrp failed (not a TTY?): %v - falling back to simple exec", err)
		return p.ExecProcess(cmd)
	}

	// Temporarily ignore SIGTTOU to prevent being stopped
	// when we call tcsetpgrp() from a background process group
	signal.Ignore(syscall.SIGTTOU)
	defer signal.Reset(syscall.SIGTTOU)

	// STEP 1: Suspend TUI (stop input, exit raw mode, exit alt screen, show cursor)
	if err := p.Suspend(); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	// STEP 2: Configure process group
	if opts.CreateProcessGroup {
		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		cmd.SysProcAttr.Setpgid = true
		cmd.SysProcAttr.Pgid = 0 // Use child PID as PGID
	}

	// STEP 3: Setup command I/O - give it full terminal control
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// STEP 4: Start child process
	if err := cmd.Start(); err != nil {
		// Restore TUI before returning error
		if resumeErr := p.Resume(); resumeErr != nil {
			return fmt.Errorf("exec: failed to start command (%v) and %w", err, resumeErr)
		}
		return fmt.Errorf("exec: failed to start command: %w", err)
	}

	childPid := cmd.Process.Pid

	// STEP 5: Transfer foreground to child (if requested)
	// CRITICAL: This MUST be called from parent, not child!
	if opts.TransferForeground {
		if err := unix.Tcsetpgrp(ttyFD, int32(childPid)); err != nil {
			// Failed to transfer - kill child and restore TUI
			_ = cmd.Process.Kill()
			if resumeErr := p.Resume(); resumeErr != nil {
				return fmt.Errorf("exec: failed to transfer foreground (%v) and %w", err, resumeErr)
			}
			return fmt.Errorf("exec: failed to transfer foreground: %w", err)
		}
	}

	// STEP 6: Wait for child to complete (BLOCKING)
	cmdErr := cmd.Wait()

	// STEP 7: Reclaim foreground (restore parent as foreground)
	if opts.TransferForeground {
		if err := unix.Tcsetpgrp(ttyFD, parentPgid); err != nil {
			// Critical error - log but continue with Resume
			// This can happen if child process group is orphaned
			log.Printf("WARNING: failed to restore TTY foreground: %v", err)
		}
	}

	// STEP 8: Resume TUI (restore raw mode, alt screen, restart input, redraw)
	// ALWAYS resume, even if command failed
	if err := p.Resume(); err != nil {
		if cmdErr != nil {
			return fmt.Errorf("exec: command failed (%v) and %w", cmdErr, err)
		}
		return fmt.Errorf("exec: %w", err)
	}

	// Return original command error (if any)
	return cmdErr
}
