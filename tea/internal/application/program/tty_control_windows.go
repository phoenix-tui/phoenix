//go:build windows

// Package program provides TTY control for Windows platforms.
package program

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// TTYOptions configures TTY control behavior for external commands.
type TTYOptions struct {
	// TransferForeground is ignored on Windows (no process groups like Unix).
	// Windows uses console mode control instead.
	TransferForeground bool

	// CreateProcessGroup creates a new process group for the child.
	// On Windows, this uses CREATE_NEW_PROCESS_GROUP flag.
	CreateProcessGroup bool
}

// execWithTTYControl executes a command with enhanced console control on Windows.
//
// This provides Level 2 TTY control with console mode management.
// Windows doesn't have Unix-style process groups and tcsetpgrp(), so we use:
//   - SetConsoleMode() to control console behavior
//   - CREATE_NEW_PROCESS_GROUP flag for process isolation
//
// Implementation:
//  1. Suspend TUI (exit raw mode, exit alt screen, show cursor)
//  2. Get console handle and save current mode
//  3. Disable virtual terminal processing (let child control console)
//  4. Create new process group (if requested)
//  5. Setup command I/O
//  6. Run child process
//  7. Resume TUI (restore raw mode, alt screen, restart input)
//
// Returns error if command execution or console control fails.
func (p *Program[T]) execWithTTYControl(cmd *exec.Cmd, opts TTYOptions) error {
	// Validate command
	if cmd == nil {
		return fmt.Errorf("exec: cmd is nil")
	}

	// STEP 1: Suspend TUI (stop input, exit raw mode, exit alt screen, show cursor)
	if err := p.Suspend(); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	// STEP 2: Get console handle (stdin)
	consoleHandle, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		// Not a console - fall back to simple ExecProcess
		log.Printf("WARNING: GetStdHandle failed (not a console?): %v - falling back to simple exec", err)
		return p.ExecProcess(cmd)
	}

	// STEP 3: Save current console mode
	var originalMode uint32
	if err := windows.GetConsoleMode(consoleHandle, &originalMode); err != nil {
		// Not a console (e.g., piped input) - fall back to simple ExecProcess
		log.Printf("WARNING: GetConsoleMode failed (not a console?): %v - falling back to simple exec", err)
		return p.ExecProcess(cmd)
	}

	// STEP 4: Disable virtual terminal processing (let child control console)
	childMode := originalMode &^ windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(consoleHandle, childMode); err != nil {
		// Non-fatal - continue anyway
		log.Printf("WARNING: SetConsoleMode failed: %v - continuing with original mode", err)
	}

	// STEP 5: Configure child process group (if requested)
	if opts.CreateProcessGroup {
		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
	}

	// STEP 6: Setup command I/O - give it full terminal control
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// STEP 7: Run child process (BLOCKING)
	cmdErr := cmd.Run()

	// STEP 8: Resume TUI (restore raw mode, alt screen, restart input, redraw)
	// ALWAYS resume, even if command failed.
	// NOTE: No defer SetConsoleMode(originalMode) — Resume() calls EnterRawMode()
	// which sets the correct console flags. A defer would fire AFTER Resume(),
	// restoring cooked-mode flags and undoing EnterRawMode (see bug report).
	if err := p.Resume(); err != nil {
		if cmdErr != nil {
			return fmt.Errorf("exec: command failed (%w) and %w", cmdErr, err)
		}
		return fmt.Errorf("exec: %w", err)
	}

	// Return original command error (if any)
	return cmdErr
}
