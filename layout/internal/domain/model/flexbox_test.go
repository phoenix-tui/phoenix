package model

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// Test helper to create a simple node
func createNode(content string) *Node {
	return NewNode(NewBox(content))
}

func TestNewFlexContainer(t *testing.T) {
	tests := []struct {
		name      string
		direction value2.FlexDirection
		wantPanic bool
	}{
		{
			name:      "Row direction",
			direction: value2.FlexDirectionRow,
			wantPanic: false,
		},
		{
			name:      "Column direction",
			direction: value2.FlexDirectionColumn,
			wantPanic: false,
		},
		{
			name:      "Invalid direction panics",
			direction: value2.FlexDirection(99),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("NewFlexContainer() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			container := NewFlexContainer(tt.direction)
			if tt.wantPanic {
				return
			}

			// Check defaults
			if container.Direction() != tt.direction {
				t.Errorf("Direction() = %v, want %v", container.Direction(), tt.direction)
			}
			if container.JustifyContent() != value2.JustifyContentStart {
				t.Errorf("JustifyContent() = %v, want %v", container.JustifyContent(), value2.JustifyContentStart)
			}
			if container.AlignItems() != value2.AlignItemsStretch {
				t.Errorf("AlignItems() = %v, want %v", container.AlignItems(), value2.AlignItemsStretch)
			}
			if container.Gap() != 0 {
				t.Errorf("Gap() = %v, want 0", container.Gap())
			}
			if !container.IsEmpty() {
				t.Errorf("IsEmpty() = false, want true")
			}
			if !container.Size().IsUnconstrained() {
				t.Errorf("Size() should be unconstrained")
			}
		})
	}
}

func TestFlexContainer_ItemCount(t *testing.T) {
	container := NewFlexContainer(value2.FlexDirectionRow)

	if container.ItemCount() != 0 {
		t.Errorf("ItemCount() = %v, want 0", container.ItemCount())
	}

	container = container.AddItem(createNode("Item 1"))
	if container.ItemCount() != 1 {
		t.Errorf("ItemCount() = %v, want 1", container.ItemCount())
	}

	container = container.AddItem(createNode("Item 2"))
	if container.ItemCount() != 2 {
		t.Errorf("ItemCount() = %v, want 2", container.ItemCount())
	}
}

func TestFlexContainer_IsEmpty(t *testing.T) {
	container := NewFlexContainer(value2.FlexDirectionRow)

	if !container.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true")
	}

	container = container.AddItem(createNode("Item"))
	if container.IsEmpty() {
		t.Errorf("IsEmpty() = true, want false")
	}
}

func TestFlexContainer_WithDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction value2.FlexDirection
		wantPanic bool
	}{
		{
			name:      "Change to row",
			direction: value2.FlexDirectionRow,
			wantPanic: false,
		},
		{
			name:      "Change to column",
			direction: value2.FlexDirectionColumn,
			wantPanic: false,
		},
		{
			name:      "Invalid direction panics",
			direction: value2.FlexDirection(99),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("WithDirection() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			container := NewFlexContainer(value2.FlexDirectionRow)
			newContainer := container.WithDirection(tt.direction)

			if tt.wantPanic {
				return
			}

			// Check immutability
			if container.Direction() != value2.FlexDirectionRow {
				t.Errorf("Original container was mutated")
			}

			if newContainer.Direction() != tt.direction {
				t.Errorf("Direction() = %v, want %v", newContainer.Direction(), tt.direction)
			}
		})
	}
}

func TestFlexContainer_WithJustifyContent(t *testing.T) {
	tests := []struct {
		name      string
		justify   value2.JustifyContent
		wantPanic bool
	}{
		{
			name:      "Set to start",
			justify:   value2.JustifyContentStart,
			wantPanic: false,
		},
		{
			name:      "Set to end",
			justify:   value2.JustifyContentEnd,
			wantPanic: false,
		},
		{
			name:      "Set to center",
			justify:   value2.JustifyContentCenter,
			wantPanic: false,
		},
		{
			name:      "Set to space-between",
			justify:   value2.JustifyContentSpaceBetween,
			wantPanic: false,
		},
		{
			name:      "Invalid justify panics",
			justify:   value2.JustifyContent(99),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("WithJustifyContent() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			container := NewFlexContainer(value2.FlexDirectionRow)
			newContainer := container.WithJustifyContent(tt.justify)

			if tt.wantPanic {
				return
			}

			if newContainer.JustifyContent() != tt.justify {
				t.Errorf("JustifyContent() = %v, want %v", newContainer.JustifyContent(), tt.justify)
			}
		})
	}
}

func TestFlexContainer_WithAlignItems(t *testing.T) {
	tests := []struct {
		name      string
		align     value2.AlignItems
		wantPanic bool
	}{
		{
			name:      "Set to stretch",
			align:     value2.AlignItemsStretch,
			wantPanic: false,
		},
		{
			name:      "Set to start",
			align:     value2.AlignItemsStart,
			wantPanic: false,
		},
		{
			name:      "Set to end",
			align:     value2.AlignItemsEnd,
			wantPanic: false,
		},
		{
			name:      "Set to center",
			align:     value2.AlignItemsCenter,
			wantPanic: false,
		},
		{
			name:      "Invalid align panics",
			align:     value2.AlignItems(99),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("WithAlignItems() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			container := NewFlexContainer(value2.FlexDirectionRow)
			newContainer := container.WithAlignItems(tt.align)

			if tt.wantPanic {
				return
			}

			if newContainer.AlignItems() != tt.align {
				t.Errorf("AlignItems() = %v, want %v", newContainer.AlignItems(), tt.align)
			}
		})
	}
}

func TestFlexContainer_WithGap(t *testing.T) {
	tests := []struct {
		name      string
		gap       int
		wantPanic bool
	}{
		{
			name:      "Zero gap",
			gap:       0,
			wantPanic: false,
		},
		{
			name:      "Positive gap",
			gap:       5,
			wantPanic: false,
		},
		{
			name:      "Negative gap panics",
			gap:       -1,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("WithGap() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			container := NewFlexContainer(value2.FlexDirectionRow)
			newContainer := container.WithGap(tt.gap)

			if tt.wantPanic {
				return
			}

			if newContainer.Gap() != tt.gap {
				t.Errorf("Gap() = %v, want %v", newContainer.Gap(), tt.gap)
			}
		})
	}
}

func TestFlexContainer_WithSize(t *testing.T) {
	container := NewFlexContainer(value2.FlexDirectionRow)
	size := value2.NewSizeExact(80, 24)

	newContainer := container.WithSize(size)

	// Check immutability
	if !container.Size().IsUnconstrained() {
		t.Errorf("Original container was mutated")
	}

	if newContainer.Size() != size {
		t.Errorf("Size() = %v, want %v", newContainer.Size(), size)
	}
}

func TestFlexContainer_AddItem(t *testing.T) {
	t.Run("Add valid item", func(_ *testing.T) {
		container := NewFlexContainer(value2.FlexDirectionRow)
		item := createNode("Item 1")

		newContainer := container.AddItem(item)

		// Check immutability
		if container.ItemCount() != 0 {
			t.Errorf("Original container was mutated")
		}

		if newContainer.ItemCount() != 1 {
			t.Errorf("ItemCount() = %v, want 1", newContainer.ItemCount())
		}

		items := newContainer.Items()
		if items[0] != item {
			t.Errorf("Item not added correctly")
		}
	})

	t.Run("Add nil item panics", func(_ *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("AddItem(nil) should panic")
			}
		}()

		container := NewFlexContainer(value2.FlexDirectionRow)
		container.AddItem(nil)
	})
}

func TestFlexContainer_AddItems(t *testing.T) {
	container := NewFlexContainer(value2.FlexDirectionRow)
	item1 := createNode("Item 1")
	item2 := createNode("Item 2")
	item3 := createNode("Item 3")

	newContainer := container.AddItems(item1, item2, item3)

	if newContainer.ItemCount() != 3 {
		t.Errorf("ItemCount() = %v, want 3", newContainer.ItemCount())
	}

	items := newContainer.Items()
	if items[0] != item1 || items[1] != item2 || items[2] != item3 {
		t.Errorf("Items not added in correct order")
	}
}

func TestFlexContainer_RemoveItem(t *testing.T) {
	t.Run("Remove valid index", func(_ *testing.T) {
		container := NewFlexContainer(value2.FlexDirectionRow).
			AddItems(
				createNode("Item 1"),
				createNode("Item 2"),
				createNode("Item 3"),
			)

		newContainer := container.RemoveItem(1) // Remove middle item

		if newContainer.ItemCount() != 2 {
			t.Errorf("ItemCount() = %v, want 2", newContainer.ItemCount())
		}

		items := newContainer.Items()
		if items[0].Box().Content() != "Item 1" || items[1].Box().Content() != "Item 3" {
			t.Errorf("Items not removed correctly")
		}
	})

	t.Run("Remove negative index panics", func(_ *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("RemoveItem(-1) should panic")
			}
		}()

		container := NewFlexContainer(value2.FlexDirectionRow).AddItem(createNode("Item"))
		container.RemoveItem(-1)
	})

	t.Run("Remove out of bounds index panics", func(_ *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("RemoveItem(out of bounds) should panic")
			}
		}()

		container := NewFlexContainer(value2.FlexDirectionRow).AddItem(createNode("Item"))
		container.RemoveItem(5)
	})
}

func TestFlexContainer_ClearItems(t *testing.T) {
	t.Run("Clear non-empty container", func(_ *testing.T) {
		container := NewFlexContainer(value2.FlexDirectionRow).
			AddItems(
				createNode("Item 1"),
				createNode("Item 2"),
			)

		newContainer := container.ClearItems()

		// Check immutability
		if container.ItemCount() != 2 {
			t.Errorf("Original container was mutated")
		}

		if !newContainer.IsEmpty() {
			t.Errorf("IsEmpty() = false, want true after clear")
		}
	})

	t.Run("Clear empty container returns self", func(_ *testing.T) {
		container := NewFlexContainer(value2.FlexDirectionRow)
		newContainer := container.ClearItems()

		if newContainer != container {
			t.Errorf("ClearItems() on empty container should return self")
		}
	})
}

func TestFlexContainer_TotalGap(t *testing.T) {
	tests := []struct {
		name      string
		itemCount int
		gap       int
		wantTotal int
	}{
		{
			name:      "No items",
			itemCount: 0,
			gap:       5,
			wantTotal: 0,
		},
		{
			name:      "One item",
			itemCount: 1,
			gap:       5,
			wantTotal: 0,
		},
		{
			name:      "Two items",
			itemCount: 2,
			gap:       5,
			wantTotal: 5,
		},
		{
			name:      "Three items",
			itemCount: 3,
			gap:       2,
			wantTotal: 4, // 2 gaps * 2 = 4
		},
		{
			name:      "Five items with gap 3",
			itemCount: 5,
			gap:       3,
			wantTotal: 12, // 4 gaps * 3 = 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			container := NewFlexContainer(value2.FlexDirectionRow).WithGap(tt.gap)

			for i := 0; i < tt.itemCount; i++ {
				container = container.AddItem(createNode("Item"))
			}

			got := container.TotalGap()
			if got != tt.wantTotal {
				t.Errorf("TotalGap() = %v, want %v", got, tt.wantTotal)
			}
		})
	}
}

func TestFlexContainer_IsHorizontal(t *testing.T) {
	row := NewFlexContainer(value2.FlexDirectionRow)
	column := NewFlexContainer(value2.FlexDirectionColumn)

	if !row.IsHorizontal() {
		t.Errorf("Row container should be horizontal")
	}

	if column.IsHorizontal() {
		t.Errorf("Column container should not be horizontal")
	}
}

func TestFlexContainer_IsVertical(t *testing.T) {
	row := NewFlexContainer(value2.FlexDirectionRow)
	column := NewFlexContainer(value2.FlexDirectionColumn)

	if row.IsVertical() {
		t.Errorf("Row container should not be vertical")
	}

	if !column.IsVertical() {
		t.Errorf("Column container should be vertical")
	}
}

func TestFlexContainer_String(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *FlexContainer
		wantParts []string
	}{
		{
			name: "Default row container",
			setup: func() *FlexContainer {
				return NewFlexContainer(value2.FlexDirectionRow)
			},
			wantParts: []string{"direction=row", "items=0"},
		},
		{
			name: "Column with justify center",
			setup: func() *FlexContainer {
				return NewFlexContainer(value2.FlexDirectionColumn).
					WithJustifyContent(value2.JustifyContentCenter)
			},
			wantParts: []string{"direction=column", "justify=center", "items=0"},
		},
		{
			name: "Row with gap and items",
			setup: func() *FlexContainer {
				return NewFlexContainer(value2.FlexDirectionRow).
					WithGap(2).
					AddItems(createNode("Item 1"), createNode("Item 2"))
			},
			wantParts: []string{"direction=row", "gap=2", "items=2"},
		},
		{
			name: "Full configuration",
			setup: func() *FlexContainer {
				return NewFlexContainer(value2.FlexDirectionColumn).
					WithJustifyContent(value2.JustifyContentSpaceBetween).
					WithAlignItems(value2.AlignItemsCenter).
					WithGap(3).
					WithSize(value2.NewSizeExact(80, 24)).
					AddItem(createNode("Item"))
			},
			wantParts: []string{"direction=column", "justify=space-between", "align=center", "gap=3", "size=", "items=1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			container := tt.setup()
			str := container.String()

			for _, part := range tt.wantParts {
				if !contains(str, part) {
					t.Errorf("String() = %q, should contain %q", str, part)
				}
			}

			if !contains(str, "FlexContainer{") {
				t.Errorf("String() should start with FlexContainer{")
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestFlexContainer_Items_Immutability(t *testing.T) {
	container := NewFlexContainer(value2.FlexDirectionRow).
		AddItems(
			createNode("Item 1"),
			createNode("Item 2"),
		)

	items := container.Items()
	originalCount := len(items)

	// Try to modify returned slice (should not affect container)
	items[0] = createNode("Modified")
	_ = append(items, createNode("Extra"))

	// Check container unchanged
	if container.ItemCount() != originalCount {
		t.Errorf("Container was mutated through Items() slice")
	}

	newItems := container.Items()
	if newItems[0].Box().Content() != "Item 1" {
		t.Errorf("Container item was mutated")
	}
}
