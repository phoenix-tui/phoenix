package model

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// Node represents a node in the layout tree using the Composite pattern.
//
// Tree Structure:
//
//	Root Node
//	├── Child Node 1
//	│   ├── Grandchild 1.1
//	│   └── Grandchild 1.2
//	└── Child Node 2
//	    └── Grandchild 2.1
//
// Design Philosophy:
//   - Immutable composite pattern
//   - Delegates box properties to Box aggregate
//   - Position is set by layout engine (not user)
//   - Supports arbitrary tree depth
//   - Fluent API for tree building
//
// Layout Phases:
//  1. Tree Construction: Build tree structure (AddChild)
//  2. Size Calculation: Calculate box sizes (via Box)
//  3. Position Assignment: Layout engine sets positions (SetPosition)
//  4. Rendering: Renderer traverses tree and draws
//
// Example:
//
//	root := NewNode(NewBox("Root")).
//		AddChild(NewNode(NewBox("Child 1"))).
//		AddChild(NewNode(NewBox("Child 2")))
//
//	// Later, layout engine sets positions
//	positioned := root.SetPosition(value.NewPosition(10, 5))
type Node struct {
	box      *Box           // Box properties (aggregate root)
	children []*Node        // Child nodes (composite pattern)
	position value.Position // Position after layout (set by engine)
}

// NewNode creates a Node with the given box.
// Box cannot be nil (panics on nil box).
//
// Default values:
//   - Children: Empty slice
//   - Position: Origin (0, 0)
//
// Example:
//
//	node := NewNode(NewBox("Content"))
func NewNode(box *Box) *Node {
	if box == nil {
		panic("node: box cannot be nil")
	}
	return &Node{
		box:      box,
		children: []*Node{},
		position: value.Origin(),
	}
}

// Box returns the box properties (delegate to Box aggregate).
func (n *Node) Box() *Box {
	return n.box
}

// Children returns a copy of the children slice (immutable).
// Modifications to the returned slice do not affect the node.
func (n *Node) Children() []*Node {
	// Return a copy to maintain immutability
	result := make([]*Node, len(n.children))
	copy(result, n.children)
	return result
}

// Position returns the node position (set by layout engine).
func (n *Node) Position() value.Position {
	return n.position
}

// ChildCount returns the number of direct children.
func (n *Node) ChildCount() int {
	return len(n.children)
}

// IsLeaf returns true if node has no children.
func (n *Node) IsLeaf() bool {
	return len(n.children) == 0
}

// HasChildren returns true if node has one or more children.
func (n *Node) HasChildren() bool {
	return len(n.children) > 0
}

// AddChild returns a new Node with the given child appended.
// Panics if child is nil or if child is the node itself (cycle detection).
//
// Example:
//
//	parent := NewNode(NewBox("Parent"))
//	child := NewNode(NewBox("Child"))
//	parent = parent.AddChild(child)
func (n *Node) AddChild(child *Node) *Node {
	if child == nil {
		panic("node: child cannot be nil")
	}
	if child == n {
		panic("node: cannot add self as child (cycle detected)")
	}

	result := *n
	result.children = make([]*Node, len(n.children)+1)
	copy(result.children, n.children)
	result.children[len(result.children)-1] = child

	return &result
}

// AddChildren returns a new Node with multiple children appended.
// This is a convenience method for adding multiple children at once.
//
// Example:
//
//	parent := NewNode(NewBox("Parent")).
//		AddChildren(
//			NewNode(NewBox("Child 1")),
//			NewNode(NewBox("Child 2")),
//			NewNode(NewBox("Child 3")),
//		)
func (n *Node) AddChildren(children ...*Node) *Node {
	result := n
	for _, child := range children {
		result = result.AddChild(child)
	}
	return result
}

// RemoveChild returns a new Node with the child at the given index removed.
// Panics if index is out of bounds.
//
// Example:
//
//	parent := parent.RemoveChild(0) // Remove first child
func (n *Node) RemoveChild(index int) *Node {
	if index < 0 || index >= len(n.children) {
		panic(fmt.Sprintf("node: index %d out of bounds (0-%d)", index, len(n.children)-1))
	}

	result := *n
	result.children = make([]*Node, len(n.children)-1)

	// Copy children before removed index
	for i := 0; i < index; i++ {
		result.children[i] = n.children[i]
	}

	// Copy children after removed index
	for i := index + 1; i < len(n.children); i++ {
		result.children[i-1] = n.children[i]
	}

	return &result
}

// ClearChildren returns a new Node with all children removed.
//
// Example:
//
//	parent := parent.ClearChildren()
func (n *Node) ClearChildren() *Node {
	if len(n.children) == 0 {
		return n // Already empty, return self
	}

	result := *n
	result.children = []*Node{}
	return &result
}

// SetPosition returns a new Node with the given position.
// This is typically called by the layout engine during layout pass.
//
// Example:
//
//	positioned := node.SetPosition(value.NewPosition(10, 5))
func (n *Node) SetPosition(p value.Position) *Node {
	result := *n
	result.position = p
	return &result
}

// SetBox returns a new Node with the given box.
// Panics if box is nil.
//
// Example:
//
//	updated := node.SetBox(NewBox("Updated content"))
func (n *Node) SetBox(box *Box) *Node {
	if box == nil {
		panic("node: box cannot be nil")
	}

	result := *n
	result.box = box
	return &result
}

// Walk performs a depth-first traversal of the tree.
// The function fn is called for each node, starting with this node.
//
// Example:
//
//	node.Walk(func(n *Node) {
//		fmt.Println(n.Box().Content())
//	})
func (n *Node) Walk(fn func(*Node)) {
	fn(n)
	for _, child := range n.children {
		child.Walk(fn)
	}
}

// WalkWithDepth performs a depth-first traversal with depth tracking.
// The function fn is called with each node and its depth (0-based).
//
// Example:
//
//	node.WalkWithDepth(func(n *Node, depth int) {
//		indent := strings.Repeat("  ", depth)
//		fmt.Printf("%s%s\n", indent, n.Box().Content())
//	})
func (n *Node) WalkWithDepth(fn func(*Node, int)) {
	n.walkWithDepth(fn, 0)
}

func (n *Node) walkWithDepth(fn func(*Node, int), depth int) {
	fn(n, depth)
	for _, child := range n.children {
		child.walkWithDepth(fn, depth+1)
	}
}

// Depth calculates the maximum depth of the tree.
// A leaf node has depth 0.
//
// Example:
//
//	depth := root.Depth() // Returns 2 for a tree with root -> child -> grandchild
func (n *Node) Depth() int {
	if n.IsLeaf() {
		return 0
	}

	maxChildDepth := 0
	for _, child := range n.children {
		childDepth := child.Depth()
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	return maxChildDepth + 1
}

// NodeCount returns the total number of nodes in the tree (including this node).
//
// Example:
//
//	count := root.NodeCount() // Total nodes in tree
func (n *Node) NodeCount() int {
	count := 1 // Count self
	for _, child := range n.children {
		count += child.NodeCount()
	}
	return count
}

// String returns a human-readable debug representation.
// Shows the tree structure with indentation.
//
// Example output:
//
//	Node{Root [pos=0,0]
//	  Node{Child 1 [pos=0,0]
//	    Node{Grandchild 1.1 [pos=0,0]}
//	  }
//	  Node{Child 2 [pos=0,0]}
//	}
func (n *Node) String() string {
	var sb strings.Builder
	n.buildString(&sb, 0)
	return sb.String()
}

func (n *Node) buildString(sb *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)

	// Get content preview from box
	content := n.box.Content()
	if len(content) > 20 {
		content = content[:17] + "..."
	}
	// Replace newlines for readability
	content = strings.ReplaceAll(content, "\n", "\\n")

	// Write node info
	sb.WriteString(indent)
	fmt.Fprintf(sb, "Node{%q [pos=%d,%d]", content, n.position.X(), n.position.Y())

	if len(n.children) == 0 {
		sb.WriteString("}")
		return
	}

	// Write children
	sb.WriteString("\n")
	for i, child := range n.children {
		child.buildString(sb, depth+1)
		if i < len(n.children)-1 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
	sb.WriteString(indent)
	sb.WriteString("}")
}
