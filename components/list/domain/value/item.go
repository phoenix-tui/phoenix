package value

// Item represents a single list item with associated data
type Item struct {
	value    interface{}            // Any data associated with the item
	label    string                 // Display text (used by default rendering)
	metadata map[string]interface{} // Additional data for custom rendering
}

// NewItem creates a new list item with the given value and label
func NewItem(value interface{}, label string) *Item {
	return &Item{
		value:    value,
		label:    label,
		metadata: make(map[string]interface{}),
	}
}

// NewItemWithMetadata creates a new list item with value, label, and metadata
func NewItemWithMetadata(value interface{}, label string, metadata map[string]interface{}) *Item {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	return &Item{
		value:    value,
		label:    label,
		metadata: metadata,
	}
}

// Value returns the underlying value of the item
func (i *Item) Value() interface{} {
	return i.value
}

// Label returns the display label of the item
func (i *Item) Label() string {
	return i.label
}

// Metadata returns all metadata associated with the item
func (i *Item) Metadata() map[string]interface{} {
	// Return a copy to prevent external modification
	result := make(map[string]interface{}, len(i.metadata))
	for k, v := range i.metadata {
		result[k] = v
	}
	return result
}

// GetMetadata retrieves a specific metadata value by key
func (i *Item) GetMetadata(key string) (interface{}, bool) {
	val, ok := i.metadata[key]
	return val, ok
}

// WithMetadata returns a new Item with the given metadata key-value pair added
func (i *Item) WithMetadata(key string, value interface{}) *Item {
	newMetadata := make(map[string]interface{}, len(i.metadata)+1)
	for k, v := range i.metadata {
		newMetadata[k] = v
	}
	newMetadata[key] = value

	return &Item{
		value:    i.value,
		label:    i.label,
		metadata: newMetadata,
	}
}
