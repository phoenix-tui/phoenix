package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/table/domain/value"
)

func TestNewColumn(t *testing.T) {
	col := NewColumn("id", "ID", 10)

	if col.Key() != "id" {
		t.Errorf("Key() = %v, want %v", col.Key(), "id")
	}
	if col.Title() != "ID" {
		t.Errorf("Title() = %v, want %v", col.Title(), "ID")
	}
	if col.Width() != 10 {
		t.Errorf("Width() = %v, want %v", col.Width(), 10)
	}
	if !col.Alignment().IsLeft() {
		t.Errorf("Alignment() should be left by default")
	}
	if col.IsSortable() {
		t.Errorf("IsSortable() = true, want false")
	}
	if col.Renderer() != nil {
		t.Errorf("Renderer() should be nil by default")
	}
}

func TestNewColumnWithAlignment(t *testing.T) {
	col := NewColumnWithAlignment("score", "Score", 8, value.AlignmentRight)

	if col.Key() != "score" {
		t.Errorf("Key() = %v, want %v", col.Key(), "score")
	}
	if !col.Alignment().IsRight() {
		t.Errorf("Alignment() should be right")
	}
}

func TestColumn_WithWidth(t *testing.T) {
	col := NewColumn("name", "Name", 20)
	col2 := col.WithWidth(30)

	// Original unchanged.
	if col.Width() != 20 {
		t.Errorf("Original column width should be 20, got %v", col.Width())
	}

	// New column has new width.
	if col2.Width() != 30 {
		t.Errorf("New column width = %v, want 30", col2.Width())
	}

	// Other properties preserved.
	if col2.Key() != "name" {
		t.Errorf("Key should be preserved")
	}
	if col2.Title() != "Name" {
		t.Errorf("Title should be preserved")
	}
}

func TestColumn_WithAlignment(t *testing.T) {
	col := NewColumn("value", "Value", 15)
	col2 := col.WithAlignment(value.AlignmentCenter)

	// Original unchanged.
	if !col.Alignment().IsLeft() {
		t.Errorf("Original alignment should be left")
	}

	// New column has new alignment.
	if !col2.Alignment().IsCenter() {
		t.Errorf("New alignment should be center")
	}
}

func TestColumn_WithSortable(t *testing.T) {
	col := NewColumn("date", "Date", 12)
	col2 := col.WithSortable(true)

	// Original unchanged.
	if col.IsSortable() {
		t.Errorf("Original should not be sortable")
	}

	// New column is sortable.
	if !col2.IsSortable() {
		t.Errorf("New column should be sortable")
	}
}

func TestColumn_WithRenderer(t *testing.T) {
	col := NewColumn("status", "Status", 10)
	renderer := func(v interface{}) string {
		return "custom"
	}
	col2 := col.WithRenderer(renderer)

	// Original unchanged.
	if col.Renderer() != nil {
		t.Errorf("Original should have no renderer")
	}

	// New column has renderer.
	if col2.Renderer() == nil {
		t.Errorf("New column should have renderer")
	}

	// Test renderer works.
	if got := col2.Renderer()("test"); got != "custom" {
		t.Errorf("Renderer() = %v, want custom", got)
	}
}

func TestColumn_Immutability(t *testing.T) {
	// Create original column.
	col1 := NewColumn("id", "ID", 10)

	// Chain multiple operations.
	col2 := col1.
		WithWidth(20).
		WithAlignment(value.AlignmentRight).
		WithSortable(true).
		WithRenderer(func(v interface{}) string { return "test" })

	// Verify original is unchanged.
	if col1.Width() != 10 {
		t.Errorf("Original width changed")
	}
	if !col1.Alignment().IsLeft() {
		t.Errorf("Original alignment changed")
	}
	if col1.IsSortable() {
		t.Errorf("Original sortable changed")
	}
	if col1.Renderer() != nil {
		t.Errorf("Original renderer changed")
	}

	// Verify new column has all changes.
	if col2.Width() != 20 {
		t.Errorf("New width = %v, want 20", col2.Width())
	}
	if !col2.Alignment().IsRight() {
		t.Errorf("New alignment should be right")
	}
	if !col2.IsSortable() {
		t.Errorf("New column should be sortable")
	}
	if col2.Renderer() == nil {
		t.Errorf("New column should have renderer")
	}
}
