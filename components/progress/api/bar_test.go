package progress

import (
	"strings"
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestNewBar(t *testing.T) {
	bar := NewBar(40)
	if bar == nil {
		t.Fatal("NewBar() returned nil")
	}
	if bar.Progress() != 0 {
		t.Errorf("Progress() = %d, expected 0", bar.Progress())
	}
}

func TestNewBarWithProgress(t *testing.T) {
	bar := NewBarWithProgress(40, 50)
	if bar.Progress() != 50 {
		t.Errorf("Progress() = %d, expected 50", bar.Progress())
	}
}

func TestBarFluentInterface(t *testing.T) {
	// Pointer chaining still works (returns pointer from NewBar)
	bar := NewBar(40).
		FillChar('▓').
		EmptyChar('▒').
		ShowPercent(true).
		Label("Test").
		SetProgress(75)

	if bar.Progress() != 75 {
		t.Errorf("Progress() = %d, expected 75", bar.Progress())
	}

	view := bar.View()
	if !strings.Contains(view, "Test") {
		t.Errorf("View() missing label: %s", view)
	}
	if !strings.Contains(view, "075%") {
		t.Errorf("View() missing percentage: %s", view)
	}
	if !strings.Contains(view, "▓") {
		t.Errorf("View() missing fill char: %s", view)
	}
	if !strings.Contains(view, "▒") {
		t.Errorf("View() missing empty char: %s", view)
	}
}

func TestBarSetProgress(t *testing.T) {
	bar := *NewBar(40) // Dereference to get value

	bar = bar.SetProgress(50) // Reassignment!
	if bar.Progress() != 50 {
		t.Errorf("SetProgress(50): Progress() = %d", bar.Progress())
	}

	bar = bar.SetProgress(100) // Reassignment!
	if bar.Progress() != 100 {
		t.Errorf("SetProgress(100): Progress() = %d", bar.Progress())
	}

	// Clamping.
	bar = bar.SetProgress(150) // Reassignment!
	if bar.Progress() != 100 {
		t.Errorf("SetProgress(150): Progress() = %d, expected 100", bar.Progress())
	}

	bar = bar.SetProgress(-10) // Reassignment!
	if bar.Progress() != 0 {
		t.Errorf("SetProgress(-10): Progress() = %d, expected 0", bar.Progress())
	}
}

func TestBarIncrement(t *testing.T) {
	bar := *NewBar(40) // Dereference to get value

	bar = bar.Increment(10) // Reassignment!
	if bar.Progress() != 10 {
		t.Errorf("After Increment(10): Progress() = %d", bar.Progress())
	}

	bar = bar.Increment(20) // Reassignment!
	if bar.Progress() != 30 {
		t.Errorf("After Increment(20): Progress() = %d", bar.Progress())
	}

	// Clamping at 100.
	bar = bar.Increment(100) // Reassignment!
	if bar.Progress() != 100 {
		t.Errorf("After Increment(100): Progress() = %d, expected 100", bar.Progress())
	}
}

func TestBarDecrement(t *testing.T) {
	bar := *NewBarWithProgress(40, 100) // Dereference to get value

	bar = bar.Decrement(10) // Reassignment!
	if bar.Progress() != 90 {
		t.Errorf("After Decrement(10): Progress() = %d", bar.Progress())
	}

	bar = bar.Decrement(20) // Reassignment!
	if bar.Progress() != 70 {
		t.Errorf("After Decrement(20): Progress() = %d", bar.Progress())
	}

	// Clamping at 0.
	bar = bar.Decrement(100) // Reassignment!
	if bar.Progress() != 0 {
		t.Errorf("After Decrement(100): Progress() = %d, expected 0", bar.Progress())
	}
}

func TestBarIsComplete(t *testing.T) {
	bar := *NewBar(40) // Dereference to get value

	if bar.IsComplete() {
		t.Errorf("IsComplete() = true at 0%%")
	}

	bar = bar.SetProgress(50) // Reassignment!
	if bar.IsComplete() {
		t.Errorf("IsComplete() = true at 50%%")
	}

	bar = bar.SetProgress(100) // Reassignment!
	if !bar.IsComplete() {
		t.Errorf("IsComplete() = false at 100%%")
	}
}

func TestBarInit(t *testing.T) {
	bar := NewBar(40)
	cmd := bar.Init()
	if cmd != nil {
		t.Errorf("Init() returned non-nil cmd")
	}
}

func TestBarUpdate(t *testing.T) {
	// NewBar().SetProgress() returns Bar after first method call.
	bar := NewBar(40).SetProgress(50)

	// Update should return self and no cmd (value semantics)
	updated, cmd := bar.Update(tea.KeyMsg{})
	if cmd != nil {
		t.Errorf("Update() returned non-nil cmd")
	}

	// Should return same value (already Bar type)
	if updated.Progress() != 50 {
		t.Errorf("Update() changed progress: %d != 50", updated.Progress())
	}
}

func TestBarView(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() Bar // Return value, not pointer
		contains    []string
		notContains []string
	}{
		{
			name: "Simple bar",
			setup: func() Bar {
				return *NewBar(10) // Dereference
			},
			contains:    []string{"░"},
			notContains: []string{"█", "%"},
		},
		{
			name: "Bar with progress",
			setup: func() Bar {
				return *NewBarWithProgress(10, 50) // Dereference
			},
			contains:    []string{"█", "░"},
			notContains: []string{"%"},
		},
		{
			name: "Bar with percentage",
			setup: func() Bar {
				// First call returns Bar after method chaining.
				return NewBarWithProgress(10, 75).ShowPercent(true)
			},
			contains: []string{"█", "░", "075%"},
		},
		{
			name: "Bar with label",
			setup: func() Bar {
				// First call returns Bar after method chaining.
				return NewBar(10).Label("Loading")
			},
			contains: []string{"Loading", "░"},
		},
		{
			name: "Full bar",
			setup: func() Bar {
				// Method chaining returns Bar (value)
				return NewBarWithProgress(10, 100).
					Label("Done").
					ShowPercent(true)
			},
			contains:    []string{"Done", "█", "100%"},
			notContains: []string{"░"},
		},
		{
			name: "Custom chars",
			setup: func() Bar {
				// Method chaining returns Bar (value)
				return NewBarWithProgress(10, 50).
					FillChar('▓').
					EmptyChar('▒')
			},
			contains:    []string{"▓", "▒"},
			notContains: []string{"█", "░"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := tt.setup()
			view := bar.View()

			for _, substr := range tt.contains {
				if !strings.Contains(view, substr) {
					t.Errorf("View() missing %q: %s", substr, view)
				}
			}

			for _, substr := range tt.notContains {
				if strings.Contains(view, substr) {
					t.Errorf("View() should not contain %q: %s", substr, view)
				}
			}
		})
	}
}

func TestBarMethodChaining(t *testing.T) {
	// Verify all methods return Bar for chaining (value semantics)
	// Pointer chaining from NewBar.
	bar := NewBar(40).
		FillChar('▓').
		EmptyChar('▒').
		ShowPercent(true).
		Label("Test").
		SetProgress(50)

	if bar.Progress() != 50 {
		t.Errorf("Progress() = %d, expected 50", bar.Progress())
	}

	// Value chaining (must reassign)
	// After first method call, we have Bar value.
	barValue := bar.Increment(10) // bar is *Bar, returns Bar
	if barValue.Progress() != 60 {
		t.Errorf("After Increment: Progress() = %d, expected 60", barValue.Progress())
	}

	barValue = barValue.Decrement(5)
	if barValue.Progress() != 55 {
		t.Errorf("After Decrement: Progress() = %d, expected 55", barValue.Progress())
	}
}

func TestBarTeaModelContract(_ *testing.T) {
	// Verify Bar implements tea model contract (Init, Update, View)
	bar := *NewBar(40) // Dereference to get value

	// Init returns Cmd.
	cmd := bar.Init()
	_ = cmd

	// Update returns (Bar, Cmd) - value semantics!
	updated, cmd2 := bar.Update(tea.KeyMsg{})
	_ = updated
	_ = cmd2

	// View returns string.
	view := bar.View()
	_ = view
}

func TestBarProgressRange(t *testing.T) {
	bar := *NewBar(100) // Dereference to get value

	// Test all percentages 0-100.
	for pct := 0; pct <= 100; pct++ {
		bar = bar.SetProgress(pct) // Reassignment!
		if bar.Progress() != pct {
			t.Errorf("SetProgress(%d): Progress() = %d", pct, bar.Progress())
		}

		view := bar.View()
		// Should have 100 characters (all filled or empty)
		totalChars := strings.Count(view, "█") + strings.Count(view, "░")
		if totalChars != 100 {
			t.Errorf("At %d%%: total chars = %d, expected 100", pct, totalChars)
		}
	}
}
