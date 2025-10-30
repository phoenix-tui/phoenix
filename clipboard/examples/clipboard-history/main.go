// clipboard-history demonstrates clipboard history tracking and management.
//
// This example shows how to:
// - Enable/disable clipboard history
// - View history entries
// - Monitor memory usage
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard"
	"github.com/phoenix-tui/phoenix/tea"
)

type model struct {
	clipboard      *clipboard.Clipboard
	historyEnabled bool
	entries        []clipboard.HistoryEntry
	selected       int
	message        string
	width          int
	height         int
}

func initialModel() model {
	clip, err := clipboard.New()
	if err != nil {
		panic(err)
	}

	// Enable history: 100 entries, 24-hour retention
	clip.EnableHistory(100, 24*time.Hour)

	return model{
		clipboard:      clip,
		historyEnabled: true,
		entries:        []clipboard.HistoryEntry{},
		selected:       0,
		message:        "Welcome! Press 'h' for help",
		width:          80,
		height:         24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit()

		case "?":
			m.message = "Keys: [v]iew [c]lear [e]nable/disable [m]emory [q]uit"
			return m, nil

		case "v":
			m.entries = m.clipboard.GetHistory()
			m.message = fmt.Sprintf("Refreshed history: %d entries", len(m.entries))
			return m, nil

		case "c":
			m.clipboard.ClearHistory()
			m.entries = []clipboard.HistoryEntry{}
			m.selected = 0
			m.message = "History cleared"
			return m, nil

		case "e":
			if m.historyEnabled {
				m.clipboard.DisableHistory()
				m.historyEnabled = false
				m.message = "History disabled"
			} else {
				m.clipboard.EnableHistory(100, 24*time.Hour)
				m.historyEnabled = true
				m.message = "History enabled (100 entries, 24h)"
			}
			return m, nil

		case "m":
			size := m.clipboard.GetHistorySize()
			totalSize := m.clipboard.GetHistoryTotalSize()
			removed := m.clipboard.RemoveExpiredHistory()
			m.message = fmt.Sprintf("Entries: %d | Memory: %d bytes | Expired removed: %d", size, totalSize, removed)
			m.entries = m.clipboard.GetHistory()
			return m, nil

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
			return m, nil

		case "down", "j":
			if m.selected < len(m.entries)-1 {
				m.selected++
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Header
	b.WriteString("=== Clipboard History Manager ===\n\n")

	// Status
	status := "ENABLED"
	if !m.historyEnabled {
		status = "DISABLED"
	}
	b.WriteString(fmt.Sprintf("Status: %s | Entries: %d | Provider: %s\n",
		status, len(m.entries), m.clipboard.GetProviderName()))
	b.WriteString(fmt.Sprintf("Memory: %d bytes total\n\n",
		m.clipboard.GetHistoryTotalSize()))

	// History list
	b.WriteString("History (newest first):\n")
	b.WriteString(strings.Repeat("-", min(m.width, 80)) + "\n")

	if len(m.entries) == 0 {
		b.WriteString("  No history entries yet\n")
		b.WriteString("  Use clipboard normally, then press 'v' to view history\n")
	} else {
		displayCount := min(len(m.entries), 10)
		for i := 0; i < displayCount; i++ {
			entry := m.entries[i]

			// Selection indicator
			if i == m.selected {
				b.WriteString("> ")
			} else {
				b.WriteString("  ")
			}

			// Entry details
			age := time.Since(entry.Timestamp)
			ageStr := formatDuration(age)

			// Content preview
			preview := string(entry.Content)
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", " ")
			preview = strings.ReplaceAll(preview, "\r", " ")

			b.WriteString(fmt.Sprintf("[%s] %s (%s, %d bytes)\n",
				ageStr, preview, entry.MIMEType, entry.Size))
		}

		if len(m.entries) > displayCount {
			b.WriteString(fmt.Sprintf("  ... and %d more entries\n", len(m.entries)-displayCount))
		}
	}

	// Footer
	b.WriteString(strings.Repeat("-", min(m.width, 80)) + "\n")
	b.WriteString(m.message + "\n")
	b.WriteString("Press '?' for help, 'q' to quit\n")

	return b.String()
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	p := tea.New(initialModel(), tea.WithAltScreen[model]())
	if err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
