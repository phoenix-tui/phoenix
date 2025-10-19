package progress

import (
	"strings"
	"testing"
	"time"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestNewSpinner(t *testing.T) {
	tests := []string{
		"dots",
		"line",
		"arrow",
		"circle",
		"bounce",
	}

	for _, style := range tests {
		t.Run(style, func(t *testing.T) {
			spinner := NewSpinner(style)
			if spinner == nil {
				t.Fatalf("NewSpinner(%s) returned nil", style)
			}
		})
	}
}

func TestNewSpinnerUnknownStyle(t *testing.T) {
	// Should fallback to default (dots)
	spinner := NewSpinner("unknown-style-name")
	if spinner == nil {
		t.Fatal("NewSpinner() with unknown style returned nil")
	}
}

func TestSpinnerLabel(t *testing.T) {
	spinner := *NewSpinner("dots")        // Dereference to get value
	spinner = spinner.Label("Loading...") // Reassignment needed!

	view := spinner.View()
	if !strings.Contains(view, "Loading...") {
		t.Errorf("View() missing label: %s", view)
	}
}

func TestSpinnerInit(t *testing.T) {
	spinner := NewSpinner("dots")
	cmd := spinner.Init()

	if cmd == nil {
		t.Fatal("Init() returned nil cmd")
	}

	// Execute cmd to get tea.TickMsg
	msg := cmd()
	if _, ok := msg.(tea.TickMsg); !ok {
		t.Errorf("Init() cmd didn't return tea.TickMsg, got %T", msg)
	}
}

func TestSpinnerUpdate(t *testing.T) {
	spinner := NewSpinner("line") // 4 frames: | / - \

	// Initial frame
	view1 := spinner.View()

	// Send tea.TickMsg
	updated, cmd := spinner.Update(tea.TickMsg{Time: time.Now()})

	// Should return spinner and tick cmd
	if cmd == nil {
		t.Error("Update() returned nil cmd")
	}

	// View should change (different frame) - updated is already *Spinner type
	view2 := updated.View()
	if view1 == view2 {
		// Note: might be same if single frame, but "line" has 4 frames
		t.Logf("Warning: views are same: %s", view1)
	}
}

func TestSpinnerUpdateOtherMsg(t *testing.T) {
	spinner := *NewSpinner("dots") // Dereference to get value

	// Send non-TickMsg
	updated, cmd := spinner.Update(tea.KeyMsg{})

	if cmd != nil {
		t.Errorf("Update() with KeyMsg returned cmd")
	}

	// Should return same spinner (value semantics - comparison works)
	if updated.View() != spinner.View() {
		t.Errorf("Update() with KeyMsg returned different spinner")
	}
}

func TestSpinnerView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() Spinner // Return value, not pointer
		check func(t *testing.T, view string)
	}{
		{
			name: "Simple spinner",
			setup: func() Spinner {
				return *NewSpinner("dots") // Dereference!
			},
			check: func(t *testing.T, view string) {
				if view == "" {
					t.Error("View() is empty")
				}
				// Should contain Unicode spinner char
				if len(view) < 1 {
					t.Errorf("View() too short: %q", view)
				}
			},
		},
		{
			name: "Spinner with label",
			setup: func() Spinner {
				s := *NewSpinner("dots")     // Dereference!
				return s.Label("Loading...") // Reassignment via return
			},
			check: func(t *testing.T, view string) {
				if !strings.Contains(view, "Loading...") {
					t.Errorf("View() missing label: %s", view)
				}
				// Should have format: "⠋ Loading..."
				parts := strings.Split(view, " ")
				if len(parts) < 2 {
					t.Errorf("View() format incorrect: %s", view)
				}
			},
		},
		{
			name: "Line spinner",
			setup: func() Spinner {
				return *NewSpinner("line") // Dereference!
			},
			check: func(t *testing.T, view string) {
				// Should be one of: | / - \
				validFrames := []string{"|", "/", "-", "\\"}
				valid := false
				for _, frame := range validFrames {
					if view == frame {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("View() = %q, expected one of %v", view, validFrames)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spinner := tt.setup()
			view := spinner.View()
			tt.check(t, view)
		})
	}
}

func TestSpinnerAnimation(t *testing.T) {
	spinner := *NewSpinner("line") // Dereference! 4 frames: | / - \

	frames := make([]string, 0, 5)
	frames = append(frames, spinner.View())

	// Advance 4 times (full cycle)
	for i := 0; i < 4; i++ {
		updated, _ := spinner.Update(tea.TickMsg{Time: time.Now()})
		spinner = updated // Reassignment (value semantics)
		frames = append(frames, spinner.View())
	}

	// First and last should be the same (wrapped around)
	if frames[0] != frames[4] {
		t.Errorf("Animation didn't wrap: frame[0]=%s, frame[4]=%s",
			frames[0], frames[4])
	}

	// All frames in between should be different
	// (for "line" spinner with 4 distinct frames)
	seen := make(map[string]bool)
	for i := 0; i < 4; i++ {
		if seen[frames[i]] {
			t.Logf("Note: frame %d repeated: %s", i, frames[i])
		}
		seen[frames[i]] = true
	}

	// Should have seen all 4 frames
	if len(seen) < 4 {
		t.Errorf("Only saw %d unique frames, expected 4: %v", len(seen), frames[:4])
	}
}

func TestSpinnerMethodChaining(t *testing.T) {
	spinner := NewSpinner("dots").
		Label("Processing...")

	view := spinner.View()
	if !strings.Contains(view, "Processing...") {
		t.Errorf("Method chaining failed: %s", view)
	}
}

func TestSpinnerTeaModelContract(t *testing.T) {
	// Verify Spinner implements tea model contract (Init, Update, View)
	spinner := NewSpinner("dots")

	// Init returns Cmd
	cmd := spinner.Init()
	_ = cmd

	// Update returns (*Spinner, Cmd)
	updated, cmd2 := spinner.Update(tea.TickMsg{})
	_ = updated
	_ = cmd2

	// View returns string
	view := spinner.View()
	_ = view
}

func TestSpinnerTickTiming(t *testing.T) {
	spinner := NewSpinner("dots") // 10 FPS

	cmd := spinner.Init()
	if cmd == nil {
		t.Fatal("Init() returned nil")
	}

	// Verify tick message
	msg := cmd()
	tickMsg, ok := msg.(tea.TickMsg)
	if !ok {
		t.Fatalf("Expected tea.TickMsg, got %T", msg)
	}

	// TickMsg should have time
	if tickMsg.Time.IsZero() {
		t.Error("TickMsg has zero time")
	}
}

func TestSpinnerEmptyLabel(t *testing.T) {
	spinner := NewSpinner("dots").Label("")

	view := spinner.View()

	// Should not have trailing space
	if strings.HasSuffix(view, " ") {
		t.Errorf("View() has trailing space: %q", view)
	}
}

func TestSpinnerAllStyles(t *testing.T) {
	// Test all pre-defined styles
	styles := []string{
		"dots", "line", "arrow", "circle", "bounce",
		"dot-pulse", "grow-vertical", "grow-horizontal",
		"box-bounce", "simple-dots", "clock", "earth",
		"moon", "toggle", "hamburger",
	}

	for _, style := range styles {
		t.Run(style, func(t *testing.T) {
			spinner := NewSpinner(style)

			// Should initialize
			cmd := spinner.Init()
			if cmd == nil {
				t.Errorf("Init() returned nil for %s", style)
			}

			// Should render
			view := spinner.View()
			if view == "" {
				t.Errorf("View() empty for %s", style)
			}

			// Should update (value semantics - cannot be nil)
			updated, _ := spinner.Update(tea.TickMsg{Time: time.Now()})
			_ = updated // updated is Spinner value (never nil)
		})
	}
}

func TestSpinnerViewFormat(t *testing.T) {
	// Without label
	spinner1 := NewSpinner("dots")
	view1 := spinner1.View()

	// Should not have leading/trailing spaces
	trimmed1 := strings.TrimSpace(view1)
	if trimmed1 != view1 {
		t.Errorf("View() without label has extra spaces: %q", view1)
	}

	// With label
	spinner2 := NewSpinner("dots").Label("Test")
	view2 := spinner2.View()

	// Should have format: "frame label" (one space)
	parts := strings.SplitN(view2, " ", 2)
	if len(parts) != 2 {
		t.Errorf("View() with label format incorrect: %q", view2)
	}
	if parts[1] != "Test" {
		t.Errorf("Label in view = %q, expected 'Test'", parts[1])
	}
}
