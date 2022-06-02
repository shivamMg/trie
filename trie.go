package trie

import (
	"github.com/shivamMg/ppds/tree"
)

const (
	RootKeyPart    = "^"
	TerminalSuffix = "($)"
)

type Trie struct {
	root *Node
}

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

func (n *Node) KeyPart() string {
	return n.keyPart
}

func (n *Node) IsTerminal() bool {
	return n.isTerminal
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) ChildNodes() []*Node {
	return n.childNodes()
}

func (n *Node) Data() interface{} {
	data := n.keyPart
	if n.isTerminal {
		data += " " + TerminalSuffix
	}
	return data
}

func (n *Node) Children() []tree.Node {
	children := n.childNodes()
	result := make([]tree.Node, len(children))
	for i, child := range children {
		result[i] = tree.Node(child)
	}
	return result
}

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

func New() *Trie {
	return &Trie{root: newNode(RootKeyPart)}
}

func (t *Trie) Root() *Node {
	return t.root
}

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
