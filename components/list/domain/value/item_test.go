package value

import (
	"reflect"
	"testing"
)

func TestNewItem(t *testing.T) {
	value := "test-value"
	label := "Test Label"

	item := NewItem(value, label)

	if item.Value() != value {
		t.Errorf("NewItem() value = %v, want %v", item.Value(), value)
	}
	if item.Label() != label {
		t.Errorf("NewItem() label = %v, want %v", item.Label(), label)
	}
	if item.metadata == nil {
		t.Error("NewItem() metadata should be initialized")
	}
}

func TestNewItemWithMetadata(t *testing.T) {
	value := 42
	label := "Answer"
	metadata := map[string]interface{}{
		"color": "blue",
		"size":  10,
	}

	item := NewItemWithMetadata(value, label, metadata)

	if item.Value() != value {
		t.Errorf("NewItemWithMetadata() value = %v, want %v", item.Value(), value)
	}
	if item.Label() != label {
		t.Errorf("NewItemWithMetadata() label = %v, want %v", item.Label(), label)
	}

	// Check metadata
	if color, ok := item.GetMetadata("color"); !ok || color != "blue" {
		t.Errorf("NewItemWithMetadata() color = %v, want blue", color)
	}
	if size, ok := item.GetMetadata("size"); !ok || size != 10 {
		t.Errorf("NewItemWithMetadata() size = %v, want 10", size)
	}
}

func TestNewItemWithMetadata_NilMetadata(t *testing.T) {
	item := NewItemWithMetadata("value", "label", nil)

	if item.metadata == nil {
		t.Error("NewItemWithMetadata() should initialize metadata even if nil is passed")
	}
}

func TestItem_Value(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"string value", "test"},
		{"int value", 42},
		{"struct value", struct{ Name string }{"Alice"}},
		{"nil value", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewItem(tt.value, "label")
			if got := item.Value(); !reflect.DeepEqual(got, tt.value) {
				t.Errorf("Item.Value() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestItem_Label(t *testing.T) {
	tests := []struct {
		name  string
		label string
	}{
		{"simple label", "Simple"},
		{"empty label", ""},
		{"unicode label", "Hello 世界 🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewItem("value", tt.label)
			if got := item.Label(); got != tt.label {
				t.Errorf("Item.Label() = %v, want %v", got, tt.label)
			}
		})
	}
}

func TestItem_Metadata(t *testing.T) {
	metadata := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	item := NewItemWithMetadata("value", "label", metadata)

	// Get metadata
	got := item.Metadata()

	// Verify contents
	if !reflect.DeepEqual(got, metadata) {
		t.Errorf("Item.Metadata() = %v, want %v", got, metadata)
	}

	// Verify it's a copy (modifying returned map shouldn't affect original)
	got["key3"] = "value3"
	if _, ok := item.metadata["key3"]; ok {
		t.Error("Item.Metadata() should return a copy, not the original map")
	}
}

func TestItem_GetMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"color": "red",
		"size":  100,
	}
	item := NewItemWithMetadata("value", "label", metadata)

	tests := []struct {
		name      string
		key       string
		wantValue interface{}
		wantOk    bool
	}{
		{
			name:      "existing key",
			key:       "color",
			wantValue: "red",
			wantOk:    true,
		},
		{
			name:      "another existing key",
			key:       "size",
			wantValue: 100,
			wantOk:    true,
		},
		{
			name:      "non-existing key",
			key:       "nonexistent",
			wantValue: nil,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := item.GetMetadata(tt.key)
			if gotOk != tt.wantOk {
				t.Errorf("Item.GetMetadata() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Item.GetMetadata() value = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestItem_WithMetadata(t *testing.T) {
	original := NewItemWithMetadata("value", "label", map[string]interface{}{
		"key1": "value1",
	})

	// Add metadata
	updated := original.WithMetadata("key2", "value2")

	// Original should be unchanged
	if _, ok := original.GetMetadata("key2"); ok {
		t.Error("Item.WithMetadata() should not modify original item")
	}

	// Updated should have both keys
	if val, ok := updated.GetMetadata("key1"); !ok || val != "value1" {
		t.Error("Item.WithMetadata() should preserve original metadata")
	}
	if val, ok := updated.GetMetadata("key2"); !ok || val != "value2" {
		t.Error("Item.WithMetadata() should add new metadata")
	}

	// Values and labels should be the same
	if updated.Value() != original.Value() {
		t.Error("Item.WithMetadata() should preserve value")
	}
	if updated.Label() != original.Label() {
		t.Error("Item.WithMetadata() should preserve label")
	}
}

func TestItem_WithMetadata_Overwrite(t *testing.T) {
	original := NewItemWithMetadata("value", "label", map[string]interface{}{
		"key": "old-value",
	})

	updated := original.WithMetadata("key", "new-value")

	// Original should still have old value
	if val, _ := original.GetMetadata("key"); val != "old-value" {
		t.Error("Item.WithMetadata() should not modify original item")
	}

	// Updated should have new value
	if val, _ := updated.GetMetadata("key"); val != "new-value" {
		t.Error("Item.WithMetadata() should overwrite existing key")
	}
}
