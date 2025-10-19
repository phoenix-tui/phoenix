package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/list/domain/value"
)

// FilterService handles list filtering logic
type FilterService struct{}

// NewFilterService creates a new filter service
func NewFilterService() *FilterService {
	return &FilterService{}
}

// Filter applies the filter function to all items and returns matching items
func (s *FilterService) Filter(items []*value.Item, query string,
	filterFunc func(*value.Item, string) bool) []*value.Item {

	if query == "" {
		return items
	}

	if filterFunc == nil {
		filterFunc = s.DefaultFilter
	}

	result := make([]*value.Item, 0, len(items))
	for _, item := range items {
		if filterFunc(item, query) {
			result = append(result, item)
		}
	}

	return result
}

// DefaultFilter performs case-insensitive substring matching on item labels
func (s *FilterService) DefaultFilter(item *value.Item, query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(
		strings.ToLower(item.Label()),
		strings.ToLower(query),
	)
}
