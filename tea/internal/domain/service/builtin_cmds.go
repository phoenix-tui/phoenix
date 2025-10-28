// Package service provides domain services for Phoenix Tea event loop.
// This includes built-in commands for common operations.
package service

import (
	"fmt"
	"time"

	model2 "github.com/phoenix-tui/phoenix/tea/internal/domain/model"
)

// Quit returns a command that sends QuitMsg to signal program termination.
//
// This is the standard way to exit a Phoenix application. When the event loop
// receives QuitMsg, it will perform cleanup and terminate gracefully.
//
// Example:
//
//	func (m AppModel) Update(msg Msg) (Model[AppModel], Cmd) {
//		switch msg := msg.(type) {
//		case KeyMsg:
//			if msg.String() == "q" || msg.String() == "ctrl+c" {
//				return m, Quit()
//			}
//		}
//		return m, nil
//	}
func Quit() model2.Cmd {
	return func() model2.Msg {
		return model2.QuitMsg{}
	}
}

// PrintlnMsg is sent by the Println command.
//
// This is useful for debugging - you can print messages without breaking
// the TUI rendering. The event loop can handle these specially (e.g., log to file).
type PrintlnMsg struct {
	Message string
}

// String returns a human-readable representation.
func (p PrintlnMsg) String() string {
	return fmt.Sprintf("println: %s", p.Message)
}

// Println returns a command that sends a PrintlnMsg with the given message.
//
// This is primarily for debugging during development. You can use it to
// log events without disrupting the TUI display.
//
// Example:
//
//	func (m AppModel) Update(msg Msg) (Model[AppModel], Cmd) {
//		switch msg := msg.(type) {
//		case DataLoadedMsg:
//			return m, Println("Data loaded successfully")
//		}
//		return m, nil
//	}
//
// Note: In production, consider using proper logging instead.
func Println(message string) model2.Cmd {
	return func() model2.Msg {
		return PrintlnMsg{Message: message}
	}
}

// TickMsg is sent by the Tick command after the specified duration.
//
// This is useful for animations, periodic updates, or anything that needs
// to happen on a schedule. The Time field contains when the tick occurred.
type TickMsg struct {
	Time time.Time
}

// String returns a human-readable representation.
func (t TickMsg) String() string {
	return fmt.Sprintf("tick at %s", t.Time.Format(time.RFC3339))
}

// Tick returns a command that waits for the specified duration then sends TickMsg.
//
// This runs in a goroutine, so it won't block the event loop. After the duration
// elapses, a TickMsg is sent to Update with the current time.
//
// Example - Simple animation:
//
//	func (m AnimModel) Init() Cmd {
//		return Tick(100 * time.Millisecond)
//	}
//
//	func (m AnimModel) Update(msg Msg) (Model[AnimModel], Cmd) {
//		switch msg := msg.(type) {
//		case TickMsg:
//			m.frame = (m.frame + 1) % len(m.frames)
//			return m, Tick(100 * time.Millisecond) // Next frame
//		}
//		return m, nil
//	}
//
// Example - Periodic data refresh:
//
//	func (m DashboardModel) Update(msg Msg) (Model[DashboardModel], Cmd) {
//		switch msg := msg.(type) {
//		case TickMsg:
//			return m, Batch(
//				LoadData(),
//				Tick(5 * time.Second), // Refresh every 5s
//			)
//		}
//		return m, nil
//	}
func Tick(duration time.Duration) model2.Cmd {
	return func() model2.Msg {
		time.Sleep(duration)
		return TickMsg{Time: time.Now()}
	}
}
