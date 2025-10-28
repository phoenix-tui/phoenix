package model

import (
	"strings"
	"testing"

	value2 "github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// TestNode_Creation tests node creation and default values
func TestNode_Creation(t *testing.T) {
	tests := []struct {
		name        string
		box         *Box
		shouldPanic bool
	}{
		{
			name:        "valid box",
			box:         NewBox("Test"),
			shouldPanic: false,
		},
		{
			name:        "nil box panics",
			box:         nil,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but got none")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			node := NewNode(tt.box)

			if tt.shouldPanic {
				return // Test passed if we panicked
			}

			// Check default values
			if node.Box() != tt.box {
				t.Error("Box not set correctly")
			}

			if len(node.Children()) != 0 {
				t.Errorf("Expected no children, got %d", len(node.Children()))
			}

			if !node.Position().IsOrigin() {
				t.Errorf("Expected origin position, got %s", node.Position())
			}

			if !node.IsLeaf() {
				t.Error("New node should be a leaf")
			}

			if node.HasChildren() {
				t.Error("New node should not have children")
			}
		})
	}
}

// TestNode_AddChild tests adding single child
func TestNode_AddChild(t *testing.T) {
	parent := NewNode(NewBox("Parent"))
	child := NewNode(NewBox("Child"))

	// Add child
	modified := parent.AddChild(child)

	// Verify immutability
	if len(parent.Children()) != 0 {
		t.Error("Original parent was mutated")
	}

	if len(modified.Children()) != 1 {
		t.Errorf("Expected 1 child, got %d", len(modified.Children()))
	}

	if modified.Children()[0] != child {
		t.Error("Child not added correctly")
	}

	if modified.IsLeaf() {
		t.Error("Node with children should not be a leaf")
	}

	if !modified.HasChildren() {
		t.Error("Node should have children")
	}

	if modified.ChildCount() != 1 {
		t.Errorf("Expected child count 1, got %d", modified.ChildCount())
	}
}

// TestNode_AddChild_Panics tests panic conditions for AddChild
func TestNode_AddChild_Panics(t *testing.T) {
	parent := NewNode(NewBox("Parent"))

	// Test nil child
	t.Run("nil child panics", func(_ *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil child")
			}
		}()
		parent.AddChild(nil)
	})

	// Test self as child (cycle detection)
	t.Run("self as child panics", func(_ *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for self as child")
			}
		}()
		parent.AddChild(parent)
	})
}

// TestNode_AddChildren tests adding multiple children at once
func TestNode_AddChildren(t *testing.T) {
	parent := NewNode(NewBox("Parent"))
	child1 := NewNode(NewBox("Child 1"))
	child2 := NewNode(NewBox("Child 2"))
	child3 := NewNode(NewBox("Child 3"))

	modified := parent.AddChildren(child1, child2, child3)

	if len(modified.Children()) != 3 {
		t.Errorf("Expected 3 children, got %d", len(modified.Children()))
	}

	children := modified.Children()
	if children[0] != child1 || children[1] != child2 || children[2] != child3 {
		t.Error("Children not added in correct order")
	}
}

// TestNode_RemoveChild tests child removal
func TestNode_RemoveChild(t *testing.T) {
	parent := NewNode(NewBox("Parent")).
		AddChildren(
			NewNode(NewBox("Child 1")),
			NewNode(NewBox("Child 2")),
			NewNode(NewBox("Child 3")),
		)

	// Remove middle child
	modified := parent.RemoveChild(1)

	// Verify immutability
	if len(parent.Children()) != 3 {
		t.Error("Original parent was mutated")
	}

	if len(modified.Children()) != 2 {
		t.Errorf("Expected 2 children after removal, got %d", len(modified.Children()))
	}

	// Verify correct children remain
	children := modified.Children()
	if children[0].Box().Content() != "Child 1" {
		t.Error("First child incorrect after removal")
	}
	if children[1].Box().Content() != "Child 3" {
		t.Error("Second child incorrect after removal")
	}

	// Remove first child
	modified2 := modified.RemoveChild(0)
	if len(modified2.Children()) != 1 {
		t.Errorf("Expected 1 child, got %d", len(modified2.Children()))
	}
	if modified2.Children()[0].Box().Content() != "Child 3" {
		t.Error("Wrong child remained")
	}

	// Remove last child
	modified3 := modified2.RemoveChild(0)
	if len(modified3.Children()) != 0 {
		t.Errorf("Expected 0 children, got %d", len(modified3.Children()))
	}
	if !modified3.IsLeaf() {
		t.Error("Should be a leaf after removing all children")
	}
}

// TestNode_RemoveChild_Panics tests panic conditions for RemoveChild
func TestNode_RemoveChild_Panics(t *testing.T) {
	parent := NewNode(NewBox("Parent")).
		AddChild(NewNode(NewBox("Child")))

	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"index too large", 1},
		{"index way too large", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic for invalid index")
				}
			}()
			parent.RemoveChild(tt.index)
		})
	}
}

// TestNode_ClearChildren tests removing all children
func TestNode_ClearChildren(t *testing.T) {
	parent := NewNode(NewBox("Parent")).
		AddChildren(
			NewNode(NewBox("Child 1")),
			NewNode(NewBox("Child 2")),
			NewNode(NewBox("Child 3")),
		)

	modified := parent.ClearChildren()

	// Verify immutability
	if len(parent.Children()) != 3 {
		t.Error("Original parent was mutated")
	}

	if len(modified.Children()) != 0 {
		t.Errorf("Expected 0 children, got %d", len(modified.Children()))
	}

	if !modified.IsLeaf() {
		t.Error("Should be a leaf after clearing children")
	}

	// Clear already empty node (should return self)
	alreadyEmpty := NewNode(NewBox("Empty"))
	cleared := alreadyEmpty.ClearChildren()
	if cleared != alreadyEmpty {
		t.Error("Clearing empty node should return self")
	}
}

// TestNode_SetPosition tests position modification
func TestNode_SetPosition(t *testing.T) {
	original := NewNode(NewBox("Test"))
	position := value2.NewPosition(10, 5)

	modified := original.SetPosition(position)

	// Verify immutability
	if !original.Position().IsOrigin() {
		t.Error("Original node was mutated")
	}

	if !modified.Position().Equals(position) {
		t.Errorf("Expected position %s, got %s", position, modified.Position())
	}
}

// TestNode_SetBox tests box modification
func TestNode_SetBox(t *testing.T) {
	originalBox := NewBox("Original")
	node := NewNode(originalBox)

	newBox := NewBox("Modified")
	modified := node.SetBox(newBox)

	// Verify immutability
	if node.Box().Content() != "Original" {
		t.Error("Original node was mutated")
	}

	if modified.Box().Content() != "Modified" {
		t.Error("Box not set correctly")
	}

	// Test nil box panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil box")
		}
	}()
	node.SetBox(nil)
}

// TestNode_Walk tests depth-first traversal
func TestNode_Walk(t *testing.T) {
	// Build tree:
	//   Root
	//   ├── Child 1
	//   │   ├── Grandchild 1.1
	//   │   └── Grandchild 1.2
	//   └── Child 2
	root := NewNode(NewBox("Root")).
		AddChildren(
			NewNode(NewBox("Child 1")).
				AddChildren(
					NewNode(NewBox("Grandchild 1.1")),
					NewNode(NewBox("Grandchild 1.2")),
				),
			NewNode(NewBox("Child 2")),
		)

	// Collect visited nodes
	var visited []string
	root.Walk(func(n *Node) {
		visited = append(visited, n.Box().Content())
	})

	expected := []string{
		"Root",
		"Child 1",
		"Grandchild 1.1",
		"Grandchild 1.2",
		"Child 2",
	}

	if len(visited) != len(expected) {
		t.Errorf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, content := range expected {
		if visited[i] != content {
			t.Errorf("At index %d: expected %q, got %q", i, content, visited[i])
		}
	}
}

// TestNode_WalkWithDepth tests depth-first traversal with depth tracking
func TestNode_WalkWithDepth(t *testing.T) {
	// Build tree:
	//   Root (depth 0)
	//   ├── Child 1 (depth 1)
	//   │   └── Grandchild 1.1 (depth 2)
	//   └── Child 2 (depth 1)
	root := NewNode(NewBox("Root")).
		AddChildren(
			NewNode(NewBox("Child 1")).
				AddChild(NewNode(NewBox("Grandchild 1.1"))),
			NewNode(NewBox("Child 2")),
		)

	// Collect nodes with depths
	type nodeInfo struct {
		content string
		depth   int
	}
	var visited []nodeInfo

	root.WalkWithDepth(func(n *Node, depth int) {
		visited = append(visited, nodeInfo{
			content: n.Box().Content(),
			depth:   depth,
		})
	})

	expected := []nodeInfo{
		{"Root", 0},
		{"Child 1", 1},
		{"Grandchild 1.1", 2},
		{"Child 2", 1},
	}

	if len(visited) != len(expected) {
		t.Errorf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, exp := range expected {
		if visited[i].content != exp.content || visited[i].depth != exp.depth {
			t.Errorf("At index %d: expected %v, got %v", i, exp, visited[i])
		}
	}
}

// TestNode_Depth tests tree depth calculation
func TestNode_Depth(t *testing.T) {
	tests := []struct {
		name          string
		buildTree     func() *Node
		expectedDepth int
	}{
		{
			name: "leaf node",
			buildTree: func() *Node {
				return NewNode(NewBox("Leaf"))
			},
			expectedDepth: 0,
		},
		{
			name: "one level",
			buildTree: func() *Node {
				return NewNode(NewBox("Root")).
					AddChild(NewNode(NewBox("Child")))
			},
			expectedDepth: 1,
		},
		{
			name: "two levels",
			buildTree: func() *Node {
				return NewNode(NewBox("Root")).
					AddChild(
						NewNode(NewBox("Child")).
							AddChild(NewNode(NewBox("Grandchild"))),
					)
			},
			expectedDepth: 2,
		},
		{
			name: "unbalanced tree",
			buildTree: func() *Node {
				return NewNode(NewBox("Root")).
					AddChildren(
						NewNode(NewBox("Child 1")).
							AddChild(
								NewNode(NewBox("Grandchild 1.1")).
									AddChild(NewNode(NewBox("Great-grandchild"))),
							),
						NewNode(NewBox("Child 2")),
					)
			},
			expectedDepth: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tree := tt.buildTree()
			depth := tree.Depth()

			if depth != tt.expectedDepth {
				t.Errorf("Expected depth %d, got %d", tt.expectedDepth, depth)
			}
		})
	}
}

// TestNode_NodeCount tests total node count calculation
func TestNode_NodeCount(t *testing.T) {
	tests := []struct {
		name          string
		buildTree     func() *Node
		expectedCount int
	}{
		{
			name: "single node",
			buildTree: func() *Node {
				return NewNode(NewBox("Single"))
			},
			expectedCount: 1,
		},
		{
			name: "parent with one child",
			buildTree: func() *Node {
				return NewNode(NewBox("Parent")).
					AddChild(NewNode(NewBox("Child")))
			},
			expectedCount: 2,
		},
		{
			name: "three generations",
			buildTree: func() *Node {
				return NewNode(NewBox("Root")).
					AddChildren(
						NewNode(NewBox("Child 1")).
							AddChildren(
								NewNode(NewBox("Grandchild 1.1")),
								NewNode(NewBox("Grandchild 1.2")),
							),
						NewNode(NewBox("Child 2")),
					)
			},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tree := tt.buildTree()
			count := tree.NodeCount()

			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}

// TestNode_String tests debug representation
func TestNode_String(t *testing.T) {
	// Build simple tree
	root := NewNode(NewBox("Root")).
		AddChildren(
			NewNode(NewBox("Child 1")),
			NewNode(NewBox("Child 2")),
		)

	str := root.String()

	// Check for expected content
	expected := []string{
		"Node{\"Root\"",
		"Node{\"Child 1\"",
		"Node{\"Child 2\"",
		"[pos=0,0]",
	}

	for _, substr := range expected {
		if !strings.Contains(str, substr) {
			t.Errorf("Expected string to contain %q, got: %s", substr, str)
		}
	}
}

// TestNode_String_LongContent tests content truncation in String
func TestNode_String_LongContent(t *testing.T) {
	node := NewNode(NewBox("This is a very long content string that should be truncated"))

	str := node.String()

	if !strings.Contains(str, "...") {
		t.Error("Expected long content to be truncated with '...'")
	}
}

// TestNode_FluentAPI tests method chaining for tree building
func TestNode_FluentAPI(t *testing.T) {
	tree := NewNode(NewBox("Root")).
		AddChildren(
			NewNode(NewBox("Child 1")).
				AddChild(NewNode(NewBox("Grandchild 1.1"))),
			NewNode(NewBox("Child 2")),
		).
		SetPosition(value2.NewPosition(10, 5))

	// Verify structure
	if tree.Position().X() != 10 || tree.Position().Y() != 5 {
		t.Error("Position not set in chain")
	}

	if len(tree.Children()) != 2 {
		t.Errorf("Expected 2 children, got %d", len(tree.Children()))
	}

	if len(tree.Children()[0].Children()) != 1 {
		t.Error("First child should have 1 child")
	}
}

// TestNode_Immutability tests that all operations return new instances
func TestNode_Immutability(t *testing.T) {
	original := NewNode(NewBox("Original")).
		AddChild(NewNode(NewBox("Child"))).
		SetPosition(value2.NewPosition(5, 5))

	// Perform modifications
	_ = original.AddChild(NewNode(NewBox("New Child")))
	_ = original.RemoveChild(0)
	_ = original.ClearChildren()
	_ = original.SetPosition(value2.NewPosition(10, 10))
	_ = original.SetBox(NewBox("Modified"))

	// Verify original is unchanged
	if original.Box().Content() != "Original" {
		t.Error("Box was mutated")
	}

	if len(original.Children()) != 1 {
		t.Error("Children were mutated")
	}

	if original.Position().X() != 5 || original.Position().Y() != 5 {
		t.Error("Position was mutated")
	}
}

// TestNode_ChildrenImmutability tests that Children() returns a copy
func TestNode_ChildrenImmutability(t *testing.T) {
	node := NewNode(NewBox("Parent")).
		AddChildren(
			NewNode(NewBox("Child 1")),
			NewNode(NewBox("Child 2")),
		)

	// Get children slice
	children := node.Children()

	// Modify the slice
	children[0] = NewNode(NewBox("Modified"))
	_ = append(children, NewNode(NewBox("Added")))

	// Verify original node is unchanged
	originalChildren := node.Children()
	if len(originalChildren) != 2 {
		t.Error("Children count was affected by external modification")
	}

	if originalChildren[0].Box().Content() != "Child 1" {
		t.Error("Child content was affected by external modification")
	}
}

// TestIntegration_BoxAndNode tests Box and Node working together
func TestIntegration_BoxAndNode(t *testing.T) {
	// Create a box with full styling
	box := NewBox("Hello, World!").
		WithPadding(value2.NewSpacingAll(1)).
		WithMargin(value2.NewSpacingAll(1)).
		WithBorder(true).
		WithAlignment(value2.NewAlignmentCenter())

	// Create node with box
	node := NewNode(box).
		AddChildren(
			NewNode(NewBox("Child 1").WithPadding(value2.NewSpacingAll(1))),
			NewNode(NewBox("Child 2").WithBorder(true)),
		).
		SetPosition(value2.NewPosition(10, 5))

	// Verify node properties
	if node.Box() != box {
		t.Error("Box not preserved in node")
	}

	if node.ChildCount() != 2 {
		t.Errorf("Expected 2 children, got %d", node.ChildCount())
	}

	// Verify box size calculations work
	totalSize := node.Box().TotalSize()
	if totalSize.Width() <= 0 || totalSize.Height() <= 0 {
		t.Error("Box size calculation failed")
	}

	// Verify tree structure
	nodeCount := node.NodeCount()
	if nodeCount != 3 {
		t.Errorf("Expected 3 total nodes, got %d", nodeCount)
	}

	// Verify position
	if node.Position().X() != 10 || node.Position().Y() != 5 {
		t.Errorf("Expected position (10,5), got %s", node.Position())
	}
}
