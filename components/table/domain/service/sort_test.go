package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/table/domain/model"
	"github.com/phoenix-tui/phoenix/components/table/domain/value"
)

func TestSortService_Sort_String(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"name": "Charlie"},
		{"name": "Alice"},
		{"name": "Bob"},
	}

	// Ascending.
	sorted := svc.Sort(rows, "name", value.SortDirectionAsc)
	if sorted[0]["name"] != "Alice" {
		t.Errorf("First row = %v, want Alice", sorted[0]["name"])
	}
	if sorted[1]["name"] != "Bob" {
		t.Errorf("Second row = %v, want Bob", sorted[1]["name"])
	}
	if sorted[2]["name"] != "Charlie" {
		t.Errorf("Third row = %v, want Charlie", sorted[2]["name"])
	}

	// Descending.
	sorted = svc.Sort(rows, "name", value.SortDirectionDesc)
	if sorted[0]["name"] != "Charlie" {
		t.Errorf("First row = %v, want Charlie", sorted[0]["name"])
	}
	if sorted[2]["name"] != "Alice" {
		t.Errorf("Third row = %v, want Alice", sorted[2]["name"])
	}
}

func TestSortService_Sort_Int(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"age": 30},
		{"age": 25},
		{"age": 35},
	}

	// Ascending.
	sorted := svc.Sort(rows, "age", value.SortDirectionAsc)
	if sorted[0]["age"] != 25 {
		t.Errorf("First row age = %v, want 25", sorted[0]["age"])
	}
	if sorted[2]["age"] != 35 {
		t.Errorf("Third row age = %v, want 35", sorted[2]["age"])
	}

	// Descending.
	sorted = svc.Sort(rows, "age", value.SortDirectionDesc)
	if sorted[0]["age"] != 35 {
		t.Errorf("First row age = %v, want 35", sorted[0]["age"])
	}
	if sorted[2]["age"] != 25 {
		t.Errorf("Third row age = %v, want 25", sorted[2]["age"])
	}
}

func TestSortService_Sort_Float64(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"score": 95.5},
		{"score": 87.2},
		{"score": 92.8},
	}

	sorted := svc.Sort(rows, "score", value.SortDirectionAsc)
	if sorted[0]["score"] != 87.2 {
		t.Errorf("First row score = %v, want 87.2", sorted[0]["score"])
	}
}

func TestSortService_Sort_Bool(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"active": true},
		{"active": false},
		{"active": true},
	}

	sorted := svc.Sort(rows, "active", value.SortDirectionAsc)
	if sorted[0]["active"] != false {
		t.Errorf("First row should be false")
	}
}

func TestSortService_Sort_None(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"id": 3},
		{"id": 1},
		{"id": 2},
	}

	sorted := svc.Sort(rows, "id", value.SortDirectionNone)

	// Should return unchanged.
	if sorted[0]["id"] != 3 {
		t.Errorf("Rows should be unchanged with SortDirectionNone")
	}
}

func TestSortService_Sort_Empty(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{}
	sorted := svc.Sort(rows, "id", value.SortDirectionAsc)

	if len(sorted) != 0 {
		t.Errorf("Sorted should be empty")
	}
}

func TestSortService_Sort_Immutability(t *testing.T) {
	svc := NewSortService()

	rows := []model.Row{
		{"id": 3},
		{"id": 1},
		{"id": 2},
	}

	sorted := svc.Sort(rows, "id", value.SortDirectionAsc)

	// Original unchanged.
	if rows[0]["id"] != 3 {
		t.Errorf("Original rows should be unchanged")
	}

	// Sorted is different.
	if sorted[0]["id"] != 1 {
		t.Errorf("Sorted first row = %v, want 1", sorted[0]["id"])
	}
}

func TestSortService_Compare_String(t *testing.T) {
	svc := NewSortService()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"Equal", "abc", "abc", 0},
		{"Less", "abc", "def", -1},
		{"Greater", "def", "abc", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortService_Compare_Int(t *testing.T) {
	svc := NewSortService()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"Equal", 5, 5, 0},
		{"Less", 3, 7, -1},
		{"Greater", 10, 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortService_Compare_Float64(t *testing.T) {
	svc := NewSortService()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"Equal", 5.5, 5.5, 0},
		{"Less", 3.2, 7.8, -1},
		{"Greater", 10.1, 2.9, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortService_Compare_Bool(t *testing.T) {
	svc := NewSortService()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"BothTrue", true, true, 0},
		{"BothFalse", false, false, 0},
		{"FalseLessThanTrue", false, true, -1},
		{"TrueGreaterThanFalse", true, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortService_Compare_Nil(t *testing.T) {
	svc := NewSortService()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int
	}{
		{"BothNil", nil, nil, 0},
		{"NilLess", nil, "abc", -1},
		{"NilLess2", nil, 123, -1},
		{"NonNilGreater", "abc", nil, 1},
		{"NonNilGreater2", 123, nil, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortService_Compare_MixedTypes(t *testing.T) {
	svc := NewSortService()

	// Mixed types should fall back to string comparison.
	result := svc.Compare(123, "456")
	if result == 0 {
		t.Errorf("Compare should handle mixed types")
	}
}

func TestSortService_Sort_StableSort(t *testing.T) {
	svc := NewSortService()

	// Rows with same sort key should preserve order.
	rows := []model.Row{
		{"group": "A", "id": 1},
		{"group": "A", "id": 2},
		{"group": "B", "id": 3},
		{"group": "A", "id": 4},
	}

	sorted := svc.Sort(rows, "group", value.SortDirectionAsc)

	// All "A" should come first, in original order.
	if sorted[0]["id"] != 1 {
		t.Errorf("First A should have id=1")
	}
	if sorted[1]["id"] != 2 {
		t.Errorf("Second A should have id=2")
	}
	if sorted[2]["id"] != 4 {
		t.Errorf("Third A should have id=4")
	}
	if sorted[3]["id"] != 3 {
		t.Errorf("B should have id=3")
	}
}
