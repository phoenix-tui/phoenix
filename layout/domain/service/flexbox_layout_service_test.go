package service

import (
	"testing"

	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/layout/domain/model"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

func TestNewFlexboxLayoutService(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)

	service := NewFlexboxLayoutService(measureService)

	if service == nil {
		t.Fatal("NewFlexboxLayoutService() returned nil")
	}

	if service.measureService == nil {
		t.Error("measureService is nil")
	}
}

func TestFlexboxLayoutService_Layout_EmptyContainer(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow)

	result := service.Layout(container, 80, 24)

	if !result.IsEmpty() {
		t.Error("Empty container should remain empty after layout")
	}
}

func TestFlexboxLayoutService_Layout_HorizontalStart(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	// Create container with 3 items, no gap
	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentStart).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("AA")),   // Width: 2
			model.NewNode(model.NewBox("BBB")),  // Width: 3
			model.NewNode(model.NewBox("CCCC")), // Width: 4
		)

	result := service.Layout(container, 80, 24)

	items := result.Items()
	if len(items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(items))
	}

	// Check positions (should be packed at start)
	// Item 0: x=0
	// Item 1: x=2
	// Item 2: x=5
	expectedPositions := []struct{ x, y int }{
		{0, 0},
		{2, 0},
		{5, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_HorizontalEnd(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	// Create container with 2 items
	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentEnd).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("AA")),  // Width: 2
			model.NewNode(model.NewBox("BBB")), // Width: 3
		)

	// Total item width: 2 + 3 = 5
	// Container width: 20
	// Remaining space: 20 - 5 = 15
	// Items should start at x=15
	result := service.Layout(container, 20, 10)

	items := result.Items()

	// Item 0: x=15
	// Item 1: x=17
	expectedPositions := []struct{ x, y int }{
		{15, 0},
		{17, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_HorizontalCenter(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentCenter).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("AA")),  // Width: 2
			model.NewNode(model.NewBox("BBB")), // Width: 3
		)

	// Total item width: 2 + 3 = 5
	// Container width: 20
	// Remaining space: 20 - 5 = 15
	// Start offset: 15 / 2 = 7
	result := service.Layout(container, 20, 10)

	items := result.Items()

	// Item 0: x=7
	// Item 1: x=9
	expectedPositions := []struct{ x, y int }{
		{7, 0},
		{9, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_HorizontalSpaceBetween(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentSpaceBetween).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("AA")), // Width: 2
			model.NewNode(model.NewBox("BB")), // Width: 2
			model.NewNode(model.NewBox("CC")), // Width: 2
		)

	// Total item width: 2 + 2 + 2 = 6
	// Container width: 20
	// Remaining space: 20 - 6 = 14
	// Gap between items (2 gaps): 14 / 2 = 7
	result := service.Layout(container, 20, 10)

	items := result.Items()

	// Item 0: x=0
	// Item 1: x=2 + 7 = 9
	// Item 2: x=9 + 2 + 7 = 18
	expectedPositions := []struct{ x, y int }{
		{0, 0},
		{9, 0},
		{18, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_HorizontalSpaceBetween_SingleItem(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentSpaceBetween).
		AddItem(model.NewNode(model.NewBox("AA")))

	result := service.Layout(container, 20, 10)

	items := result.Items()

	// Single item should be at start
	pos := items[0].Position()
	if pos.X() != 0 || pos.Y() != 0 {
		t.Errorf("Single item position = (%d, %d), want (0, 0)", pos.X(), pos.Y())
	}
}

func TestFlexboxLayoutService_Layout_WithGap(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentStart).
		WithGap(3).
		AddItems(
			model.NewNode(model.NewBox("AA")), // Width: 2
			model.NewNode(model.NewBox("BB")), // Width: 2
			model.NewNode(model.NewBox("CC")), // Width: 2
		)

	result := service.Layout(container, 80, 24)

	items := result.Items()

	// Item 0: x=0
	// Item 1: x=2 + 3 = 5
	// Item 2: x=5 + 2 + 3 = 10
	expectedPositions := []struct{ x, y int }{
		{0, 0},
		{5, 0},
		{10, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_VerticalStart(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionColumn).
		WithJustifyContent(value.JustifyContentStart).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("A")), // Height: 1
			model.NewNode(model.NewBox("B")), // Height: 1
			model.NewNode(model.NewBox("C")), // Height: 1
		)

	result := service.Layout(container, 80, 24)

	items := result.Items()

	// Items stacked vertically at y=0, 1, 2
	expectedPositions := []struct{ x, y int }{
		{0, 0},
		{0, 1},
		{0, 2},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_VerticalCenter(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionColumn).
		WithJustifyContent(value.JustifyContentCenter).
		WithAlignItems(value.AlignItemsStart).
		AddItems(
			model.NewNode(model.NewBox("A")), // Height: 1
			model.NewNode(model.NewBox("B")), // Height: 1
		)

	// Total height: 1 + 1 = 2
	// Container height: 10
	// Remaining space: 10 - 2 = 8
	// Start offset: 8 / 2 = 4
	result := service.Layout(container, 80, 10)

	items := result.Items()

	// Item 0: y=4
	// Item 1: y=5
	expectedPositions := []struct{ x, y int }{
		{0, 4},
		{0, 5},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}

func TestFlexboxLayoutService_Layout_CrossAxisAlignment(t *testing.T) {
	tests := []struct {
		name          string
		alignItems    value.AlignItems
		containerSize int
		itemSize      int
		expectedY     int // For horizontal layout
	}{
		{
			name:          "Align start",
			alignItems:    value.AlignItemsStart,
			containerSize: 10,
			itemSize:      3,
			expectedY:     0,
		},
		{
			name:          "Align end",
			alignItems:    value.AlignItemsEnd,
			containerSize: 10,
			itemSize:      3,
			expectedY:     7, // 10 - 3 = 7
		},
		{
			name:          "Align center",
			alignItems:    value.AlignItemsCenter,
			containerSize: 10,
			itemSize:      3,
			expectedY:     3, // (10 - 3) / 2 = 3
		},
		{
			name:          "Align stretch",
			alignItems:    value.AlignItemsStretch,
			containerSize: 10,
			itemSize:      3,
			expectedY:     0, // Positioned at 0 (stretching handled elsewhere)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unicodeService := coreService.NewUnicodeService()
			measureService := NewMeasureService(unicodeService)
			service := NewFlexboxLayoutService(measureService)

			// Create multi-line content to control item height
			content := "A\nB\nC" // Height: 3

			container := model.NewFlexContainer(value.FlexDirectionRow).
				WithAlignItems(tt.alignItems).
				AddItem(model.NewNode(model.NewBox(content)))

			result := service.Layout(container, 80, tt.containerSize)

			items := result.Items()
			pos := items[0].Position()

			if pos.Y() != tt.expectedY {
				t.Errorf("Y position = %d, want %d", pos.Y(), tt.expectedY)
			}
		})
	}
}

func TestFlexboxLayoutService_LayoutWithDetails(t *testing.T) {
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		AddItems(
			model.NewNode(model.NewBox("AA")),
			model.NewNode(model.NewBox("BBB")),
		)

	result := service.LayoutWithDetails(container, 80, 24)

	if result.Container == nil {
		t.Fatal("Container is nil")
	}

	if len(result.ItemPositions) != 2 {
		t.Errorf("ItemPositions length = %d, want 2", len(result.ItemPositions))
	}

	if len(result.ItemSizes) != 2 {
		t.Errorf("ItemSizes length = %d, want 2", len(result.ItemSizes))
	}

	// Check that positions match container items
	items := result.Container.Items()
	for i := range items {
		if items[i].Position() != result.ItemPositions[i] {
			t.Errorf("Item %d position mismatch", i)
		}
	}
}

func TestFlexboxLayoutService_Layout_Overflow(t *testing.T) {
	// Test when items exceed container size
	unicodeService := coreService.NewUnicodeService()
	measureService := NewMeasureService(unicodeService)
	service := NewFlexboxLayoutService(measureService)

	container := model.NewFlexContainer(value.FlexDirectionRow).
		WithJustifyContent(value.JustifyContentStart).
		AddItems(
			model.NewNode(model.NewBox("AAAAAAAAAA")), // Width: 10
			model.NewNode(model.NewBox("BBBBBBBBBB")), // Width: 10
		)

	// Container width: 5 (smaller than items)
	result := service.Layout(container, 5, 10)

	items := result.Items()

	// Should still position items (they'll overflow)
	// Item 0: x=0
	// Item 1: x=10
	expectedPositions := []struct{ x, y int }{
		{0, 0},
		{10, 0},
	}

	for i, expected := range expectedPositions {
		pos := items[i].Position()
		if pos.X() != expected.x || pos.Y() != expected.y {
			t.Errorf("Item %d: position = (%d, %d), want (%d, %d)",
				i, pos.X(), pos.Y(), expected.x, expected.y)
		}
	}
}
