//go:build linux

// Package program provides TTY control for Linux platforms.
// Uses ioctl TIOCGPGRP/TIOCSPGRP for process group control.
package program

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"
)

// TTYOptions configures TTY control behavior for external commands.
type TTYOptions struct {
	// TransferForeground controls whether to transfer foreground process group.
	// Linux: Uses ioctl TIOCSPGRP to make child foreground.
	// This enables proper job control (Ctrl+Z in child won't affect parent).
	TransferForeground bool

	// CreateProcessGroup creates a new process group for the child.
	// Recommended when TransferForeground is true.
	CreateProcessGroup bool
}

// TIOCGPGRP and TIOCSPGRP ioctl constants for Linux
const (
	TIOCGPGRP = 0x540F // Get process group (linux)
	TIOCSPGRP = 0x5410 // Set process group (linux)
)

// tcgetpgrp gets the foreground process group of the terminal.
func tcgetpgrp(fd int) (int32, error) {
	var pgrp int32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TIOCGPGRP, uintptr(unsafe.Pointer(&pgrp)))
	if errno != 0 {
		return 0, errno
	}
	return pgrp, nil
}

// tcsetpgrp sets the foreground process group of the terminal.
func tcsetpgrp(fd int, pgrp int32) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TIOCSPGRP, uintptr(unsafe.Pointer(&pgrp)))
	if errno != 0 {
		return errno
	}
	return nil
}

// execWithTTYControl executes a command with full TTY control on Linux.
//
// This provides Level 2 TTY control with proper foreground process group transfer.
// Uses ioctl TIOCSPGRP to give the child process full terminal control, enabling:
//   - Proper job control signals (Ctrl+Z suspends child, not parent)
//   - Child can use its own signal handlers
//   - Proper TTY ownership for nested shells
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
	parentPgid, err := tcgetpgrp(ttyFD)
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
		if err := tcsetpgrp(ttyFD, int32(childPid)); err != nil {
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
		if err := tcsetpgrp(ttyFD, parentPgid); err != nil {
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
