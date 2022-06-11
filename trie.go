package trie

import (
	"github.com/shivamMg/ppds/tree"
)

const (
	RootKeyPart    = "^"
	terminalSuffix = "($)"
)

// Trie is the trie data structure.
type Trie struct {
	root *Node
}

// Node is a tree node inside Trie.
type Node struct {
	keyPart     string
	isTerminal  bool
	value       interface{}
	dllNode     *dllNode
	children    map[string]*Node
	childrenDLL *doublyLinkedList
}

func newNode(keyPart string) *Node {
	return &Node{
		keyPart:     keyPart,
		children:    make(map[string]*Node),
		childrenDLL: &doublyLinkedList{},
	}
}

// KeyPart returns the part (string) of the key ([]string) that this Node represents.
func (n *Node) KeyPart() string {
	return n.keyPart
}

// IsTerminal returns a boolean that tells whether a key ends at this Node.
func (n *Node) IsTerminal() bool {
	return n.isTerminal
}

// Value returns the value stored for the key ending at this Node. If Node is not a terminal, it returns nil.
func (n *Node) Value() interface{} {
	return n.value
}

// ChildNodes returns the child-nodes of this Node.
func (n *Node) ChildNodes() []*Node {
	return n.childNodes()
}

// Data is used in Print(). Use Value() to get value at this Node.
func (n *Node) Data() interface{} {
	data := n.keyPart
	if n.isTerminal {
		data += " " + terminalSuffix
	}
	return data
}

// Children is used in Print(). Use ChildNodes() to get child-nodes of this Node.
func (n *Node) Children() []tree.Node {
	children := n.childNodes()
	result := make([]tree.Node, len(children))
	for i, child := range children {
		result[i] = tree.Node(child)
	}
	return result
}

// Print prints the tree rooted at this Node. A Trie's root node is printed as RootKeyPart.
// All the terminal nodes are suffixed with ($).
func (n *Node) Print() {
	tree.PrintHrn(n)
}

func (n *Node) Sprint() string {
	return tree.SprintHrn(n)
}

func (n *Node) childNodes() []*Node {
	children := make([]*Node, 0, len(n.children))
	dllNode := n.childrenDLL.head
	for dllNode != nil {
		children = append(children, dllNode.trieNode)
		dllNode = dllNode.next
	}
	return children
}

// New returns a new instance of Trie.
func New() *Trie {
	return &Trie{root: newNode(RootKeyPart)}
}

// Root returns the root node of the Trie.
func (t *Trie) Root() *Node {
	return t.root
}

// Put upserts value the given key in the Trie. It returns a boolean depending on
// whether the key already existed or not.
func (t *Trie) Put(key []string, value interface{}) (existed bool) {
	node := t.root
	for i, part := range key {
		child, ok := node.children[part]
		if !ok {
			child = newNode(part)
			child.dllNode = newDLLNode(child)
			node.children[part] = child
			node.childrenDLL.append(child.dllNode)
		}
		if i == len(key)-1 {
			existed = child.isTerminal
			child.isTerminal = true
			child.value = value
		}
		node = child
	}
	return existed
}

// Delete deletes key-value for the given key in the Trie. It returns (value, true) if the key existed,
// else (nil, false).
func (t *Trie) Delete(key []string) (value interface{}, existed bool) {
	node := t.root
	parent := make(map[*Node]*Node)
	for _, keyPart := range key {
		child, ok := node.children[keyPart]
		if !ok {
			return nil, false
		}
		parent[child] = node
		node = child
	}
	if !node.isTerminal {
		return nil, false
	}
	node.isTerminal = false
	value = node.value
	node.value = nil
	for node != nil && !node.isTerminal && len(node.children) == 0 {
		delete(parent[node].children, node.keyPart)
		parent[node].childrenDLL.pop(node.dllNode)
		node = parent[node]
	}
	return value, true
}
