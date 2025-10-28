// Package application contains application services that coordinate domain logic.
// Application services orchestrate domain services and expose use cases.
package application

import (
	"time"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/service"
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// EventProcessor processes raw mouse events and enriches them with
// higher-level semantics (clicks, drags, etc.).
type EventProcessor struct {
	clickDetector    *service2.ClickDetector
	dragTracker      *service2.DragTracker
	scrollCalculator *service2.ScrollCalculator
}

// NewEventProcessor creates a new EventProcessor with default settings.
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		clickDetector:    service2.NewClickDetector(500*1000000, 1), // 500ms, 1 cell tolerance
		dragTracker:      service2.NewDragTracker(2),                // 2 cell threshold
		scrollCalculator: service2.NewScrollCalculator(3),           // 3 lines per scroll
	}
}

// NewEventProcessorWithConfig creates a new EventProcessor with custom configuration.
func NewEventProcessorWithConfig(
	clickTimeout time.Duration, // timeout between clicks
	clickTolerance int,
	dragThreshold int,
	linesPerScroll int,
) *EventProcessor {
	return &EventProcessor{
		clickDetector:    service2.NewClickDetector(clickTimeout, clickTolerance),
		dragTracker:      service2.NewDragTracker(dragThreshold),
		scrollCalculator: service2.NewScrollCalculator(linesPerScroll),
	}
}

// ProcessEvent processes a raw mouse event and returns enriched events.
// May return multiple events (e.g., both a release and a click event).
func (p *EventProcessor) ProcessEvent(event model.MouseEvent) []model.MouseEvent {
	var events []model.MouseEvent

	switch event.Type() {
	case value.EventPress:
		// Start drag tracking
		p.dragTracker.ProcessPress(event)
		events = append(events, event)

	case value.EventRelease:
		// Check for drag end
		wasDrag, _, _ := p.dragTracker.ProcessRelease(event)

		// Only detect clicks if it wasn't a drag
		if !wasDrag {
			if clickEvent := p.clickDetector.DetectClick(event); clickEvent != nil {
				events = append(events, *clickEvent)
			}
		}

		events = append(events, event)

	case value.EventMotion:
		// Check for drag
		if dragEvent := p.dragTracker.ProcessMotion(event); dragEvent != nil {
			events = append(events, *dragEvent)
		} else {
			// Normal motion
			events = append(events, event)
		}

	case value.EventScroll:
		// Scroll events pass through
		events = append(events, event)

	default:
		// Unknown event type, pass through
		events = append(events, event)
	}

	return events
}

// Reset resets the event processor state (useful for testing or state cleanup).
func (p *EventProcessor) Reset() {
	p.clickDetector.Reset()
	p.dragTracker.Reset()
}

// ScrollDelta calculates the scroll delta for a scroll event.
func (p *EventProcessor) ScrollDelta(event model.MouseEvent) int {
	return p.scrollCalculator.CalculateDelta(event)
}

// IsScrollUp checks if the event is a scroll up event.
func (p *EventProcessor) IsScrollUp(event model.MouseEvent) bool {
	return p.scrollCalculator.IsScrollUp(event)
}

// IsScrollDown checks if the event is a scroll down event.
func (p *EventProcessor) IsScrollDown(event model.MouseEvent) bool {
	return p.scrollCalculator.IsScrollDown(event)
}

// IsDragging returns true if a drag is currently in progress.
func (p *EventProcessor) IsDragging() bool {
	return p.dragTracker.IsDrag()
}

// ClickCount returns the current click count (for debugging/testing).
func (p *EventProcessor) ClickCount() int {
	return p.clickDetector.ClickCount()
}
