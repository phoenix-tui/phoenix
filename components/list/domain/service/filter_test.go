package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/list/domain/value"
)

func TestFilterService_Filter(t *testing.T) {
	svc := NewFilterService()

	items := []*value.Item{
		value.NewItem("file1.go", "file1.go"),
		value.NewItem("file2.txt", "file2.txt"),
		value.NewItem("main.go", "main.go"),
		value.NewItem("readme.md", "readme.md"),
	}

	tests := []struct {
		name       string
		items      []*value.Item
		query      string
		filterFunc func(*value.Item, string) bool
		wantCount  int
		wantLabels []string
	}{
		{
			name:       "empty query returns all",
			items:      items,
			query:      "",
			filterFunc: nil,
			wantCount:  4,
			wantLabels: []string{"file1.go", "file2.txt", "main.go", "readme.md"},
		},
		{
			name:       "default filter - substring match",
			items:      items,
			query:      ".go",
			filterFunc: nil,
			wantCount:  2,
			wantLabels: []string{"file1.go", "main.go"},
		},
		{
			name:       "default filter - case insensitive",
			items:      items,
			query:      "FILE",
			filterFunc: nil,
			wantCount:  2,
			wantLabels: []string{"file1.go", "file2.txt"},
		},
		{
			name:  "custom filter - exact match",
			items: items,
			query: "main.go",
			filterFunc: func(item *value.Item, query string) bool {
				return item.Label() == query
			},
			wantCount:  1,
			wantLabels: []string{"main.go"},
		},
		{
			name:  "custom filter - prefix match",
			items: items,
			query: "file",
			filterFunc: func(item *value.Item, query string) bool {
				return strings.HasPrefix(item.Label(), query)
			},
			wantCount:  2,
			wantLabels: []string{"file1.go", "file2.txt"},
		},
		{
			name:       "no matches",
			items:      items,
			query:      "xyz",
			filterFunc: nil,
			wantCount:  0,
			wantLabels: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Filter(tt.items, tt.query, tt.filterFunc)

			if len(got) != tt.wantCount {
				t.Errorf("FilterService.Filter() returned %d items, want %d", len(got), tt.wantCount)
			}

			gotLabels := make([]string, len(got))
			for i, item := range got {
				gotLabels[i] = item.Label()
			}

			if !reflect.DeepEqual(gotLabels, tt.wantLabels) {
				t.Errorf("FilterService.Filter() labels = %v, want %v", gotLabels, tt.wantLabels)
			}
		})
	}
}

func TestFilterService_DefaultFilter(t *testing.T) {
	svc := NewFilterService()

	tests := []struct {
		name  string
		item  *value.Item
		query string
		want  bool
	}{
		{
			name:  "empty query matches all",
			item:  value.NewItem("value", "Test Item"),
			query: "",
			want:  true,
		},
		{
			name:  "exact match",
			item:  value.NewItem("value", "Test Item"),
			query: "Test Item",
			want:  true,
		},
		{
			name:  "substring match",
			item:  value.NewItem("value", "Test Item"),
			query: "Test",
			want:  true,
		},
		{
			name:  "case insensitive",
			item:  value.NewItem("value", "Test Item"),
			query: "test",
			want:  true,
		},
		{
			name:  "no match",
			item:  value.NewItem("value", "Test Item"),
			query: "xyz",
			want:  false,
		},
		{
			name:  "partial match in middle",
			item:  value.NewItem("value", "Hello World"),
			query: "lo Wo",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.DefaultFilter(tt.item, tt.query); got != tt.want {
				t.Errorf("FilterService.DefaultFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterService_Filter_EmptyItems(t *testing.T) {
	svc := NewFilterService()

	items := []*value.Item{}
	got := svc.Filter(items, "query", nil)

	if len(got) != 0 {
		t.Errorf("FilterService.Filter() on empty items should return empty slice, got %d items", len(got))
	}
}

func TestFilterService_Filter_PreservesOrder(t *testing.T) {
	svc := NewFilterService()

	items := []*value.Item{
		value.NewItem(1, "apple"),
		value.NewItem(2, "apricot"),
		value.NewItem(3, "banana"),
		value.NewItem(4, "avocado"),
	}

	got := svc.Filter(items, "a", nil)

	expectedOrder := []string{"apple", "apricot", "banana", "avocado"}
	for i, item := range got {
		if item.Label() != expectedOrder[i] {
			t.Errorf("FilterService.Filter() order mismatch at index %d: got %s, want %s",
				i, item.Label(), expectedOrder[i])
		}
	}
}
