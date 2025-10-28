package service

import (
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

func TestClickDetector_SingleClick(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)

	// Create a release event
	event := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)

	clickEvent := detector.DetectClick(event)
	if clickEvent == nil {
		t.Fatal("Expected click event, got nil")
	}

	if clickEvent.Type() != value2.EventClick {
		t.Errorf("Expected EventClick, got %v", clickEvent.Type())
	}

	if detector.ClickCount() != 1 {
		t.Errorf("Expected click count 1, got %d", detector.ClickCount())
	}
}

func TestClickDetector_DoubleClick(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)

	// First click
	event1 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
	)
	detector.DetectClick(event1)

	// Second click (immediately after)
	event2 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
	)
	clickEvent := detector.DetectClick(event2)

	if clickEvent == nil {
		t.Fatal("Expected click event, got nil")
	}

	if clickEvent.Type() != value2.EventDoubleClick {
		t.Errorf("Expected EventDoubleClick, got %v", clickEvent.Type())
	}

	if detector.ClickCount() != 2 {
		t.Errorf("Expected click count 2, got %d", detector.ClickCount())
	}
}

func TestClickDetector_TripleClick(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)

	// Three rapid clicks
	for i := 0; i < 3; i++ {
		event := model.NewMouseEvent(
			value2.EventRelease,
			value2.ButtonLeft,
			pos,
			value2.ModifierNone,
		)
		detector.DetectClick(event)
	}

	if detector.ClickCount() != 3 {
		t.Errorf("Expected click count 3, got %d", detector.ClickCount())
	}
}

func TestClickDetector_TimeoutReset(t *testing.T) {
	detector := NewClickDetector(100*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)

	// First click
	event1 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		time.Now(),
	)
	detector.DetectClick(event1)

	// Second click after timeout
	event2 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		time.Now().Add(200*time.Millisecond),
	)
	clickEvent := detector.DetectClick(event2)

	// Should be single click (reset after timeout)
	if clickEvent.Type() != value2.EventClick {
		t.Errorf("Expected EventClick after timeout, got %v", clickEvent.Type())
	}

	if detector.ClickCount() != 1 {
		t.Errorf("Expected click count reset to 1, got %d", detector.ClickCount())
	}
}

func TestClickDetector_DifferentPosition(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)

	// First click at (10, 5)
	event1 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)
	detector.DetectClick(event1)

	// Second click at (15, 5) - beyond tolerance
	event2 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		value2.NewPosition(15, 5),
		value2.ModifierNone,
	)
	clickEvent := detector.DetectClick(event2)

	// Should be single click (different position)
	if clickEvent.Type() != value2.EventClick {
		t.Errorf("Expected EventClick for different position, got %v", clickEvent.Type())
	}

	if detector.ClickCount() != 1 {
		t.Errorf("Expected click count reset to 1, got %d", detector.ClickCount())
	}
}

func TestClickDetector_Reset(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)

	// Create a click
	event := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)
	detector.DetectClick(event)

	// Reset
	detector.Reset()

	if detector.ClickCount() != 0 {
		t.Errorf("Expected click count 0 after reset, got %d", detector.ClickCount())
	}
}

func TestClickDetector_IgnoreNonRelease(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)

	// Try with press event
	event := model.NewMouseEvent(
		value2.EventPress,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)

	clickEvent := detector.DetectClick(event)
	if clickEvent != nil {
		t.Errorf("Expected nil for press event, got %v", clickEvent.Type())
	}
}

// Edge case: Rapid clicking (< 500ms)
func TestClickDetector_RapidClicking(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)
	baseTime := time.Now()

	// Click at 0ms
	event1 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime,
	)
	click1 := detector.DetectClick(event1)
	if click1.Type() != value2.EventClick {
		t.Errorf("First click: expected EventClick, got %v", click1.Type())
	}

	// Click at 100ms (well within timeout)
	event2 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime.Add(100*time.Millisecond),
	)
	click2 := detector.DetectClick(event2)
	if click2.Type() != value2.EventDoubleClick {
		t.Errorf("Second click: expected EventDoubleClick, got %v", click2.Type())
	}

	// Click at 250ms (still within timeout from last click)
	event3 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime.Add(250*time.Millisecond),
	)
	click3 := detector.DetectClick(event3)
	if click3.Type() != value2.EventTripleClick {
		t.Errorf("Third click: expected EventTripleClick, got %v", click3.Type())
	}

	// Click at 400ms (should reset after triple click)
	event4 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime.Add(400*time.Millisecond),
	)
	click4 := detector.DetectClick(event4)
	if click4.Type() != value2.EventClick {
		t.Errorf("Fourth click: expected EventClick (reset after triple), got %v", click4.Type())
	}
}

// Edge case: Exactly 500ms boundary
func TestClickDetector_ExactTimeoutBoundary(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)
	baseTime := time.Now()

	// First click
	event1 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime,
	)
	detector.DetectClick(event1)

	// Second click exactly at 500ms (should still be within timeout, <= comparison)
	event2 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime.Add(500*time.Millisecond),
	)
	click2 := detector.DetectClick(event2)
	if click2.Type() != value2.EventDoubleClick {
		t.Errorf("Expected EventDoubleClick at exactly 500ms, got %v", click2.Type())
	}

	// Reset for next test
	detector.Reset()

	// Third click at 501ms (should be outside timeout)
	event3 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime,
	)
	detector.DetectClick(event3)

	event4 := model.NewMouseEventWithTimestamp(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
		baseTime.Add(501*time.Millisecond),
	)
	click4 := detector.DetectClick(event4)
	if click4.Type() != value2.EventClick {
		t.Errorf("Expected EventClick at 501ms (timeout), got %v", click4.Type())
	}
}

// Edge case: Position tolerance (exactly at boundary)
func TestClickDetector_PositionTolerance(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	basePos := value2.NewPosition(10, 10)

	// First click at (10, 10)
	event1 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		basePos,
		value2.ModifierNone,
	)
	detector.DetectClick(event1)

	tests := []struct {
		name     string
		pos      value2.Position
		expected value2.EventType
	}{
		{
			name:     "Same position",
			pos:      value2.NewPosition(10, 10),
			expected: value2.EventDoubleClick,
		},
		{
			name:     "Within tolerance: +1 X",
			pos:      value2.NewPosition(11, 10),
			expected: value2.EventDoubleClick,
		},
		{
			name:     "Within tolerance: -1 X",
			pos:      value2.NewPosition(9, 10),
			expected: value2.EventDoubleClick,
		},
		{
			name:     "Within tolerance: +1 Y",
			pos:      value2.NewPosition(10, 11),
			expected: value2.EventDoubleClick,
		},
		{
			name:     "Within tolerance: -1 Y",
			pos:      value2.NewPosition(10, 9),
			expected: value2.EventDoubleClick,
		},
		{
			name:     "Outside tolerance: +2 X",
			pos:      value2.NewPosition(12, 10),
			expected: value2.EventClick,
		},
		{
			name:     "Outside tolerance: -2 X",
			pos:      value2.NewPosition(8, 10),
			expected: value2.EventClick,
		},
		{
			name:     "Outside tolerance: +2 Y",
			pos:      value2.NewPosition(10, 12),
			expected: value2.EventClick,
		},
		{
			name:     "Outside tolerance: -2 Y",
			pos:      value2.NewPosition(10, 8),
			expected: value2.EventClick,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and create first click
			detector.Reset()
			detector.DetectClick(event1)

			// Second click at test position
			event2 := model.NewMouseEvent(
				value2.EventRelease,
				value2.ButtonLeft,
				tt.pos,
				value2.ModifierNone,
			)
			click := detector.DetectClick(event2)
			if click.Type() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, click.Type())
			}
		})
	}
}

// Edge case: Different button types
func TestClickDetector_DifferentButtons(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)

	// Left button click
	event1 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
	)
	detector.DetectClick(event1)

	// Right button click (should reset)
	event2 := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonRight,
		pos,
		value2.ModifierNone,
	)
	click := detector.DetectClick(event2)

	if click.Type() != value2.EventClick {
		t.Errorf("Expected EventClick for different button, got %v", click.Type())
	}

	if detector.ClickCount() != 1 {
		t.Errorf("Expected click count reset to 1, got %d", detector.ClickCount())
	}
}

// Edge case: Constructor with invalid parameters
func TestClickDetector_Constructor(t *testing.T) {
	tests := []struct {
		name              string
		timeout           time.Duration
		tolerance         int
		expectedTimeout   time.Duration
		expectedTolerance int
	}{
		{
			name:              "Valid parameters",
			timeout:           300 * time.Millisecond,
			tolerance:         2,
			expectedTimeout:   300 * time.Millisecond,
			expectedTolerance: 2,
		},
		{
			name:              "Zero timeout (should default to 500ms)",
			timeout:           0,
			tolerance:         1,
			expectedTimeout:   500 * time.Millisecond,
			expectedTolerance: 1,
		},
		{
			name:              "Negative timeout (should default to 500ms)",
			timeout:           -100 * time.Millisecond,
			tolerance:         1,
			expectedTimeout:   500 * time.Millisecond,
			expectedTolerance: 1,
		},
		{
			name:              "Negative tolerance (should default to 1)",
			timeout:           500 * time.Millisecond,
			tolerance:         -5,
			expectedTimeout:   500 * time.Millisecond,
			expectedTolerance: 1,
		},
		{
			name:              "Zero tolerance (should default to 1)",
			timeout:           500 * time.Millisecond,
			tolerance:         -1,
			expectedTimeout:   500 * time.Millisecond,
			expectedTolerance: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewClickDetector(tt.timeout, tt.tolerance)

			if detector.Timeout() != tt.expectedTimeout {
				t.Errorf("Expected timeout %v, got %v", tt.expectedTimeout, detector.Timeout())
			}

			if detector.Tolerance() != tt.expectedTolerance {
				t.Errorf("Expected tolerance %d, got %d", tt.expectedTolerance, detector.Tolerance())
			}
		})
	}
}

// Edge case: ClickCount after reset
func TestClickDetector_ClickCountAfterReset(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)
	pos := value2.NewPosition(10, 5)

	// Create triple click
	for i := 0; i < 3; i++ {
		event := model.NewMouseEvent(
			value2.EventRelease,
			value2.ButtonLeft,
			pos,
			value2.ModifierNone,
		)
		detector.DetectClick(event)
	}

	if detector.ClickCount() != 3 {
		t.Errorf("Expected click count 3, got %d", detector.ClickCount())
	}

	// Reset
	detector.Reset()

	if detector.ClickCount() != 0 {
		t.Errorf("Expected click count 0 after reset, got %d", detector.ClickCount())
	}

	// New click should start fresh
	event := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		pos,
		value2.ModifierNone,
	)
	click := detector.DetectClick(event)

	if click.Type() != value2.EventClick {
		t.Errorf("Expected EventClick after reset, got %v", click.Type())
	}

	if detector.ClickCount() != 1 {
		t.Errorf("Expected click count 1 after first click post-reset, got %d", detector.ClickCount())
	}
}

// Edge case: Multiple resets
func TestClickDetector_MultipleResets(t *testing.T) {
	detector := NewClickDetector(500*time.Millisecond, 1)

	// Reset without any clicks
	detector.Reset()
	if detector.ClickCount() != 0 {
		t.Errorf("Expected click count 0, got %d", detector.ClickCount())
	}

	// Create a click
	event := model.NewMouseEvent(
		value2.EventRelease,
		value2.ButtonLeft,
		value2.NewPosition(10, 5),
		value2.ModifierNone,
	)
	detector.DetectClick(event)

	// Multiple resets
	for i := 0; i < 5; i++ {
		detector.Reset()
		if detector.ClickCount() != 0 {
			t.Errorf("Reset %d: expected click count 0, got %d", i, detector.ClickCount())
		}
	}
}
